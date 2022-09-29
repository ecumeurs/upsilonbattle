package def

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position/pattern"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
)

// Prepare default Properties.

func DefaultTargetNumber() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.TargetNumber, 0, property.FriendlyController, property.Skill)
}

func DefaultAccuracy() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Accuracy, 100, property.FriendlyController, property.Skill)
}

func DefaultDodge() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Dodge, 0, property.FriendlyController, property.Skill)
}

func DefaultParry() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Parry, 0, property.FriendlyController, property.Skill)
}

func DefaultDamage() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Damage, 0, property.FriendlyController, property.Skill)
}

func DefaultHeal() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Heal, 0, property.FriendlyController, property.Skill)
}

func DefaultShieldPower() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.ShieldPower, 0, property.FriendlyController, property.Skill)
}

func DefaultStunPower() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.StunPower, 0, property.FriendlyController, property.Skill)
}

func DefaultStunChance() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.StunChance, 0, property.FriendlyController, property.Skill)
}

func DefaultCriticialChance() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.CriticialChance, 0, property.FriendlyController, property.Skill)
}

func DefaultCriticialDamage() *defaultproperty.DefaultFloatProperty {
	return defaultproperty.MakeFloatProperty(property.CriticialDamage, 0, property.FriendlyController, property.Skill)
}

func DefaultDuration() *defaultproperty.DefaultIntCounterProperty {
	return defaultproperty.MakeIntCounterProperty(property.Duration, 0, 0, property.FriendlyController, property.Skill)
}

func DefaultPoisonPower() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.PoisonPower, 0, property.FriendlyController, property.Skill)
}

func DefaultPoisonChance() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.PoisonChance, 0, property.FriendlyController, property.Skill)
}

func DefaultDelay() *defaultproperty.DefaultIntCounterProperty {
	return defaultproperty.MakeIntCounterProperty(property.Delay, 0, 500, property.FriendlyController, property.Skill)
}

func DefaultHPLeech() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.HPLeech, 0, property.FriendlyController, property.Skill)
}

func DefaultMPLeech() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.MPLeech, 0, property.FriendlyController, property.Skill)
}

func DefaultSPLeech() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.SPLeech, 0, property.FriendlyController, property.Skill)
}

func DefaultCooldown() *defaultproperty.DefaultIntCounterProperty {
	return defaultproperty.MakeIntCounterProperty(property.Cooldown, 0, 3, property.FriendlyController, property.Skill)
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
