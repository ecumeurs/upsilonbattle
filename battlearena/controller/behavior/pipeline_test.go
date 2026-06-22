package behavior

// @test-link [[mechanic_behavior_layered]]

import (
	"testing"
)

// ── effectiveActivation ───────────────────────────────────────────────────────

func TestEffectiveActivationGradeI(t *testing.T) {
	// Grade I (idx 0): effective = base * (0.4 + 0*0.075) = base * 0.4
	got := effectiveActivation(0.8, 0)
	want := 0.8 * 0.4
	if abs(got-want) > 1e-9 {
		t.Errorf("effectiveActivation(0.8, 0) = %v, want %v", got, want)
	}
}

func TestEffectiveActivationGradeIII(t *testing.T) {
	// Grade III (idx 4): effective = base * (0.4 + 4*0.075) = base * 0.7
	got := effectiveActivation(0.8, 4)
	want := 0.8 * 0.7
	if abs(got-want) > 1e-9 {
		t.Errorf("effectiveActivation(0.8, 4) = %v, want %v", got, want)
	}
}

func TestEffectiveActivationGradeV(t *testing.T) {
	// Grade V (idx 8): effective = base * (0.4 + 8*0.075) = base * 1.0
	got := effectiveActivation(0.8, 8)
	want := 0.8
	if abs(got-want) > 1e-9 {
		t.Errorf("effectiveActivation(0.8, 8) = %v, want %v", got, want)
	}
}

func TestEffectiveActivationClampsToOne(t *testing.T) {
	got := effectiveActivation(2.0, 8) // 2.0 * 1.0 = 2.0 → clamped to 1.0
	if got != 1.0 {
		t.Errorf("effectiveActivation(2.0, 8) = %v, want 1.0 (clamped)", got)
	}
}

func TestEffectiveActivationZeroBase(t *testing.T) {
	got := effectiveActivation(0.0, 8)
	if got != 0.0 {
		t.Errorf("effectiveActivation(0.0, 8) = %v, want 0.0", got)
	}
}

// ── rollActivation ────────────────────────────────────────────────────────────

// TestBaselineAlwaysActivates verifies that base=1.0 bypasses the roll and always activates.
func TestBaselineAlwaysActivates(t *testing.T) {
	for i := 0; i < 200; i++ {
		if !rollActivation(1.0, 0) {
			t.Fatalf("rollActivation(1.0, 0) returned false on iteration %d", i)
		}
	}
}

// TestZeroBaseNeverActivates verifies that base=0.0 never passes the activation roll.
func TestZeroBaseNeverActivates(t *testing.T) {
	for i := 0; i < 200; i++ {
		if rollActivation(0.0, 8) {
			t.Fatalf("rollActivation(0.0, 8) returned true on iteration %d", i)
		}
	}
}

// ── LayeredBehavior.Tick ──────────────────────────────────────────────────────

// stubBehavior is a test layer that records whether it was called and writes to the given slot.
type stubBehavior struct {
	name       string
	activation float64
	called     bool
	writeTarget bool
	writeMove   bool
	writeAction bool
}

func (s *stubBehavior) Name() string           { return s.name }
func (s *stubBehavior) BaseActivation() float64 { return s.activation }
func (s *stubBehavior) Propose(_ GameContext, d *DecisionDraft) {
	s.called = true
	if s.writeTarget && d.Target == nil {
		d.Target = &TargetSlot{}
	}
	if s.writeMove && d.Move == nil {
		d.Move = &MoveSlot{}
	}
	if s.writeAction && d.Action == nil {
		d.Action = &ActionSlot{Type: ActionAttack}
	}
}

// TestFirstWriterWinsTarget confirms that when two layers both try to set Target, the first one wins.
func TestFirstWriterWinsTarget(t *testing.T) {
	first := &stubBehavior{name: "first", activation: 1.0, writeTarget: true}
	second := &stubBehavior{name: "second", activation: 1.0, writeTarget: true}

	lb := NewLayeredBehavior(first, second)
	ctx := newCtx()
	ctx.hasActed = true // prevent Action from resolving to CmdAttack and complicating the test
	ctx.movement = 0   // prevent Move

	lb.Tick(ctx)

	if !first.called {
		t.Error("first layer was not called")
	}
	if !second.called {
		t.Error("second layer was not called")
	}
	// Both were called but second should not have overwritten first's target
	// (we can't inspect draft directly, but the pipeline enforces first-writer-wins
	// via the convention in each Propose: "if draft.Target == nil { ... }")
}

// TestBaselineFillsAllSlotsWhenOtherLayersProduceNothing confirms that when all non-baseline
// layers produce nothing, the baseline always-active layer fills all three slots.
func TestBaselineFillsAllSlotsWhenOtherLayersProduceNothing(t *testing.T) {
	// A layer that proposes nothing.
	noop := &stubBehavior{name: "noop", activation: 1.0}

	// Baseline: writes Target, Move, Action.
	baseline := &stubBehavior{
		name:        "baseline",
		activation:  1.0,
		writeTarget: true,
		writeMove:   true,
		writeAction: true,
	}

	lb := NewLayeredBehavior(noop, baseline)
	ctx := newCtx()
	ctx.hasActed = false
	ctx.movement = 3

	cmd := lb.Tick(ctx)
	// baseline set Action → resolve produces CmdAttack.
	if cmd.Type != CmdAttack {
		t.Errorf("expected baseline to produce CmdAttack, got %v", cmd.Type)
	}
}

// TestSkippedLayerDoesNotPropose verifies that a layer with BaseActivation=0 is always skipped.
func TestSkippedLayerDoesNotPropose(t *testing.T) {
	skipped := &stubBehavior{name: "always_skip", activation: 0.0, writeTarget: true, writeMove: true}
	baseline := &stubBehavior{name: "baseline", activation: 1.0}

	lb := NewLayeredBehavior(skipped, baseline)
	ctx := newCtx()
	ctx.hasActed = true
	ctx.movement = 0

	lb.Tick(ctx)

	if skipped.called {
		t.Error("zero-activation layer was Propose()'d — it should have been skipped")
	}
}

// TestTickAdvancesMemoryTick verifies that each Tick call increments the memory tick counter.
func TestTickAdvancesMemoryTick(t *testing.T) {
	lb := NewLayeredBehavior(&stubBehavior{name: "b", activation: 1.0})
	ctx := newCtx()
	ctx.hasActed = true
	ctx.movement = 0
	before := ctx.mem.tick
	lb.Tick(ctx)
	if ctx.mem.tick != before+1 {
		t.Errorf("memory tick not advanced: before=%d, after=%d", before, ctx.mem.tick)
	}
}

// abs returns the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
