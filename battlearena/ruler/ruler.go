package ruler

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/cell"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/gridgenerator"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/controllervalidator"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type GameState int

const (
	WaitingForControllers GameState = 1
	InProgress            GameState = 2
	Finished              GameState = 3
)

// Convert GameState to String
func (g GameState) String() string {
	switch g {
	case WaitingForControllers:
		return "WaitingForControllers"
	case InProgress:
		return "InProgress"
	case Finished:
		return "Finished"
	default:
		return fmt.Sprintf("Unknown GameState %d", g)
	}
}

type Ruler struct {
	act          *actor.Actor
	Controllers  map[uuid.UUID]*controller.Controller
	Turner       turner.Turner
	Grid         *grid.Grid
	Entities     map[uuid.UUID]*entity.Entity
	CurrentState GameState

	NbControllers           int
	NbEntitiesPerController int
}

func NewRuler() Ruler {
	r := Ruler{
		act:          actor.New("Ruler"),
		Controllers:  make(map[uuid.UUID]*controller.Controller),
		Turner:       turner.NewTurner(),
		Grid:         nil,
		Entities:     make(map[uuid.UUID]*entity.Entity),
		CurrentState: WaitingForControllers,
	}
	r.act.SetReceiveMessageHandler(r.handleMessage)
	r.act.SetReplyMessageHandler(r.handleReply)

	gg := gridgenerator.GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50)
	gg.Length = tools.NewIntRange(20, 50)
	gg.Height = tools.NewIntRange(10, 15)
	gg.GenerateObstrcution = true
	gg.Type = gridgenerator.Flat
	gg.ObstructionRate = tools.NewIntRange(10, 50)

	r.Grid = gg.Generate()

	r.NbControllers = 2
	r.NbEntitiesPerController = tools.NewIntRange(2, 3).Random()
	nbEntities := r.NbEntitiesPerController * r.NbControllers

	for i := 0; i < nbEntities; i++ {
		e := entity.NewEntity()
		e.Type = entity.Character
		e.Name = fmt.Sprintf("Entity %d", i)
		e.CurrentDelay = tools.NewIntRange(1000, 2000).Random()
		e.Position = r.Grid.RandomPosition()
		r.Grid.MoveEntity(position.New(0, 0, 0), e.Position, e.ID)
		r.Entities[e.ID] = &e
		r.Turner.AddEntity(e.ID, e.CurrentDelay)
	}

	r.act.Start()

	return r
}

// Access to the Message Queue ...
func (r *Ruler) GetMQ() *messagequeue.MessageQueue {
	return r.act.GetQueue()
}

func (r *Ruler) handleReply(msg message.Message) bool {
	// Handle reply from other actors
	switch msg.TargetMethod.(type) {
	}

	return false
}

func (r *Ruler) handleMessage(msg message.Message) bool {
	switch msg.TargetMethod.(type) {
	case rulermethods.AddController:
		r.addController(msg, msg.TargetMethod.(rulermethods.AddController))
		return true
	case rulermethods.GetState:
		r.getState(msg)
		return true
	case rulermethods.GetGridState:
		r.getGridState(msg)
		return true
	case rulermethods.GetEntitiesState:
		r.getEntitiesState(msg, msg.TargetMethod.(rulermethods.GetEntitiesState))
		return true
	case rulermethods.ControllerMove:
		r.controllerMove(msg, msg.TargetMethod.(rulermethods.ControllerMove))
		return true
	case rulermethods.ControllerAttack:
		r.controllerAttack(msg, msg.TargetMethod.(rulermethods.ControllerAttack))
		return true
	case rulermethods.NotifyController:
		r.notifyController(msg, msg.TargetMethod.(rulermethods.NotifyController))
		return true
	case rulermethods.EndOfTurn:
		r.endOfTurn(msg, msg.TargetMethod.(rulermethods.EndOfTurn))
		return true
	}
	return false
}

