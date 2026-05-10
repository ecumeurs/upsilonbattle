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
// @test-link [[mechanic_mech_arena_lifecycle]]
// @test-link [[uc_combat_turn]]

import (
	"reflect"
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
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

// TestRulerControllerCanMoveAttackAndEndTurn simulates a basic player turn sequence, including movement, combat, and turn finalization.
func TestRulerControllerCanMoveAttackAndEndTurn(t *testing.T) {
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

	// Flatten grid to avoid height-blocking in 4x4 tests
	for x := 0; x < ruler.GameState.Grid.Width; x++ {
		for y := 0; y < ruler.GameState.Grid.Length; y++ {
			z := ruler.GameState.Grid.TopMostGroundAt(x, y)
			if cell, ok := ruler.GameState.Grid.CellAt(position.New(x, y, z)); ok {
				if z != 0 {
					delete(ruler.GameState.Grid.Cells, cell.Position)
					cell.Position.Z = 0
					ruler.GameState.Grid.Cells[cell.Position] = cell
				}
			}
		}
	}
	for id, e := range ruler.GameState.Entities {
		e.Position.Z = 0
		ruler.GameState.Entities[id] = e
	}

	ruler.Start()
	defer ruler.Stop()
	dChan := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	replyChan := make(chan *message.Message)

	ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), replyChan)
	msg := <-replyChan
	grd := msg.Content.(rulermethods.GetGridStateReply).Grid

	ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), replyChan)
	msg = <-replyChan

	var entities []entity.Entity
	var foeEntities []entity.Entity
	for _, ent := range msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		if ent.ControllerID == ctrl.ID {
			entities = append(entities, ent)
		} else {
			foeEntities = append(foeEntities, ent)
		}
	}

	attacker := entities[0]
	target := foeEntities[0]

	done := false
	turnTimeout := time.After(15 * time.Second)

	hasTestedMove := false
	hasTestedAttack := false

	for !done {
		select {
		case msg := <-ctrl.Inbox:
			if msg.TargetMethod != nil {
				switch m := msg.TargetMethod.(type) {
				case rulermethods.ControllerNextTurn:
					if m.Entity.ID == attacker.ID {
						if !hasTestedMove {
							var nextPos position.Position
							found := false

							for dx := -1; dx <= 1; dx++ {
								for dy := -1; dy <= 1; dy++ {
									if dx == 0 && dy == 0 {
										continue
									}
									// Only orthogonal adjacency is valid in many turn-based grids, but just in case, check proper neighbors.
									if dx != 0 && dy != 0 {
										continue
									}

									p := position.New(attacker.Position.X+dx, attacker.Position.Y+dy, 0)
									p.Z = grd.TopMostGroundAt(p.X, p.Y)
									logrus.Debugf("Attacker Position: %s destination %s", attacker.Position.String(), p.String())

									if !p.IsAdjacent(attacker.Position, 2) {
										continue
									}

									zDiff := p.Z - attacker.Position.Z
									if zDiff < 0 {
										zDiff = -zDiff
									}

									if c, ok := grd.CellAt(p); ok && !c.IsOccupied() && zDiff <= 2 {
										nextPos = p
										found = true
										break
									}
								}
								if found {
									break
								}
							}

							movePath := []position.Position{}
							if found {
								movePath = append(movePath, nextPos)
							}

							ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
								EntityID:     attacker.ID,
								Path:         movePath,
								ControllerID: ctrl.ID,
							}, rulermethods.ControllerMoveReply{}), ctrl.GetCallbackChan())
						} else if !hasTestedAttack {
							ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
								EntityID:     attacker.ID,
								Target:       target.Position,
								ControllerID: ctrl.ID,
							}, rulermethods.ControllerAttackReply{}), ctrl.GetCallbackChan())
						} else {
							ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
								EntityID:     m.Entity.ID,
								ControllerID: ctrl.ID,
							}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
							done = true
						}
					} else {
						ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
							EntityID:     m.Entity.ID,
							ControllerID: ctrl.ID,
						}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
					}
				case rulermethods.ControllerMoveReply:
					hasTestedMove = true
					ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
						EntityID:     attacker.ID,
						ControllerID: ctrl.ID,
					}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
				case rulermethods.ControllerAttackReply:
					hasTestedAttack = true
					done = true // short-circuit end test!
					ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
						EntityID:     attacker.ID,
						ControllerID: ctrl.ID,
					}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
				}
			}
		case msg := <-ctrl2.Inbox:
			if msg.TargetMethod != nil {
				switch m := msg.TargetMethod.(type) {
				case rulermethods.ControllerNextTurn:
					ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
						EntityID:     m.Entity.ID,
						ControllerID: ctrl2.ID,
					}, rulermethods.EndOfTurn{}), ctrl2.GetCallbackChan())
				}
			}
		case <-turnTimeout:
			t.Fatal("Timeout in battle simulation")
		}
	}

	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{ControllerID: ctrl.ID}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{ControllerID: ctrl2.ID}, nil))

	ctrl.Stop()
	ctrl2.Stop()
}
