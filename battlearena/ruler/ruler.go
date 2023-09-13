package ruler

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/entitygenerator"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/gridgenerator"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
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
	ID uuid.UUID
	*actor.Actor
	GameState    *rules.GameState
	CurrentState ArenaState
	logger       *logrus.Entry

	NbControllers           int
	NbEntitiesPerController int

	ControllerBattleReady map[uuid.UUID]bool
}

func NewRuler() Ruler {
	id := uuid.New()
	r := Ruler{
		ID:                    id,
		Actor:                 actor.New("Ruler"),
		CurrentState:          WaitingForControllers,
		ControllerBattleReady: make(map[uuid.UUID]bool),
	}
	r.GameState = rules.New(r.ID)
	r.logger = logrus.WithFields(logrus.Fields{
		"component": "Ruler",
		"name":      r.Name()})

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
		e := entitygenerator.GenerateRandomEntity()
		e.Type = entity.Character
		e.Name = fmt.Sprintf("Entity %d", i)
		e.CurrentDelay = tools.NewIntRange(1000, 2000).Random()
		e.Position = r.GameState.Grid.RandomPosition()
		r.GameState.Grid.MoveEntity(position.New(0, 0, 0), e.Position, e.ID)
		r.GameState.Entities[e.ID] = e
		r.GameState.Turner.AddEntity(e.ID, e.CurrentDelay)
	}

	r.AddMethod(rulermethods.AddController{}, r.addController, nil)
	r.AddMethod(rulermethods.GetState{}, r.getState, nil)
	r.AddMethod(rulermethods.GetGridState{}, r.getGridState, nil)
	r.AddMethod(rulermethods.GetEntitiesState{}, r.getEntitiesState, nil)
	r.AddMethod(rulermethods.ControllerMove{}, r.controllerMove, nil)
	r.AddMethod(rulermethods.ControllerAttack{}, r.controllerAttack, nil)
	r.AddMethod(rulermethods.ControllerUseSkill{}, r.controllerUseSkill, nil)
	r.AddMethod(rulermethods.NotifyController{}, r.notifyController, nil)
	r.AddMethod(rulermethods.EndOfTurn{}, r.endOfTurn, nil)
	r.AddMethod(rulermethods.ControllerQuit{}, r.controllerQuit, nil)
	r.AddMethod(rulermethods.BattleStart{}, r.battleStart, nil)
	r.AddMethod(rulermethods.ControllerBattleReady{}, r.controllerBattleReady, nil)
	r.AddMethod(rulermethods.ControllerTurnReady{}, r.controllerTurnReady, nil)

	r.Start()

	return r
}

func (r *Ruler) PrintStack() {
	r.GetQueue().PrintStack()
}

func (r *Ruler) addController(msg *message.Message) bool {
	req := msg.Content.(rulermethods.AddController)
	r.RequestLogger.WithFields(logrus.Fields{
		"ControllerID": req.ControllerID.String()[0:8]}).Info("AddController")

	// reject if already registered
	if _, ok := r.GameState.Controllers[req.ControllerID]; ok {
		r.RequestLogger.Warn("Controller already registered")
		r.Reply(msg, msg.ReplyWithError(fmt.Sprintf("Controller %s already registered", req.ControllerID), "controller.already.registered"))
		return true
	}

	req.Controller.NotifyActor(message.Create(nil, controllermethods.SetQueue{Ruler: r}, nil))

	// Assign the controller to the designated number of entities
	for i := 0; i < r.NbEntitiesPerController; i++ {
		for idx, e := range r.GameState.Entities {
			if e.ControllerID == uuid.Nil {
				e.ControllerID = req.ControllerID
				r.GameState.Entities[idx] = e
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
	r.Reply(msg, reply)

	// check if game is ready to begin.
	r.RequestLogger.WithFields(logrus.Fields{
		"nbControllers": len(r.GameState.Controllers),
		"expected":      r.NbControllers}).Debug("Controller added")

	if len(r.GameState.Controllers) == r.NbControllers {
		r.NotifyActor(message.Create(nil, rulermethods.BattleStart{}, nil))
	}
	return true
}

func (r *Ruler) controllerBattleReady(msg *message.Message) bool {
	req := msg.Content.(rulermethods.ControllerBattleReady)
	r.RequestLogger.Info("ControllerBattleReady")
	r.ControllerBattleReady[req.ControllerID] = true
	if len(r.ControllerBattleReady) == r.NbControllers {
		// All controllers are ready, start the game

		entID := r.GameState.Turner.CurrentEntityTurn
		r.RequestLogger.WithFields(logrus.Fields{
			"entityID": entID.String()[0:8]}).Info("First entity to play")

		ent := r.GameState.Entities[entID]
		ent.CurrentDelay = 0
		// notify controller of his turn
		r.GameState.Controllers[ent.ControllerID].NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
			Entity: ent,
			Turn:   r.GameState.Turner.GetTurnState(),
		}, nil))

	}
	r.NoReply(msg)
	return true
}

func (r *Ruler) controllerTurnReady(msg *message.Message) bool {
	r.RequestLogger.Info("ControllerTurnReady")
	r.NoReply(msg)
	return true
}

