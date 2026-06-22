package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// TestRuleSkillFailMvt ensures that a skill fails if the entity does not have enough Movement points.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill Movement cost higher than the entity's current Movement points.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 11 /*Default Mvt is 3*/, property.FriendlyController, property.Skill)
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when Movement points are insufficient.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.mvt" {
		t.Errorf("Expected error 'skill.cost.mvt', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailHP ensures that a skill fails if the entity does not have enough Health Points (HP) to pay the cost.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill HP cost higher than the entity's current HP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 11 /*Default HP is 10*/, property.FriendlyController, property.Skill)
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when HP is insufficient.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.hp" {
		t.Errorf("Expected error 'skill.cost.hp', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillDeduceMvt verifies that using a skill correctly deducts the Movement point cost from the entity.
// @test-link [[mech_skill_validation]]
func TestRuleSkillDeduceMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a Movement cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 1, property.FriendlyController, property.Skill)
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

	// Verify that Movement was deducted in the GameState (from 3 to 2).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.Movement)
	if prop.GetValue() != 2 { 
		t.Errorf("Expected Movement(in GameState) to be 2, got %d", prop.GetValue())
	}

	// Verify that Movement was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.Movement)
	if prop.GetValue() != 2 { 
		t.Errorf("Expected Movement to be 2, got %d", prop.GetValue())
	}
}

// TestRuleSkillDeduceHP verifies that using a skill correctly deducts the Health Point (HP) cost from the entity.
// @test-link [[mech_skill_validation]]
func TestRuleSkillDeduceHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an HP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 1, property.FriendlyController, property.Skill)
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

	// Verify that HP was deducted in the GameState (from 10 to 9).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.HP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected HP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// Verify that HP was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.HP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected HP to be 9, got %d", prop.GetValue())
	}
}
