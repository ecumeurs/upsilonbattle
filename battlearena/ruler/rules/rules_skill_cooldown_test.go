package rules

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// TestRuleSkillFailCooldown ensures that a skill fails if it is currently on cooldown.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set the skill cooldown to 1.
	fake.Skill.Cooldown = 1
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill that is on cooldown.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cooldown" {
		t.Errorf("Expected error 'skill.cooldown', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillSetCooldown verifies that using a skill correctly applies its cooldown to the entity.
// @test-link [[mech_skill_validation]]
func TestRuleSkillSetCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Default cooldown cost for the test skill is 3.
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Expected no error, got '%s'", reply.ErrorKey)
	}

	// Verify that the cooldown was updated in the GameState.
	skill := gs.Entities[fake.Attacker].Skills[fake.SkillID]
	if skill.Cooldown != 3 {
		t.Errorf("Expected cooldown to be 3(in GameState), got %d", skill.Cooldown)
	}

	// Verify that the cooldown was updated in the reply message.
	skill = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.Skills[fake.SkillID]
	if skill.Cooldown != 3 {
		t.Errorf("Expected cooldown to be 3, got %d", skill.Cooldown)
	}
}
