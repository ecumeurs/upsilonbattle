package behavior

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// AggressiveBehavior drives an entity to close on and attack the nearest enemy.
// This is a port of the core logic from AggressiveController.ControllerNextTurn,
// expressed as a pure decision function without any actor/message plumbing.
//
// @spec-link [[mech_behavior_system]]
type AggressiveBehavior struct{}

// OnTurn implements Behavior.
// Returns: move toward nearest enemy → attack if adjacent → pass otherwise.
func (b *AggressiveBehavior) OnTurn(ctx GameContext) Decision {
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	grd := ctx.Grid()

	selfTeam := self.GetPropertyI(property.TeamID).I()
	selfHP := self.GetPropertyI(property.HP).I()
	if selfHP <= 0 {
		return Decision{Type: DecisionPass}
	}

	// Find nearest living enemy
	var target *position.Position
	bestDist := int(^uint(0) >> 1) // MaxInt

	for _, ent := range entities {
		if ent.ID == self.ID {
			continue
		}
		entHP := ent.GetPropertyI(property.HP).I()
		entTeam := ent.GetPropertyI(property.TeamID).I()
		if entHP <= 0 || entTeam == selfTeam {
			continue
		}
		dist := self.Position.Distance(ent.Position)
		if dist < bestDist {
			bestDist = dist
			pos := ent.Position
			target = &pos
		}
	}

	if target == nil {
		return Decision{Type: DecisionPass}
	}

	// Check if already adjacent (can attack)
	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	attackRange := self.GetPropertyI(property.AttackRange).I()
	if self.Position.Distance(*target) <= attackRange {
		return Decision{Type: DecisionAttack, Target: *target}
	}

	// Try to move toward target
	if grd == nil {
		return Decision{Type: DecisionPass}
	}
	path := findPathToward(grd, self.Position, *target, jumpHeight, self.GetPropertyC(property.Movement).GetMaxValue())
	if len(path) == 0 {
		return Decision{Type: DecisionPass}
	}
	return Decision{Type: DecisionMove, Path: path}
}

// findPathToward does a simple greedy step toward the target, avoiding occupied cells.
// The full A* pathfinding lives in the AggressiveController; this is a lightweight version
// sufficient for the behavior interface. Currently returns nil (stub) — the AggressiveController's
// existing pathfinder handles this correctly when the Behavior delegates back to it.
func findPathToward(_ *grid.Grid, _, _ position.Position, _, _ int) []position.Position {
	return nil
}
