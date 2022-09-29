package skillgenerator

import (
	"log"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
	"github.com/ecumeurs/upsilontools/tools"
)

var propertiesTargetingRandomizers = []func() property.Property{
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.Accuracy, tools.RandomInt(50, 150), property.Public, property.Skill)
	},
}
var propertiesEffectRandomizers = []func() property.Property{
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.Damage, tools.RandomInt(1, 3), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.Heal, tools.RandomInt(1, 3), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.ShieldPower, tools.RandomInt(1, 3), property.Public, property.Skill)
	},
}
var propertiesCostRandomizers = []func() property.Property{
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.Cooldown, tools.RandomInt(1, 3), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.HPLeech, tools.RandomInt(1, 3), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.MPLeech, tools.RandomInt(1, 3), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.SPLeech, tools.RandomInt(1, 3), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.Delay, tools.RandomInt(500, 1000), property.Public, property.Skill)
	},
}

func GenerateRandomSkill() *skill.Skill {
	sk := skill.New()

	for _, v := range propertiesTargetingRandomizers {
		if tools.RandomInt(0, 100) > 50 {
			sk.Targeting = append(sk.Targeting, v())
		}
	}
	for len(sk.Effect) == 0 {
		for _, v := range propertiesEffectRandomizers {
			if tools.RandomInt(0, 100) > 50 {
				skp := v()
				log.Println(skp)
				sk.Effect = append(sk.Effect, skp)
				log.Println("Effect", sk.Effect)
				sk.Name = sk.Effect[0].Name(property.GameMaster)
				break // only one effect for now.
			}
		}
	}
	// might have multiple costs... for fun.
	for _, v := range propertiesCostRandomizers {
		if tools.RandomInt(0, 100) > 50 {
			sk.Cost = append(sk.Cost, v())
		}
	}

	return sk
}
