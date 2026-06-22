/*
Package ruler contains the tests for the Ruler actor.

# Testing Pattern

The tests in this file validate the core logic of the `Ruler` by creating a real
`Ruler` instance and attaching fake controllers (`FakeController`) to it via the
actor framework.

Because the Ruler and Controllers run in separate goroutines and communicate
asynchronously via message queues, synchronizing the assertions in these tests
is done using an "EventStream" (Inbox) pattern.

# The FakeController

`FakeController` acts as a dummy controller that wiretaps the messages the Ruler
broadcasts.
 1. When the fake controller receives a message, its dummy handlers immediately push
    that message onto a buffered Inbox channel.
 2. The test explicitly waits for messages in this channel using ExpectMessage, which
    deals with out-of-order events from async dispatch cleanly.
 3. For complex loops, tests can select directly over the Inboxes.

This allows tests to deterministically simulate a multi-turn battle loop without
race conditions or sleep timers.
*/
package ruler
// @test-link [[mechanic_arena_lifecycle]]
// @test-link [[uc_combat_turn]]

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"

	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

// init initializes the logrus formatter and level for the test suite.
func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	tools.SeedWith(42)
}


// TestRulerBattleBegin verifies that the battle starts correctly when multiple controllers are added.
func TestRulerBattleBegin(t *testing.T) {
	ruler := NewCompleteRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// Manually assign entities as per modern usage
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

	dChan1 := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan1)
	<-dChan1

	dChan2 := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan2)
	<-dChan2

	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// TestRulerBattleBeginNextTurn verifies that the ruler correctly transitions to the next turn after starting.
func TestRulerBattleBeginNextTurn(t *testing.T) {
	ruler := NewCompleteRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")
	// Manually assign entities
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

	dChan1 := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan1)
	<-dChan1

	dChan2 := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan2)
	<-dChan2

	timeout := time.After(5 * time.Second)
	foundTurn := false
	for !foundTurn {
		select {
		case msg := <-ctrl.Inbox:
			if msg.TargetMethod != nil {
				if m, ok := msg.TargetMethod.(rulermethods.ControllerNextTurn); ok {
					if m.Entity.ControllerID != ctrl.ID {
						t.Error("NextTurn received by ctrl but ControllerID mismatch")
					}
					foundTurn = true
				}
			}
		case msg := <-ctrl2.Inbox:
			if msg.TargetMethod != nil {
				if m, ok := msg.TargetMethod.(rulermethods.ControllerNextTurn); ok {
					if m.Entity.ControllerID != ctrl2.ID {
						t.Error("NextTurn received by ctrl2 but ControllerID mismatch")
					}
					foundTurn = true
				}
			}
		case <-timeout:
			t.Fatal("Timeout waiting for ControllerNextTurn")
		}
	}

	ctrl.Stop()
	ctrl2.Stop()
}

// TestRulerBattleBeginNextTurnFetchGridAndEntities ensures that controllers can successfully fetch the current grid and entity state after the battle starts.
func TestRulerBattleBeginNextTurnFetchGridAndEntities(t *testing.T) {
	ruler := NewCompleteRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")
	// Manually assign entities
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

	dChan1 := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan1)
	<-dChan1

	dChan2 := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan2)
	<-dChan2

	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	replyChan := make(chan *message.Message)

	ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), replyChan)
	msg := <-replyChan
	grd := msg.Content.(rulermethods.GetGridStateReply).Grid
	if grd == nil {
		t.Error("Grid should not be nil")
	}

	ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), replyChan)
	msg = <-replyChan
	entities := msg.Content.(rulermethods.GetEntitiesStateReply).Entities
	if len(entities) < 2 {
		t.Error("Entities should have at least 2 entities")
	}

	appropriateEntitiesCounter := 0
	var myEntities []entity.Entity
	for _, ent := range entities {
		if ent.ControllerID == ctrl.ID {
			appropriateEntitiesCounter++
			myEntities = append(myEntities, ent)
		}
	}
	if appropriateEntitiesCounter < 1 {
		t.Error("Entities should have at least 1 entity for controller")
	}

	for _, ent := range myEntities {
		c, found := grd.CellAt(ent.Position)
		if !found {
			t.Error("Position should be on board")
		}
		if !c.IsOccupied() {
			t.Error("Cell should have an entity")
		}
		if !c.HasEntity(ent.ID) {
			t.Error("Cell should contain this entity")
		}
	}

	ctrl.Stop()
	ctrl2.Stop()
}

