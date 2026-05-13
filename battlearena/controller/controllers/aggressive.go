package controllers

import (
	"fmt"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// AIController is an actor-based controller that runs an entity's turn through a
// LayeredBehavior pipeline, emitting one EngineCommand per tick.
//
// AggressiveController is retained as an alias for backward compatibility.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type AIController struct {
	*actor.Actor
	ID             uuid.UUID
	KnownEntities  map[uuid.UUID]entity.Entity
	Grid           *grid.Grid
	ruler          actor.Communication
	BattleFinished chan bool
	battleready    bool

	pipeline *behavior.LayeredBehavior
	memory   *behavior.DecisionMemory
	gradeIdx int

	// per-turn transient state (reset each ControllerNextTurn)
	currentEntityID uuid.UUID
	turnCtx         *controllerGameContext
	lastSentPath    []position.Position
}

// AggressiveController is kept as a type alias so existing call sites compile unchanged.
type AggressiveController = AIController

// NewAggressiveController returns an AIController with the baseline-only stack.
// Use NewAIController to supply an archetype-specific pipeline.
func NewAggressiveController(id uuid.UUID, name string) *AIController {
	return NewAIController(id, name, behavior.NewLayeredBehavior(&behavior.AggressiveBehavior{}))
}

// NewRdAggressiveController creates a baseline AIController with a random UUID.
func NewRdAggressiveController(name string) *AIController {
	return NewAggressiveController(uuid.New(), name)
}

// NewAIController creates an AIController driven by the supplied pipeline.
func NewAIController(id uuid.UUID, name string, pipeline *behavior.LayeredBehavior) *AIController {
	ctrl := &AIController{
		ID:             id,
		Actor:          actor.New(name),
		KnownEntities:  make(map[uuid.UUID]entity.Entity),
		BattleFinished: make(chan bool, 1),
		battleready:    false,
		pipeline:       pipeline,
		memory:         behavior.NewDecisionMemory(),
		gradeIdx:       0,
	}
	ctrl.Logger = ctrl.Logger.WithFields(logrus.Fields{
		"ControllerID": id.String()[0:8],
		"Controller":   name,
		"Component":    "controller",
	})

	ctrl.AddNotificationHandler(controllermethods.SetQueue{}, ctrl.SetQueue, nil)
	ctrl.AddNotificationHandler(controllermethods.Send{}, ctrl.Send, nil)
	ctrl.AddNotificationHandler(controllermethods.ReceiveAPIMessage{}, ctrl.ReceiveAPIMessage, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerNextTurn{}, ctrl.ControllerNextTurn, nil)
	ctrl.AddNotificationHandler(rulermethods.BattleStart{}, ctrl.BattleStart, nil)
	ctrl.AddNotificationHandler(rulermethods.BattleEnd{}, ctrl.BattleEnd, nil)
	ctrl.AddNotificationHandler(rulermethods.EntitiesStateChanged{}, ctrl.EntitiesStateChanged, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerAttacked{}, ctrl.ControllerAttacked, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerSkillUsed{}, ctrl.NoOp, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerMoved{}, ctrl.NoOp, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerPassed{}, ctrl.NoOp, nil)

	ctrl.AddReplyHandler(rulermethods.GetStateReply{}, ctrl.GetStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.GetGridStateReply{}, ctrl.GetGridStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.GetEntitiesStateReply{}, ctrl.GetEntitiesStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.ControllerMoveReply{}, ctrl.ControllerMoveReply, nil)
	ctrl.AddReplyHandler(rulermethods.ControllerAttackReply{}, ctrl.ControllerAttackReply, nil)
	ctrl.AddReplyHandler(rulermethods.ControllerUseSkillReply{}, ctrl.ControllerUseSkillReply, nil)
	ctrl.AddReplyHandler(rulermethods.EndOfTurn{}, ctrl.EndOfTurnReply, nil)

	return ctrl
}

// PrintStack implements actor.Manageable.
func (ctl *AIController) PrintStack() {
	ctl.GetQueue().PrintStack()
}

// ── Notification handlers ─────────────────────────────────────────────────────

