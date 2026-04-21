package controllers

import (
	"errors"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AggressiveController struct {
	*actor.Actor
	ID             uuid.UUID
	KnownEntities  map[uuid.UUID]entity.Entity
	Grid           *grid.Grid
	ruler          actor.Communication
	latestTarget   entity.Entity
	BattleFinished chan bool
	battleready    bool
}

func NewRdAggressiveController(name string) *AggressiveController {
	return NewAggressiveController(uuid.New(), name)
}

// New
func NewAggressiveController(id uuid.UUID, name string) *AggressiveController {
	ctrl := &AggressiveController{
		ID:             id,
		Actor:          actor.New(name),
		KnownEntities:  make(map[uuid.UUID]entity.Entity),
		BattleFinished: make(chan bool, 1),
		battleready:    false,
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
	ctrl.AddNotificationHandler(rulermethods.ControllerMoved{}, ctrl.NoOp, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerPassed{}, ctrl.NoOp, nil)

	ctrl.AddReplyHandler(rulermethods.GetStateReply{}, ctrl.GetStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.GetGridStateReply{}, ctrl.GetGridStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.GetEntitiesStateReply{}, ctrl.GetEntitiesStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.ControllerMoveReply{}, ctrl.ControllerMoveReply, nil)
	ctrl.AddReplyHandler(rulermethods.ControllerAttackReply{}, ctrl.ControllerAttackReply, nil)
	ctrl.AddReplyHandler(rulermethods.EndOfTurn{}, ctrl.EndOfTurnReply, nil)

	return ctrl
}

// implement actor.Manageable

func (ctl *AggressiveController) PrintStack() {
	ctl.GetQueue().PrintStack()
}

// implement all AggressiveController methods handlers.

// @spec-link [[mech_controller_handshake]]
func (ctl *AggressiveController) SetQueue(ctx actor.NotificationContext) {
	m := ctx.Msg.TargetMethod.(controllermethods.SetQueue)
	ctl.ruler = m.Ruler
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), ctl.GetCallbackChan())
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), ctl.GetCallbackChan())
}

func (ctl *AggressiveController) Send(ctx actor.NotificationContext) {

}

func (ctl *AggressiveController) ReceiveAPIMessage(ctx actor.NotificationContext) {

}

