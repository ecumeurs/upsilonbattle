package controllers

import (
	"errors"
	"reflect"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/position/pattern"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AggressiveController struct {
	act            *actor.Actor
	ID             uuid.UUID
	KnownEntities  map[uuid.UUID]entity.Entity
	Grid           *grid.Grid
	ruler          actor.Communication
	latestTarget   entity.Entity
	BattleFinished chan bool
}

// New
func NewAggressiveController(name string) *AggressiveController {
	ctrl := &AggressiveController{
		ID:             uuid.New(),
		KnownEntities:  make(map[uuid.UUID]entity.Entity),
		BattleFinished: make(chan bool),
	}

	ctrl.act = actor.New(name)
	ctrl.act.SetReceiveMessageHandler(ctrl.handleMessage)
	ctrl.act.SetReplyMessageHandler(ctrl.handleReply)

	return ctrl
}

// implement actor.Manageable

func (c *AggressiveController) Start() {
	c.act.Start()

}
func (c *AggressiveController) Stop() {
	c.act.Stop()
}
func (c *AggressiveController) PrepareToStop() chan bool {
	return c.act.PrepareToStop()
}

func (c *AggressiveController) handleMessage(msg message.Message) bool {
	logrus.WithFields(logrus.Fields{
		"Controller":   c.act.Name(),
		"message_type": reflect.TypeOf(msg.TargetMethod).String(),
		"controllerID": c.ID.String()[0:8]}).Info("Controller received message: ", reflect.TypeOf(msg.TargetMethod).String())

	if msg.HasError {
		logrus.WithFields(logrus.Fields{
			"error":      msg.ErrorMessage,
			"Controller": c.act.Name,
		}).Error("Error received")
	}

	switch msg.TargetMethod.(type) {
	case controllermethods.SetValidatorAndQueue:
		c.SetValidatorAndQueue(msg, msg.TargetMethod.(controllermethods.SetValidatorAndQueue))
		return true
	case controllermethods.Send:
		c.Send(msg)
		return true
	case controllermethods.ReceiveAPIMessage:
		c.ReceiveAPIMessage(msg)
		return true
	case rulermethods.ControllerNextTurn:
		c.ControllerNextTurn(msg)
		return true
	case rulermethods.BattleStart:
		c.BattleStart(msg)
		return true
	case rulermethods.BattleEnd:
		c.BattleEnd(msg)
		return true
	case rulermethods.EntitiesStateChanged:
		c.EntitiesStateChanged(msg)
		return true
	case rulermethods.ControllerAttacked:
		c.ControllerAttacked(msg)
		return true
	}
	return false
}

func (c *AggressiveController) handleReply(msg message.Message) bool {
	logrus.WithFields(logrus.Fields{
		"Controller": c.act.Name()}).Info("Controller received reply: ", reflect.TypeOf(msg.TargetMethod).String())

	if msg.HasError {
		logrus.WithFields(logrus.Fields{
			"Controller": c.act.Name(),
			"RequestId":  msg.RequestId.String()[0:8],
			"Error":      msg.ErrorMessage,
		}).Error("Error")
	}

	switch msg.TargetMethod.(type) {
	case rulermethods.GetStateReply:
		c.GetStateReply(msg)
		return true
	case rulermethods.GetGridStateReply:
		c.GetGridStateReply(msg)
		return true
	case rulermethods.GetEntitiesStateReply:
		c.GetEntitiesStateReply(msg)
		return true
	case rulermethods.ControllerMoveReply:
		c.ControllerMoveReply(msg)
		return true
	case rulermethods.ControllerAttackReply:
		c.ControllerAttackReply(msg)
		return true
	}
	return true // replies are always handled. we don't care about them.
}

// implement all AggressiveController methods handlers.

