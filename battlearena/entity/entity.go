package entity

import (
	"time"

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
	Up    EntityOrientation = 0
	Right EntityOrientation = 1
	Down  EntityOrientation = 2
	Left  EntityOrientation = 3
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
	}
}