// @spec-link [[mech_controller_handshake]]
func (ctl *AIController) SetQueue(ctx actor.NotificationContext) {
	m := ctx.Msg.TargetMethod.(controllermethods.SetQueue)
	ctl.ruler = m.Ruler
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), ctl.GetCallbackChan())
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), ctl.GetCallbackChan())
}

// Send is a no-op — AI controllers do not send API messages.
func (ctl *AIController) Send(ctx actor.NotificationContext) {
}

// ReceiveAPIMessage is a no-op — AI controllers do not process external API messages.
func (ctl *AIController) ReceiveAPIMessage(ctx actor.NotificationContext) {
}

// ControllerNextTurn starts a new entity turn. It initialises the per-turn context
// and runs the first pipeline tick, dispatching the resulting EngineCommand.
func (ctl *AIController) ControllerNextTurn(ctx actor.NotificationContext) {
	controllerData := ctx.Msg.TargetMethod.(rulermethods.ControllerNextTurn)

	// GUARD: broadcasted to all controllers — only handle our own entity.
	if controllerData.Entity.ControllerID != ctl.ID {
		return
	}

	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn":     controllerData.Turn.String(),
		"EntityID": controllerData.Entity.String(),
	}).Info("##### Turn BEGIN #####")
	time.Sleep(100 * time.Millisecond)

	ent := ctl.KnownEntities[controllerData.Entity.ID]
	ctl.memory.AdvanceTurn()
	ctl.currentEntityID = ent.ID
	ctl.turnCtx = &controllerGameContext{
		self:         ent,
		entities:     ctl.KnownEntities,
		grd:          ctl.Grid,
		hasActed:     false,
		remainingMvt: ent.GetPropertyC(property.Movement).GetValue(),
		lastOutcome:  behavior.TickNone,
		mem:          ctl.memory,
		gradeIdx:     ctl.gradeIdx,
	}

	ctl.dispatchTick()
}

// BattleStart fetches the initial grid and entity snapshots so the AI is ready for its first turn.
func (ctl *AIController) BattleStart(ctx actor.NotificationContext) {
	ctl.RequestLogger.Info("##### BattleStart #####")
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), ctl.GetCallbackChan())
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), ctl.GetCallbackChan())
}

// BattleEnd signals the battle is over and unblocks any caller waiting on BattleFinished.
//
// @spec-link [[mechanic_mech_ai_termination]]
func (ctl *AIController) BattleEnd(ctx actor.NotificationContext) {
	ctl.RequestLogger.Info("##### BattleEnd #####")
	select {
	case ctl.BattleFinished <- true:
	default:
	}
}

// ControllerAttacked removes dead entities from KnownEntities.
func (ctl *AIController) ControllerAttacked(ctx actor.NotificationContext) {
	attacked := ctx.Msg.TargetMethod.(rulermethods.ControllerAttacked)
	ctl.RequestLogger.WithFields(logrus.Fields{
		"EntityID":   attacked.Entity.ID.String()[0:8],
		"AttackerID": attacked.Attacker.ID.String()[0:8],
		"Position":   attacked.Entity.Position,
	}).Debug("ControllerAttacked")
	if attacked.Dead {
		delete(ctl.KnownEntities, attacked.Entity.ID)
	}
}

// EntitiesStateChanged refreshes the full entity snapshot.
func (ctl *AIController) EntitiesStateChanged(ctx actor.NotificationContext) {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn": ctx.Msg.TargetMethod.(rulermethods.EntitiesStateChanged).Turn.String(),
	}).Info("New Turn Received")
	ctl.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range ctx.Msg.TargetMethod.(rulermethods.EntitiesStateChanged).Entities {
		ctl.KnownEntities[e.ID] = e
	}
}

// NoOp is an empty handler for events this controller doesn't need to act on.
func (ctl *AIController) NoOp(ctx actor.NotificationContext) {
}

// ── Reply handlers ────────────────────────────────────────────────────────────

// GetStateReply is a no-op — state is fetched through the entities/grid-specific replies.
func (ctl *AIController) GetStateReply(ctx actor.ReplyContext) {
}