func (c *AggressiveController) SetValidatorAndQueue(msg message.Message, m controllermethods.SetValidatorAndQueue) {
	c.ruler = m.Ruler
	c.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), c.act.CallbackChan)
	c.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), c.act.CallbackChan)

	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) Send(msg message.Message) {
	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) ReceiveAPIMessage(msg message.Message) {
	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) ControllerNextTurn(msg message.Message) {
	c.nextTurn(msg) // will move or attack.
	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) BattleStart(msg message.Message) {
	logrus.Info("##### BattleStart #####")
	c.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), c.act.CallbackChan)
	c.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), c.act.CallbackChan)

	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) BattleEnd(msg message.Message) {
	logrus.Info("##### BattleEnd #####")
	c.BattleFinished <- true

	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) ControllerAttacked(msg message.Message) {
	logrus.WithFields(logrus.Fields{
		"EntityID":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.ID.String()[0:8],
		"AttackerID": msg.TargetMethod.(rulermethods.ControllerAttacked).Attacker.ID.String()[0:8],
		"Position":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.Position}).Debug("ControllerAttacked")
	// nothing to do post attack

	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) EntitiesStateChanged(msg message.Message) {
	logrus.WithFields(logrus.Fields{
		"Controller": c.act.Name(),
		"RequestId":  msg.RequestId.String()[0:8],
		"Turn":       msg.TargetMethod.(rulermethods.EntitiesStateChanged).Turn.String()}).Info("New Turn Received")
	// fill in known entities (clear them beforehand.)
	c.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range msg.TargetMethod.(rulermethods.EntitiesStateChanged).Entities {
		c.KnownEntities[e.ID] = e
	}
	c.act.NoReply(msg.Reply())
}

func (c *AggressiveController) GetStateReply(msg message.Message) {
}

func (c *AggressiveController) GetGridStateReply(msg message.Message) {
	c.Grid = msg.Content.(rulermethods.GetGridStateReply).Grid
}

func (c *AggressiveController) GetEntitiesStateReply(msg message.Message) {
	logrus.WithFields(logrus.Fields{
		"Controller": c.act.Name(),
		"RequestId":  msg.RequestId.String()[0:8],
		"Turn":       msg.Content.(rulermethods.GetEntitiesStateReply).Turn.String()}).Info("New Turn Received")
	// fill in known entities (clear them beforehand.)
	c.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		c.KnownEntities[e.ID] = e
	}
}

func (c *AggressiveController) ControllerMoveReply(msg message.Message) {
	c.afterMoveAttack(c.KnownEntities, msg)
}

func (c *AggressiveController) ControllerAttackReply(msg message.Message) {
	logrus.WithFields(logrus.Fields{
		"Controller":     c.ID.String()[0:8],
		"ControllerName": c.act.Name(),
		"Error":          msg.HasError,
		"Message":        msg.ErrorMessage,
	}).Info("Attack done, ending turn")
	c.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     msg.TargetMethod.(rulermethods.ControllerAttackReply).Entity.ID,
		ControllerID: c.ID,
	}, rulermethods.EndOfTurn{}), c.act.CallbackChan)
}

