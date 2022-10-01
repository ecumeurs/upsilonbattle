package rules

import (
	"reflect"
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/cell"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

func TestRuleMoveFailOutOfTurn(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity2)

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

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.turn.missmatch" {
		t.Errorf("Expected error 'entity.turn.missmatch', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailWrongController(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	// Move entity 1
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller2,
			Path: []position.Position{
				{X: 0, Y: 1, Z: 3},
				{X: 0, Y: 2, Z: 3},
				{X: 0, Y: 3, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.controller.missmatch" {
		t.Errorf("Expected error 'entity.controller.missmatch', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailUnknownEntity(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	// Move entity 1
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     uuid.New(),
			ControllerID: fake.Controller1,
			Path: []position.Position{
				{X: 0, Y: 1, Z: 3},
				{X: 0, Y: 2, Z: 3},
				{X: 0, Y: 3, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.notfound" {
		t.Errorf("Expected error 'entity.notfound', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailOccupied(t *testing.T) {

	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	// Move entity 1 to the right.
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{ // expect begin a 0,0,3 -> 3,0,3
				{X: 1, Y: 0, Z: 3}, // should fail: this one is occupied.
				{X: 2, Y: 0, Z: 3},
				{X: 3, Y: 0, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.occupied" {
		t.Errorf("Expected error 'entity.path.occupied', got '%s'", reply.ErrorKey)
	}

	// expect state of entity 1 to be unchanged.
	ent1 := gs.Entities[fake.Entity1]
	if !ent1.Position.Equals(position.Position{X: 0, Y: 0, Z: 3}) {
		t.Errorf("Expected entity 1 to be at 0,0,3, got %v", ent1.Position)
	}
}

func TestRuleMoveFailObstacle(t *testing.T) {

	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	gs.Grid.ReplaceCellType(position.Position{X: 0, Y: 1, Z: 3}, cell.Obstacle)

	// Move entity 1 to the right.
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{ // expect begin a 0,0,3 -> 3,0,3
				{X: 0, Y: 1, Z: 3}, // should fail: this one is occupied.
				{X: 0, Y: 2, Z: 3},
				{X: 0, Y: 3, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.obstacle" {
		t.Errorf("Expected error 'entity.path.obstacle', got '%s'", reply.ErrorKey)
	}

	// expect state of entity 1 to be unchanged.
	ent1 := gs.Entities[fake.Entity1]
	if !ent1.Position.Equals(position.Position{X: 0, Y: 0, Z: 3}) {
		t.Errorf("Expected entity 1 to be at 0,0,3, got %v", ent1.Position)
	}
}

func TestRuleMoveFailMovementPropertyNotEnoughForPath(t *testing.T) {
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
				{X: 0, Y: 4, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.too.long" {
		t.Errorf("Expected error 'entity.path.too.long', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailNotAdjascent(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	// Move entity 1
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{
				{X: 0, Y: 1, Z: 3},
				{X: 0, Y: 7, Z: 3},
				{X: 0, Y: 3, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.notvalid" {
		t.Errorf("Expected error 'entity.path.notvalid', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailStartNotAdjascent(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	// Move entity 1
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{
				{X: 0, Y: 5, Z: 3},
				{X: 0, Y: 6, Z: 3},
				{X: 0, Y: 7, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.notadjascent" {
		t.Errorf("Expected error 'entity.path.notadjascent', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailNoMovementCredits(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)
	p := gs.Entities[fake.Entity1].GetPropertyC(property.Movement)
	p.SetValue(0)
	gs.Entities[fake.Entity1].UpdateProperty(p)

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

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.movement.nocredits" {
		t.Errorf("Expected error 'entity.movement.nocredits', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailNotEnoughMovementCredits(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)
	p := gs.Entities[fake.Entity1].GetPropertyC(property.Movement)
	p.SetValue(2)
	gs.Entities[fake.Entity1].UpdateProperty(p)

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

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.movement.credits" {
		t.Errorf("Expected error 'entity.movement.credits', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailNotAdjascentJumpHeight(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)
	c, _ := gs.Grid.CellAt(position.Position{X: 0, Y: 2, Z: 3})
	c.Position.Z = 6 // default jump height is 2, so this should be sufficient
	gs.Grid.ReplaceCell(position.New(0, 2, 3), c)

	// Move entity 1
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{
				{X: 0, Y: 1, Z: 3},
				{X: 0, Y: 2, Z: 6},
				{X: 0, Y: 3, Z: 3},
			},
		}, nil)

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.notvalid" {
		t.Errorf("Expected error 'entity.path.notvalid', got '%s'", reply.ErrorKey)
	}
}

func TestRuleMoveFailAlreadyMoved(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	ent := gs.Entities[fake.Entity1]
	// this flag is set when the entity has acted or has been interrupted during movement by reaction or traps.
	prop := ent.GetProperty(property.HasMoved)
	prop.Set(true)
	ent.UpdateProperty(prop)
	gs.Entities[fake.Entity1] = ent

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

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.movement.already" {
		t.Errorf("Expected error 'entity.movement.already', got '%s'", reply.ErrorKey)
	}
}

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

	reply := gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

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
	if c.EntityID != fake.Entity1 {
		t.Errorf("Expected entity id '%s' at position '%s', got '%s'", fake.Entity1, content.Entity.Position, c.EntityID)
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
