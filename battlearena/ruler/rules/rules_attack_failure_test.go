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

// @test-link [[mech_combat_standard_attack_computation]]
// @test-link [[rule_friendly_fire]]

// TestRuleAttackFailOutOfTurn verifies that an attack fails if it's not the entity's turn.
func TestRuleAttackFailOutOfTurn(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()
	gs.Turner.ForceTurn(fake.Entity2)

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.turn.missmatch" {
		t.Errorf("Expected error 'entity.turn.missmatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailWrongController verifies that an attack fails if the controller ID doesn't match.
func TestRuleAttackFailWrongController(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.FoeControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.controller.missmatch" {
		t.Errorf("Expected error 'entity.controller.missmatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailUnknownEntity verifies that an attack fails if the entity ID is unknown.
func TestRuleAttackFailUnknownEntity(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     uuid.New(),
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.notfound" {
		t.Errorf("Expected error 'entity.notfound', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailTargetNotFound verifies that an attack fails if the target position is out of bounds.
func TestRuleAttackFailTargetNotFound(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(11, 11, 3), // board is 10x10
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.attack.target.invalid" {
		t.Errorf("Expected error 'entity.attack.target.invalid', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailTargetNotGround verifies that an attack fails if the target cell type is invalid (e.g. Water).
func TestRuleAttackFailTargetNotGround(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()
	gs.Grid.ReplaceCellType(fake.FoePosition, cell.Water)

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.attack.celltype" {
		t.Errorf("Expected error 'entity.attack.celltype', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailTargetNotEntity verifies that an attack fails if there is no entity at the target position.
func TestRuleAttackFailTargetNotEntity(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()
	gs.Grid.MoveEntity(fake.FoePosition, position.New(0, 0, 3), fake.Foe)

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.attack.noentity" {
		t.Errorf("Expected error 'entity.attack.noentity', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailTargetNotInRange verifies that an attack fails if the target is beyond the attacker's range.
func TestRuleAttackFailTargetNotInRange(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()
	foeTempPosition := fake.FoePosition.Add(position.New(1, 0, 0)) // default range for basic attack is 1
	gs.Grid.MoveEntity(fake.FoePosition, foeTempPosition, fake.Foe)

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       foeTempPosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.attack.outofrange" {
		t.Errorf("Expected error 'entity.attack.outofrange', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailAlreadyActed verifies that an attack fails if the entity has already performed an action this turn.
func TestRuleAttackFailAlreadyActed(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()

	ent := gs.Entities[fake.Attacker]
	// this flag is set when the entity has acted or has been interrupted during movement by reaction or traps.
	prop := ent.GetProperty(property.HasActed)
	prop.Set(true)
	ent.UpdateProperty(prop)
	gs.Entities[fake.Attacker] = ent

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.hasacted" {
		t.Errorf("Expected error 'entity.hasacted', got '%s'", reply.ErrorKey)
	}
}

// TestRuleAttackFailFriendlyFire verifies that an attack fails if the target is on the same team.
func TestRuleAttackFailFriendlyFire(t *testing.T) {
	gs, fake := makeGameStateForTwoAttack()

	// Set both entities on the same team
	ent1 := gs.Entities[fake.Attacker]
	ent3 := gs.Entities[fake.Foe]

	prop1 := ent1.GetProperty(property.TeamID)
	prop1.Set(1)
	ent1.UpdateProperty(prop1)

	prop3 := ent3.GetProperty(property.TeamID)
	prop3.Set(1)
	ent3.UpdateProperty(prop3)

	gs.Entities[fake.Attacker] = ent1
	gs.Entities[fake.Foe] = ent3

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.attack.friendlyfire" {
		t.Errorf("Expected error 'entity.attack.friendlyfire', got '%s'", reply.ErrorKey)
	}
}
