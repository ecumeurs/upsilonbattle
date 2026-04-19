package ruler

import (
	"fmt"

	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/entitygenerator"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rules"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilonmapmaker/gridgenerator"
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
	SetQueueAcks          map[uuid.UUID]bool

	shotClock         *time.Timer
	shotClockVersion  int64
	ShotClockDuration time.Duration
	firstTurnSent     bool
}

// NewCompleteRuler creates a new Ruler with a complete GameState.
// WARNING: Calling this DOES NOT start the actor. You MUST call Start() manually
// after initial setup is complete. Failure to do so will result in an unresponsive arena.
func NewCompleteRuler() *Ruler {
	tools.Seed()
	id := uuid.New()
	r := Ruler{
		ID:                    id,
		Actor:                 actor.New("Ruler"),
		CurrentState:          WaitingForControllers,
		ControllerBattleReady: make(map[uuid.UUID]bool),
		SetQueueAcks:          make(map[uuid.UUID]bool),
		ShotClockDuration:     30 * time.Second,
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

	r.init()
	return &r
}

// NewRuler creates a new Ruler with an empty GameState.
// WARNING: Calling this DOES NOT start the actor. You MUST call Start() manually
// after initial setup is complete (e.g., SetGrid, AddEntity).
// Once Start() is called, you must never modify GameState directly again.
func NewRuler(id uuid.UUID) *Ruler {
	tools.Seed()
	r := Ruler{
		ID:                    id,
		Actor:                 actor.New("Ruler"),
		CurrentState:          WaitingForControllers,
		ControllerBattleReady: make(map[uuid.UUID]bool),
		SetQueueAcks:          make(map[uuid.UUID]bool),
		ShotClockDuration:     30 * time.Second,
	}
	r.GameState = rules.New(r.ID)
	r.logger = logrus.WithFields(logrus.Fields{
		"component": "Ruler",
		"name":      r.Name()})

	r.init()
	return &r
}

func (r *Ruler) SetNbControllers(nb int) {
	r.NbControllers = nb
}

func (r *Ruler) SetGrid(g *grid.Grid) {
	r.GameState.Grid = g
}

func (r *Ruler) AddEntity(e entity.Entity) {
	e.CurrentDelay = tools.NewIntRange(1000, 1500).Random()

	// ISS-047: Guard against nil grid during initialization races or misconfiguration
	if r.GameState.Grid == nil {
		r.logger.WithFields(logrus.Fields{
			"entityID": e.ID.String()[0:8]}).Error("Cannot add entity: Grid is not initialized. Ensure SetGrid is called before AddEntity.")
		return
	}

	e.Position = r.GameState.Grid.RandomPosition()
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), e.Position, e.ID)
	r.GameState.Entities[e.ID] = e
	r.GameState.Turner.AddEntity(e.ID, e.CurrentDelay)
}

func (r *Ruler) init() {
	r.AddCallHandler(rulermethods.AddController{}, r.addController, nil)
	r.AddCallHandler(rulermethods.GetState{}, r.getState, nil)
	r.AddCallHandler(rulermethods.GetGridState{}, r.getGridState, nil)
	r.AddCallHandler(rulermethods.GetEntitiesState{}, r.getEntitiesState, nil)
	r.AddCallHandler(rulermethods.ControllerMove{}, r.controllerMove, nil)
	r.AddCallHandler(rulermethods.ControllerAttack{}, r.controllerAttack, nil)
	r.AddCallHandler(rulermethods.ControllerUseSkill{}, r.controllerUseSkill, nil)
	r.AddNotificationHandler(rulermethods.NotifyController{}, r.notifyController, nil)
	r.AddCallHandler(rulermethods.EndOfTurn{}, r.endOfTurn, nil)
	r.AddNotificationHandler(rulermethods.ControllerQuit{}, r.controllerQuit, nil)
	r.AddNotificationHandler(rulermethods.BattleStart{}, r.battleStart, nil)
	r.AddNotificationHandler(rulermethods.ControllerBattleReady{}, r.controllerBattleReady, nil)
	r.AddNotificationHandler(rulermethods.ControllerTurnReady{}, r.controllerTurnReady, nil)
	r.AddNotificationHandler(rulermethods.ControllerPassed{}, r.controllerPassed, nil)
	r.AddCallHandler(rulermethods.ControllerForfeit{}, r.controllerForfeit, nil)
	r.AddNotificationHandler(rulermethods.Timeout{}, r.timeout, nil)
	r.AddNotificationHandler(actor.ActorAboutToStop{}, r.actorAboutToStop, nil)
	r.AddReplyHandler(controllermethods.SetQueueReply{}, r.handleSetQueueReply, nil)

	// Testing-only handlers
	r.AddNotificationHandler(rulermethods.TestingDeleteEntity{}, r.testingDeleteEntity, nil)
	r.AddCallHandler(rulermethods.TestingGetState{}, r.testingGetState, nil)

	// r.Start() REMOVED: Callers must start manually after setup.
}

