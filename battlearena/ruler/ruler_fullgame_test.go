package ruler

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

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

type FakeFGController struct {
	act            *actor.Actor
	ID             uuid.UUID
	KnownEntities  map[uuid.UUID]entity.Entity
	Grid           *grid.Grid
	ruler          actor.Communication
	latestTarget   entity.Entity
	BattleFinished chan bool
}

// New
func NewFakeFG(name string) *FakeFGController {
	ctrl := &FakeFGController{
		ID:             uuid.New(),
		KnownEntities:  make(map[uuid.UUID]entity.Entity),
		BattleFinished: make(chan bool),
	}

	ctrl.act = actor.New(name)
	ctrl.act.SetReceiveMessageHandler(ctrl.handleMessage)
	ctrl.act.SetReplyMessageHandler(ctrl.handleReply)
	ctrl.act.Start()

	return ctrl
}

func (c *FakeFGController) handleMessage(msg message.Message) bool {
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

func (c *FakeFGController) handleReply(msg message.Message) bool {
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

// implement all FakeFGController methods handlers.

func (c *FakeFGController) SetValidatorAndQueue(msg message.Message, m controllermethods.SetValidatorAndQueue) {
	c.ruler = m.Ruler
	c.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), c.act.CallbackChan)
	c.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), c.act.CallbackChan)

	c.act.NoReply(msg.Reply())
}

func (c *FakeFGController) Send(msg message.Message) {
	c.act.NoReply(msg.Reply())
}

func (c *FakeFGController) ReceiveAPIMessage(msg message.Message) {
	c.act.NoReply(msg.Reply())
}

func (c *FakeFGController) ControllerNextTurn(msg message.Message) {
	c.nextTurn(msg) // will move or attack.
	c.act.NoReply(msg.Reply())
}

func (c *FakeFGController) BattleStart(msg message.Message) {
	logrus.Info("##### BattleStart #####")
	c.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), c.act.CallbackChan)
	c.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), c.act.CallbackChan)

	c.act.NoReply(msg.Reply())
}

func (c *FakeFGController) BattleEnd(msg message.Message) {
	logrus.Info("##### BattleEnd #####")
	c.BattleFinished <- true

	c.act.NoReply(msg.Reply())
}

func (c *FakeFGController) ControllerAttacked(msg message.Message) {
	logrus.WithFields(logrus.Fields{
		"EntityID":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.ID.String()[0:8],
		"AttackerID": msg.TargetMethod.(rulermethods.ControllerAttacked).Attacker.ID.String()[0:8],
		"Position":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.Position}).Debug("ControllerAttacked")
	// nothing to do post attack

	c.act.NoReply(msg.Reply())
}

