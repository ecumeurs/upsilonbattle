package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/properties"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/cell"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var defaultAttackProp = properties.DefaultIntProperty(1)
var defaultAttackRangeProp = properties.DefaultIntProperty(1)
var defaultDefenseProp = properties.DefaultIntProperty(0)
var defaultHPProp = properties.DefaultIntProperty(1)

type localAttackCtx struct {
	*GameState
	log *logrus.Entry
}

func (gs *GameState) Attack(msg message.Message, req rulermethods.ControllerAttack) (reply message.Message) {
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
	attackerAttack := ent.GetPropertyI("Attack", &defaultAttackProp)

	foe := gs.Entities[target.EntityID]
	foeDefense := foe.GetPropertyI("Defense", &defaultDefenseProp)
	foeHP := foe.GetPropertyI("HP", &defaultHPProp)

	computedDamage := attackerAttack.I() - foeDefense.I()

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + 500

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

	// notify foe controller of the attack.

	foectrl, found := gs.Controllers[foe.ControllerID]
	if !found {
		ctx.log.WithFields(logrus.Fields{
			"foeID":           foe.ID.String()[0:8],
			"foeControllerID": foe.ControllerID.String()[0:8]}).Error("Foe controller not found")

	} else {
		foectrl.NotifyActor(message.Create(nil, rulermethods.ControllerAttacked{
			Entity:   foe,
			Attacker: ent,
		}, nil))
	}

	// reply with the new entities state (opaque to the client)
	reply = msg.Reply()

	reply.Content = rulermethods.ControllerAttackReply{
		Entity: ent,
	}

	return reply
}

func (ctx *localAttackCtx) preAttackChecks(msg message.Message, req rulermethods.ControllerAttack) (ok bool, reply message.Message) {

	// Check if the controller is allowed to use the entity
	if !ctx.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		ctx.log.Error("Controller is not allowed to use this entity")
		return false, msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.missmatch")
	}

	if ctx.Turner.CurrentEntityTurn != req.EntityID {
		ctx.log.Error("It is not this entity turn")
		return false, msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch")
	}

	ent, found := ctx.Entities[req.EntityID]
	if !found {
		ctx.log.Error("Entity not found")
		return false, msg.ReplyWithError("Entity not found", "entity.notfound")
	}

	// Check if the attack is valid
	target, found := ctx.Grid.CellAt(req.Target)
	if !found {
		ctx.log.Error("Target is not found")
		return false, msg.ReplyWithError("Invalid target", "entity.attack.target.invalid")
	}

	if target.Type != cell.Ground {
		ctx.log.Error("Target is not valid")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.invalid")
	}

	if target.EntityID == uuid.Nil {
		ctx.log.Error("Target has no entities")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.invalid")
	}

	// range check.

	attackerRange := ent.GetPropertyI("AttackRange", &defaultAttackRangeProp)
	distance := ent.Position.Distance(target.Position)
	if attackerRange.I() <= distance {
		ctx.log.WithFields(logrus.Fields{
			"attackrange": attackerRange.I(),
			"distance":    distance,
		}).Error("Target is out of range")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.outofrange")
	}

	return true, reply
}
