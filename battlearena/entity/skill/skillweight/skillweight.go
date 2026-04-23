package skillweight

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
)

// Calculate computes the positive, negative, and net Skill Weight (SW) for a given skill.
func Calculate(s skill.Skill) (positiveSW int, negativeSW int, netSW int) {
	positiveSW = 0
	negativeSW = 0

	// -- Benefits (Positive SW) --
	
	// Damage Multiplier: +10 SW per 10% (100% damage = +100 SW)
	damageProp := s.GetPropertyI(property.Damage).I()
	if damageProp > 0 {
		positiveSW += damageProp
	}

	// Critical Chance: +2 SW per 1%
	critChance := s.GetPropertyI(property.CriticalChance).I()
	if critChance > 0 {
		positiveSW += critChance * 2
	}

	// Range Extension: +10 SW per cell > 1
	rangeProp := s.GetProperty(property.Range)
	if rp, ok := rangeProp.(*def.RangeProperty); ok {
		if rp.MaxRange > 1 {
			positiveSW += (rp.MaxRange - 1) * 10
		}
	}

	// Zone/AoE: +50 SW per extra cell
	zoneProp := s.GetProperty(property.Zone)
	if zp, ok := zoneProp.(*def.ZoneProperty); ok {
		cellCount := len(zp.ZonePattern)
		if cellCount > 1 {
			positiveSW += (cellCount - 1) * 50
		}
	}

	// Target Anywhere: +40 SW
	targetingMech := s.GetProperty(property.TargetingMechanics)
	if tm, ok := targetingMech.(*def.TargetingMechanicsProperty); ok {
		if tm.TargetingMechanics == def.TargetingMechanicsAnywhere {
			positiveSW += 40
		}
	}

	// Stun Chance: +2 SW per 1%
	stunChance := s.GetPropertyI(property.StunChance).I()
	if stunChance > 0 {
		positiveSW += stunChance * 2
	}

	// Poison Power: +15 SW per 1 dmg/turn
	poisonPower := s.GetPropertyI(property.PoisonPower).I()
	if poisonPower > 0 {
		positiveSW += poisonPower * 15
	}

	// Heal: +15 SW per 10% base
	healProp := s.GetPropertyI(property.Heal).I()
	if healProp > 0 {
		positiveSW += (healProp / 10) * 15
	}

	// Shield: +10 SW per point
	shieldPower := s.GetPropertyI(property.ShieldPower).I()
	if shieldPower > 0 {
		positiveSW += shieldPower * 10
	}

	// -- Payments (Negative SW) --
	
	// Delay (+100): -100 SW (baseline). Extra Delay: -10 SW per +10 delay
	// i.e., -1 SW per 1 Delay
	delayProp := s.GetPropertyC(property.Delay).GetMaxValue()
	if delayProp > 0 {
		negativeSW -= delayProp
	}

	// Channeling: -15 SW per 10 delay
	channelingProp := s.GetPropertyC(property.Channeling).GetMaxValue()
	if channelingProp > 0 {
		negativeSW -= (channelingProp / 10) * 15
	}

	// MP Leech: -15 SW per 1 MP
	mpLeech := s.GetPropertyI(property.MPLeech).I()
	if mpLeech > 0 {
		negativeSW -= mpLeech * 15
	}

	// SP Leech: -10 SW per 1 SP
	spLeech := s.GetPropertyI(property.SPLeech).I()
	if spLeech > 0 {
		negativeSW -= spLeech * 10
	}

	// HP Leech: -20 SW per 1 HP
	hpLeech := s.GetPropertyI(property.HPLeech).I()
	if hpLeech > 0 {
		negativeSW -= hpLeech * 20
	}

	// Cooldown: -25 SW per turn
	cooldown := s.GetPropertyC(property.Cooldown).GetMaxValue()
	if cooldown > 0 {
		negativeSW -= cooldown * 25
	}

	return positiveSW, negativeSW, positiveSW + negativeSW
}

// GetGrade returns the skill grade (I-V) based on the Positive SW.
func GetGrade(positiveSW int) string {
	if positiveSW <= 150 {
		return "Grade I"
	} else if positiveSW <= 300 {
		return "Grade II"
	} else if positiveSW <= 500 {
		return "Grade III"
	} else if positiveSW <= 750 {
		return "Grade IV"
	}
	return "Grade V"
}

// GetCreditCost returns the cost of the skill based on its Positive SW.
func GetCreditCost(positiveSW int) int {
	return positiveSW * 2
}