func (c *FakeFGController) EntitiesStateChanged(msg message.Message) {
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

func (c *FakeFGController) GetStateReply(msg message.Message) {
}

func (c *FakeFGController) GetGridStateReply(msg message.Message) {
	c.Grid = msg.Content.(rulermethods.GetGridStateReply).Grid
}

func (c *FakeFGController) GetEntitiesStateReply(msg message.Message) {
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

func (c *FakeFGController) ControllerMoveReply(msg message.Message) {
	c.afterMoveAttack(c.KnownEntities, msg)
}

func (c *FakeFGController) ControllerAttackReply(msg message.Message) {
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
func selectNearestFoe(currentEntity entity.Entity, entities map[uuid.UUID]entity.Entity) (entity.Entity, error) {
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
		return entity.Entity{}, errors.New("No nearest entity")
	}
}

func preparePathToEntity(pos position.Position, grd *grid.Grid, ent entity.Entity) []position.Position {
	path := pattern.PathTo2D(ent.Position.Substract(pos)).Apply2D(pos)
	// ensure Z is at the right place all along.
	for i, p := range path {
		path[i].Z = grd.TopMostGroundAt(p.X, p.Y)
	}

	return path
}

func (ctl *FakeFGController) nextTurn(msg message.Message) {
	controllerData := msg.TargetMethod.(rulermethods.ControllerNextTurn)
	logrus.WithFields(logrus.Fields{
		"Controller":     ctl.ID.String()[0:8],
		"ControllerName": ctl.act.Name(),
		"RequestID":      msg.RequestId.String()[0:8],
		"Turn":           controllerData.Turn.String(),
		"EntityID":       controllerData.Entity.ID.String()[0:8]}).Info("##### Turn BEGIN #####")
	target, err := selectNearestFoe(controllerData.Entity, ctl.KnownEntities)
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

	path := preparePathToEntity(controllerData.Entity.Position, ctl.Grid, target)
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

func (ctl *FakeFGController) afterMoveAttack(knownEntities map[uuid.UUID]entity.Entity, msg message.Message) {
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
func (c *FakeFGController) NotifyActor(msg message.Message) {
	c.act.Notify(msg)
}

func (c *FakeFGController) SendActor(msg message.Message, callback chan message.Message) {
	c.act.Send(msg, callback)
}

func TestFindNextAppropriateTarget(t *testing.T) {
	ctrl := &FakeFGController{
		ID:             uuid.New(),
		KnownEntities:  make(map[uuid.UUID]entity.Entity),
		BattleFinished: make(chan bool),
	}
	for i := 0; i < 10; i++ {
		e := entity.Entity{
			ID:       uuid.New(),
			Position: position.Position{X: i * 10, Y: i * 3},
		}
		if i%2 == 0 {
			e.ControllerID = uuid.New()
		} else {
			e.ControllerID = ctrl.ID
		}
		ctrl.KnownEntities[e.ID] = e
	}

	// find first entity of this controller
	for _, e := range ctrl.KnownEntities {
		if e.ControllerID == ctrl.ID {
			ent, err := selectNearestFoe(e, ctrl.KnownEntities)

			logrus.WithFields(logrus.Fields{
				"found_pos":        ent.Position,
				"found_entity":     ent.ID.String()[0:8],
				"found_controller": ent.ControllerID.String()[0:8]}).Info("found")
			if err != nil {
				t.Error(err)
			} else {
				if ent.ControllerID == ctrl.ID {
					t.Error("Found own entity")
				}
				if ent.ID == e.ID {
					t.Error("Found same entity")
				}
			}
		}
	}

}

func TestRulerControllerFullGame(t *testing.T) {

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	ruler := NewRuler()
	ctrl := NewFakeFG("Fake1")
	ctrl2 := NewFakeFG("Fake2")

	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl,
		ControllerID: ctrl.ID,
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl2,
		ControllerID: ctrl2.ID,
	}, nil))

	// when the messages queues are stuck, this is good :) (like wicked)
	go func() {
		<-time.After(20 * time.Second)
		ctrl.act.GetQueue().PrintStack()
		ctrl2.act.GetQueue().PrintStack()
		ruler.act.GetQueue().PrintStack()
	}()

	<-ctrl.BattleFinished
	<-ctrl2.BattleFinished

	logrus.Info("Battle Finished, doing end of game Checks")

	replyChan := make(chan message.Message)
	ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, nil), replyChan)
	msg := <-replyChan

	ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, nil), replyChan)
	msg = <-replyChan

	entities := msg.Content.(rulermethods.GetEntitiesStateReply).Entities
	// there must be only entities of one controller left.
	foundCtrlers := make(map[uuid.UUID]bool)
	for _, ent := range entities {
		foundCtrlers[ent.ControllerID] = true
	}
	if len(foundCtrlers) != 1 {
		logrus.Error("Expected only one controller left, found ", len(foundCtrlers))
		t.Fail()
	}

	logrus.Info("END OF LAST CHECKS")

	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{
		ControllerID: ctrl.ID,
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{
		ControllerID: ctrl2.ID,
	}, nil))

	ctrl1stop := ctrl.act.PrepareToStop()
	ctrl2stop := ctrl2.act.PrepareToStop()

	go func() {
		<-ctrl1stop
		logrus.Info("ctrl1 stopped")
		ctrl.act.Stop()
	}()

	go func() {
		<-ctrl2stop
		logrus.Info("ctrl2 stopped")
		ctrl2.act.Stop()
	}()

	after := time.After(1 * time.Second)
	<-after

	// add t.Fail() if you want to check the logs.
}
