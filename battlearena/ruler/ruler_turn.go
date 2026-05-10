package ruler

import (
	"fmt"
	"time"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/behavior"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
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
		
		hasBehavior := ent.HasProperty(property.AIBehavior) && ent.GetProperty(property.AIBehavior).Get().(string) != "none"
		if !ok || (ent.ControllerID == uuid.Nil && !hasBehavior) {
			r.RequestLogger.WithFields(logrus.Fields{
				"entityID": candidateID.String()}).Debug("First entity in queue is uncontrolled and has no behavior, waiting for readiness")
			return
		}

		r.RequestLogger.Info("Picking first entity to play")
		entID = r.GameState.Turner.NextTurn()
	}

	ent, ok := r.GameState.Entities[entID]
	hasBehavior := ent.HasProperty(property.AIBehavior) && ent.GetProperty(property.AIBehavior).Get().(string) != "none"
	if !ok || (ent.ControllerID == uuid.Nil && !hasBehavior) {
		r.RequestLogger.WithFields(logrus.Fields{
			"entityID": entID.String()}).Warn("Current entity turn is invalid/uncontrolled/no-behavior, skipping (waiting for recovery)")
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

	// @spec-link [[mech_behavior_system]]
	behaviorProp := ent.GetProperty(property.AIBehavior)
	behaviorSlug := "none"
	if behaviorProp != nil {
		behaviorSlug = behaviorProp.Get().(string)
	}

	if behaviorSlug != "none" || ent.ControllerID == uuid.Nil {
		r.logger.WithFields(logrus.Fields{
			"entityID": entID.String()[0:8],
			"behavior": behaviorSlug,
		}).Info("Executing automated behavior")

		b := behavior.GetBehavior(behaviorSlug)
		msg := b.Decide(r.GameState, ent)
		r.SelfDispatchMessageDelayed(msg, 50*time.Millisecond)
		return
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
// @spec-link [[mech_action_economy_action_cost_rules]]
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

	ok, reply := r.GameState.EndOfTurn(ctx.Msg, req, r.GameState.Entities[req.EntityID])
	if !ok {
		ctx.Reply(reply)
		return
	}

	nextTurnEnt := r.GameState.Turner.NextTurn()

	// @spec-link [[mech_initiative_active_state]]
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
			r.GameState.BeginingOfTurn(beg)
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

// evaluateVictory checks if a win condition is met and transitions the state if necessary.
func (r *Ruler) evaluateVictory(nextTurnEnt uuid.UUID) {
	remainingTeams := make(map[int]bool)
	winningTeamID := 0
	for _, ent := range r.GameState.Entities {
		remainingTeams[ent.GetPropertyI(property.TeamID).I()] = true
		winningTeamID = ent.GetPropertyI(property.TeamID).I()
	}

	if len(remainingTeams) <= 1 || nextTurnEnt == uuid.Nil {
		// @spec-link [[rule_team_mechanics]]
		r.RequestLogger.Info("##### END OF BATTLE! #####")
		r.CurrentState = Finished
		r.GameState.WinnerTeamID = winningTeamID
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerTeamID: winningTeamID,
				Version:      r.GameState.Version,
			}, nil))
		}
	}
}

// startShotClock initializes and starts the turn timer.
// @spec-link [[rule_turn_clock]]
func (r *Ruler) startShotClock() {
	r.stopShotClock()

	if r.ShotClockDuration <= 0 {
		return
	}

	turn := uint32(r.GameState.GetTurn())

	r.logger.WithFields(logrus.Fields{
		"turn":    turn,
		"version": fmt.Sprintf("%d.%d", turn, r.GameState.GetAction()),
		"timeout": r.ShotClockDuration.String()}).Info("Starting turn shot clock")

	r.shotClock = time.AfterFunc(r.ShotClockDuration, func() {
		r.NotifyActor(message.Create(nil, rulermethods.Timeout{TurnIndex: turn}, nil))
	})
}

// timeout handles the turn expiration safely within the actor loop.
func (r *Ruler) timeout(ctx actor.NotificationContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.Timeout)

	if uint32(r.GameState.GetTurn()) != req.TurnIndex {
		r.logger.WithFields(logrus.Fields{
			"capturedTurn": req.TurnIndex,
			"currentTurn":  r.GameState.GetTurn()}).Debug("Shot clock expired but turn already progressed, ignoring.")
		return
	}

	r.logger.Warn("Turn timeout detected! Forcing EndOfTurn.")

	currentEntityID := r.GameState.Turner.CurrentEntityTurn
	if currentEntityID == uuid.Nil {
		return
	}

	ent, found := r.GameState.Entities[currentEntityID]
	if !found {
		return
	}

	r.endOfTurn(actor.CallContext{
		Msg: message.Create(nil, rulermethods.EndOfTurn{
			ControllerID: ent.ControllerID,
			EntityID:     currentEntityID,
			IsTimeout:    true,
			TurnIndex:    uint32(r.GameState.GetTurn()),
		}, nil),
	})
}

// stopShotClock stops the active shot clock timer.
func (r *Ruler) stopShotClock() {
	if r.shotClock != nil {
		r.shotClock.Stop()
		r.shotClock = nil
	}
}