func (ctl *AggressiveController) ControllerNextTurn(ctx actor.NotificationContext) {
	controllerData := ctx.Msg.TargetMethod.(rulermethods.ControllerNextTurn)

	// GUARD: Since ControllerNextTurn is now broadcasted to all controllers, 
	// we must only take action if the entity belongs to US.
	if controllerData.Entity.ControllerID != ctl.ID {
		return
	}

	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn":     controllerData.Turn.String(),
		"EntityID": controllerData.Entity.String()}).Info("##### Turn BEGIN #####")
	time.Sleep(100 * time.Millisecond)
	target, err := ctl.selectNearestFoe(controllerData.Entity, ctl.KnownEntities)
	if err != nil {
		ctl.RequestLogger.Debug("Nothing to attack, ending turn")
		// No target ... Might have won the game!
		ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
			EntityID:     controllerData.Entity.ID,
			ControllerID: ctl.ID,
		}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())

		return
	}
	ctl.latestTarget = target
	ctl.RequestLogger.Debug("Moving To Attack")
	jumpHeight := ctl.KnownEntities[controllerData.Entity.ID].GetPropertyI(property.JumpHeight).I()

	path, found := ctl.preparePathToEntity(controllerData.Entity.Position, ctl.Grid, target, jumpHeight, controllerData.Entity.ID)
	// can't be on the same cell as target.
	if found && len(path) > 1 {
		movement := ctl.KnownEntities[controllerData.Entity.ID].GetProperty(property.Movement)
		mvt := movement.(*defaultproperty.DefaultIntCounterProperty).Value
		atkrng := ctl.KnownEntities[controllerData.Entity.ID].GetPropertyI(property.AttackRange).I()
		// 1. Determine how far we WANT to go to be in range
		// atkrng 1 means we stop 1 cell before the target.
		// AStarPath returns [p1, p2, ..., target], so len(path) includes target.
		limit := len(path) - atkrng
		if limit < 0 {
			limit = 0
		}

		// 2. Determine how far we CAN go based on movement
		if limit > mvt {
			limit = mvt
		}

		// 3. Determine how far we CAN go based on occupancy and obstacles
		// We must stop before the first blocked cell in the path.
		actualLimit := 0
		for i := 0; i < limit; i++ {
			// AStarPath now avoids entities, but we check again for extra safety or intermediate changes
			if ctl.isPathStepBlocked(path[i], controllerData.Entity.ID) {
				ctl.RequestLogger.WithFields(logrus.Fields{
					"blockedAt": path[i].String(),
					"step":      i,
				}).Debug("Path step blocked by entity, truncating movement")
				break
			}
			actualLimit = i + 1
		}

		if actualLimit > 1 {
			// Path includes start position at path[0]. actualLimit=1 means only start is reached.
			// actualLimit=2 means [start, next]. We move to next (path[1]).
			movePath := path[1:actualLimit]
			ctl.RequestLogger.WithFields(logrus.Fields{
				"actualLimit": actualLimit,
				"pathLen":     len(path),
				"movingTo":    movePath[len(movePath)-1].String(),
			}).Info("Moving toward target")

			ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
				ControllerID: ctl.ID,
				EntityID:     controllerData.Entity.ID,
				Path:         movePath,
			}, rulermethods.ControllerMoveReply{}), ctl.GetCallbackChan())
		} else {
			// Cannot move closer or already at range. 
			// If we are precisely at reach (len(path) <= atkrng + 1), attempt attack.
			if len(path) > 0 && len(path) <= atkrng+1 {
				ctl.RequestLogger.Info("Target in range, initiating attack")
				ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
					ControllerID: ctl.ID,
					EntityID:     controllerData.Entity.ID,
					Target:       target.Position,
				}, rulermethods.ControllerAttackReply{}), ctl.GetCallbackChan())
			} else {
				ctl.RequestLogger.WithFields(logrus.Fields{
					"atRange": len(path) <= atkrng+1,
					"path":    len(path),
				}).Warn("Target unreachable or blocked, passing turn")
				ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
					EntityID:     controllerData.Entity.ID,
					ControllerID: ctl.ID,
				}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
			}
		}
	} else {
		// it is already in place. Send attack

		if !found || len(path) == 0 {
			// Unable to find a path to target ...
			ctl.RequestLogger.WithFields(logrus.Fields{
				"EntityID": controllerData.Entity.ID.String()[0:8],
				"Position": controllerData.Entity.Position,
				"Expected": target.Position}).Debug("Unable to find path, ending turn")
			ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
				EntityID:     controllerData.Entity.ID,
				ControllerID: ctl.ID,
			}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
		} else {
			// right next to target.
			ctl.RequestLogger.WithFields(logrus.Fields{
				"EntityID":       controllerData.Entity.ID.String()[0:8],
				"Position":       controllerData.Entity.Position,
				"TargetPosition": target.Position,
				"TargetEntity":   target.ID.String()[0:8]}).Info("Attacking")
			ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
				EntityID:     controllerData.Entity.ID,
				Target:       target.Position,
				ControllerID: ctl.ID,
			}, rulermethods.ControllerAttackReply{
				Entity: controllerData.Entity,
			}), ctl.GetCallbackChan())
		}
	}
}

func (ctl *AggressiveController) BattleStart(ctx actor.NotificationContext) {
	ctl.RequestLogger.Info("##### BattleStart #####")
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), ctl.GetCallbackChan())
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), ctl.GetCallbackChan())
}

// @spec-link [[mech_ai_termination]]
func (ctl *AggressiveController) BattleEnd(ctx actor.NotificationContext) {
	ctl.RequestLogger.Info("##### BattleEnd #####")
	select {
	case ctl.BattleFinished <- true:
	default:
	}
}

func (ctl *AggressiveController) ControllerAttacked(ctx actor.NotificationContext) {
	attacked := ctx.Msg.TargetMethod.(rulermethods.ControllerAttacked)

	ctl.RequestLogger.WithFields(logrus.Fields{
		"EntityID":   attacked.Entity.ID.String()[0:8],
		"AttackerID": attacked.Attacker.ID.String()[0:8],
		"Position":   attacked.Entity.Position}).Debug("ControllerAttacked")
	// nothing to do post attack

	if attacked.Dead {
		// remove it from the known entities.
		delete(ctl.KnownEntities, attacked.Entity.ID)
	}

}

