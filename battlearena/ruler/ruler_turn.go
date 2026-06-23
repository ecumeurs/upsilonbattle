package ruler

import (

	"fmt"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rules"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// processTriggerFirstTurn is a notification handler that initiates the first turn of the battle.
func (r *Ruler) processTriggerFirstTurn(ctx actor.NotificationContext) {
	r.triggerFirstTurn()
}

// triggerFirstTurn identifies and hands the turn to the first eligible entity.
func (r *Ruler) triggerFirstTurn() {
	if r.firstTurnSent {
		return
	}

	entID := r.GameState.Turner.CurrentEntityTurn
	if entID == uuid.Nil {
		if len(r.GameState.Turner.Turns) == 0 {
			r.RequestLogger.Warn("Turner is empty, cannot trigger first turn")
			return
		}

		candidateID := r.GameState.Turner.Turns[0].EntityId
		ent, ok := r.GameState.Entities[candidateID]

		if !ok || ent.ControllerID == uuid.Nil {
			r.RequestLogger.WithFields(logrus.Fields{
				"entityID": candidateID.String()}).Debug("First entity in queue is uncontrolled, waiting for readiness")
			return
		}

		r.RequestLogger.Info("Picking first entity to play")
		entID = r.GameState.Turner.NextTurn()
	}

	ent, ok := r.GameState.Entities[entID]
	if !ok || ent.ControllerID == uuid.Nil {
		r.RequestLogger.WithFields(logrus.Fields{
			"entityID": entID.String()}).Warn("Current entity turn is invalid/uncontrolled, skipping (waiting for recovery)")
		return
	}

	r.RequestLogger.WithFields(logrus.Fields{
		"entityID": entID.String()[0:8]}).Info("Handing turn to first entity")

	r.firstTurnSent = true
	r.GameState.IncTurn() // @spec-link [[mech_game_state_versioning]]
	
	r.handTurn(entID)
}

// controllerTurnReady is a notification handler that confirms a controller is ready for its turn.
func (r *Ruler) controllerTurnReady(ctx actor.NotificationContext) {
	r.RequestLogger.Info("ControllerTurnReady")
}

// handTurn transfers turn authority to an entity and starts its shot clock.
func (r *Ruler) handTurn(entID uuid.UUID) {
	ent, found := r.GameState.Entities[entID]
	if !found {
		r.logger.WithField("entityID", entID.String()[0:8]).Error("Cannot hand turn: Entity not found")
		return
	}

	ent.CurrentDelay = 0
	r.GameState.Entities[entID] = ent

	// Every entity must have an owning controller (ISS-101) — AddEntity
	// rejects uuid.Nil at registration time, so reaching this point with a
	// Nil ControllerID means that invariant was bypassed somewhere. This
	// must never be treated as "automated"; surface it loudly instead of
	// silently falling through to ExpirationBehavior.
	if ent.ControllerID == uuid.Nil {
		panic(fmt.Sprintf("handTurn: entity %s has no ControllerID (uuid.Nil) — every entity must have an owning controller (ISS-101)", entID))
	}

	_, found = r.GameState.Controllers[ent.ControllerID]
	if !found {
		r.RequestLogger.WithFields(logrus.Fields{
			"entityID":     entID.String()[0:8],
			"controllerID": ent.ControllerID.String()[0:8]}).Error("Controller not found for entity")
		return
	}

	// @spec-link [[rule_turn_clock]]
	r.startShotClock()
	for _, c := range r.GameState.Controllers {
		c.NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
			Entity:  ent,
			Turn:    r.GameState.Turner.GetTurnState(),
			Version: r.GameState.Version,
		}, nil))
	}
}

// endOfTurn finalizes the current entity's turn and picks the next one.
// @spec-link [[mech_action_economy]]
func (r *Ruler) endOfTurn(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.EndOfTurn)
	r.RequestLogger = r.RequestLogger.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8]})
	r.RequestLogger.Debug("End of turn request")

	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	// shot clock validation
	if req.IsTimeout && req.TurnIndex != 0 && req.TurnIndex != r.GameState.GetTurn() {
		r.logger.WithFields(logrus.Fields{
			"reqTurn":     req.TurnIndex,
			"currentTurn": r.GameState.GetTurn()}).Debug("Ignoring late timeout message")
		ctx.Reply(ctx.Msg.Reply())
		return
	}

	r.stopShotClock()

	ok, reply := rules.EndOfTurn(r.GameState, ctx.Msg, req, r.GameState.Entities[req.EntityID])
	if !ok {
		ctx.Reply(reply)
		return
	}

	nextTurnEnt := r.GameState.Turner.NextTurn()

	// @spec-link [[mech_initiative]]
	for nextTurnEnt != uuid.Nil {
		if _, alive := r.GameState.Entities[nextTurnEnt]; alive {
			break
		}
		r.RequestLogger.WithFields(logrus.Fields{
			"skippedEntityID": nextTurnEnt.String()[0:8],
		}).Warn("Next-turn entity was already dead; skipping to next in queue")
		nextTurnEnt = r.GameState.Turner.NextTurn()
	}

	if nextTurnEnt == uuid.Nil {
		r.RequestLogger.Info("##### END OF BATTLE! (WEIRD) #####")
	} else {
		if beg, found := r.GameState.Entities[nextTurnEnt]; found {
			rules.BeginingOfTurn(r.GameState, beg)
		}
	}

	ent := make([]entity.Entity, 0)
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	if nextTurnEnt == uuid.Nil {
		r.evaluateVictory(nextTurnEnt)
	} else {
		remainingTeams := make(map[int]bool)
		for _, ent := range r.GameState.Entities {
			remainingTeams[ent.GetPropertyI(property.TeamID).I()] = true
		}

		if len(remainingTeams) <= 1 {
			r.evaluateVictory(nextTurnEnt)
		} else {
			r.RequestLogger.Info("##### END OF TURN #####")
			for _, ctrl := range r.GameState.Controllers {
				ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
					Entities: ent,
					Turn:     r.GameState.Turner.GetTurnState(),
					Version:  r.GameState.Version,
				}, nil))
			}

			if nextTurnEnt != uuid.Nil {
				r.handTurn(nextTurnEnt)
			}
		}
	}

	ctx.Reply(ctx.Msg.Reply())
}




