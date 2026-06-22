package battletest_test

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/battletest"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// TestSandboxDamageSkill re-expresses a basic single-target damage skill through the
// scenario sandbox to validate ergonomics.
// @test-link [[module_skill_sandbox]]
func TestSandboxDamageSkill(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)

	attacker := s.Place(1, position.New(5, 5, 3), battletest.WithStat(property.Attack, 5))
	foe := s.Place(2, position.New(5, 6, 3), battletest.WithHP(20), battletest.WithStat(property.Defense, 0))

	strike := attacker.Give(
		battletest.NewSkill("Strike").TargetType("Entity").Range(1, 1).Damage(100).Build(),
	)

	s.Turn(attacker)
	reply, damaged, _ := s.UseSkill(attacker, strike, s.Pos(foe))

	if reply.HasError {
		t.Fatalf("expected no error, got %q", reply.ErrorKey)
	}
	if got := s.HP(foe); got != 15 {
		t.Errorf("expected foe HP 15 (20 - 5), got %d", got)
	}
	if len(damaged) != 1 {
		t.Errorf("expected 1 damaged entity, got %d", len(damaged))
	}
	if !s.HasActed(attacker) {
		t.Errorf("expected attacker to have acted")
	}
}

// TestSandboxTrapOnEnter re-expresses an OnEnter poison trap triggered by movement.
// @test-link [[module_skill_sandbox]]
func TestSandboxTrapOnEnter(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)

	walker := s.Place(1, position.New(0, 0, 3))
	trapPos := position.New(0, 1, 3)
	s.Trap(trapPos, battletest.PoisonTrap(10, property.TriggerOnEnter, true))

	s.Turn(walker)
	s.Move(walker, trapPos)

	if got := s.Poison(walker); got != 10 {
		t.Errorf("expected poison 10 after stepping on trap, got %d", got)
	}
	if s.TrapAt(trapPos) {
		t.Errorf("expected RemoveOnTrigger trap to be consumed")
	}
}
