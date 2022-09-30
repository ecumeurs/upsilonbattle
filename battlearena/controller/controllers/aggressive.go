package controllers

import (
	"errors"
	"reflect"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"

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
	battleready    bool
}

// New
func NewAggressiveController(name string) *AggressiveController {
	ctrl := &AggressiveController{
		ID:             uuid.New(),
		KnownEntities:  make(map[uuid.UUID]entity.Entity),
		BattleFinished: make(chan bool),
		battleready:    false,
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

func (c *AggressiveController) PrintStack() {
	c.act.GetQueue().PrintStack()
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
	case controllermethods.SetQueue:
		c.SetQueue(msg, msg.TargetMethod.(controllermethods.SetQueue))
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

func (c *AggressiveController) SetQueue(msg message.Message, m controllermethods.SetQueue) {
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

func (ctl *AggressiveController) ControllerNextTurn(msg message.Message) {
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

			logrus.WithFields(logrus.Fields{
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
			}), ctl.act.CallbackChan)
		} else {
			// it is already in place. Send attack
			logrus.WithFields(logrus.Fields{
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
			}), ctl.act.CallbackChan)
		}
	} else {
		// it is already in place. Send attack

		if len(path) == 0 {
			// Unable to find a path to target ...
		} else {
			// right next to target.
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
	ctl.act.NoReply(msg.Reply())
}

func (c *AggressiveController) BattleStart(msg message.Message) {
	logrus.Info("##### BattleStart #####")
	c.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), c.act.CallbackChan)
	c.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), c.act.CallbackChan)

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
	if !c.battleready {
		c.battleready = true
		c.ruler.NotifyActor(message.Create(nil, rulermethods.ControllerBattleReady{
			ControllerID: c.ID,
		}, nil))
	}
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

func (ctl *AggressiveController) ControllerMoveReply(msg message.Message) {
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

		attacker := ctl.KnownEntities[msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID]

		logrus.Info(" Attacker: ", attacker.PrettyString())

		atkrng := attacker.GetPropertyI(property.AttackRange).I()
		if msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.Position.Distance(target.Position) <= atkrng {

			// it is already in place. Send attack
			logrus.WithFields(logrus.Fields{
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
			}), ctl.act.CallbackChan)
		} else {
			// end of turn

			logrus.WithFields(logrus.Fields{
				"EntityID":       ControllerData.Entity.ID.String()[0:8],
				"Position":       ControllerData.Entity.Position,
				"Controller":     ctl.ID.String()[0:8],
				"ControllerName": ctl.act.Name(),
				"Expected":       target.Position}).Debug("Too far away from target, ending turn")

			ctl.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
				EntityID:     msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID,
				ControllerID: ctl.ID,
			}, rulermethods.EndOfTurn{}), ctl.act.CallbackChan)
		}
	}

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

func (ctl *AggressiveController) preparePathToEntity(pos position.Position, grd *grid.Grid, ent entity.Entity, jumpHeight int) []position.Position {
	path, found := grd.AStarPath(pos, ent.Position, jumpHeight)
	if !found {
		return nil
	}

	return path
}

// Implement the actor.Communication interface
// Notify Actor
func (c *AggressiveController) NotifyActor(msg message.Message) {
	c.act.Notify(msg)
}

func (c *AggressiveController) SendActor(msg message.Message, callback chan message.Message) {
	c.act.Send(msg, callback)
}
