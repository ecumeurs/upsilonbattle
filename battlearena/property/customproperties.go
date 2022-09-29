package property

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position/pattern"
)

// use this file to define custom properties that aren't basic Int/Float properties.

// Behavior Property: 	Behavior SkillProperties = "Behavior" Skill broad category: Direct, Reaction, Passive, Counter

type BehaviorType int

const (
	BehaviorTypeDirect BehaviorType = iota
	BehaviorTypeReaction
	BehaviorTypePassive
	BehaviorTypeCounter
)

type BehaviorProperty struct {
	Property
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

func (bh *BehaviorProperty) Name(i InformationLevel) string {
	return "Behavior"
}

func (bh *BehaviorProperty) UserFriendlyGet(i InformationLevel) interface{} {
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

func (bh *BehaviorProperty) GetType() PropertyType {
	return Skill
}

// Range Property: 	Range TargetingProperties = "Range" // Range of the Skill

type RangeProperty struct {
	Property
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

func (bh *RangeProperty) Name(i InformationLevel) string {
	return "Range"
}

func (bh *RangeProperty) UserFriendlyGet(i InformationLevel) interface{} {
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

func (bh *RangeProperty) GetType() PropertyType {
	return Skill
}

// Zone         TargetingProperties = "Zone"  // Area of Effect
// ZoneProperty expects to be casted to be used.
type ZoneProperty struct {
	Property
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

func (bh *ZoneProperty) Name(i InformationLevel) string {
	return "Zone"
}

func (bh *ZoneProperty) UserFriendlyGet(i InformationLevel) interface{} {
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

func (bh *ZoneProperty) GetType() PropertyType {
	return Skill
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
	Property
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

func (bh *TargetTypeProperty) Name(i InformationLevel) string {
	return "TargetType"
}

func (bh *TargetTypeProperty) UserFriendlyGet(i InformationLevel) interface{} {
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

func (bh *TargetTypeProperty) GetType() PropertyType {
	return Skill
}

// TargetingMechanics TargetingProperties = "TargetingMechanics" // Anywhere, Line of Sight, and maybe other mechanics later.

type TargetingMechanicsType string

const (
	TargetingMechanicsAnywhere TargetingMechanicsType = "Anywhere"
	TargetingMechanicsLOS      TargetingMechanicsType = "Line of Sight"
)

type TargetingMechanicsProperty struct {
	Property
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

func (bh *TargetingMechanicsProperty) Name(i InformationLevel) string {
	return "TargetingMechanics"
}

func (bh *TargetingMechanicsProperty) UserFriendlyGet(i InformationLevel) interface{} {
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

func (bh *TargetingMechanicsProperty) GetType() PropertyType {
	return Skill
}
