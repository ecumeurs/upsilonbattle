package battletest_test

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/battletest"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// The reposition trap matrix is the acceptance criteria for movement skills: tiles flown over
// must NOT fire positional effects; only the landing tile fires (OnEnter).
// @test-link [[mech_movement_reposition]]

// TestRepositionDashOverTrap: a self-dash flies over an intermediate trap without triggering it.
func TestRepositionDashOverTrap(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3), battletest.WithMovement(3))

	// Non-removing OnEnter poison traps on the flown-over tiles.
	s.Trap(position.New(5, 6, 3), battletest.PoisonTrap(10, property.TriggerOnEnter, false))
	s.Trap(position.New(5, 7, 3), battletest.PoisonTrap(10, property.TriggerOnEnter, false))

	dash := caster.Give(battletest.NewSkill("Dash").
		TargetType(def.TargetTypeTile).Range(1, 3).
		Reposition(def.RepositionSubjectSelf, 3).Build())

	s.Turn(caster)
	reply, _, _ := s.UseSkill(caster, dash, position.New(5, 8, 3))

	if reply.HasError {
		t.Fatalf("expected no error, got %q", reply.ErrorKey)
	}
	if got := s.Pos(caster); !got.Equals(position.New(5, 8, 3)) {
		t.Errorf("expected caster at (5,8,3), got %s", got)
	}
	if got := s.Poison(caster); got != 0 {
		t.Errorf("expected no poison (flew over traps), got %d", got)
	}
}

// TestRepositionDashOntoTrap: a self-dash that lands on a trap triggers it.
func TestRepositionDashOntoTrap(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3), battletest.WithMovement(3))
	s.Trap(position.New(5, 8, 3), battletest.PoisonTrap(10, property.TriggerOnEnter, true))

	dash := caster.Give(battletest.NewSkill("Dash").
		TargetType(def.TargetTypeTile).Range(1, 3).
		Reposition(def.RepositionSubjectSelf, 3).Build())

	s.Turn(caster)
	s.UseSkill(caster, dash, position.New(5, 8, 3))

	if got := s.Poison(caster); got != 10 {
		t.Errorf("expected poison 10 (landed on trap), got %d", got)
	}
	if s.TrapAt(position.New(5, 8, 3)) {
		t.Errorf("expected landing trap to be consumed")
	}
}

// TestRepositionPushOverTrap: a push displaces the target over an intermediate trap (no trigger).
func TestRepositionPushOverTrap(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3))
	foe := s.Place(2, position.New(5, 6, 3), battletest.WithHP(50))
	s.Trap(position.New(5, 7, 3), battletest.PoisonTrap(10, property.TriggerOnEnter, false))

	push := caster.Give(battletest.NewSkill("Kick").
		TargetType(def.TargetTypeEnemyOnly).Range(1, 1).
		Reposition(def.RepositionSubjectTarget, 2).Build())

	s.Turn(caster)
	s.UseSkill(caster, push, position.New(5, 6, 3))

	if got := s.Pos(foe); !got.Equals(position.New(5, 8, 3)) {
		t.Errorf("expected foe pushed to (5,8,3), got %s", got)
	}
	if got := s.Poison(foe); got != 0 {
		t.Errorf("expected foe unpoisoned (flew over trap), got %d", got)
	}
}

// TestRepositionPushOntoTrap: a push that lands the target on a trap triggers it.
func TestRepositionPushOntoTrap(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3))
	foe := s.Place(2, position.New(5, 6, 3), battletest.WithHP(50))
	s.Trap(position.New(5, 8, 3), battletest.PoisonTrap(10, property.TriggerOnEnter, false))

	push := caster.Give(battletest.NewSkill("Kick").
		TargetType(def.TargetTypeEnemyOnly).Range(1, 1).
		Reposition(def.RepositionSubjectTarget, 2).Build())

	s.Turn(caster)
	s.UseSkill(caster, push, position.New(5, 6, 3))

	if got := s.Poison(foe); got != 10 {
		t.Errorf("expected foe poisoned (landed on trap), got %d", got)
	}
}

