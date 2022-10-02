package rules

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/cell"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type FakeStateAttack struct {
	FakeState

	Attacker             uuid.UUID
	Foe                  uuid.UUID
	AttackerControllerID uuid.UUID
	FoeControllerID      uuid.UUID
	FoePosition          position.Position
	AttackerPosition     position.Position
}

func makeGameStateForTwoAttack() (*GameState, FakeStateAttack) {
	gs, f := makeGameStateForTwo()
	fake := FakeStateAttack{
		FakeState:            f,
		Attacker:             f.Entity1,
		Foe:                  f.Entity3,
		AttackerControllerID: f.Controller1,
		FoeControllerID:      f.Controller2,
		FoePosition:          position.Position{X: 5, Y: 6, Z: 3},
		AttackerPosition:     position.Position{X: 5, Y: 5, Z: 3},
	}

	gs.Turner.ForceTurn(fake.Attacker)

	// the attacker
	ent1 := gs.Entities[fake.Attacker]
	ent1.Position = fake.AttackerPosition
	ent1.CurrentDelay = 0 // it's his turn.
	ent1.FaceToward(fake.FoePosition)
	// the foe
	ent3 := gs.Entities[fake.Foe]
	ent3.Position = fake.FoePosition
	ent3.FaceToward(fake.AttackerPosition)

	gs.Entities[fake.Attacker] = ent1
	gs.Entities[fake.Foe] = ent3

	gs.Grid.MoveEntity(position.New(0, 0, 3), fake.AttackerPosition, fake.Attacker)
	gs.Grid.MoveEntity(position.New(9, 9, 3), fake.FoePosition, fake.Foe)

	return gs, fake
}

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

func TestRuleAttackSucceed(t *testing.T) {

	gs, fake := makeGameStateForTwoAttack()
	gs.Turner.ForceTurn(fake.Attacker)

	// Attack Entity 3
	msg := message.Create(nil,
		rulermethods.ControllerAttack{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
		}, nil)

	reply := gs.Attack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))

	if reply.HasError {
		t.Errorf("Expected no error, got '%s'", reply.ErrorKey)
	}

	attacker := gs.Entities[fake.Attacker]
	foe := gs.Entities[fake.Foe]

	// Attacker has Acted.
	if !attacker.HasActed() {
		t.Errorf("Expected Attacker to have acted(in GameState).")
	}

	if !attacker.HasMoved() {
		t.Errorf("Expected Attacker to be prevented from moving again(in GameState).")
	}

	// Foe has been damaged.
	foehp := foe.GetPropertyC(property.HP)
	attackerattack := attacker.GetPropertyI(property.Attack)
	foedefense := foe.GetPropertyI(property.Defense)

	if foehp.GetValue() != foehp.GetMaxValue()-(attackerattack.I()-foedefense.I()) {
		t.Errorf("Expected Foe to have taken damage(in GameState).")
	}

	// attacker has been delayed

	if attacker.CurrentDelay != 500 {
		t.Errorf("Expected Attacker to be delayed(in GameState).")
	}

	// Check reply entity's to have the same values as game state.
	content := reply.Content.(rulermethods.ControllerAttackReply)

	attacker = content.Entity

	// Attacker has Acted.
	if !attacker.HasActed() {
		t.Errorf("Expected Attacker to have acted.")
	}

	if !attacker.HasMoved() {
		t.Errorf("Expected Attacker to be prevented from moving again.")
	}

	// attacker has been delayed

	if attacker.CurrentDelay != 500 {
		t.Errorf("Expected Attacker to be delayed.")
	}

	// expect foe's inbox to have been filled with a notification of the attack.
	if len(fake.FakeController2.NotifyMessages) != 1 {
		t.Errorf("Expected Foe to have received a notification of the attack.")
	} else {
		notif := fake.FakeController2.NotifyMessages[0].TargetMethod.(rulermethods.ControllerAttacked)

		if notif.AttackerControllerID != fake.AttackerControllerID {
			t.Errorf("Expected AttackerControllerID to be %d, got %d", fake.AttackerControllerID, notif.AttackerControllerID)
		}
		if notif.Attacker.ID != fake.Attacker {
			t.Errorf("Expected AttackerID to be %d, got %d", fake.Attacker, notif.Attacker.ID)
		}
		if notif.ControllerID != fake.FoeControllerID {
			t.Errorf("Expected ControllerID to be %d, got %d", fake.FoeControllerID, notif.ControllerID)
		}
		if notif.Entity.ID != fake.Foe {
			t.Errorf("Expected EntityID to be %d, got %d", fake.Foe, notif.Entity.ID)
		}

		// expect entity in notification to have it's HP deduced.
		if notif.Entity.GetPropertyC(property.HP).GetValue() != foe.GetPropertyC(property.HP).GetValue() {
			t.Errorf("Expected Entity in notification to have it's HP deduced.")
		}
	}

}
