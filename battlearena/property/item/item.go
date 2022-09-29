package item

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
)

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
)

func DefaultDurability() *property.DefaultIntProperty {
	return property.MakeIntProperty(Durability, 0, property.Public, property.Item)
}

func DefaultWeight() *property.DefaultIntProperty {
	return property.MakeIntProperty(Weight, 0, property.Public, property.Item)
}

func DefaultArmorRating() *property.DefaultIntProperty {
	return property.MakeIntProperty(ArmorRating, 0, property.Public, property.Item)
}

func DefaultWeaponRange() *property.DefaultIntProperty {
	return property.MakeIntProperty(WeaponRange, 0, property.Public, property.Item)
}

func DefaultWeaponBaseDamage() *property.DefaultIntProperty {
	return property.MakeIntProperty(WeaponBaseDamage, 0, property.Public, property.Item)
}

func DefaultStackable() *property.DefaultBoolProperty {
	return property.MakeBoolProperty(Stackable, false, property.Public, property.Item)
}

func DefaultStackSize() *property.DefaultIntProperty {
	return property.MakeIntProperty(StackSize, 0, property.Public, property.Item)
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
func MakeItemTypeProperty(it ItemTypes, minInfoLevel property.InformationLevel) *ItemTypeProperty {
	return &ItemTypeProperty{
		ItemType:            it,
		minInformationLevel: minInfoLevel,
	}
}

// DefaultItemTypeProperty
func DefaultItemTypeProperty() *ItemTypeProperty {
	return MakeItemTypeProperty(ItemTypeNone, property.OwnController)
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
	Effect              *skill.Skill
	minInformationLevel property.InformationLevel
}

// MakeEffectProperty
func MakeEffectProperty(e *skill.Skill, minInfoLevel property.InformationLevel) *EffectProperty {
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
