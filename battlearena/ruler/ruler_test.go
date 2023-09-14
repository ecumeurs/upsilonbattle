package ruler

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilonmapdata/grid/position/pattern"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FakeController struct {
	*actor.Actor
	ID              uuid.UUID
	Stoppers        map[string]chan *message.Message
	KnownEntities   map[uuid.UUID]entity.Entity
	StopperCallback chan bool
}

// New
func NewFake(name string) *FakeController {
	ctrl := &FakeController{
		Actor:           actor.New(name),
		Stoppers:        make(map[string]chan *message.Message),
		ID:              uuid.New(),
		KnownEntities:   make(map[uuid.UUID]entity.Entity),
		StopperCallback: make(chan bool),
	}

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

	ctrl.Start()

	return ctrl
}

// AddStopper
func (c *FakeController) AddStopper(method interface{}) chan *message.Message {
	stopper := make(chan *message.Message)
	c.Stoppers[reflect.TypeOf(method).String()] = stopper
	return stopper
}

// AddStoppers
func (c *FakeController) AddStoppers(methods ...interface{}) {
	for _, method := range methods {
		c.AddStopper(method)
	}
}

// GetStopper
func (c *FakeController) GetStopper(method interface{}) chan *message.Message {
	return c.Stoppers[reflect.TypeOf(method).String()]
}

// Close all the stoppers
func (c *FakeController) Close() {
	for _, stopper := range c.Stoppers {
		close(stopper)
	}
}

// implement all FakeController methods handlers.

func (c *FakeController) SetQueue(msg *message.Message) bool {
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) Send(msg *message.Message) bool {
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) ReceiveAPIMessage(msg *message.Message) bool {
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) ControllerNextTurn(msg *message.Message) bool {
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) BattleStart(msg *message.Message) bool {
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) BattleEnd(msg *message.Message) bool {
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) ControllerAttacked(msg *message.Message) bool {
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) EntitiesStateChanged(msg *message.Message) bool {
	c.RequestLogger.WithFields(logrus.Fields{
		"Turn": msg.TargetMethod.(rulermethods.EntitiesStateChanged).Turn.String()}).Info("New Turn Received")
	// fill in known entities (clear them beforehand.)
	c.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range msg.TargetMethod.(rulermethods.EntitiesStateChanged).Entities {
		c.KnownEntities[e.ID] = e
	}
	c.NoReply(msg.Reply())
	return true
}

func (c *FakeController) GetStateReply(msg *message.Message) bool {
	return true
}

func (c *FakeController) GetGridStateReply(msg *message.Message) bool {
	return true
}

func (c *FakeController) GetEntitiesStateReply(msg *message.Message) bool {
	return true
}

func (c *FakeController) ControllerMoveReply(msg *message.Message) bool {
	return true
}

func (c *FakeController) ControllerAttackReply(msg *message.Message) bool {
	return true
}

func TestRulerBattleBegin(t *testing.T) {
	ruler := NewRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// expect a battle start
	ctrl.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()] = make(chan *message.Message)
	ctrl2.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()] = make(chan *message.Message)
	// defer close of stoppers
	defer func() {
		close(ctrl.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()])
		close(ctrl2.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()])
	}()

	endchan := make(chan bool)
	defer close(endchan)

	go func() {
		<-ctrl.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()]
		ctrl.StopperCallback <- true
		fmt.Println("BattleStart received by ctrl")
		<-ctrl2.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()]
		ctrl2.StopperCallback <- true
		fmt.Println("BattleStart received by ctrl2")
		endchan <- true
	}()

	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl,
		ControllerID: uuid.New(),
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl2,
		ControllerID: uuid.New(),
	}, nil))

	<-endchan
}

