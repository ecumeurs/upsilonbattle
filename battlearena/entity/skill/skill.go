package skill

import "github.com/ecumeurs/upsilonbattle/battlearena/property"

type Skill struct {
	Name string

	Targeting []property.Property
	Cost      []property.Property
	Effect    []property.Property
}

// New
func New() *Skill {
	return &Skill{
		Name:      "New Skill",
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
