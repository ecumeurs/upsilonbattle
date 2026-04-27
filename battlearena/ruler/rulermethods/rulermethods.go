package rulermethods

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/google/uuid"
)

// Input struct

// @spec-link [[api_ruler_methods]]
type CreditAward struct {
	PlayerID uuid.UUID `json:"player_id"`
	Amount   int       `json:"amount"`
	Source   string    `json:"source"`
}

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

type InternalTriggerFirstTurn struct {
	actor.NoReply
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

type ActionResult struct {
	Target       entity.Entity `json:"target"`
	TargetID     uuid.UUID     `json:"target_id"`
	Damage       int           `json:"damage"`
	Heal         int           `json:"heal"`
	PrevHP       int           `json:"prev_hp"`
	NewHP        int           `json:"new_hp"`
	CreditAwards []CreditAward `json:"credits,omitempty"`
}

type ControllerMoveReply struct {
	Entity entity.Entity `json:"entity"`
}

type ControllerPassReply struct {
	Entity  entity.Entity `json:"entity"`
	Version int64         `json:"version"`
}

type ControllerAttackReply struct {
	Attacker entity.Entity  `json:"attacker"`
	Results  []ActionResult `json:"results"`
}

type ControllerUseSkillReply struct {
	Attacker entity.Entity  `json:"attacker"`
	Results  []ActionResult `json:"results"`
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
	CreditAwards        []CreditAward `json:"credits,omitempty"`

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
	Dead				 bool // is set to true, mark this entity as dead. 
	Version              int64
	CreditAwards         []CreditAward `json:"credits,omitempty"`

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
