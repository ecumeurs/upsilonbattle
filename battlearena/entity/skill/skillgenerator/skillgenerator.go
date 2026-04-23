package skillgenerator

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill/skillweight"
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
		return defaultproperty.MakeIntProperty(property.Damage, tools.RandomInt(50, 200), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.Heal, tools.RandomInt(50, 150), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.ShieldPower, tools.RandomInt(10, 50), property.Public, property.Skill)
	},
}

var propertiesCostRandomizers = []func() property.Property{
	func() property.Property {
		return defaultproperty.MakeIntCounterProperty(property.Cooldown, 0, tools.RandomInt(1, 5), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.HPLeech, tools.RandomInt(1, 10), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.MPLeech, tools.RandomInt(1, 10), property.Public, property.Skill)
	},
	func() property.Property {
		return defaultproperty.MakeIntProperty(property.SPLeech, tools.RandomInt(1, 10), property.Public, property.Skill)
	},
}

func GenerateRandomSkill() skill.Skill {
	sk := skill.New()

	for _, v := range propertiesTargetingRandomizers {
		if tools.RandomInt(0, 100) > 50 {
			sk.Targeting[v().Name(property.GameMaster)] = v()
		}
	}
	for len(sk.Effect.Properties) == 0 {
		for _, v := range propertiesEffectRandomizers {
			if tools.RandomInt(0, 100) > 50 {
				skp := v()
				sk.Effect.Properties = append(sk.Effect.Properties, skp)
				sk.Name = sk.Effect.Properties[0].Name(property.GameMaster)
				break // only one effect for now.
			}
		}
	}
	
	// Add some random costs
	for _, v := range propertiesCostRandomizers {
		if tools.RandomInt(0, 100) > 50 {
			skp := v()
			sk.Costs[skp.Name(property.GameMaster)] = skp
		}
	}

	// Balance skill using Skill Weight
	_, _, netSW := skillweight.Calculate(sk)
	
	currentDelay := sk.GetPropertyC(property.Delay).GetMaxValue()
	newDelay := currentDelay + netSW
	if newDelay < 0 {
		extraDamage := -newDelay
		
		damageVal := sk.GetPropertyI(property.Damage).I()
		
		var damageProp property.IntProperty
		for _, p := range sk.Effect.Properties {
			if p.Name(property.GameMaster) == property.Damage.String() {
				damageProp = p.(property.IntProperty)
				break
			}
		}
		if damageProp != nil {
			damageProp.SetI(damageProp.I() + extraDamage)
		} else {
			sk.Effect.Properties = append(sk.Effect.Properties, defaultproperty.MakeIntProperty(property.Damage, damageVal + extraDamage, property.Public, property.Skill))
		}
		
		newDelay = 0
	}
	
	sk.Costs[property.Delay.String()] = defaultproperty.MakeIntCounterProperty(property.Delay, 0, newDelay, property.Public, property.Skill)

	return sk
}
