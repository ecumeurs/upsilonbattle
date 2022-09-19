package ruler

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/controllervalidator"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type GameState int

const (
	WaitingForControllers GameState = 1
	InProgress            GameState = 2
	Finished              GameState = 3
)

type Ruler struct {
	act                  actor.Actor
	ControllerValidators map[uuid.UUID]controllervalidator.ControllerValidator
	Turner               turner.Turner
	Grid                 *grid.Grid
	Entities             map[uuid.UUID]*entity.Entity
	CurrentState         GameState
}

func NewRuler() Ruler {
	r := Ruler{
		ControllerValidators: make(map[uuid.UUID]controllervalidator.ControllerValidator),
		Turner:               turner.NewTurner(),
		Grid:                 nil,
		Entities:             make(map[uuid.UUID]*entity.Entity),
		CurrentState:         WaitingForControllers,
	}
	r.act.SetReceiveMessageHandler(r.handleMessage)

	// Generate a new Map
	// Generate a new set of entities
	// Set a number of controllers to be expected and assign them to entities

	return r
}

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

func (r *Ruler) handleMessage(msg message.Message) {
	switch msg.TargetMethod.(type) {
	case AddController:
		r.AddController(msg, msg.TargetMethod.(AddController))
	case GetState:
		r.GetState(msg)
	case GetGridState:
		r.GetGridState(msg)
	case GetEntitiesState:
		r.GetEntitiesState(msg, msg.TargetMethod.(GetEntitiesState))
	case ControllerMove:
		r.ControllerMove(msg, msg.TargetMethod.(ControllerMove))
	case ControllerAttack:
		r.ControllerAttack(msg, msg.TargetMethod.(ControllerAttack))
	case NotifyController:
		r.NotifyController(msg, msg.TargetMethod.(NotifyController))
	}
}

func (r *Ruler) AddController(msg message.Message, req AddController) {
	// Create a new controller validator
	// reply with current map state (as visible for controller)
	// reply with current entities state (as visible for controller)
	// reply with current turn state (as visible for controller)
	// if all controllers are ready, start the game
}

func (r *Ruler) GetState(msg message.Message) {
	// return current game state (waiting for controllers, in progress, finished)
	// will also return number of controllers expected, number of controllers ready, number of controllers not ready
	// will also return, when in progress, the current turn state (as visible for controller)
}

func (r *Ruler) GetGridState(msg message.Message) {
	// reply with the grid state (opaque to the client)
}

func (r *Ruler) GetEntitiesState(msg message.Message, req GetEntitiesState) {
	// reply with the entities state (opaque to the client)
	// will also return the current turn state (as visible for controller)
	// especially if one of these entity should act
}

func (r *Ruler) ControllerMove(msg message.Message, req ControllerMove) {
	// Check if the controller is allowed to move the entity
	// Check if the path is valid
	// Move the entity
	// Compute the new delay
	// reply with the new entities state (opaque to the client)
}

func (r *Ruler) ControllerAttack(msg message.Message, req ControllerAttack) {
	// Check if the controller is allowed to attack with the entity
	// Check if the attack is valid
	// Attack
	// Compute the new delay
	// reply with the new entities state (opaque to the client)
	// reply with the oppenent entity state (opaque to the client)
}

func (r *Ruler) NotifyController(msg message.Message, req NotifyController) {
	// Notify the controller of a message
}

func (r *Ruler) EndOfTurn(msg message.Message, req EndOfTurn) {
	// Check if the controller is allowed to end the turn
	// Check if the entity is allowed to end the turn
	// End the turn
	// Based on the entity delay, compute the next turn
	// Trigger the next turn and notify all controllers
}
