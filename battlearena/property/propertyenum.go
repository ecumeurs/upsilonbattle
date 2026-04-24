package property

import "fmt"

type EntityProperties string

const (
	// all these are counters
	HP       EntityProperties = "HP"       // Absence means 10
	Movement EntityProperties = "Movement" // Absence means 3; reset at end of turn
	SP       EntityProperties = "SP"       // Absence means 10
	MP       EntityProperties = "MP"       // Absence means 10

	Attack     EntityProperties = "Attack"     // Absence means 1, basic attack
	Defense    EntityProperties = "Defense"    // Absence means 0, basic defense
	JumpHeight EntityProperties = "JumpHeight" // Absence means 2
	TeamID     EntityProperties = "TeamID"     // Team affiliation

	// Flags
	// Absence means -1 (not dying), when HP reaches 0, IsDying is set to 3, and at each end of entity's turn, will be reduced by one. When it reaches 0, the entity is removed from game.
	IsDying  EntityProperties = "IsDying"
	HasMoved EntityProperties = "HasMoved" // Absence means false, can move multiple times up to acting, once acted, can't move. Trap triggered or reaction can also mark the target as having moved.
	HasActed EntityProperties = "HasActed" // Absence means false, will be set after the target has attacked or used a skill.

	// Buffs and other applied, inherited properties ... ?
	AttackRange EntityProperties = "AttackRange" // Absence means 1, Basic attack range. Altered with items, mostly.

	// @spec-link [[mech_entity_expiration]]
	// EntityDuration tracks how many turns a temporary entity lives. 0 means permanent.
	// NOTE: distinct from SkillProperties.Duration which is used for buff duration.
	EntityDuration      EntityProperties = "EntityDuration"
	// ExpiresWithCaster: if true, entity's owned positional effects are removed when their caster entity dies.
	// @spec-link [[mechanic_effect_caster_tracking]]
	ExpiresWithCaster   EntityProperties = "ExpiresWithCaster"
	// WalkThrough: if true, other entities may occupy the same cell as this entity.
	// @spec-link [[mechanic_multi_entity_cell_system]]
	WalkThrough         EntityProperties = "WalkThrough"
	// Invisible: if true, the entity is not sent in client-facing state snapshots.
	Invisible           EntityProperties = "Invisible"
	// AIBehavior stores the slug for the automated behavior (e.g. "aggressive", "spooked").
	// @spec-link [[mech_behavior_interface]]
	AIBehavior          EntityProperties = "AIBehavior"

	// (counters) Absence means 0,0, can have overshield (when applied through healing and buffs...) Max shield is only used at initialisation of battle. Allowed to twice HP Max
	Shield EntityProperties = "Shield"
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

	Damage             SkillProperties = "Damage"             // Absence means 100%
	Heal               SkillProperties = "Heal"               // Absence means 0
	ShieldPower        SkillProperties = "Shield"             // Absence means 0 , can be negative or positive.
	StunPower          SkillProperties = "Stun"               // Absence means 0 , can be negative or positive.
	StunChance         SkillProperties = "StunChance"         // Absence means 0%
	CriticalChance     SkillProperties = "CriticalChance"     // Absence means 0%
	CriticalMultiplier SkillProperties = "CriticalMultiplier" // Absence means 0%
	Duration           SkillProperties = "Duration"           // Absence means 0
	PoisonPower        SkillProperties = "Poison"             // Absence means 0 , can be negative or positive.
	PoisonChance       SkillProperties = "PoisonChance"       // Absence means 0%

	Delay      SkillProperties = "Delay"      // Absence means 500
	Channeling SkillProperties = "Channeling" // Absence means 0
	HPLeech SkillProperties = "HPLeech" // Absence means 0
	MPLeech SkillProperties = "MPLeech" // Absence means 0
	SPLeech SkillProperties = "SPLeech" // Absence means 0
	MvtCost SkillProperties = "MvtCost" // Absence means 0

	Cooldown SkillProperties = "Cooldown" // Absence means 3 turns. Special note: Cool down is stored as a counter, minValue represent initial cooldown at battle start. MaxValue represent the cooldown value when used.

	// @spec-link [[mech_trigger_system]]
	// TriggerType defines when a positional effect fires. Value is a TriggerTypeValue string.
	TriggerType     SkillProperties = "TriggerType"
	// RemoveOnTrigger: if true, the positional effect is consumed after firing once.
	RemoveOnTrigger SkillProperties = "RemoveOnTrigger"
	// TriggerCount: how many times the effect can fire (0 = unlimited).
	TriggerCount    SkillProperties = "TriggerCount"

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
	Damage:             true,
	Heal:               true,
	ShieldPower:        true,
	StunPower:          true,
	StunChance:         true,
	CriticalChance:     true,
	CriticalMultiplier: true,
	Duration:           true,
	PoisonPower:        true,
	PoisonChance:       true,
}

var SkillCostProperties = map[SkillProperties]bool{
	Delay:      true,
	Channeling: true,
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