// selectNearestFoe find nearest foe, based on controller id.
func (ctl *AggressiveController) selectNearestFoe(currentEntity entity.Entity, entities map[uuid.UUID]entity.Entity) (entity.Entity, error) {
	nearestid := uuid.Nil
	minDist := 10000
	pos := currentEntity.Position
	currentCtrl := currentEntity.ControllerID

	logrus.WithFields(logrus.Fields{
		"pos":        pos,
		"entity":     currentEntity.ID.String()[0:8],
		"controller": currentEntity.ControllerID.String()[0:8],
	}).Info("selectNearestFoe")

	for id, ent := range entities {
		if ent.ControllerID != currentCtrl {
			if currentEntity.ID != ent.ID {
				logrus.WithFields(logrus.Fields{
					"candidate_pos":        ent.Position,
					"candidate_entity":     ent.ID.String()[0:8],
					"candidate_controller": ent.ControllerID.String()[0:8]}).Info("candidate")

				dist := tools.Distance(pos.X, pos.Y, ent.Position.X, ent.Position.Y)
				if nearestid == uuid.Nil || dist < minDist {
					logrus.WithFields(logrus.Fields{
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
		logrus.WithFields(logrus.Fields{
			"nearest_pos":        nearest.Position,
			"nearest_entity":     nearest.ID.String()[0:8],
			"nearest_controller": nearest.ControllerID.String()[0:8]}).Info("nearest")
		return nearest, nil
	} else {
		return entity.Entity{}, errors.New("no nearest entity")
	}
}

func (ctl *AggressiveController) preparePathToEntity(pos position.Position, grd *grid.Grid, ent entity.Entity) []position.Position {
	path := pattern.PathTo2D(ent.Position.Substract(pos)).Apply2D(pos)
	// ensure Z is at the right place all along.
	for i, p := range path {
		path[i].Z = grd.TopMostGroundAt(p.X, p.Y)
	}

	return path
}

func (ctl *AggressiveController) nextTurn(msg message.Message) {
	controllerData := msg.TargetMethod.(rulermethods.ControllerNextTurn)
	logrus.WithFields(logrus.Fields{
		"Controller":     ctl.ID.String()[0:8],
		"ControllerName": ctl.act.Name(),
		"RequestID":      msg.RequestId.String()[0:8],
		"Turn":           controllerData.Turn.String(),
		"EntityID":       controllerData.Entity.ID.String()[0:8]}).Info("##### Turn BEGIN #####")
	target, err := ctl.selectNearestFoe(controllerData.Entity, ctl.KnownEntities)
	if err != nil {
		logrus.Debug("Nothing to attack, ending turn")
		// No target ... Might have won the game!
		ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
			EntityID:     controllerData.Entity.ID,
			ControllerID: ctl.ID,
		}, rulermethods.EndOfTurn{}), ctl.act.CallbackChan)

		return
	}
	ctl.latestTarget = target
	logrus.Debug("Moving To Attack")

	path := ctl.preparePathToEntity(controllerData.Entity.Position, ctl.Grid, target)
	// can't be on the same cell as target.
	if len(path) > 1 {
		path = path[:len(path)-1]

		if len(path) != 0 {
			logrus.WithFields(logrus.Fields{
				"EntityID":       controllerData.Entity.ID.String()[0:8],
				"Position":       controllerData.Entity.Position,
				"Expected":       path[len(path)-1],
				"TargetPosition": target.Position,
				"TargetEntity":   target.ID.String()[0:8],
			}).Info("Moving attacker")

			ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
				EntityID:     controllerData.Entity.ID,
				Path:         path,
				ControllerID: ctl.ID,
			}, rulermethods.ControllerMoveReply{
				Entity: controllerData.Entity,
			}), ctl.act.CallbackChan)
		} else {
			// it is already in place. Send attack
			logrus.WithFields(logrus.Fields{
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
			}), ctl.act.CallbackChan)
		}
	} else {
		// it is already in place. Send attack
		logrus.WithFields(logrus.Fields{
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
		}), ctl.act.CallbackChan)
	}
}

func (ctl *AggressiveController) afterMoveAttack(knownEntities map[uuid.UUID]entity.Entity, msg message.Message) {
	if msg.HasError {
		ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
			EntityID:     msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID,
			ControllerID: ctl.ID,
		}, rulermethods.EndOfTurn{}), ctl.act.CallbackChan)
	} else {
		ControllerData := msg.Content.(rulermethods.ControllerMoveReply)
		target := ctl.latestTarget

		logrus.WithFields(logrus.Fields{
			"EntityID":       ControllerData.Entity.ID.String()[0:8],
			"Position":       ControllerData.Entity.Position,
			"Controller":     ctl.ID.String()[0:8],
			"ControllerName": ctl.act.Name(),
			"Expected":       target.Position}).Debug("Move Succesfull")

		// it is already in place. Send attack
		logrus.WithFields(logrus.Fields{
			"EntityID":       ControllerData.Entity.ID.String()[0:8],
			"Position":       ControllerData.Entity.Position,
			"TargetPosition": target.Position,
			"TargetEntity":   target.ID.String()[0:8]}).Info("Attacking")
		ctl.ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
			EntityID:     ControllerData.Entity.ID,
			Target:       target.Position,
			ControllerID: ctl.ID,
		}, rulermethods.ControllerAttackReply{
			Entity: ControllerData.Entity,
		}), ctl.act.CallbackChan)
	}
}

// Implement the actor.Communication interface
// Notify Actor
func (c *AggressiveController) NotifyActor(msg message.Message) {
	c.act.Notify(msg)
}

func (c *AggressiveController) SendActor(msg message.Message, callback chan message.Message) {
	c.act.Send(msg, callback)
}
