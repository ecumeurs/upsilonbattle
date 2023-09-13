package def

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position/pattern"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
)

// Prepare default Properties.

func TargetNumber() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.TargetNumber, 0, property.FriendlyController, property.Skill)
}

func Accuracy() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Accuracy, 100, property.FriendlyController, property.Skill)
}

func Dodge() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Dodge, 0, property.FriendlyController, property.Skill)
}

func Parry() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Parry, 0, property.FriendlyController, property.Skill)
}

func Damage() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Damage, 100, property.FriendlyController, property.Skill)
}

func Heal() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Heal, 0, property.FriendlyController, property.Skill)
}

func ShieldPower() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.ShieldPower, 0, property.FriendlyController, property.Skill)
}

func StunPower() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.StunPower, 0, property.FriendlyController, property.Skill)
}

func StunChance() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.StunChance, 0, property.FriendlyController, property.Skill)
}

func CriticalChance() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.CriticalChance, 0, property.FriendlyController, property.Skill)
}

func CriticalMultiplier() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.CriticalMultiplier, 100, property.FriendlyController, property.Skill)
}

func Duration() *defaultproperty.DefaultIntCounterProperty {
	return defaultproperty.MakeIntCounterProperty(property.Duration, 0, 0, property.FriendlyController, property.Skill)
}

func PoisonPower() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.PoisonPower, 0, property.FriendlyController, property.Skill)
}

func PoisonChance() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.PoisonChance, 0, property.FriendlyController, property.Skill)
}

func Delay() *defaultproperty.DefaultIntCounterProperty {
	return defaultproperty.MakeIntCounterProperty(property.Delay, 0, 500, property.FriendlyController, property.Skill)
}

func HPLeech() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.HPLeech, 0, property.FriendlyController, property.Skill)
}

func MPLeech() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.MPLeech, 0, property.FriendlyController, property.Skill)
}

func SPLeech() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.SPLeech, 0, property.FriendlyController, property.Skill)
}

func MvtCost() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.MvtCost, 0, property.FriendlyController, property.Skill)
}

// Cooldown() Default to 0, 3 ;Special note: Cool down is stored as a counter, minValue represent initial cooldown at battle start. MaxValue represent the cooldown value when used.
func Cooldown() *defaultproperty.DefaultIntCounterProperty {

	return defaultproperty.MakeIntCounterProperty(property.Cooldown, 0, 3, property.FriendlyController, property.Skill)
}

// Behavior property.Property: 	Behavior SkillProperties = "Behavior" property.Skill broad category: Direct, Reaction, Passive, Counter

// Note: Behavior is a property because ...
// Expect buffs to have an impact on behavior (e.g. a buff that makes a skill a counter skill)

type BehaviorType int

const (
	BehaviorTypeDirect BehaviorType = iota
	BehaviorTypeReaction
	BehaviorTypePassive
	BehaviorTypeCounter
	BehaviorTypeTrap
)

type BehaviorProperty struct {
	property.Property
	BehaviorType BehaviorType
}

// IsDirect
func (bh *BehaviorProperty) IsDirect() bool {
	return bh.BehaviorType == BehaviorTypeDirect
}

// IsReaction
func (bh *BehaviorProperty) IsReaction() bool {
	return bh.BehaviorType == BehaviorTypeReaction
}

// IsPassive
func (bh *BehaviorProperty) IsPassive() bool {
	return bh.BehaviorType == BehaviorTypePassive
}

// IsCounter
func (bh *BehaviorProperty) IsCounter() bool {
	return bh.BehaviorType == BehaviorTypeCounter
}

// IsTrap
func (bh *BehaviorProperty) IsTrap() bool {
	return bh.BehaviorType == BehaviorTypeTrap
}

// MakeBehaviorProperty creates a BehaviorProperty
func MakeBehaviorProperty(bh BehaviorType) *BehaviorProperty {
	return &BehaviorProperty{
		BehaviorType: bh,
	}
}

// DefaultBehaviorProperty
func DefaultBehavior() *BehaviorProperty {
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

func (bh *BehaviorProperty) UnapplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*BehaviorProperty)
	// TODO ... :)
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
func DefaultRange() *RangeProperty {
	return MakeRangeProperty(1, 1)
}

func (bh *RangeProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*RangeProperty)
	// replace
	res.MinRange = p.(*RangeProperty).MinRange
	res.MaxRange = p.(*RangeProperty).MaxRange
	return res
}
func (bh *RangeProperty) UnapplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*RangeProperty)
	// TODO :)
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
func DefaultZone() *ZoneProperty {
	return MakeZoneProperty(pattern.Single())
}

func (bh *ZoneProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*ZoneProperty)
	// replace
	res.ZonePattern = p.(*ZoneProperty).ZonePattern
	return res
}

func (bh *ZoneProperty) UnapplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*ZoneProperty)
	// TODO
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
	TargetTypeEntity       TargetTypes = "Entity"
	TargetTypeFriendOnly   TargetTypes = "FriendOnly"
	TargetTypeEnemyOnly    TargetTypes = "EnemyOnly"
	TargetTypeTile         TargetTypes = "Tile"
	TargetTypeEntityOrTile TargetTypes = "EntityOrTile"
	TargetTypeSelf         TargetTypes = "Self"
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
func DefaultTargetType() *TargetTypeProperty {
	return MakeTargetTypeProperty(TargetTypeEntity)
}

func (bh *TargetTypeProperty) ApplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*TargetTypeProperty)
	// replace
	res.TargetType = p.(*TargetTypeProperty).TargetType
	return res
}

func (bh *TargetTypeProperty) UnapplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*TargetTypeProperty)
	// TODO
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
	return bh.TargetType
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
func DefaultTargetingMechanics() *TargetingMechanicsProperty {
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

func (bh *TargetingMechanicsProperty) UnapplyBuff(p property.Property) property.Property {
	res := bh.Duplicate().(*TargetingMechanicsProperty)
	// TODO
	res.TargetingMechanics = p.(*TargetingMechanicsProperty).TargetingMechanics
	return res
}

func SkillProperty(ps property.SkillProperties) property.Property {
	switch ps {
	case property.Accuracy:
		return Accuracy()
	case property.Behavior:
		return DefaultBehavior()
	case property.Range:
		return DefaultRange()
	case property.Zone:
		return DefaultZone()
	case property.TargetType:
		return DefaultTargetType()
	case property.TargetingMechanics:
		return DefaultTargetingMechanics()
	case property.Dodge:
		return Dodge()
	case property.Parry:
		return Parry()
	case property.Damage:
		return Damage()
	case property.Heal:
		return Heal()
	case property.ShieldPower:
		return ShieldPower()
	case property.StunPower:
		return StunPower()
	case property.StunChance:
		return StunChance()
	case property.CriticalChance:
		return CriticalChance()
	case property.CriticalMultiplier:
		return CriticalMultiplier()
	case property.Duration:
		return Duration()
	case property.PoisonPower:
		return PoisonPower()
	case property.PoisonChance:
		return PoisonChance()
	case property.Delay:
		return Delay()
	case property.Cooldown:
		return Cooldown()
	case property.HPLeech:
		return HPLeech()
	case property.MPLeech:
		return MPLeech()
	case property.SPLeech:
		return SPLeech()
	case property.MvtCost:
		return MvtCost()
	}
	return nil
}
