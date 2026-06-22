package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// TestRuleSkillFailTargetOutOfRange ensures that a skill fails if the target is outside the skill's defined range.
// @test-link [[mech_skill_validation]]
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
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.range" {
		t.Errorf("Expected error 'skill.target.range', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetOutOfGrid ensures that a skill fails if the target coordinate is outside the board boundaries.
// @test-link [[mech_skill_validation]]
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
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.outofgrid" {
		t.Errorf("Expected error 'skill.target.outofgrid', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Entity ensures that a skill fails if no valid entity target is present at the target location.
// @test-link [[mech_skill_validation]]
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
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Self ensures that a skill with 'Self' targeting fails if used on another entity.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailTargetNoApplicableTarget_Self(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'Self' target only.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeSelf))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	rng := fake.Skill.GetProperty(property.Range).(*defaultproperty.DefaultIntCounterProperty)
	rng.Value = 0
	fake.Skill.Targeting[rng.Name(property.GameMaster)] = rng
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // Attempting to target someone else.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting another entity fails.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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

	reply, _, _ = UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Ennemy ensures that an 'EnemyOnly' skill fails if used on an ally or empty cell.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailTargetNoApplicableTarget_Ennemy(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'EnemyOnly'.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeEnemyOnly))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // Empty cell.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an empty cell fails.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
	reply, _, _ = UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
	reply, _, _ = UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Ally ensures that a 'FriendOnly' skill fails if used on an enemy or empty cell.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailTargetNoApplicableTarget_Ally(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'FriendOnly'.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeFriendOnly))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // Empty cell.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an empty cell fails.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
	reply, _, _ = UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

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
	reply, _, _ = UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}
}

// TestRuleSkillFailTargetNoApplicableTarget_Cell ensures that a 'Tile' target skill fails if used on a cell occupied by an entity.
// @test-link [[mech_skill_validation]]
func TestRuleSkillFailTargetNoApplicableTarget_Cell(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// Modify skill to be 'Tile' target only.
	prop := def.SkillProperty(property.TargetType)
	prop.Set(string(def.TargetTypeTile))
	fake.Skill.Targeting[property.TargetType.String()] = prop
	addSkillToEntity(gs, fake.Attacker, fake.Skill) 

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // Target an occupied cell.
			SkillID:      fake.SkillID,
		}, nil)

	// Verify that targeting an occupied cell fails for a 'Tile' only skill.
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}
}
