package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
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

	foe := gs.Entities[target.EntityID]
	foeDefense := foe.GetPropertyI(property.Defense)
	foeHP := foe.GetPropertyI(property.HP)

	// @spec-link [[mech_combat_standard_attack_computation]]
	computedDamage := tools.Max(1,attackerAttack.I() - foeDefense.I())

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + 100

	ctx.log.WithFields(logrus.Fields{
		"entityID":    req.EntityID.String()[0:8],
		"foeID":       target.EntityID.String()[0:8],
		"damage":      computedDamage,
		"Attack":      attackerAttack.I(),
		"Defense":     foeDefense.I(),
		"HP":          foeHP.I(),
		"ResultingHP": foeHP.I() - computedDamage,
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

	if foeHP.I() <= 0 {

		gs.Grid.RemoveEntity(foe.Position)
		delete(gs.Entities, foe.ID)
		gs.Turner.RemoveEntity(foe.ID)

		ctx.log.WithFields(logrus.Fields{
			"entityID": foe.ID.String()[0:8],
			"position": foe.Position}).Info("##### Entity removed #####")
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
	}

	for _, ctrl := range gs.Controllers {
		ctrl.NotifyActor(message.Create(nil, notification, nil))
	}

	gs.IncVersion()

	// reply with the new entities state (opaque to the client)
	reply = msg.Reply()

	reply.Content = rulermethods.ControllerAttackReply{
		Entity: ent,
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

	if target.EntityID == uuid.Nil {
		ctx.log.Error("Target has no entities")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.noentity")
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
	targetEntity := ctx.Entities[target.EntityID]
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
