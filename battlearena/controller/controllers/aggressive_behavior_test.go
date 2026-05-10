package controllers
// @test-link [[mechanic_mech_arena_lifecycle]]
// @test-link [[rule_team_mechanics]]
// @test-link [[uc_combat_turn]]

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type MockRuler struct {
	*actor.Actor
	ReceivedMessages chan *message.Message
}

// NewMockRuler creates a new MockRuler for testing purposes.
func NewMockRuler() *MockRuler {
	m := &MockRuler{
		Actor:            actor.New("MockRuler"),
		ReceivedMessages: make(chan *message.Message, 10),
	}
	m.Start()
	return m
}

// SendActor is a mock implementation for sending messages to the ruler actor.
func (m *MockRuler) SendActor(msg *message.Message, replyChan chan *message.Message) {
	m.ReceivedMessages <- msg
}

// NotifyActor is a mock implementation for notifying the ruler actor.
func (m *MockRuler) NotifyActor(msg *message.Message) {
	m.ReceivedMessages <- msg
}

// TestAggressiveControllerStall verifies that the AI controller correctly finds a path around blockers instead of stalling.
func TestAggressiveControllerStall(t *testing.T) {
	// Setup environment
	ruler := NewMockRuler()
	defer ruler.Stop()

	aiID := uuid.New()
	ctl := NewAggressiveController(aiID, "TestAI")
	ctl.ruler = ruler
	ctl.Start()
	defer ctl.Stop()

	grd := grid.NewGrid(5, 5, 1)
	ctl.Grid = grd

	// Entities
	aiEnt := entity.Entity{
		ID:           uuid.New(),
		ControllerID: aiID,
		Position:     position.Position{X: 0, Y: 0, Z: 1},
	}
	aiEnt.Properties = make(map[string]property.Property)
	aiEnt.RepsertPropertyValue(property.TeamID, 1)
	aiEnt.RepsertPropertyValue(property.AttackRange, 1)

	foeEnt := entity.Entity{
		ID:           uuid.New(),
		ControllerID: uuid.New(),
		Position:     position.Position{X: 2, Y: 0, Z: 1},
	}
	foeEnt.Properties = make(map[string]property.Property)
	foeEnt.RepsertPropertyValue(property.TeamID, 2)

	blockEnt := entity.Entity{
		ID:           uuid.New(),
		ControllerID: aiID, // Ally blocker
		Position:     position.Position{X: 1, Y: 0, Z: 1},
	}
	blockEnt.Properties = make(map[string]property.Property)
	blockEnt.RepsertPropertyValue(property.TeamID, 1)

	// Populate KnownEntities
	ctl.KnownEntities[aiEnt.ID] = aiEnt
	ctl.KnownEntities[foeEnt.ID] = foeEnt
	ctl.KnownEntities[blockEnt.ID] = blockEnt

	// Trigger Turn
	ctl.NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
		Entity: aiEnt,
	}, nil))

	// Expectation: AI should attempt to find a path AROUND the blocker.
	// In a 5x5 grid, it should find plenty of space.
	
	select {
	case msg := <-ruler.ReceivedMessages:
		switch m := msg.TargetMethod.(type) {
		case rulermethods.ControllerMove:
			t.Logf("AI correctly attempted to move: %v", m.Path)
			// Verify that the path does NOT contain (1,0)
			for _, p := range m.Path {
				if p.X == 1 && p.Y == 0 {
					t.Errorf("AI attempted to move THROUGH blocker at (1,0)!")
				}
			}
		case rulermethods.ControllerAttack:
			t.Errorf("AI attempted to attack from current position while target is out of range! Target: %v, AI: %v", m.Target, aiEnt.Position)
		case rulermethods.EndOfTurn:
			t.Errorf("AI passed turn immediately without moving or attacking")
		default:
			t.Errorf("AI sent unexpected message: %T", m)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout: AI didn't send any action to Ruler")
	}
}
