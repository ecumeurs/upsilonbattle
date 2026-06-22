package ruler
// @test-link [[mechanic_arena_lifecycle]]
// @test-link [[uc_combat_turn]]

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestRulerEntityLeak specifically reproduces the bug where triggerFirstTurn
// pops an entity from the Turner but returns early if the entity is uncontrolled,
// effectively leaking the entity and eventually draining the Turner.
func TestRulerEntityLeak(t *testing.T) {
	r := NewCompleteRuler()
	r.Start()
	defer r.Stop()

	// 1. Setup two controllers
	ctrl1 := NewFake("Controller1")
	ctrl2 := NewFake("Controller2")
	defer ctrl1.Stop()
	defer ctrl2.Stop()

	dChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl1, ControllerID: ctrl1.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	r.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// 2. Intentionally mess up the ControllerID of all entities to trigger the leak guard
	// This simulates a race where entities are not yet fully assigned when triggerFirstTurn is called
	for id, ent := range r.GameState.Entities {
		ent.ControllerID = uuid.Nil
		r.GameState.Entities[id] = ent
	}

	initialTurnerCount := len(r.GameState.Turner.Turns)
	t.Logf("Initial turner count: %d", initialTurnerCount)

	// 3. Manually trigger BattleStart multiple times
	// Since we messed up ControllerID, triggerFirstTurn will return early EACH TIME
	// and should NOT leak entities.
	for i := 0; i < 3; i++ {
		r.NotifyActor(message.Create(nil, rulermethods.BattleStart{}, nil))
		time.Sleep(100 * time.Millisecond) // Give actor some time to process
	}

	// 4. Check if we leaked entities
	// Without the fix, turner_count will have reduced by 3.
	// With the fix, it should remain unchanged.
	currentTurnerCount := len(r.GameState.Turner.Turns)
	t.Logf("Current turner count: %d", currentTurnerCount)

	assert.Equal(t, initialTurnerCount, currentTurnerCount, "Turner should not have leaked entities due to failed triggerFirstTurn calls")
}