// TestRepositionPull: a negative distance pulls the target toward the caster.
func TestRepositionPull(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3))
	foe := s.Place(2, position.New(5, 9, 3), battletest.WithHP(50))

	pull := caster.Give(battletest.NewSkill("Grapple").
		TargetType(def.TargetTypeEnemyOnly).Range(1, 4).
		Reposition(def.RepositionSubjectTarget, -2).Build())

	s.Turn(caster)
	s.UseSkill(caster, pull, position.New(5, 9, 3))

	if got := s.Pos(foe); !got.Equals(position.New(5, 7, 3)) {
		t.Errorf("expected foe pulled to (5,7,3), got %s", got)
	}
}

// TestRepositionRetreatWithBuff: a self-dash composes with a shield buff applied to the caster.
func TestRepositionRetreatWithBuff(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3), battletest.WithMovement(3))

	retreat := caster.Give(battletest.NewSkill("Retreat").
		TargetType(def.TargetTypeTile).Range(1, 3).
		ShieldPower(5).
		Reposition(def.RepositionSubjectSelf, 2).Build())

	s.Turn(caster)
	s.UseSkill(caster, retreat, position.New(5, 3, 3))

	if got := s.Pos(caster); !got.Equals(position.New(5, 3, 3)) {
		t.Errorf("expected caster retreated to (5,3,3), got %s", got)
	}
	if got := s.Shield(caster); got != 5 {
		t.Errorf("expected shield 5 from composed buff, got %d", got)
	}
}

// TestRepositionMvtCost: a movement skill consumes movement points.
func TestRepositionMvtCost(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3), battletest.WithMovement(3))

	dash := caster.Give(battletest.NewSkill("Dash").
		TargetType(def.TargetTypeTile).Range(1, 3).
		MvtCost(1).
		Reposition(def.RepositionSubjectSelf, 2).Build())

	s.Turn(caster)
	s.UseSkill(caster, dash, position.New(5, 7, 3))

	if got := s.Movement(caster); got != 2 {
		t.Errorf("expected 2 movement left (3 - 1 cost), got %d", got)
	}
}

// TestRepositionBlockedLanding: landing on a blocked or out-of-grid tile fails with the right key.
func TestRepositionBlockedLanding(t *testing.T) {
	cases := []struct {
		name     string
		casterAt position.Position
		blockAt  *position.Position
		aim      position.Position
		dist     int
		errKey   string
	}{
		{
			name:     "occupied",
			casterAt: position.New(5, 5, 3),
			blockAt:  ptr(position.New(5, 7, 3)),
			aim:      position.New(5, 8, 3),
			dist:     2,
			errKey:   "skill.reposition.blocked",
		},
		{
			name:     "out of grid",
			casterAt: position.New(8, 5, 3),
			aim:      position.New(9, 5, 3),
			dist:     3,
			errKey:   "skill.reposition.outofgrid",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := battletest.New(t, 10, 10, 3)
			caster := s.Place(1, tc.casterAt, battletest.WithMovement(3))
			if tc.blockAt != nil {
				s.Place(2, *tc.blockAt, battletest.WithHP(10))
			}
			dash := caster.Give(battletest.NewSkill("Dash").
				TargetType(def.TargetTypeTile).Range(1, 3).
				Reposition(def.RepositionSubjectSelf, tc.dist).Build())

			s.Turn(caster)
			reply, _, _ := s.UseSkill(caster, dash, tc.aim)

			if !reply.HasError {
				t.Fatalf("expected error, got success")
			}
			if reply.ErrorKey != tc.errKey {
				t.Errorf("expected error key %q, got %q", tc.errKey, reply.ErrorKey)
			}
			if got := s.Pos(caster); !got.Equals(tc.casterAt) {
				t.Errorf("expected caster to stay at %s on blocked reposition, got %s", tc.casterAt, got)
			}
		})
	}
}

func ptr(p position.Position) *position.Position { return &p }
