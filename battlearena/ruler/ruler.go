package ruler

import (
	"fmt"
	"reflect"
	"time"

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
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	Controllers  map[uuid.UUID]actor.Communication
	Turner       turner.Turner
	Grid         *grid.Grid
	Entities     map[uuid.UUID]*entity.Entity
	CurrentState GameState
	logger       *logrus.Entry

	NbControllers           int
	NbEntitiesPerController int
}

func NewRuler() Ruler {
	r := Ruler{
		act:          actor.New("Ruler"),
		Controllers:  make(map[uuid.UUID]actor.Communication),
		Turner:       turner.NewTurner(),
		Grid:         nil,
		Entities:     make(map[uuid.UUID]*entity.Entity),
		CurrentState: WaitingForControllers,
	}
	r.act.SetReceiveMessageHandler(r.handleMessage)
	r.act.SetReplyMessageHandler(r.handleReply)
	r.logger = logrus.WithFields(logrus.Fields{
		"component": "Ruler",
		"name":      r.act.Name()})

	gg := gridgenerator.GridGenerator{}
	gg.Width = tools.NewIntRange(20, 50)
	gg.Length = tools.NewIntRange(20, 50)
	gg.Height = tools.NewIntRange(10, 15)
	gg.GenerateObstrcution = false
	gg.Type = gridgenerator.Flat
	gg.ObstructionRate = tools.NewIntRange(0, 0)

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

// Implement actor.Communication Interface

func (r *Ruler) NotifyActor(msg message.Message) {
	r.act.Notify(msg)
}

func (r *Ruler) SendActor(msg message.Message, callback chan message.Message) {
	r.act.Send(msg, callback)
}

func (r *Ruler) handleReply(msg message.Message) bool {
	// Handle reply from other actors
	switch msg.TargetMethod.(type) {
	}

	return false
}

func (r *Ruler) handleMessage(msg message.Message) bool {
	logrus.WithFields(logrus.Fields{
		"RequestID":    msg.RequestId.String()[0:8],
		"component":    "Ruler",
		"message_type": reflect.TypeOf(msg.TargetMethod).String()}).Info("Received message")

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
	case rulermethods.ControllerQuit:
		r.controllerQuit(msg, msg.TargetMethod.(rulermethods.ControllerQuit))
		return true
	case rulermethods.BattleStart:
		r.battleStart(msg)
		return true
	}
	return false
}

func (r *Ruler) addController(msg message.Message, req rulermethods.AddController) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})
	loclog.WithFields(logrus.Fields{
		"ControllerID": req.ControllerID.String()[0:8]}).Info("AddController")

	// reject if already registered
	if _, ok := r.Controllers[req.ControllerID]; ok {
		loclog.Warn("Controller already registered")
		r.act.Reply(msg.ReplyWithError(fmt.Sprintf("Controller %s already registered", req.ControllerID), "controller.already.registered"))
		return
	}

	// Create a new controller validator
	validator := controllervalidator.ControllerValidator{}
	req.Controller.NotifyActor(message.Create(nil, controllermethods.SetValidatorAndQueue{Validator: validator, Ruler: r}, nil))

	// Assign the controller to the designated number of entities
	for i := 0; i < r.NbEntitiesPerController; i++ {
		for _, e := range r.Entities {
			if e.ControllerID == uuid.Nil {
				e.ControllerID = req.ControllerID
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
		ControllerID: req.ControllerID,
		Grid:         r.Grid,
		TurnState:    r.Turner,
		Entities:     ent,
	}

	r.Controllers[req.ControllerID] = req.Controller

	// Controller added
	r.act.Reply(reply)

	// check if game is ready to begin.
	loclog.WithFields(logrus.Fields{
		"nbControllers": len(r.Controllers),
		"expected":      r.NbControllers}).Debug("Controller added")
	if len(r.Controllers) == r.NbControllers {
		go func() {
			<-time.After(2 * time.Second)
			r.NotifyActor(message.Create(nil, rulermethods.BattleStart{}, nil))
		}()
	}
}

func (r *Ruler) battleStart(msg message.Message) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})
	loclog.Info("Game started")
	r.CurrentState = InProgress
	// Select first entity to play
	entID := r.Turner.NextTurn()
	loclog.WithFields(logrus.Fields{
		"entityID": entID.String()[0:8]}).Info("First entity to play")

	ent := r.Entities[entID]
	ent.CurrentDelay = 0
	r.Entities[entID] = ent

	// notify all controller that the game is about to start.
	for _, c := range r.Controllers {
		c.NotifyActor(message.Create(nil, rulermethods.BattleStart{
			Turn: r.Turner.GetTurnState(),
		}, nil))
	}

	go func() {
		<-time.After(2 * time.Second)
		// notify controller of his turn
		r.Controllers[ent.ControllerID].NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
			Entity: *ent,
			Turn:   r.Turner.GetTurnState(),
		}, nil))
	}()

	r.act.NoReply(msg.Reply())
}

