package entity

import (
	"fmt"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

type EntityType int

const (
	Character EntityType = 0
	Monster   EntityType = 1
	Others    EntityType = 2
)

type EntityOrientation int

const (
	Up    EntityOrientation = 0 // X:0 Y:1
	Right EntityOrientation = 1 // X:1 Y:0
	Down  EntityOrientation = 2 // X:0 Y:-1
	Left  EntityOrientation = 3 // X:-1 Y:0
)

// String
func (e EntityOrientation) String() string {
	return [...]string{"Up", "Right", "Down", "Left"}[e]
}

type Entity struct {
	ID           uuid.UUID
	ControllerID uuid.UUID
	Type         EntityType
	LastActivity time.Time
	Position     position.Position
	Name         string
	CurrentDelay int
	Orientation  EntityOrientation
	Properties   map[string]property.Property
	Buffs        []property.TemporaryProperties
	Skills       map[uuid.UUID]skill.Skill
}

// NewEntity
func New() Entity {
	return Entity{
		ID:           uuid.New(),
		ControllerID: uuid.Nil,
		Type:         Others,
		LastActivity: time.Now(),
		CurrentDelay: 0,
		Orientation:  Up,
		Properties:   make(map[string]property.Property),
		Buffs:        make([]property.TemporaryProperties, 0),
		Skills:       make(map[uuid.UUID]skill.Skill),
	}
}

// String
func (e Entity) String() string {
	return fmt.Sprintf("[E %s-%s]", e.ID.String()[0:8], e.Name)
}

func (e Entity) PrettyString() string {
	res := fmt.Sprintf("%s %s %s\n", e.String(), e.Position.String(), e.Orientation.String())
	for _, v := range e.Properties {
		res += fmt.Sprintf(" - %s\n", property.PrettyPrint(v, property.GameMaster))
	}

	return res
}

// FaceToward will change the orientation of the entity to face given position.
// Decide facing based on angle toward position, with 0 being facing straight up.
// Allow UP to be set from -45 to 45 degrees, RIGHT from 45 to 135, DOWN from 135 to 225, LEFT from 225 to 315.
func (e *Entity) FaceToward(p position.Position) {
	angle := e.Position.AngleTo(p)
	if angle < 45 || angle > 315 {
		e.Orientation = Up
	} else if angle < 135 {
		e.Orientation = Right
	} else if angle < 225 {
		e.Orientation = Down
	} else {
		e.Orientation = Left
	}
}

// getBasePropertyOrDefault
func (e Entity) getBasePropertyOrDefault(name interface{}) property.Property {
	nname := property.PropertyToString(name)

	var prop property.Property

	if _, found := e.Properties[nname]; found {
		prop = e.Properties[nname].Duplicate()
	} else {
		prop = def.DefaultProperty(name)
	}

	return prop
}

// GetProperty will return the property with the given name, or default
func (e Entity) GetProperty(name interface{}) property.Property {
	prop := e.getBasePropertyOrDefault(name)

	buffs := e.GetBuffsFor(name)

	for _, buff := range buffs {
		prop = prop.ApplyBuff(buff)
	}

	return prop
}

// GetPropertyI will return the property with the given name, or default
func (e Entity) GetPropertyI(name interface{}) property.IntProperty {
	return e.GetProperty(name).(property.IntProperty)
}

// GetPropertyF will return the property with the given name, or default
func (e Entity) GetPropertyF(name interface{}) property.FloatProperty {
	return e.GetProperty(name).(property.FloatProperty)
}

// GetPropertyC will return the property with the given name, or default
func (e Entity) GetPropertyC(name interface{}) property.IntCounterProperty {
	return e.GetProperty(name).(property.IntCounterProperty)
}

func (e *Entity) UpdateProperty(p property.Property) {
	e.Properties[p.Name(property.GameMaster)] = p
}

func (e *Entity) RegisterSkill(s skill.Skill) {
	e.Skills[s.ID] = s
}

func (e *Entity) RegisterBuff(b property.TemporaryProperties) {
	e.Buffs = append(e.Buffs, b)
}

func (e Entity) GetBuffsFor(name interface{}) []property.Property {
	nname := property.PropertyToString(name)

	res := make([]property.Property, 0)
	for _, v := range e.Buffs {
		if _, found := v.Properties[nname]; found {
			res = append(res, v.Properties[nname])
		}
	}
	return res
}

func (e *Entity) BuffTickDown() {
	nbbuf := make([]property.TemporaryProperties, 0)
	for _, buff := range e.Buffs {
		if !buff.TickDown() {
			nbbuf = append(nbbuf, buff)
		}
	}
	e.Buffs = nbbuf
}

func (e Entity) HasActed() bool {
	return e.GetProperty(property.HasActed).Get().(bool)
}

func (e Entity) HasMoved() bool {
	return e.GetProperty(property.HasMoved).Get().(bool)
}

func (e Entity) HasProperty(name interface{}) bool {
	_, found := e.Properties[property.PropertyToString(name)]
	return found
}

// RepsertPropertyValue will insert property if unknown!
func (e *Entity) RepsertPropertyValue(p interface{}, value interface{}) {
	prop := e.GetProperty(p)
	prop.Set(value)
	e.Properties[prop.Name(property.GameMaster)] = prop
}

// RepsertPropertyCMaxValue will insert property if unknown!
func (e *Entity) RepsertPropertyCMaxValue(p interface{}, maxvalue int) {
	prop := e.GetProperty(p).(property.IntCounterProperty)
	prop.SetMaxValue(maxvalue)
	e.Properties[prop.Name(property.GameMaster)] = prop
}

// UpdatePropertyValue Will only update value if known to the entity (wont affect buffs)
func (e *Entity) UpdatePropertyValue(p interface{}, value interface{}) {
	prop := e.GetProperty(p)
	prop.Set(value)
	e.UpdateProperty(prop)
}

// UpdatePropertyCMaxValue Will only update value if known to the entity (wont affect buffs)
func (e *Entity) UpdatePropertyCMaxValue(p interface{}, maxvalue int) {
	prop := e.GetProperty(p).(property.IntCounterProperty)
	prop.SetMaxValue(maxvalue)
	e.UpdateProperty(prop)
}
