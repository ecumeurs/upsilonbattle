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

// TestAggressiveControllerStall verifies that the AI controller finds a path around
// blockers rather than stalling or moving through allied entities.
func TestAggressiveControllerStall(t *testing.T) {
	ruler := NewMockRuler()
	defer ruler.Stop()

	aiID := uuid.New()
	ctl := NewAggressiveController(aiID, "TestAI")
	ctl.ruler = ruler
	ctl.Start()
	defer ctl.Stop()

	ctl.Grid = grid.NewGrid(5, 5, 1)

	aiEnt := makeTestEntity(aiID, position.Position{X: 0, Y: 0, Z: 1}, 1)
	aiEnt.RepsertPropertyValue(property.AttackRange, 1)
	foeEnt := makeTestEntity(uuid.New(), position.Position{X: 2, Y: 0, Z: 1}, 2)
	blockEnt := makeTestEntity(aiID, position.Position{X: 1, Y: 0, Z: 1}, 1) // ally blocker

	ctl.KnownEntities[aiEnt.ID] = aiEnt
	ctl.KnownEntities[foeEnt.ID] = foeEnt
	ctl.KnownEntities[blockEnt.ID] = blockEnt

	ctl.NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{Entity: aiEnt}, nil))

	select {
	case msg := <-ruler.ReceivedMessages:
		assertStallAction(t, msg, aiEnt)
	case <-time.After(5 * time.Second):
		t.Error("Timeout: AI didn't send any action to Ruler")
	}
}

func makeTestEntity(controllerID uuid.UUID, pos position.Position, teamID int) entity.Entity {
	e := entity.Entity{ID: uuid.New(), ControllerID: controllerID, Position: pos}
	e.Properties = make(map[string]property.Property)
	e.RepsertPropertyValue(property.TeamID, teamID)
	return e
}

func assertStallAction(t *testing.T, msg *message.Message, aiEnt entity.Entity) {
	t.Helper()
	switch m := msg.TargetMethod.(type) {
	case rulermethods.ControllerMove:
		t.Logf("AI correctly attempted to move: %v", m.Path)
		for _, p := range m.Path {
			if p.X == 1 && p.Y == 0 {
				t.Error("AI moved THROUGH blocker at (1,0)")
			}
		}
	case rulermethods.ControllerAttack:
		t.Errorf("AI attacked from out-of-range position: target=%v ai=%v", m.Target, aiEnt.Position)
	case rulermethods.EndOfTurn:
		t.Error("AI passed turn immediately without moving or attacking")
	default:
		t.Errorf("AI sent unexpected message: %T", m)
	}
}
