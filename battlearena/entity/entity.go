package entity

import (
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity/properties.go"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
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

type Entity struct {
	ID           uuid.UUID
	ControllerID uuid.UUID
	Type         EntityType
	LastActivity time.Time
	Position     position.Position
	Name         string
	CurrentDelay int
	Orientation  EntityOrientation
	Properties   map[string]properties.Property
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
		Properties:   make(map[string]properties.Property),
	}
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
func (e Entity) GetProperty(name string, defaultProperty properties.Property) properties.Property {
	if _, found := e.Properties[name]; found {
		return e.Properties[name]
	}
	return defaultProperty
}

// GetPropertyI will return the property with the given name, or default
func (e Entity) GetPropertyI(name string, defaultProperty properties.IntProperty) properties.IntProperty {
	if _, found := e.Properties[name]; found {
		return e.Properties[name].(properties.IntProperty)
	}
	return defaultProperty
}

// GetPropertyF will return the property with the given name, or default
func (e Entity) GetPropertyF(name string, defaultProperty properties.FloatProperty) properties.FloatProperty {
	if _, found := e.Properties[name]; found {
		return e.Properties[name].(properties.FloatProperty)
	}
	return defaultProperty
}

func (e Entity) UpdateProperty(p properties.Property) {
	if _, found := e.Properties[p.Name(properties.OwnController)]; found {
		e.Properties[p.Name(properties.OwnController)] = p
	}
}
