package entitygenerator

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"

	"github.com/ecumeurs/upsilontools/tools"
)

var propertyRandomizers = map[string]tools.IntRange{
	"HP":          tools.IntRange{Start: 3, End: 20},
	"Attack":      tools.IntRange{Start: 1, End: 5},
	"Defense":     tools.IntRange{Start: 0, End: 3},
	"Movement":    tools.IntRange{Start: 3, End: 7},
	"AttackRange": tools.IntRange{Start: 1, End: 3},
	"JumpHeight":  tools.IntRange{Start: 1, End: 3},
}

func GenerateRandomEntity() entity.Entity {
	ent := entity.NewEntity()

	for _, v := range property.DefaultPropertiesForCharacter() {
		ent.Properties[v.Name(property.OwnController)] = v
		// use randomizer here
		ent.Properties[v.Name(property.OwnController)].Set(propertyRandomizers[v.Name(property.OwnController)].Random())
	}

	return ent
}
