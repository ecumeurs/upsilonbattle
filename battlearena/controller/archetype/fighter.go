package archetype

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior/micro"
	"github.com/ecumeurs/upsilontypes/entity/skill"
)

// FighterArchetype — aggressive melee bruiser that locks onto a target and charges in.
//
// Stack: BattleFocus(0.8) → ChargeIn(0.7) → Baseline(1.0)
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type FighterArchetype struct{}

// Slug returns the canonical identifier used in the archetype registry.
func (FighterArchetype) Slug() string  {
	return "fighter"
}

// Behavior returns the layered behavior stack for this archetype.
func (FighterArchetype) Behavior() *behavior.LayeredBehavior {
	return behavior.NewLayeredBehavior(
		&micro.BattleFocus{},
		&micro.ChargeIn{},
		&behavior.AggressiveBehavior{},
	)
}

// AllowedSkillTags returns the skill tags this archetype may equip.
func (FighterArchetype) AllowedSkillTags() []string {
	return []string{"melee", "buff", "aoe"}
}

// ForbidSkillTags returns the skill tags this archetype may not equip.
func (FighterArchetype) ForbidSkillTags() []string {
	return []string{"heal", "shield", "trap"}
}

// BuildSkillBundle generates numSkills skills filtered by this archetype's tag constraints.
func (FighterArchetype) BuildSkillBundle(grade string, numSkills int) []skill.Skill {
	return buildBundle(FighterArchetype{}.AllowedSkillTags(), FighterArchetype{}.ForbidSkillTags(), grade, numSkills)
}

// StatWeights returns the relative CP allocation weights for each stat.
func (FighterArchetype) StatWeights() StatWeights {
	return StatWeights{
		HP:      0.40,
		Attack:  0.35,
		Defense: 0.15,
		SP:      0.05,
		MP:      0.05,
	}
}