func TestRulerBattleBeginNextTurn(t *testing.T) {
	ruler := NewRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// expect a battle start
	ctrl.Stoppers[reflect.TypeOf(rulermethods.ControllerNextTurn{}).String()] = make(chan *message.Message)
	ctrl2.Stoppers[reflect.TypeOf(rulermethods.ControllerNextTurn{}).String()] = make(chan *message.Message)
	defer func() {
		close(ctrl.Stoppers[reflect.TypeOf(rulermethods.ControllerNextTurn{}).String()])
		close(ctrl2.Stoppers[reflect.TypeOf(rulermethods.ControllerNextTurn{}).String()])
	}()

	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl,
		ControllerID: ctrl.ID,
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl2,
		ControllerID: ctrl2.ID,
	}, nil))

	select {
	case msg := <-ctrl.Stoppers[reflect.TypeOf(rulermethods.ControllerNextTurn{}).String()]:
		fmt.Println("ControllerNextTurn received by ctrl")

		if msg.Content.(rulermethods.ControllerNextTurn).Entity.ControllerID != ctrl.ID {
			t.Error("ControllerNextTurn received by ctrl but ControllerID is not ctrl.ID")
		}
		ctrl.StopperCallback <- true
	case msg := <-ctrl2.Stoppers[reflect.TypeOf(rulermethods.ControllerNextTurn{}).String()]:
		fmt.Println("ControllerNextTurn received by ctrl2")

		if msg.Content.(rulermethods.ControllerNextTurn).Entity.ControllerID != ctrl2.ID {
			t.Error("ControllerNextTurn received by ctrl but ControllerID is not ctrl2.ID")
		}
		ctrl2.StopperCallback <- true
	}
}

func TestRulerBattleBeginNextTurnFetchGridAndEntities(t *testing.T) {
	ruler := NewRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// expect a battle start
	ctrl.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()] = make(chan *message.Message)
	ctrl2.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()] = make(chan *message.Message)
	defer func() {
		close(ctrl.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()])
		close(ctrl2.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()])
	}()
	endchan := make(chan bool)
	defer close(endchan)
	go func() {
		<-ctrl.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()]
		fmt.Println("BattleStart received by ctrl")
		ctrl.StopperCallback <- true
		<-ctrl2.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()]
		fmt.Println("BattleStart received by ctrl2")
		ctrl2.StopperCallback <- true
		endchan <- true
	}()

	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl,
		ControllerID: ctrl.ID,
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl2,
		ControllerID: ctrl2.ID,
	}, nil))

	<-endchan

	// now either controller should be able to access the grid and entities
	replyChan := make(chan *message.Message)
	defer close(replyChan)

	var grd *grid.Grid
	grd = nil

	ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, replyChan), replyChan)
	msg := <-replyChan
	if msg.Content.(rulermethods.GetGridStateReply).Grid == nil {
		t.Error("Grid should not be nil")
	}
	grd = msg.Content.(rulermethods.GetGridStateReply).Grid

	ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, replyChan), replyChan)
	msg = <-replyChan
	if msg.Content.(rulermethods.GetEntitiesStateReply).Entities == nil {
		t.Error("Entities should not be nil")
	}

	// controller should have at least more than one entities
	if len(msg.Content.(rulermethods.GetEntitiesStateReply).Entities) < 2 {
		t.Error("Entities should have at least 2 entities")
	}
	appropriateEntitiesCounter := 0
	entities := make([]entity.Entity, 0)
	for _, ent := range msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		if ent.ControllerID == ctrl.ID {
			appropriateEntitiesCounter++
			entities = append(entities, ent)
		}
	}
	if appropriateEntitiesCounter < 2 {
		t.Error("Entities should have at least 2 entities for controller")
	}

	//entities should be on board.
	for _, ent := range entities {
		pos := ent.Position
		c, found := grd.CellAt(pos)
		if !found {
			t.Error("Position should be on board")
		}
		if c.EntityID == uuid.Nil {
			t.Error("Cell should have an entity")
		}
		if c.EntityID != ent.ID {
			t.Error("Cell should have the same entity")
		}
	}
}

