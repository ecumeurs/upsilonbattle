package battletest

import (
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilontypes/property/effect"
)

// SkillSpec is a fluent builder for skills used in tests. It starts from skill.New()
// (Direct behavior, range 1, Single zone, Entity target) and lets tests override
// targeting, behavior, costs, and effect properties.
type SkillSpec struct {
	sk skill.Skill
}

// NewSkill starts a SkillSpec with the given diegetic name.
func NewSkill(name string) *SkillSpec {
	sk := skill.New()
	sk.Name = name
	return &SkillSpec{sk: sk}
}

// Behavior sets the skill behavior (Direct, Trap, Passive, …).
func (b *SkillSpec) Behavior(bt def.BehaviorType) *SkillSpec {
	beh := def.DefaultBehavior()
	beh.Set(string(bt))
	b.sk.Behavior = beh
	return b
}

// TargetType sets the skill's target type (Entity, Self, Tile, EnemyOnly, …).
func (b *SkillSpec) TargetType(tt def.TargetTypes) *SkillSpec {
	tp := def.DefaultTargetType()
	tp.Set(string(tt))
	b.sk.Targeting[string(property.TargetType)] = tp
	return b
}

// Range sets the skill's minimum and maximum range.
func (b *SkillSpec) Range(min, max int) *SkillSpec {
	b.sk.Targeting[string(property.Range)] =
		defaultproperty.MakeIntCounterProperty(property.Range, min, max, property.FriendlyController, property.Skill)
	return b
}

// Zone sets the skill's area-of-effect pattern (e.g. "Neighbours", "Line:3").
func (b *SkillSpec) Zone(patternStr string) *SkillSpec {
	z := def.DefaultZone()
	z.Set(patternStr)
	b.sk.Targeting[string(property.Zone)] = z
	return b
}

// Cost adds an arbitrary cost property (HPLeech, MvtCost, Cooldown, …).
func (b *SkillSpec) Cost(p property.Property) *SkillSpec {
	b.sk.Costs[p.Name(property.GameMaster)] = p
	return b
}

// MvtCost sets the skill's movement-point cost.
func (b *SkillSpec) MvtCost(v int) *SkillSpec {
	return b.Cost(defaultproperty.MakeIntProperty(property.MovementCost, v, property.FriendlyController, property.Skill))
}

// Effect adds an arbitrary effect property.
func (b *SkillSpec) Effect(p property.Property) *SkillSpec {
	b.sk.Effect.Properties = append(b.sk.Effect.Properties, p)
	return b
}

// Damage adds a Damage effect property (percentage of Attack).
func (b *SkillSpec) Damage(v int) *SkillSpec {
	return b.Effect(defaultproperty.MakeIntProperty(property.DamageScale, v, property.FriendlyController, property.Skill))
}

// Heal adds a Heal effect property.
func (b *SkillSpec) Heal(v int) *SkillSpec {
	return b.Effect(defaultproperty.MakeIntProperty(property.Heal, v, property.FriendlyController, property.Skill))
}

// ShieldPower adds a ShieldPower effect property.
func (b *SkillSpec) ShieldPower(v int) *SkillSpec {
	return b.Effect(defaultproperty.MakeIntProperty(property.ShieldPower, v, property.FriendlyController, property.Skill))
}

// Reposition makes the skill a movement skill: it displaces the given subject (Self or Target)
// by a signed distance along the casting ray (positive = forward/away, negative = pull).
func (b *SkillSpec) Reposition(subject def.RepositionSubjectType, dist int) *SkillSpec {
	b.Effect(def.RepositionSubject(subject))
	return b.Effect(def.RepositionDistance(dist))
}

// Build returns the assembled skill.
func (b *SkillSpec) Build() skill.Skill {
	return b.sk
}

// PoisonTrap builds a positional-effect that applies guaranteed poison when triggered.
// It is the canonical observable for fly-over vs landing assertions: poison applied
// (or trap consumed) means the cell fired; unchanged poison means it was flown over.
func PoisonTrap(power int, trigger property.TriggerTypeValue, removeOnTrigger bool) effect.Effect {
	props := []property.Property{
		defaultproperty.MakeIntProperty(property.PoisonPower, power, property.GameMaster, property.Skill),
		defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.GameMaster, property.Skill),
		def.TriggerType(trigger),
	}
	if removeOnTrigger {
		props = append(props, def.RemoveOnTrigger(true))
	}
	return effect.Effect{
		Name:       "Poison Trap",
		Properties: props,
	}
}
