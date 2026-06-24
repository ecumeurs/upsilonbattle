package rules

import (
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)


type localAttackCtx struct {
	*gamestate.GameState
	log *logrus.Entry
}

// Attack handles the execution of a basic physical attack between two entities.
// It calculates damage based on attacker stats and weapon damage, applies defense/armor,
// handles backstabbing bonuses, and updates the action state and version of the game.
// Intent: Execute standard combat interactions between entities.
// Inputs: msg (*message.Message) - the incoming request message, req (rulermethods.ControllerAttack) - the attack parameters.
// Outputs: reply (*message.Message) - the attack results including damage dealt and credit awards.
// @spec-link [[mech_action_economy]]
func Attack(gs *gamestate.GameState, msg *message.Message, req rulermethods.ControllerAttack) (reply *message.Message) {
	ctx := localAttackCtx{
		GameState: gs,
		log: gs.Logger.WithFields(logrus.Fields{
			"RequestID":    msg.RequestId.String()[0:8],
			"controllerID": req.ControllerID.String()[0:8],
			"entityID":     req.EntityID.String()[0:8],
			"target":       req.Target,
			"rule":         "attack",
		}),
	}
	ctx.log.Debug("Controller attack request")

	ok, reply := ctx.preAttackChecks(msg, req)
	if !ok {
		return reply
	}

	// Attack handler.

	target, _ := ctx.Grid.CellAt(req.Target)

	ent := gs.Entities[req.EntityID]
	attackerAttack := ent.GetPropertyI(property.Attack)
	weaponDmg := ent.GetPropertyI(property.WeaponBaseDamage)
	totalAttack := attackerAttack.I() + weaponDmg.I()

	// In multi-entity cells, attack only targets Characters/Monsters, not WalkThrough entities.
	// @spec-link [[mechanic_multi_entity_cell_system]]
	foe, found := gs.FindCharacterInCell(target.EntityIDs)
	if !found {
		return msg.ReplyWithError("No attackable entity at target position", "entity.attack.noentity")
	}
	foeDefense := foe.GetPropertyI(property.Defense)
	foeArmor := foe.GetPropertyI(property.ArmorRating)
	foeHP := foe.GetPropertyI(property.HP)

	// @spec-link [[mech_combat_attack_computation]]
	// @spec-link [[mechanic_backstab_detection_algorithm]]
	multiplier := 1.0
	effectiveDefense := foeDefense.I() + foeArmor.I()

	if ent.IsBackstabbing(foe) {
		multiplier = 1.5
		// 50% armor penetration
		effectiveDefense = int(float64(effectiveDefense) * 0.5)
		ctx.log.Debug("Backstab detected! 150% damage and 50% armor penetration applied.")
	}

	computedDamage := tools.Max(1, int(float64(totalAttack)*multiplier)-effectiveDefense)

	// Apply shield (not penetrated)
	foeShield := foe.GetPropertyC(property.Shield)
	if foeShield.GetValue() > 0 {
		if foeShield.GetValue() >= computedDamage {
			ctx.log.Debugf("Shield absorbed all %d damage", computedDamage)
			foeShield.SetValue(foeShield.GetValue() - computedDamage)
			computedDamage = 0
		} else {
			ctx.log.Debugf("Shield absorbed %d of %d damage", foeShield.GetValue(), computedDamage)
			computedDamage -= foeShield.GetValue()
			foeShield.SetValue(0)
		}
		foe.UpdateProperty(foeShield)
	}

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + 100

	ctx.log.WithFields(logrus.Fields{
		"entityID":    req.EntityID.String()[0:8],
		"foeID":       foe.ID.String()[0:8],
		"damage":      computedDamage,
		"Attack":      attackerAttack.I(),
		"WeaponDmg":   weaponDmg.I(),
		"Defense":     foeDefense.I(),
		"Armor":       foeArmor.I(),
		"HP":          foeHP.I(),
		"ResultingHP": foeHP.I() - computedDamage,
		"Shield":      foeShield.GetValue(),
	}).Debug("Entity attack")

	// Update the entity
	gs.Entities[req.EntityID] = ent

	hasActed := ent.GetProperty(property.HasActed)
	hasActed.Set(true)
	ent.UpdateProperty(hasActed)
	hasMoved := ent.GetProperty(property.HasMoved)
	hasMoved.Set(true)
	ent.UpdateProperty(hasMoved)

	// apply damage
	foeHP.SetI(foeHP.I() - computedDamage)
	// update foe's HP
	foe.UpdateProperty(foeHP)
	// Damage fills a channeling foe's interruption gauge (may break its channel).
	// @spec-link [[mechanic_channeling_mechanic]]
	ApplyInterruption(gs, &foe, computedDamage)
	gs.Entities[foe.ID] = foe

	credits := []rulermethods.CreditAward{}
	if computedDamage > 0 {
		credits = append(credits, rulermethods.CreditAward{
			PlayerID: ent.ControllerID,
			Amount:   computedDamage,
			Source:   "damage",
		})
	}

	if foeHP.I() <= 0 {
		ctx.log.WithField("foeID", foe.ID.String()[0:8]).Info("Entity killed in combat")
		gs.RemoveEntity(foe.ID)
	}

	// notify all controllers of the attack.
	notification := rulermethods.ControllerAttacked{
		Entity:               foe,
		Attacker:             ent,
		ControllerID:         foe.ControllerID,
		AttackerControllerID: ent.ControllerID,
		Damage:               computedDamage,
		PrevHP:               foeHP.I() + computedDamage,
		NewHP:                foeHP.I(),
		CreditAwards:         credits,
		Version:              gs.Version,
	}

	gs.IncVersion()
	notification.Version = gs.Version

	for _, ctrl := range gs.Controllers {
		ctrl.NotifyActor(message.Create(nil, notification, nil))
	}

	// reply with the new entities state (opaque to the client)
	reply = msg.Reply()

	reply.Content = rulermethods.ControllerAttackReply{
		Attacker: ent,
		Results: []rulermethods.ActionResult{
			{
				Target:       foe,
				TargetID:     foe.ID,
				Damage:       computedDamage,
				PrevHP:       foeHP.I() + computedDamage,
				NewHP:        foeHP.I(),
				CreditAwards: credits,
			},
		},
	}

	return reply
}


