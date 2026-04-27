package effectapplicator

import (
	"math"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/sirupsen/logrus"
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
func ApplyDirectEffect(logger *logrus.Entry, ent *entity.Entity, eff effect.Effect, target position.Position, cells []position.Position, grd *grid.Grid, targetedEntities []entity.Entity) (damaged []entity.Entity, affected []entity.Entity, credits []rulermethods.CreditAward, err string, errkey string) {
	logger.WithFields(logrus.Fields{}).Info("ApplyDirectEffect")
	credits = []rulermethods.CreditAward{}

	// Hit test!
	if eff.IsDamaging() {
		// @spec-link [[mech_combat_attack_computation]]
		logger.WithFields(logrus.Fields{}).Info("Damaging")

		damageTargets := []entity.Entity{}
		if len(targetedEntities) > 0 {
			accuracy := ent.GetPropertyI(property.Accuracy).I()
			logger.WithFields(logrus.Fields{"accuracy": accuracy}).Info("Accuracy")

			for _, target := range targetedEntities {
				dodge := ent.GetPropertyI(property.Dodge).I()
				hittest := tools.RandomInt(0, 100)
				logger.WithFields(logrus.Fields{"dodge": dodge,
					"hittest": hittest,
					"test":    accuracy - dodge}).Info("Dodge ?")
				if hittest < accuracy-dodge {
					damageTargets = append(damageTargets, target)
				}
			}
		}

		attack := ent.GetPropertyI(property.Attack).I()
		damage := getPropertyOrDefaultI(eff, property.Damage).I()
		shieldPower := getPropertyOrDefaultI(eff, property.ShieldPower).I()
		stunpwr := getPropertyOrDefaultI(eff, property.StunPower).I()
		stunchance := getPropertyOrDefaultI(eff, property.StunChance).I()
		poisonpwr := getPropertyOrDefaultI(eff, property.PoisonPower).I()
		poisonchance := getPropertyOrDefaultI(eff, property.PoisonChance).I()

		critChance := getPropertyOrDefaultI(eff, property.CriticalChance).I()
		critMultiplier := getPropertyOrDefaultI(eff, property.CriticalMultiplier).I()

		logger.WithFields(logrus.Fields{
			"attack":         attack,
			"damage":         damage,
			"shieldPower":    shieldPower,
			"stunpwr":        stunpwr,
			"stunchance":     stunchance,
			"poisonpwr":      poisonpwr,
			"poisonchance":   poisonchance,
			"critChance":     critChance,
			"critMultiplier": critMultiplier,
		}).Info("Attacker")

		totalDamageCredits := 0
		for _, target := range damageTargets {
			hp := target.GetPropertyC(property.HP).GetValue()
			defense := target.GetPropertyI(property.Defense).I()
			shield := target.GetPropertyC(property.Shield).GetValue()
			armor := target.GetPropertyI(property.ArmorRating).I()

			truepoison := tools.Max(poisonpwr-defense, 0)
			truestun := tools.Max(stunpwr-armor, 0)

			multiplier := 1.0
			// roll for crit
			if critChance > 0 && tools.RandomInt(0, 100) < critChance {
				multiplier = float64(critMultiplier) / 100.0
			}

			truedmg := tools.Max((attack*damage/100)-defense-armor, 0) + truepoison + truestun
			truedmg = tools.Max(int(math.Floor(float64(truedmg)*multiplier)), 0)

			logger.WithFields(logrus.Fields{
				"hp":          hp,
				"defense":     defense,
				"shield":      shield,
				"armor":       armor,
				"truepoison":  truepoison,
				"truestun":    truestun,
				"truedmg":     truedmg,
				"shieldPower": shieldPower,
				"multiplier":  multiplier,
			}).Info("Target")

			// apply shield damage (only if negative! otherwise it's healing the shield)
			if shieldPower < 0 {
				target.UpdatePropertyValue(property.Shield, tools.Max(shield+shieldPower, 0))
				shield = tools.Max(shield+shieldPower, 0)
				logger.WithFields(logrus.Fields{
					"shield": shield,
				}).Info("Shield")
			}

			// check poisoning & stunning.
			if truepoison > 0 {
				if tools.RandomInt(0, 100) < poisonchance {
					// target is poisoned
					// update poison value.
					poison := target.GetPropertyI(property.Poison).I()
					// Don't expect poison property to be preset ... this isn't something that sticks between game
					target.RepsertPropertyValue(property.Poison, truepoison+poison)
					logger.WithFields(logrus.Fields{"oldpoison": poison, "newpoison": (poison + truepoison)}).Info("Poisoned")
					
					// Flat rate credit for status effect application
					// Since we don't have the skill weight here easily without the full skill,
					// we'll rely on the caller (UseSkill) to calculate status credits if it's a skill.
					// For direct damage, we only count HP loss.
				}
			}
			if truestun > 0 {
				if tools.RandomInt(0, 100) < stunchance {
					// target is stunned
					// update stun value.
					stun := target.GetPropertyI(property.Stun).I()
					target.RepsertPropertyValue(property.Stun, stun+truestun)
					// Don't expect stun property to be preset ... this isn't something that sticks between game
					logger.WithFields(logrus.Fields{"oldstun": stun, "newstun": (stun + truestun)}).Info("Stunned")
				}
			}
			
			actionDamage := 0
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
					target.UpdatePropertyValue(property.Shield, shield)
					logger.WithFields(logrus.Fields{"shield": shield}).Info("Shield")
				}
				// then kill off hp.
				if hp > 0 {
					if hp > truedmg {
						actionDamage = truedmg
						hp -= truedmg
						truedmg = 0
					} else {
						actionDamage = hp
						truedmg -= hp
						hp = 0
					}
					//update hp
					target.UpdatePropertyValue(property.HP, hp)
					logger.WithFields(logrus.Fields{"hp": hp}).Info("HP")
				}
			}
			totalDamageCredits += actionDamage
			damaged = append(damaged, target)
		}

		if totalDamageCredits > 0 {
			credits = append(credits, rulermethods.CreditAward{
				PlayerID: ent.ControllerID,
				Amount:   totalDamageCredits,
				Source:   "damage",
			})
		}
	}

	if eff.IsHealing() {
		heal := getPropertyOrDefaultI(eff, property.Heal).I()
		shieldPower := getPropertyOrDefaultI(eff, property.ShieldPower).I()
		poisonPower := getPropertyOrDefaultI(eff, property.PoisonPower).I()
		stunPower := getPropertyOrDefaultI(eff, property.StunPower).I()

		totalHealCredits := 0
		for _, target := range targetedEntities {
			hp := target.GetPropertyC(property.HP).GetValue()
			maxhp := target.GetPropertyC(property.HP).GetMaxValue()
			if hp < maxhp {
				actualHeal := 0
				if hp+heal > maxhp {
					actualHeal = maxhp - hp
					hp = maxhp
				} else {
					actualHeal = heal
					hp += heal
				}
				totalHealCredits += actualHeal
				logger.WithFields(logrus.Fields{"hp": hp}).Info("HP")
				target.UpdatePropertyValue(property.HP, hp)
			}

			shield := target.GetPropertyC(property.Shield).GetValue()
			poison := target.GetPropertyI(property.Poison).I()
			stun := target.GetPropertyI(property.Stun).I()

			// Can have overshield.
			if shieldPower > 0 {
				target.UpdatePropertyValue(property.Shield, tools.Min(shield+shieldPower, maxhp*2))
			}
			// ONLY IF NEGATIVE! Otherwise it's a damaging effect
			if poisonPower < 0 {
				target.UpdatePropertyValue(property.Poison, tools.Max(poison+poisonPower, 0))
			}
			// ONLY IF NEGATIVE! Otherwise it's a damaging effect
			if stunPower < 0 {
				target.UpdatePropertyValue(property.Stun, tools.Max(stun+stunPower, 0))
			}

			affected = append(affected, target)
		}

		if totalHealCredits > 0 {
			credits = append(credits, rulermethods.CreditAward{
				PlayerID: ent.ControllerID,
				Amount:   totalHealCredits,
				Source:   "healing",
			})
		}
	}

	return damaged, affected, credits, "", ""
}
