package turner

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type entityturn struct {
	entityid uuid.UUID
	delay    int
}

type Turner struct {
	Turns             []entityturn
	CurrentEntityTurn uuid.UUID
}

type TurnState struct {
	CurrentEntityTurn uuid.UUID
	RemainingTurns    []entityturn
}

func NewTurner() Turner {
	return Turner{
		Turns:             make([]entityturn, 0),
		CurrentEntityTurn: uuid.Nil,
	}
}

func (t *Turner) ForceTurn(entityid uuid.UUID) {
	if t.CurrentEntityTurn == entityid {
		return
	}
	if t.CurrentEntityTurn != uuid.Nil {
		// reinsert it at 200
		t.AddEntity(t.CurrentEntityTurn, 200)
	}
	t.CurrentEntityTurn = entityid

	// remove entity from turns
	for i := range t.Turns {
		if t.Turns[i].entityid == entityid {
			t.Turns = append(t.Turns[:i], t.Turns[i+1:]...)
			break
		}
	}
	// next turn removes entity from stack.
}

// GetEntityDelay returns the delay of the entity or error if not present
func (t Turner) GetEntityDelay(entityid uuid.UUID) (int, error) {
	for i := range t.Turns {
		if t.Turns[i].entityid == entityid {
			return t.Turns[i].delay, nil
		}
	}
	return 0, errors.New("entity not found")
}

func (t *Turner) AddEntity(entityid uuid.UUID, delay int) {
	//insert entity at the right place according to delay
	if len(t.Turns) == 0 {
		t.Turns = append(t.Turns, entityturn{entityid: entityid, delay: delay})
		return
	}

	for i := range t.Turns {
		if t.Turns[i].delay > delay {
			t.Turns = append(t.Turns[:i], append([]entityturn{{entityid, delay}}, t.Turns[i:]...)...)
			return
		}
	}
	// wasn't inserted, so it's the last one
	t.Turns = append(t.Turns, entityturn{entityid: entityid, delay: delay})
}

func (t *Turner) RemoveEntity(entityid uuid.UUID) {
	for i := range t.Turns {
		if t.Turns[i].entityid == entityid {
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
	// remove from remaining turns the delay of current turn
	for i, ot := range t.Turns {
		ot.delay -= turn.delay
		t.Turns[i] = ot
	}

	t.CurrentEntityTurn = turn.entityid
	return turn.entityid
}

func (t Turner) String() string {
	s := t.CurrentEntityTurn.String()[0:8] + " 0,"
	for i := range t.Turns {
		s += t.Turns[i].entityid.String()[0:8] + " " + fmt.Sprint("", t.Turns[i].delay) + ","
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
		s += st.RemainingTurns[i].entityid.String()[0:8] + " " + fmt.Sprint("", st.RemainingTurns[i].delay) + ","
	}
	return s
}

// GetEntityDelay
func (st TurnState) GetEntityDelay(entityid uuid.UUID) (int, error) {
	for i := range st.RemainingTurns {
		if st.RemainingTurns[i].entityid == entityid {
			return st.RemainingTurns[i].delay, nil
		}
	}
	return 0, errors.New("entity not found")
}
