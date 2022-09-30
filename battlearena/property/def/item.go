package def

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect"
)

func Durability() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Durability, 0, property.Public, property.Item)
}

func Weight() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Weight, 0, property.Public, property.Item)
}

func ArmorRating() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.ArmorRating, 0, property.Public, property.Item)
}

func WeaponRange() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.WeaponRange, 0, property.Public, property.Item)
}

func WeaponBaseDamage() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.WeaponBaseDamage, 0, property.Public, property.Item)
}

func Stackable() *defaultproperty.DefaultBoolProperty {
	return defaultproperty.MakeBoolProperty(property.Stackable, false, property.Public, property.Item)
}

func StackSize() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.StackSize, 0, property.Public, property.Item)
}

func Value() *defaultproperty.DefaultIntProperty {
	return defaultproperty.MakeIntProperty(property.Value, 0, property.Public, property.Item)
}

// ItemType         ItemProperties = "ItemType"         // Absence means None (out of Wearable, Consumable, Usable, Throwable, Ammunitions and None)

type ItemTypes string

const (
	ItemTypeWearable    ItemTypes = "Wearable"
	ItemTypeConsumable  ItemTypes = "Consumable"
	ItemTypeUsable      ItemTypes = "Usable"
	ItemTypeThrowable   ItemTypes = "Throwable"
	ItemTypeAmmunitions ItemTypes = "Ammunitions"
	ItemTypeNone        ItemTypes = "Misc"
)

type ItemTypeProperty struct {
	property.Property
	ItemType            ItemTypes
	minInformationLevel property.InformationLevel
}

// MakeItemTypeProperty
func MakeItemType(it ItemTypes, minInfoLevel property.InformationLevel) *ItemTypeProperty {
	return &ItemTypeProperty{
		ItemType:            it,
		minInformationLevel: minInfoLevel,
	}
}

// DefaultItemTypeProperty
func DefaultItemType() *ItemTypeProperty {
	return MakeItemType(ItemTypeNone, property.OwnController)
}

func (bh *ItemTypeProperty) Name(i property.InformationLevel) string {
	if i >= bh.minInformationLevel {
		return "ItemType"
	}
	return ""
}

func (bh *ItemTypeProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	if i >= bh.minInformationLevel {
		return bh.ItemType
	}
	return ItemTypeNone
}

func (bh *ItemTypeProperty) Get() interface{} {
	return bh
}

func (bh *ItemTypeProperty) Set(p interface{}) {
	// will be altered directly.
}

func (bh *ItemTypeProperty) Increase() {
}

func (bh *ItemTypeProperty) GetType() property.PropertyType {
	return property.Item
}

func (bh *ItemTypeProperty) Duplicate() property.Property {
	return &ItemTypeProperty{
		ItemType: bh.ItemType,
	}
}

func (bh *ItemTypeProperty) ApplyBuff(p property.Property) property.Property {
	return bh.Duplicate() // property.Item can't be buffed.
}

// 	Effect           ItemProperties = "Effect"           // Absence means nil: No effect. Effects are Skills. (except None)

type EffectProperty struct {
	property.Property
	Effect              *effect.Effect
	minInformationLevel property.InformationLevel
}

// MakeEffectProperty
func MakeEffectProperty(e *effect.Effect, minInfoLevel property.InformationLevel) *EffectProperty {
	return &EffectProperty{
		Effect:              e,
		minInformationLevel: minInfoLevel,
	}
}

// DefaultEffectProperty
func DefaultEffectProperty() *EffectProperty {
	return MakeEffectProperty(nil, property.Analyser)
}

func (bh *EffectProperty) Name(i property.InformationLevel) string {
	if i >= bh.minInformationLevel {
		return "Effect"
	}
	return ""
}

func (bh *EffectProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	if i >= bh.minInformationLevel {
		return bh.Effect
	}
	return nil
}

func (bh *EffectProperty) Get() interface{} {
	return bh.Effect
}

func (bh *EffectProperty) Set(p interface{}) {
}

func (bh *EffectProperty) Increase() {
}

func (bh *EffectProperty) GetType() property.PropertyType {
	return property.Item
}

func (bh *EffectProperty) Duplicate() property.Property {
	return &EffectProperty{
		Effect:              bh.Effect,
		minInformationLevel: bh.minInformationLevel,
	}
}

func (bh *EffectProperty) ApplyBuff(p property.Property) property.Property {
	return bh.Duplicate() // property.Item can't be buffed.
}

// 	WeaponType       ItemProperties = "WeaponType"       // Absence means 0: no weapon type (only for Wearable)

type WeaponTypes string

const (
	NoWeapon WeaponTypes = "None"
	// Melee
	OneHandedMelee WeaponTypes = "One-Handed Melee"
	TwoHandedMelee WeaponTypes = "Two-Handed Melee"
	// Ranged
	OneHandedRanged WeaponTypes = "One-Handed Ranged"
	TwoHandedRanged WeaponTypes = "Two-Handed Ranged"
)

type WeaponTypeProperty struct {
	property.Property
	WeaponType WeaponTypes
}

// MakeWeaponTypeProperty
func MakeWeaponTypeProperty(wt WeaponTypes) *WeaponTypeProperty {
	return &WeaponTypeProperty{
		WeaponType: wt,
	}
}

