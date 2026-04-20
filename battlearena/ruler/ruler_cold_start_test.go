package ruler

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestRulerFirstTurnReadiness proves that the first turn triggered by the Ruler
// must contain the correct active entity and a non-zero version.
func TestRulerFirstTurnReadiness(t *testing.T) {
	r := NewCompleteRuler()
	r.NbControllers = 1
	r.Start()
	defer r.Stop()

	ctrl := NewFake("Controller1")
	defer ctrl.Stop()

	// 1. Assign entities to this controller
	for id, ent := range r.GameState.Entities {
		ent.ControllerID = ctrl.ID
		r.GameState.Entities[id] = ent
	}

	dChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// 2. Trigger BattleStart
	r.NotifyActor(message.Create(nil, rulermethods.BattleStart{}, nil))

	// 3. Expect BattleStart (Version 0)
	msg := ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	start := msg.TargetMethod.(rulermethods.BattleStart)
	t.Logf("Received BattleStart: Version=%d", start.Version)

	// 3.5 Pulse readiness (mirroring HTTPController behavior)
	r.NotifyActor(message.Create(nil, rulermethods.ControllerBattleReady{
		ControllerID: ctrl.ID,
	}, nil))

	// 4. Expect ControllerNextTurn (Version 1)
	msg = ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 5*time.Second)
	nextTurn := msg.TargetMethod.(rulermethods.ControllerNextTurn)
	
	t.Logf("Received ControllerNextTurn: EntityID=%s, Version=%d", nextTurn.Entity.ID, nextTurn.Version)

	// ASSERTION: The core of the 30s hang. Turn 1 MUST HAVE an active entity.
	assert.NotEqual(t, uuid.Nil, nextTurn.Entity.ID, "First turn must have a valid active entity ID")
	assert.Greater(t, nextTurn.Version, int64(0), "First turn must have a version > 0 (Turn 1)")
}
