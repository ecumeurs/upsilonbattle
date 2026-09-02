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

// @spec-link [[mech_combat_attack_computation]]
// @spec-link [[mechanic_ai_termination]]

// getPropertyOrDefault retrieves a property from the provided effect. 
// If the property is not found, it returns the global default value for that property type.
// This ensures that effect application is resilient to missing optional properties in skill definitions.
func getPropertyOrDefault(eff effect.Effect, p property.Key) property.Property {
	res := eff.GetProperty(p)
	if res == nil {
		res = def.DefaultProperty(p)
	}
	return res
}

// getPropertyOrDefaultI is a typed helper that returns an IntProperty from an effect.
// It uses getPropertyOrDefault and performs a type assertion to property.IntProperty.
// Use this when you are certain the property exists or has a valid integer default.
func getPropertyOrDefaultI(eff effect.Effect, p property.Key) property.IntProperty {
	return getPropertyOrDefault(eff, p).(property.IntProperty)
}

// getPropertyOrDefaultF is a typed helper that returns a FloatProperty from an effect.
// It uses getPropertyOrDefault and performs a type assertion to property.FloatProperty.
// Use this for properties like multipliers or rates where precision is required.
func getPropertyOrDefaultF(eff effect.Effect, p property.Key) property.FloatProperty {
	return getPropertyOrDefault(eff, p).(property.FloatProperty)
}

// getPropertyOrDefaultC is a typed helper that returns an IntCounterProperty from an effect.
// It uses getPropertyOrDefault and performs a type assertion to property.IntCounterProperty.
// Useful for properties that track current/max values like HP, SP, or Shield.
func getPropertyOrDefaultC(eff effect.Effect, p property.Key) property.IntCounterProperty {
	return getPropertyOrDefault(eff, p).(property.IntCounterProperty)
}

// ApplyDirectEffect is the main entry point for applying an effect's results to a list of entities.
// It separates the logic into damaging and healing branches, executing both if the effect qualifies.
// It returns detailed ActionResults for both damage and healing, along with any credits awarded to the actor.
func ApplyDirectEffect(logger *logrus.Entry, ent *entity.Entity, eff effect.Effect, target position.Position, cells []position.Position, grd *grid.Grid, targetedEntities []entity.Entity) (damaged []rulermethods.ActionResult, affected []rulermethods.ActionResult, credits []rulermethods.CreditAward, err string, errkey string) {
	logger.Info("ApplyDirectEffect: Starting application of tactical effect")
	credits = []rulermethods.CreditAward{}

	// Process damaging component of the effect
	if eff.IsDamaging() {
		logger.Debug("Effect has damaging properties, executing damaging branch")
		dmg, creds := applyDamagingEffect(logger, ent, eff, targetedEntities)
		damaged = append(damaged, dmg...)
		credits = append(credits, creds...)
	}

	// Process healing component of the effect
	if eff.IsHealing() {
		logger.Debug("Effect has healing properties, executing healing branch")
		aff, creds := applyHealingEffect(logger, ent, eff, targetedEntities)
		affected = append(affected, aff...)
		credits = append(credits, creds...)
	}

	return damaged, affected, credits, "", ""
}

