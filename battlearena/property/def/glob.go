package def

import "github.com/ecumeurs/upsilonbattle/battlearena/property"

// DefaultProperty(name interface{}) Property
func DefaultProperty(name interface{}) property.Property {
	switch convname := name.(type) {
	case property.SkillProperties:
		return SkillProperty(convname)
	case property.EntityProperties:
		return EntityProperty(convname)
	case property.ItemProperties:
		return ItemProperty(convname)
	default:
		return nil
	}
}
