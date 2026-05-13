package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/google/uuid"
)

// CounterCharge proposes the entity that attacked us last tick as the primary target.
// Reads the sticky target set by the ControllerAttacked handler. Only activates when
// an attacker is recorded in memory and is still alive.
// Used by the Fighter archetype to retaliate aggressively.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type CounterCharge struct{}

// Name returns the stable identifier used in memory records and logs.
func (*CounterCharge) Name() string             {
	return "counter_charge"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*CounterCharge) BaseActivation() float64  {
	return 0.7
}

// Propose targets the entity that attacked self last turn, enabling retaliation priority.
func (b *CounterCharge) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Target != nil {
		return
	}
	attackerID := ctx.Memory().CurrentTarget()
	if attackerID == uuid.Nil {
		return
	}
	entities := ctx.KnownEntities()
	attacker, ok := entities[attackerID]
	if !ok || attacker.GetPropertyI(property.HP).I() <= 0 {
		return
	}
	draft.Target = &behavior.TargetSlot{EntityID: attackerID}
}

var _ behavior.Behavior = (*CounterCharge)(nil)
