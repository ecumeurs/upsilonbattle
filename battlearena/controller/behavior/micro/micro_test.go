package micro_test

// @test-link [[mechanic_ai_controller_archetypes]]

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior/micro"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/google/uuid"
)

// ── micro test harness ────────────────────────────────────────────────────────

type testCtx struct {
	self     entity.Entity
	known    map[uuid.UUID]entity.Entity
	grd      *grid.Grid
	hasActed bool
	movement int
	mem      *behavior.DecisionMemory
	grade    int
}

func (c *testCtx) SelfEntity() entity.Entity                  { return c.self }
func (c *testCtx) KnownEntities() map[uuid.UUID]entity.Entity { return c.known }
func (c *testCtx) Grid() *grid.Grid                           { return c.grd }
func (c *testCtx) HasActed() bool                             { return c.hasActed }
func (c *testCtx) RemainingMovement() int                     { return c.movement }
func (c *testCtx) LastTickOutcome() behavior.TickOutcome       { return behavior.TickNone }
func (c *testCtx) Memory() *behavior.DecisionMemory           { return c.mem }
func (c *testCtx) Grade() int                                  { return c.grade }

func newTestCtx(self entity.Entity) *testCtx {
	return &testCtx{
		self:     self,
		known:    map[uuid.UUID]entity.Entity{self.ID: self},
		mem:      behavior.NewDecisionMemory(),
		movement: 3,
		grade:    8,
	}
}

func makeEnt(x, y, teamID, hp, maxHP int) entity.Entity {
	e := entity.Entity{
		ID:       uuid.New(),
		Position: position.New(x, y, 1),
	}
	e.Properties = make(map[string]property.Property)
	e.RepsertPropertyValue(property.TeamID, teamID)
	e.RepsertPropertyCMaxValue(property.HP, maxHP)
	e.RepsertPropertyCValue(property.HP, hp)
	return e
}

func emptyDraft() *behavior.DecisionDraft {
	return &behavior.DecisionDraft{Notes: make(map[string]any)}
}

// ── FocusWeakest ──────────────────────────────────────────────────────────────

// TestFocusWeakestChoosesLowestHPEnemy verifies that FocusWeakest selects the enemy
// with the lowest current HP when multiple enemies are present.
func TestFocusWeakestChoosesLowestHPEnemy(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	ctx := newTestCtx(self)

	strong := makeEnt(1, 0, 2, 8, 10) // 8 HP
	weak := makeEnt(2, 0, 2, 3, 10)   // 3 HP — should be selected
	ctx.known[strong.ID] = strong
	ctx.known[weak.ID] = weak

	d := emptyDraft()
	(&micro.FocusWeakest{}).Propose(ctx, d)

	if d.Target == nil {
		t.Fatal("FocusWeakest did not set Target")
	}
	if d.Target.EntityID != weak.ID {
		t.Errorf("FocusWeakest selected wrong target: got %v, want weakest %v", d.Target.EntityID, weak.ID)
	}
}

// TestFocusWeakestDoesNotOverwriteExistingTarget confirms first-writer-wins: if Target is already set,
// FocusWeakest does not replace it.
func TestFocusWeakestDoesNotOverwriteExistingTarget(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	ctx := newTestCtx(self)
	enemy := makeEnt(1, 0, 2, 2, 10)
	ctx.known[enemy.ID] = enemy

	priorID := uuid.New()
	d := emptyDraft()
	d.Target = &behavior.TargetSlot{EntityID: priorID}

	(&micro.FocusWeakest{}).Propose(ctx, d)

	if d.Target.EntityID != priorID {
		t.Error("FocusWeakest overwrote an existing Target (violates first-writer-wins)")
	}
}

// TestFocusWeakestNoEnemiesLeavesTargetNil confirms that with no enemies, Target is not set.
func TestFocusWeakestNoEnemiesLeavesTargetNil(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	ctx := newTestCtx(self)
	// Only self in known entities (same team)

	d := emptyDraft()
	(&micro.FocusWeakest{}).Propose(ctx, d)

	if d.Target != nil {
		t.Error("FocusWeakest set Target despite having no enemies")
	}
}

// ── KiteAway ─────────────────────────────────────────────────────────────────

// TestKiteAwayProposesMoveWhenFoeInMeleeRange verifies that KiteAway proposes a move
// when an enemy is within the entity's attack range and movement budget is available.
func TestKiteAwayProposesMoveWhenFoeInMeleeRange(t *testing.T) {
	self := makeEnt(3, 3, 1, 10, 10)
	self.RepsertPropertyValue(property.AttackRange, 2)
	self.RepsertPropertyValue(property.JumpHeight, 1)
	ctx := newTestCtx(self)

	// Enemy is adjacent (distance 1, within attack range 2).
	foe := makeEnt(4, 3, 2, 10, 10)
	ctx.known[foe.ID] = foe

	// Provide a large enough grid so A* can find a flee path.
	ctx.grd = grid.NewGrid(10, 10, 1)

	d := emptyDraft()
	(&micro.KiteAway{}).Propose(ctx, d)

	if d.Move == nil {
		t.Error("KiteAway did not propose a Move despite foe within attack range")
	}
}

