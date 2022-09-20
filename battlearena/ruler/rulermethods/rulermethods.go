package rulermethods

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
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

type AddControllerReply struct {
	ControllerID uuid.UUID
	Grid         *grid.Grid
	Entities     []entity.Entity
	TurnState    turner.Turner
}

type GetStateReply struct {
	GameState               string
	NbControllers           int
	NbControllersExpected   int
	NbEntitiesPerController int
	CurrentEntityTurn       uuid.UUID
}

type GetGridStateReply struct {
	Grid *grid.Grid
}

type GetEntitiesStateReply struct {
	Entities  []entity.Entity
	TurnState turner.Turner
}

type ControllerMoveReply struct {
	Entity entity.Entity
}

type ControllerAttackReply struct {
	Entity entity.Entity
}

type ControllerAttacked struct {
	Entity   entity.Entity
	Attacker entity.Entity
}

type ControllerNextTurn struct {
	Entity    entity.Entity
	TurnState turner.Turner
}
