package behavior

// Shared test helpers for the behavior package tests.
// Defines the mockCtx GameContext implementation and entity/draft factories.

import (
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/google/uuid"
)

// mockCtx satisfies GameContext for deterministic unit tests.
type mockCtx struct {
	self     entity.Entity
	known    map[uuid.UUID]entity.Entity
	grd      *grid.Grid
	hasActed bool
	movement int
	outcome  TickOutcome
	mem      *DecisionMemory
	grade    int
}

func (m *mockCtx) SelfEntity() entity.Entity                  { return m.self }
func (m *mockCtx) KnownEntities() map[uuid.UUID]entity.Entity { return m.known }
func (m *mockCtx) Grid() *grid.Grid                           { return m.grd }
func (m *mockCtx) HasActed() bool                             { return m.hasActed }
func (m *mockCtx) RemainingMovement() int                     { return m.movement }
func (m *mockCtx) LastTickOutcome() TickOutcome                { return m.outcome }
func (m *mockCtx) Memory() *DecisionMemory                    { return m.mem }
func (m *mockCtx) Grade() int                                 { return m.grade }

// newCtx returns a mockCtx with sensible defaults: grade V, 3 movement, no grid.
func newCtx() *mockCtx {
	self := newEnt(0, 0, 1)
	return &mockCtx{
		self:     self,
		known:    map[uuid.UUID]entity.Entity{self.ID: self},
		mem:      NewDecisionMemory(),
		movement: 3,
		grade:    8, // Grade V: maximum activation rate
	}
}

// newEnt builds a minimal test entity at position (x,y,z=1) on teamID with full HP.
func newEnt(x, y, teamID int) entity.Entity {
	e := entity.Entity{
		ID:       uuid.New(),
		Position: position.New(x, y, 1),
	}
	e.Properties = make(map[string]property.Property)
	e.RepsertPropertyValue(property.TeamID, teamID)
	e.RepsertPropertyCMaxValue(property.HP, 10)
	e.RepsertPropertyCValue(property.HP, 10)
	return e
}

// newEntHP builds an entity with a specific HP value (useful for testing below-threshold conditions).
func newEntHP(x, y, teamID, hp, maxHP int) entity.Entity {
	e := newEnt(x, y, teamID)
	e.RepsertPropertyCMaxValue(property.HP, maxHP)
	e.RepsertPropertyCValue(property.HP, hp)
	return e
}