func (r *Ruler) battleStart(msg *message.Message) bool {
	r.RequestLogger.Info("Game started")
	r.CurrentState = InProgress
	// Select first entity to play
	entID := r.GameState.Turner.NextTurn()
	r.RequestLogger.WithFields(logrus.Fields{
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

	// expect all controller to send a battle ready message when they are ready to play.

	r.NoReply(msg)
	return true
}

func (r *Ruler) getState(msg *message.Message) bool {
	r.RequestLogger.Debug("GetState")
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

	r.Reply(msg, reply)
	return true
}

func (r *Ruler) getGridState(msg *message.Message) bool {
	r.RequestLogger.Debug("GetGridState")
	// reply with the grid state (opaque to the client)
	reply := msg.Reply()
	reply.Content = rulermethods.GetGridStateReply{
		Grid: r.GameState.Grid,
	}

	r.Reply(msg, reply)
	return true
}

func (r *Ruler) getEntitiesState(msg *message.Message) bool {
	r.RequestLogger.Debug("GetEntitiesState")
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

	r.Reply(msg, reply)
	return true
}

func (r *Ruler) controllerMove(msg *message.Message) bool {
	req := msg.Content.(rulermethods.ControllerMove)
	// check gamestate
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		r.Reply(msg, msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return true
	}

	r.Reply(msg, r.GameState.Move(msg, req))
	return true
}

func (r *Ruler) controllerAttack(msg *message.Message) bool {
	req := msg.Content.(rulermethods.ControllerAttack)
	// check gamestate

	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		r.Reply(msg, msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return true
	}

	r.Reply(msg, r.GameState.Attack(msg, req))
	return true
}

func (r *Ruler) controllerUseSkill(msg *message.Message) bool {
	req := msg.Content.(rulermethods.ControllerUseSkill)
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		r.Reply(msg, msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return true
	}

	reply, damaged, affected := r.GameState.UseSkill(msg, req)

	r.Reply(msg, reply)

	for _, d := range damaged {
		foectrlid := d.ControllerID
		// notify foe controller of the attack.
		foectrl, found := r.GameState.Controllers[foectrlid]
		if !found {
			r.RequestLogger.WithFields(logrus.Fields{
				"foeControllerID": foectrlid.String()[0:8]}).Error("Foe controller not found")

		} else {
			foectrl.NotifyActor(message.Create(nil, d, nil))
		}
	}

	for _, d := range affected {
		targetctrlid := d.ControllerID
		// notify target controller of the skill use.
		targetctrl, found := r.GameState.Controllers[targetctrlid]
		if !found {
			r.RequestLogger.WithFields(logrus.Fields{
				"targetControllerID": targetctrlid.String()[0:8]}).Error("target controller not found")

		} else {
			targetctrl.NotifyActor(message.Create(nil, d, nil))
		}
	}

	return true
}

func (r *Ruler) notifyController(msg *message.Message) bool {
	// NOTE: Not implemented?
	// Notify the controller of a message

	r.NoReply(msg)
	return true
}

func (r *Ruler) endOfTurn(msg *message.Message) bool {
	req := msg.Content.(rulermethods.EndOfTurn)
	r.RequestLogger = r.RequestLogger.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8]})
	r.RequestLogger.Debug("End of turn request")

	// check gamestate
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		r.Reply(msg, msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return true
	}

	ok, reply := r.GameState.EndOfTurn(msg, req, r.GameState.Entities[req.EntityID])
	if !ok {
		r.Reply(msg, reply)
		return true
	}

	nextTurnEnt := r.GameState.Turner.NextTurn()
	if nextTurnEnt == uuid.Nil {
		r.RequestLogger.Info("##### END OF BATTLE! (WEIRD) #####")
	} else {
		// if entity is in gamestate.
		if beg, found := r.GameState.Entities[nextTurnEnt]; found {
			r.GameState.BeginingOfTurn(beg)
		}
	}

	// checks begining of turn for next entity.

	ent := make([]entity.Entity, 0)
	// fill entities
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	// check end of game. End of game is decided when all remaining entities are from the same controller
	remainingController := make(map[uuid.UUID]bool)
	remainingControllerID := uuid.Nil
	for _, ent := range r.GameState.Entities {
		r.RequestLogger.WithFields(logrus.Fields{
			"entityID":     ent.ID.String()[0:8],
			"controllerID": ent.ControllerID.String()[0:8],
			"delay":        ent.CurrentDelay,
			"position":     ent.Position}).Debug("Remaining entity")
		remainingController[ent.ControllerID] = true
		remainingControllerID = ent.ControllerID
	}

	if len(remainingController) <= 1 || nextTurnEnt == uuid.Nil {
		r.RequestLogger.Info("##### END OF BATTLE! #####")
		// End of game
		r.CurrentState = Finished

		// notify all controllers of the end of the game
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerControllerID: remainingControllerID,
			}, nil))
		}
	} else {

		r.RequestLogger.Info("##### END OF TURN #####")
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

		if nextTurnEnt == uuid.Nil {
			// No more turn, end of the game
		} else {
			// Notify the controller of the next turn
			ent := r.GameState.Entities[nextTurnEnt]
			ent.CurrentDelay = 0
			r.GameState.Entities[nextTurnEnt] = ent

			ctrl, found := r.GameState.Controllers[ent.ControllerID]
			if !found {
				r.RequestLogger.WithFields(logrus.Fields{
					"entityID":     nextTurnEnt.String()[0:8],
					"controllerID": ent.ControllerID.String()[0:8]}).Error("Controller not found")
			} else {
				ctrl.NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
					Entity: ent,
					Turn:   r.GameState.Turner.GetTurnState(),
				}, nil))
			}
		}
	}

	r.Reply(msg, msg.Reply())
	return true
}

func (r *Ruler) controllerQuit(msg *message.Message) bool {
	req := msg.Content.(rulermethods.ControllerQuit)
	r.RequestLogger.Debug("Controller quit notification")

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

	r.NoReply(msg)
	return true
}
