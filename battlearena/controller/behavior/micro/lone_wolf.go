package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
)

const loneWolfBuffer = 2

// LoneWolf steers away from ally clumping. If more than one ally is within
// loneWolfBuffer cells, proposes a move away from the nearest ally.
// Used by the Sneak archetype to maintain separation.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type LoneWolf struct{}

// Name returns the stable identifier used in memory records and logs.
func (*LoneWolf) Name() string             {
	return "lone_wolf"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*LoneWolf) BaseActivation() float64  {
	return 0.6
}

// Propose repositions away from clustered allies to avoid target-concentration penalties.
func (b *LoneWolf) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Move != nil || ctx.RemainingMovement() <= 0 {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	var nearestAllyPos *position.Position
	nearbyAllies := 0
	bestDist := int(^uint(0) >> 1)

	for _, ent := range entities {
		if ent.ID == self.ID || ent.GetPropertyI(property.TeamID).I() != selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		d := tools.Distance(self.Position.X, self.Position.Y, ent.Position.X, ent.Position.Y)
		if d <= loneWolfBuffer {
			nearbyAllies++
			if d < bestDist {
				bestDist = d
				p := ent.Position
				nearestAllyPos = &p
			}
		}
	}

	if nearbyAllies <= 1 || nearestAllyPos == nil {
		return // no clumping
	}

	// Move away from the nearest ally.
	grd := ctx.Grid()
	if grd == nil {
		return
	}
	fleeTarget := position.New(
		2*self.Position.X-nearestAllyPos.X,
		2*self.Position.Y-nearestAllyPos.Y,
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

var _ behavior.Behavior = (*LoneWolf)(nil)
