package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
)

const ambushCooldownTurns = 3

// Ambush makes the entity wait adjacent to an obstacle for up to N turns, then
// strikes once the target comes within attack range. Throttled to once every
// ambushCooldownTurns turns via DecisionMemory.
// Used by the Sneak archetype.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type Ambush struct{}

// Name returns the stable identifier used in memory records and logs.
func (*Ambush) Name() string             {
	return "ambush"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*Ambush) BaseActivation() float64  {
	return 0.65
}

// Propose moves self adjacent to an obstacle and holds until a target enters attack range.
func (b *Ambush) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	// Don't ambush too frequently.
	if ctx.Memory().TurnsSince(b.Name()) < ambushCooldownTurns {
		return
	}

	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()
	atkRange := self.GetPropertyI(property.AttackRange).I()

	target := nearestLivingEnemy(self, selfTeam, entities)
	if target == nil {
		return
	}

	dist := tools.Distance(self.Position.X, self.Position.Y, target.Position.X, target.Position.Y)

	// Phase 1: target is not yet in range — move adjacent to an obstacle to hide.
	if dist > atkRange && ctx.RemainingMovement() > 0 && draft.Move == nil {
		grd := ctx.Grid()
		if grd == nil {
			return
		}
		hidePos := findObstacleAdjacentCell(self.Position, grd)
		if hidePos == nil {
			return
		}
		jumpHeight := self.GetPropertyI(property.JumpHeight).I()
		path, found := grd.AStarPath(self.Position, *hidePos, jumpHeight, nil)
		if !found || len(path) <= 1 {
			return
		}
		budget := ctx.RemainingMovement()
		end := budget + 1
		if end > len(path) {
			end = len(path)
		}
		draft.Move = &behavior.MoveSlot{Path: path[1:end]}
		return
	}

	// Phase 2: target is in range — strike.
	if dist <= atkRange && !ctx.HasActed() && draft.Action == nil {
		draft.Target = &behavior.TargetSlot{EntityID: target.ID}
		draft.Action = &behavior.ActionSlot{Type: behavior.ActionAttack, Target: target.Position}
	}
}

type cellProvider interface {
	CellsForPositions([]position.Position) []*cell.Cell
}

// findObstacleAdjacentCell returns a ground cell adjacent to an obstacle near the entity,
// or nil if none found within a 2-cell search radius.
func findObstacleAdjacentCell(self position.Position, grd cellProvider) *position.Position {
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			candidate := position.New(self.X+dx, self.Y+dy, self.Z)
			if !isGroundCell(candidate, grd) {
				continue
			}
			if hasAdjacentObstacle(candidate, grd) {
				return &candidate
			}
		}
	}
	return nil
}

func isGroundCell(pos position.Position, grd cellProvider) bool {
	cells := grd.CellsForPositions([]position.Position{pos})
	return len(cells) > 0 && cells[0] != nil && cells[0].Type == cell.Ground
}

func hasAdjacentObstacle(pos position.Position, grd cellProvider) bool {
	for _, adj := range []position.Position{
		position.New(pos.X+1, pos.Y, pos.Z),
		position.New(pos.X-1, pos.Y, pos.Z),
		position.New(pos.X, pos.Y+1, pos.Z),
		position.New(pos.X, pos.Y-1, pos.Z),
	} {
		cells := grd.CellsForPositions([]position.Position{adj})
		if len(cells) > 0 && cells[0] != nil && cells[0].Type != cell.Ground {
			return true
		}
	}
	return false
}

var _ behavior.Behavior = (*Ambush)(nil)
