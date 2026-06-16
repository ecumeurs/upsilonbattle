package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// @test-link [[mech_move_validation]]
// @test-link [[mech_move_validation_move_validation_already_moved]]

// TestRuleMoveFailOutOfTurn verifies that movement fails if it's not the entity's turn.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.turn.mismatch" {
		t.Errorf("Expected error 'entity.turn.mismatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailWrongController verifies that movement fails if the controller ID is incorrect.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.controller.mismatch" {
		t.Errorf("Expected error 'entity.controller.mismatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailUnknownEntity verifies that movement fails for an unknown entity ID.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.notfound" {
		t.Errorf("Expected error 'entity.notfound', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailOccupied verifies that movement fails if a cell in the path is occupied.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

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

// TestRuleMoveFailObstacle verifies that movement fails if a cell in the path is an obstacle.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

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

// TestRuleMoveFailMovementPropertyNotEnoughForPath verifies failure when the path exceeds movement range.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.too.long" {
		t.Errorf("Expected error 'entity.path.too.long', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailNotAdjascent verifies that movement fails if the path steps are not adjacent.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.notvalid" {
		t.Errorf("Expected error 'entity.path.notvalid', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailStartNotAdjascent verifies that movement fails if the first step is not adjacent to the start.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.notadjacent" {
		t.Errorf("Expected error 'entity.path.notadjacent', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailNoMovementCredits verifies movement fails if the entity has zero movement credits.
func TestRuleMoveFailNoMovementCredits(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)
	ent1 := gs.Entities[fake.Entity1]
	p := ent1.GetPropertyC(property.Movement)
	p.SetValue(0)
	ent1.UpdateProperty(p)
	gs.Entities[fake.Entity1] = ent1

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

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.movement.nocredits" {
		t.Errorf("Expected error 'entity.movement.nocredits', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailNotEnoughMovementCredits verifies movement fails if the path is longer than remaining credits.
func TestRuleMoveFailNotEnoughMovementCredits(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)
	ent1 := gs.Entities[fake.Entity1]
	p := ent1.GetPropertyC(property.Movement)
	p.SetValue(2)
	ent1.UpdateProperty(p)
	gs.Entities[fake.Entity1] = ent1

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

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.movement.credits" {
		t.Errorf("Expected error 'entity.movement.credits', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailNotAdjascentJumpHeight verifies movement fails if the vertical difference exceeds jump height.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.path.notvalid" {
		t.Errorf("Expected error 'entity.path.notvalid', got '%s'", reply.ErrorKey)
	}
}

// TestRuleMoveFailAlreadyMoved verifies that movement fails if the entity has already moved this turn.
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

	reply := Move(gs, msg, msg.TargetMethod.(rulermethods.ControllerMove))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.movement.already" {
		t.Errorf("Expected error 'entity.movement.already', got '%s'", reply.ErrorKey)
	}
}
