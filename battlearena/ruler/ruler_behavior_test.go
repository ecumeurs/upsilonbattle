package ruler
// @test-link [[mechanic_arena_lifecycle]]
// @test-link [[mech_behavior_system]]

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// TestRulerAggressiveBehavior validates that entities with the "aggressive" behavior correctly move toward enemies.
func TestRulerAggressiveBehavior(t *testing.T) {
	tools.Seed()
	r := NewRuler(uuid.New())
	r.NbControllers = 1
	r.NbEntitiesPerController = 1
	r.ShotClockDuration = 0 // Disable shot clock for testing

	r.GameState.Grid = grid.NewGrid(10, 10, 3)

	// Create a character
	char := entity.New()
	char.ID = uuid.New()
	char.Type = entity.Character
	char.ControllerID = uuid.New() // Player controller
	char.Position = position.Position{X: 0, Y: 0, Z: 3}
	char.RepsertPropertyValue(property.TeamID, 1)
	r.AddEntity(char)
	r.GameState.Turner.RemoveEntity(char.ID)
	char.CurrentDelay = 2000
	r.GameState.Entities[char.ID] = char
	r.GameState.Turner.AddEntity(char.ID, char.CurrentDelay)

	// Create an aggressive monster
	monster := entity.New()
	monster.ID = uuid.New()
	monster.Type = entity.Monster
	// Automation is driven by AIBehavior alone (ISS-101); even an
	// automated monster must have an owning controller — uuid.Nil is
	// rejected by AddEntity.
	monster.ControllerID = uuid.New()
	monster.Position = position.Position{X: 2, Y: 0, Z: 3}
	monster.RepsertPropertyValue(property.TeamID, 2)
	monster.RepsertPropertyValue(property.AIBehavior, "aggressive")
	r.AddEntity(monster)
	r.GameState.Turner.RemoveEntity(monster.ID)
	monster.CurrentDelay = 0
	r.GameState.Entities[monster.ID] = monster
	r.GameState.Turner.AddEntity(monster.ID, monster.CurrentDelay)

	r.Start()
	defer r.Stop()

	// Add the player controller
	replyChan := make(chan *message.Message, 1)
	addReplyChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.AddController{
		ControllerID: char.ControllerID,
		Controller:   &fakeController{replyChan: replyChan},
	}, &rulermethods.AddControllerReply{}), addReplyChan)

	// Wait for AddControllerReply
	select {
	case <-addReplyChan:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for AddControllerReply")
	}

	// Wait for BattleStart
	for {
		select {
		case msg := <-replyChan:
			if _, ok := msg.TargetMethod.(rulermethods.BattleStart); ok {
				goto started
			}
			t.Logf("Skipping message %T", msg.TargetMethod)
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for BattleStart")
		}
	}

started:
	// Drain replyChan in the background from here on. Past this point nothing
	// else reads it, and it's buffer-1: a single undrained notification (e.g.
	// the aggressive monster's move) permanently blocks the Ruler's actor
	// goroutine on the send, wedging the whole actor.
	go func() {
		for range replyChan {
		}
	}()

	// Signal readiness
	r.NotifyActor(message.Create(nil, rulermethods.ControllerBattleReady{
		ControllerID: char.ControllerID,
	}, nil))

	// Read state back through the actor's message queue (testingFetchEntities),
	// not r.GameState directly — the Ruler owns GameState once Start() has
	// been called (see domain_ruler_state.atom.md), and direct reads race
	// with the actor loop under -race.
	_, turnState := testingFetchEntities(r)
	t.Logf("Turner state: %s", turnState.String())

	// Wait for the first turn
	time.Sleep(1 * time.Second) // Give it time to execute behavior

	entities, _ := testingFetchEntities(r)
	var ent entity.Entity
	for _, e := range entities {
		if e.ID == monster.ID {
			ent = e
			break
		}
	}
	t.Logf("Monster position: %v", ent.Position)
	if ent.Position.X != 1 {
		t.Errorf("Monster should have moved toward player (X:0), expected X:1, got %d", ent.Position.X)
	}
}

type fakeController struct {
	replyChan chan *message.Message
}

// NotifyActor is a mock implementation for the actor interface.
func (f *fakeController) NotifyActor(msg *message.Message) {
	f.replyChan <- msg
}

// SendActor is a mock implementation for the actor interface.
func (f *fakeController) SendActor(msg *message.Message, callback chan *message.Message) {
	f.replyChan <- msg
}
