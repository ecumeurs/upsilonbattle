package rules

import (
	"reflect"
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// @test-link [[mech_move_validation]]
// @test-link [[mech_move_validation]]



// TestRuleMoveSucceed verifies that a valid movement path correctly updates
// the entity's position, reduces movement credits, and updates the grid.
func TestRuleMoveSucceed(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	// Move entity 1
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{
				{X: 0, Y: 1, Z: 3},
				{X: 0, Y: 2, Z: 3},
				{X: 0, Y: 3, Z: 3},
			},
		}, nil)

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if reply.HasError {
		t.Errorf("Expected no error, got '%s'", reply.ErrorKey)
	}

	if reflect.TypeOf(reply.Content) != reflect.TypeOf(rulermethods.ControllerMoveReply{}) {
		t.Errorf("Expected reply's content to be of type ControllerMoveReply, got '%s'", reflect.TypeOf(reply.Content))
	}

	content := reply.Content.(rulermethods.ControllerMoveReply)

	// test reply's entity position.
	if content.Entity.ID != fake.Entity1 {
		t.Errorf("Expected entity id '%s', got '%s'", fake.Entity1, reply.Content.(rulermethods.ControllerMoveReply).Entity.ID)
	}

	endpos := msg.TargetMethod.(rulermethods.ControllerMove).Path[len(msg.TargetMethod.(rulermethods.ControllerMove).Path)-1]
	if !content.Entity.Position.Equals(endpos) {
		t.Errorf("Expected entity position '%s', got '%s'", endpos, content.Entity.Position)
	}

	// Checks position on grid
	c, _ := gs.Grid.CellAt(content.Entity.Position)
	if !c.HasEntity(fake.Entity1) {
		t.Errorf("Expected entity id '%s' at position '%s'", fake.Entity1, content.Entity.Position)
	}

	// Checks entity's position in gamestate
	e := gs.Entities[fake.Entity1]
	if !e.Position.Equals(content.Entity.Position) {
		t.Errorf("Expected entity position '%s', got '%s'", content.Entity.Position, e.Position)
	}

	// expect entity's movement credit to have been reduced by len(path)
	prop := content.Entity.GetPropertyC(property.Movement)
	if prop == nil {
		t.Errorf("Expected entity to have movement property, got none.")
	}

	if prop.GetValue() != prop.GetMaxValue()-len(msg.TargetMethod.(rulermethods.ControllerMove).Path) {
		t.Errorf("Expected entity's movement to be %d, got %d", prop.GetMaxValue()-len(msg.TargetMethod.(rulermethods.ControllerMove).Path), prop.GetValue())
	}

	// check's prop has been updated in game state as well.
	prop = gs.Entities[fake.Entity1].GetPropertyC(property.Movement)
	if prop == nil {
		t.Errorf("Expected entity to have movement property(in GameState), got none.")
	}

	if prop.GetValue() != prop.GetMaxValue()-len(msg.TargetMethod.(rulermethods.ControllerMove).Path) {
		t.Errorf("Expected entity's movement(in GameState) to be %d, got %d", prop.GetMaxValue()-len(msg.TargetMethod.(rulermethods.ControllerMove).Path), prop.GetValue())
	}
}