func (r *Ruler) addController(msg message.Message, req rulermethods.AddController) {
	// Create a new controller validator
	validator := controllervalidator.ControllerValidator{}
	req.Controller.NotifyActor(message.Create(nil, controllermethods.SetValidatorAndQueue{Validator: validator, Queue: r.act.GetQueue()}, nil))

	// Assign the controller to the designated number of entities
	for i := 0; i < r.NbEntitiesPerController; i++ {
		for _, e := range r.Entities {
			if e.ControllerID == uuid.Nil {
				e.ControllerID = req.Controller.ID
				break
			}
		}
	}

	// reply with current map state (as visible for controller)
	// reply with current entities state (as visible for controller)
	// reply with current turn state (as visible for controller)
	// if all controllers are ready, start the game

	reply := msg.Reply()
	ent := make([]entity.Entity, 0)
	// fill entities
	for _, e := range r.Entities {
		ent = append(ent, *e)
	}

	reply.Content = rulermethods.AddControllerReply{
		ControllerID: req.Controller.ID,
		Grid:         r.Grid,
		TurnState:    r.Turner,
		Entities:     ent,
	}

	r.act.Reply(reply)

	// check if game is ready to begin.
	if len(r.Controllers) == r.NbControllers {
		r.CurrentState = InProgress
		// Select first entity to play
		entID := r.Turner.NextTurn()
		ent := r.Entities[entID]
		ent.CurrentDelay = 0
		r.Entities[entID] = ent

		// notify all controller that the game is about to start.
		for _, c := range r.Controllers {
			c.NotifyActor(message.Create(controllermethods.Send{}, rulermethods.BattleStart{
				TurnState: r.Turner,
			}, nil))
		}

		// notify controller of his turn
		r.Controllers[ent.ControllerID].NotifyActor(message.Create(controllermethods.Send{}, rulermethods.ControllerNextTurn{
			Entity:    *ent,
			TurnState: r.Turner,
		}, nil))
	}
}

func (r *Ruler) getState(msg message.Message) {
	// return current game state (waiting for controllers, in progress, finished)
	// will also return number of controllers expected, number of controllers ready, number of controllers not ready
	// will also return, when in progress, the current turn state (as visible for controller)

	reply := msg.Reply()
	reply.Content = rulermethods.GetStateReply{
		GameState:               r.CurrentState.String(),
		NbControllers:           len(r.Controllers),
		NbControllersExpected:   r.NbControllers,
		NbEntitiesPerController: r.NbEntitiesPerController,
		CurrentEntityTurn:       r.Turner.CurrentEntityTurn,
	}

	r.act.Reply(reply)
}

func (r *Ruler) getGridState(msg message.Message) {
	// reply with the grid state (opaque to the client)
	reply := msg.Reply()
	reply.Content = rulermethods.GetGridStateReply{
		Grid: r.Grid,
	}

	r.act.Reply(reply)
}

func (r *Ruler) getEntitiesState(msg message.Message, req rulermethods.GetEntitiesState) {
	// reply with the entities state (opaque to the client)
	// will also return the current turn state (as visible for controller)
	// especially if one of these entity should act

	reply := msg.Reply()
	ent := make([]entity.Entity, 0)
	// fill entities
	for _, e := range r.Entities {
		ent = append(ent, *e)
	}

	reply.Content = rulermethods.GetEntitiesStateReply{
		Entities:  ent,
		TurnState: r.Turner,
	}

	r.act.Reply(reply)
}

func (r *Ruler) checkControllerForEntity(controllerID uuid.UUID, entityID uuid.UUID) bool {
	for _, e := range r.Entities {
		if e.ID == entityID {
			return e.ControllerID == controllerID
		}
	}
	return false
}

