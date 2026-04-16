package turner

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// @spec-link [[mech_initiative]]
// @spec-link [[mech_action_economy]]
type EntityTurn struct {
	EntityId uuid.UUID
	Delay    int
}

// Turner: The actual context Turner lives in.
type Turner struct {
	Turns             []EntityTurn
	CurrentEntityTurn uuid.UUID
}

// TurnState: is a DTO (a snapshot) that is only used to transit information here and there.
type TurnState struct {
	CurrentEntityTurn uuid.UUID
	RemainingTurns    []EntityTurn
}

func NewTurner() Turner {
	return Turner{
		Turns:             make([]EntityTurn, 0),
		CurrentEntityTurn: uuid.Nil,
	}
}

func (t *Turner) ForceTurn(EntityId uuid.UUID) {
	if t.CurrentEntityTurn == EntityId {
		return
	}
	if t.CurrentEntityTurn != uuid.Nil {
		// reinsert it at 200
		t.AddEntity(t.CurrentEntityTurn, 200)
	}
	t.CurrentEntityTurn = EntityId

	// remove entity from turns
	for i := range t.Turns {
		if t.Turns[i].EntityId == EntityId {
			t.Turns = append(t.Turns[:i], t.Turns[i+1:]...)
			break
		}
	}
	// next turn removes entity from stack.
}

// GetEntityDelay returns the Delay of the entity or error if not present
func (t Turner) GetEntityDelay(EntityId uuid.UUID) (int, error) {
	for i := range t.Turns {
		if t.Turns[i].EntityId == EntityId {
			return t.Turns[i].Delay, nil
		}
	}
	return 0, errors.New("entity not found")
}

func (t *Turner) AddEntity(EntityId uuid.UUID, Delay int) {
	//insert entity at the right place according to Delay
	if len(t.Turns) == 0 {
		t.Turns = append(t.Turns, EntityTurn{EntityId: EntityId, Delay: Delay})
		return
	}

	for i := range t.Turns {
		if t.Turns[i].Delay > Delay {
			t.Turns = append(t.Turns[:i], append([]EntityTurn{{EntityId, Delay}}, t.Turns[i:]...)...)
			return
		}
	}
	// wasn't inserted, so it's the last one
	t.Turns = append(t.Turns, EntityTurn{EntityId: EntityId, Delay: Delay})
}

func (t *Turner) RemoveEntity(EntityId uuid.UUID) {
	for i := range t.Turns {
		if t.Turns[i].EntityId == EntityId {
			t.Turns = append(t.Turns[:i], t.Turns[i+1:]...)
			return
		}
	}
}

func (t *Turner) NextTurn() uuid.UUID {
	if len(t.Turns) == 0 {
		return uuid.Nil
	}
	turn := t.Turns[0]
	t.Turns = t.Turns[1:]
	// remove from remaining turns the Delay of current turn
	for i, ot := range t.Turns {
		ot.Delay -= turn.Delay
		t.Turns[i] = ot
	}

	t.CurrentEntityTurn = turn.EntityId
	return turn.EntityId
}

func (t Turner) String() string {
	s := t.CurrentEntityTurn.String()[0:8] + " 0,"
	for i := range t.Turns {
		s += t.Turns[i].EntityId.String()[0:8] + " " + fmt.Sprint("", t.Turns[i].Delay) + ","
	}
	return s
}

func (t Turner) GetTurnState() TurnState {
	return TurnState{
		CurrentEntityTurn: t.CurrentEntityTurn,
		RemainingTurns:    t.Turns,
	}
}

func (st TurnState) String() string {
	s := st.CurrentEntityTurn.String()[0:8] + " 0,"
	for i := range st.RemainingTurns {
		s += st.RemainingTurns[i].EntityId.String()[0:8] + " " + fmt.Sprint("", st.RemainingTurns[i].Delay) + ","
	}
	return s
}

// GetEntityDelay
func (st TurnState) GetEntityDelay(EntityId uuid.UUID) (int, error) {
	for i := range st.RemainingTurns {
		if st.RemainingTurns[i].EntityId == EntityId {
			return st.RemainingTurns[i].Delay, nil
		}
	}
	return 0, errors.New("entity not found")
}
