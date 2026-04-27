package rules

import (
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)


type localAttackCtx struct {
	*GameState
	log *logrus.Entry
}

// @spec-link [[mech_action_economy_action_cost_rules]]
func (gs *GameState) Attack(msg *message.Message, req rulermethods.ControllerAttack) (reply *message.Message) {
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

	// In multi-entity cells, attack only targets Characters/Monsters, not WalkThrough entities.
	// @spec-link [[mechanic_multi_entity_cell_system]]
	foe, found := gs.FindCharacterInCell(target.EntityIDs)
	if !found {
		return msg.ReplyWithError("No attackable entity at target position", "entity.attack.noentity")
	}
	foeDefense := foe.GetPropertyI(property.Defense)
	foeHP := foe.GetPropertyI(property.HP)

	// @spec-link [[mech_combat_standard_attack_computation]]
	// @spec-link [[mec_backstabbing_mechanic]]
	multiplier := 1.0
	effectiveDefense := foeDefense.I()

	if ent.IsBackstabbing(foe) {
		multiplier = 1.5
		// 50% armor penetration
		effectiveDefense = int(float64(effectiveDefense) * 0.5)
		ctx.log.Debug("Backstab detected! 150% damage and 50% armor penetration applied.")
	}

	computedDamage := tools.Max(1, int(float64(attackerAttack.I())*multiplier)-effectiveDefense)

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
		"Defense":     foeDefense.I(),
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

func (ctx *localAttackCtx) preAttackChecks(msg *message.Message, req rulermethods.ControllerAttack) (ok bool, reply *message.Message) {

	ent, found := ctx.Entities[req.EntityID]
	if !found {
		ctx.log.Error("Entity not found")
		return false, msg.ReplyWithError("Entity not found", "entity.notfound")
	}

	// Check if the controller is allowed to use the entity
	if !ctx.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		ctx.log.Error("Controller is not allowed to use this entity")
		return false, msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.missmatch")
	}

	if ctx.Turner.CurrentEntityTurn != req.EntityID {
		ctx.log.Error("It is not this entity turn")
		return false, msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch")
	}

	// Check if the attack is valid
	target, found := ctx.Grid.CellAt(req.Target)
	if !found {
		ctx.log.Error("Target is not found")
		return false, msg.ReplyWithError("Invalid target", "entity.attack.target.invalid")
	}

	if target.Type != cell.Ground {
		ctx.log.Error("Target is not valid")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.celltype")
	}

	if !target.IsOccupied() {
		ctx.log.Error("Target has no entities")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.noentity")
	}

	// Use FindCharacterInCell to ensure we only target attackable entities.
	targetEntity, found := ctx.FindCharacterInCell(target.EntityIDs)
	if !found {
		ctx.log.Error("Target cell has no attackable character")
		return false, msg.ReplyWithError("No attackable entity", "entity.attack.noentity")
	}

	// range check.
	attackerRange := ent.GetPropertyI(property.AttackRange)
	distance := ent.Position.Distance(target.Position)
	if attackerRange.I() < distance {
		ctx.log.WithFields(logrus.Fields{
			"attackrange": attackerRange.I(),
			"distance":    distance,
			"attacker":    ent.Position,
		}).Error("Target is out of range")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.outofrange")
	}

	// @spec-link [[rule_friendly_fire_team_validation]]
	attackerTeam := ent.GetPropertyI(property.TeamID)
	targetTeam := targetEntity.GetPropertyI(property.TeamID)

	if attackerTeam.I() == targetTeam.I() {
		ctx.log.WithFields(logrus.Fields{
			"attackerTeam": attackerTeam.I(),
			"targetTeam":   targetTeam.I(),
		}).Error("Friendly fire is not allowed")
		return false, msg.ReplyWithError("Friendly fire is not allowed", "entity.attack.friendlyfire")
	}

	propHasActed := ent.GetProperty(property.HasActed).Get().(bool)
	if propHasActed {
		ctx.log.Error("Entity has already acted")
		return false, msg.ReplyWithError("Entity has already acted", "entity.hasacted")
	}

	return true, reply
}
