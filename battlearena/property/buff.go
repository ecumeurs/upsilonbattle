package property

import "github.com/google/uuid"

type TemporaryProperties struct {
	Properties map[string]Property
	Duration   int
	Forever    bool
	// conditions?

	// Origin Informations
	OriginEntityID uuid.UUID
	OriginSkillID  uuid.UUID
}

// MakeTemporaryProperties will return a new TemporaryProperties
func MakeTemporaryProperties(duration int) TemporaryProperties {
	return TemporaryProperties{
		Properties: make(map[string]Property),
		Duration:   duration,
		Forever:    false,
	}
}

// TickDown will decrease duration by 1, if duration reaches 0 return true
func (t *TemporaryProperties) TickDown() bool {
	if t.Forever {
		return false
	}

	t.Duration--
	return t.Duration <= 0
}