// GetGridStateReply stores the received grid and marks the controller as battle-ready.
//
// @spec-link [[mech_controller_handshake]]
func (ctl *AIController) GetGridStateReply(ctx actor.ReplyContext) {
	ctl.Grid = ctx.Msg.Content.(rulermethods.GetGridStateReply).Grid
	if !ctl.battleready {
		ctl.battleready = true
		ctl.ruler.NotifyActor(message.Create(nil, rulermethods.ControllerBattleReady{
			ControllerID: ctl.ID,
		}, nil))
	}
}

// GetEntitiesStateReply refreshes the full entity snapshot from the ruler.
func (ctl *AIController) GetEntitiesStateReply(ctx actor.ReplyContext) {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn": ctx.Msg.Content.(rulermethods.GetEntitiesStateReply).Turn.String(),
	}).Info("New Turn Info Received")
	ctl.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range ctx.Msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		ctl.KnownEntities[e.ID] = e
	}
}

// ControllerMoveReply updates turn state after a move and runs the next pipeline tick.
func (ctl *AIController) ControllerMoveReply(ctx actor.ReplyContext) {
	if ctl.turnCtx == nil {
		return
	}
	if ctx.Msg.HasError {
		ctl.RequestLogger.WithFields(logrus.Fields{
			"error": ctx.Msg.ErrorMessage,
		}).Error("Move failed, ending turn")
		ctl.endTurn()
		return
	}
	reply, ok := ctx.Msg.Content.(rulermethods.ControllerMoveReply)
	if !ok {
		ctl.RequestLogger.Error("Invalid MoveReply content type")
		ctl.endTurn()
		return
	}

	// Update entity position in KnownEntities and turn context.
	ctl.KnownEntities[reply.Entity.ID] = reply.Entity
	ctl.turnCtx.self = reply.Entity
	ctl.turnCtx.entities = ctl.KnownEntities
	ctl.turnCtx.remainingMvt -= len(ctl.lastSentPath)
	if ctl.turnCtx.remainingMvt < 0 {
		ctl.turnCtx.remainingMvt = 0
	}
	ctl.turnCtx.lastOutcome = behavior.TickSuccess

	ctl.RequestLogger.WithFields(logrus.Fields{
		"EntityID":     reply.Entity.ID.String()[0:8],
		"Position":     reply.Entity.Position,
		"RemainingMvt": ctl.turnCtx.remainingMvt,
	}).Debug("Move successful")

	time.Sleep(100 * time.Millisecond)
	ctl.dispatchTick()
}

// ControllerAttackReply updates turn state after an attack and runs the next pipeline tick.
func (ctl *AIController) ControllerAttackReply(ctx actor.ReplyContext) {
	if ctl.turnCtx == nil {
		return
	}
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Error":   ctx.Msg.HasError,
		"Message": ctx.Msg.ErrorMessage,
	}).Info("Attack done")

	if ctx.Msg.HasError {
		ctl.endTurn()
		return
	}
	reply, ok := ctx.Msg.Content.(rulermethods.ControllerAttackReply)
	if !ok {
		ctl.RequestLogger.Error("Invalid AttackReply content type")
		ctl.endTurn()
		return
	}
	ctl.KnownEntities[reply.Attacker.ID] = reply.Attacker
	ctl.turnCtx.self = reply.Attacker
	ctl.turnCtx.entities = ctl.KnownEntities
	ctl.turnCtx.hasActed = true
	ctl.turnCtx.lastOutcome = behavior.TickSuccess

	time.Sleep(100 * time.Millisecond)
	ctl.dispatchTick()
}

// ControllerUseSkillReply updates turn state after a skill and runs the next pipeline tick.
func (ctl *AIController) ControllerUseSkillReply(ctx actor.ReplyContext) {
	if ctl.turnCtx == nil {
		return
	}
	ctl.RequestLogger.Info("Skill used")

	if ctx.Msg.HasError {
		ctl.endTurn()
		return
	}
	reply, ok := ctx.Msg.Content.(rulermethods.ControllerUseSkillReply)
	if !ok {
		ctl.RequestLogger.Error("Invalid UseSkillReply content type")
		ctl.endTurn()
		return
	}
	ctl.KnownEntities[reply.Attacker.ID] = reply.Attacker
	ctl.turnCtx.self = reply.Attacker
	ctl.turnCtx.entities = ctl.KnownEntities
	ctl.turnCtx.hasActed = true
	ctl.turnCtx.lastOutcome = behavior.TickSuccess

	time.Sleep(100 * time.Millisecond)
	ctl.dispatchTick()
}