// TestKiteAwayDoesNotMoveWhenFoeOutOfRange confirms that KiteAway does not propose a move
// when the enemy is already beyond the entity's attack range.
func TestKiteAwayDoesNotMoveWhenFoeOutOfRange(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	self.RepsertPropertyValue(property.AttackRange, 1)
	self.RepsertPropertyValue(property.JumpHeight, 1)
	ctx := newTestCtx(self)

	// Enemy is 5 cells away — beyond attack range 1.
	foe := makeEnt(5, 0, 2, 10, 10)
	ctx.known[foe.ID] = foe
	ctx.grd = grid.NewGrid(10, 10, 1)

	d := emptyDraft()
	(&micro.KiteAway{}).Propose(ctx, d)

	if d.Move != nil {
		t.Error("KiteAway proposed a Move even though foe is out of range")
	}
}

// TestKiteAwayDoesNotOverwriteExistingMove confirms first-writer-wins for the Move slot.
func TestKiteAwayDoesNotOverwriteExistingMove(t *testing.T) {
	self := makeEnt(3, 3, 1, 10, 10)
	self.RepsertPropertyValue(property.AttackRange, 2)
	self.RepsertPropertyValue(property.JumpHeight, 1)
	ctx := newTestCtx(self)

	foe := makeEnt(4, 3, 2, 10, 10)
	ctx.known[foe.ID] = foe
	ctx.grd = grid.NewGrid(10, 10, 1)

	priorPath := []position.Position{position.New(3, 3, 1), position.New(3, 4, 1)}
	d := emptyDraft()
	d.Move = &behavior.MoveSlot{Path: priorPath}

	(&micro.KiteAway{}).Propose(ctx, d)

	if len(d.Move.Path) != len(priorPath) {
		t.Error("KiteAway overwrote an existing Move (violates first-writer-wins)")
	}
}

// ── HealAlly ─────────────────────────────────────────────────────────────────

// TestHealAllyNoSkillDoesNothing confirms that HealAlly produces no action when the entity
// has no heal skill equipped.
func TestHealAllyNoSkillDoesNothing(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	ctx := newTestCtx(self)

	// Ally at low HP (below 50%).
	ally := makeEnt(1, 0, 1, 3, 10)
	ctx.known[ally.ID] = ally

	d := emptyDraft()
	(&micro.HealAlly{}).Propose(ctx, d)

	if d.Action != nil || d.Target != nil {
		t.Error("HealAlly proposed an action despite having no heal skill")
	}
}

// TestHealAllySkipsWhenAlreadyActed confirms that HealAlly does not propose when the entity
// has already acted this turn.
func TestHealAllySkipsWhenAlreadyActed(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	ctx := newTestCtx(self)
	ctx.hasActed = true

	ally := makeEnt(1, 0, 1, 2, 10)
	ctx.known[ally.ID] = ally

	d := emptyDraft()
	(&micro.HealAlly{}).Propose(ctx, d)

	if d.Action != nil {
		t.Error("HealAlly proposed an action after the entity already acted")
	}
}

// TestHealAllySkipsWhenAllyAboveThreshold confirms that HealAlly does not act when the
// wounded ally is above 50% HP.
func TestHealAllySkipsWhenAllyAboveThreshold(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	ctx := newTestCtx(self)

	// Ally at 60% HP — above the 50% threshold.
	ally := makeEnt(1, 0, 1, 6, 10)
	ctx.known[ally.ID] = ally

	d := emptyDraft()
	(&micro.HealAlly{}).Propose(ctx, d)

	if d.Action != nil {
		t.Error("HealAlly proposed heal for ally above 50% HP threshold")
	}
}

// ── Ambush ────────────────────────────────────────────────────────────────────

// TestAmbushRespectsMemoryCooldown confirms that Ambush does not propose when it was
// used recently (within the cooldown window).
func TestAmbushRespectsMemoryCooldown(t *testing.T) {
	self := makeEnt(0, 0, 1, 10, 10)
	self.RepsertPropertyValue(property.AttackRange, 1)
	self.RepsertPropertyValue(property.JumpHeight, 1)
	ctx := newTestCtx(self)

	foe := makeEnt(5, 0, 2, 10, 10)
	ctx.known[foe.ID] = foe
	ctx.grd = grid.NewGrid(10, 10, 1)

	// Record ambush as having just fired on the current turn.
	draft := &behavior.DecisionDraft{Notes: make(map[string]any)}
	ctx.mem.Record("ambush", draft)
	// TurnsSince("ambush") == 0 → within 3-turn cooldown → should skip.

	d := emptyDraft()
	(&micro.Ambush{}).Propose(ctx, d)

	if d.Move != nil || d.Action != nil {
		t.Error("Ambush proposed despite cooldown not yet expired")
	}
}
