package turner

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// EntityTurn represents a single turn entry in the initiative queue.
// It links an entity to its relative delay in the timeline.
type EntityTurn struct {
	EntityId uuid.UUID
	Delay    int
}

// Turner manages the sequencing of entity turns within a battle arena.
// It maintains a sorted queue of upcoming turns based on temporal delay.
// @spec-link [[mech_initiative]]
type Turner struct {
	Turns             []EntityTurn
	CurrentEntityTurn uuid.UUID
}

// TurnState is a read-only snapshot of the current initiative state.
// It is used for serializing and transmitting turn data to controllers or UI.
type TurnState struct {
	CurrentEntityTurn uuid.UUID
	RemainingTurns    []EntityTurn
}

// NewTurner initializes an empty Turner context.
// By default, it has no turns and no active entity.
// @spec-link [[mech_initiative]]
func NewTurner() Turner {
	return Turner{
		Turns:             make([]EntityTurn, 0),
		CurrentEntityTurn: uuid.Nil,
	}
}

// ForceTurn immediately overrides the current turn with the specified entity.
// If another entity was currently active, it is requeued with a default delay.
// This is typically used for interrupts or special action sequencing.
// @spec-link [[mech_initiative]]
func (t *Turner) ForceTurn(EntityId uuid.UUID) {
	if t.CurrentEntityTurn == EntityId {
		return
	}
	if t.CurrentEntityTurn != uuid.Nil {
		// Reinsert the displaced entity with a fixed penalty delay (200).
		t.AddEntity(t.CurrentEntityTurn, 200)
	}
	t.CurrentEntityTurn = EntityId

	// Remove the entity from any future positions in the turn queue.
	for i := range t.Turns {
		if t.Turns[i].EntityId == EntityId {
			t.Turns = append(t.Turns[:i], t.Turns[i+1:]...)
			break
		}
	}
}

// GetEntityDelay retrieves the accumulated delay for a specific entity in the queue.
// Returns an error if the entity is not currently scheduled for a turn.
// @spec-link [[mech_initiative_delay_costs]]
func (t Turner) GetEntityDelay(EntityId uuid.UUID) (int, error) {
	for i := range t.Turns {
		if t.Turns[i].EntityId == EntityId {
			return t.Turns[i].Delay, nil
		}
	}
	return 0, errors.New("entity not found")
}

// AddEntity inserts a new entity into the initiative queue at a position determined by its delay.
// The queue is maintained in ascending order of delay (earliest turns first).
// @spec-link [[mech_initiative]]
func (t *Turner) AddEntity(EntityId uuid.UUID, Delay int) {
	// If the queue is empty, simply append the new turn.
	if len(t.Turns) == 0 {
		t.Turns = append(t.Turns, EntityTurn{EntityId: EntityId, Delay: Delay})
		return
	}

	// Iterate through the queue to find the correct insertion point based on delay.
	for i := range t.Turns {
		if t.Turns[i].Delay > Delay {
			t.Turns = append(t.Turns[:i], append([]EntityTurn{{EntityId, Delay}}, t.Turns[i:]...)...)
			return
		}
	}
	// If the delay is greater than all existing entries, append to the end.
	t.Turns = append(t.Turns, EntityTurn{EntityId: EntityId, Delay: Delay})
}

// RemoveEntity permanently deletes an entity from the initiative queue.
// If the removed entity was the current active entity, the active slot is cleared.
// This is used when an entity dies or is removed from the match.
// @spec-link [[mech_initiative]]
func (t *Turner) RemoveEntity(EntityId uuid.UUID) {
	// Clear CurrentEntityTurn if it matches the removed entity to prevent shot-clock issues.
	if t.CurrentEntityTurn == EntityId {
		t.CurrentEntityTurn = uuid.Nil
	}

	for i := range t.Turns {
		if t.Turns[i].EntityId == EntityId {
			t.Turns = append(t.Turns[:i], t.Turns[i+1:]...)
			return
		}
	}
}

// NextTurn advances the initiative queue to the next entity.
// It subtracts the elapsed delay from all remaining entities and sets the new active entity.
// Returns the ID of the entity whose turn it now is, or uuid.Nil if the queue is empty.
// @spec-link [[mech_initiative_requeue_calculation]]
func (t *Turner) NextTurn() uuid.UUID {
	if len(t.Turns) == 0 {
		return uuid.Nil
	}
	turn := t.Turns[0]
	t.Turns = t.Turns[1:]
	
	// Normalize the delays of all remaining entities by subtracting the current turn's delay.
	for i, ot := range t.Turns {
		ot.Delay -= turn.Delay
		t.Turns[i] = ot
	}

	t.CurrentEntityTurn = turn.EntityId
	return turn.EntityId
}

// String returns a human-readable representation of the initiative queue.
func (t Turner) String() string {
	s := t.CurrentEntityTurn.String()[0:8] + " 0,"
	for i := range t.Turns {
		s += t.Turns[i].EntityId.String()[0:8] + " " + fmt.Sprint("", t.Turns[i].Delay) + ","
	}
	return s
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

// String returns a human-readable representation of the turn state snapshot.
func (st TurnState) String() string {
	s := st.CurrentEntityTurn.String()[0:8] + " 0,"
	for i := range st.RemainingTurns {
		s += st.RemainingTurns[i].EntityId.String()[0:8] + " " + fmt.Sprint("", st.RemainingTurns[i].Delay) + ","
	}
	return s
}

// GetEntityDelay retrieves the delay for a specific entity from a static TurnState snapshot.
// @spec-link [[mech_initiative_delay_costs]]
func (st TurnState) GetEntityDelay(EntityId uuid.UUID) (int, error) {
	for i := range st.RemainingTurns {
		if st.RemainingTurns[i].EntityId == EntityId {
			return st.RemainingTurns[i].Delay, nil
		}
	}
	return 0, errors.New("entity not found")
}
