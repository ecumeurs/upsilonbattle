package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
)

const resourceThresholdPct = 25

// PreserveResource vetoes skill usage when SP or MP is critically low (< 25% of max).
// It does this by writing a Pass action to the Action slot, overriding any queued skill.
// Used across archetypes that rely heavily on skills.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type PreserveResource struct{}

// Name returns the stable identifier used in memory records and logs.
func (*PreserveResource) Name() string             {
	return "preserve_resource"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*PreserveResource) BaseActivation() float64  {
	return 0.7
}

// Propose blocks the Action slot to suppress skill use when SP/MP is critically low.
func (b *PreserveResource) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if ctx.HasActed() {
		return
	}
	// Only relevant when a skill action has already been proposed by a prior layer.
	if draft.Action == nil || draft.Action.Type != behavior.ActionSkill {
		return
	}

	self := ctx.SelfEntity()
	sp := self.GetPropertyC(property.SP).GetValue()
	maxSP := self.GetPropertyC(property.SP).GetMaxValue()
	mp := self.GetPropertyC(property.MP).GetValue()
	maxMP := self.GetPropertyC(property.MP).GetMaxValue()

	spLow := maxSP > 0 && sp*100/maxSP < resourceThresholdPct
	mpLow := maxMP > 0 && mp*100/maxMP < resourceThresholdPct

	if spLow || mpLow {
		// Clear the skill action — let the baseline fall back to a basic attack.
		draft.Action = nil
	}
}

var _ behavior.Behavior = (*PreserveResource)(nil)