// applyDamagingEffect handles the multi-target damage distribution logic.
// It first performs an accuracy/dodge check for every target, then calculates and applies damage.
func applyDamagingEffect(logger *logrus.Entry, ent *entity.Entity, eff effect.Effect, targetedEntities []entity.Entity) ([]rulermethods.ActionResult, []rulermethods.CreditAward) {
	logger.Info("applyDamagingEffect: Processing damage for targets")
	damaged := []rulermethods.ActionResult{}
	credits := []rulermethods.CreditAward{}

	// 1. Perform Hit Test for all entities in range
	accuracy := ent.GetPropertyI(property.Accuracy).I()
	damageTargets := []entity.Entity{}
	for _, target := range targetedEntities {
		// Dodge is the TARGET's evasion, not the caster's — an entity's own
		// Dodge must never suppress its own attacks. ISS-145.
		dodge := target.GetPropertyI(property.Dodge).I()
		if tools.RandomInt(0, 100) < accuracy-dodge {
			damageTargets = append(damageTargets, target)
		} else {
			logger.WithField("targetID", target.ID).Debug("Target dodged the effect")
		}
	}

	// 2. Resolve damage for each hit target
	totalDamageCredits := 0
	for _, target := range damageTargets {
		res := applyDamageToSingleTarget(logger, ent, eff, target)
		totalDamageCredits += res.Damage
		damaged = append(damaged, res)
	}

	// 3. Award credits based on total HP loss caused
	if totalDamageCredits > 0 {
		credits = append(credits, rulermethods.CreditAward{
			PlayerID: ent.ControllerID,
			Amount:   totalDamageCredits,
			Source:   "damage",
		})
	}
	return damaged, credits
}

// applyDamageToSingleTarget computes and applies the specific damage values to one entity.
// It handles crits, defense/armor reduction, shield depletion, and status effect application.
func applyDamageToSingleTarget(logger *logrus.Entry, ent *entity.Entity, eff effect.Effect, target entity.Entity) rulermethods.ActionResult {
	attack := ent.GetPropertyI(property.Attack).I()
	damage := getPropertyOrDefaultI(eff, property.DamageScale).I()
	shieldPower := getPropertyOrDefaultI(eff, property.ShieldPower).I()
	stunpwr := getPropertyOrDefaultI(eff, property.StunPower).I()
	stunchance := getPropertyOrDefaultI(eff, property.StunChance).I()
	poisonpwr := getPropertyOrDefaultI(eff, property.PoisonPower).I()
	poisonchance := getPropertyOrDefaultI(eff, property.PoisonChance).I()
	critChance := getPropertyOrDefaultI(eff, property.CriticalChance).I()
	critMultiplier := getPropertyOrDefaultI(eff, property.CriticalMultiplier).I()

	hp := target.GetPropertyC(property.HP).GetValue()
	defense := target.GetPropertyI(property.Defense).I()
	shield := target.GetPropertyC(property.Shield).GetValue()
	armor := target.GetPropertyI(property.ArmorRating).I()

	// Roll for Critical Hit
	multiplier := 1.0
	if critChance > 0 && tools.RandomInt(0, 100) < critChance {
		multiplier = float64(critMultiplier) / 100.0
		logger.Debug("Critical hit triggered!")
	}

	// Calculate base raw damage after defense/armor reduction
	truepoison := tools.Max(poisonpwr-defense, 0)
	truestun := tools.Max(stunpwr-armor, 0)
	truedmg := tools.Max((attack*damage/100)-defense-armor, 0) + truepoison + truestun
	truedmg = tools.Max(int(math.Floor(float64(truedmg)*multiplier)), 0)

	prevHP := hp
	// Deplete shield if the effect has shield-damaging properties.
	// Shield is buffable: persist the change as a base-level delta (never the
	// composed absolute), so an active buff's contribution is never folded
	// into base.
	// @spec-link [[rule_entity_property_write_isolation]]
	if shieldPower < 0 {
		newShield := tools.Max(shield+shieldPower, 0)
		target.AdjustPropertyCValue(property.Shield, newShield-shield)
		shield = newShield
	}

	// Roll for and apply Poison status effect
	if truepoison > 0 && tools.RandomInt(0, 100) < poisonchance {
		poison := target.GetPropertyI(property.Poison).I()
		target.RepsertPropertyValue(property.Poison, truepoison+poison)
	}
	// Roll for and apply Stun status effect
	if truestun > 0 && tools.RandomInt(0, 100) < stunchance {
		stun := target.GetPropertyI(property.Stun).I()
		target.RepsertPropertyValue(property.Stun, stun+truestun)
	}

	// Distribute remaining damage between Shield and HP. Both are buffable:
	// persist each change as a base-level delta, never a composed absolute.
	// @spec-link [[rule_entity_property_write_isolation]]
	actionDamage := 0
	if truedmg > 0 {
		if shield > 0 {
			shieldBefore := shield
			if shield > truedmg {
				shield -= truedmg
				truedmg = 0
			} else {
				truedmg -= shield
				shield = 0
			}
			target.AdjustPropertyCValue(property.Shield, shield-shieldBefore)
		}
		if hp > 0 {
			hpBefore := hp
			if hp > truedmg {
				actionDamage = truedmg
				hp -= truedmg
			} else {
				actionDamage = hp
				hp = 0
			}
			target.AdjustPropertyCValue(property.HP, hp-hpBefore)
		}
	}

	return rulermethods.ActionResult{
		Target: target, TargetID: target.ID, Damage: actionDamage, PrevHP: prevHP, NewHP: hp,
	}
}

