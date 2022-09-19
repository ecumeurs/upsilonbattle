package turner

import (
	"github.com/google/uuid"
)

type entityturn struct {
	entityid uuid.UUID
	delay    int
}

type Turner struct {
	Turns []entityturn
}

func NewTurner() Turner {
	return Turner{
		Turns: make([]entityturn, 0),
	}
}

func (t *Turner) AddEntity(entityid uuid.UUID, delay int) {
	t.Turns = append(t.Turns, entityturn{
		entityid: entityid,
		delay:    delay,
	})
}

func (t *Turner) NextTurn() uuid.UUID {
	if len(t.Turns) == 0 {
		return uuid.Nil
	}
	turn := t.Turns[0]
	t.Turns = t.Turns[1:]
	return turn.entityid
}
