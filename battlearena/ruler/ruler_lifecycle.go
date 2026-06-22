package ruler

import (
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// isBattleReadyToStart checks if all required components (Grid, Controllers) are present.
// @spec-link [[rule_spec_match_format_ready_to_start_rule]]
func (r *Ruler) isBattleReadyToStart() bool {
	if r.GameState.Grid == nil {
		r.logger.Debug("isBattleReadyToStart: Grid is nil")
		return false
	}
	if len(r.GameState.Controllers) != r.NbControllers {
		r.logger.WithFields(logrus.Fields{
			"current": len(r.GameState.Controllers),
			"target":  r.NbControllers,
		}).Debug("isBattleReadyToStart: Not all controllers added yet")
		return false
	}
	return true
}

// isBattleReadyToExecute checks if the battle is in progress and all controllers have signaled readiness.
// @spec-link [[rule_battle_readiness]]
func (r *Ruler) isBattleReadyToExecute() bool {
	if r.CurrentState != InProgress {
		return false
	}
	if !r.isBattleReadyToStart() {
		return false
	}
	if len(r.ControllerBattleReady) != r.NbControllers {
		return false
	}
	for _, ready := range r.ControllerBattleReady {
		if !ready {
			return false
		}
	}
	return true
}

// battleStart handles the transition to the InProgress state and notifies controllers.
func (r *Ruler) battleStart(ctx actor.NotificationContext) {
	r.RequestLogger.Info("Processing BattleStart internal notification")
	r.CurrentState = InProgress

	// @spec-link [[rule_turn_clock]]
	r.startShotClock()

	r.RequestLogger.Info("Broadcasting BattleStart to all controllers")
	for id, c := range r.GameState.Controllers {
		r.RequestLogger.WithFields(logrus.Fields{"target": id}).Debug("Sending BattleStart")
		c.NotifyActor(message.Create(nil, rulermethods.BattleStart{
			Turn:    r.GameState.Turner.GetTurnState(),
			Version: r.GameState.Version,
		}, nil))
	}

	if r.isBattleReadyToExecute() {
		r.RequestLogger.Info("All controllers already ready, triggering first turn immediately")
		r.SelfNotifyDelayed(rulermethods.InternalTriggerFirstTurn{}, 100*time.Millisecond)
	} else {
		r.RequestLogger.Info("Waiting for controllers to signal readiness before first turn")
	}
}

// controllerBattleReady marks a controller as ready to begin the match.
// @spec-link [[rule_battle_readiness]]
func (r *Ruler) controllerBattleReady(ctx actor.NotificationContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerBattleReady)
	r.RequestLogger.Info("ControllerBattleReady")

	if r.CurrentState == Finished {
		r.RequestLogger.Warn("Received ControllerBattleReady after battle finished")
		return
	}

	r.ControllerBattleReady[req.ControllerID] = true

	if r.isBattleReadyToExecute() {
		r.SelfNotifyDelayed(rulermethods.InternalTriggerFirstTurn{}, 100*time.Millisecond)
	}
}

// Resurrect configures the Ruler to resume an arena from persisted state (ISS-054).
func (r *Ruler) Resurrect(turns []turner.EntityTurn, currentEntityID uuid.UUID, version int64) {
	r.GameState.Turner.Turns = turns
	r.GameState.Turner.CurrentEntityTurn = currentEntityID
	r.GameState.Version = version
	r.GameState.TurnIndex = uint32(version >> 32)
	r.GameState.ActionIndex = uint32(version & 0xFFFFFFFF)
	r.CurrentState = InProgress
	r.firstTurnSent = true
}

// resurrect is the notification handler triggered after all controllers register during resurrection.
func (r *Ruler) resurrect(ctx actor.NotificationContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.Resurrect)
	r.RequestLogger.WithFields(logrus.Fields{
		"entityID": req.CurrentEntityID.String()[0:8],
	}).Info("Resuming resurrected arena \u2014 handing turn to current entity")
	r.startShotClock()
	r.handTurn(req.CurrentEntityID)
}

// actorAboutToStop handles the graceful shutdown of the ruler actor.
// @spec-link [[mechanic_arena_lifecycle]]
func (r *Ruler) actorAboutToStop(ctx actor.NotificationContext) {
	r.logger.Info("Ruler is about to stop, stopping all controllers and timers")
	r.stopShotClock()
	stopped := make(map[actor.Communication]bool)
	for id, ctrl := range r.GameState.Controllers {
		if stopped[ctrl] {
			continue
		}
		r.logger.Infof("Stopping controller %s", id)
		ctrl.NotifyActor(message.Create(nil, actor.ActorStop{}, nil))
		stopped[ctrl] = true
	}
}
