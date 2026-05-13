package archetype

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior/micro"
	"github.com/ecumeurs/upsilontypes/entity/skill"
)

// SupportArchetype — healer/shielder that stays behind the front line and sustains allies.
//
// Stack: HealAlly(0.85) → ShieldAlly(0.75) → StayBehindFront(0.7) → Baseline(1.0)
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type SupportArchetype struct{}

// Slug returns the canonical identifier used in the archetype registry.
func (SupportArchetype) Slug() string  {
	return "support"
}

// Behavior returns the layered behavior stack for this archetype.
func (SupportArchetype) Behavior() *behavior.LayeredBehavior {
	return behavior.NewLayeredBehavior(
		&micro.HealAlly{},
		&micro.ShieldAlly{},
		&micro.StayBehindFront{},
		&behavior.AggressiveBehavior{},
	)
}

// AllowedSkillTags returns the skill tags this archetype may equip.
func (SupportArchetype) AllowedSkillTags() []string {
	return []string{"heal", "shield", "buff"}
}

// ForbidSkillTags returns the skill tags this archetype may not equip.
func (SupportArchetype) ForbidSkillTags() []string {
	return []string{"melee", "dot"}
}

// BuildSkillBundle generates numSkills skills filtered by this archetype's tag constraints.
func (SupportArchetype) BuildSkillBundle(grade string, numSkills int) []skill.Skill {
	return buildBundle(SupportArchetype{}.AllowedSkillTags(), SupportArchetype{}.ForbidSkillTags(), grade, numSkills)
}

// StatWeights returns the relative CP allocation weights for each stat.
func (SupportArchetype) StatWeights() StatWeights {
	return StatWeights{
		SP:      0.35,
		MP:      0.25,
		HP:      0.20,
		Defense: 0.10,
		Attack:  0.10,
	}
}
