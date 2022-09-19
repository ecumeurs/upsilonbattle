package rulermethods

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/google/uuid"
)

// Input struct

type AddController struct {
	Controller *controller.Controller
}

type GetState struct{}

type GetGridState struct {
	AsController uuid.UUID // if left to nil, will reply with full display of the grid.
}

type GetEntitiesState struct {
	AsController uuid.UUID // if left to nil, will reply with all informations.
}

type ControllerMove struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
	Path         []position.Position
}

type ControllerAttack struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
	Target       position.Position
}

type NotifyController struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
	Message      string
}

type EndOfTurn struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
}

// Output struct
