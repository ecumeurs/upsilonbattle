// @test-link [[rule_entity_property_write_isolation]]
package rules

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
)

// TestRuleSkillDeduceHP_UnderActiveBuffDoesNotEscalate is ISS-144 scenario (c):
// paySkillCost reads HP composed (base + buffs) and must strip the buff's
// contribution before persisting to base. With a +10/+10 HP buff active and a
// skill costing 5 HP, the composed HP before cast is 20/20 (base 10/10 + buff
// 10/10); after paying the cost it must read 15/20 (base correctly reduced to
// 5, buff still applied on top) — not 25/30, which is what the pre-fix
// composed-read/base-write bug produces (it writes the full composed 20 minus
// 5 back to base, and the buff re-applies on top of that on the next read;
// the composed max escalates too, 20 -> 30, hence 25/30 and not 25/20).
// Pre-fix figure independently reproduced against unmodified HEAD, 2026-08-28.
// @test-link [[mech_skill_validation]]
func TestRuleSkillDeduceHP_UnderActiveBuffDoesNotEscalate(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	registerForeverBuff(gs, fake.Attacker, property.HP, 10, 10)

	// Sanity-check the composed pre-cast HP is 20/20 (base 10/10 + buff 10/10).
	preCast := gs.Entities[fake.Attacker].GetPropertyC(property.HP)
	if preCast.GetValue() != 20 || preCast.GetMaxValue() != 20 {
		t.Fatalf("precondition failed: expected composed HP 20/20 before cast, got %d/%d",
			preCast.GetValue(), preCast.GetMaxValue())
	}

	// Set an HP cost of 5.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 5, property.FriendlyController, property.Skill)
	addSkillToEntity(gs, fake.Attacker, fake.Skill)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))
	if reply.HasError {
		t.Fatalf("Expected no error, got '%s'", reply.ErrorKey)
	}

	postCast := gs.Entities[fake.Attacker].GetPropertyC(property.HP)
	if postCast.GetValue() != 15 || postCast.GetMaxValue() != 20 {
		t.Fatalf("expected composed HP 15/20 after paying a 5-cost skill under an active +10/+10 buff, got %d/%d",
			postCast.GetValue(), postCast.GetMaxValue())
	}
}
