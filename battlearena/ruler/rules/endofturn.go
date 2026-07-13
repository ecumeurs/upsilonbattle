package rules

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// EndOfTurn handles the completion of an entity's turn.
// It processes turn-end costs (delays, penalties), status effects (poison),
// resets movement credits, decrements temporary entity durations, and advances the global turn counter.
// @spec-link [[rule_turn_clock]]
// @spec-link [[mech_action_economy]]
func EndOfTurn(gs *gamestate.GameState, msg *message.Message, req rulermethods.EndOfTurn, ent entity.Entity) (ok bool, reply *message.Message) {
	loclog := gs.Logger.WithFields(logrus.Fields{
		"RequestID":    msg.RequestId.String()[0:8],
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8],
		"rule":         "endofturn",
	})
	loclog.Debug("End of turn request")

	if req.EntityID == uuid.Nil {
		loclog.Error("Can't work with nil entity")
		return false, msg.ReplyWithError("Can't work with nil entity", "entity.nil")
	}

	if _, found := gs.Entities[req.EntityID]; !found {
		loclog.WithFields(logrus.Fields{
			"RequestID":    msg.RequestId.String()[0:8],
			"controllerID": req.ControllerID.String()[0:8],
			"entityID":     req.EntityID.String()[0:8]}).Error("Can't work with absent entity")

		return false, msg.ReplyWithError("Can't work with absent entity", "entity.absent")
	}

	// Check if the controller is allowed to end the turn
	// Check if the controller is allowed to use the entity
	if !gs.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		loclog.WithFields(logrus.Fields{
			"controllerID":       req.ControllerID.String()[0:8],
			"entityID":           req.EntityID.String()[0:8],
			"entityControllerID": gs.Entities[req.EntityID].ControllerID.String()[0:8],
		}).Error("Controller is not allowed to use this entity")
		return false, msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.mismatch")
	}
	if gs.Turner.CurrentEntityTurn != req.EntityID {
		loclog.Error("It is not this entity turn")
		return false, msg.ReplyWithError("It is not this entity turn", "entity.turn.mismatch")
	}

	loclog.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8],
		"isTimeout": req.IsTimeout}).Debug("Entity end of turn")

	// Channeled casters auto-pass: instead of the flat Pass delay, reschedule the
	// caster out by the channel delay and keep it locked (HasActed/HasMoved retained,
	// no movement restore, no poison tick — those happen at the resolution's normal
	// end of turn). The effect resolves when the dormant caster is re-picked.
	// @spec-link [[mechanic_channeling_mechanic]]
	if ent.IsChanneling() {
		chSkill := ent.Skills[ent.IsCasting.SkillID]
		ent.CurrentDelay += channelDelay(chSkill)
		gs.Entities[req.EntityID] = ent
		gs.Turner.AddEntity(req.EntityID, ent.CurrentDelay)
		gs.IncTurn()

		for _, ctrl := range gs.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.ControllerPassed{
				EntityID:     req.EntityID,
				ControllerID: req.ControllerID,
				Version:      gs.Version,
			}, nil))
		}
		return true, msg
	}

	delay := 300 // Base Pass cost as per [[mech_action_economy]]
	if req.IsTimeout {
		delay += 100 // Penalty as per [[us_take_combat_turn]]
	}
	ent.CurrentDelay += delay

	// check poisonned status!

	if ent.GetPropertyI(property.Poison).I() > 0 {
		poison := ent.GetPropertyI(property.Poison).I()
		ent.UpdatePropertyValue(property.HP, ent.GetPropertyC(property.HP).GetValue()-poison)

		poison = poison / 2
		if poison <= 1 {
			poison = 0
		}
		ent.UpdatePropertyValue(property.Poison, poison)
		if ent.GetPropertyC(property.HP).GetValue() <= 0 {
			ent.UpdatePropertyValue(property.HP, 1)
			// can't die from poison.... Just suffer greatly.
		}
	}

	// restore movement points
	ent.UpdatePropertyValue(property.Movement, ent.GetPropertyC(property.Movement).GetMaxValue())
	// unset HasActed, HasMoved
	ent.UpdatePropertyValue(property.HasActed, false)
	ent.UpdatePropertyValue(property.HasMoved, false)
	// tick every equipped skill's cooldown down by one (floored at 0) so a cast
	// skill becomes castable again once its cooldown counter reaches 0. A skill
	// cast this turn was set to its max cooldown just now, so a single decrement
	// here leaves it > 0 — still rejected on the owner's immediately-following
	// turn — as [[mech_skill_validation]]'s cooldown gate requires. (ISS-111)
	// @spec-link [[mech_skill_validation]]
	ent.SkillCooldownTickDown()

	// update entity in state
	gs.Entities[req.EntityID] = ent

	// @spec-link [[mechanic_expiration_controller]]
	// Decrement EntityDuration for temporary entities.
	// EntityDuration == 0 means permanent (no expiry).
	if durProp := ent.GetProperty(property.EntityDuration); durProp != nil {
		dur := ent.GetPropertyC(property.EntityDuration).GetValue()
		if dur > 0 {
			dur--
			if dur == 0 {
				loclog.WithFields(logrus.Fields{
					"entityID": req.EntityID.String()[0:8],
					"reason":   "duration_expired",
				}).Info("Temporary entity expired — removing from arena")
				// Do NOT re-add to Turner. RemoveEntity handles all cleanup.
				gs.RemoveEntity(req.EntityID)
				gs.IncTurn()
				// Notify controllers: skip the ControllerPassed broadcast for dead entity.
				return true, msg
			}
			ent.UpdatePropertyValue(property.EntityDuration, dur)
			gs.Entities[req.EntityID] = ent
		}
	}

	gs.Turner.AddEntity(req.EntityID, gs.Entities[req.EntityID].CurrentDelay) // well ...end of turn delay
	gs.IncTurn()

	// notify all controllers of the pass/turn end.
	for _, ctrl := range gs.Controllers {
		ctrl.NotifyActor(message.Create(nil, rulermethods.ControllerPassed{
			EntityID:     req.EntityID,
			ControllerID: req.ControllerID,
			Version:      gs.Version,
		}, nil))
	}

	return true, msg
}
