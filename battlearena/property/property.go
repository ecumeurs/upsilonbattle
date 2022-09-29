package property

import (
	"fmt"
)

type PropertyType int

const (
	None      PropertyType = 0
	Character PropertyType = 1
	Skill     PropertyType = 2
	Item      PropertyType = 3
)

type InformationLevel int

const (
	Public             InformationLevel = 0
	ArenaObserver      InformationLevel = 1
	ForeignController  InformationLevel = 2
	FriendlyController InformationLevel = 3

	OwnController InformationLevel = 4

	Analyser          InformationLevel = 5
	ExpertAnalyst     InformationLevel = 6
	SpecialistAnalyst InformationLevel = 7
	MasterAnalyst     InformationLevel = 8

	GameMaster InformationLevel = 9
)

type Property interface {
	Name(i InformationLevel) string                 // most will always reply with a name, some might be hidden by restrictions of with a scrambled name.
	UserFriendlyGet(i InformationLevel) interface{} // most will be expected to return an int (float will be frowned upon) but might return a string if appropriate (status for example) may return nil in which case the information won't be displayed.
	Get() interface{}                               // this will be used mostly internally to compute values from rules.
	Set(p interface{})                              // this will be used mostly internally to compute values from rules.
	Increase()                                      // This won't be used in v0.0.2 but later on when we implement leveling.
	GetType() PropertyType
	Duplicate() Property
	ApplyBuff(p Property) Property
}

// PrettyPrint
func PrettyPrint(p Property, i InformationLevel) string {
	return fmt.Sprintf("%s: %v", p.Name(i), p.UserFriendlyGet(i))
}

// these interface are here for rules ...
type IntProperty interface {
	Property
	I() int
	SetI(int)
}

type FloatProperty interface {
	Property
	F() float64
	SetF(float64)
}

type BoolProperty interface {
	Property
	B() bool
	SetB(bool)
}

type DefaultIntProperty struct {
	value               int
	name                string
	minInformationLevel InformationLevel
	propertytype        PropertyType
}

// MakeIntProperty
func MakeIntProperty(name interface{}, value int, minInformationLevel InformationLevel, t PropertyType) *DefaultIntProperty {
	switch name.(type) {
	case string:
		return &DefaultIntProperty{
			name:                name.(string),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertytype:        t,
		}
	default:
		return nil
	}
}

// implements IntProperty
func (d DefaultIntProperty) Name(i InformationLevel) string {
	if i >= d.minInformationLevel {
		return d.name
	}
	return ""
}

func (d DefaultIntProperty) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= d.minInformationLevel {
		return d.value
	}
	return nil
}

func (d DefaultIntProperty) Get() interface{} {
	return d.value
}

func (d *DefaultIntProperty) Set(p interface{}) {
	d.value = p.(int)
}

func (d DefaultIntProperty) Increase() {
	// do nothing
}

func (d DefaultIntProperty) GetType() PropertyType {
	return d.propertytype
}

func (d DefaultIntProperty) I() int {
	return d.value
}

func (d *DefaultIntProperty) SetI(i int) {
	d.value = i
}

func (d DefaultIntProperty) ApplyBuff(p Property) Property {
	res := d.Duplicate()
	res.Set(d.Get().(int) + p.Get().(int))
	return res
}

func (d DefaultIntProperty) Duplicate() Property {
	return &DefaultIntProperty{
		value:               d.value,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertytype:        d.propertytype,
	}
}

type DefaultIntCounterProperty struct {
	Value    int
	MaxValue int

	name                string
	minInformationLevel InformationLevel
	propertytype        PropertyType
}

// MakeIntProperty
func MakeIntCounterProperty(name interface{}, value, maxvalue int, minInformationLevel InformationLevel, t PropertyType) *DefaultIntCounterProperty {
	switch name.(type) {
	case string:
		return &DefaultIntCounterProperty{
			name:                name.(string),
			Value:               value,
			MaxValue:            maxvalue,
			minInformationLevel: minInformationLevel,
			propertytype:        t,
		}
	default:
		return nil
	}
}

// implements IntProperty
func (d DefaultIntCounterProperty) Name(i InformationLevel) string {
	if i >= d.minInformationLevel {
		return d.name
	}
	return ""
}