func (ctl *AggressiveController) EntitiesStateChanged(ctx actor.NotificationContext) {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn": ctx.Msg.TargetMethod.(rulermethods.EntitiesStateChanged).Turn.String()}).Info("New Turn Received")
	// fill in known entities (clear them beforehand.)
	ctl.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range ctx.Msg.TargetMethod.(rulermethods.EntitiesStateChanged).Entities {
		ctl.KnownEntities[e.ID] = e
	}
}

func (ctl *AggressiveController) GetStateReply(ctx actor.ReplyContext) {
}

func (ctl *AggressiveController) GetGridStateReply(ctx actor.ReplyContext) {
	ctl.Grid = ctx.Msg.Content.(rulermethods.GetGridStateReply).Grid
	if !ctl.battleready {
		ctl.battleready = true
		ctl.ruler.NotifyActor(message.Create(nil, rulermethods.ControllerBattleReady{
			ControllerID: ctl.ID,
		}, nil))
	}
}

func (ctl *AggressiveController) GetEntitiesStateReply(ctx actor.ReplyContext) {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn": ctx.Msg.Content.(rulermethods.GetEntitiesStateReply).Turn.String()}).Info("New Turn Info Received")
	// fill in known entities (clear them beforehand.)
	ctl.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range ctx.Msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		ctl.KnownEntities[e.ID] = e
	}
}

func (ctl *AggressiveController) ControllerMoveReply(ctx actor.ReplyContext) {
	if ctx.Msg.HasError {
		ctl.RequestLogger.WithFields(logrus.Fields{
			"error": ctx.Msg.ErrorMessage,
		}).Error("Move failed, ending turn")
		// If move fails, we don't have the updated entity state in the reply.
		// We should have it in KnownEntities though.
		// However, without a valid Entity in the reply, we might not know which entity was moving
		// if multiple were active (not the case for AggressiveController but good practice).
		return
	}

	ControllerData, ok := ctx.Msg.Content.(rulermethods.ControllerMoveReply)
	if !ok {
		ctl.RequestLogger.Error("Invalid MoveReply content type")
		return
	}
	target := ctl.latestTarget

	ctl.RequestLogger.WithFields(logrus.Fields{
		"EntityID": ControllerData.Entity.ID.String()[0:8],
		"Position": ControllerData.Entity.Position,
		"Expected": target.Position}).Debug("Move Succesfull")
	time.Sleep(100 * time.Millisecond)

	attacker := ctl.KnownEntities[ControllerData.Entity.ID]

	ctl.RequestLogger.Info(" Attacker: ", attacker.PrettyString())

	atkrng := attacker.GetPropertyI(property.AttackRange).I()
	if ControllerData.Entity.Position.Distance(target.Position) <= atkrng {

		// it is already in place. Send attack
		ctl.RequestLogger.WithFields(logrus.Fields{
			"EntityID":       ControllerData.Entity.ID.String()[0:8],
			"Position":       ControllerData.Entity.Position,
			"TargetPosition": target.Position,
			"AttackRange":    atkrng,
			"TargetEntity":   target.ID.String()[0:8]}).Info("Attacking")
		ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
			EntityID:     ControllerData.Entity.ID,
			Target:       target.Position,
			ControllerID: ctl.ID,
		}, rulermethods.ControllerAttackReply{
			Entity: ControllerData.Entity,
		}), ctl.GetCallbackChan())
	} else {
		// end of turn

		ctl.RequestLogger.WithFields(logrus.Fields{
			"EntityID": ControllerData.Entity.ID.String()[0:8],
			"Position": ControllerData.Entity.Position,
			"Expected": target.Position}).Debug("Too far away from target, ending turn")

		ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
			EntityID:     ControllerData.Entity.ID,
			ControllerID: ctl.ID,
		}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
	}
}

func (ctl *AggressiveController) ControllerAttackReply(ctx actor.ReplyContext) {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Error":   ctx.Msg.HasError,
		"Message": ctx.Msg.ErrorMessage,
	}).Info("Attack done, ending turn")
	
	if ctx.Msg.HasError {
		// Just end turn if attack failed
		return
	}

	reply, ok := ctx.Msg.Content.(rulermethods.ControllerAttackReply)
	if !ok {
		ctl.RequestLogger.Error("Invalid AttackReply content type")
		return
	}

	time.Sleep(100 * time.Millisecond)
	ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     reply.Entity.ID,
		ControllerID: ctl.ID,
	}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
}

