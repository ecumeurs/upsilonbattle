package effect

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/google/uuid"
)

type Effect struct {
	Properties []property.Property
	Name       string
	CasterID   uuid.UUID
}

// New
func New() *Effect {
	return &Effect{
		Properties: []property.Property{},
		Name:       "New Effect",
		CasterID:   uuid.Nil,
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

// GetProperty
func (e Effect) GetProperty(p interface{}) property.Property {
	pstr := property.PropertyToString(p)
	for _, v := range e.Properties {
		if v.Name(property.GameMaster) == pstr {
			return v
		}
	}

	return nil
}

// GetProperty
func (e Effect) GetPropertyI(p interface{}) property.IntProperty {
	return e.GetProperty(p).(property.IntProperty)
}

// GetProperty
func (e Effect) GetPropertyF(p interface{}) property.FloatProperty {
	return e.GetProperty(p).(property.FloatProperty)
}

// GetProperty
func (e Effect) GetPropertyC(p interface{}) property.IntCounterProperty {
	return e.GetProperty(p).(property.IntCounterProperty)
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
