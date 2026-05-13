package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
)

// BattleFocus locks onto the current sticky target from memory and maintains
// focus across turns, even if a different enemy is technically closer.
// Used by the Fighter archetype to prevent target-switching mid-combat.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type BattleFocus struct{}

// Name returns the stable identifier used in memory records and logs.
func (*BattleFocus) Name() string             {
	return "battle_focus"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*BattleFocus) BaseActivation() float64  {
	return 0.8
}

// Propose locks the Target slot onto the current combat target for sustained focus.
func (b *BattleFocus) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Target != nil {
		return
	}
	stickyID := ctx.Memory().CurrentTarget()
	if stickyID == [16]byte{} {
		return
	}
	entities := ctx.KnownEntities()
	target, ok := entities[stickyID]
	if !ok || target.GetPropertyI(property.HP).I() <= 0 {
		ctx.Memory().ClearTarget()
		return
	}
	draft.Target = &behavior.TargetSlot{EntityID: stickyID}
}

var _ behavior.Behavior = (*BattleFocus)(nil)
