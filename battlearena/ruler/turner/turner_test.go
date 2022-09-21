package turner

import (
	"testing"

	"github.com/google/uuid"
)

func TestTurner(t *testing.T) {
	turn := NewTurner()
	turn.AddEntity(uuid.New(), 1)
	turn.AddEntity(uuid.New(), 2)
	turn.AddEntity(uuid.New(), 3)
	turn.AddEntity(uuid.New(), 4)

	currententity := turn.NextTurn()
	if currententity == uuid.Nil {
		t.Errorf("currententity should not be nil")
	}
	if turn.CurrentEntityTurn != currententity {
		t.Errorf("currententity should be %v", currententity)
	}
	if len(turn.Turns) != 3 {
		t.Errorf("len(turn.Turns) should be 3")
	}
	if turn.Turns[0].delay != 1 {
		t.Errorf("turn.Turns[0].delay should be 1")
	}
}

// TestTurnerRemove AddEntity AddEntity RemoveEntity NextTurn
func TestTurnerRemove(t *testing.T) {
	turn := NewTurner()
	entity1 := uuid.New()
	entity2 := uuid.New()
	entity3 := uuid.New()
	entity4 := uuid.New()
	turn.AddEntity(entity1, 1)
	turn.AddEntity(entity2, 2)
	turn.AddEntity(entity3, 3)
	turn.AddEntity(entity4, 4)

	turn.RemoveEntity(entity2)
	currententity := turn.NextTurn()
	if currententity == uuid.Nil {
		t.Errorf("currententity should not be nil")
	}
	if turn.CurrentEntityTurn != currententity {
		t.Errorf("currententity should be %v", currententity)
	}
	if currententity != entity1 {
		t.Errorf("currententity should be %v", entity1)
	}
	if len(turn.Turns) != 2 {
		t.Errorf("len(turn.Turns) should be 2")
	}
	if turn.Turns[0].delay != 2 {
		t.Errorf("turn.Turns[0].delay should be 2")
	}
	if turn.Turns[0].entityid != entity3 {
		t.Errorf("turn.Turns[0].entityid should be %v", entity3)
	}
}
