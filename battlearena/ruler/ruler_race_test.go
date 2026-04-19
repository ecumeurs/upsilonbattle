package ruler

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapmaker/gridgenerator"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// TestShotClockRace highlights the race condition in startShotClock (ISS-047)
// where the ShotClock goroutine accesses GameState without synchronization.
// @test-link [[rule_turn_clock]]
func TestShotClockRace(t *testing.T) {
	r := NewRuler(uuid.New())
	r.Start()
	defer r.Stop()

	// Set a very short duration to trigger the goroutine quickly
	r.ShotClockDuration = 1 * time.Microsecond

	// In a loop, start the shot clock and simultaneously modify the turn
	// This should reliably trigger the race detector
	for i := 0; i < 100; i++ {
		r.startShotClock()
		r.GameState.IncTurn()
	}
}

// TestInitRace highlights the race condition between Ruler initialization (ISS-047)
// and direct state modification (as done in bridge.go).
// @test-link [[module_upsilonapi]]
func TestInitRace(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := uuid.New()
		r := NewRuler(id)
		// r.Start() // Uncomment this to see the race condition once setup is no longer sacrosanct
		
		// These calls happen before the Ruler actor loop is started.
		// Following commit 488bca1, setup is done synchronously.
		r.SetGrid(gridgenerator.GeneratePlainSquare(10, 10))
		r.SetNbControllers(2)
		r.AddEntity(entity.Entity{ID: uuid.New()})
		
		r.Start()
		r.Stop()
	}
}

// TestStartArenaWithoutGrid highlights ISS-010 where a battle can start
// even if the grid is missing from the GameState.
// @spec-link [[spec_match_format_ready_to_start_rule]]
func TestStartArenaWithoutGrid(t *testing.T) {
	matchID := uuid.New()
	pID := uuid.New()

	r := NewRuler(matchID)
	// Force Grid to nil to test Readiness Guard
	r.GameState.Grid = nil
	r.SetNbControllers(1)
	r.Start()
	defer r.Stop()
	
	// Setup a fake controller
	// Since we are in package ruler, we can use NewFake from ruler_test.go
	ctrl := NewFake("FakeController")
	defer ctrl.Stop()
	
	// 2. Add the controller. This should trigger BattleStart if NbControllers is reached.
	dChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.AddController{
		Controller:   ctrl,
		ControllerID: pID,
	}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	
	// 3. Wait for BattleStart.
	// In a FIXED version, this should NOT be received because Grid is nil.
	timeout := time.After(2 * time.Second)
	select {
	case msg := <-ctrl.Inbox:
		if _, ok := msg.TargetMethod.(rulermethods.BattleStart); ok {
			t.Fatal("Battle should NOT have started without a grid (ISS-010 failure)")
		}
	case <-timeout:
		// Success: No BattleStart message received
	}
}
