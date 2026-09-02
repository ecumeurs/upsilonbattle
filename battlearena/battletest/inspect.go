package battletest

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// Entity returns the actor's current entity snapshot from the game state.
func (s *Scenario) Entity(a *Actor) entity.Entity {
	return s.GS.Entities[a.ID]
}

// Alive reports whether the actor still exists in the game state (not killed/removed).
func (s *Scenario) Alive(a *Actor) bool {
	_, ok := s.GS.Entities[a.ID]
	return ok
}

// Pos returns the actor's current grid position.
func (s *Scenario) Pos(a *Actor) position.Position {
	return s.Entity(a).Position
}

// HP returns the actor's current HP.
func (s *Scenario) HP(a *Actor) int {
	return s.Entity(a).GetPropertyC(property.HP).GetValue()
}

// Movement returns the actor's remaining movement points.
func (s *Scenario) Movement(a *Actor) int {
	return s.Entity(a).GetPropertyC(property.Movement).GetValue()
}

// Shield returns the actor's current shield value.
func (s *Scenario) Shield(a *Actor) int {
	return s.Entity(a).GetPropertyC(property.Shield).GetValue()
}

// Poison returns the actor's current accumulated poison.
func (s *Scenario) Poison(a *Actor) int {
	return s.Entity(a).GetPropertyI(property.Poison).I()
}

// Stat returns an integer property value for the actor.
func (s *Scenario) Stat(a *Actor, p property.Key) int {
	return s.Entity(a).GetPropertyI(p).I()
}

// HasActed reports whether the actor has consumed its action this turn.
func (s *Scenario) HasActed(a *Actor) bool {
	return s.Entity(a).HasActed()
}

// TrapAt reports whether any positional effect is registered at pos. Combined with a
// RemoveOnTrigger trap, this lets tests assert whether a trap fired (consumed) or was
// flown over (still present).
func (s *Scenario) TrapAt(pos position.Position) bool {
	return len(s.GS.PositionalEffects[pos]) > 0
}
