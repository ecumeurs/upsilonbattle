package turner

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

// TurnState is a read-only snapshot of the current initiative state.
// It is used for serializing and transmitting turn data to controllers or UI.
type TurnState struct {
	CurrentEntityTurn uuid.UUID
	RemainingTurns    []EntityTurn
}

// GetTurnState creates a snapshot of the current Turner state for synchronization.
// @spec-link [[mech_initiative]]
func (t Turner) GetTurnState() TurnState {
	turns := make([]EntityTurn, len(t.Turns))
	copy(turns, t.Turns)
	return TurnState{
		CurrentEntityTurn: t.CurrentEntityTurn,
		RemainingTurns:    turns,
	}
}

// String returns a human-readable representation of the initiative queue.
func (t Turner) String() string {
	s := t.CurrentEntityTurn.String()[0:8] + " 0,"
	for i := range t.Turns {
		s += t.Turns[i].EntityId.String()[0:8] + " " + fmt.Sprint("", t.Turns[i].Delay) + ","
	}
	return s
}

// String returns a human-readable representation of the turn state snapshot.
func (st TurnState) String() string {
	s := st.CurrentEntityTurn.String()[0:8] + " 0,"
	for i := range st.RemainingTurns {
		s += st.RemainingTurns[i].EntityId.String()[0:8] + " " + fmt.Sprint("", st.RemainingTurns[i].Delay) + ","
	}
	return s
}

// GetEntityDelay retrieves the delay for a specific entity from a static TurnState snapshot.
// @spec-link [[mech_initiative]]
func (st TurnState) GetEntityDelay(EntityId uuid.UUID) (int, error) {
	for i := range st.RemainingTurns {
		if st.RemainingTurns[i].EntityId == EntityId {
			return st.RemainingTurns[i].Delay, nil
		}
	}
	return 0, errors.New("entity not found")
}