func (r *Ruler) PrintStack() {
	r.GetQueue().PrintStack()
}

// addController handles the addition of a controller to the battle.
func (r *Ruler) addController(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.AddController)
	r.RequestLogger.WithFields(logrus.Fields{
		"ControllerID": req.ControllerID.String()[0:8]}).Info("AddController")

	// reject if already registered
	if _, ok := r.GameState.Controllers[req.ControllerID]; ok {
		r.RequestLogger.Warn("Controller already registered")
		ctx.Reply(ctx.Msg.ReplyWithError(fmt.Sprintf("Controller %s already registered", req.ControllerID), "controller.already.registered"))
		return
	}

	// ISS-010: Guard against starting without a grid
	if r.GameState.Grid == nil {
		r.RequestLogger.Error("Cannot add controller: Grid is not initialized")
		ctx.Reply(ctx.Msg.ReplyWithError("Grid not initialized", "arena.not_ready.no_grid"))
		return
	}

	// @spec-link [[mech_controller_handshake]]
	// @spec-link [[mech_controller_communication_sequence]]
	req.Controller.SendActor(message.Create(nil, controllermethods.SetQueue{
		ControllerID: req.ControllerID,
		Ruler:        r,
	}, controllermethods.SetQueueReply{}), r.GetCallbackChan())

	// Assign the controller to the designated number of entities
	for i := 0; i < r.NbEntitiesPerController; i++ {
		for idx, e := range r.GameState.Entities {
			if e.ControllerID == uuid.Nil {
				e.ControllerID = req.ControllerID
				// @spec-link [[rule_team_mechanics]]
				// Surgical Team Fallback: Seed a unique team if one isn't already assigned (for tests/mock)
				if e.GetPropertyI(property.TeamID).I() == 0 {
					// Use Join Order as Team ID (1-indexed)
					e.RepsertPropertyValue(property.TeamID, len(r.GameState.Controllers)+1)
				}

				r.GameState.Entities[idx] = e
				break
			}
		}
	}

	reply := ctx.Msg.Reply()
	ent := make([]entity.Entity, 0)
	// fill entities
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	reply.Content = rulermethods.AddControllerReply{
		ControllerID: req.ControllerID,
		Grid:         r.GameState.Grid,
		TurnState:    r.GameState.Turner.GetTurnState(),
		Entities:     ent,
	}

	r.GameState.Controllers[req.ControllerID] = req.Controller
	r.SetQueueAcks[req.ControllerID] = false
	r.ControllerBattleReady[req.ControllerID] = false

	ctx.Reply(reply)

	// In rare cases (like tests with 1 controller), the handshake might already be complete
	r.checkBattleStart()
}

func (r *Ruler) handleSetQueueReply(ctx actor.ReplyContext) {
	msg := ctx.Msg.TargetMethod.(controllermethods.SetQueueReply)
	r.RequestLogger.WithFields(logrus.Fields{
		"ControllerID": msg.ControllerID.String()[0:8]}).Info("SetQueue Handshake ACK received")

	r.SetQueueAcks[msg.ControllerID] = true
	r.checkBattleStart()
}

func (r *Ruler) checkBattleStart() {
	if r.CurrentState != WaitingForControllers {
		return
	}
	if r.allControllersAckedSetQueue() && r.isBattleReadyToStart() {
		r.RequestLogger.Info("All handshakes complete, scheduling BattleStart")
		// Delay slightly to ensure any other concurrent initialization completes
		r.SelfNotifyDelayed(rulermethods.BattleStart{}, 20*time.Millisecond)
	}
}

func (r *Ruler) testingDeleteEntity(ctx actor.NotificationContext) {
	msg := ctx.Msg.TargetMethod.(rulermethods.TestingDeleteEntity)
	r.logger.WithField("entityID", msg.EntityID.String()[0:8]).Info("Testing-only removal of entity")
	delete(r.GameState.Entities, msg.EntityID)
	r.GameState.Turner.RemoveEntity(msg.EntityID)
}

