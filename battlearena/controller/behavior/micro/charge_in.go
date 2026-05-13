package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontools/tools"
)

// ChargeIn uses all available movement to close distance to the target as fast as
// possible, ignoring attack-range caution (melee focus). It proposes a Move that
// uses the full movement budget right up to the target.
// Used by the Fighter archetype.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type ChargeIn struct{}

// Name returns the stable identifier used in memory records and logs.
func (*ChargeIn) Name() string             {
	return "charge_in"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*ChargeIn) BaseActivation() float64  {
	return 0.7
}

// Propose moves the entity directly toward the nearest foe using the fastest available path.
func (b *ChargeIn) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Move != nil || ctx.RemainingMovement() <= 0 {
		return
	}

	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()
	atkRange := self.GetPropertyI(property.AttackRange).I()

	// Use draft target if set, else nearest foe.
	foe := nearestLivingEnemy(self, selfTeam, entities)
	if draft.Target != nil {
		if ent, ok := entities[draft.Target.EntityID]; ok && ent.GetPropertyI(property.HP).I() > 0 {
			e := ent
			foe = &e
		}
	}
	if foe == nil {
		return
	}

	dist := tools.Distance(self.Position.X, self.Position.Y, foe.Position.X, foe.Position.Y)
	if dist <= atkRange {
		return // already in range — nothing to charge toward
	}

	grd := ctx.Grid()
	if grd == nil {
		return
	}
	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	path, found := grd.AStarPath(self.Position, foe.Position, jumpHeight, nil)
	if !found || len(path) <= 1 {
		return
	}

	// Charge: use full movement budget (stop 1 cell before the target to allow attack).
	limit := len(path) - 1
	if limit > ctx.RemainingMovement() {
		limit = ctx.RemainingMovement()
	}
	if limit < 2 {
		return
	}
	draft.Move = &behavior.MoveSlot{Path: path[1:limit]}
}

var _ behavior.Behavior = (*ChargeIn)(nil)
