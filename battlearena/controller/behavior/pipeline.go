package behavior

import "math/rand"

// LayeredBehavior orchestrates an ordered stack of Behavior layers.
//
// Layers are evaluated top to bottom. Each layer rolls its activation rate against
// a grade-scaled threshold before being allowed to Propose. If the roll fails, the
// layer is skipped and a RecordSkipped entry is written to memory.
//
// First-writer-wins per slot: each layer only writes to an unset slot.
//
// The last layer in the stack must be the always-active baseline (AggressiveBehavior,
// BaseActivation = 1.0) to guarantee a decision is always reached.
//
// @spec-link [[mechanic_behavior_layered]]
type LayeredBehavior struct {
	Layers []Behavior
}

// NewLayeredBehavior creates a LayeredBehavior with the given ordered stack.
// The baseline (AggressiveBehavior) must be the last element.
func NewLayeredBehavior(layers ...Behavior) *LayeredBehavior {
	return &LayeredBehavior{Layers: layers}
}

// Tick runs one pipeline iteration and returns the single EngineCommand to execute.
// It advances the memory tick counter before running layers.
func (l *LayeredBehavior) Tick(ctx GameContext) EngineCommand {
	ctx.Memory().AdvanceTick()
	draft := &DecisionDraft{Notes: make(map[string]any)}
	for _, layer := range l.Layers {
		if !rollActivation(layer.BaseActivation(), ctx.Grade()) {
			ctx.Memory().RecordSkipped(layer.Name())
			continue
		}
		layer.Propose(ctx, draft)
		ctx.Memory().Record(layer.Name(), draft)
	}
	return draft.Resolve(ctx)
}

// effectiveActivation scales a layer's declared rate by the entity's grade index.
//
//	Grade I  (idx 0) → 0.40 × base
//	Grade III (idx 4) → 0.70 × base
//	Grade V  (idx 8) → 1.00 × base
//
// Clamped to [0, 1].
func effectiveActivation(base float64, gradeIdx int) float64 {
	rate := base * (0.4 + float64(gradeIdx)*0.075)
	if rate > 1.0 {
		return 1.0
	}
	if rate < 0 {
		return 0
	}
	return rate
}

// rollActivation returns true when the layer should run this tick.
// A base of 1.0 bypasses the grade-scaling roll — the layer always activates.
// This is the contract for the always-active baseline (AggressiveBehavior).
func rollActivation(base float64, gradeIdx int) bool {
	if base >= 1.0 {
		return true
	}
	return rand.Float64() < effectiveActivation(base, gradeIdx)
}
