package rules

import (
	"testing"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// TestVersioningAudit verifies that the GameState version is correctly incremented 
// and propagated in notifications for all major controller actions (Move, Attack, Pass).
// @test-link [[mech_game_state_versioning]]
// @test-link [[mech_version_bit_packing]]
func TestVersioningAudit(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	// Baseline
	v0 := gs.Version

	// 1. Move
	moveMsg := message.Create(nil, rulermethods.ControllerMove{
		EntityID:     fake.Entity1,
		ControllerID: fake.Controller1,
		Path:         []position.Position{{X: 0, Y: 1, Z: 3}},
	}, nil)
	Move(gs, moveMsg, moveMsg.TargetMethod.(rulermethods.ControllerMove))

	if gs.Version <= v0 {
		t.Errorf("Version not bumped after Move. Got %d, want > %d", gs.Version, v0)
	}

	// Check notifications for Move
	foundMoveNotif := false
	for _, m := range fake.FakeController1.NotifyMessages {
		if notif, ok := m.TargetMethod.(rulermethods.ControllerMoved); ok {
			foundMoveNotif = true
			if notif.Version != gs.Version {
				t.Errorf("Move notification has wrong version. Got %d, want %d", notif.Version, gs.Version)
			}
		}
	}
	if !foundMoveNotif {
		t.Errorf("Move notification not found")
	}

	// 2. Attack
	v1 := gs.Version
	gs.Turner.ForceTurn(fake.Entity1)
	ent1 := gs.Entities[fake.Entity1]
	ent1.UpdatePropertyValue(property.HasMoved, false)
	ent1.UpdatePropertyValue(property.HasActed, false)
	// Move ent1 back to 0,0,3 to be adjacent to ent2 at 1,0,3
	gs.Grid.MoveEntity(ent1.Position, position.Position{X: 0, Y: 0, Z: 3}, ent1.ID)
	ent1.Position = position.Position{X: 0, Y: 0, Z: 3}
	gs.Entities[fake.Entity1] = ent1

	ent2 := gs.Entities[fake.Entity2]
	ent2.UpdatePropertyValue(property.TeamID, 2)
	gs.Entities[fake.Entity2] = ent2

	attackMsg := message.Create(nil, rulermethods.ControllerAttack{
		EntityID:     fake.Entity1,
		ControllerID: fake.Controller1,
		Target:       position.Position{X: 1, Y: 0, Z: 3},
	}, nil)
	Attack(gs, attackMsg, attackMsg.TargetMethod.(rulermethods.ControllerAttack))

	if gs.Version <= v1 {
		t.Errorf("Version not bumped after Attack. Got %d, want > %d", gs.Version, v1)
	}

	// Check notifications for Attack
	foundAttackNotif := false
	for _, m := range fake.FakeController1.NotifyMessages {
		if notif, ok := m.TargetMethod.(rulermethods.ControllerAttacked); ok {
			foundAttackNotif = true
			if notif.Version != gs.Version {
				t.Errorf("Attack notification has wrong version. Got %d, want %d", notif.Version, gs.Version)
			}
		}
	}
	if !foundAttackNotif {
		t.Errorf("Attack notification not found")
	}

	// 3. Pass (EndOfTurn)
	v2 := gs.Version
	gs.Turner.ForceTurn(fake.Entity1)
	passMsg := message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     fake.Entity1,
		ControllerID: fake.Controller1,
	}, nil)
	ent1 = gs.Entities[fake.Entity1]
	EndOfTurn(gs, passMsg, passMsg.TargetMethod.(rulermethods.EndOfTurn), ent1)

	if gs.Version <= v2 {
		t.Errorf("Version not bumped after Pass. Got %d, want > %d", gs.Version, v2)
	}

	// Check notifications for Pass
	foundPassNotif := false
	for _, m := range fake.FakeController1.NotifyMessages {
		if notif, ok := m.TargetMethod.(rulermethods.ControllerPassed); ok {
			foundPassNotif = true
			if notif.Version != gs.Version {
				t.Errorf("Pass notification has wrong version. Got %d, want %d", notif.Version, gs.Version)
			}
		}
	}
	if !foundPassNotif {
		t.Errorf("Pass notification not found")
	}
}
