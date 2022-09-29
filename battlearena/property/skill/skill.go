package skill

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position/pattern"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
)

type SkillProperties string

const (
	Behavior SkillProperties = "Behavior" // property.Skill broad category: Direct, Reaction, Passive, Counter; Absence means Direct

	Range              SkillProperties = "Range"              // Range of the property.Skill // Absence means 1
	Zone               SkillProperties = "Zone"               // Area of Effect // Absence means 1 tile effect
	TargetNumber       SkillProperties = "TargetNumber"       // Absence means all targets within the targeted zone.
	Accuracy           SkillProperties = "Accuracy"           // Absence means 100%
	Dodge              SkillProperties = "Dodge"              // Absence means 0%
	Parry              SkillProperties = "Parry"              // Absence means 0%
	TargetType         SkillProperties = "TargetType"         // Entity, Tile, Both, Self
	TargetingMechanics SkillProperties = "TargetingMechanics" // Anywhere, Line of Sight, and maybe other mechanics later.

	Damage          SkillProperties = "Damage"          // Absence means 0
	Heal            SkillProperties = "Heal"            // Absence means 0
	ShieldPower     SkillProperties = "Shield"          // Absence means 0 , can be negative or positive.
	StunPower       SkillProperties = "Stun"            // Absence means 0 , can be negative or positive.
	StunChance      SkillProperties = "StunChance"      // Absence means 0%
	CriticialChance SkillProperties = "CriticialChance" // Absence means 0%
	CriticialDamage SkillProperties = "CriticialDamage" // Absence means 0%
	Duration        SkillProperties = "Duration"        // Absence means 0
	PoisonPower     SkillProperties = "Poison"          // Absence means 0 , can be negative or positive.
	PoisonChance    SkillProperties = "PoisonChance"    // Absence means 0%

	Delay    SkillProperties = "Delay"    // Absence means 500
	HPLeech  SkillProperties = "HPLeech"  // Absence means 0
	MPLeech  SkillProperties = "MPLeech"  // Absence means 0
	SPLeech  SkillProperties = "SPLeech"  // Absence means 0
	Cooldown SkillProperties = "Cooldown" // Absence means 3 turns
)

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
	Damage:          true,
	Heal:            true,
	ShieldPower:     true,
	StunPower:       true,
	StunChance:      true,
	CriticialChance: true,
	CriticialDamage: true,
	Duration:        true,
	PoisonPower:     true,
	PoisonChance:    true,
}

var SkillCostProperties = map[SkillProperties]bool{
	Delay:    true,
	HPLeech:  true,
	MPLeech:  true,
	SPLeech:  true,
	Cooldown: true,
}

// Prepare default Properties.

func DefaultTargetNumber() *property.DefaultIntProperty {
	return property.MakeIntProperty(TargetNumber, 0, property.FriendlyController, property.Skill)
}

func DefaultAccuracy() *property.DefaultIntProperty {
	return property.MakeIntProperty(Accuracy, 100, property.FriendlyController, property.Skill)
}

func DefaultDodge() *property.DefaultIntProperty {
	return property.MakeIntProperty(Dodge, 0, property.FriendlyController, property.Skill)
}

func DefaultParry() *property.DefaultIntProperty {
	return property.MakeIntProperty(Parry, 0, property.FriendlyController, property.Skill)
}

func DefaultDamage() *property.DefaultIntProperty {
	return property.MakeIntProperty(Damage, 0, property.FriendlyController, property.Skill)
}

func DefaultHeal() *property.DefaultIntProperty {
	return property.MakeIntProperty(Heal, 0, property.FriendlyController, property.Skill)
}

func DefaultShieldPower() *property.DefaultIntProperty {
	return property.MakeIntProperty(ShieldPower, 0, property.FriendlyController, property.Skill)
}

func DefaultStunPower() *property.DefaultIntProperty {
	return property.MakeIntProperty(StunPower, 0, property.FriendlyController, property.Skill)
}

func DefaultStunChance() *property.DefaultIntProperty {
	return property.MakeIntProperty(StunChance, 0, property.FriendlyController, property.Skill)
}

func DefaultCriticialChance() *property.DefaultIntProperty {
	return property.MakeIntProperty(CriticialChance, 0, property.FriendlyController, property.Skill)
}

func DefaultCriticialDamage() *property.DefaultFloatProperty {
	return property.MakeFloatProperty(CriticialDamage, 0, property.FriendlyController, property.Skill)
}

func DefaultDuration() *property.DefaultIntCounterProperty {
	return property.MakeIntCounterProperty(Duration, 0, 0, property.FriendlyController, property.Skill)
}

