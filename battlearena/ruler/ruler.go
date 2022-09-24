package ruler

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/gridgenerator"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/controllervalidator"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rules"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ArenaState int

const (
	WaitingForControllers ArenaState = 1
	InProgress            ArenaState = 2
	Finished              ArenaState = 3
)

// Convert GameState to String
func (g ArenaState) String() string {
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
	ID           uuid.UUID
	act          *actor.Actor
	GameState    *rules.GameState
	CurrentState ArenaState
	logger       *logrus.Entry

	NbControllers           int
	NbEntitiesPerController int
}

func NewRuler() Ruler {
	r := Ruler{
		ID:           uuid.New(),
		act:          actor.New("Ruler"),
		CurrentState: WaitingForControllers,
	}
	r.GameState = rules.NewGameState(r.ID)
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

	r.GameState.Grid = gg.Generate()

	r.NbControllers = 2
	r.NbEntitiesPerController = tools.NewIntRange(2, 3).Random()
	nbEntities := r.NbEntitiesPerController * r.NbControllers

	for i := 0; i < nbEntities; i++ {
		e := entity.NewEntity()
		e.Type = entity.Character
		e.Name = fmt.Sprintf("Entity %d", i)
		e.CurrentDelay = tools.NewIntRange(1000, 2000).Random()
		e.Position = r.GameState.Grid.RandomPosition()
		r.GameState.Grid.MoveEntity(position.New(0, 0, 0), e.Position, e.ID)
		r.GameState.Entities[e.ID] = e
		r.GameState.Turner.AddEntity(e.ID, e.CurrentDelay)
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
	if _, ok := r.GameState.Controllers[req.ControllerID]; ok {
		loclog.Warn("Controller already registered")
		r.act.Reply(msg.ReplyWithError(fmt.Sprintf("Controller %s already registered", req.ControllerID), "controller.already.registered"))
		return
	}

	// Create a new controller validator
	validator := controllervalidator.ControllerValidator{}
	req.Controller.NotifyActor(message.Create(nil, controllermethods.SetValidatorAndQueue{Validator: validator, Ruler: r}, nil))

	// Assign the controller to the designated number of entities
	for i := 0; i < r.NbEntitiesPerController; i++ {
		for _, e := range r.GameState.Entities {
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
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	reply.Content = rulermethods.AddControllerReply{
		ControllerID: req.ControllerID,
		Grid:         r.GameState.Grid,
		TurnState:    r.GameState.Turner,
		Entities:     ent,
	}

	r.GameState.Controllers[req.ControllerID] = req.Controller

	// Controller added
	r.act.Reply(reply)

	// check if game is ready to begin.
	loclog.WithFields(logrus.Fields{
		"nbControllers": len(r.GameState.Controllers),
		"expected":      r.NbControllers}).Debug("Controller added")
	if len(r.GameState.Controllers) == r.NbControllers {
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
	entID := r.GameState.Turner.NextTurn()
	loclog.WithFields(logrus.Fields{
		"entityID": entID.String()[0:8]}).Info("First entity to play")

	ent := r.GameState.Entities[entID]
	ent.CurrentDelay = 0
	r.GameState.Entities[entID] = ent

	// notify all controller that the game is about to start.
	for _, c := range r.GameState.Controllers {
		c.NotifyActor(message.Create(nil, rulermethods.BattleStart{
			Turn: r.GameState.Turner.GetTurnState(),
		}, nil))
	}

	go func() {
		<-time.After(2 * time.Second)
		// notify controller of his turn
		r.GameState.Controllers[ent.ControllerID].NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
			Entity: ent,
			Turn:   r.GameState.Turner.GetTurnState(),
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
		NbControllers:           len(r.GameState.Controllers),
		NbControllersExpected:   r.NbControllers,
		NbEntitiesPerController: r.NbEntitiesPerController,
		CurrentEntityTurn:       r.GameState.Turner.CurrentEntityTurn,
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
		Grid: r.GameState.Grid,
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
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	reply.Content = rulermethods.GetEntitiesStateReply{
		Entities: ent,
		Turn:     r.GameState.Turner.GetTurnState(),
	}

	r.act.Reply(reply)
}

func (r *Ruler) controllerMove(msg message.Message, req rulermethods.ControllerMove) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID":    msg.RequestId.String()[0:8],
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8]})
	// check gamestate
	if r.CurrentState != InProgress {
		loclog.Error("Game is not in progress")
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	r.act.Reply(r.GameState.Move(msg, req))
}

func (r *Ruler) controllerAttack(msg message.Message, req rulermethods.ControllerAttack) {
	loclog := r.logger.WithFields(logrus.Fields{
		"RequestID":    msg.RequestId.String()[0:8],
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8],
		"target":       req.Target})
	// check gamestate

	if r.CurrentState != InProgress {
		loclog.Error("Game is not in progress")
		r.act.Reply(msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	r.act.Reply(r.GameState.Attack(msg, req))
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
	if _, found := r.GameState.Entities[req.EntityID]; !found {
		loclog.WithFields(logrus.Fields{
			"RequestID":    msg.RequestId.String()[0:8],
			"controllerID": req.ControllerID.String()[0:8],
			"entityID":     req.EntityID.String()[0:8]}).Error("Can't work with absent entity")
		r.act.Reply(msg.ReplyWithError("Can't work with absent entity", "entity.absent"))
		return
	}

	// Check if the controller is allowed to end the turn
	// Check if the controller is allowed to use the entity
	if !r.GameState.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		loclog.WithFields(logrus.Fields{
			"controllerID":       req.ControllerID.String()[0:8],
			"entityID":           req.EntityID.String()[0:8],
			"entityControllerID": r.GameState.Entities[req.EntityID].ControllerID.String()[0:8],
		}).Error("Controller is not allowed to use this entity")
		r.act.Reply(msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.missmatch"))
		return
	}
	if r.GameState.Turner.CurrentEntityTurn != req.EntityID {
		loclog.Error("It is not this entity turn")

		r.act.Reply(msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch"))
		return
	}

	loclog.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8],
		"newDelay": r.GameState.Entities[req.EntityID].CurrentDelay + 500}).Debug("Entity end of turn, reinserting entity in the turn")
	r.GameState.Turner.AddEntity(req.EntityID, r.GameState.Entities[req.EntityID].CurrentDelay+500) // well ...end of turn delay

	ent := make([]entity.Entity, 0)
	// fill entities
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	// check end of game. End of game is decided when all remaining entities are from the same controller
	remainingController := make(map[uuid.UUID]bool)
	remainingControllerID := uuid.Nil
	for _, ent := range r.GameState.Entities {
		loclog.WithFields(logrus.Fields{
			"entityID":     ent.ID.String()[0:8],
			"controllerID": ent.ControllerID.String()[0:8],
			"delay":        ent.CurrentDelay,
			"position":     ent.Position}).Debug("Remaining entity")
		remainingController[ent.ControllerID] = true
		remainingControllerID = ent.ControllerID
	}

	if len(remainingController) <= 1 {
		loclog.Info("##### END OF BATTLE! #####")
		// End of game
		r.CurrentState = Finished

		// notify all controllers of the end of the game
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerControllerID: remainingControllerID,
			}, nil))
		}
	} else {

		loclog.Info("##### END OF TURN #####")
		// notify all other controllers of the new turn
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.GameState.Turner.GetTurnState(),
			}, nil))
		}

		// Check if the entity is allowed to end the turn
		// End the turn
		// Based on the entity delay, compute the next turn
		// Trigger the next turn and notify all controllers

		nextTurnEnt := r.GameState.Turner.NextTurn()
		if nextTurnEnt == uuid.Nil {
			// No more turn, end of the game
		} else {
			// Notify the controller of the next turn
			ent := r.GameState.Entities[nextTurnEnt]
			ent.CurrentDelay = 0
			r.GameState.Entities[nextTurnEnt] = ent

			ctrl, found := r.GameState.Controllers[ent.ControllerID]
			if !found {
				loclog.WithFields(logrus.Fields{
					"entityID":     nextTurnEnt.String()[0:8],
					"controllerID": ent.ControllerID.String()[0:8]}).Error("Controller not found")
			} else {

				go func() {
					<-time.After(2 * time.Second)
					// notify controller of his turn
					ctrl.NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
						Entity: ent,
						Turn:   r.GameState.Turner.GetTurnState(),
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

	_, found := r.GameState.Controllers[req.ControllerID]
	if found {
		delete(r.GameState.Controllers, req.ControllerID)

		// remove all entities from the controller
		for _, ent := range r.GameState.Entities {
			if ent.ControllerID == req.ControllerID {
				r.GameState.Grid.RemoveEntity(ent.Position)
				delete(r.GameState.Entities, ent.ID)
				r.GameState.Turner.RemoveEntity(ent.ID)
			}
		}

	}

	r.act.NoReply(msg.Reply())
}
