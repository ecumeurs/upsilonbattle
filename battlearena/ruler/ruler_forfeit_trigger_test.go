package ruler

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRulerForfeitTriggersWinnerID(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.Start()
	ctrl1 := NewFake("P1")
	ctrl2 := NewFake("P2")

	// 1. Join controllers
	dChan := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl1, ControllerID: ctrl1.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	// 2. Wait for battle start
	ctrl1.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	// 3. Forfeit P1 (Controller 1)
	// We use ControllerForfeit which normally takes an EntityID, but should work with uuid.Nil for the whole team
	ruler.SendActor(message.Create(nil, rulermethods.ControllerForfeit{
		ControllerID: ctrl1.ID,
		EntityID:     uuid.Nil,
	}, rulermethods.ControllerForfeit{}), dChan)
	<-dChan

	// 4. Assert BattleEnd is received with P2 as winner
	msg1 := ctrl1.ExpectMessage(t, rulermethods.BattleEnd{}, 5*time.Second)
	msg2 := ctrl2.ExpectMessage(t, rulermethods.BattleEnd{}, 5*time.Second)

	endEvent1 := msg1.TargetMethod.(rulermethods.BattleEnd)
	endEvent2 := msg2.TargetMethod.(rulermethods.BattleEnd)

	assert.Equal(t, 2, endEvent1.WinnerTeamID, "P1 should see Team 2 as winner")
	assert.Equal(t, 2, endEvent2.WinnerTeamID, "P2 should see Team 2 as winner")

	// 5. Assert internal Ruler state
	// Use TestingGetState to avoid race
	replyChan := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.TestingGetState{}, rulermethods.TestingGetStateReply{}), replyChan)
	replyMsg := <-replyChan
	st := replyMsg.TargetMethod.(rulermethods.TestingGetStateReply)

	assert.Equal(t, Finished.String(), st.CurrentState, "Ruler should be in Finished state")
	assert.Equal(t, 2, st.WinnerTeamID, "GameState should persist WinnerTeamID")

	ctrl1.Stop()
	ctrl2.Stop()
}