func (r *Ruler) testingGetState(ctx actor.CallContext) {
	ctx.Reply(message.Create(nil, rulermethods.TestingGetStateReply{
		CurrentEntityTurn: r.GameState.Turner.CurrentEntityTurn,
		CurrentState:      r.CurrentState.String(),
		WinnerTeamID:      r.GameState.WinnerTeamID,
	}, nil))
}

func (r *Ruler) isBattleReadyToStart() bool {
	if r.GameState.Grid == nil {
		r.logger.Debug("isBattleReadyToStart: Grid is nil")
		return false
	}
	if len(r.GameState.Controllers) != r.NbControllers {
		r.logger.WithFields(logrus.Fields{
			"current": len(r.GameState.Controllers),
			"target":  r.NbControllers,
		}).Debug("isBattleReadyToStart: Not all controllers added yet")
		return false
	}
	// All controllers must have sent ControllerBattleReady
	for id, ready := range r.ControllerBattleReady {
		if !ready {
			r.logger.WithField("controllerID", id.String()[0:8]).Debug("isBattleReadyToStart: Controller not ready yet")
			return false
		}
	}
	// Map must have all controllers
	if len(r.ControllerBattleReady) < r.NbControllers {
		r.logger.Debug("isBattleReadyToStart: ControllerBattleReady map too small")
		return false
	}

	return true
}

func (r *Ruler) allControllersAckedSetQueue() bool {
	if len(r.SetQueueAcks) < r.NbControllers {
		r.logger.WithFields(logrus.Fields{
			"current": len(r.SetQueueAcks),
			"target":  r.NbControllers,
		}).Debug("allControllersAckedSetQueue: SetQueueAcks map too small")
		return false
	}
	for id, acked := range r.SetQueueAcks {
		if !acked {
			r.logger.WithField("controllerID", id.String()[0:8]).Debug("allControllersAckedSetQueue: ACK missing")
			return false
		}
	}
	return true
}

// @spec-link [[rule_battle_readiness]]
func (r *Ruler) isBattleReadyToExecute() bool {
	if r.CurrentState != InProgress {
		return false
	}
	if !r.isBattleReadyToStart() {
		return false
	}
	if !r.allControllersAckedSetQueue() {
		return false
	}
	if len(r.ControllerBattleReady) != r.NbControllers {
		return false
	}
	return true
}

func (r *Ruler) triggerFirstTurn() {
	if r.firstTurnSent {
		return
	}

	entID := r.GameState.Turner.CurrentEntityTurn
	if entID == uuid.Nil {
		r.RequestLogger.Info("Picking first entity to play")
		entID = r.GameState.Turner.NextTurn()
	}

	// GUARD: Ensure the entity still exists and is controlled
	ent, ok := r.GameState.Entities[entID]
	if !ok || ent.ControllerID == uuid.Nil {
		r.RequestLogger.WithFields(logrus.Fields{
			"entityID": entID.String()}).Debug("Current entity turn is missing or uncontrolled during ready check (waiting for full assignment)")
		return
	}

	r.RequestLogger.WithFields(logrus.Fields{
		"entityID": entID.String()[0:8]}).Info("Handing turn to first entity")

	r.firstTurnSent = true
	r.GameState.IncTurn() // @spec-link [[mech_game_state_versioning]]
	// @spec-link [[rule_turn_clock]]
	r.startShotClock()

	r.RequestLogger.WithFields(logrus.Fields{
		"entityID": entID.String()[0:8]}).Info("First entity to play")

	ent.CurrentDelay = 0
	r.GameState.Entities[entID] = ent

	r.GameState.Controllers[ent.ControllerID].NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
		Entity: ent,
		Turn:   r.GameState.Turner.GetTurnState(),
	}, nil))
}

func (r *Ruler) controllerBattleReady(ctx actor.NotificationContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerBattleReady)
	r.RequestLogger.Info("ControllerBattleReady")

	// GUARD: Don't process ready checks if the game is already finished
	if r.CurrentState == Finished {
		r.RequestLogger.Warn("Received ControllerBattleReady after battle finished")
		return
	}

	r.ControllerBattleReady[req.ControllerID] = true
	r.checkBattleStart()

	if r.isBattleReadyToExecute() {
		r.triggerFirstTurn()
	}
}

func (r *Ruler) controllerTurnReady(ctx actor.NotificationContext) {
	r.RequestLogger.Info("ControllerTurnReady")
}

