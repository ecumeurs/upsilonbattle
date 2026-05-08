package turner

import (
	"testing"

	"github.com/google/uuid"
)

// TestTurner verifies the basic initiative sequencing logic.
// It ensures entities are correctly queued by delay and that NextTurn advances the timeline.
// @test-link [[mech_initiative]]
// @test-link [[mech_initiative_requeue_calculation]]
func TestTurner(t *testing.T) {
	// Initialize a new turner instance.
	turn := NewTurner()
	
	// Add four entities with sequential delays (1, 2, 3, 4).
	turn.AddEntity(uuid.New(), 1)
	turn.AddEntity(uuid.New(), 2)
	turn.AddEntity(uuid.New(), 3)
	turn.AddEntity(uuid.New(), 4)

	// Advance to the first turn.
	currententity := turn.NextTurn()
	
	// Verify that an entity was actually selected.
	if currententity == uuid.Nil {
		t.Errorf("currententity should not be nil")
	}
	
	// Ensure the turner state correctly tracks the active entity.
	// This confirms the bi-directional link between current turn and entity ID.
	if turn.CurrentEntityTurn != currententity {
		t.Errorf("currententity should be %v", currententity)
	}
	
	// Check that one entity was removed from the pending queue.
	if len(turn.Turns) != 3 {
		t.Errorf("len(turn.Turns) should be 3")
	}
	
	// Validate that the relative delay of the next entity was correctly adjusted.
	if turn.Turns[0].Delay != 1 {
		t.Errorf("turn.Turns[0].Delay should be 1")
	}
}

// TestTurnerRemove validates that removing an entity from the queue correctly updates the remaining turns.
// It ensures that the relative delays of other entities are preserved after removal.
// @test-link [[mech_initiative]]
func TestTurnerRemove(t *testing.T) {
	// Setup turner with multiple entities to test removal logic.
	turn := NewTurner()
	entity1 := uuid.New()
	entity2 := uuid.New()
	entity3 := uuid.New()
	entity4 := uuid.New()
	
	// Add entities with sequential delays to establish a predictable timeline.
	turn.AddEntity(entity1, 1)
	turn.AddEntity(entity2, 2)
	turn.AddEntity(entity3, 3)
	turn.AddEntity(entity4, 4)

	// Remove an entity that is in the middle of the queue (entity2 with delay 2).
	// This tests the queue's ability to collapse and maintain order after an entry is removed.
	turn.RemoveEntity(entity2)
	
	// Advance the turn to the first remaining entity in the adjusted timeline.
	currententity := turn.NextTurn()
	
	// Verify turn transition integrity to ensure the system didn't enter an invalid state.
	if currententity == uuid.Nil {
		t.Errorf("currententity should not be nil")
	}
	// The CurrentEntityTurn property must match the ID returned by NextTurn.
	if turn.CurrentEntityTurn != currententity {
		t.Errorf("currententity should be %v", currententity)
	}
	
	// Ensure entity1 was the one selected (it was first in queue before and after removal).
	if currententity != entity1 {
		t.Errorf("currententity should be %v", entity1)
	}
	
	// Verify that the queue length was reduced correctly. 
	// We started with 4, removed 1 (entity2), and consumed 1 (entity1), so 2 should remain.
	if len(turn.Turns) != 2 {
		t.Errorf("len(turn.Turns) should be 2")
	}
	
	// Check that entity3 is now first in queue and its delay is correctly calculated relative to the previous turn.
	// The new delay should be entity3.Delay (3) minus the consumed entity1.Delay (1), which equals 2.
	if turn.Turns[0].Delay != 2 {
		t.Errorf("turn.Turns[0].Delay should be 2")
	}
	// Also verify that the identity of the next entity in queue is indeed entity3.
	if turn.Turns[0].EntityId != entity3 {
		t.Errorf("turn.Turns[0].EntityId should be %v", entity3)
	}
	
	// This concludes the removal and transition test, confirming the timeline remains consistent.
}
