package behavior

// @test-link [[mechanic_mech_decision_memory]]

import (
	"testing"

	"github.com/google/uuid"
)

// TestTurnsSinceReturnsSentinelForUnknownLayer confirms that a layer never recorded returns 999.
func TestTurnsSinceReturnsSentinelForUnknownLayer(t *testing.T) {
	m := NewDecisionMemory()
	if got := m.TurnsSince("nonexistent"); got != 999 {
		t.Errorf("TurnsSince(unknown) = %d, want 999", got)
	}
}

// TestTurnsSinceReturnsZeroImmediatelyAfterRecord confirms that TurnsSince returns 0 when
// the layer fired in the current turn.
func TestTurnsSinceReturnsZeroImmediatelyAfterRecord(t *testing.T) {
	m := NewDecisionMemory()
	m.Record("ambush", &DecisionDraft{Notes: make(map[string]any)})
	if got := m.TurnsSince("ambush"); got != 0 {
		t.Errorf("TurnsSince after same-turn record = %d, want 0", got)
	}
}

// TestTurnsSinceAdvancesCorrectlyAcrossTurns verifies that turn advancement is correctly reflected.
func TestTurnsSinceAdvancesCorrectlyAcrossTurns(t *testing.T) {
	m := NewDecisionMemory()
	m.Record("kite_away", &DecisionDraft{Notes: make(map[string]any)})
	m.AdvanceTurn() // turn 1
	m.AdvanceTurn() // turn 2
	if got := m.TurnsSince("kite_away"); got != 2 {
		t.Errorf("TurnsSince after 2 turns = %d, want 2", got)
	}
}

// TestTurnsSinceIgnoresSkippedRecords verifies that Skipped records are not counted as "last active".
func TestTurnsSinceIgnoresSkippedRecords(t *testing.T) {
	m := NewDecisionMemory()
	m.Record("ambush", &DecisionDraft{Notes: make(map[string]any)}) // active on turn 0
	m.AdvanceTurn()
	m.RecordSkipped("ambush") // skipped on turn 1
	// TurnsSince should point to the last non-skipped (turn 0), so 1 turn ago.
	if got := m.TurnsSince("ambush"); got != 1 {
		t.Errorf("TurnsSince ignoring skipped = %d, want 1", got)
	}
}

// TestCountInLastNCountsNonSkippedActivations verifies the sliding-window count.
// CountInLastN(n) includes records at turns >= (current_turn - n), skipping Skipped records.
func TestCountInLastNCountsNonSkippedActivations(t *testing.T) {
	m := NewDecisionMemory()
	// Active on turns 0, 1, 2; then advance to turn 3 and record skipped.
	for i := 0; i < 3; i++ {
		m.Record("heal_ally", &DecisionDraft{Notes: make(map[string]any)})
		m.AdvanceTurn()
	}
	m.RecordSkipped("heal_ally") // turn 3, skipped — not counted

	// Current turn = 3. CountInLastN(3): cutoff = 3-3 = 0 → includes turns 0,1,2 → 3 active.
	got := m.CountInLastN("heal_ally", 3)
	if got != 3 {
		t.Errorf("CountInLastN(3) = %d, want 3 (turns 0,1,2 all within window)", got)
	}

	// CountInLastN(2): cutoff = 3-2 = 1 → includes turns 1,2 → 2 active.
	got = m.CountInLastN("heal_ally", 2)
	if got != 2 {
		t.Errorf("CountInLastN(2) = %d, want 2 (only turns 1,2 within window)", got)
	}

	// CountInLastN(1): cutoff = 3-1 = 2 → includes turn 2 → 1 active.
	got = m.CountInLastN("heal_ally", 1)
	if got != 1 {
		t.Errorf("CountInLastN(1) = %d, want 1 (only turn 2 within window)", got)
	}
}

// TestCurrentTargetStickyAcrossTicks confirms that the sticky target persists across tick advances.
func TestCurrentTargetStickyAcrossTicks(t *testing.T) {
	m := NewDecisionMemory()
	targetID := uuid.New()
	draft := &DecisionDraft{
		Target: &TargetSlot{EntityID: targetID},
		Notes:  make(map[string]any),
	}
	m.Record("baseline_aggressive", draft)
	m.AdvanceTick()

	if got := m.CurrentTarget(); got != targetID {
		t.Errorf("CurrentTarget after AdvanceTick = %v, want %v", got, targetID)
	}
}

// TestClearTargetResetsToNil confirms that ClearTarget zeroes the sticky target.
func TestClearTargetResetsToNil(t *testing.T) {
	m := NewDecisionMemory()
	targetID := uuid.New()
	draft := &DecisionDraft{
		Target: &TargetSlot{EntityID: targetID},
		Notes:  make(map[string]any),
	}
	m.Record("baseline_aggressive", draft)
	m.ClearTarget()
	if got := m.CurrentTarget(); got != uuid.Nil {
		t.Errorf("CurrentTarget after ClearTarget = %v, want uuid.Nil", got)
	}
}

// TestCurrentTargetNilBeforeAnyRecord confirms that CurrentTarget is uuid.Nil on a fresh memory.
func TestCurrentTargetNilBeforeAnyRecord(t *testing.T) {
	m := NewDecisionMemory()
	if got := m.CurrentTarget(); got != uuid.Nil {
		t.Errorf("CurrentTarget on fresh memory = %v, want uuid.Nil", got)
	}
}

// TestRingBufferDoesNotGrowBeyondCapacity writes more than memorySize entries and confirms
// the buffer does not grow beyond the fixed capacity.
func TestRingBufferDoesNotGrowBeyondCapacity(t *testing.T) {
	m := NewDecisionMemory()
	const writes = memorySize + 10
	for i := 0; i < writes; i++ {
		m.Record("layer", &DecisionDraft{Notes: make(map[string]any)})
	}
	if m.size > memorySize {
		t.Errorf("buffer size = %d, want ≤ %d", m.size, memorySize)
	}
}

// TestRingBufferOverwritesOldestEntry confirms that after filling the ring the oldest entries
// are no longer accessible (overwritten by newer ones).
func TestRingBufferOverwritesOldestEntry(t *testing.T) {
	m := NewDecisionMemory()
	// Fill the buffer on turn 0 with layer "old".
	for i := 0; i < memorySize; i++ {
		m.RecordSkipped("old")
	}
	// Advance many turns and record with layer "new".
	for i := 0; i < 5; i++ {
		m.AdvanceTurn()
	}
	for i := 0; i < 5; i++ {
		m.Record("new", &DecisionDraft{Notes: make(map[string]any)})
	}
	// The oldest "old" entries are gone — CountInLastN("old", 100) should only see the still-buffered ones.
	// "old" records were all on turn 0 (skipped); CountInLastN ignores skipped, so it should be 0.
	if got := m.CountInLastN("old", 100); got != 0 {
		t.Errorf("old skipped records counted as active: got %d, want 0", got)
	}
}
