package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type FakeStateSkill struct {
	FakeStateAttack
	SkillID uuid.UUID
	Skill   skill.Skill
}

func (gs *GameState) addSkillToEntity(entityID uuid.UUID, skill skill.Skill) {
	e := gs.Entities[entityID]
	e.Skills[skill.ID] = skill
	gs.Entities[entityID] = e
}

func makeGameStateForTwoSkill() (*GameState, FakeStateSkill) {
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

	gs.addSkillToEntity(fake.Attacker, fake.Skill)

	return gs, fake
}

// TestRuleSkillFailOutOfTurn ensures that a skill usage fails if it is not the entity's turn.
// @test-link [[mech_skill_validation_turn_controller_identity_verification]]
func TestRuleSkillFailOutOfTurn(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	// Force the turn to another entity.
	gs.Turner.ForceTurn(fake.Entity2)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when it is not the attacker's turn.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.turn.missmatch" {
		t.Errorf("Expected error 'entity.turn.missmatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailWrongController ensures that a skill usage fails if the controller ID does not match the entity's controller.
// @test-link [[mech_skill_validation_turn_controller_identity_verification]]
func TestRuleSkillFailWrongController(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.FoeControllerID, // Wrong controller!
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill with an unauthorized controller.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.controller.missmatch" {
		t.Errorf("Expected error 'entity.controller.missmatch', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailUnknownEntity ensures that a skill usage fails if the entity ID is invalid.
// @test-link [[mech_skill_validation_existence_verification]]
func TestRuleSkillFailUnknownEntity(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     uuid.New(), // Random UUID
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill for a non-existent entity.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.notfound" {
		t.Errorf("Expected error 'entity.notfound', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailUnknownSkill ensures that a skill usage fails if the skill ID is invalid for the entity.
// @test-link [[mech_skill_validation_existence_verification]]
func TestRuleSkillFailUnknownSkill(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      uuid.New(), // Random UUID
		}, nil)

	// Attempt to use a skill that the entity does not possess.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.notfound" {
		t.Errorf("Expected error 'skill.notfound', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailAlreadyActed ensures that an entity cannot use a skill if it has already acted this turn.
// @test-link [[mech_skill_validation_action_state_verification]]
func TestRuleSkillFailAlreadyActed(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	ent := gs.Entities[fake.Attacker]
	// Manually set the HasActed property to true.
	prop := ent.GetProperty(property.HasActed)
	prop.Set(true)
	ent.UpdateProperty(prop)
	gs.Entities[fake.Attacker] = ent

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when the entity is already in 'acted' state.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.alreadyacted" {
		t.Errorf("Expected error 'entity.alreadyacted', got '%s'", reply.ErrorKey)
	}
}

// Targeting checks

// TestRuleSkillFailTargetOutOfRange ensures that a skill fails if the target is outside the skill's defined range.
// @test-link [[mech_skill_validation_range_limit_verification]]
func TestRuleSkillFailTargetOutOfRange(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Default skill has range 1, targeting (5, 7, 3) from (5, 5, 3) is distance 2.
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(5, 7, 3), 
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill on an out-of-range target.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.range" {
		t.Errorf("Expected error 'skill.target.range', got '%s'", reply.ErrorKey)
	}

}

// TestRuleSkillFailTargetOutOfGrid ensures that a skill fails if the target coordinate is outside the board boundaries.
// @test-link [[mech_skill_validation_grid_boundaries_verification]]
func TestRuleSkillFailTargetOutOfGrid(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Move attacker to the grid border (0, 0).
	attacker := gs.Entities[fake.Attacker]
	attacker.Position = position.Position{X: 0, Y: 0, Z: 3}
	gs.Entities[fake.Attacker] = attacker
	gs.Grid.MoveEntity(position.New(5, 5, 3), position.New(0, 0, 3), fake.Attacker)
	fake.AttackerPosition = position.New(0, 0, 3)

	// Targeting (-1, 0, 3) is outside the grid.
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(-1, 0, 3),
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill on a target outside grid boundaries.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.outofgrid" {
		t.Errorf("Expected error 'skill.target.outofgrid', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Entity ensures that a skill fails if no valid entity target is present at the target location.
// @test-link [[mech_skill_validation_entity_targeting_rules_verification]]
func TestRuleSkillFailTargetNoApplicableTarget_Entity(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Move the foe away so that the target cell (5, 6, 3) is empty.
	foe := gs.Entities[fake.Foe]
	foe.Position = position.Position{X: 5, Y: 7, Z: 3}
	gs.Entities[fake.Foe] = foe
	gs.Grid.MoveEntity(position.New(5, 6, 3), position.New(5, 7, 3), fake.Foe)
	fake.FoePosition = position.New(5, 7, 3)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(5, 6, 3), // Target an empty cell.
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill on an empty cell when an entity target is required.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Self ensures that a skill with 'Self' targeting fails if used on another entity.
// @test-link [[mech_skill_validation_entity_targeting_rules_verification]]
func TestRuleSkillFailTargetNoApplicableTarget_Self(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'Self' target only.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeSelf))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	rng := fake.Skill.GetProperty(property.Range).(*defaultproperty.DefaultIntCounterProperty)
	rng.Value = 0
	fake.Skill.Targeting[rng.Name(property.GameMaster)] = rng
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // Attempting to target someone else.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting another entity fails.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.self" {
		t.Errorf("Expected error 'skill.target.self', got '%s'", reply.ErrorKey)
	}

	// Verify that targeting self succeeds.
	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.AttackerPosition, 
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Ennemy ensures that an 'EnemyOnly' skill fails if used on an ally or empty cell.
// @test-link [[mech_skill_validation_entity_targeting_rules_verification]]
func TestRuleSkillFailTargetNoApplicableTarget_Ennemy(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'EnemyOnly'.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeEnemyOnly))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // Empty cell.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an empty cell fails.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	// Move an ally to the target position.
	ally := gs.Entities[fake.Entity2]
	ally.Position = position.Position{X: 4, Y: 5, Z: 3}
	gs.Entities[fake.Entity2] = ally
	gs.Grid.MoveEntity(position.New(1, 0, 3), position.New(4, 5, 3), fake.Entity2)

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // Target an ally.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an ally fails for an 'EnemyOnly' skill.
	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // Target an actual enemy.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an actual enemy succeeds.
	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Ally ensures that a 'FriendOnly' skill fails if used on an enemy or empty cell.
// @test-link [[mech_skill_validation_entity_targeting_rules_verification]]
func TestRuleSkillFailTargetNoApplicableTarget_Ally(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'FriendOnly'.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeFriendOnly))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // Empty cell.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an empty cell fails.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	// Attempt to target an enemy.
	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, 
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an enemy fails for a 'FriendOnly' skill.
	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	// Move an ally to the target position.
	ally := gs.Entities[fake.Entity2]
	ally.Position = position.Position{X: 4, Y: 5, Z: 3}
	gs.Entities[fake.Entity2] = ally
	gs.Grid.MoveEntity(position.New(1, 0, 3), position.New(4, 5, 3), fake.Entity2)

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // Target an ally.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an actual ally succeeds.
	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Cell ensures that a 'Tile' target skill fails if used on a cell occupied by an entity.
// @test-link [[mech_skill_validation_entity_targeting_rules_verification]]
func TestRuleSkillFailTargetNoApplicableTarget_Cell(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'Tile' target only.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeTile))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // Target an occupied cell.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an occupied cell fails for a 'Tile' only skill.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}
}

// Check ability to pay for the skill usage.

// TestRuleSkillFailCooldown ensures that a skill fails if it is currently on cooldown.
// @test-link [[mech_skill_validation_economic_cost_verification_cooldown_check]]
func TestRuleSkillFailCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set the skill cooldown to 1.
	fake.Skill.Cooldown = 1
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill that is on cooldown.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cooldown" {
		t.Errorf("Expected error 'skill.cooldown', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailMP ensures that a skill fails if the entity does not have enough Mana Points (MP).
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill MP cost higher than the entity's current MP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 11 /*Default MP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when MP is insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.mp" {
		t.Errorf("Expected error 'skill.cost.mp', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailSP ensures that a skill fails if the entity does not have enough Stamina Points (SP).
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailSP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill SP cost higher than the entity's current SP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.SPLeech, 11 /*Default HP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when SP is insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.sp" {
		t.Errorf("Expected error 'skill.cost.sp', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailMvt ensures that a skill fails if the entity does not have enough Movement points.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill Movement cost higher than the entity's current Movement points.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 11 /*Default Mvt is 3*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when Movement points are insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.mvt" {
		t.Errorf("Expected error 'skill.cost.mvt', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailHP ensures that a skill fails if the entity does not have enough Health Points (HP) to pay the cost.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillFailHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a skill HP cost higher than the entity's current HP.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 11 /*Default HP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Attempt to use a skill when HP is insufficient.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.hp" {
		t.Errorf("Expected error 'skill.cost.hp', got '%s'", reply.ErrorKey)
	}
}

// Check skill cost

// TestRuleSkillSetCooldown verifies that using a skill correctly applies its cooldown to the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_cooldown_check]]
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
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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

// TestRuleSkillDeduceMP verifies that using a skill correctly deducts the Mana Point (MP) cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an MP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// Verify that MP was deducted in the GameState (from 10 to 9).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.MP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected MP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// Verify that MP was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.MP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected MP to be 9, got %d", prop.GetValue())
	}
}

// TestRuleSkillDeduceSP verifies that using a skill correctly deducts the Stamina Point (SP) cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceSP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an SP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.SPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// Verify that SP was deducted in the GameState (from 10 to 9).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.SP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected SP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// Verify that SP was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.SP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected SP to be 9, got %d", prop.GetValue())
	}
}

// TestRuleSkillDeduceMvt verifies that using a skill correctly deducts the Movement point cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set a Movement cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// Verify that Movement was deducted in the GameState (from 3 to 2).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.Movement)
	if prop.GetValue() != 2 { 
		t.Errorf("Expected Movement(in GameState) to be 2, got %d", prop.GetValue())
	}

	// Verify that Movement was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.Movement)
	if prop.GetValue() != 2 { 
		t.Errorf("Expected Movement to be 2, got %d", prop.GetValue())
	}
}

// TestRuleSkillDeduceHP verifies that using a skill correctly deducts the Health Point (HP) cost from the entity.
// @test-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
func TestRuleSkillDeduceHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Set an HP cost of 1.
	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	// Execute skill usage.
	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// Verify that HP was deducted in the GameState (from 10 to 9).
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.HP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected HP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// Verify that HP was deducted in the reply message.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Attacker.GetPropertyC(property.HP)
	if prop.GetValue() != 9 { 
		t.Errorf("Expected HP to be 9, got %d", prop.GetValue())
	}
}
