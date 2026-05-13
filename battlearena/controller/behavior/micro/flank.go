package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontools/tools"
)

// Flank proposes a movement that approaches the target from the side or rear
// rather than head-on. If the entity is already in a flanking position, it attacks.
// Used by the Sneak archetype as a light alternative to BackstabSeeker.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type Flank struct{}

// Name returns the stable identifier used in memory records and logs.
func (*Flank) Name() string             {
	return "flank"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*Flank) BaseActivation() float64  {
	return 0.7
}

// Propose routes self to a lateral or rear cell relative to the current target.
func (b *Flank) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
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

	dist := tools.Distance(self.Position.X, self.Position.Y, target.Position.X, target.Position.Y)
	if dist <= atkRange && !ctx.HasActed() && draft.Action == nil {
		if self.IsBackstabbing(*target) {
			draft.Action = &behavior.ActionSlot{Type: behavior.ActionAttack, Target: target.Position}
		}
		return
	}

	if draft.Move != nil || ctx.RemainingMovement() <= 0 {
		return
	}

	flankDest := rearPosition(target)
	grd := ctx.Grid()
	if grd == nil {
		return
	}
	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	path, found := grd.AStarPath(self.Position, flankDest, jumpHeight, nil)
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

var _ behavior.Behavior = (*Flank)(nil)
