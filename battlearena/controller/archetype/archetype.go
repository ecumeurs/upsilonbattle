// Package archetype defines the four AI archetypes (Fighter, Ranger, Support, Sneak),
// their layered behavior stacks, skill-generation tag filters, and stat-weight profiles.
//
// Usage:
//
//	arch, ok := archetype.Get("fighter")
//	controller := controllers.NewAIController(id, name, arch.Behavior())
//	skills := arch.BuildSkillBundle(grade, 3)
//
// @spec-link [[mec_ai_archetype_system]]
package archetype

import (
	"math/rand"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/entity/skill/skillgenerator"
)

// Archetype defines the tactical profile of an AI entity.
//
// @spec-link [[mec_ai_archetype_system]]
type Archetype interface {
	// Slug returns the canonical identifier ("fighter", "ranger", "support", "sneak").
	Slug() string
	// Behavior returns the layered behavior stack for this archetype.
	// The always-active baseline (AggressiveBehavior) must be the last element.
	Behavior() *behavior.LayeredBehavior
	// AllowedSkillTags returns the skill category tags this archetype may equip.
	AllowedSkillTags() []string
	// ForbidSkillTags returns category tags explicitly excluded for this archetype.
	ForbidSkillTags() []string
	// BuildSkillBundle generates numSkills skills suitable for this archetype at grade.
	BuildSkillBundle(grade string, numSkills int) []skill.Skill
	// StatWeights describes relative priority for each stat when distributing a CP pool.
	StatWeights() StatWeights
}

// StatWeights describes how to apportion a CP pool across entity statistics.
// Values are relative — the bridge normalises the sum to 1.0.
//
// @spec-link [[mechanic_ai_progression_matching]]
type StatWeights struct {
	HP          float64
	SP          float64
	MP          float64
	Attack      float64
	Defense     float64
	Movement    float64
	AttackRange float64
}

// buildBundle is the shared skill-generation loop used by all concrete archetypes.
// It retries up to maxRetries times per skill slot on generation error.
func buildBundle(allowed, forbid []string, grade string, numSkills int) []skill.Skill {
	const maxRetries = 5
	result := make([]skill.Skill, 0, numSkills)
	for i := 0; i < numSkills; i++ {
		var sk skill.Skill
		for attempt := 0; attempt < maxRetries; attempt++ {
			generated, _, err := skillgenerator.Generate(skillgenerator.GenerateRequest{
				TargetGrade: grade,
				AllowedTags: allowed,
				ForbidTags:  forbid,
			})
			if err == nil {
				sk = generated
				break
			}
		}
		result = append(result, sk)
	}
	return result
}

// ── Registry ──────────────────────────────────────────────────────────────────

var registry = map[string]Archetype{
	"fighter": FighterArchetype{},
	"ranger":  RangerArchetype{},
	"support": SupportArchetype{},
	"sneak":   SneakArchetype{},
}

// Get returns the Archetype registered under slug, or false if unknown.
func Get(slug string) (Archetype, bool) {
	a, ok := registry[slug]
	return a, ok
}

// All returns all registered archetypes in undefined order.
func All() []Archetype {
	result := make([]Archetype, 0, len(registry))
	for _, a := range registry {
		result = append(result, a)
	}
	return result
}

// RandomFor returns a random archetype slug respecting team-composition constraints:
//   - at most 1 "support" per team
//   - at most 1 "sneak" per team
//
// existingTeam is the list of slugs already assigned to teammates.
//
// @spec-link [[rule_ai_team_composition_rules]]
func RandomFor(existingTeam []string) string {
	hasSupport := false
	hasSneak := false
	for _, s := range existingTeam {
		switch s {
		case "support":
			hasSupport = true
		case "sneak":
			hasSneak = true
		}
	}

	available := make([]string, 0, 4)
	for slug := range registry {
		if slug == "support" && hasSupport {
			continue
		}
		if slug == "sneak" && hasSneak {
			continue
		}
		available = append(available, slug)
	}
	if len(available) == 0 {
		return "fighter" // fallback: always valid
	}
	return available[rand.Intn(len(available))]
}
