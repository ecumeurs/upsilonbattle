// @test-link [[mech_combat_attack_computation]]
package effectapplicator

import (
	"testing"

	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/sirupsen/logrus"
)

// TestEffectApplicatorDodge_ReadsFromTargetNotCaster is the ISS-145 regression
// guard: the hit test must consult the TARGET's Dodge, not the caster's. The
// caster is given a Dodge of 100 (which, if wrongly read as the defender's,
// would suppress the caster's own attack) while the target keeps the default
// Dodge of 0. With the RNG roll fixed at 50 (below Accuracy's default 100, but
// not below the caster's own Dodge of 100), the attack must land: a target
// Dodge of 0 makes accuracy(100)-dodge(0)=100, so 50 < 100 is a hit. If Dodge
// were still (incorrectly) read from the caster, accuracy(100)-dodge(100)=0
// and 50 < 0 is false, so the target would incorrectly evade.
func TestEffectApplicatorDodge_ReadsFromTargetNotCaster(t *testing.T) {
	fake := makeTestingEnvironment()

	// Caster has a high Dodge; this must have zero bearing on whether the
	// caster's own attack lands against the target.
	fake.Caster.RepsertPropertyValue(property.Dodge, 100)
	// Target keeps the default Dodge of 0 (untouched).

	eff := effect.New()
	eff.Properties = append(eff.Properties, defaultproperty.MakeIntProperty(property.Damage, 100, property.Public, property.Skill))

	tools.TesterRand(func(n int) int {
		return 50 // fixed roll: hits iff accuracy-dodge > 50
	})

	log := logrus.WithField("test", "TestEffectApplicatorDodge_ReadsFromTargetNotCaster")

	damaged, _, _, err, errkey := ApplyDirectEffect(log, &fake.Caster, *eff, fake.TargetPos, fake.Pos, fake.Grid, fake.Entities)

	if err != "" {
		t.Fatalf("Error applying effect: %s (%s)", err, errkey)
	}

	if len(damaged) != 1 {
		t.Fatalf("expected the target's own (0) Dodge to let the attack land, got %d damaged targets (caster's Dodge of 100 wrongly suppressed the hit)", len(damaged))
	}

	if damaged[0].Target.ID != fake.Target.ID {
		t.Fatalf("expected target to be damaged, got %s", damaged[0].Target.ID)
	}

	hp := damaged[0].Target.GetPropertyC(property.HP).GetValue()
	if hp != 7 {
		t.Fatalf("expected target to have 7 HP after a 3-damage hit, got %d", hp)
	}
}
