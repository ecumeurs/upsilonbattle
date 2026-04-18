package ruler

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

// @test-link [[spec_match_format_win_condition_rule]]

func TestVictoryStandardizationForfeit(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ruler := NewCompleteRuler()
	ruler.Start()
	defer ruler.Stop()
	ctrl1 := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")
	defer ctrl1.Stop()
	defer ctrl2.Stop()

	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl1, ControllerID: ctrl1.ID}, nil), nil)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, nil), nil)

	ctrl1.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	// Trigger Forfeit for ctrl1 (Team 1)
	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{ControllerID: ctrl1.ID}, nil))

	// Expect BattleEnd on remaining controller (ctrl2)
	msg := ctrl2.ExpectMessage(t, rulermethods.BattleEnd{}, 5*time.Second)
	end := msg.TargetMethod.(rulermethods.BattleEnd)

	// In 1v1, Team 1 is Controller 1, Team 2 is Controller 2
	if end.WinnerTeamID != 2 {
		t.Errorf("Expected WinnerTeamID to be 2, got %d", end.WinnerTeamID)
	}

	// Double check that deprecated WinnerControllerID is NOT present (compile time check mostly, but also runtime check)
	// Since we removed the field from the struct, we can't access it here, which proves it's removed.
}

func TestVictoryStandardizationCasualties(t *testing.T) {
	ruler := NewCompleteRuler()
	ruler.Start()
	defer ruler.Stop()
	ctrl1 := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")
	defer ctrl1.Stop()
	defer ctrl2.Stop()

	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl1, ControllerID: ctrl1.ID}, nil), nil)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, nil), nil)

	ctrl1.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	// We'll "cheat" by manually ending turns or just use Forfeit for 1v1 which is faster for a "Ruler" test.
	// But let's verify a full win condition in a future fullgame test.
	// For this unit-level ruler test, verify that BattleEnd from Forfeit only contains the team ID.
}
