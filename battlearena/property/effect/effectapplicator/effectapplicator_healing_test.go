package effectapplicator

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/sirupsen/logrus"
)

func TestEffectApplicatorHeal(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.Heal, 5, property.Public, property.Skill))

	log := logrus.WithField("test", "TestEffectApplicatorHeal")

	// ensure target has something to heal.
	fake.Target.UpdatePropertyValue(property.HP, 5)

	// Apply the effect
	_, affected, _, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(affected) != 1 {
		t.Errorf("Expected 1 target to be affected, got %d", len(affected))
	}

	if affected[0].ID != fake.Target.ID {
		t.Errorf("Expected target to be affected, got %s", affected[0].ID)
	}

	hp := affected[0].GetPropertyC(property.HP).GetValue()

	if hp != 10 {
		t.Errorf("Expected target to have 10 HP, got %d", hp)
	}
}

func TestEffectApplicatorOverheal(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.Heal, 15, property.Public, property.Skill))

	log := logrus.WithField("test", "TestEffectApplicatorOverheal")

	// ensure target has something to heal.
	fake.Target.UpdatePropertyValue(property.HP, 5)

	// Apply the effect
	_, affected, _, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(affected) != 1 {
		t.Errorf("Expected 1 target to be affected, got %d", len(affected))
	}

	if affected[0].ID != fake.Target.ID {
		t.Errorf("Expected target to be affected, got %s", affected[0].ID)
	}

	hp := affected[0].GetPropertyC(property.HP).GetValue()

	if hp != 10 {
		t.Errorf("Expected target to have 10 HP, got %d", hp)
	}
}

func TestEffectApplicatorCurePoison(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.PoisonPower, -15, property.Public, property.Skill))

	log := logrus.WithField("test", "TestEffectApplicatorCurePoison")

	// ensure target has something to heal.
	fake.Target.RepsertPropertyValue(property.Poison, 5)

	// Apply the effect
	_, affected, _, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(affected) != 1 {
		t.Errorf("Expected 1 target to be affected, got %d", len(affected))
	}

	if affected[0].ID != fake.Target.ID {
		t.Errorf("Expected target to be affected, got %s", affected[0].ID)
	}

	poison := affected[0].GetPropertyI(property.Poison).I()

	if poison != 0 {
		t.Errorf("Expected target to have 0 poison, got %d", poison)
	}
}

func TestEffectApplicatorCureStun(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.StunPower, -15, property.Public, property.Skill))

	log := logrus.WithField("test", "TestEffectApplicatorCureStun")

	// ensure target has something to heal.
	fake.Target.RepsertPropertyValue(property.Stun, 5)

	// Apply the effect
	_, affected, _, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(affected) != 1 {
		t.Errorf("Expected 1 target to be affected, got %d", len(affected))
	}

	if affected[0].ID != fake.Target.ID {
		t.Errorf("Expected target to be affected, got %s", affected[0].ID)
	}

	stun := affected[0].GetPropertyI(property.Stun).I()

	if stun != 0 {
		t.Errorf("Expected target to have 0 stun, got %d", stun)
	}
}

func TestEffectApplicatorShielding(t *testing.T) {
	fake := makeTestingEnvironment()

	// update the effect to be applied. Default ... has nothing.
	eff := effect.New()
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.ShieldPower, 5, property.Public, property.Skill))

	log := logrus.WithField("test", "TestEffectApplicatorShielding")

	// Apply the effect
	_, affected, _, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Errorf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(affected) != 1 {
		t.Errorf("Expected 1 target to be affected, got %d", len(affected))
	}

	if affected[0].ID != fake.Target.ID {
		t.Errorf("Expected target to be affected, got %s", affected[0].ID)
	}

	shield := affected[0].GetPropertyC(property.Shield).GetValue()

	if shield != 5 {
		t.Errorf("Expected target to have 5 shield, got %d", shield)
	}
}
