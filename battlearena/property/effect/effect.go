package effect

import "github.com/ecumeurs/upsilonbattle/battlearena/property"

type Effect struct {
	Properties []property.Property
	Name       string
}

// New
func New() *Effect {
	return &Effect{
		Properties: []property.Property{},
		Name:       "New Effect",
	}
}

// HasProperty
func (e Effect) HasProperty(p interface{}) bool {
	pstr := property.PropertyToString(p)
	for _, v := range e.Properties {
		if v.Name(property.GameMaster) == pstr {
			return true
		}
	}
	return false
}

// HasPositiveProperty
func (s Effect) HasPositiveProperty(p interface{}) bool {
	pstr := property.PropertyToString(p)
	for _, v := range s.Properties {
		if v.Name(property.GameMaster) == pstr {
			if v.(property.IntProperty).I() >= 0 {
				return true
			}
		}
	}
	return false
}

// HasNegativeProperty
func (s Effect) HasNegativeProperty(p interface{}) bool {
	pstr := property.PropertyToString(p)
	for _, v := range s.Properties {
		if v.Name(property.GameMaster) == pstr {
			if v.(property.IntProperty).I() < 0 {
				return true
			}
		}
	}
	return false
}

// IsDamaging
func (s Effect) IsDamaging() bool {
	return (s.HasPositiveProperty(property.Damage) ||
		s.HasPositiveProperty(property.StunPower) ||
		s.HasPositiveProperty(property.PoisonPower) ||
		s.HasNegativeProperty(property.ShieldPower))
}

// IsHealing
func (s Effect) IsHealing() bool {
	return (!s.HasProperty(property.Damage) ||
		s.HasNegativeProperty(property.StunPower) ||
		s.HasNegativeProperty(property.PoisonPower) ||
		s.HasPositiveProperty(property.ShieldPower) ||
		s.HasPositiveProperty(property.Heal))
}

// IsOverTime (poison, stun, etc) (Buff/Curse)
func (s Effect) IsOverTime() bool {
	return s.HasPositiveProperty(property.Duration)
}