// EndOfTurnReply acknowledges the ruler's confirmation that the turn has ended.
func (ctl *AIController) EndOfTurnReply(ctx actor.ReplyContext) {
}

// ── Internal helpers ──────────────────────────────────────────────────────────

// dispatchTick runs one pipeline tick and sends the resulting command to the ruler.
func (ctl *AIController) dispatchTick() {
	cmd := ctl.pipeline.Tick(ctl.turnCtx)
	entityID := ctl.currentEntityID

	switch cmd.Type {
	case behavior.CmdMove:
		ctl.lastSentPath = cmd.Path
		ctl.RequestLogger.WithFields(logrus.Fields{
			"EntityID": entityID.String()[0:8],
			"PathLen":  len(cmd.Path),
			"MovingTo": cmd.Path[len(cmd.Path)-1].String(),
		}).Info("Moving toward target")
		ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
			ControllerID: ctl.ID,
			EntityID:     entityID,
			Path:         cmd.Path,
		}, rulermethods.ControllerMoveReply{}), ctl.GetCallbackChan())

	case behavior.CmdAttack:
		ctl.RequestLogger.WithFields(logrus.Fields{
			"EntityID": entityID.String()[0:8],
			"Target":   cmd.Target.String(),
		}).Info("Attacking")
		ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
			ControllerID: ctl.ID,
			EntityID:     entityID,
			Target:       cmd.Target,
		}, rulermethods.ControllerAttackReply{}), ctl.GetCallbackChan())

	case behavior.CmdSkill:
		ctl.RequestLogger.WithFields(logrus.Fields{
			"EntityID": entityID.String()[0:8],
			"SkillID":  cmd.SkillID.String()[0:8],
			"Target":   cmd.Target.String(),
		}).Info("Using skill")
		ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerUseSkill{
			ControllerID: ctl.ID,
			EntityID:     entityID,
			SkillID:      cmd.SkillID,
			Target:       cmd.Target,
		}, rulermethods.ControllerUseSkillReply{}), ctl.GetCallbackChan())

	default: // CmdEndOfTurn
		ctl.endTurn()
	}
}

// endTurn sends EndOfTurn to the ruler and clears the turn context.
func (ctl *AIController) endTurn() {
	ctl.RequestLogger.Debug("Ending turn")
	ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     ctl.currentEntityID,
		ControllerID: ctl.ID,
	}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
	ctl.turnCtx = nil
	ctl.lastSentPath = nil
}

// SetGrade configures the AI grade index (0 = Grade I, 8 = Grade V) that scales
// each behavior layer's activation rate.
func (ctl *AIController) SetGrade(gradeIdx int) {
	ctl.gradeIdx = gradeIdx
}

// selectNearestFoe returns the living enemy closest to currentEntity by Manhattan distance.
// @spec-link [[rule_team_mechanics]]
func (ctl *AIController) selectNearestFoe(currentEntity entity.Entity, entities map[uuid.UUID]entity.Entity) (entity.Entity, error) {
	var nearest entity.Entity
	found := false
	bestDist := int(^uint(0) >> 1)
	currentTeam := currentEntity.GetPropertyI(property.TeamID).I()

	for _, ent := range entities {
		if ent.ID == currentEntity.ID {
			continue
		}
		if ent.GetPropertyI(property.TeamID).I() == currentTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		d := tools.Distance(currentEntity.Position.X, currentEntity.Position.Y, ent.Position.X, ent.Position.Y)
		if !found || d < bestDist {
			bestDist = d
			nearest = ent
			found = true
		}
	}
	if !found {
		return entity.Entity{}, fmt.Errorf("no living enemy found")
	}
	return nearest, nil
}