func TestRulerControllerCanMoveAttackAndEndTurn(t *testing.T) {

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	ruler := NewRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// expect a battle start
	ctrl.AddStoppers(rulermethods.BattleStart{}, rulermethods.ControllerNextTurn{}, rulermethods.EndOfTurn{}, rulermethods.ControllerMoveReply{}, rulermethods.ControllerAttackReply{}, rulermethods.ControllerAttackReply{})
	ctrl2.AddStoppers(rulermethods.BattleStart{}, rulermethods.ControllerNextTurn{}, rulermethods.EndOfTurn{}, rulermethods.ControllerMoveReply{}, rulermethods.ControllerAttackReply{}, rulermethods.ControllerAttacked{})

	endchan := make(chan bool)
	defer func() {
		logrus.Info("Shutting Down")
		ctrl.Close()
		ctrl2.Close()
		close(endchan)
	}()

	// wait for game to begin ...
	go func() {
		<-ctrl.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()]
		logrus.Info("BattleStart received by ctrl")
		ctrl.StopperCallback <- true
		<-ctrl2.Stoppers[reflect.TypeOf(rulermethods.BattleStart{}).String()]
		logrus.Info("BattleStart received by ctrl2")
		ctrl2.StopperCallback <- true
		endchan <- true
	}()

	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl,
		ControllerID: ctrl.ID,
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl2,
		ControllerID: ctrl2.ID,
	}, nil))

	<-endchan

	// now either controller should be able to access the grid and entities
	replyChan := make(chan *message.Message)
	defer close(replyChan)

	var grd *grid.Grid
	grd = nil

	ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, replyChan), replyChan)
	msg := <-replyChan
	grd = msg.Content.(rulermethods.GetGridStateReply).Grid

	ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, replyChan), replyChan)
	msg = <-replyChan
	entities := make([]entity.Entity, 0)
	foeEntities := make([]entity.Entity, 0)
	for _, ent := range msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		if ent.ControllerID == ctrl.ID {
			entities = append(entities, ent)
		} else {
			foeEntities = append(foeEntities, ent)
		}
	}

	attacker := entities[0]
	target := foeEntities[0]

	// designate a meeting point (attacker centric)
	attackerNewPos := position.New(grd.Length/2, grd.Width/2, 0)
	attackerNewPos.Z = grd.TopMostGroundAt(attackerNewPos.X, attackerNewPos.Y)
	attackerMovePath := pattern.PathTo2D(attackerNewPos.Substract(attacker.Position)).Apply2D(attacker.Position)
	// ensure Z is at the right place all along.
	for i, p := range attackerMovePath {
		attackerMovePath[i].Z = grd.TopMostGroundAt(p.X, p.Y)
	}

	targetNewPos := position.New(grd.Length/2+1, grd.Width/2, 0)
	targetNewPos.Z = grd.TopMostGroundAt(targetNewPos.X, targetNewPos.Y)
	targetMovePath := pattern.PathTo2D(targetNewPos.Substract(target.Position)).Apply2D(target.Position)
	// ensure Z is at the right place all along.
	for i, p := range targetMovePath {
		targetMovePath[i].Z = grd.TopMostGroundAt(p.X, p.Y)
	}

	logrus.WithFields(logrus.Fields{
		"attacker":       attacker.Position,
		"target":         target.Position,
		"attackerNewPos": attackerNewPos,
		"targetNewPos":   targetNewPos}).Info("Preparing to move")
	// move both so that they are adjacent

	done := false
	for !done {
		select {
		case msg := <-ctrl.GetStopper(rulermethods.ControllerNextTurn{}):
			logrus.WithFields(logrus.Fields{
				"RequestID": msg.RequestId.String()[0:8],
				"Turn":      msg.TargetMethod.(rulermethods.ControllerNextTurn).Turn.String(),
				"EntityID":  msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID.String()[0:8]}).Info("ControllerNextTurn received by ctrl")
			if msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID == attacker.ID {
				logrus.Info("Moving attacker")
				if !msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.Position.Equals(attackerNewPos) {
					logrus.WithFields(logrus.Fields{
						"EntityID": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID.String()[0:8],
						"Position": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.Position,
						"Expected": attackerNewPos}).Info("Moving attacker")
					ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
						EntityID:     attacker.ID,
						Path:         attackerMovePath,
						ControllerID: ctrl.ID,
					}, rulermethods.ControllerMoveReply{}), ctrl.CallbackChan)
				} else {
					// it is already in place. Send attack
					logrus.WithFields(logrus.Fields{
						"EntityID": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID.String()[0:8],
						"Position": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.Position}).Info("Attacking")
					ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
						EntityID:     attacker.ID,
						Target:       targetNewPos,
						ControllerID: ctrl.ID,
					}, rulermethods.ControllerAttackReply{}), ctrl.CallbackChan)
				}
			} else {
				// that's another entity, end turn.
				logrus.Info("Not attacker's turn, ending turn")
				ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
					EntityID:     msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID,
					ControllerID: ctrl.ID,
				}, rulermethods.EndOfTurn{}), ctrl.CallbackChan)
			}
			ctrl.StopperCallback <- true
		case msg := <-ctrl.GetStopper(rulermethods.ControllerMoveReply{}):
			if msg.HasError {
				t.Error("Error while moving attacker", msg.ErrorMessage)
				done = true
			} else {
				logrus.WithFields(logrus.Fields{
					"EntityID": msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID.String()[0:8],
					"Position": msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.Position,
					"Expected": attackerNewPos}).Info("ControllerMoveReply received by ctrl")

				if msg.Content.(rulermethods.ControllerMoveReply).Entity.ID != attacker.ID {
					t.Error("Wrong entity moved")
				}
				if !msg.Content.(rulermethods.ControllerMoveReply).Entity.Position.Equals(attackerNewPos) {
					t.Error("Entity not moved to the right position")
				}
				logrus.Info("Movement done, ending turn")
				ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
					EntityID:     msg.Content.(rulermethods.ControllerMoveReply).Entity.ID,
					ControllerID: ctrl.ID,
				}, rulermethods.EndOfTurn{}), ctrl.CallbackChan)
			}
			ctrl.StopperCallback <- true
		case msg := <-ctrl.GetStopper(rulermethods.ControllerAttackReply{}):
			logrus.WithFields(logrus.Fields{
				"EntityID": msg.TargetMethod.(rulermethods.ControllerAttack).EntityID.String()[0:8],
				"Position": msg.Content.(rulermethods.ControllerAttackReply).Entity.Position}).Info("ControllerAttackReply received by ctrl")
			if msg.Content.(rulermethods.ControllerAttackReply).Entity.ID != attacker.ID {
				t.Error("Wrong entity attacked")
			}
			logrus.Info("Attack done, ending turn")
			ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
				EntityID:     msg.Content.(rulermethods.ControllerAttackReply).Entity.ID,
				ControllerID: ctrl.ID,
			}, rulermethods.EndOfTurn{}), ctrl.CallbackChan)
			ctrl.StopperCallback <- true
		case <-ctrl.GetStopper(rulermethods.EndOfTurn{}):
			logrus.Info("EndOfTurn received by ctrl")
			ctrl.StopperCallback <- true
		// Controller 2
		case msg := <-ctrl2.GetStopper(rulermethods.ControllerAttacked{}):
			logrus.WithFields(logrus.Fields{
				"EntityID":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.ID.String()[0:8],
				"AttackerID": msg.TargetMethod.(rulermethods.ControllerAttacked).Attacker.ID.String()[0:8],
				"Position":   msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.Position}).Info("ControllerAttacked received by ctrl2")
			if msg.TargetMethod.(rulermethods.ControllerAttacked).Entity.ID != target.ID {
				t.Error("Wrong entity attacked")
			}
			if msg.TargetMethod.(rulermethods.ControllerAttacked).Attacker.ID != attacker.ID {
				t.Error("Wrong attacker")
			}
			logrus.Info("Attacked! END OF TEST")
			done = true // test is over
			ctrl2.StopperCallback <- true
		case msg := <-ctrl2.GetStopper(rulermethods.ControllerNextTurn{}):
			logrus.WithFields(
				logrus.Fields{
					"EntityID": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID.String()[0:8],
					"Turn":     msg.TargetMethod.(rulermethods.ControllerNextTurn).Turn,
					"Position": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.Position}).Info("ControllerNextTurn received by ctrl2")
			if msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID == target.ID {
				logrus.Info("Target's turn")
				if !msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.Position.Equals(targetNewPos) {
					logrus.WithFields(logrus.Fields{
						"EntityID": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID.String()[0:8],
						"Position": msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.Position,
						"Expected": targetNewPos}).Info("Target not moved to the right position")
					ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
						EntityID:     target.ID,
						Path:         targetMovePath,
						ControllerID: ctrl2.ID,
					}, rulermethods.ControllerMoveReply{}), ctrl2.CallbackChan)
				} else {
					// it is already in place. Send Wait
					logrus.Info("Target is at the right position, waiting")
					ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
						EntityID:     msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID,
						ControllerID: ctrl2.ID,
					}, rulermethods.EndOfTurn{}), ctrl2.CallbackChan)
				}
			} else {
				// that's another entity, end turn.
				logrus.Info("Not target's turn, ending turn")
				ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
					EntityID:     msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID,
					ControllerID: ctrl2.ID,
				}, rulermethods.EndOfTurn{}), ctrl2.CallbackChan)
			}
			ctrl2.StopperCallback <- true
		case msg := <-ctrl2.GetStopper(rulermethods.ControllerMoveReply{}):
			if msg.HasError {
				t.Error("Error while moving target", msg.ErrorMessage)
				done = true
			} else {
				logrus.WithFields(logrus.Fields{
					"EntityID": msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.ID.String()[0:8],
					"Position": msg.TargetMethod.(rulermethods.ControllerMoveReply).Entity.Position,
					"Expected": targetNewPos}).Info("ControllerMoveReply received by ctrl2")

				if msg.Content.(rulermethods.ControllerMoveReply).Entity.ID != target.ID {
					t.Error("Wrong entity moved")
				}
				if !msg.Content.(rulermethods.ControllerMoveReply).Entity.Position.Equals(targetNewPos) {
					t.Error("Entity not moved to the right position")
				}
				logrus.Info("Move done, ending turn")
				ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
					EntityID:     msg.Content.(rulermethods.ControllerMoveReply).Entity.ID,
					ControllerID: ctrl2.ID,
				}, rulermethods.EndOfTurn{}), ctrl2.CallbackChan)
			}
			ctrl2.StopperCallback <- true
		case <-ctrl2.GetStopper(rulermethods.EndOfTurn{}):
			logrus.Info("EndOfTurn received by ctrl2")
			ctrl2.StopperCallback <- true
		}
	}

	ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, replyChan), replyChan)
	msg = <-replyChan
	grd = msg.Content.(rulermethods.GetGridStateReply).Grid

	// check absence of target and presence of attacker
	if c, ok := grd.CellAt(attackerNewPos); ok {
		if c.EntityID == uuid.Nil {
			t.Error("No entity at attacker's position")
		} else if c.EntityID != attacker.ID {
			t.Error("Wrong entity at attacker's position")
		}
	}
	if c, ok := grd.CellAt(targetNewPos); ok {
		if c.EntityID != uuid.Nil {
			t.Error("Entity at target's position")
		}
	}

	ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, replyChan), replyChan)
	msg = <-replyChan
	for _, ent := range msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		if ent.ID == target.ID {
			t.Error("Target entity still present")
		}
	}
	_, err := msg.Content.(rulermethods.GetEntitiesStateReply).Turn.GetEntityDelay(target.ID)
	if err == nil {
		t.Error("Target entity still present in turn state")
	}

	logrus.Info("END OF LAST CHECKS")

	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{
		ControllerID: ctrl.ID,
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{
		ControllerID: ctrl2.ID,
	}, nil))

	ctrl1stop := ctrl.PrepareToStop()
	ctrl2stop := ctrl2.PrepareToStop()

	go func() {
		<-ctrl1stop
		logrus.Info("ctrl1 stopped")
		ctrl.Stop()
	}()

	go func() {
		<-ctrl2stop
		logrus.Info("ctrl2 stopped")
		ctrl2.Stop()
	}()

	after := time.After(1 * time.Second)
	<-after

	//t.Fail()
}
