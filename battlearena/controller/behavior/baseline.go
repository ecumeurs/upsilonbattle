package behavior

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/google/uuid"
)

// AggressiveBehavior is the always-active baseline layer.
//
// It proposes to all three unset slots (Target, Move, Action) using a full A*
// pathfinder, mirroring the logic previously embedded in AggressiveController.
//
// BaseActivation = 1.0 — it always runs and sits at the bottom of every archetype stack.
// When a higher-priority layer has already set a slot, the baseline respects it.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type AggressiveBehavior struct{}

// Name returns the stable identifier used in memory records and logs.
func (*AggressiveBehavior) Name() string  {
	return "baseline_aggressive"
}

// BaseActivation returns 1.0, ensuring this layer always activates regardless of grade.
func (*AggressiveBehavior) BaseActivation() float64  {
	return 1.0
}

// Propose fills any unset slot in draft.
//   - Target: nearest living enemy (or respects an already-set target from a prior layer).
//   - Action: attack if target is in range and has-not-acted.
//   - Move:   A* path toward target (trimmed to movement budget and attack range).
func (b *AggressiveBehavior) Propose(ctx GameContext, draft *DecisionDraft) {
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	grd := ctx.Grid()

	selfTeam := self.GetPropertyI(property.TeamID).I()

	// ── 1. Resolve target ──────────────────────────────────────────────────────
	var target *entity.Entity

	if draft.Target != nil {
		if ent, ok := entities[draft.Target.EntityID]; ok && hpOf(ent) > 0 {
			target = &ent
		} else {
			// Prior layer's target is dead/gone — clear it so we pick a fresh one.
			draft.Target = nil
			ctx.Memory().ClearTarget()
		}
	}

	if target == nil {
		target = nearestEnemy(self, selfTeam, entities)
		if target == nil {
			return // no enemies — nothing to do
		}
		draft.Target = &TargetSlot{EntityID: target.ID}
	}

	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	atkRange := self.GetPropertyI(property.AttackRange).I()
	dist := tools.Distance(self.Position.X, self.Position.Y, target.Position.X, target.Position.Y)

	// ── 2. Propose Action if already in range ──────────────────────────────────
	if dist <= atkRange && !ctx.HasActed() {
		if draft.Action == nil {
			draft.Action = &ActionSlot{Type: ActionAttack, Target: target.Position}
		}
		return
	}

	// ── 3. Propose Move via A* ─────────────────────────────────────────────────
	if ctx.RemainingMovement() > 0 && draft.Move == nil && grd != nil {
		path, found := grd.AStarPath(
			self.Position,
			target.Position,
			jumpHeight,
			func(p position.Position) bool {
				return isStepBlocked(p, self.ID, entities, grd)
			},
		)
		if found && len(path) > 1 {
			movePath := trimToRange(path, ctx.RemainingMovement(), atkRange)
			if len(movePath) > 0 {
				draft.Move = &MoveSlot{Path: movePath}
			}
		}
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func hpOf(e entity.Entity) int {
	return e.GetPropertyI(property.HP).I()
}

// nearestEnemy finds the closest living enemy by Manhattan distance.
func nearestEnemy(self entity.Entity, selfTeam int, entities map[uuid.UUID]entity.Entity) *entity.Entity {
	var best *entity.Entity
	bestDist := int(^uint(0) >> 1)
	for _, ent := range entities {
		if ent.ID == self.ID {
			continue
		}
		if ent.GetPropertyI(property.TeamID).I() == selfTeam {
			continue
		}
		if hpOf(ent) <= 0 {
			continue
		}
		d := tools.Distance(self.Position.X, self.Position.Y, ent.Position.X, ent.Position.Y)
		if d < bestDist {
			bestDist = d
			e := ent
			best = &e
		}
	}
	return best
}

// isStepBlocked returns true if pos is occupied by a living entity (other than selfID)
// or is a non-ground cell.
func isStepBlocked(pos position.Position, selfID uuid.UUID, entities map[uuid.UUID]entity.Entity, grd *grid.Grid) bool {
	if entityOccupies(pos, selfID, entities) {
		return true
	}
	return cellBlocked(pos, selfID, entities, grd)
}

func entityOccupies(pos position.Position, selfID uuid.UUID, entities map[uuid.UUID]entity.Entity) bool {
	for _, ent := range entities {
		if ent.ID != selfID && hpOf(ent) > 0 && ent.Position.X == pos.X && ent.Position.Y == pos.Y {
			return true
		}
	}
	return false
}

func cellBlocked(pos position.Position, selfID uuid.UUID, entities map[uuid.UUID]entity.Entity, grd *grid.Grid) bool {
	if grd == nil {
		return false
	}
	cells := grd.CellsForPositions([]position.Position{pos})
	if len(cells) == 0 || cells[0] == nil {
		return false
	}
	c := cells[0]
	if c.Type != cell.Ground {
		return true
	}
	for _, entID := range c.EntityIDs {
		if entID != selfID {
			if _, present := entities[entID]; present {
				return true
			}
		}
	}
	return false
}

// trimToRange trims an A* path to the movement budget, stopping atkRange cells before
// the target so the entity ends up in attack range without entering the target cell.
//
// path[0] is the entity's current position; path[len-1] is the target.
// Returns the sub-slice path[1:limit] (cells to actually move through), or nil if
// no meaningful movement is possible.
func trimToRange(path []position.Position, moveBudget, atkRange int) []position.Position {
	if len(path) <= 1 {
		return nil
	}
	// How many cells do we want to traverse (stop atkRange before the target)?
	limit := len(path) - atkRange
	if limit < 0 {
		limit = 0
	}
	if limit > moveBudget {
		limit = moveBudget
	}
	if limit < 2 {
		return nil
	}
	// path[1:limit] skips the start position and walks up to limit steps.
	return path[1:limit]
}