// applyHealingEffect handles the distribution of healing and cleansing to target entities.
func applyHealingEffect(logger *logrus.Entry, ent *entity.Entity, eff effect.Effect, targetedEntities []entity.Entity) ([]rulermethods.ActionResult, []rulermethods.CreditAward) {
	logger.Info("applyHealingEffect: Processing healing for targets")
	affected := []rulermethods.ActionResult{}
	credits := []rulermethods.CreditAward{}

	totalHealCredits := 0
	for _, target := range targetedEntities {
		res := applyHealToSingleTarget(logger, eff, target)
		totalHealCredits += res.Heal
		affected = append(affected, res)
	}

	// Award credits based on total HP restored
	if totalHealCredits > 0 {
		credits = append(credits, rulermethods.CreditAward{
			PlayerID: ent.ControllerID,
			Amount:   totalHealCredits,
			Source:   "healing",
		})
	}
	return affected, credits
}

// applyHealToSingleTarget computes and applies healing, overshield, and cleansing to one entity.
func applyHealToSingleTarget(logger *logrus.Entry, eff effect.Effect, target entity.Entity) rulermethods.ActionResult {
	heal := getPropertyOrDefaultI(eff, property.Heal).I()
	shieldPower := getPropertyOrDefaultI(eff, property.ShieldPower).I()
	poisonPower := getPropertyOrDefaultI(eff, property.PoisonPower).I()
	stunPower := getPropertyOrDefaultI(eff, property.StunPower).I()

	hp := target.GetPropertyC(property.HP).GetValue()
	maxhp := target.GetPropertyC(property.HP).GetMaxValue()
	prevHP := hp
	actualHeal := 0

	// Apply HP restoration up to maximum capacity. HP is buffable: persist the
	// change as a base-level delta (never the composed absolute), so an
	// active buff's contribution is never folded into base.
	// @spec-link [[rule_entity_property_write_isolation]]
	if hp < maxhp {
		hpBefore := hp
		if hp+heal > maxhp {
			actualHeal = maxhp - hp
			hp = maxhp
		} else {
			actualHeal = heal
			hp += heal
		}
		target.AdjustPropertyCValue(property.HP, hp-hpBefore)
	}

	shield := target.GetPropertyC(property.Shield).GetValue()
	poison := target.GetPropertyI(property.Poison).I()
	stun := target.GetPropertyI(property.Stun).I()

	// Apply overshield (limit set to 2x max HP). Shield is buffable: persist
	// as a base-level delta, same as above.
	// @spec-link [[rule_entity_property_write_isolation]]
	if shieldPower > 0 {
		newShield := tools.Min(shield+shieldPower, maxhp*2)
		target.AdjustPropertyCValue(property.Shield, newShield-shield)
	}
	// Apply status effect cleansing (only if power values are negative)
	if poisonPower < 0 {
		target.UpdatePropertyValue(property.Poison, tools.Max(poison+poisonPower, 0))
	}
	if stunPower < 0 {
		target.UpdatePropertyValue(property.Stun, tools.Max(stun+stunPower, 0))
	}

	return rulermethods.ActionResult{
		Target: target, TargetID: target.ID, Heal: actualHeal, PrevHP: prevHP, NewHP: hp,
	}
}
