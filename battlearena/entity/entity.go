package entity

import (
	"fmt"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
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
	Skills       []skill.Skill
}

// NewEntity
func NewEntity() Entity {
	return Entity{
		ID:           uuid.New(),
		ControllerID: uuid.Nil,
		Type:         Others,
		LastActivity: time.Now(),
		CurrentDelay: 0,
		Orientation:  Up,
		Properties:   make(map[string]property.Property),
		Buffs:        make([]property.TemporaryProperties, 0),
		Skills:       make([]skill.Skill, 0),
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

// GetProperty will return the property with the given name, or default
func (e Entity) GetProperty(name interface{}) property.Property {
	var nname string
	switch convname := name.(type) {
	case property.SkillProperties:
		nname = convname.String()
	case property.ItemProperties:
		nname = convname.String()
	case property.EntityProperties:
		nname = convname.String()
	}

	var prop property.Property

	if _, found := e.Properties[nname]; found {
		prop = e.Properties[nname].Duplicate()
	} else {
		prop = def.DefaultProperty(name)
	}

	buffs := e.GetBuffsFor(name)

	for _, buff := range buffs {
		prop = prop.ApplyBuff(buff)
	}

	return prop
}

// GetPropertyI will return the property with the given name, or default
func (e Entity) GetPropertyI(name interface{}) property.IntProperty {
	var nname string
	switch convname := name.(type) {
	case property.SkillProperties:
		nname = convname.String()
	case property.ItemProperties:
		nname = convname.String()
	case property.EntityProperties:
		nname = convname.String()
	}

	var prop property.Property

	if _, found := e.Properties[nname]; found {
		prop = e.Properties[nname].Duplicate()
	} else {
		prop = def.DefaultProperty(name)
	}

	buffs := e.GetBuffsFor(name)

	for _, buff := range buffs {
		prop = prop.ApplyBuff(buff)
	}

	return prop.(property.IntProperty)
}

// GetPropertyF will return the property with the given name, or default
func (e Entity) GetPropertyF(name interface{}) property.FloatProperty {
	var nname string
	switch convname := name.(type) {
	case property.SkillProperties:
		nname = convname.String()
	case property.ItemProperties:
		nname = convname.String()
	case property.EntityProperties:
		nname = convname.String()
	}

	var prop property.Property

	if _, found := e.Properties[nname]; found {
		prop = e.Properties[nname].Duplicate()
	} else {
		prop = def.DefaultProperty(name)
	}

	buffs := e.GetBuffsFor(name)

	for _, buff := range buffs {
		prop = prop.ApplyBuff(buff)
	}

	return prop.(property.FloatProperty)
}

func (e Entity) UpdateProperty(p property.Property) {
	if _, found := e.Properties[p.Name(property.GameMaster)]; found {
		e.Properties[p.Name(property.GameMaster)] = p
	}
}

func (e *Entity) RegisterSkill(s skill.Skill) {
	e.Skills = append(e.Skills, s)
}

func (e *Entity) RegisterBuff(b property.TemporaryProperties) {
	e.Buffs = append(e.Buffs, b)
}

func (e Entity) GetBuffsFor(name interface{}) []property.Property {
	var nname string
	switch convname := name.(type) {
	case property.SkillProperties:
		nname = convname.String()
	case property.ItemProperties:
		nname = convname.String()
	case property.EntityProperties:
		nname = convname.String()
	case string:
		nname = convname
	}

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