func (r *Ruler) controllerMove(msg message.Message, req rulermethods.ControllerMove) {
	// check gamestate
	if r.CurrentState != InProgress {
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	// Check if the controller is allowed to move the entity
	if !r.checkControllerForEntity(req.ControllerID, req.EntityID) {
		r.act.Reply(msg.ReplyWithError("Controller is not allowed to move this entity", "entity.controller.missmatch"))
		return
	}
	if r.Turner.CurrentEntityTurn != req.EntityID {
		r.act.Reply(msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch"))
		return
	}

	// Check if the path is valid
	cells := r.Grid.CellsForPositions(req.Path)
	// a valid path is a path that contains only walkable cells and all cells must be adjascent
	for i, c := range cells {
		if c.Type == cell.Ground {
			if i > 0 && !cells[i-1].Position.IsAdjacent(c.Position, 2) {
				r.act.Reply(msg.ReplyWithError("Invalid path", "entity.path.invalid"))
				return
			}
		} else {
			r.act.Reply(msg.ReplyWithError("Invalid path", "entity.path.invalid"))
			return
		}
	}

	ent := r.Entities[req.EntityID]

	// Move the entity
	r.Grid.MoveEntity(ent.Position, req.Path[len(req.Path)-1], ent.ID)

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + len(req.Path)*200

	// Update the entity
	r.Entities[req.EntityID] = ent

	// reply with the new entities state (opaque to the client)
	reply := msg.Reply()

	reply.Content = rulermethods.ControllerMoveReply{
		Entity: *ent,
	}

	r.act.Reply(reply)

}

func (r *Ruler) controllerAttack(msg message.Message, req rulermethods.ControllerAttack) {
	// check gamestate
	if r.CurrentState != InProgress {
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	// Check if the controller is allowed to move the entity
	if !r.checkControllerForEntity(req.ControllerID, req.EntityID) {
		r.act.Reply(msg.ReplyWithError("Controller is not allowed to move this entity", "entity.controller.missmatch"))
		return
	}
	if r.Turner.CurrentEntityTurn != req.EntityID {
		r.act.Reply(msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch"))
		return
	}

	// Check if the attack is valid
	target, found := r.Grid.CellAt(req.Target)
	if !found {
		r.act.Reply(msg.ReplyWithError("Invalid target", "entity.attack.target.invalid"))
		return
	}

	if target.Type != cell.Ground {
		r.act.Reply(msg.ReplyWithError("Invalid attack", "entity.attack.invalid"))
		return
	}

	if target.EntityID == uuid.Nil {
		r.act.Reply(msg.ReplyWithError("Invalid attack", "entity.attack.invalid"))
		return
	}

	ent := r.Entities[req.EntityID]
	foe := r.Entities[target.EntityID]

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + 500

	// Update the entity
	r.Entities[req.EntityID] = ent

	r.Grid.RemoveEntity(foe.Position)
	r.Entities[foe.ID] = nil
	r.Turner.RemoveEntity(foe.ID)

	// notify foe controller of the attack.

	foectrl, found := r.Controllers[foe.ControllerID]
	if !found {
		fmt.Println("foe controller not found")
	} else {
		foectrl.NotifyActor(message.Message{
			TargetMethod: controllermethods.Send{},
			Content: rulermethods.ControllerAttacked{
				Entity:   *foe,
				Attacker: *ent,
			},
		})
	}

	// reply with the new entities state (opaque to the client)
	reply := msg.Reply()

	reply.Content = rulermethods.ControllerMoveReply{
		Entity: *ent,
	}

	r.act.Reply(reply)
}

func (r *Ruler) notifyController(msg message.Message, req rulermethods.NotifyController) {
	// Notify the controller of a message
}

func (r *Ruler) endOfTurn(msg message.Message, req rulermethods.EndOfTurn) {
	// check gamestate
	if r.CurrentState != InProgress {
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	// Check if the controller is allowed to end the turn
	// Check if the controller is allowed to move the entity
	if !r.checkControllerForEntity(req.ControllerID, req.EntityID) {
		r.act.Reply(msg.ReplyWithError("Controller is not allowed to move this entity", "entity.controller.missmatch"))
		return
	}
	if r.Turner.CurrentEntityTurn != req.EntityID {
		r.act.Reply(msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch"))
		return
	}

	r.Turner.AddEntity(req.EntityID, r.Entities[req.EntityID].CurrentDelay)

	// Check if the entity is allowed to end the turn
	// End the turn
	// Based on the entity delay, compute the next turn
	// Trigger the next turn and notify all controllers

	nextTurnEnt := r.Turner.NextTurn()
	if nextTurnEnt == uuid.Nil {
		// No more turn, end of the game
	} else {
		// Notify the controller of the next turn
		ent := r.Entities[nextTurnEnt]
		ent.CurrentDelay = 0
		r.Entities[nextTurnEnt] = ent

		ctrl, found := r.Controllers[ent.ControllerID]
		if !found {
			fmt.Println("controller not found")
		} else {
			ctrl.NotifyActor(message.Message{
				TargetMethod: controllermethods.Send{},
				Content: rulermethods.ControllerNextTurn{
					Entity:    *ent,
					TurnState: r.Turner,
				},
			})
		}

	}

	ent := make([]entity.Entity, 0)
	// fill entities
	for _, e := range r.Entities {
		ent = append(ent, *e)
	}

	// notify all other controllers of the new turn
	for _, ctrl := range r.Controllers {
		ctrl.NotifyActor(message.Message{
			TargetMethod: controllermethods.Send{},
			Content: rulermethods.GetEntitiesStateReply{
				Entities:  ent,
				TurnState: r.Turner,
			},
		})
	}

	r.act.Reply(msg.Reply())

	// check end of game. End of game is decided when all remaining entities are from the same controller
	remainingController := make(map[uuid.UUID]bool)
	remainingControllerID := uuid.Nil
	for _, ent := range r.Entities {
		if ent != nil {
			remainingController[ent.ControllerID] = true
			remainingControllerID = ent.ControllerID
		}
	}

	if len(remainingController) == 1 {
		// End of game
		r.CurrentState = Finished

		// notify all controllers of the end of the game
		for _, ctrl := range r.Controllers {
			ctrl.NotifyActor(message.Message{
				TargetMethod: controllermethods.Send{},
				Content: rulermethods.BattleEnd{
					WinnerControllerID: remainingControllerID,
					WinnerName:         r.Controllers[remainingControllerID].ControllerName,
				},
			})
		}

	}
}
