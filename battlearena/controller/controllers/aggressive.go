package controllers

import (
	"errors"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
	"github.com/ecumeurs/upsilonmapdata/grid"
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

// New
func NewAggressiveController(name string) *AggressiveController {
	id := uuid.New()
	ctrl := &AggressiveController{
		ID:             id,
		Actor:          actor.New(name),
		KnownEntities:  make(map[uuid.UUID]entity.Entity),
		BattleFinished: make(chan bool),
		battleready:    false,
	}
	ctrl.Logger = ctrl.Logger.WithFields(logrus.Fields{
		"ControllerID": id.String()[0:8],
		"Controller":   name,
		"Component":    "controller",
	})

	ctrl.AddMethod(controllermethods.SetQueue{}, ctrl.SetQueue, nil)
	ctrl.AddMethod(controllermethods.Send{}, ctrl.Send, nil)
	ctrl.AddMethod(controllermethods.ReceiveAPIMessage{}, ctrl.ReceiveAPIMessage, nil)
	ctrl.AddMethod(rulermethods.ControllerNextTurn{}, ctrl.ControllerNextTurn, nil)
	ctrl.AddMethod(rulermethods.BattleStart{}, ctrl.BattleStart, nil)
	ctrl.AddMethod(rulermethods.BattleEnd{}, ctrl.BattleEnd, nil)
	ctrl.AddMethod(rulermethods.EntitiesStateChanged{}, ctrl.EntitiesStateChanged, nil)
	ctrl.AddMethod(rulermethods.ControllerAttacked{}, ctrl.ControllerAttacked, nil)

	ctrl.AddReply(rulermethods.GetState{}, ctrl.GetStateReply, nil)
	ctrl.AddReply(rulermethods.GetGridState{}, ctrl.GetGridStateReply, nil)
	ctrl.AddReply(rulermethods.GetEntitiesState{}, ctrl.GetEntitiesStateReply, nil)
	ctrl.AddReply(rulermethods.ControllerMove{}, ctrl.ControllerMoveReply, nil)
	ctrl.AddReply(rulermethods.ControllerAttack{}, ctrl.ControllerAttackReply, nil)

	return ctrl
}

// implement actor.Manageable

func (ctl *AggressiveController) PrintStack() {
	ctl.GetQueue().PrintStack()
}

// implement all AggressiveController methods handlers.

func (ctl *AggressiveController) SetQueue(msg *message.Message) bool {
	m := msg.TargetMethod.(controllermethods.SetQueue)
	ctl.ruler = m.Ruler
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), ctl.GetCallbackChan())
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), ctl.GetCallbackChan())
	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) Send(msg *message.Message) bool {

	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) ReceiveAPIMessage(msg *message.Message) bool {

	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) ControllerNextTurn(msg *message.Message) bool {
	controllerData := msg.TargetMethod.(rulermethods.ControllerNextTurn)
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn":     controllerData.Turn.String(),
		"EntityID": controllerData.Entity.String()}).Info("##### Turn BEGIN #####")
	target, err := ctl.selectNearestFoe(controllerData.Entity, ctl.KnownEntities)
	if err != nil {
		ctl.RequestLogger.Debug("Nothing to attack, ending turn")
		// No target ... Might have won the game!
		ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
			EntityID:     controllerData.Entity.ID,
			ControllerID: ctl.ID,
		}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())

		return true
	}
	ctl.latestTarget = target
	ctl.RequestLogger.Debug("Moving To Attack")
	jumpHeight := ctl.KnownEntities[controllerData.Entity.ID].GetPropertyI(property.JumpHeight).I()

	path := ctl.preparePathToEntity(controllerData.Entity.Position, ctl.Grid, target, jumpHeight)
	// can't be on the same cell as target.
	if len(path) > 1 {
		movement := ctl.KnownEntities[controllerData.Entity.ID].GetProperty(property.Movement)
		mvt := movement.(*defaultproperty.DefaultIntCounterProperty).Value
		atkrng := ctl.KnownEntities[controllerData.Entity.ID].GetPropertyI(property.AttackRange).I()
		if len(path) > atkrng {
			path = path[:atkrng]
		} else {
			path = nil // no need to move, already in range.
		}

		// limit path according to entity's movement
		if len(path) > mvt {
			path = path[:mvt]
		}

		if len(path) != 0 {

			ctl.RequestLogger.WithFields(logrus.Fields{
				"EntityID":       controllerData.Entity.ID.String()[0:8],
				"Position":       controllerData.Entity.Position,
				"Expected":       path[len(path)-1],
				"Movement":       mvt,
				"AttackRange":    atkrng,
				"TargetPosition": target.Position,
				"TargetEntity":   target.ID.String()[0:8],
			}).Info("Moving attacker")

			ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
				EntityID:     controllerData.Entity.ID,
				Path:         path,
				ControllerID: ctl.ID,
			}, rulermethods.ControllerMoveReply{
				Entity: controllerData.Entity,
			}), ctl.GetCallbackChan())
		} else {
			// it is already in place. Send attack
			ctl.RequestLogger.WithFields(logrus.Fields{
				"EntityID":       controllerData.Entity.ID.String()[0:8],
				"Position":       controllerData.Entity.Position,
				"TargetPosition": target.Position,
				"Movement":       mvt,
				"AttackRange":    atkrng,
				"TargetEntity":   target.ID.String()[0:8]}).Info("Attacking")
			ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
				EntityID:     controllerData.Entity.ID,
				Target:       target.Position,
				ControllerID: ctl.ID,
			}, rulermethods.ControllerAttackReply{
				Entity: controllerData.Entity,
			}), ctl.GetCallbackChan())
		}
	} else {
		// it is already in place. Send attack

		if len(path) == 0 {
			// Unable to find a path to target ...
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
	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) BattleStart(msg *message.Message) bool {
	ctl.RequestLogger.Info("##### BattleStart #####")
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), ctl.GetCallbackChan())
	ctl.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), ctl.GetCallbackChan())
	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) BattleEnd(msg *message.Message) bool {
	ctl.RequestLogger.Info("##### BattleEnd #####")
	ctl.BattleFinished <- true
	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) ControllerAttacked(msg *message.Message) bool {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"EntityID":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.ID.String()[0:8],
		"AttackerID": msg.TargetMethod.(rulermethods.ControllerAttacked).Attacker.ID.String()[0:8],
		"Position":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.Position}).Debug("ControllerAttacked")
	// nothing to do post attack

	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) EntitiesStateChanged(msg *message.Message) bool {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn": msg.TargetMethod.(rulermethods.EntitiesStateChanged).Turn.String()}).Info("New Turn Received")
	// fill in known entities (clear them beforehand.)
	ctl.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range msg.TargetMethod.(rulermethods.EntitiesStateChanged).Entities {
		ctl.KnownEntities[e.ID] = e
	}
	ctl.NoReply(msg)
	return true
}

