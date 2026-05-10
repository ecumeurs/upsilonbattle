package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// TestRuleSkillFailCooldown ensures that a skill fails if it is currently on cooldown.
// @test-link [[mech_skill_validation_economic_cost_verification_cooldown_check]]
func TestRuleSkillFailCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set the skill cooldown to 1.
	fake.Skill.Cooldown = 1
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill that is on cooldown.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cooldown" {
		t.Errorf("Expected error 'skill.cooldown', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailMP ensures that a skill fails if the entity does not have enough Mana Points (MP).
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill MP cost higher than the entity's current MP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 11 /*Default MP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when MP is insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when SP is insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.sp" {
		t.Errorf("Expected error 'skill.cost.sp', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailMvt ensures that a skill fails if the entity does not have enough Movement points.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill Movement cost higher than the entity's current Movement points.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 11 /*Default Mvt is 3*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when Movement points are insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.mvt" {
		t.Errorf("Expected error 'skill.cost.mvt', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailHP ensures that a skill fails if the entity does not have enough Health Points (HP) to pay the cost.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill HP cost higher than the entity's current HP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 11 /*Default HP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when HP is insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.hp" {
		t.Errorf("Expected error 'skill.cost.hp', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillSetCooldown verifies that using a skill correctly applies its cooldown to the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_cooldown_check]]
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
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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

// TestRuleSkillDeduceMP verifies that using a skill correctly deducts the Mana Point (MP) cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an MP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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

// TestRuleSkillDeduceMvt verifies that using a skill correctly deducts the Movement point cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a Movement cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an HP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
