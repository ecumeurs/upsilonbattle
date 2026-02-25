package entitygenerator

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/sirupsen/logrus"

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
	ent := entity.New()

	for _, v := range def.PropertiesForCharacter() {
		propName := v.Name(property.OwnController)
		ent.Properties[propName] = v
		// Only randomize if we have a defined range for this property
		if r, ok := propertyRandomizers[propName]; ok {
			ent.Properties[propName].Set(r.Random())
		} else {
			logrus.Warnf("entitygenerator: Missing propertyRandomizer for '%s', using default value", propName)
		}
	}

	return ent
}
