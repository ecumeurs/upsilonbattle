package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// @test-link [[mech_combat_standard_attack_computation]]
// @test-link [[rule_friendly_fire]]

type FakeStateAttack struct {
	FakeState

	Attacker             uuid.UUID
	Foe                  uuid.UUID
	AttackerControllerID uuid.UUID
	FoeControllerID      uuid.UUID
	FoePosition          position.Position
	AttackerPosition     position.Position
}

// makeGameStateForTwoAttack initializes a GameState with two entities (attacker and foe)
// positioned for combat testing. It sets the current turn to the attacker.
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


// TestRuleAttackSucceed verifies that a valid attack correctly applies damage,
// updates action states, and notifies relevant controllers.
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
	if attacker.CurrentDelay != 100 {
		t.Errorf("Expected Attacker to be delayed(in GameState).")
	}

	// Check reply entity's to have the same values as game state.
	content := reply.Content.(rulermethods.ControllerAttackReply)

	attacker = content.Attacker

	// Attacker has Acted.
	if !attacker.HasActed() {
		t.Errorf("Expected Attacker to have acted.")
	}

	if !attacker.HasMoved() {
		t.Errorf("Expected Attacker to be prevented from moving again.")
	}

	// attacker has been delayed
	if attacker.CurrentDelay != 100 {
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


// TestRuleAttackCredits verifies that credits are awarded correctly after a successful attack.
func TestRuleAttackCredits(t *testing.T) {
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

	attackerattack := attacker.GetPropertyI(property.Attack).I()
	foedefense := foe.GetPropertyI(property.Defense).I()
	expectedDamage := attackerattack - foedefense

	// expect foe's inbox to have been filled with a notification of the attack.
	if len(fake.FakeController2.NotifyMessages) != 1 {
		t.Errorf("Expected Foe to have received a notification of the attack.")
	} else {
		notif := fake.FakeController2.NotifyMessages[0].TargetMethod.(rulermethods.ControllerAttacked)

		if len(notif.CreditAwards) != 1 {
			t.Errorf("Expected 1 credit award, got %d", len(notif.CreditAwards))
		} else {
			award := notif.CreditAwards[0]
			if award.PlayerID != fake.AttackerControllerID {
				t.Errorf("Expected credit for player %s, got %s", fake.AttackerControllerID, award.PlayerID)
			}
			if award.Amount != expectedDamage {
				t.Errorf("Expected %d credits, got %d", expectedDamage, award.Amount)
			}
			if award.Source != "damage" {
				t.Errorf("Expected source 'damage', got '%s'", award.Source)
			}
		}
	}
}
