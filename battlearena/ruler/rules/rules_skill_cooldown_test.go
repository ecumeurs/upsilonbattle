package rules

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// TestRuleSkillFailCooldown ensures that a skill fails if it is currently on cooldown.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set the skill cooldown to 1.
	fake.Skill.Cooldown = 1
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill that is on cooldown.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cooldown" {
		t.Errorf("Expected error 'skill.cooldown', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillSetCooldown verifies that using a skill correctly applies its cooldown to the entity.
// @test-link [[mech_skill_validation]]
func TestRuleSkillSetCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Default cooldown cost for the test skill is 3.
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Expected no error, got '%s'", reply.ErrorKey)
	}

	// Verify that the cooldown was updated in the GameState.
	skill := gs.Entities[fake.Attacker].Skills[fake.SkillID]
	if skill.Cooldown != 3 {
		t.Errorf("Expected cooldown to be 3(in GameState), got %d", skill.Cooldown)
	}

	// Verify that the cooldown was updated in the reply message.
	skill = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.Skills[fake.SkillID]
	if skill.Cooldown != 3 {
		t.Errorf("Expected cooldown to be 3, got %d", skill.Cooldown)
	}
}

// TestRuleSkillCooldownClearsAfterElapsedTurns is the mirror image of
// TestRuleSkillFailCooldown: it asserts that a cast skill (cooldown set to its
// max value of 3) is STILL rejected on the caster's immediately-following turn,
// but becomes castable again once enough of the caster's own turns elapse to
// tick the cooldown counter back down to 0. Regression guard for ISS-111,
// where the cooldown was set on cast but never decremented — permanently
// locking out every active skill after a single use.
// @test-link [[mech_skill_validation]]
func TestRuleSkillCooldownClearsAfterElapsedTurns(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// cast issues the skill on the caster's current turn and returns the reply.
	cast := func() *message.Message {
		gs.Turner.ForceTurn(fake.Attacker)
		msg := message.Create(nil,
			rulermethods.ControllerUseSkill{
				EntityID:     fake.Attacker,
				ControllerID: fake.AttackerControllerID,
				Target:       fake.FoePosition,
				SkillID:      fake.SkillID,
			}, nil)
		reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))
		return reply
	}

	// elapseTurn runs the caster's end-of-turn, which both clears HasActed (so a
	// re-cast reaches the cooldown gate rather than tripping entity.alreadyacted)
	// and ticks the skill cooldown down by one — the exact production per-turn
	// path exercised by advanceTurn.
	elapseTurn := func() {
		gs.Turner.ForceTurn(fake.Attacker)
		msg := message.Create(nil,
			rulermethods.EndOfTurn{
				EntityID:     fake.Attacker,
				ControllerID: fake.AttackerControllerID,
			}, nil)
		ok, reply := EndOfTurn(gs, msg, msg.TargetMethod.(rulermethods.EndOfTurn), gs.Entities[fake.Attacker])
		if !ok {
			t.Fatalf("EndOfTurn failed unexpectedly: %s", reply.ErrorKey)
		}
	}

	cooldown := func() int {
		return gs.Entities[fake.Attacker].Skills[fake.SkillID].Cooldown
	}

	// Turn 0: cast succeeds and puts the skill on its full cooldown (3).
	if reply := cast(); reply.HasError {
		t.Fatalf("Expected first cast to succeed, got '%s'", reply.ErrorKey)
	}
	if cooldown() != 3 {
		t.Fatalf("Expected cooldown 3 right after cast, got %d", cooldown())
	}

	// End of the casting turn ticks it once (3 -> 2): a single decrement must
	// NOT clear the cooldown, or the skill would be re-castable too soon.
	elapseTurn()
	if cooldown() != 2 {
		t.Fatalf("Expected cooldown 2 after one elapsed turn, got %d", cooldown())
	}

	// Turn 1 (immediately following the cast): still rejected.
	if reply := cast(); !reply.HasError || reply.ErrorKey != "skill.cooldown" {
		t.Fatalf("Expected 'skill.cooldown' on the immediately-following turn, got HasError=%v key='%s'", reply.HasError, reply.ErrorKey)
	}

	// Two more of the caster's turns elapse (2 -> 1 -> 0).
	elapseTurn()
	if cooldown() != 1 {
		t.Fatalf("Expected cooldown 1, got %d", cooldown())
	}
	elapseTurn()
	if cooldown() != 0 {
		t.Fatalf("Expected cooldown 0 after the cooldown elapses, got %d", cooldown())
	}

	// Cooldown exhausted: the skill is castable again (the gate now passes).
	if reply := cast(); reply.HasError {
		t.Fatalf("Expected skill to be castable again after cooldown elapsed, got '%s'", reply.ErrorKey)
	}
	// ...and casting re-arms the cooldown back to its max value.
	if cooldown() != 3 {
		t.Fatalf("Expected cooldown re-armed to 3 after re-cast, got %d", cooldown())
	}
}
