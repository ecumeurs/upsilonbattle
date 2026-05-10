package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// @test-link [[mech_skill_validation]]
// @test-link [[mechanic_mech_effect_damage]]
// @test-link [[mechanic_mech_effect_heal]]
// @test-link [[mechanic_mech_effect_poison]]
// @test-link [[mechanic_mech_effect_shield]]

// TestRuleSkillEffectDamage verifies that a skill with a damage effect correctly
// reduces target HP and handles entity death.
func TestRuleSkillEffectDamage(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Effect.Properties = append(fake.Skill.Effect.Properties,
		defaultproperty.MakeIntProperty(property.Damage, 400, property.Public, property.Skill))
	
	// Add 100% accuracy to ensure it hits
	fake.Skill.Targeting[property.Accuracy.String()] = defaultproperty.MakeIntProperty(property.Accuracy, 100, property.Public, property.Skill)
	
	gs.addSkillToEntity(fake.Attacker, fake.Skill)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Fatalf("Expected no error, got '%s'", reply.ErrorKey)
	}

	// Check foe is removed from GameState
	if foe, exists := gs.Entities[fake.Foe]; exists {
		hp := foe.GetPropertyC(property.HP)
		t.Errorf("Expected foe to be removed from GameState due to death, but it exists with HP: %d", hp.GetValue())
	}
}

// TestRuleSkillEffectHeal verifies that a healing skill correctly restores HP to the target.
func TestRuleSkillEffectHeal(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Damage the attacker first so we can heal them
	attacker := gs.Entities[fake.Attacker]
	hp := attacker.GetPropertyC(property.HP)
	hp.SetValue(5) // Max is 10
	attacker.UpdateProperty(hp)
	gs.Entities[fake.Attacker] = attacker

	// Add Heal 10 to skill, target Self
	fake.Skill.Effect.Properties = append(fake.Skill.Effect.Properties,
		defaultproperty.MakeIntProperty(property.Heal, 10, property.Public, property.Skill))
	
	targetProp := defaultproperty.MakeValidatedStringProperty(property.TargetType, "Self", property.Public, property.Skill, []string{"Self"})
	fake.Skill.Targeting[property.TargetType.String()] = targetProp
	
	rng := defaultproperty.MakeIntCounterProperty(property.Range, 0, 1, property.Public, property.Skill)
	fake.Skill.Targeting[property.Range.String()] = rng
	
	gs.addSkillToEntity(fake.Attacker, fake.Skill)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.AttackerPosition, // Self
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Fatalf("Expected no error, got '%s'", reply.ErrorKey)
	}

	// Check attacker HP
	attackerAfter := gs.Entities[fake.Attacker]
	hpAfter := attackerAfter.GetPropertyC(property.HP)
	if hpAfter.GetValue() != 10 { // Capped at max
		t.Errorf("Expected attacker HP to be 10, got %d", hpAfter.GetValue())
	}
}

// TestRuleSkillEffectPoisonCounter verifies that a poison skill correctly applies
// the poison property to the target.
func TestRuleSkillEffectPoisonCounter(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Add Poison 5, 100% chance
	fake.Skill.Effect.Properties = append(fake.Skill.Effect.Properties,
		defaultproperty.MakeIntProperty(property.PoisonPower, 5, property.Public, property.Skill))
	fake.Skill.Effect.Properties = append(fake.Skill.Effect.Properties,
		defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.Public, property.Skill))
	// Add 100% accuracy to ensure it hits
	fake.Skill.Targeting[property.Accuracy.String()] = defaultproperty.MakeIntProperty(property.Accuracy, 100, property.Public, property.Skill)
	
	gs.addSkillToEntity(fake.Attacker, fake.Skill)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Fatalf("Expected no error, got '%s'", reply.ErrorKey)
	}

	// Foe should have Poison property increased
	foe := gs.Entities[fake.Foe]
	poison := foe.GetPropertyI(property.Poison)
	
	if poison.I() != 5 {
		t.Errorf("Expected Poison property counter to be 5, got %d", poison.I())
	}
}

// TestRuleSkillEffectShield verifies that a skill correctly applies a shield property to the target.
func TestRuleSkillEffectShield(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Add Shield 10 to skill, target Self
	fake.Skill.Effect.Properties = append(fake.Skill.Effect.Properties,
		defaultproperty.MakeIntProperty(property.ShieldPower, 10, property.Public, property.Skill))
	// Add Heal 1 (needs to be considered a healing skill to process positive shield)
	fake.Skill.Effect.Properties = append(fake.Skill.Effect.Properties,
		defaultproperty.MakeIntProperty(property.Heal, 1, property.Public, property.Skill))

	targetProp := defaultproperty.MakeValidatedStringProperty(property.TargetType, "Self", property.Public, property.Skill, []string{"Self"})
	fake.Skill.Targeting[property.TargetType.String()] = targetProp
	rng := defaultproperty.MakeIntCounterProperty(property.Range, 0, 1, property.Public, property.Skill)
	fake.Skill.Targeting[property.Range.String()] = rng
	
	gs.addSkillToEntity(fake.Attacker, fake.Skill)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.AttackerPosition, // Self
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Fatalf("Expected no error, got '%s'", reply.ErrorKey)
	}

	// Check attacker Shield
	attackerAfter := gs.Entities[fake.Attacker]
	shieldAfter := attackerAfter.GetPropertyC(property.Shield)
	if shieldAfter.GetValue() != 10 {
		t.Errorf("Expected attacker Shield to be 10, got %d", shieldAfter.GetValue())
	}
}
