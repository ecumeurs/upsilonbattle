package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/google/uuid"
)

// @test-link [[mech_skill_validation]]
// @test-link [[mech_skill_validation]]

// FakeStateSkill extends FakeStateAttack with skill-specific context for testing.
type FakeStateSkill struct {
	FakeStateAttack
	SkillID uuid.UUID
	Skill   skill.Skill
}

// addSkillToEntity is a test helper to inject a skill into an entity's skill set.
func addSkillToEntity(gs *gamestate.GameState, entityID uuid.UUID, skill skill.Skill) {
	e := gs.Entities[entityID]
	e.Skills[skill.ID] = skill
	gs.Entities[entityID] = e
}

// makeGameStateForTwoSkill creates a test scenario with two entities and a default skill assigned to the attacker.
func makeGameStateForTwoSkill() (*gamestate.GameState, FakeStateSkill) {
	// attacker is at (5,5,3), foe is at (5,6,3) and they are facing each other.
	gs, f := makeGameStateForTwoAttack()
	fake := FakeStateSkill{
		FakeStateAttack: f,
	}

	fake.Skill = skill.New()
	fake.SkillID = fake.Skill.ID

	// Fake skill has no targeting options (1 range, 1 tile zone)
	// No Effet(at all)
	// No cost
	// default is ready to be used(cooldown = 0)
	// default target is entity (any)

	// assign skill to attacker.
	addSkillToEntity(gs, fake.Attacker, fake.Skill)

	return gs, fake
}
