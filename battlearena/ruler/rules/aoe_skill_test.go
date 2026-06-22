package rules

import (

	"testing"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// @test-link [[mech_skill_validation]]
// @test-link [[mech_skill_validation]]

// TestRuleSkill_AOE_ZoneProperty verifies that a skill with a Neighbours zone correctly
// targets multiple entities and that the effect applicator handles them.
func TestRuleSkill_AOE_ZoneProperty(t *testing.T) {
	// Setup game state with 1 attacker and 2 foes close to each other
	gs, fake := makeGameStateForTwoSkill()
	
	// Attacker is at (5,5,3)
	// Foe 1 is at (5,6,3)
	
	// Add Foe 2 at (4,6,3) - which is a neighbour of Foe 1
	foe2ID := uuid.New()
	foe2Pos := position.New(4, 6, 3)
	
	foe2Ent := entity.New()
	foe2Ent.ID = foe2ID
	foe2Ent.Name = "Foe 2"
	foe2Ent.ControllerID = fake.FoeControllerID
	foe2Ent.Position = foe2Pos
	foe2Ent.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		foe2Ent.Properties[v.Name(property.GameMaster)] = v
	}
	foe2Ent.RepsertPropertyValue(property.TeamID, 2)
	gs.Grid.MoveEntity(foe2Ent.Position, foe2Ent.Position, foe2Ent.ID)
	gs.Entities[foe2Ent.ID] = foe2Ent
	
	// Modify the skill to be an AOE skill
	sk := fake.Skill
	sk.Name = "Nova"
	
	// Set Zone to Neighbours
	zp := def.DefaultZone()
	zp.Set("Neighbours") // This triggers the new PatternType and pattern selection
	sk.Targeting[property.PropertyToString(property.Zone)] = zp
	
	// Set TargetType to EnemyOnly
	tt := def.SkillProperty(property.TargetType)
	tt.Set("EnemyOnly")
	sk.Targeting[property.PropertyToString(property.TargetType)] = tt
	
	// Set Effect to deal Damage
	sk.Effect.Properties = append(sk.Effect.Properties, defaultproperty.MakeIntProperty(property.Damage, 50, property.FriendlyController, property.Skill))
	
	addSkillToEntity(gs, fake.Attacker, sk)

	// Execute skill targeting Foe 1 at (5,6,3)
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // (5,6,3)
			SkillID:      sk.ID,
		}, nil)

	reply, damaged, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// Verification
	assert.False(t, reply.HasError, "Skill should not return an error")
	assert.Len(t, damaged, 2, "Both foes should be damaged by the AOE effect")
	
	// Check if both Foe 1 and Foe 2 are in the damaged list
	foundFoe1 := false
	foundFoe2 := false
	for _, d := range damaged {
		if d.Entity.ID == fake.Foe {
			foundFoe1 = true
		}
		if d.Entity.ID == foe2ID {
			foundFoe2 = true
		}
	}
	assert.True(t, foundFoe1, "Foe 1 should have been damaged")
	assert.True(t, foundFoe2, "Foe 2 should have been damaged")
	
	// Verify HP reduction
	ent1 := gs.Entities[fake.Foe]
	ent2 := gs.Entities[foe2ID]
	assert.Less(t, ent1.GetPropertyC(property.HP).GetValue(), 10, "Foe 1 should have lost HP")
	assert.Less(t, ent2.GetPropertyC(property.HP).GetValue(), 10, "Foe 2 should have lost HP")
}