func DefaultPoisonPower() *property.DefaultIntProperty {
	return property.MakeIntProperty(PoisonPower, 0, property.FriendlyController, property.Skill)
}

func DefaultPoisonChance() *property.DefaultIntProperty {
	return property.MakeIntProperty(PoisonChance, 0, property.FriendlyController, property.Skill)
}

func DefaultDelay() *property.DefaultIntCounterProperty {
	return property.MakeIntCounterProperty(Delay, 0, 500, property.FriendlyController, property.Skill)
}

func DefaultHPLeech() *property.DefaultIntProperty {
	return property.MakeIntProperty(HPLeech, 0, property.FriendlyController, property.Skill)
}

func DefaultMPLeech() *property.DefaultIntProperty {
	return property.MakeIntProperty(MPLeech, 0, property.FriendlyController, property.Skill)
}

func DefaultSPLeech() *property.DefaultIntProperty {
	return property.MakeIntProperty(SPLeech, 0, property.FriendlyController, property.Skill)
}

func DefaultCooldown() *property.DefaultIntCounterProperty {
	return property.MakeIntCounterProperty(Cooldown, 0, 3, property.FriendlyController, property.Skill)
}

// Behavior property.Property: 	Behavior SkillProperties = "Behavior" property.Skill broad category: Direct, Reaction, Passive, Counter

type BehaviorType int

const (
	BehaviorTypeDirect BehaviorType = iota
	BehaviorTypeReaction
	BehaviorTypePassive
	BehaviorTypeCounter
)

type BehaviorProperty struct {
	property.Property
	BehaviorType BehaviorType
}

// MakeBehaviorProperty creates a BehaviorProperty
func MakeBehaviorProperty(bh BehaviorType) *BehaviorProperty {
	return &BehaviorProperty{
		BehaviorType: bh,
	}
}

// DefaultBehaviorProperty
func DefaultBehaviorProperty() *BehaviorProperty {
	return MakeBehaviorProperty(BehaviorTypeDirect)
}

func (bh *BehaviorProperty) Name(i property.InformationLevel) string {
	return "Behavior"
}

func (bh *BehaviorProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return bh.BehaviorType
}

func (bh *BehaviorProperty) Get() interface{} {
	return bh.BehaviorType
}

func (bh *BehaviorProperty) Set(p interface{}) {
	// shouldn't be altered.
	bh.BehaviorType = p.(BehaviorType)
}

func (bh *BehaviorProperty) Increase() {
	// shouldn't be upgradable.... well maybe convert a Direct skill to a Reaction or Counter ? fun...
}

func (bh *BehaviorProperty) GetType() property.PropertyType {
	return property.Skill
}

func (db *BehaviorProperty) Duplicate() property.Property {
	return &BehaviorProperty{
		BehaviorType: db.BehaviorType,
	}
}

func (bh *BehaviorProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*BehaviorProperty)
	// replace
	res.BehaviorType = p.(*BehaviorProperty).BehaviorType
	return res
}

// Range property.Property: 	Range TargetingProperties = "Range" // Range of the property.Skill

type RangeProperty struct {
	property.Property
	MinRange int
	MaxRange int
}

// MakeRangeProperty
func MakeRangeProperty(min, max int) *RangeProperty {
	return &RangeProperty{
		MinRange: min,
		MaxRange: max,
	}
}

// DefaultRangeProperty
func DefaultRangeProperty() *RangeProperty {
	return MakeRangeProperty(1, 1)
}

func (bh *RangeProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*RangeProperty)
	// replace
	res.MinRange = p.(*RangeProperty).MinRange
	res.MaxRange = p.(*RangeProperty).MaxRange
	return res
}

func (bh *RangeProperty) Name(i property.InformationLevel) string {
	return "Range"
}

func (bh *RangeProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return fmt.Sprintf("%d - %d", bh.MinRange, bh.MaxRange)
}

func (bh *RangeProperty) Get() interface{} {
	return bh
}

func (bh *RangeProperty) Set(p interface{}) {
	// will be altered directly.
}

func (bh *RangeProperty) Increase() {
	// shouldn't be upgradable.... well maybe convert a Direct skill to a Reaction or Counter ? fun...
}

func (bh *RangeProperty) GetType() property.PropertyType {
	return property.Skill
}

func (bh *RangeProperty) Duplicate() property.Property {
	return &RangeProperty{
		MinRange: bh.MinRange,
		MaxRange: bh.MaxRange,
	}
}

// Zone         TargetingProperties = "Zone"  // Area of Effect
// ZoneProperty expects to be casted to be used.
type ZoneProperty struct {
	property.Property
	ZonePattern pattern.Pattern
}

