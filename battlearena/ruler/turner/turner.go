package turner

import (
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

func NewTurner() Turner {
	return Turner{
		Turns:             make([]entityturn, 0),
		CurrentEntityTurn: uuid.Nil,
	}
}

func (t *Turner) AddEntity(entityid uuid.UUID, delay int) {
	//insert entity at the right place according to delay
	for i := range t.Turns {
		if t.Turns[i].delay > delay {
			t.Turns = append(t.Turns[:i], append([]entityturn{{entityid, delay}}, t.Turns[i:]...)...)
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
	for i := range t.Turns {
		t.Turns[i].delay -= turn.delay
	}

	t.CurrentEntityTurn = turn.entityid
	return turn.entityid
}
