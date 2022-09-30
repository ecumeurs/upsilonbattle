package skill

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect"
	"github.com/google/uuid"
)

type Skill struct {
	ID   uuid.UUID
	Name string

	Behavior def.BehaviorProperty

	Targeting map[string]property.Property
	Costs     map[string]property.Property
	Effect    effect.Effect
	Cooldown  int
}

// New
func New() Skill {
	return Skill{
		ID:        uuid.New(),
		Name:      "New Skill",
		Behavior:  *def.DefaultBehaviorProperty(),
		Targeting: make(map[string]property.Property),
		Costs:     make(map[string]property.Property),
		Effect:    *effect.New(),
		Cooldown:  0, // current cooldown value
	}
}

// NewSkill
func NewSkill(name string, targeting, cost map[string]property.Property, effect effect.Effect) Skill {
	return Skill{
		ID:        uuid.New(),
		Name:      name,
		Behavior:  *def.DefaultBehaviorProperty(),
		Targeting: targeting,
		Costs:     cost,
		Effect:    effect,
	}
}

// HasProperty
func (s Skill) HasProperty(p interface{}) bool {
	pstr := property.PropertyToString(p)
	for _, v := range s.Targeting {
		if v.Name(property.GameMaster) == pstr {
			return true
		}
	}
	for _, v := range s.Costs {
		if v.Name(property.GameMaster) == pstr {
			return true
		}
	}
	return s.Effect.HasProperty(p)
}

// HasProperties
func (s Skill) HasAnyProperties(p ...property.SkillProperties) bool {
	one := false
	for _, v := range p {
		if s.HasProperty(v) {
			one = true
		}
	}
	return one
}

func (s Skill) HasAllProperties(p ...property.SkillProperties) bool {
	for _, v := range p {
		if !s.HasProperty(v) {
			return false
		}
	}
	return true
}

// GetProperty
func (s Skill) GetProperty(p interface{}) property.Property {
	pstr := property.PropertyToString(p)
	for _, v := range s.Targeting {
		if v.Name(property.GameMaster) == pstr {
			return v
		}
	}
	for _, v := range s.Costs {
		if v.Name(property.GameMaster) == pstr {
			return v
		}
	}
	for _, v := range s.Effect.Properties {
		if v.Name(property.GameMaster) == pstr {
			return v
		}
	}
	return nil
}

// GetProperty
func (s Skill) GetPropertyI(p interface{}) property.IntProperty {
	pstr := property.PropertyToString(p)
	for _, v := range s.Targeting {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.IntProperty)
		}
	}
	for _, v := range s.Costs {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.IntProperty)
		}
	}
	for _, v := range s.Effect.Properties {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.IntProperty)
		}
	}
	return nil
}

// GetProperty
func (s Skill) GetPropertyF(p interface{}) property.FloatProperty {
	pstr := property.PropertyToString(p)
	for _, v := range s.Targeting {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.FloatProperty)
		}
	}
	for _, v := range s.Costs {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.FloatProperty)
		}
	}
	for _, v := range s.Effect.Properties {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.FloatProperty)
		}
	}
	return nil
}

// GetProperty
func (s Skill) GetPropertyC(p interface{}) property.IntCounterProperty {
	pstr := property.PropertyToString(p)
	for _, v := range s.Targeting {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.IntCounterProperty)
		}
	}
	for _, v := range s.Costs {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.IntCounterProperty)
		}
	}
	for _, v := range s.Effect.Properties {
		if v.Name(property.GameMaster) == pstr {
			return v.(property.IntCounterProperty)
		}
	}
	return nil
}

// IsDirect
func (s *Skill) IsDirect() bool {
	return s.Behavior.IsDirect()
}

// IsReaction
func (s *Skill) IsReaction() bool {
	return s.Behavior.IsReaction()
}

// IsPassive
func (s *Skill) IsPassive() bool {
	return s.Behavior.IsPassive()
}

// IsCounter
func (s *Skill) IsCounter() bool {
	return s.Behavior.IsCounter()
}

// IsTrap
func (s *Skill) IsTrap() bool {
	return s.Behavior.IsTrap()
}
