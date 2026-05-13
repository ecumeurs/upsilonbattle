// Package behavior defines the pluggable AI behavior system for entity controllers.
//
// Architecture — two behavior systems exist in this codebase with distinct roles:
//
//   - ruler/behavior  — simple, message-producing behaviors for auto-entities
//     (traps, expiring effects, turrets). No tactical reasoning.
//
//   - controller/behavior  (this package) — layered, draft-based tactical AI for
//     player-equivalent entities (archetypes: Fighter, Ranger, Support, Sneak).
//     Each Behavior layer proposes to a shared DecisionDraft; the LayeredBehavior
//     orchestrator resolves one EngineCommand per tick.
//
// Pipeline model:
//
//	LayeredBehavior.Tick(ctx) → EngineCommand
//	  for each Behavior in Layers (top → bottom):
//	    roll activation against grade-scaled rate → skip if fails
//	    layer.Propose(ctx, draft)    // first-writer-wins per slot
//	  draft.Resolve(ctx)             // Action > Move > EndOfTurn
//
// @spec-link [[mechanic_mech_behavior_layered]]
package behavior

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/google/uuid"
)

// GameContext provides the read-only per-tick view a Behavior layer needs to make decisions.
// The concrete implementation lives in the controllers package (controllerGameContext),
// which owns the per-controller entity snapshot.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type GameContext interface {
	// SelfEntity returns a snapshot of the entity controlled by this behavior.
	SelfEntity() entity.Entity
	// KnownEntities returns all entities visible to this controller (snapshot per tick).
	KnownEntities() map[uuid.UUID]entity.Entity
	// Grid returns the arena grid (static geometry).
	Grid() *grid.Grid

	// HasActed reports whether the entity has already used its action this turn.
	HasActed() bool
	// RemainingMovement returns the movement budget left for this turn.
	RemainingMovement() int
	// LastTickOutcome reports the result of the previous engine command.
	LastTickOutcome() TickOutcome

	// Memory returns the cross-tick/cross-turn decision history for this entity.
	// @spec-link [[mechanic_mech_decision_memory]]
	Memory() *DecisionMemory
	// Grade returns the AI grade index (0 = Grade I … 8 = Grade V).
	// Used by LayeredBehavior to scale each layer's activation rate.
	Grade() int
}

// Behavior is the single interface for all AI decision layers.
//
// Each layer proposes to a shared DecisionDraft; it may write to any unset slot
// (first-writer-wins) or abstain entirely. The baseline layer (AggressiveBehavior)
// sits at the bottom of every stack with BaseActivation 1.0, ensuring a decision
// is always reached.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type Behavior interface {
	// Propose may set any unset slot on draft.
	// Convention: check draft.Target == nil before proposing Target, etc.
	Propose(ctx GameContext, draft *DecisionDraft)
	// BaseActivation is the declared activation probability [0, 1].
	// LayeredBehavior scales this by grade before rolling.
	// Use 1.0 for the always-active baseline.
	BaseActivation() float64
	// Name returns a stable identifier used in logs and DecisionMemory records.
	Name() string
}
