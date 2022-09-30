package property

type TemporaryProperties struct {
	Properties map[string]Property
	Duration   int
	// conditions?
}

// MakeTemporaryProperties will return a new TemporaryProperties
func MakeTemporaryProperties(duration int) TemporaryProperties {
	return TemporaryProperties{
		Properties: make(map[string]Property),
		Duration:   duration,
	}
}

// TickDown will decrease duration by 1, if duration reaches 0 return true
func (t *TemporaryProperties) TickDown() bool {
	t.Duration--
	return t.Duration <= 0
}
