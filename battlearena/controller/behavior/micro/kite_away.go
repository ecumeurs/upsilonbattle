package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
)

// KiteAway proposes a retreat move when a foe has entered melee range.
// Used by the Ranger archetype to preserve attack range.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type KiteAway struct{}

// Name returns the stable identifier used in memory records and logs.
func (*KiteAway) Name() string             {
	return "kite_away"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*KiteAway) BaseActivation() float64  {
	return 0.8
}

// Propose steps the entity away from any foe within melee range to maintain safe distance.
func (b *KiteAway) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Move != nil || ctx.RemainingMovement() <= 0 {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	// Find nearest enemy.
	var nearestFoe *position.Position
	bestDist := int(^uint(0) >> 1)
	for _, ent := range entities {
		if ent.ID == self.ID || ent.GetPropertyI(property.TeamID).I() == selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		d := tools.Distance(self.Position.X, self.Position.Y, ent.Position.X, ent.Position.Y)
		if d < bestDist {
			bestDist = d
			p := ent.Position
			nearestFoe = &p
		}
	}

	if nearestFoe == nil {
		return
	}

	atkRange := self.GetPropertyI(property.AttackRange).I()
	// Only kite if the foe is within melee reach of us (closer than our attack range + 1).
	if bestDist > atkRange {
		return
	}

	// Flee: move to a position that is farther from the foe.
	grd := ctx.Grid()
	if grd == nil {
		return
	}

	// Compute a flee target: mirror self position away from the foe.
	fleeTarget := position.New(
		2*self.Position.X-nearestFoe.X,
		2*self.Position.Y-nearestFoe.Y,
		self.Position.Z,
	)

	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	path, found := grd.AStarPath(self.Position, fleeTarget, jumpHeight, nil)
	if !found || len(path) <= 1 {
		return
	}
	budget := ctx.RemainingMovement()
	end := budget + 1
	if end > len(path) {
		end = len(path)
	}
	draft.Move = &behavior.MoveSlot{Path: path[1:end]}
}

var _ behavior.Behavior = (*KiteAway)(nil)
