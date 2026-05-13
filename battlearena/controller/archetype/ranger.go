package archetype

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior/micro"
	"github.com/ecumeurs/upsilontypes/entity/skill"
)

// RangerArchetype — ranged skirmisher that kites enemies and targets back-line support.
//
// Stack: KiteAway(0.8) → MaintainRange(0.75) → FocusBackline(0.7) → Baseline(1.0)
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type RangerArchetype struct{}

// Slug returns the canonical identifier used in the archetype registry.
func (RangerArchetype) Slug() string  {
	return "ranger"
}

// Behavior returns the layered behavior stack for this archetype.
func (RangerArchetype) Behavior() *behavior.LayeredBehavior {
	return behavior.NewLayeredBehavior(
		&micro.KiteAway{},
		&micro.MaintainRange{},
		&micro.FocusBackline{},
		&behavior.AggressiveBehavior{},
	)
}

// AllowedSkillTags returns the skill tags this archetype may equip.
func (RangerArchetype) AllowedSkillTags() []string {
	return []string{"ranged", "aoe", "debuff", "dot"}
}

// ForbidSkillTags returns the skill tags this archetype may not equip.
func (RangerArchetype) ForbidSkillTags() []string {
	return []string{"melee", "shield"}
}

// BuildSkillBundle generates numSkills skills filtered by this archetype's tag constraints.
func (RangerArchetype) BuildSkillBundle(grade string, numSkills int) []skill.Skill {
	return buildBundle(RangerArchetype{}.AllowedSkillTags(), RangerArchetype{}.ForbidSkillTags(), grade, numSkills)
}

// StatWeights returns the relative CP allocation weights for each stat.
func (RangerArchetype) StatWeights() StatWeights {
	return StatWeights{
		Attack:      0.35,
		AttackRange: 0.20,
		SP:          0.20,
		HP:          0.15,
		Movement:    0.10,
	}
}