func (r *Ruler) battleStart(ctx actor.NotificationContext) {
	r.RequestLogger.Info("Processing BattleStart internal notification")
	r.CurrentState = InProgress

	// @spec-link [[rule_turn_clock]]
	r.startShotClock()

	r.RequestLogger.Info("Broadcasting BattleStart to all controllers")
	for id, c := range r.GameState.Controllers {
		r.RequestLogger.WithFields(logrus.Fields{"target": id}).Debug("Sending BattleStart")
		c.NotifyActor(message.Create(nil, rulermethods.BattleStart{
			Turn: r.GameState.Turner.GetTurnState(),
		}, nil))
	}

	if r.isBattleReadyToExecute() {
		r.RequestLogger.Info("All controllers already ready, triggering first turn immediately")
		r.triggerFirstTurn()
	} else {
		r.RequestLogger.Info("Waiting for controllers to signal readiness before first turn")
	}
}

func (r *Ruler) getState(ctx actor.CallContext) {
	r.RequestLogger.Debug("GetState")
	reply := ctx.Msg.Reply()
	reply.Content = rulermethods.GetStateReply{
		GameState:               r.CurrentState.String(),
		NbControllers:           len(r.GameState.Controllers),
		NbControllersExpected:   r.NbControllers,
		NbEntitiesPerController: r.NbEntitiesPerController,
		CurrentEntityTurn:       r.GameState.Turner.CurrentEntityTurn,
	}

	ctx.Reply(reply)
}

func (r *Ruler) getGridState(ctx actor.CallContext) {
	r.RequestLogger.Debug("GetGridState")
	reply := ctx.Msg.Reply()
	reply.Content = rulermethods.GetGridStateReply{
		Grid: r.GameState.Grid,
	}

	ctx.Reply(reply)
}

func (r *Ruler) getEntitiesState(ctx actor.CallContext) {
	r.RequestLogger.Debug("GetEntitiesState")

	reply := ctx.Msg.Reply()
	ent := make([]entity.Entity, 0)
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	reply.Content = rulermethods.GetEntitiesStateReply{
		Entities: ent,
		Turn:     r.GameState.Turner.GetTurnState(),
	}

	ctx.Reply(reply)
}

// @spec-link [[mech_action_economy_action_cost_rules]]
func (r *Ruler) controllerMove(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerMove)
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	reply := r.GameState.Move(ctx.Msg, req)
	ctx.Reply(reply)

	if !reply.HasError {
		ent := make([]entity.Entity, 0, len(r.GameState.Entities))
		for _, e := range r.GameState.Entities {
			ent = append(ent, e)
		}
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.GameState.Turner.GetTurnState(),
			}, nil))
		}
	}
}

// @spec-link [[mech_action_economy_action_cost_rules]]
func (r *Ruler) controllerAttack(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerAttack)
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	reply := r.GameState.Attack(ctx.Msg, req)
	ctx.Reply(reply)

	if !reply.HasError {
		ent := make([]entity.Entity, 0, len(r.GameState.Entities))
		for _, e := range r.GameState.Entities {
			ent = append(ent, e)
		}
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.GameState.Turner.GetTurnState(),
			}, nil))
		}
	}
}

func (r *Ruler) controllerUseSkill(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerUseSkill)
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	reply, damaged, affected := r.GameState.UseSkill(ctx.Msg, req)

	ctx.Reply(reply)

	for _, d := range damaged {
		foectrlid := d.ControllerID
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
		targetctrl, found := r.GameState.Controllers[targetctrlid]
		if !found {
			r.RequestLogger.WithFields(logrus.Fields{
				"targetControllerID": targetctrlid.String()[0:8]}).Error("target controller not found")
		} else {
			targetctrl.NotifyActor(message.Create(nil, d, nil))
		}
	}

	if !reply.HasError {
		ent := make([]entity.Entity, 0, len(r.GameState.Entities))
		for _, e := range r.GameState.Entities {
			ent = append(ent, e)
		}
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.GameState.Turner.GetTurnState(),
			}, nil))
		}
	}
}

func (r *Ruler) notifyController(ctx actor.NotificationContext) {
}

func (r *Ruler) controllerPassed(ctx actor.NotificationContext) {
	r.RequestLogger.Debug("Controller passed notification ignored by ruler")
}

