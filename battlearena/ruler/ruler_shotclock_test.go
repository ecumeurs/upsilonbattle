package ruler

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// @test-link [[rule_turn_clock]]
func TestShotClockExpiry(t *testing.T) {
	ruler := NewCompleteRuler()
	defer ruler.Stop()
	// Set a very short shot clock for testing
	ruler.ShotClockDuration = 100 * time.Millisecond
	
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	// Setup controllers
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, nil), nil)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, nil), nil)

	// Wait for battle start
	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 2*time.Second)
	
	// Wait for first entity turn
	msg := ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 2*time.Second)
	entID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID
	
	t.Logf("Turn started for entity %s, waiting for timeout...", entID)

	// Wait for shot clock to expire (100ms) and trigger EndOfTurn
	// This should result in EntitiesStateChanged being broadcast
	ctrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 2*time.Second)
	
	// And then the NEXT turn should be triggered
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 2*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// @test-link [[rule_turn_clock]]
func TestShotClockCancellation(t *testing.T) {
	ruler := NewCompleteRuler()
	defer ruler.Stop()
	// Set a duration long enough to manually intervene but short enough for a fast test
	ruler.ShotClockDuration = 500 * time.Millisecond
	
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")

	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, nil), nil)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, nil), nil)

	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 2*time.Second)
	
	// First turn starts
	msg := ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 2*time.Second)
	entID := msg.TargetMethod.(rulermethods.ControllerNextTurn).Entity.ID

	// Manually end turn before timeout
	ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		ControllerID: ctrl.ID,
		EntityID:     entID,
	}, nil), nil)

	// Wait for transition
	ctrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 2*time.Second)
	
	// Verify next turn triggered
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 2*time.Second)

	ctrl.Stop()
	ctrl2.Stop()
}

// @test-link [[rule_turn_clock]]
func TestShotClockTurnProtection(t *testing.T) {
	ruler := NewCompleteRuler()
	defer ruler.Stop()
	ruler.ShotClockDuration = 100 * time.Millisecond
	
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")
	
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, nil), nil)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, nil), nil)
	
	// Start battle
	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 2*time.Second)
	
	// Turn 1.0 starts
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 2*time.Second)
	
	// Wait for timeout 1 (transitions from 1.0 to 2.0)
	ctrl.ExpectMessage(t, rulermethods.EntitiesStateChanged{}, 2*time.Second)
	
	// Turn 2.0 starts
	ctrl.ExpectMessage(t, rulermethods.ControllerNextTurn{}, 2*time.Second)
	
	ctrl.Stop()
	ctrl2.Stop()
}
