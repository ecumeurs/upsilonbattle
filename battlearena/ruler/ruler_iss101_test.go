package ruler
// @test-link [[mechanic_arena_lifecycle]]
// @test-link [[mech_initiative]]
// @test-link [[uc_combat_turn]]
// @test-link [[mech_ruler_behavior]]

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// TestRulerAddEntityRejectsNilControllerID is a regression test for ISS-101.
//
// ISS-101 was diagnosed from a CI/production trace (GitHub Actions run
// 27952883421): a live player's entity vanished from the turn queue without
// ever appearing as current_entity_id, while the match continued normally
// for every other entity until the player's team was wiped having taken
// zero actions. The root cause turned out to be structural, not a one-off
// race: nothing on the entity-creation path enforced that an entity must
// have an owning controller.
//
//   - entity.New() (upsilontypes/entity/entity.go:61-72) defaults
//     ControllerID to uuid.Nil.
//   - behavior.go's doc comment and docs/mech_ruler_behavior.atom.md
//     (status: DRAFT) described Nil ControllerID as the *intentional*
//     signal for "non-actor" entities (traps, turrets, simple summons) —
//     but no production code path actually created such entities (the
//     entity.Type values TimeBased/Trap/AreaEffect/Obstacle are defined
//     but never instantiated outside tests; "traps" today are
//     PositionalEffects keyed by grid position, not Entities), and the one
//     test exercising the closest real mechanic
//     (rules_iss066_test.go: TestRuleEntityExpiration) assigned a real,
//     non-nil ControllerID — contradicting that "design."
//
// Conclusion: every entity, with no exception (including any future
// trap/bomb/summon), must have an owning controller. uuid.Nil is not a
// valid ControllerID. AddEntity now rejects it at the door, per this
// project's Crash Early / Fail Fast doctrine (.agent/rules/COMMON.md) —
// the failure must surface at registration time, not as a silent no-op
// turn deep in a live match.
func TestRulerAddEntityRejectsNilControllerID(t *testing.T) {
	r := NewRuler(uuid.New())
	r.GameState.Grid = grid.NewGrid(10, 10, 3)

	ghost := entity.New() // ControllerID left at its zero-value: uuid.Nil.
	ghost.Type = entity.Character
	ghost.Position = position.Position{X: 1, Y: 1, Z: 3}
	ghost.RepsertPropertyValue(property.TeamID, 1)

	if ghost.ControllerID != uuid.Nil {
		t.Fatalf("test setup invariant broken: expected entity.New() to default ControllerID to uuid.Nil")
	}

	defer func() {
		if recover() == nil {
			t.Error("expected AddEntity to panic/reject an entity with a Nil ControllerID (ISS-101)")
		}
		if _, ok := r.GameState.Entities[ghost.ID]; ok {
			t.Error("AddEntity panicked but the entity was registered anyway — the guard must run before any state mutation")
		}
	}()

	r.AddEntity(ghost)
}

// TestRulerHandTurnRejectsNilControllerID is a defense-in-depth regression
// test for ISS-101, complementing TestRulerAddEntityRejectsNilControllerID.
//
// AddEntity is the sanctioned, canonical entry point for putting an entity
// on the board — and it now rejects a Nil ControllerID outright. But
// handTurn (ruler_turn.go) used to ALSO treat a Nil ControllerID as a
// legitimate "automated entity" signal, silently resolving the turn via
// ExpirationBehavior with zero notification to anyone. That branch has
// been removed; reaching handTurn with a Nil ControllerID is now always
// treated as an invariant violation, not a routing decision.
//
// This test exercises handTurn directly (bypassing AddEntity on purpose —
// the only way left to construct the precondition, now that AddEntity
// rejects it) to prove handTurn defends itself independently, rather than
// relying solely on AddEntity never letting a bad entity through.
func TestRulerHandTurnRejectsNilControllerID(t *testing.T) {
	r := NewRuler(uuid.New())
	r.GameState.Grid = grid.NewGrid(10, 10, 3)

	ghost := entity.New() // ControllerID left at its zero-value: uuid.Nil.
	ghost.Type = entity.Character
	ghost.Position = position.Position{X: 1, Y: 1, Z: 3}
	ghost.RepsertPropertyValue(property.TeamID, 1)

	// Deliberately bypass AddEntity — it would (correctly) panic here. This
	// is a focused unit test of handTurn's own internal guard, not a claim
	// that this is how an entity should ever reach the board in practice.
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), ghost.Position, ghost.ID)
	r.GameState.Entities[ghost.ID] = ghost
	r.GameState.Turner.AddEntity(ghost.ID, 0)

	defer func() {
		if recover() == nil {
			t.Error("expected handTurn to panic for an entity with a Nil ControllerID (ISS-101)")
		}
	}()

	r.handTurn(ghost.ID)
}
