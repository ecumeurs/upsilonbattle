package ruler

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
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
	act          actor.Actor
	Controllers  map[uuid.UUID]*controller.Controller
	Turner       turner.Turner
	Grid         *grid.Grid
	Entities     map[uuid.UUID]*entity.Entity
	CurrentState GameState
}

func NewRuler() Ruler {
	r := Ruler{
		Controllers:  make(map[uuid.UUID]*controller.Controller),
		Turner:       turner.NewTurner(),
		Grid:         nil,
		Entities:     make(map[uuid.UUID]*entity.Entity),
		CurrentState: WaitingForControllers,
	}
	r.act.SetReceiveMessageHandler(r.handleMessage)

	// Generate a new Map
	// Generate a new set of entities
	// Set a number of controllers to be expected and assign them to entities

	return r
}

func (r *Ruler) handleMessage(msg message.Message) {
	switch msg.TargetMethod.(type) {
	case rulermethods.AddController:
		r.AddController(msg, msg.TargetMethod.(rulermethods.AddController))
	case rulermethods.GetState:
		r.GetState(msg)
	case rulermethods.GetGridState:
		r.GetGridState(msg)
	case rulermethods.GetEntitiesState:
		r.GetEntitiesState(msg, msg.TargetMethod.(rulermethods.GetEntitiesState))
	case rulermethods.ControllerMove:
		r.ControllerMove(msg, msg.TargetMethod.(rulermethods.ControllerMove))
	case rulermethods.ControllerAttack:
		r.ControllerAttack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))
	case rulermethods.NotifyController:
		r.NotifyController(msg, msg.TargetMethod.(rulermethods.NotifyController))
	}
}

func (r *Ruler) AddController(msg message.Message, req rulermethods.AddController) {
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

func (r *Ruler) GetEntitiesState(msg message.Message, req rulermethods.GetEntitiesState) {
	// reply with the entities state (opaque to the client)
	// will also return the current turn state (as visible for controller)
	// especially if one of these entity should act
}

func (r *Ruler) ControllerMove(msg message.Message, req rulermethods.ControllerMove) {
	// Check if the controller is allowed to move the entity
	// Check if the path is valid
	// Move the entity
	// Compute the new delay
	// reply with the new entities state (opaque to the client)
}

func (r *Ruler) ControllerAttack(msg message.Message, req rulermethods.ControllerAttack) {
	// Check if the controller is allowed to attack with the entity
	// Check if the attack is valid
	// Attack
	// Compute the new delay
	// reply with the new entities state (opaque to the client)
	// reply with the oppenent entity state (opaque to the client)
}

func (r *Ruler) NotifyController(msg message.Message, req rulermethods.NotifyController) {
	// Notify the controller of a message
}

func (r *Ruler) EndOfTurn(msg message.Message, req rulermethods.EndOfTurn) {
	// Check if the controller is allowed to end the turn
	// Check if the entity is allowed to end the turn
	// End the turn
	// Based on the entity delay, compute the next turn
	// Trigger the next turn and notify all controllers
}
