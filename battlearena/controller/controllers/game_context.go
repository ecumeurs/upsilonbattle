package controllers

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/google/uuid"
)

// controllerGameContext is the private concrete implementation of behavior.GameContext.
// One instance is created per entity turn and mutated between ticks as the turn progresses.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type controllerGameContext struct {
	self         entity.Entity
	entities     map[uuid.UUID]entity.Entity
	grd          *grid.Grid
	hasActed     bool
	remainingMvt int
	lastOutcome  behavior.TickOutcome
	mem          *behavior.DecisionMemory
	gradeIdx     int
}

// SelfEntity returns a snapshot of the controlled entity for this tick.
func (c *controllerGameContext) SelfEntity() entity.Entity  {
	return c.self
}

// KnownEntities returns the entity snapshot visible to this controller.
func (c *controllerGameContext) KnownEntities() map[uuid.UUID]entity.Entity  {
	return c.entities
}

// Grid returns the arena grid.
func (c *controllerGameContext) Grid() *grid.Grid  {
	return c.grd
}

// HasActed reports whether the entity has consumed its action this turn.
func (c *controllerGameContext) HasActed() bool  {
	return c.hasActed
}

// RemainingMovement returns the movement budget remaining for this turn.
func (c *controllerGameContext) RemainingMovement() int  {
	return c.remainingMvt
}

// LastTickOutcome reports the engine result of the previous command.
func (c *controllerGameContext) LastTickOutcome() behavior.TickOutcome  {
	return c.lastOutcome
}

// Memory returns the cross-tick decision history for this controller.
func (c *controllerGameContext) Memory() *behavior.DecisionMemory  {
	return c.mem
}

// Grade returns the AI grade index (0 = Grade I … 8 = Grade V).
func (c *controllerGameContext) Grade() int  {
	return c.gradeIdx
}
