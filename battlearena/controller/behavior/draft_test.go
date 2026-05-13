package behavior

// @test-link [[mechanic_mech_behavior_layered]]

import (
	"testing"

	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// TestResolveEmptyDraftEndsTheTurn confirms that an empty draft always resolves to EndOfTurn.
func TestResolveEmptyDraftEndsTheTurn(t *testing.T) {
	ctx := newCtx()
	d := &DecisionDraft{Notes: make(map[string]any)}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdEndOfTurn {
		t.Errorf("empty draft: got %v, want CmdEndOfTurn", cmd.Type)
	}
}

// TestResolveActionProducesAttack confirms that a set Action slot produces CmdAttack when the entity has not yet acted.
func TestResolveActionProducesAttack(t *testing.T) {
	ctx := newCtx()
	target := position.New(1, 0, 1)
	d := &DecisionDraft{
		Action: &ActionSlot{Type: ActionAttack, Target: target},
		Notes:  make(map[string]any),
	}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdAttack {
		t.Errorf("got %v, want CmdAttack", cmd.Type)
	}
	if cmd.Target != target {
		t.Errorf("wrong target position")
	}
}

// TestResolveSkillProducesCmdSkill confirms that ActionSkill in the Action slot yields CmdSkill with the correct SkillID.
func TestResolveSkillProducesCmdSkill(t *testing.T) {
	ctx := newCtx()
	skillID := uuid.New()
	target := position.New(2, 0, 1)
	d := &DecisionDraft{
		Action: &ActionSlot{Type: ActionSkill, Target: target, SkillID: skillID},
		Notes:  make(map[string]any),
	}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdSkill {
		t.Errorf("got %v, want CmdSkill", cmd.Type)
	}
	if cmd.SkillID != skillID {
		t.Errorf("skill ID mismatch")
	}
}

// TestResolveActionBlockedWhenAlreadyActed confirms that when HasActed is true,
// the Action slot is ignored and the pipeline falls through to Move.
func TestResolveActionBlockedWhenAlreadyActed(t *testing.T) {
	ctx := newCtx()
	ctx.hasActed = true
	path := []position.Position{position.New(0, 0, 1), position.New(1, 0, 1)}
	d := &DecisionDraft{
		Action: &ActionSlot{Type: ActionAttack, Target: position.New(3, 0, 1)},
		Move:   &MoveSlot{Path: path},
		Notes:  make(map[string]any),
	}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdMove {
		t.Errorf("acted entity: got %v, want CmdMove (action should be skipped)", cmd.Type)
	}
}

// TestResolveMoveProducesCmdMove confirms that a Move slot with remaining budget resolves to CmdMove.
func TestResolveMoveProducesCmdMove(t *testing.T) {
	ctx := newCtx()
	ctx.movement = 2
	path := []position.Position{position.New(0, 0, 1), position.New(1, 0, 1)}
	d := &DecisionDraft{
		Move:  &MoveSlot{Path: path},
		Notes: make(map[string]any),
	}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdMove {
		t.Errorf("got %v, want CmdMove", cmd.Type)
	}
	if len(cmd.Path) != len(path) {
		t.Errorf("path length mismatch: got %d, want %d", len(cmd.Path), len(path))
	}
}

// TestResolveMoveBlockedByNoBudget confirms that a Move slot is ignored when RemainingMovement is 0.
func TestResolveMoveBlockedByNoBudget(t *testing.T) {
	ctx := newCtx()
	ctx.movement = 0
	path := []position.Position{position.New(0, 0, 1), position.New(1, 0, 1)}
	d := &DecisionDraft{
		Move:  &MoveSlot{Path: path},
		Notes: make(map[string]any),
	}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdEndOfTurn {
		t.Errorf("zero movement: got %v, want CmdEndOfTurn", cmd.Type)
	}
}

// TestResolveMoveBlockedByEmptyPath confirms that a Move slot with an empty path falls through to EndOfTurn.
func TestResolveMoveBlockedByEmptyPath(t *testing.T) {
	ctx := newCtx()
	d := &DecisionDraft{
		Move:  &MoveSlot{Path: nil},
		Notes: make(map[string]any),
	}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdEndOfTurn {
		t.Errorf("empty path: got %v, want CmdEndOfTurn", cmd.Type)
	}
}

// TestResolveActionWinsOverMoveWhenNotActed confirms that Action takes priority over Move.
func TestResolveActionWinsOverMoveWhenNotActed(t *testing.T) {
	ctx := newCtx()
	ctx.hasActed = false
	ctx.movement = 3
	path := []position.Position{position.New(0, 0, 1), position.New(1, 0, 1)}
	d := &DecisionDraft{
		Action: &ActionSlot{Type: ActionAttack, Target: position.New(2, 0, 1)},
		Move:   &MoveSlot{Path: path},
		Notes:  make(map[string]any),
	}
	cmd := d.Resolve(ctx)
	if cmd.Type != CmdAttack {
		t.Errorf("action+move set, not acted: got %v, want CmdAttack", cmd.Type)
	}
}