// @spec-link [[mech_action_economy_action_cost_rules]]
func (r *Ruler) endOfTurn(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.EndOfTurn)
	r.RequestLogger = r.RequestLogger.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8]})
	r.RequestLogger.Debug("End of turn request")

	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	// @spec-link [[rule_turn_clock]]
	// GUARD: If this is a timeout from an old turn, ignore it.
	if req.IsTimeout && req.TurnIndex != 0 && req.TurnIndex != r.GameState.GetTurn() {
		r.logger.WithFields(logrus.Fields{
			"reqTurn":     req.TurnIndex,
			"currentTurn": r.GameState.GetTurn()}).Debug("Ignoring late timeout message")
		ctx.Reply(ctx.Msg.Reply())
		return
	}

	r.stopShotClock()

	ok, reply := r.GameState.EndOfTurn(ctx.Msg, req, r.GameState.Entities[req.EntityID])
	if !ok {
		ctx.Reply(reply)
		return
	}

	nextTurnEnt := r.GameState.Turner.NextTurn()

	// @spec-link [[mech_initiative_active_state]]
	// GUARD: Skip any entities that were killed mid-turn and removed from
	// gs.Entities before NextTurn() could be called (e.g., killed by an
	// attack action right before the acting entity passed). Without this,
	// the ruler hands ControllerNextTurn to a dead entity (zero-value
	// ControllerID == uuid.Nil), no controller receives the notification,
	// and the battle permanently hangs.
	for nextTurnEnt != uuid.Nil {
		if _, alive := r.GameState.Entities[nextTurnEnt]; alive {
			break
		}
		r.RequestLogger.WithFields(logrus.Fields{
			"skippedEntityID": nextTurnEnt.String()[0:8],
		}).Warn("Next-turn entity was already dead; skipping to next in queue")
		nextTurnEnt = r.GameState.Turner.NextTurn()
	}

	if nextTurnEnt == uuid.Nil {
		r.RequestLogger.Info("##### END OF BATTLE! (WEIRD) #####")
	} else {
		if beg, found := r.GameState.Entities[nextTurnEnt]; found {
			r.GameState.BeginingOfTurn(beg)
		}
	}

	ent := make([]entity.Entity, 0)
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	if nextTurnEnt == uuid.Nil {
		r.evaluateVictory(nextTurnEnt)
	} else {
		remainingTeams := make(map[int]bool)
		for _, ent := range r.GameState.Entities {
			remainingTeams[ent.GetPropertyI(property.TeamID).I()] = true
		}

		if len(remainingTeams) <= 1 {
			r.evaluateVictory(nextTurnEnt)
		} else {
			// @spec-link [[rule_turn_clock]]
			r.RequestLogger.Info("##### END OF TURN #####")
			for _, ctrl := range r.GameState.Controllers {
				ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
					Entities: ent,
					Turn:     r.GameState.Turner.GetTurnState(),
				}, nil))
			}

			if nextTurnEnt != uuid.Nil {
				ent := r.GameState.Entities[nextTurnEnt]
				ent.CurrentDelay = 0
				r.GameState.Entities[nextTurnEnt] = ent

				ctrl, found := r.GameState.Controllers[ent.ControllerID]
				if !found {
					r.RequestLogger.WithFields(logrus.Fields{
						"entityID":     nextTurnEnt.String()[0:8],
						"controllerID": ent.ControllerID.String()[0:8]}).Error("Controller not found")
				} else {
					// @spec-link [[rule_turn_clock]]
					r.startShotClock()
					ctrl.NotifyActor(message.Create(nil, rulermethods.ControllerNextTurn{
						Entity: ent,
						Turn:   r.GameState.Turner.GetTurnState(),
					}, nil))
				}
			}
		}

	}

	ctx.Reply(ctx.Msg.Reply())
}

// controllerForfeit handles the forfeiture of a controller.
// @spec-link [[rule_forfeit_battle]]
func (r *Ruler) controllerForfeit(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerForfeit)
	r.RequestLogger.WithFields(logrus.Fields{
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8]}).Info("ControllerForfeit")

	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	_, winnerTeamID, finished := r.GameState.Forfeit(req.ControllerID)

	if finished {
		r.CurrentState = Finished
		r.GameState.WinnerTeamID = winnerTeamID
		r.RequestLogger.Info("##### END OF BATTLE! #####")

		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerTeamID: winnerTeamID,
			}, nil))
		}
	}

	ctx.Reply(ctx.Msg.Reply())
}