func (d DefaultIntCounterProperty) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= d.minInformationLevel {
		return d.Value
	}
	return nil
}

func (d DefaultIntCounterProperty) Get() interface{} {
	return d.Value
}

func (d *DefaultIntCounterProperty) Set(p interface{}) {
	d.Value = p.(int)
}

func (d DefaultIntCounterProperty) Increase() {
	// do nothing
}

func (d DefaultIntCounterProperty) GetType() PropertyType {
	return d.propertytype
}

func (d DefaultIntCounterProperty) I() int {
	return d.Value
}

func (d *DefaultIntCounterProperty) SetI(i int) {
	d.Value = i
}

func (d DefaultIntCounterProperty) ApplyBuff(p Property) Property {
	res := d.Duplicate().(*DefaultIntCounterProperty)
	res.Value = d.Value + p.(*DefaultIntCounterProperty).Value
	res.MaxValue = d.MaxValue + p.(*DefaultIntCounterProperty).MaxValue
	return res
}

func (d DefaultIntCounterProperty) Duplicate() Property {
	return &DefaultIntCounterProperty{
		Value:               d.Value,
		MaxValue:            d.MaxValue,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertytype:        d.propertytype,
	}
}

type DefaultFloatProperty struct {
	value               float64
	name                string
	minInformationLevel InformationLevel
	propertytype        PropertyType
}

// MakeIntProperty
func MakeFloatProperty(name interface{}, value float64, minInformationLevel InformationLevel, pt PropertyType) *DefaultFloatProperty {
	switch name.(type) {
	case string:

		return &DefaultFloatProperty{
			name:                name.(string),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertytype:        pt,
		}
	default:
		return nil
	}
}

// implements IntProperty
func (d DefaultFloatProperty) Name(i InformationLevel) string {
	return d.name
}

func (d DefaultFloatProperty) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= d.minInformationLevel {
		return d.value
	}
	return nil
}

func (d DefaultFloatProperty) Get() interface{} {
	return d.value
}

func (d *DefaultFloatProperty) Set(p interface{}) {
	d.value = p.(float64)
}

func (d DefaultFloatProperty) Increase() {
	// do nothing
}

func (d DefaultFloatProperty) GetType() PropertyType {
	return d.propertytype
}

func (d DefaultFloatProperty) I() float64 {
	return d.value
}

func (d *DefaultFloatProperty) SetI(f float64) {
	d.value = f
}

func (d DefaultFloatProperty) Duplicate() Property {
	return &DefaultFloatProperty{
		value:               d.value,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertytype:        d.propertytype,
	}
}

func (d DefaultFloatProperty) ApplyBuff(p Property) Property {
	res := d.Duplicate()
	res.Set(d.Get().(float64) + p.Get().(float64))
	return res
}

// Bool default value are essentially flags ...

type DefaultBoolProperty struct {
	value               bool
	name                string
	minInformationLevel InformationLevel
	propertytype        PropertyType
}

// MakeIntProperty
func MakeBoolProperty(name interface{}, value bool, minInformationLevel InformationLevel, pt PropertyType) *DefaultBoolProperty {
	switch name.(type) {
	case string:
		return &DefaultBoolProperty{
			name:                name.(string),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertytype:        pt,
		}
	default:
		return nil
	}
}

// implements IntProperty
func (d DefaultBoolProperty) Name(i InformationLevel) string {
	return d.name
}

func (d DefaultBoolProperty) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= d.minInformationLevel {
		return d.value
	}
	return nil
}

func (d DefaultBoolProperty) Get() interface{} {
	return d.value
}

func (d *DefaultBoolProperty) Set(p interface{}) {
	d.value = p.(bool)
}

func (d DefaultBoolProperty) Increase() {
	// do nothing
}

func (d DefaultBoolProperty) GetType() PropertyType {
	return d.propertytype
}

func (d DefaultBoolProperty) I() bool {
	return d.value
}

func (d *DefaultBoolProperty) SetI(f bool) {
	d.value = f
}

func (d DefaultBoolProperty) Duplicate() Property {
	return &DefaultBoolProperty{
		value:               d.value,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertytype:        d.propertytype,
	}
}

func (d DefaultBoolProperty) ApplyBuff(p Property) Property {
	res := d.Duplicate()
	res.Set(d.Get().(bool) && p.Get().(bool))
	return res
}
