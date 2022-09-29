package skill

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
)

type Skill struct {
	Name string

	Behavior def.BehaviorProperty

	Targeting []property.Property
	Cost      []property.Property
	Effect    []property.Property
}

// New
func New() *Skill {
	return &Skill{
		Name:      "New Skill",
		Behavior:  *def.DefaultBehaviorProperty(),
		Targeting: []property.Property{},
		Cost:      []property.Property{},
		Effect:    []property.Property{},
	}
}

// NewSkill
func NewSkill(name string, targeting, cost, effect []property.Property) *Skill {
	return &Skill{
		Name:      name,
		Targeting: targeting,
		Cost:      cost,
		Effect:    effect,
	}
}
