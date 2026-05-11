package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// TestRuleSkillFailMP ensures that a skill fails if the entity does not have enough Mana Points (MP).
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill MP cost higher than the entity's current MP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 11 /*Default MP is 10*/, property.FriendlyController, property.Skill)
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when MP is insufficient.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.mp" {
		t.Errorf("Expected error 'skill.cost.mp', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailSP ensures that a skill fails if the entity does not have enough Stamina Points (SP).
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailSP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill SP cost higher than the entity's current SP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.SPLeech, 11 /*Default HP is 10*/, property.FriendlyController, property.Skill)
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when SP is insufficient.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.sp" {
		t.Errorf("Expected error 'skill.cost.sp', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillDeduceMP verifies that using a skill correctly deducts the Mana Point (MP) cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an MP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 1, property.FriendlyController, property.Skill)
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// Verify that MP was deducted in the GameState (from 10 to 9).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.MP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected MP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// Verify that MP was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.MP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected MP to be 9, got %d", prop.GetValue())
	}
}

// TestRuleSkillDeduceSP verifies that using a skill correctly deducts the Stamina Point (SP) cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceSP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an SP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.SPLeech, 1, property.FriendlyController, property.Skill)
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// Verify that SP was deducted in the GameState (from 10 to 9).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.SP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected SP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// Verify that SP was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.SP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected SP to be 9, got %d", prop.GetValue())
	}
}
