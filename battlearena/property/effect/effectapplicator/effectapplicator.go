package effectapplicator

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect"
	"github.com/ecumeurs/upsilontools/tools"
)

// GetPropertyOrDefault(eff Effect, p interface{}) property.Property
func getPropertyOrDefault(eff effect.Effect, p interface{}) property.Property {
	res := eff.GetProperty(p)
	if res == nil {
		res = def.DefaultProperty(p)
	}
	return res
}

func getPropertyOrDefaultI(eff effect.Effect, p interface{}) property.IntProperty {
	return getPropertyOrDefault(eff, p).(property.IntProperty)
}

func getPropertyOrDefaultF(eff effect.Effect, p interface{}) property.FloatProperty {
	return getPropertyOrDefault(eff, p).(property.FloatProperty)
}

func getPropertyOrDefaultC(eff effect.Effect, p interface{}) property.IntCounterProperty {
	return getPropertyOrDefault(eff, p).(property.IntCounterProperty)
}

// Machine that apply effects
func ApplyDirectEffect(ent *entity.Entity, eff effect.Effect, pos []position.Position, grd *grid.Grid, targetedEntities []entity.Entity) (damaged []entity.Entity, affected []entity.Entity, err string, errkey string) {

	// Hit test!
	if eff.IsDamaging() {

		damageTargets := []entity.Entity{}
		if len(targetedEntities) > 0 {
			accuracy := ent.GetPropertyI(property.Accuracy).I()

			for _, target := range targetedEntities {
				dodge := ent.GetPropertyI(property.Dodge).I()
				if tools.RandomInt(0, 100) < accuracy-dodge {
					damageTargets = append(damageTargets, target)
				}
			}
		}

		damage := getPropertyOrDefaultI(eff, property.Damage).I()
		stunpwr := getPropertyOrDefaultI(eff, property.StunPower).I()
		stunchance := getPropertyOrDefaultI(eff, property.StunChance).I()
		poisonpwr := getPropertyOrDefaultI(eff, property.PoisonPower).I()
		poisonchance := getPropertyOrDefaultI(eff, property.PoisonChance).I()

		for _, target := range damageTargets {
			hp := target.GetPropertyC(property.HP).GetValue()
			defense := target.GetPropertyI(property.Defense).I()
			shield := target.GetPropertyC(property.Shield).GetValue()
			armor := target.GetPropertyI(property.ArmorRating).I()

			truepoison := poisonpwr - defense
			truestun := stunpwr - armor

			truedmg := damage - defense - armor + truepoison + truestun

			// check poisoning & stunning.
			if truepoison > 0 {
				if tools.RandomInt(0, 100) < poisonchance {
					// target is poisoned
					// update poison value.
					poison := target.GetPropertyI(property.Poison).I()
					target.GetPropertyI(property.Poison).Set(poison + truepoison)
				}
			}
			if truestun > 0 {
				if tools.RandomInt(0, 100) < stunchance {
					// target is stunned
					// update stun value.
					stun := target.GetPropertyI(property.Stun).I()
					target.GetPropertyI(property.Stun).Set(stun + truestun)
				}
			}
			if truedmg > 0 {
				// first killoff shield.
				if shield > 0 {
					if shield > truedmg {
						shield -= truedmg
						truedmg = 0
					} else {
						truedmg -= shield
						shield = 0
					}
					//update shield
					target.GetPropertyC(property.Shield).SetValue(shield)
				}
				// then kill off hp.
				if hp > 0 {
					if hp > truedmg {
						hp -= truedmg
						truedmg = 0
					} else {
						truedmg -= hp
						hp = 0
					}
					//update hp
					target.GetPropertyC(property.HP).SetValue(hp)
				}
			}
			damaged = append(damaged, target)
		}

	}

	if eff.IsHealing() {
		heal := getPropertyOrDefaultI(eff, property.Heal).I()
		for _, target := range targetedEntities {
			hp := target.GetPropertyC(property.HP).GetValue()
			maxhp := target.GetPropertyC(property.HP).GetMaxValue()
			if hp < maxhp {
				if hp+heal > maxhp {
					hp = maxhp
				} else {
					hp += heal
				}
				target.GetPropertyC(property.HP).SetValue(hp)
			}
			affected = append(affected, target)
		}

	}

	return damaged, affected, "", ""
}