// @spec-link [[rule_team_mechanics]]
// selectNearestFoe find nearest foe, based on team id.
// We rely on ctl.KnownEntities for dynamic status (life/death) and positions, 
// as ctl.Grid is a static snapshot of the map layout.
func (ctl *AggressiveController) selectNearestFoe(currentEntity entity.Entity, entities map[uuid.UUID]entity.Entity) (entity.Entity, error) {
	nearestid := uuid.Nil
	minDist := 10000
	pos := currentEntity.Position
	currentTeam := currentEntity.GetPropertyI(property.TeamID).I()

	ctl.RequestLogger.WithFields(logrus.Fields{
		"pos":     pos,
		"entity":  currentEntity.ID.String()[0:8],
		"team_id": currentTeam,
	}).Info("selectNearestFoe")

	for id, ent := range entities {
		hp := ent.GetPropertyI(property.HP).I()
		if ent.GetPropertyI(property.TeamID).I() != currentTeam && hp > 0 {
			if currentEntity.ID != ent.ID {
				ctl.RequestLogger.WithFields(logrus.Fields{
					"candidate_pos":    ent.Position,
					"candidate_entity": ent.ID.String()[0:8],
					"candidate_team":   ent.GetPropertyI(property.TeamID).I()}).Info("candidate")

				dist := tools.Distance(pos.X, pos.Y, ent.Position.X, ent.Position.Y)
				if nearestid == uuid.Nil || dist < minDist {
					ctl.RequestLogger.WithFields(logrus.Fields{
						"selected_pos":    ent.Position,
						"selected_entity": ent.ID.String()[0:8],
						"selected_team":   ent.GetPropertyI(property.TeamID).I()}).Info("selected")
					nearestid = id
					minDist = dist
				}
			}
		}
	}

	if nearestid != uuid.Nil {
		nearest := entities[nearestid]
		ctl.RequestLogger.WithFields(logrus.Fields{
			"nearest_pos":        nearest.Position,
			"nearest_entity":     nearest.ID.String()[0:8],
			"nearest_controller": nearest.ControllerID.String()[0:8]}).Info("nearest")
		return nearest, nil
	} else {
		return entity.Entity{}, errors.New("no nearest entity")
	}
}

// preparePathToEntity calculates a path using A* that avoids blocked cells.
// NOTE: ctl.Grid provides static map geometry, while the exclusion callback 
// utilizes ctl.KnownEntities to identify dynamic obstacles (other players).
func (ctl *AggressiveController) preparePathToEntity(pos position.Position, grd *grid.Grid, ent entity.Entity, jumpHeight int, selfID uuid.UUID) ([]position.Position, bool) {
	path, found := grd.AStarPath(pos, ent.Position, jumpHeight, func(p position.Position) bool {
		return ctl.isPathStepBlocked(p, selfID)
	})
	return path, found
}

func (ctl *AggressiveController) EndOfTurnReply(ctx actor.ReplyContext) {}

func (ctl *AggressiveController) NoOp(ctx actor.NotificationContext) {}

func (ctl *AggressiveController) isPathStepBlocked(pos position.Position, selfID uuid.UUID) bool {
	// 1. Check if occupied by another ALIVE entity in KnownEntities.
	// This is the primary source for dynamic occupancy since ctl.Grid is static.
	for _, ent := range ctl.KnownEntities {
		hp := ent.GetPropertyI(property.HP).I()
		if ent.ID != selfID && ent.Position.X == pos.X && ent.Position.Y == pos.Y && hp > 0 {
			return true
		}
	}

	// 2. Check Grid for obstacles or occupancy 
	if ctl.Grid != nil {
		cells := ctl.Grid.CellsForPositions([]position.Position{pos})
		if len(cells) > 0 {
			c := cells[0]
			if c.Type != cell.Ground {
				return true
			}

			
			if c.EntityID != uuid.Nil && c.EntityID != selfID {
				// check if entity is in KnownEntities; if not present, then it's a dead or ghost entity on a stale grid.
				if _, present := ctl.KnownEntities[c.EntityID]; present {
					return true
				}
			}
		}
	}

	return false
}
