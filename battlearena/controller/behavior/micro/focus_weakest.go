package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
)

// FocusWeakest proposes the lowest-HP enemy as the primary target.
// Used by the Sneak and Fighter archetypes to finish off damaged enemies quickly.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type FocusWeakest struct{}

// Name returns the stable identifier used in memory records and logs.
func (*FocusWeakest) Name() string             {
	return "focus_weakest"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*FocusWeakest) BaseActivation() float64  {
	return 0.7
}

// Propose sets the Target slot to the living enemy with the lowest HP.
func (b *FocusWeakest) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Target != nil {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	target := lowestHPEnemy(self, selfTeam, entities)
	if target == nil {
		return
	}
	draft.Target = &behavior.TargetSlot{EntityID: target.ID}
}

var _ behavior.Behavior = (*FocusWeakest)(nil)
