package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// TestRuleSkillFailOutOfTurn ensures that a skill usage fails if it is not the entity's turn.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailOutOfTurn(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	// Force the turn to another entity.
	gs.Turner.ForceTurn(fake.Entity2)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when it is not the attacker's turn.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.turn.mismatch" {
		t.Errorf("Expected error 'entity.turn.mismatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailWrongController ensures that a skill usage fails if the controller ID does not match the entity's controller.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailWrongController(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.FoeControllerID, // Wrong controller!
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill with an authorized controller.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.controller.mismatch" {
		t.Errorf("Expected error 'entity.controller.mismatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailUnknownEntity ensures that a skill usage fails if the entity ID is invalid.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailUnknownEntity(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     uuid.New(), // Random UUID
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill for a non-existent entity.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.notfound" {
		t.Errorf("Expected error 'entity.notfound', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailUnknownSkill ensures that a skill usage fails if the skill ID is invalid for the entity.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailUnknownSkill(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      uuid.New(), // Random UUID
		}, nil)

	// Attempt to use a skill that the entity does not possess.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.notfound" {
		t.Errorf("Expected error 'skill.notfound', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailAlreadyActed ensures that an entity cannot use a skill if it has already acted this turn.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailAlreadyActed(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	ent := gs.Entities[fake.Attacker]
	// Manually set the HasActed property to true.
	prop := ent.GetProperty(property.HasActed)
	prop.Set(true)
	ent.UpdateProperty(prop)
	gs.Entities[fake.Attacker] = ent

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when the entity is already in 'acted' state.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.alreadyacted" {
		t.Errorf("Expected error 'entity.alreadyacted', got '%s'", reply.ErrorKey)
	}
}