func (r *Ruler) getState(msg message.Message) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})
	loclog.Debug("GetState")
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

	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})
	loclog.Debug("GetGridState")
	// reply with the grid state (opaque to the client)
	reply := msg.Reply()
	reply.Content = rulermethods.GetGridStateReply{
		Grid: r.Grid,
	}

	r.act.Reply(reply)
}

func (r *Ruler) getEntitiesState(msg message.Message, req rulermethods.GetEntitiesState) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})
	loclog.Debug("GetEntitiesState")
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
		Entities: ent,
		Turn:     r.Turner.GetTurnState(),
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
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})

	loclog.WithFields(logrus.Fields{
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8],
		"path":         req.Path}).Debug("Controller move request")

	// check gamestate
	if r.CurrentState != InProgress {
		loclog.Error("Game is not in progress")
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	// Check if the controller is allowed to move the entity
	if !r.checkControllerForEntity(req.ControllerID, req.EntityID) {
		loclog.Error("Controller is not allowed to move this entity")
		r.act.Reply(msg.ReplyWithError("Controller is not allowed to move this entity", "entity.controller.missmatch"))
		return
	}
	if r.Turner.CurrentEntityTurn != req.EntityID {
		loclog.Error("It is not this entity turn")
		r.act.Reply(msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch"))
		return
	}

	// Check if the path is valid
	cells := r.Grid.CellsForPositions(req.Path)
	// a valid path is a path that contains only walkable cells and all cells must be adjascent
	for i, c := range cells {
		if c.Type == cell.Ground {
			if i > 0 && !cells[i-1].Position.IsAdjacent(c.Position, 2) {
				loclog.Error("Path is not valid")
				r.act.Reply(msg.ReplyWithError("Invalid path", "entity.path.invalid"))
				return
			}
		} else {
			loclog.Error("Path is not valid")
			r.act.Reply(msg.ReplyWithError("Invalid path(wrong type)", "entity.path.invalid"))
			return
		}
	}

	ent := r.Entities[req.EntityID]

	// Move the entity
	r.Grid.MoveEntity(ent.Position, req.Path[len(req.Path)-1], ent.ID)

	loclog.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8],
		"from":     ent.Position,
		"to":       req.Path[len(req.Path)-1]}).Debug("Entity moved")

	// update entity position
	ent.Position = req.Path[len(req.Path)-1]

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
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})
	loclog.WithFields(logrus.Fields{
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8],
		"target":       req.Target}).Debug("Controller attack request")

	// check gamestate

	if r.CurrentState != InProgress {
		loclog.Error("Game is not in progress")
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	// Check if the controller is allowed to use the entity
	if !r.checkControllerForEntity(req.ControllerID, req.EntityID) {
		loclog.Error("Controller is not allowed to use this entity")
		r.act.Reply(msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.missmatch"))
		return
	}
	if r.Turner.CurrentEntityTurn != req.EntityID {
		loclog.Error("It is not this entity turn")
		r.act.Reply(msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch"))
		return
	}

	// Check if the attack is valid
	target, found := r.Grid.CellAt(req.Target)
	if !found {
		loclog.Error("Target is not found")
		r.act.Reply(msg.ReplyWithError("Invalid target", "entity.attack.target.invalid"))
		return
	}

	if target.Type != cell.Ground {
		loclog.Error("Target is not valid")
		r.act.Reply(msg.ReplyWithError("Invalid attack", "entity.attack.invalid"))
		return
	}

	if target.EntityID == uuid.Nil {
		loclog.Error("Target has no entities")
		r.act.Reply(msg.ReplyWithError("Invalid attack", "entity.attack.invalid"))
		return
	}

	ent := r.Entities[req.EntityID]
	foe := r.Entities[target.EntityID]

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + 500

	loclog.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8],
		"foeID":    target.EntityID.String()[0:8]}).Debug("Entity attack")

	// Update the entity
	r.Entities[req.EntityID] = ent

	r.Grid.RemoveEntity(foe.Position)
	delete(r.Entities, foe.ID)
	r.Turner.RemoveEntity(foe.ID)

	loclog.WithFields(logrus.Fields{
		"entityID": foe.ID.String()[0:8],
		"position": foe.Position}).Info("##### Entity removed #####")

	// notify foe controller of the attack.

	foectrl, found := r.Controllers[foe.ControllerID]
	if !found {
		loclog.WithFields(logrus.Fields{
			"foeID":           foe.ID.String()[0:8],
			"foeControllerID": foe.ControllerID.String()[0:8]}).Error("Foe controller not found")

	} else {
		foectrl.NotifyActor(message.Create(nil, rulermethods.ControllerAttacked{
			Entity:   *foe,
			Attacker: *ent,
		}, nil))
	}

	// reply with the new entities state (opaque to the client)
	reply := msg.Reply()

	reply.Content = rulermethods.ControllerAttackReply{
		Entity: *ent,
	}

	r.act.Reply(reply)
}

