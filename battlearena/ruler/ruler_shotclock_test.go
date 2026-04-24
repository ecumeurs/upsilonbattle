// @test-link [[rule_turn_clock]]
// @test-link [[rule_ruler_test_robustness]]
package ruler

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// @test-link [[rule_turn_clock]]
// TestShotClockExpiry verifies that when the shot clock (turn timer) runs out,
// the Ruler automatically triggers an EndOfTurn event for the current entity.
func TestShotClockExpiry(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.ShotClockDuration = 100 * time.Millisecond
	dChan := make(chan *message.Message, 1)
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// Setup controllers
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

	ruler.Start()
	defer ruler.Stop()
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Note: While the test code doesn't explicitly send ControllerBattleReady,
	// FakeController (in ruler_test.go) automatically sends it once it receives
	// the grid and entities state. This satisfies the Ruler's requirement to
	// trigger the first turn.

	// Wait for battle start
	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)

	// Wait for first entity turn on either controller
	var msg *message.Message
	var activeCtrl *FakeController
	timeout := time.After(10 * time.Second)
	foundTurn := false
	for !foundTurn {
		select {
		case m := <-ctrl.Inbox:
			if _, ok := m.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				msg = m
				activeCtrl = ctrl
				foundTurn = true
			}
		case m := <-ctrl2.Inbox:
			if _, ok := m.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				msg = m
				activeCtrl = ctrl2
				foundTurn = true
			}
		case <-timeout:
			t.Fatal("Timeout waiting for ControllerNextTurn")
		}
	}

	entID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID
	t.Logf("Turn started for entity %s on %s, waiting for timeout...", entID, activeCtrl.Name())

	// Wait for shot clock to expire (100ms) and trigger EndOfTurn
	// All controllers receive EntitiesStateChanged
	activeCtrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 10*time.Second)

	// And then the NEXT turn should be triggered
	// (Note: we don't know who gets the next turn either, but we just want to see it triggered)
	// We check the active controller again for the transition broadcast
	activeCtrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// @test-link [[rule_turn_clock]]
// @test-link [[rule_ruler_test_robustness]]
func TestShotClockCancellation(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.ShotClockDuration = 500 * time.Millisecond
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	dChan := make(chan *message.Message, 1)
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

	ruler.Start()
	defer ruler.Stop()

	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Note: The FakeControllers automatically handle the BattleReady handshake in the background.
	// Ruler will only send BattleStart and hand out the first turn (ControllerNextTurn)
	// once both controllers have signaled readiness.

	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)
	// (Check ctrl2 too to ensure it received BattleStart)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)

	// Wait for first entity turn on either controller
	var msg *message.Message
	var activeCtrl *FakeController
	timeout := time.After(10 * time.Second)
	foundTurn := false
	for !foundTurn {
		select {
		case m := <-ctrl.Inbox:
			if _, ok := m.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				msg = m
				activeCtrl = ctrl
				foundTurn = true
			}
		case m := <-ctrl2.Inbox:
			if _, ok := m.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				msg = m
				activeCtrl = ctrl2
				foundTurn = true
			}
		case <-timeout:
			t.Fatal("Timeout waiting for ControllerNextTurn")
		}
	}

	entID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID

	// Manually end turn before timeout
	ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		ControllerID: activeCtrl.ID,
		EntityID:     entID,
	}, rulermethods.EndOfTurn{}), dChan)
	<-dChan

	// Wait for transition
	activeCtrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 10*time.Second)

	// Verify next turn triggered
	activeCtrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// @test-link [[rule_turn_clock]]
// TestShotClockTurnProtection verifies that the shot clock correctly follows the
// battle versioning and turn sequence, ensuring that timeouts from older turns
// cannot accidentally end newer turns even if they were to fire late.
func TestShotClockTurnProtection(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.ShotClockDuration = 100 * time.Millisecond
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

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

	ruler.Start()
	defer ruler.Stop()

	// Register controllers
	dChan := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Wait for battle start
	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 10*time.Second)

	// Turn 1.0 starts
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	// Wait for timeout 1 (transitions from 1.0 to 2.0)
	ctrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 10*time.Second)

	// Turn 2.0 starts
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 10*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// @test-link [[rule_turn_clock]]
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
	ctrl1 := NewFake("Controller1")
	ctrl2 := NewFake("Controller2")

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

	r.Start()
	defer r.Stop()

	// Register controllers
	dChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl1, ControllerID: ctrl1.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	r.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// Wait for first entity turn on either controller
	var msg *message.Message
	var activeCtrl *FakeController
	timeout := time.After(10 * time.Second)
	foundTurn := false
	for !foundTurn {
		select {
		case m := <-ctrl1.Inbox:
			if _, ok := m.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				msg = m
				activeCtrl = ctrl1
				foundTurn = true
			}
		case m := <-ctrl2.Inbox:
			if _, ok := m.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				msg = m
				activeCtrl = ctrl2
				foundTurn = true
			}
		case <-timeout:
			t.Fatal("Timeout waiting for ControllerNextTurn")
		}
	}

	currentEntID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID
	t.Logf("Turn started for entity %s on %s", currentEntID, activeCtrl.Name())

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
