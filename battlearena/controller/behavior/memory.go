package behavior

import "github.com/google/uuid"

const memorySize = 20

// DecisionRecord captures what a single behavior layer did on a given tick.
//
// @spec-link [[mechanic_mech_decision_memory]]
type DecisionRecord struct {
	LayerName string
	Skipped   bool      // true if the layer's activation roll failed
	Turn      int
	Tick      int
	TargetID  uuid.UUID // uuid.Nil when Target slot was not touched
}

// DecisionMemory is a fixed-size ring buffer of per-tick layer records, shared across
// all ticks within a turn and across turns for a single entity.
//
// Layers use it to self-throttle (e.g. "don't Ambush again for 3 turns") and to
// carry a sticky target across ticks within a turn.
//
// @spec-link [[mechanic_mech_decision_memory]]
type DecisionMemory struct {
	entries [memorySize]DecisionRecord
	head    int
	size    int
	target  uuid.UUID // most recently committed target, persists across ticks
	tick    int
	turn    int
}

// NewDecisionMemory creates an empty DecisionMemory.
func NewDecisionMemory() *DecisionMemory {
	return &DecisionMemory{}
}

// Record saves a layer's contribution for the current tick.
// If the layer proposed a Target, that target becomes the sticky current target.
func (m *DecisionMemory) Record(layerName string, draft *DecisionDraft) {
	targetID := uuid.Nil
	if draft.Target != nil {
		targetID = draft.Target.EntityID
		m.target = targetID
	}
	m.add(DecisionRecord{
		LayerName: layerName,
		Skipped:   false,
		Turn:      m.turn,
		Tick:      m.tick,
		TargetID:  targetID,
	})
}

// RecordSkipped saves that a layer was skipped (activation roll failed) this tick.
func (m *DecisionMemory) RecordSkipped(layerName string) {
	m.add(DecisionRecord{
		LayerName: layerName,
		Skipped:   true,
		Turn:      m.turn,
		Tick:      m.tick,
	})
}

// AdvanceTick should be called by the controller each time a new pipeline tick begins.
func (m *DecisionMemory) AdvanceTick() {
	m.tick++
}

// AdvanceTurn should be called by the controller when a new game turn begins for this entity.
func (m *DecisionMemory) AdvanceTurn() {
	m.turn++
	m.tick = 0
}

// CurrentTarget returns the most recently committed target entity ID.
// Returns uuid.Nil if no target has been selected yet.
func (m *DecisionMemory) CurrentTarget() uuid.UUID {
	return m.target
}

// ClearTarget resets the sticky target. Call when the current target is dead or gone.
func (m *DecisionMemory) ClearTarget() {
	m.target = uuid.Nil
}

// TurnsSince returns how many turns have elapsed since layerName last actively participated
// (i.e. was not skipped). Returns 999 if the layer has never participated.
func (m *DecisionMemory) TurnsSince(layerName string) int {
	for i := m.size - 1; i >= 0; i-- {
		e := m.at(i)
		if e.LayerName == layerName && !e.Skipped {
			return m.turn - e.Turn
		}
	}
	return 999
}

// CountInLastN returns how many turns layerName actively participated within the last n turns.
func (m *DecisionMemory) CountInLastN(layerName string, n int) int {
	cutoff := m.turn - n
	count := 0
	for i := 0; i < m.size; i++ {
		e := m.at(i)
		if e.LayerName == layerName && !e.Skipped && e.Turn >= cutoff {
			count++
		}
	}
	return count
}

func (m *DecisionMemory) add(r DecisionRecord) {
	m.entries[m.head] = r
	m.head = (m.head + 1) % memorySize
	if m.size < memorySize {
		m.size++
	}
}

// at returns the i-th entry where 0 is the oldest and m.size-1 is the newest.
func (m *DecisionMemory) at(i int) DecisionRecord {
	idx := (m.head - m.size + i + memorySize*2) % memorySize
	return m.entries[idx]
}
