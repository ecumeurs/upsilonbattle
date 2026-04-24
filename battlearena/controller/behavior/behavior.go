// Package behavior defines the pluggable AI behavior system for entity controllers.
//
// The behavior system decouples AI decision-making from controller plumbing.
// Each entity has exactly one Controller, but may use a composite Behavior that
// combines multiple atomic Behaviors (AggressiveBehavior, SpookedBehavior, etc.).
//
// NOTE: This is currently a one-controller-per-entity model.
// TODO: Evolve toward a single EntityController that holds a map[entityID]Behavior
//       so that all non-player, non-AI entities can be managed by one actor, reducing
//       goroutine count and simplifying coordination. See ISS-066 preparation doc §Phase 3.
//
// @spec-link [[mech_behavior_system]]
package behavior

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// DecisionType enumerates the actions a Behavior can request.
type DecisionType string

const (
	// DecisionPass means the entity should end its turn immediately.
	DecisionPass DecisionType = "pass"
	// DecisionMove means the entity should move along Decision.Path.
	DecisionMove DecisionType = "move"
	// DecisionAttack means the entity should attack Decision.Target.
	DecisionAttack DecisionType = "attack"
	// DecisionSkill means the entity should use skill Decision.SkillID at Decision.Target.
	DecisionSkill DecisionType = "skill"
)

// Decision is the output of a Behavior evaluation: what the entity should do next.
type Decision struct {
	Type    DecisionType
	Path    []position.Position // populated for DecisionMove
	Target  position.Position   // populated for DecisionAttack / DecisionSkill
	SkillID uuid.UUID           // populated for DecisionSkill
}

// GameContext provides the read-only game state a Behavior needs to make decisions.
// Using an interface here keeps behaviors testable without a real GameState.
//
// @spec-link [[mech_behavior_system]]
type GameContext interface {
	// SelfEntity returns the entity controlled by this behavior.
	SelfEntity() entity.Entity
	// KnownEntities returns all entities visible to this controller.
	KnownEntities() map[uuid.UUID]entity.Entity
	// Grid returns the current arena grid.
	Grid() *grid.Grid
}

// Behavior is the single decision point for AI-controlled entities.
// Each call to OnTurn should return the next Decision for the entity.
//
// @spec-link [[mech_behavior_system]]
type Behavior interface {
	// OnTurn is called when it is the entity's turn. It returns what to do next.
	OnTurn(ctx GameContext) Decision
}