func (r *Ruler) controllerQuit(ctx actor.NotificationContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerQuit)
	r.RequestLogger.Debug("Controller quit notification")

	if r.CurrentState != Finished {
	}

	_, found := r.GameState.Controllers[req.ControllerID]
	if found {
		delete(r.GameState.Controllers, req.ControllerID)
		r.RequestLogger.Info("Controller removed from match")

		for id, ent := range r.GameState.Entities {
			if ent.ControllerID == req.ControllerID {
				r.GameState.Grid.RemoveEntity(ent.Position)
				delete(r.GameState.Entities, id)
				r.GameState.Turner.RemoveEntity(id)
			}
		}

		// Re-evaluate victory if the match is still active
		if r.CurrentState != Finished {
			r.evaluateVictory(r.GameState.Turner.CurrentEntityTurn)
		}
	}
}

func (r *Ruler) evaluateVictory(nextTurnEnt uuid.UUID) {
	remainingTeams := make(map[int]bool)
	winningTeamID := 0
	for _, ent := range r.GameState.Entities {
		remainingTeams[ent.GetPropertyI(property.TeamID).I()] = true
		winningTeamID = ent.GetPropertyI(property.TeamID).I()
	}

	if len(remainingTeams) <= 1 || nextTurnEnt == uuid.Nil {
		// @spec-link [[rule_team_mechanics]]
		// @spec-link [[spec_match_format_win_condition_rule]]
		r.RequestLogger.Info("##### END OF BATTLE! #####")
		r.CurrentState = Finished
		r.GameState.WinnerTeamID = winningTeamID
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerTeamID: winningTeamID,
			}, nil))
		}
	}
}

// startShotClock initializes and starts the turn timer.
// @spec-link [[rule_turn_clock]]
// @spec-link [[mech_action_economy_timeout_penalty_rules]]
func (r *Ruler) startShotClock() {
	r.stopShotClock()

	if r.ShotClockDuration <= 0 {
		return
	}

	// Capture the turn index this shot clock is intended for
	turn := uint32(r.GameState.GetTurn())

	r.logger.WithFields(logrus.Fields{
		"turn":    turn,
		"version": fmt.Sprintf("%d.%d", turn, r.GameState.GetAction()),
		"timeout": r.ShotClockDuration.String()}).Info("Starting turn shot clock")

	r.shotClock = time.AfterFunc(r.ShotClockDuration, func() {
		// Send a notification to self to handle the timeout safely within the actor loop.
		// @spec-link [[mech_game_state_versioning]]
		// @spec-link [[mech_action_economy_time_constraint_rules]]
		r.NotifyActor(message.Create(nil, rulermethods.Timeout{TurnIndex: turn}, nil))
	})
}

// timeout handles the turn expiration safely within the actor loop.
// @spec-link [[mech_game_state_versioning]]
func (r *Ruler) timeout(ctx actor.NotificationContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.Timeout)

	// Verify if the turn has changed since the timer was started (Race Prevention)
	if uint32(r.GameState.GetTurn()) != req.TurnIndex {
		r.logger.WithFields(logrus.Fields{
			"capturedTurn": req.TurnIndex,
			"currentTurn":  r.GameState.GetTurn()}).Debug("Shot clock expired but turn already progressed, ignoring.")
		return
	}

	r.logger.Warn("Turn timeout detected! Forcing EndOfTurn.")

	// Validate that the current entity still exists
	currentEntityID := r.GameState.Turner.CurrentEntityTurn
	if currentEntityID == uuid.Nil {
		return
	}

	ent, found := r.GameState.Entities[currentEntityID]
	if !found {
		return
	}

	// Trigger end of turn as a timeout
	r.endOfTurn(actor.CallContext{
		Msg: message.Create(nil, rulermethods.EndOfTurn{
			ControllerID: ent.ControllerID,
			EntityID:     currentEntityID,
			IsTimeout:    true,
			TurnIndex:    uint32(r.GameState.GetTurn()),
		}, nil),
	})
}

// @spec-link [[rule_turn_clock]]
func (r *Ruler) stopShotClock() {
	if r.shotClock != nil {
		r.shotClock.Stop()
		r.shotClock = nil
	}
}

// @spec-link [[mech_arena_lifecycle]]
func (r *Ruler) actorAboutToStop(ctx actor.NotificationContext) {
	r.logger.Info("Ruler is about to stop, stopping all controllers and timers")
	r.stopShotClock()
	for id, ctrl := range r.GameState.Controllers {
		r.logger.Infof("Stopping controller %s", id)
		ctrl.NotifyActor(message.Create(nil, actor.ActorStop{}, nil))
	}
}