func (r *Ruler) notifyController(msg message.Message, req rulermethods.NotifyController) {
	// Notify the controller of a message
}

func (r *Ruler) endOfTurn(msg message.Message, req rulermethods.EndOfTurn) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID":    msg.RequestId.String()[0:8],
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8]})
	loclog.Debug("End of turn request")

	// check gamestate
	if r.CurrentState != InProgress {
		loclog.Error("Game is not in progress")
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	if req.EntityID == uuid.Nil {
		loclog.Error("Can't work with nil entity")
		r.act.Reply(msg.ReplyWithError("Can't work with nil entity", "entity.nil"))
		return
	}
	if _, found := r.Entities[req.EntityID]; !found {
		loclog.WithFields(logrus.Fields{
			"RequestID":    msg.RequestId.String()[0:8],
			"controllerID": req.ControllerID.String()[0:8],
			"entityID":     req.EntityID.String()[0:8]}).Error("Can't work with absent entity")
		r.act.Reply(msg.ReplyWithError("Can't work with absent entity", "entity.absent"))
		return
	}

	// Check if the controller is allowed to end the turn
	// Check if the controller is allowed to use the entity
	if !r.checkControllerForEntity(req.ControllerID, req.EntityID) {
		loclog.WithFields(logrus.Fields{
			"controllerID":       req.ControllerID.String()[0:8],
			"entityID":           req.EntityID.String()[0:8],
			"entityControllerID": r.Entities[req.EntityID].ControllerID.String()[0:8],
		}).Error("Controller is not allowed to use this entity")
		r.act.Reply(msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.missmatch"))
		return
	}
	if r.Turner.CurrentEntityTurn != req.EntityID {
		loclog.Error("It is not this entity turn")

		r.act.Reply(msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch"))
		return
	}

	loclog.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8],
		"newDelay": r.Entities[req.EntityID].CurrentDelay + 500}).Debug("Entity end of turn, reinserting entity in the turn")
	r.Turner.AddEntity(req.EntityID, r.Entities[req.EntityID].CurrentDelay+500) // well ...end of turn delay

	ent := make([]entity.Entity, 0)
	// fill entities
	for _, e := range r.Entities {
		ent = append(ent, *e)
	}

	// check end of game. End of game is decided when all remaining entities are from the same controller
	remainingController := make(map[uuid.UUID]bool)
	remainingControllerID := uuid.Nil
	for _, ent := range r.Entities {
		if ent != nil {
			loclog.WithFields(logrus.Fields{
				"entityID":     ent.ID.String()[0:8],
				"controllerID": ent.ControllerID.String()[0:8],
				"delay":        ent.CurrentDelay,
				"position":     ent.Position}).Debug("Remaining entity")
			remainingController[ent.ControllerID] = true
			remainingControllerID = ent.ControllerID
		}
	}

	if len(remainingController) <= 1 {
		loclog.Info("##### END OF BATTLE! #####")
		// End of game
		r.CurrentState = Finished

		// notify all controllers of the end of the game
		for _, ctrl := range r.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerControllerID: remainingControllerID,
			}, nil))
		}
	} else {

		loclog.Info("##### END OF TURN #####")
		// notify all other controllers of the new turn
		for _, ctrl := range r.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.Turner.GetTurnState(),
			}, nil))
		}

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
				loclog.WithFields(logrus.Fields{
					"entityID":     nextTurnEnt.String()[0:8],
					"controllerID": ent.ControllerID.String()[0:8]}).Error("Controller not found")
			} else {

				go func() {
					<-time.After(2 * time.Second)
					// notify controller of his turn
					ctrl.NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
						Entity: *ent,
						Turn:   r.Turner.GetTurnState(),
					}, nil))
				}()
			}
		}
	}

	r.act.Reply(msg.Reply())
}

func (r *Ruler) controllerQuit(msg message.Message, req rulermethods.ControllerQuit) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID": msg.RequestId.String()[0:8]})
	loclog.WithFields(logrus.Fields{
		"controllerID": req.ControllerID.String()[0:8]}).Debug("Controller quit request")

	// expect this to happend when the game is done.

	if r.CurrentState != Finished {
		// Can't do much if it isn't the case.
	}

	_, found := r.Controllers[req.ControllerID]
	if found {
		delete(r.Controllers, req.ControllerID)

		// remove all entities from the controller
		for _, ent := range r.Entities {
			if ent.ControllerID == req.ControllerID {
				r.Grid.RemoveEntity(ent.Position)
				delete(r.Entities, ent.ID)
				r.Turner.RemoveEntity(ent.ID)
			}
		}

	}

	r.act.NoReply(msg.Reply())
}