func (ctl *AggressiveController) GetStateReply(msg *message.Message) bool {
	return true
}

func (ctl *AggressiveController) GetGridStateReply(msg *message.Message) bool {
	ctl.Grid = msg.Content.(rulermethods.GetGridStateReply).Grid
	if !ctl.battleready {
		ctl.battleready = true
		ctl.ruler.NotifyActor(message.Create(nil, rulermethods.ControllerBattleReady{
			ControllerID: ctl.ID,
		}, nil))
	}
	return true
}

func (ctl *AggressiveController) GetEntitiesStateReply(msg *message.Message) bool {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Turn": msg.Content.(rulermethods.GetEntitiesStateReply).Turn.String()}).Info("New Turn Info Received")
	// fill in known entities (clear them beforehand.)
	ctl.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		ctl.KnownEntities[e.ID] = e
	}
	return true
}

func (ctl *AggressiveController) ControllerMoveReply(msg *message.Message) bool {
	if msg.HasError {
		ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
			EntityID:     msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID,
			ControllerID: ctl.ID,
		}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
	} else {
		ControllerData := msg.Content.(rulermethods.ControllerMoveReply)
		target := ctl.latestTarget

		ctl.RequestLogger.WithFields(logrus.Fields{
			"EntityID": ControllerData.Entity.ID.String()[0:8],
			"Position": ControllerData.Entity.Position,
			"Expected": target.Position}).Debug("Move Succesfull")

		attacker := ctl.KnownEntities[msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID]

		ctl.RequestLogger.Info(" Attacker: ", attacker.PrettyString())

		atkrng := attacker.GetPropertyI(property.AttackRange).I()
		if msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.Position.Distance(target.Position) <= atkrng {

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
				EntityID:     msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID,
				ControllerID: ctl.ID,
			}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
		}
	}

	return true
}

func (ctl *AggressiveController) ControllerAttackReply(msg *message.Message) bool {
	ctl.RequestLogger.WithFields(logrus.Fields{
		"Error":   msg.HasError,
		"Message": msg.ErrorMessage,
	}).Info("Attack done, ending turn")
	ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     msg.TargetMethod.(rulermethods.ControllerAttackReply).Entity.ID,
		ControllerID: ctl.ID,
	}, rulermethods.EndOfTurn{}), ctl.GetCallbackChan())
	return true
}

// selectNearestFoe find nearest foe, based on controller id.
func (ctl *AggressiveController) selectNearestFoe(currentEntity entity.Entity, entities map[uuid.UUID]entity.Entity) (entity.Entity, error) {
	nearestid := uuid.Nil
	minDist := 10000
	pos := currentEntity.Position
	currentCtrl := currentEntity.ControllerID

	ctl.RequestLogger.WithFields(logrus.Fields{
		"pos":        pos,
		"entity":     currentEntity.ID.String()[0:8],
		"controller": currentEntity.ControllerID.String()[0:8],
	}).Info("selectNearestFoe")

	for id, ent := range entities {
		if ent.ControllerID != currentCtrl {
			if currentEntity.ID != ent.ID {
				ctl.RequestLogger.WithFields(logrus.Fields{
					"candidate_pos":        ent.Position,
					"candidate_entity":     ent.ID.String()[0:8],
					"candidate_controller": ent.ControllerID.String()[0:8]}).Info("candidate")

				dist := tools.Distance(pos.X, pos.Y, ent.Position.X, ent.Position.Y)
				if nearestid == uuid.Nil || dist < minDist {
					ctl.RequestLogger.WithFields(logrus.Fields{
						"selected_pos":        ent.Position,
						"selected_entity":     ent.ID.String()[0:8],
						"selected_controller": ent.ControllerID.String()[0:8]}).Info("selected")
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

func (ctl *AggressiveController) preparePathToEntity(pos position.Position, grd *grid.Grid, ent entity.Entity, jumpHeight int) []position.Position {
	path, found := grd.AStarPath(pos, ent.Position, jumpHeight)
	if !found {
		return nil
	}

	return path
}
