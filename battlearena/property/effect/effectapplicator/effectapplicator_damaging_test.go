package effectapplicator

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/sirupsen/logrus"
)

/**
	Default properties for a character

// note: futher properties may be added per entity basis.
func PropertiesForCharacter() []property.Property {
	return []property.Property{
		defaultproperty.MakeIntCounterProperty(property.HP, 10, 10, property.Public, property.Character),
		defaultproperty.MakeIntCounterProperty(property.Movement, 3, 3, property.Public, property.Character),
		defaultproperty.MakeIntCounterProperty(property.SP, 10, 10, property.Public, property.Character),
		defaultproperty.MakeIntCounterProperty(property.MP, 10, 10, property.Public, property.Character),
		defaultproperty.MakeIntCounterProperty(property.Shield, 0, 0, property.Public, property.Character),
		defaultproperty.MakeIntProperty(property.Attack, 3, property.Public, property.Character),
		defaultproperty.MakeIntProperty(property.AttackRange, 1, property.Public, property.Character),
		defaultproperty.MakeIntProperty(property.Defense, 0, property.Public, property.Character),
		defaultproperty.MakeIntProperty(property.JumpHeight, 2, property.Public, property.Character),
		defaultproperty.MakeIntProperty(property.IsDying, -1, property.Public, property.Character),
		defaultproperty.MakeBoolProperty(property.HasMoved, false, property.GameMaster, property.Character),
		defaultproperty.MakeBoolProperty(property.HasActed, false, property.GameMaster, property.Character),
	}
}

*/

func TestEffectApplicatorDamage1Target(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// This effect will double the damage (3 -> 6)
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.Damage, 200, property.Public, property.Skill))

	// Base defense is 0, No armoring provided, so damage should be 6

	log := logrus.WithField("test", "TestEffectApplicatorDamage1Target")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount != 6 {
		t.Errorf("Expected 6 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 4 {
		t.Errorf("Expected target to have 4 HP, got %d", hp)
	}
}

func TestEffectApplicatorDamage1TargetShield(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// This effect will double the damage (3 -> 6)
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.Damage, 200, property.Public, property.Skill))

	fake.Target.UpdatePropertyValue(property.Shield, 5)

	// Base defense is 0, No armoring provided, so damage should be 6
	// Shield should absorb 5, leaving 1 damage

	log := logrus.WithField("test", "TestEffectApplicatorDamage1TargetShield")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount !=  1 {
		t.Errorf("Expected  1 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 9 {
		t.Errorf("Expected target to have 9 HP, got %d", hp)
	}

	shield := damaged[0].Target.GetPropertyC(property.Shield).GetValue()
	if shield != 0 {
		t.Errorf("Expected shield to be 0, got %d", shield)
	}
}

func TestEffectApplicatorDamage1Defense(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// This effect will double the damage (3 -> 6)
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.Damage, 200, property.Public, property.Skill))

	fake.Target.RepsertPropertyValue(property.ArmorRating, 5)

	// Base defense is 0, 5 armore provided, so damage should be 1

	log := logrus.WithField("test", "TestEffectApplicatorDamage1Defense")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount !=  1 {
		t.Errorf("Expected  1 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 9 {
		t.Errorf("Expected target to have 9 HP, got %d", hp)
	}
}

func TestEffectApplicatorPoison1Defense(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// this effect adds 2 flat poison damage.
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.PoisonPower, 2, property.Public, property.Skill))

	fake.Target.RepsertPropertyValue(property.ArmorRating, 5)
	// Base attack is 3, without multiplier, so damage should be 2 = (3 - 5) +2 poison

	log := logrus.WithField("test", "TestEffectApplicatorPoison1Defense")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount !=  2 {
		t.Errorf("Expected  2 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 8 {
		t.Errorf("Expected target to have 8 HP, got %d", hp)
	}
}

func TestEffectApplicatorStun1Armor(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// this effect adds 2 flat poison damage.
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.StunPower, 2, property.Public, property.Skill))

	fake.Target.RepsertPropertyValue(property.Defense, 5)
	// Base attack is 3, without multiplier, so damage should be 2 = (3 - 5) +2 stun

	log := logrus.WithField("test", "TestEffectApplicatorStun1Armor")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount !=  2 {
		t.Errorf("Expected  2 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 8 {
		t.Errorf("Expected target to have 8 HP, got %d", hp)
	}
}

func TestEffectApplicatorPoisoned(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// this effect adds 2 flat poison damage. Total damage is 5
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.PoisonPower, 2, property.Public, property.Skill))
	// forces chance to 100%
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.Public, property.Skill))

	tools.TesterRand(func(n int) int {
		return 0 // always hit, always poison or stun
	})

	log := logrus.WithField("test", "TestEffectApplicatorPoisoned")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount !=  5 {
		t.Errorf("Expected 6 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 5 {
		t.Errorf("Expected target to have 5 HP, got %d", hp)
	}

	// Poisoned!
	poisoned := damaged[0].Target.GetPropertyI(property.Poison).I()

	if poisoned != 2 {
		t.Errorf("Expected target to be poisoned, got %d", poisoned)
	}
}

func TestEffectApplicatorPoisonedTwice(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// this effect adds 2 flat poison damage. Total damage is 5
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.PoisonPower, 2, property.Public, property.Skill))
	// forces chance to 100%
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.Public, property.Skill))

	tools.TesterRand(func(n int) int {
		return 0 // always hit, always poison or stun
	})

	log := logrus.WithField("test", "TestEffectApplicatorPoisoned")

	// Apply the effect
	ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)
	damaged, _, _, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 0 { // oops it's also dead ;) But that's for the rules to decide. Maybe.
		t.Errorf("Expected target to have 0 HP, got %d", hp)
	}

	// Poisoned!
	poisoned := damaged[0].Target.GetPropertyI(property.Poison).I()

	if poisoned != 4 {
		t.Errorf("Expected target to be poisoned, got %d", poisoned)
	}
}

func TestEffectApplicatorStuned(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	// this effect adds 2 flat stun damage. Total damage is 5
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.StunPower, 2, property.Public, property.Skill))
	// forces chance to 100%
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.StunChance, 100, property.Public, property.Skill))

	tools.TesterRand(func(n int) int {
		return 0 // always hit, always stun
	})

	log := logrus.WithField("test", "TestEffectApplicatorStuned")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount !=  5 {
		t.Errorf("Expected 6 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 5 {
		t.Errorf("Expected target to have 5 HP, got %d", hp)
	}

	// Stuned!
	stuned := damaged[0].Target.GetPropertyI(property.Stun).I()

	if stuned != 2 {
		t.Errorf("Expected target to be stuned, got %d", stuned)
	}
}

func TestEffectApplicatorCrit(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()

	tools.TesterRand(func(n int) int {
		return 0 // always hit, always stun
	})

	// This effect only allow critic: 3 x 200% (crit mutliplier) = 6
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.Damage, 100, property.Public, property.Skill))
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.CriticalChance, 100, property.Public, property.Skill))
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.CriticalMultiplier, 200, property.Public, property.Skill))

	// Base defense is 0, No armoring provided, so damage should be 6

	log := logrus.WithField("test", "TestEffectApplicatorCrit")

	// Apply the effect
	damaged, _, credits, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(credits) != 1 || credits[0].Amount != 6 {
		t.Errorf("Expected 6 credits, got %v", credits)
	}

	if len(damaged) != 1 {
		t.Errorf("Expected 1 target to be damaged, got %d", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Errorf("Expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()

	if hp != 4 {
		t.Errorf("Expected target to have 4 HP, got %d", hp)
	}
}
