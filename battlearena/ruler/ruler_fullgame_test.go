package ruler

import (
	"os"
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllers"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestRulerControllerFullGame(t *testing.T) {

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	ruler := NewRuler()
	ctrl := controllers.NewAggressiveController("Fake1")
	ctrl2 := controllers.NewAggressiveController("Fake2")
	ctrl.Start()
	ctrl2.Start()

	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl,
		ControllerID: ctrl.ID,
	}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl2,
		ControllerID: ctrl2.ID,
	}, nil))

	go func() {
		<-time.After(20 * time.Second)
		ctrl.PrintStack()
		ctrl2.PrintStack()
		ruler.PrintStack()
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

	// add t.Fail() if you want to check the logs.
}
