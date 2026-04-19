package ruler

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// @test-link [[rule_turn_clock]]
func TestShotClockExpiry(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.ShotClockDuration = 100 * time.Millisecond
	ruler.Start()
	defer ruler.Stop()
	
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// Setup controllers
	dChan := make(chan *message.Message, 1)
	// Manually assign entities to controllers to ensure win conditions/timeouts work
	i := 0
	for id, e := range ruler.GameState.Entities {
		if i == 0 {
			e.ControllerID = ctrl.ID
		} else {
			e.ControllerID = ctrl2.ID
		}
		ruler.GameState.Entities[id] = e
		i++
	}
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Wait for battle start
	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)
	
	// Wait for first entity turn
	msg := ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)
	entID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID
	
	t.Logf("Turn started for entity %s, waiting for timeout...", entID)

	// Wait for shot clock to expire (100ms) and trigger EndOfTurn
	// This should result in EntitiesStateChanged being broadcast
	ctrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 10*time.Second)
	
	// And then the NEXT turn should be triggered
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// @test-link [[rule_turn_clock]]
func TestShotClockCancellation(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.ShotClockDuration = 500 * time.Millisecond
	ruler.Start()
	defer ruler.Stop()
	
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	dChan := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Manually assign entities to controllers
	i := 0
	for id, e := range ruler.GameState.Entities {
		if i == 0 {
			e.ControllerID = ctrl.ID
		} else {
			e.ControllerID = ctrl2.ID
		}
		ruler.GameState.Entities[id] = e
		i++
	}
	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)
	
	// First turn starts
	msg := ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)
	entID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID

	// Manually end turn before timeout
	ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		ControllerID: ctrl.ID,
		EntityID:     entID,
	}, rulermethods.EndOfTurn{}), dChan)
	<-dChan

	// Wait for transition
	ctrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 10*time.Second)
	
	// Verify next turn triggered
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// @test-link [[rule_turn_clock]]
func TestShotClockTurnProtection(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.ShotClockDuration = 100 * time.Millisecond
	ruler.Start()
	defer ruler.Stop()

	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	dChan := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Start battle
	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)
	// Manually assign entities to controllers
	i := 0
	for id, e := range ruler.GameState.Entities {
		if i == 0 {
			e.ControllerID = ctrl.ID
		} else {
			e.ControllerID = ctrl2.ID
		}
		ruler.GameState.Entities[id] = e
		i++
	}

	// Turn 1.0 starts
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	// Wait for timeout 1 (transitions from 1.0 to 2.0)
	ctrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 10*time.Second)

	// Turn 2.0 starts
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// TestShotClockWithDeadEntity tests that the shot clock properly handles
// the case where the current entity is killed after the shot clock has started
// but before it fires. This is a regression test for the complete ISS-046 fix.
//
// Scenario:
// 1. Entity A's turn starts, shot clock begins
// 2. Another entity kills Entity A during its turn
// 3. Entity A is removed from the game state
// 4. Shot clock fires and tries to send timeout EndOfTurn
// 5. Without fixes, the shot clock would try to send timeout to the dead entity
// 6. With fixes, the shot clock should detect the entity is dead and skip safely
func TestShotClockWithDeadEntity(t *testing.T) {
	r := NewCompleteRuler()
	r.Start()
	defer r.Stop()

	// Set a short shot clock for testing
	r.ShotClockDuration = 200 * time.Millisecond

	ctrl1 := NewFake("Controller1")
	ctrl2 := NewFake("Controller2")

	// Setup controllers
	dChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl1, ControllerID: ctrl1.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	r.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Wait for battle start
	ctrl1.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)
	// Manually assign entities to controllers
	i := 0
	for id, e := range r.GameState.Entities {
		if i == 0 {
			e.ControllerID = ctrl1.ID
		} else {
			e.ControllerID = ctrl2.ID
		}
		r.GameState.Entities[id] = e
		i++
	}

	// Get first entity turn
	msg := ctrl1.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)
	currentEntID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID

	// Now simulate the current entity being killed by an external action
	// This tests the scenario where an entity dies after the shot clock has started
	// but before it fires
	r.NotifyActor(message.Create(nil, rulermethods.TestingDeleteEntity{EntityID: currentEntID}, nil))

	// Allow some time for the removal notification to be processed on the actor thread
	time.Sleep(50 * time.Millisecond)

	// Verify CurrentEntityTurn was cleared when the entity was removed
	// We use TestingGetState to avoid data races
	replyChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.TestingGetState{}, rulermethods.TestingGetStateReply{}), replyChan)
	replyMsg := <-replyChan
	
	if replyMsg.TargetMethod.(rulermethods.TestingGetStateReply).CurrentEntityTurn == currentEntID {
		t.Fatal("CurrentEntityTurn should have been cleared when the entity was removed")
	}

	// Wait for shot clock to fire (200ms duration)
	time.Sleep(300 * time.Millisecond)

	// The shot clock should have fired and detected the entity was dead
	// Without the fix, this would have caused an error or panic when trying to
	// look up the dead entity in r.GameState.Entities
	// With the fix, it should have logged an error but not crashed

	// Test passes if we get here without panic
	// The battle should continue normally without trying to send timeout to dead entity

	ctrl1.Stop()
	ctrl2.Stop()
}
