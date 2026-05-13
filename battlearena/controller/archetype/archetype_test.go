package archetype_test

// @test-link [[mec_ai_archetype_system]]
// @test-link [[rule_ai_team_composition_rules]]

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/archetype"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
)

// TestRegistryCoversAllFour verifies that all four canonical archetypes are present in the registry.
func TestRegistryCoversAllFour(t *testing.T) {
	for _, slug := range []string{"fighter", "ranger", "support", "sneak"} {
		a, ok := archetype.Get(slug)
		if !ok {
			t.Errorf("archetype %q not in registry", slug)
			continue
		}
		if a.Slug() != slug {
			t.Errorf("slug mismatch: got %q, want %q", a.Slug(), slug)
		}
	}
}

// TestBehaviorStackHasBaselineAsLastLayer ensures each archetype stack terminates with the always-active baseline.
func TestBehaviorStackHasBaselineAsLastLayer(t *testing.T) {
	for _, slug := range []string{"fighter", "ranger", "support", "sneak"} {
		a, _ := archetype.Get(slug)
		lb := a.Behavior()
		if len(lb.Layers) == 0 {
			t.Errorf("%s: empty behavior stack", slug)
			continue
		}
		last := lb.Layers[len(lb.Layers)-1]
		if _, ok := last.(*behavior.AggressiveBehavior); !ok {
			t.Errorf("%s: last layer is %T, want *behavior.AggressiveBehavior", slug, last)
		}
		if last.BaseActivation() != 1.0 {
			t.Errorf("%s: baseline BaseActivation = %v, want 1.0", slug, last.BaseActivation())
		}
	}
}

// TestBehaviorStackHasAtLeastTwoLayers checks that each archetype has at least one specialized layer plus the baseline.
func TestBehaviorStackHasAtLeastTwoLayers(t *testing.T) {
	for _, slug := range []string{"fighter", "ranger", "support", "sneak"} {
		a, _ := archetype.Get(slug)
		if len(a.Behavior().Layers) < 2 {
			t.Errorf("%s: expected ≥ 2 layers, got %d", slug, len(a.Behavior().Layers))
		}
	}
}

// TestStatWeightsArePositive verifies that each archetype's stat weights sum to a positive value.
func TestStatWeightsArePositive(t *testing.T) {
	for _, slug := range []string{"fighter", "ranger", "support", "sneak"} {
		a, _ := archetype.Get(slug)
		w := a.StatWeights()
		total := w.HP + w.SP + w.MP + w.Attack + w.Defense + w.Movement + w.AttackRange
		if total <= 0 {
			t.Errorf("%s: StatWeights sum to zero", slug)
		}
	}
}

// TestRandomForRespectsCompositionConstraints confirms that RandomFor never returns a constrained archetype when the team already has one.
func TestRandomForRespectsCompositionConstraints(t *testing.T) {
	// With support + sneak on team, only fighter or ranger should be returned.
	for i := 0; i < 50; i++ {
		got := archetype.RandomFor([]string{"support", "sneak"})
		if got == "support" || got == "sneak" {
			t.Errorf("RandomFor returned %q despite team already having support+sneak", got)
		}
	}
}

// TestRandomForEmptyTeam confirms RandomFor returns a valid archetype for an empty team without panicking.
func TestRandomForEmptyTeam(t *testing.T) {
	// Any archetype is valid for an empty team — just check no panic.
	got := archetype.RandomFor([]string{})
	if _, ok := archetype.Get(got); !ok {
		t.Errorf("RandomFor returned unknown slug %q", got)
	}
}

// TestBuildSkillBundleReturnsRequestedCount verifies that BuildSkillBundle returns exactly the requested number of skills.
func TestBuildSkillBundleReturnsRequestedCount(t *testing.T) {
	for _, slug := range []string{"fighter", "ranger", "support", "sneak"} {
		a, _ := archetype.Get(slug)
		skills := a.BuildSkillBundle("I", 3)
		if len(skills) != 3 {
			t.Errorf("%s: BuildSkillBundle(I, 3) returned %d skills, want 3", slug, len(skills))
		}
	}
}
