package defaultproperty

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilontools/tools"
)

type DefaultIntProperty struct {
	value               int
	name                string
	minInformationLevel property.InformationLevel
	propertyType        property.PropertyType
}

// MakeIntProperty
func MakeIntProperty(name interface{}, value int, minInformationLevel property.InformationLevel, t property.PropertyType) *DefaultIntProperty {
	switch convname := name.(type) {
	case property.EntityProperties:
		return &DefaultIntProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        t,
		}
	case property.SkillProperties:
		return &DefaultIntProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        t,
		}
	case property.ItemProperties:
		return &DefaultIntProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        t,
		}
	default:
		return nil
	}
}

// implements IntProperty
func (d DefaultIntProperty) Name(i property.InformationLevel) string {
	if i >= d.minInformationLevel {
		return d.name
	}
	return ""
}

func (d DefaultIntProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
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

func (d DefaultIntProperty) GetType() property.PropertyType {
	return d.propertyType
}

func (d DefaultIntProperty) I() int {
	return d.value
}

func (d *DefaultIntProperty) SetI(i int) {
	d.value = i
}

func (d DefaultIntProperty) ApplyBuff(p property.Property) property.Property {
	res := d.Duplicate()
	res.Set(d.Get().(int) + p.Get().(int))
	return res
}

func (d DefaultIntProperty) UnapplyBuff(p property.Property) property.Property {
	res := d.Duplicate()
	res.Set(d.Get().(int) - p.Get().(int))
	return res
}

func (d DefaultIntProperty) Duplicate() property.Property {
	return &DefaultIntProperty{
		value:               d.value,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertyType:        d.propertyType,
	}
}

type DefaultIntCounterProperty struct {
	Value    int
	MaxValue int

	name                string
	minInformationLevel property.InformationLevel
	propertyType        property.PropertyType
}

// MakeIntProperty
func MakeIntCounterProperty(name interface{}, value, maxvalue int, minInformationLevel property.InformationLevel, t property.PropertyType) *DefaultIntCounterProperty {
	nname := property.PropertyToString(name)
	if nname == "" {
		return nil
	}

	return &DefaultIntCounterProperty{
		name:                nname,
		Value:               value,
		MaxValue:            maxvalue,
		minInformationLevel: minInformationLevel,
		propertyType:        t,
	}
}

func (d DefaultIntCounterProperty) GetValue() int {
	return d.Value
}

func (d DefaultIntCounterProperty) GetMaxValue() int {
	return d.MaxValue
}

func (d *DefaultIntCounterProperty) SetValue(i int) {
	d.Value = i
}
func (d *DefaultIntCounterProperty) SetMaxValue(i int) {
	d.MaxValue = i
}

// implements IntProperty
func (d DefaultIntCounterProperty) Name(i property.InformationLevel) string {
	if i >= d.minInformationLevel {
		return d.name
	}
	return ""
}

func (d DefaultIntCounterProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
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

func (d DefaultIntCounterProperty) GetType() property.PropertyType {
	return d.propertyType
}

func (d DefaultIntCounterProperty) I() int {
	return d.Value
}

func (d *DefaultIntCounterProperty) SetI(i int) {
	d.Value = i
}

func (d DefaultIntCounterProperty) ApplyBuff(p property.Property) property.Property {
	res := d.Duplicate().(*DefaultIntCounterProperty)
	res.Value = d.Value + p.(*DefaultIntCounterProperty).Value
	res.MaxValue = d.MaxValue + p.(*DefaultIntCounterProperty).MaxValue
	return res
}

func (d DefaultIntCounterProperty) UnapplyBuff(p property.Property) property.Property {
	res := d.Duplicate().(*DefaultIntCounterProperty)

	res.MaxValue = d.MaxValue - p.(*DefaultIntCounterProperty).MaxValue
	res.Value = tools.Max(d.Value-p.(*DefaultIntCounterProperty).Value, res.MaxValue)
	return res
}

func (d DefaultIntCounterProperty) Duplicate() property.Property {
	return &DefaultIntCounterProperty{
		Value:               d.Value,
		MaxValue:            d.MaxValue,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertyType:        d.propertyType,
	}
}

type DefaultFloatProperty struct {
	value               float64
	name                string
	minInformationLevel property.InformationLevel
	propertyType        property.PropertyType
}

// MakeIntProperty
func MakeFloatProperty(name interface{}, value float64, minInformationLevel property.InformationLevel, pt property.PropertyType) *DefaultFloatProperty {
	switch convname := name.(type) {
	case property.EntityProperties:
		return &DefaultFloatProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        pt,
		}
	case property.SkillProperties:
		return &DefaultFloatProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        pt,
		}
	case property.ItemProperties:
		return &DefaultFloatProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        pt,
		}
	default:
		return nil
	}
}

// implements IntProperty
func (d DefaultFloatProperty) Name(i property.InformationLevel) string {
	return d.name
}

func (d DefaultFloatProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
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

func (d DefaultFloatProperty) GetType() property.PropertyType {
	return d.propertyType
}

func (d DefaultFloatProperty) I() float64 {
	return d.value
}

func (d *DefaultFloatProperty) SetI(f float64) {
	d.value = f
}

func (d DefaultFloatProperty) Duplicate() property.Property {
	return &DefaultFloatProperty{
		value:               d.value,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertyType:        d.propertyType,
	}
}

func (d DefaultFloatProperty) ApplyBuff(p property.Property) property.Property {
	res := d.Duplicate()
	res.Set(d.Get().(float64) + p.Get().(float64))
	return res
}
func (d DefaultFloatProperty) UnapplyBuff(p property.Property) property.Property {
	res := d.Duplicate()
	res.Set(d.Get().(float64) - p.Get().(float64))
	return res
}

// Bool default value are essentially flags ...

type DefaultBoolProperty struct {
	value               bool
	name                string
	minInformationLevel property.InformationLevel
	propertyType        property.PropertyType
}

// MakeIntProperty
func MakeBoolProperty(name interface{}, value bool, minInformationLevel property.InformationLevel, pt property.PropertyType) *DefaultBoolProperty {
	switch convname := name.(type) {
	case property.EntityProperties:
		return &DefaultBoolProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        pt,
		}
	case property.SkillProperties:
		return &DefaultBoolProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        pt,
		}
	case property.ItemProperties:
		return &DefaultBoolProperty{
			name:                convname.String(),
			value:               value,
			minInformationLevel: minInformationLevel,
			propertyType:        pt,
		}
	default:
		return nil
	}
}

// implements IntProperty
func (d DefaultBoolProperty) Name(i property.InformationLevel) string {
	return d.name
}

func (d DefaultBoolProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
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

func (d DefaultBoolProperty) GetType() property.PropertyType {
	return d.propertyType
}

func (d DefaultBoolProperty) B() bool {
	return d.value
}

func (d *DefaultBoolProperty) SetB(f bool) {
	d.value = f
}

func (d DefaultBoolProperty) Duplicate() property.Property {
	return &DefaultBoolProperty{
		value:               d.value,
		name:                d.name,
		minInformationLevel: d.minInformationLevel,
		propertyType:        d.propertyType,
	}
}

func (d DefaultBoolProperty) ApplyBuff(p property.Property) property.Property {
	res := d.Duplicate()
	res.Set(d.Get().(bool) && p.Get().(bool))
	return res
}

func (d DefaultBoolProperty) UnapplyBuff(p property.Property) property.Property {
	res := d.Duplicate()
	res.Set(!p.Get().(bool))
	return res
}
