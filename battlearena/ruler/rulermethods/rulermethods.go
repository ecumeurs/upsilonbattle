package rulermethods

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/google/uuid"
)

// Input struct

// @spec-link [[api_ruler_methods]]
type AddController struct {
	Controller   actor.Communication
	ControllerID uuid.UUID
}

type ControllerBattleReady struct {
	ControllerID uuid.UUID
	actor.NoReply
}
type ControllerTurnReady struct {
	ControllerID uuid.UUID
	actor.NoReply
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

type ControllerUseSkill struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
	SkillID      uuid.UUID
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
	IsTimeout    bool
	TurnIndex    uint32
	actor.NoReply
}

type ControllerQuit struct {
	ControllerID uuid.UUID
	actor.NoReply
}

type ControllerForfeit struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
}

// Timeout is an internal notification triggered by the ShotClock.
// @spec-link [[api_ruler_methods]]
type Timeout struct {
	TurnIndex uint32
	actor.NoReply
}

// Output struct (often found in Content!)
// need to standardize this ...

type AddControllerReply struct {
	ControllerID uuid.UUID
	Grid         *grid.Grid
	Entities     []entity.Entity
	TurnState    turner.TurnState
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
	Entities []entity.Entity
	Turn     turner.TurnState
}

type GetBoardState struct {
	ActionContext interface{} // Parrotted back in reply to keep context
}

type GetBoardStateReply struct {
	Grid          *grid.Grid
	Entities      []entity.Entity
	TurnState     turner.TurnState
	WinnerTeamID  int
	Version       int64
	ActionContext interface{}
}

// especially for those ... Are they in the CallbackTarget field or content ? or both ? :'(

type ControllerMoveReply struct {
	Entity entity.Entity
}

type ControllerAttackReply struct {
	Entity entity.Entity
}

type ControllerUseSkillReply struct {
	Entity entity.Entity
}

// Broadcasted messages

type ControllerNextTurn struct {
	Entity  entity.Entity
	Turn    turner.TurnState
	Version int64
	actor.NoReply
}

type BattleStart struct {
	Turn    turner.TurnState
	Version int64
	actor.NoReply
}

type BattleEnd struct {
	WinnerTeamID int
	WinnerName   string
	Version      int64
	actor.NoReply
}

type ControllerSkillUsed struct {
	ControllerID        uuid.UUID // the controller that was affected
	EmitterControllerID uuid.UUID // the controller that skilled
	Entity              entity.Entity
	Emitter             entity.Entity
	SkillID             uuid.UUID
	Version             int64

	actor.NoReply
}

type ControllerAttacked struct {
	ControllerID         uuid.UUID // the controller that was attacked
	AttackerControllerID uuid.UUID // the controller that attacked
	Entity               entity.Entity
	Attacker             entity.Entity
	SkillID              uuid.UUID
	Damage               int
	PrevHP               int
	NewHP                int
	Version              int64

	actor.NoReply
}

type ControllerMoved struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
	Path         []position.Position
	Version      int64

	actor.NoReply
}

type ControllerPassed struct {
	ControllerID uuid.UUID
	EntityID     uuid.UUID
	Version      int64

	actor.NoReply
}

type EntitiesStateChanged struct {
	Entities []entity.Entity
	Turn     turner.TurnState
	Version  int64
	actor.NoReply
}

// Testing-only messages to avoid direct access to actor state from tests
// which causes data races.

type TestingDeleteEntity struct {
	EntityID uuid.UUID
}

type TestingGetState struct {
}

type TestingGetStateReply struct {
	CurrentEntityTurn uuid.UUID
	CurrentState      string // Use string to avoid circular dependency if needed, or just let it be.
	WinnerTeamID      int
}
