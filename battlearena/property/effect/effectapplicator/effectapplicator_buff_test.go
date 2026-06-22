package effectapplicator

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/sirupsen/logrus"
)

// @test-link [[mechanic_effect_shield]]
// @test-link [[mechanic_effect_heal]]

func TestShieldOvershield_CappedAt2xMaxHP(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	// Add massive shield (100 points). Max HP is 10, so shield cap is 2×10 = 20.
	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.ShieldPower, 100, property.Public, property.Skill),
	)

	log := logrus.WithField("test", "TestShieldOvershield_CappedAt2xMaxHP")
	_, affected, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(affected) != 1 {
		t.Fatalf("expected 1 affected target, got %d", len(affected))
	}

	shield := affected[0].Target.GetPropertyC(property.Shield).GetValue()
	maxHP := affected[0].Target.GetPropertyC(property.HP).GetMaxValue()
	if shield > maxHP*2 {
		t.Errorf("shield %d exceeds 2×maxHP cap (%d)", shield, maxHP*2)
	}
	if shield != 20 {
		t.Errorf("expected shield=20 (capped at 2×10), got %d", shield)
	}
}

func TestHealAndShield_Combined(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.Heal, 5, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.ShieldPower, 3, property.Public, property.Skill),
	)

	// Reduce target HP first
	fake.Target.UpdatePropertyValue(property.HP, 5)

	log := logrus.WithField("test", "TestHealAndShield_Combined")
	_, affected, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(affected) != 1 {
		t.Fatalf("expected 1 affected target, got %d", len(affected))
	}

	hp := affected[0].Target.GetPropertyC(property.HP).GetValue()
	if hp != 10 {
		t.Errorf("expected HP=10 after healing 5 from 5, got %d", hp)
	}

	shield := affected[0].Target.GetPropertyC(property.Shield).GetValue()
	if shield != 3 {
		t.Errorf("expected Shield=3, got %d", shield)
	}
}

func TestCleanse_PoisonAndStun(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	// Negative values cleanse
	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.PoisonPower, -10, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.StunPower, -10, property.Public, property.Skill),
	)

	// Apply poison and stun to target first
	fake.Target.RepsertPropertyValue(property.Poison, 5)
	fake.Target.RepsertPropertyValue(property.Stun, 3)

	log := logrus.WithField("test", "TestCleanse_PoisonAndStun")
	_, affected, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(affected) != 1 {
		t.Fatalf("expected 1 affected target, got %d", len(affected))
	}

	poison := affected[0].Target.GetPropertyI(property.Poison).I()
	if poison != 0 {
		t.Errorf("expected poison=0 after cleanse, got %d", poison)
	}

	stun := affected[0].Target.GetPropertyI(property.Stun).I()
	if stun != 0 {
		t.Errorf("expected stun=0 after cleanse, got %d", stun)
	}
}
