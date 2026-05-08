package effectapplicator

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/sirupsen/logrus"
)

// @test-link [[mechanic_mech_effect_stun]]
// @test-link [[mechanic_mech_effect_poison]]

// --- Stun pairing tests ---

func TestStunApplied_WhenBothPowerAndChance(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.StunPower, 5, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.StunChance, 100, property.Public, property.Skill),
	)

	// Force all random rolls to succeed (return 0).
	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestStunApplied_WhenBothPowerAndChance")
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(damaged) != 1 {
		t.Fatalf("expected 1 damaged target, got %d", len(damaged))
	}

	stun := damaged[0].Target.GetPropertyI(property.Stun).I()
	if stun <= 0 {
		t.Errorf("expected stun > 0 when both StunPower and StunChance are set, got %d", stun)
	}
}

func TestStunNotApplied_WhenChanceZero(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.StunPower, 5, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.StunChance, 0, property.Public, property.Skill),
	)

	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestStunNotApplied_WhenChanceZero")
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(damaged) != 1 {
		t.Fatalf("expected 1 damaged target, got %d", len(damaged))
	}

	stun := damaged[0].Target.GetPropertyI(property.Stun).I()
	if stun != 0 {
		t.Errorf("expected stun=0 when StunChance=0, got %d", stun)
	}
}

func TestStunNotApplied_WhenPowerZero(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	// StunPower=0 means IsDamaging() returns false (no positive damage properties).
	// The effect won't enter the damage branch at all.
	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.StunPower, 0, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.StunChance, 100, property.Public, property.Skill),
	)

	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestStunNotApplied_WhenPowerZero")
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	// With StunPower=0, IsDamaging() is false → no damage branch runs.
	if len(damaged) != 0 {
		t.Errorf("expected 0 damaged targets (StunPower=0 is not damaging), got %d", len(damaged))
	}

	// Verify stun was NOT applied via any path.
	stun := fake.Entities[0].GetPropertyI(property.Stun).I()
	if stun != 0 {
		t.Errorf("expected stun=0 when StunPower=0, got %d", stun)
	}
}

// --- Poison pairing tests ---

func TestPoisonApplied_WhenBothPowerAndChance(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.PoisonPower, 3, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.Public, property.Skill),
	)

	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestPoisonApplied_WhenBothPowerAndChance")
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(damaged) != 1 {
		t.Fatalf("expected 1 damaged target, got %d", len(damaged))
	}

	poison := damaged[0].Target.GetPropertyI(property.Poison).I()
	if poison <= 0 {
		t.Errorf("expected poison > 0 when both PoisonPower and PoisonChance are set, got %d", poison)
	}
}

func TestPoisonNotApplied_WhenChanceZero(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.PoisonPower, 3, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.PoisonChance, 0, property.Public, property.Skill),
	)

	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestPoisonNotApplied_WhenChanceZero")
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(damaged) != 1 {
		t.Fatalf("expected 1 damaged target, got %d", len(damaged))
	}

	poison := damaged[0].Target.GetPropertyI(property.Poison).I()
	if poison != 0 {
		t.Errorf("expected poison=0 when PoisonChance=0, got %d", poison)
	}
}

func TestPoisonNotApplied_WhenPowerZero(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	// PoisonPower=0 means IsDamaging() returns false (no positive damage properties).
	// The effect won't enter the damage branch at all.
	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.PoisonPower, 0, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.Public, property.Skill),
	)

	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestPoisonNotApplied_WhenPowerZero")
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	// With PoisonPower=0, IsDamaging() is false → no damage branch runs.
	if len(damaged) != 0 {
		t.Errorf("expected 0 damaged targets (PoisonPower=0 is not damaging), got %d", len(damaged))
	}

	// Verify poison was NOT applied via any path.
	poison := fake.Entities[0].GetPropertyI(property.Poison).I()
	if poison != 0 {
		t.Errorf("expected poison=0 when PoisonPower=0, got %d", poison)
	}
}

// --- Stacking tests ---

func TestPoisonStacking(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.PoisonPower, 3, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.Public, property.Skill),
	)

	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestPoisonStacking")

	// Apply twice
	ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(damaged) != 1 {
		t.Fatalf("expected 1 damaged target, got %d", len(damaged))
	}

	poison := damaged[0].Target.GetPropertyI(property.Poison).I()
	// Defense is 0, so truepoison = max(3 - 0, 0) = 3 each time. Stack = 6.
	if poison != 6 {
		t.Errorf("expected poison=6 after two applications (3+3), got %d", poison)
	}
}

func TestStunStacking(t *testing.T) {
	fake := makeTestingEnvironment()
	eff := effect.New()

	eff.Properties = append(eff.Properties,
		defaultproperty.MakeIntProperty(property.StunPower, 2, property.Public, property.Skill),
		defaultproperty.MakeIntProperty(property.StunChance, 100, property.Public, property.Skill),
	)

	tools.TesterRand(func(n int) int { return 0 })
	defer tools.TesterRand(nil)

	log := logrus.WithField("test", "TestStunStacking")

	// Apply twice
	ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)
	damaged, _, _, err, _ := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(damaged) != 1 {
		t.Fatalf("expected 1 damaged target, got %d", len(damaged))
	}

	stun := damaged[0].Target.GetPropertyI(property.Stun).I()
	// ArmorRating is 0, so truestun = max(2 - 0, 0) = 2 each time. Stack = 4.
	if stun != 4 {
		t.Errorf("expected stun=4 after two applications (2+2), got %d", stun)
	}
}