// DefaultWeaponTypeProperty
func DefaultWeaponTypeProperty() *WeaponTypeProperty {
	return MakeWeaponTypeProperty(NoWeapon)
}

func (bh *WeaponTypeProperty) Name(i property.InformationLevel) string {
	return "WeaponType"
}

func (bh *WeaponTypeProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return bh.WeaponType
}

func (bh *WeaponTypeProperty) Get() interface{} {
	return bh.WeaponType
}

func (bh *WeaponTypeProperty) Set(p interface{}) {
}

func (bh *WeaponTypeProperty) Increase() {
}

func (bh *WeaponTypeProperty) GetType() property.PropertyType {
	return property.Item
}

func (bh *WeaponTypeProperty) Duplicate() property.Property {
	return &WeaponTypeProperty{
		WeaponType: bh.WeaponType,
	}
}

func (bh *WeaponTypeProperty) ApplyBuff(p property.Property) property.Property {
	return bh.Duplicate() // property.Item can't be buffed.
}

//ArmorType        ItemProperties = "ArmorType"        // Absence means 0: no armor type (only for Wearable)

type ArmorTypes string

const (
	NoArmor   ArmorTypes = "None"
	HeadSlot  ArmorTypes = "Head"
	BodySlot  ArmorTypes = "Body"
	HandsSlot ArmorTypes = "Hands"
	LegsSlot  ArmorTypes = "Legs"
	FeetSlot  ArmorTypes = "Feet"
	BeltSlot  ArmorTypes = "Belt"
	NeckSlot  ArmorTypes = "Neck"
	RingSlot  ArmorTypes = "Ring"
)

type ArmorTypeProperty struct {
	property.Property
	ArmorType ArmorTypes
}

// MakeArmorTypeProperty
func MakeArmorTypeProperty(at ArmorTypes) *ArmorTypeProperty {
	return &ArmorTypeProperty{
		ArmorType: at,
	}
}

// DefaultArmorTypeProperty
func DefaultArmorTypeProperty() *ArmorTypeProperty {
	return MakeArmorTypeProperty(NoArmor)
}

func (bh *ArmorTypeProperty) Name(i property.InformationLevel) string {
	return "ArmorType"
}

func (bh *ArmorTypeProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return bh.ArmorType
}

func (bh *ArmorTypeProperty) Get() interface{} {
	return bh.ArmorType
}

func (bh *ArmorTypeProperty) Set(p interface{}) {
}

func (bh *ArmorTypeProperty) Increase() {
}

func (bh *ArmorTypeProperty) GetType() property.PropertyType {
	return property.Item
}

func (bh *ArmorTypeProperty) Duplicate() property.Property {
	return &ArmorTypeProperty{
		ArmorType: bh.ArmorType,
	}
}

func (bh *ArmorTypeProperty) ApplyBuff(p property.Property) property.Property {
	return bh.Duplicate() // property.Item can't be buffed.
}

//ToolType         ItemProperties = "ToolType"         // Absence means 0: no tool type (only for Wearable)

type ToolTypes string

const (
	NoTool   ToolTypes = "None"
	SomeTool ToolTypes = "SomeTool"
)

type ToolTypeProperty struct {
	property.Property
	ToolType ToolTypes
}

// MakeToolTypeProperty
func MakeToolTypeProperty(tt ToolTypes) *ToolTypeProperty {
	return &ToolTypeProperty{
		ToolType: tt,
	}
}

// DefaultToolTypeProperty
func DefaultToolTypeProperty() *ToolTypeProperty {
	return MakeToolTypeProperty(NoTool)
}

func (bh *ToolTypeProperty) Name(i property.InformationLevel) string {
	return "ToolType"
}

func (bh *ToolTypeProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return bh.ToolType
}

func (bh *ToolTypeProperty) Get() interface{} {
	return bh.ToolType
}

func (bh *ToolTypeProperty) Set(p interface{}) {
}

func (bh *ToolTypeProperty) Increase() {
}

func (bh *ToolTypeProperty) GetType() property.PropertyType {
	return property.Item
}

func (bh *ToolTypeProperty) Duplicate() property.Property {
	return &ToolTypeProperty{
		ToolType: bh.ToolType,
	}
}

func (bh *ToolTypeProperty) ApplyBuff(p property.Property) property.Property {
	return bh.Duplicate() // property.Item can't be buffed.
}

func ItemProperty(name property.ItemProperties) property.Property {
	switch name {
	case property.Effect:
		return DefaultEffectProperty()
	case property.WeaponType:
		return DefaultWeaponTypeProperty()
	case property.ArmorType:
		return DefaultArmorTypeProperty()
	case property.ToolType:
		return DefaultToolTypeProperty()
	case property.Durability:
		return Durability()
	case property.Weight:
		return Weight()
	case property.Value:
		return Value()
	case property.ItemType:
		return DefaultItemType()
	case property.WeaponBaseDamage:
		return WeaponBaseDamage()
	case property.WeaponRange:
		return WeaponRange()
	case property.StackSize:
		return StackSize()
	case property.Stackable:
		return Stackable()

	}

	return nil
}