// MakeZoneProperty
func MakeZoneProperty(zp pattern.Pattern) *ZoneProperty {
	return &ZoneProperty{
		ZonePattern: zp,
	}
}

// DefaultZoneProperty
func DefaultZoneProperty() *ZoneProperty {
	return MakeZoneProperty(pattern.Single())
}

func (bh *ZoneProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*ZoneProperty)
	// replace
	res.ZonePattern = p.(*ZoneProperty).ZonePattern
	return res
}

func (bh *ZoneProperty) Name(i property.InformationLevel) string {
	return "Zone"
}

func (bh *ZoneProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return "" // will be implemented later
}

func (bh *ZoneProperty) Get() interface{} {
	return bh
}

func (bh *ZoneProperty) Set(p interface{}) {
	// will be altered directly.
}

func (bh *ZoneProperty) Increase() {
	// shouldn't be upgradable.... well maybe convert a Direct skill to a Reaction or Counter ? fun...
}

func (bh *ZoneProperty) GetType() property.PropertyType {
	return property.Skill
}

func (bh *ZoneProperty) Duplicate() property.Property {
	return &ZoneProperty{
		ZonePattern: bh.ZonePattern,
	}
}

// TargetType         TargetingProperties = "TargetType"         // Entity, Tile, Both, Self

type TargetTypes string

const (
	TargetTypeEntity TargetTypes = "Entity"
	TargetTypeTile   TargetTypes = "Tile"
	TargetTypeBoth   TargetTypes = "Both"
	TargetTypeSelf   TargetTypes = "Self"
)

type TargetTypeProperty struct {
	property.Property
	TargetType TargetTypes
}

// MakeTargetTypeProperty
func MakeTargetTypeProperty(tt TargetTypes) *TargetTypeProperty {
	return &TargetTypeProperty{
		TargetType: tt,
	}
}

// DefaultTargetTypeProperty
func DefaultTargetTypeProperty() *TargetTypeProperty {
	return MakeTargetTypeProperty(TargetTypeEntity)
}

func (bh *TargetTypeProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*TargetTypeProperty)
	// replace
	res.TargetType = p.(*TargetTypeProperty).TargetType
	return res
}

func (bh *TargetTypeProperty) Name(i property.InformationLevel) string {
	return "TargetType"
}

func (bh *TargetTypeProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return bh.TargetType
}

func (bh *TargetTypeProperty) Get() interface{} {
	return bh
}

func (bh *TargetTypeProperty) Set(p interface{}) {
	// will be altered directly.
}

func (bh *TargetTypeProperty) Increase() {
}

func (bh *TargetTypeProperty) GetType() property.PropertyType {
	return property.Skill
}

func (bh *TargetTypeProperty) Duplicate() property.Property {
	return &TargetTypeProperty{
		TargetType: bh.TargetType,
	}
}

// TargetingMechanics TargetingProperties = "TargetingMechanics" // Anywhere, Line of Sight, and maybe other mechanics later.

type TargetingMechanicsType string

const (
	TargetingMechanicsAnywhere TargetingMechanicsType = "Anywhere"
	TargetingMechanicsLOS      TargetingMechanicsType = "Line of Sight"
)

type TargetingMechanicsProperty struct {
	property.Property
	TargetingMechanics TargetingMechanicsType
}

// MakeTargetingMechanicsProperty
func MakeTargetingMechanicsProperty(tm TargetingMechanicsType) *TargetingMechanicsProperty {
	return &TargetingMechanicsProperty{
		TargetingMechanics: tm,
	}
}

// DefaultTargetingMechanicsProperty
func DefaultTargetingMechanicsProperty() *TargetingMechanicsProperty {
	return MakeTargetingMechanicsProperty(TargetingMechanicsAnywhere)
}

func (bh *TargetingMechanicsProperty) Name(i property.InformationLevel) string {
	return "TargetingMechanics"
}

func (bh *TargetingMechanicsProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return bh.TargetingMechanics
}

func (bh *TargetingMechanicsProperty) Get() interface{} {
	return bh
}

func (bh *TargetingMechanicsProperty) Set(p interface{}) {
	// will be altered directly.
}

func (bh *TargetingMechanicsProperty) Increase() {
}

func (bh *TargetingMechanicsProperty) GetType() property.PropertyType {
	return property.Skill
}

func (bh *TargetingMechanicsProperty) Duplicate() property.Property {
	return &TargetingMechanicsProperty{
		TargetingMechanics: bh.TargetingMechanics,
	}
}

func (bh *TargetingMechanicsProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*TargetingMechanicsProperty)
	// replace
	res.TargetingMechanics = p.(*TargetingMechanicsProperty).TargetingMechanics
	return res
}
