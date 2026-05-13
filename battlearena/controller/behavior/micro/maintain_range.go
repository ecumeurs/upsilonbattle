package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontools/tools"
)

// MaintainRange proposes a move toward the current target, stopping at ideal attack range.
// Works with KiteAway (which handles the too-close case).
// Used by the Ranger archetype.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type MaintainRange struct{}

// Name returns the stable identifier used in memory records and logs.
func (*MaintainRange) Name() string             {
	return "maintain_range"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*MaintainRange) BaseActivation() float64  {
	return 0.75
}

// Propose adjusts position to keep the current target within optimal attack range.
func (b *MaintainRange) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Move != nil || ctx.RemainingMovement() <= 0 {
		return
	}

	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()
	atkRange := self.GetPropertyI(property.AttackRange).I()

	target := nearestLivingEnemy(self, selfTeam, entities)
	if draft.Target != nil {
		if ent, ok := entities[draft.Target.EntityID]; ok && ent.GetPropertyI(property.HP).I() > 0 {
			e := ent
			target = &e
		}
	}
	if target == nil {
		return
	}

	if tools.Distance(self.Position.X, self.Position.Y, target.Position.X, target.Position.Y) <= atkRange {
		return // already in range; KiteAway handles too-close
	}

	grd := ctx.Grid()
	if grd == nil {
		return
	}
	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	path, found := grd.AStarPath(self.Position, target.Position, jumpHeight, nil)
	if !found || len(path) <= 1 {
		return
	}

	// Stop atkRange cells short of the target position.
	limit := len(path) - atkRange
	if limit < 0 {
		limit = 0
	}
	if limit > ctx.RemainingMovement() {
		limit = ctx.RemainingMovement()
	}
	if limit < 2 {
		return
	}
	draft.Move = &behavior.MoveSlot{Path: path[1:limit]}
}

var _ behavior.Behavior = (*MaintainRange)(nil)
