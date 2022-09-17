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

type Entity struct {
	Uuid           uuid.UUID
	ControllerUuid uuid.UUID
	Type           EntityType
	LastActivity   time.Time
	Position       position.Position
}
