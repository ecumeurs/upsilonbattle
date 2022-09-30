package property

import "fmt"

type EntityProperties string

const (
	HP       EntityProperties = "HP"       // Absence means 10 /!\ this is a counter (has two value: current, max)
	Movement EntityProperties = "Movement" // Absence means 5 /!\ this is a counter (has two value: current, max), reset at end of turn
	SP       EntityProperties = "SP"       // Absence means 10
	MP       EntityProperties = "MP"       // Absence means 10

	Attack     EntityProperties = "Attack"     // Absence means 1, basic attack
	Defense    EntityProperties = "Defense"    // Absence means 0, basic defense
	JumpHeight EntityProperties = "JumpHeight" // Absence means 2

	// Flags

	IsDying     EntityProperties = "IsDying"     // Absence means false
	HasMoved    EntityProperties = "HasMoved"    // Absence means false
	HasAttacked EntityProperties = "HasAttacked" // Absence means false

	// Buffs and other applied, inherited properties ... ?
	AttackRange EntityProperties = "AttackRange" // Absence means 1, Basic attack range. Altered with items, mostly.

	// Computed at end of turn.
	Shield EntityProperties = "Shield" // Absence means 0,0 (counter)
	Poison EntityProperties = "Poison" // Absence means 0, Poisoned state. Sum of all poison damage taken per turn. Mostly a temporary status. Negative PoisonPower will cure; /2 each turn, remove at <1
	Stun   EntityProperties = "Stun"   // Absence means 0, Stunned state. Sum of all stun taken per turn. Mostly a temporary status. Negative StunPower will cure; /2 each turn, remove at <1
)

// String
func (e EntityProperties) String() string {
	return string(e)
}

type SkillProperties string

const (
	Behavior SkillProperties = "Behavior" // property.Skill broad category: Direct, Reaction, Passive, Counter, Trap; Absence means Direct

	Range              SkillProperties = "Range"              // Range of the property.Skill // Absence means 1
	Zone               SkillProperties = "Zone"               // Area of Effect // Absence means 1 tile effect
	TargetNumber       SkillProperties = "TargetNumber"       // Absence means all targets within the targeted zone.
	Accuracy           SkillProperties = "Accuracy"           // Absence means 100%
	Dodge              SkillProperties = "Dodge"              // Absence means 0%
	Parry              SkillProperties = "Parry"              // Absence means 0%
	TargetType         SkillProperties = "TargetType"         // Entity, Tile, Both, Self
	TargetingMechanics SkillProperties = "TargetingMechanics" // Anywhere, Line of Sight, and maybe other mechanics later.

	Damage         SkillProperties = "Damage"         // Absence means 0
	Heal           SkillProperties = "Heal"           // Absence means 0
	ShieldPower    SkillProperties = "Shield"         // Absence means 0 , can be negative or positive.
	StunPower      SkillProperties = "Stun"           // Absence means 0 , can be negative or positive.
	StunChance     SkillProperties = "StunChance"     // Absence means 0%
	CriticalChance SkillProperties = "CriticalChance" // Absence means 0%
	CriticalDamage SkillProperties = "CriticalDamage" // Absence means 0%
	Duration       SkillProperties = "Duration"       // Absence means 0
	PoisonPower    SkillProperties = "Poison"         // Absence means 0 , can be negative or positive.
	PoisonChance   SkillProperties = "PoisonChance"   // Absence means 0%

	Delay   SkillProperties = "Delay"   // Absence means 500
	HPLeech SkillProperties = "HPLeech" // Absence means 0
	MPLeech SkillProperties = "MPLeech" // Absence means 0
	SPLeech SkillProperties = "SPLeech" // Absence means 0
	MvtCost SkillProperties = "MvtCost" // Absence means 0

	Cooldown SkillProperties = "Cooldown" // Absence means 3 turns. Special note: Cool down is stored as a counter, minValue represent initial cooldown at battle start. MaxValue represent the cooldown value when used.

)

// String
func (sp SkillProperties) String() string {
	return string(sp)
}

var SkillTargetingProperties = map[SkillProperties]bool{
	Range:              true,
	Zone:               true,
	TargetNumber:       true,
	Accuracy:           true,
	Dodge:              true,
	Parry:              true,
	TargetType:         true,
	TargetingMechanics: true,
}

var SkillEffectProperties = map[SkillProperties]bool{
	Damage:         true,
	Heal:           true,
	ShieldPower:    true,
	StunPower:      true,
	StunChance:     true,
	CriticalChance: true,
	CriticalDamage: true,
	Duration:       true,
	PoisonPower:    true,
	PoisonChance:   true,
}

var SkillCostProperties = map[SkillProperties]bool{
	Delay:    true,
	HPLeech:  true,
	MPLeech:  true,
	SPLeech:  true,
	Cooldown: true,
}

type ItemProperties string

const (
	Durability       ItemProperties = "Durability"       // Absence means 0: invulnerable
	Weight           ItemProperties = "Weight"           // Absence means 0: no weight
	ItemType         ItemProperties = "ItemType"         // Absence means None (out of Wearable, Consumable, Usable, Throwable, Ammunitions and None)
	ArmorRating      ItemProperties = "Armor"            // Absence means 0: no armor (only for Wearable)
	WeaponType       ItemProperties = "WeaponType"       // Absence means 0: no weapon type (only for Wearable)
	ArmorType        ItemProperties = "ArmorType"        // Absence means 0: no armor type (only for Wearable)
	ToolType         ItemProperties = "ToolType"         // Absence means 0: no tool type (only for Wearable)
	WeaponRange      ItemProperties = "WeaponRange"      // Absence means 0: no weapon range (only for Wearable)
	WeaponBaseDamage ItemProperties = "WeaponBaseDamage" // Absence means 0: no weapon base damage (only for Wearable)
	Stackable        ItemProperties = "Stackable"        // Absence means 0: not stackable
	StackSize        ItemProperties = "StackSize"        // Absence means 0: no stack size
	Effect           ItemProperties = "Effect"           // Absence means nil: No effect. Effects are Skills. (except None)
	Value            ItemProperties = "Value"            // Absence means 0: no value
)

// String
func (ip ItemProperties) String() string {
	return string(ip)
}

// PropertyToString
func PropertyToString(p interface{}) string {
	switch pconv := p.(type) {
	case EntityProperties:
		return pconv.String()
	case SkillProperties:
		return pconv.String()
	case ItemProperties:
		return pconv.String()
	case string:
		return pconv
	default:
		// Abort
		panic(fmt.Sprintf("PropertyToString: Unknown property type: %T", p))
	}
}
