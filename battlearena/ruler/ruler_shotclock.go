package ruler

import (
	"fmt"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

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
