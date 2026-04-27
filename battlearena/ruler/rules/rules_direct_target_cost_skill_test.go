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

func TestRuleSkillFailOutOfTurn(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	gs.Turner.ForceTurn(fake.Entity2)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.turn.missmatch" {
		t.Errorf("Expected error 'entity.turn.missmatch', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailWrongController(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.FoeControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.controller.missmatch" {
		t.Errorf("Expected error 'entity.controller.missmatch', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailUnknownEntity(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     uuid.New(),
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.notfound" {
		t.Errorf("Expected error 'entity.notfound', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailUnknownSkill(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      uuid.New(),
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.notfound" {
		t.Errorf("Expected error 'skill.notfound', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailAlreadyActed(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	ent := gs.Entities[fake.Attacker]
	// this flag is set when the entity has acted or has been interrupted during movement by reaction or traps.
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

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "entity.alreadyacted" {
		t.Errorf("Expected error 'entity.alreadyacted', got '%s'", reply.ErrorKey)
	}
}

// Targeting checks

func TestRuleSkillFailTargetOutOfRange(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// default skill has range 1, so target is out of range.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(5, 7, 3), // just target somewhere else ;)
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.range" {
		t.Errorf("Expected error 'skill.target.range', got '%s'", reply.ErrorKey)
	}

}

func TestRuleSkillFailTargetOutOfGrid(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// move attacker at the border of the grid.
	attacker := gs.Entities[fake.Attacker]
	attacker.Position = position.Position{X: 0, Y: 0, Z: 3}
	gs.Entities[fake.Attacker] = attacker
	gs.Grid.MoveEntity(position.New(5, 5, 3), position.New(0, 0, 3), fake.Attacker)
	fake.AttackerPosition = position.New(0, 0, 3)

	// default skill has range 1, zone 1, so targeting toward -1,0,3 should be out of grid.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(-1, 0, 3),
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.outofgrid" {
		t.Errorf("Expected error 'skill.target.outofgrid', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailTargetNoApplicableTarget_Entity(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// move foe further away.
	foe := gs.Entities[fake.Foe]
	foe.Position = position.Position{X: 5, Y: 7, Z: 3}
	gs.Entities[fake.Foe] = foe
	gs.Grid.MoveEntity(position.New(5, 6, 3), position.New(5, 7, 3), fake.Foe)
	fake.FoePosition = position.New(5, 7, 3)

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(5, 6, 3), // skill has range 1 and zone 1, so there is noone under its scope.
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailTargetNoApplicableTarget_Self(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = def.MakeTargetTypeProperty(def.TargetTypeSelf)
	rng := fake.Skill.GetProperty(property.Range).(*def.RangeProperty)
	rng.MinRange = 0
	fake.Skill.Targeting[rng.Name(property.GameMaster)] = rng
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // target somebody else.
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.self" {
		t.Errorf("Expected error 'skill.target.self', got '%s'", reply.ErrorKey)
	}

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.AttackerPosition, // target self!
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}
}

func TestRuleSkillFailTargetNoApplicableTarget_Ennemy(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = def.MakeTargetTypeProperty(def.TargetTypeEnemyOnly)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // there is nobody there, so can't target the ennemy.
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	// move an ally (entity 2) to the target position.
	ally := gs.Entities[fake.Entity2]
	ally.Position = position.Position{X: 4, Y: 5, Z: 3}
	gs.Entities[fake.Entity2] = ally
	gs.Grid.MoveEntity(position.New(1, 0, 3), position.New(4, 5, 3), fake.Entity2)

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // there is an ally there
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	// and test again, shouldn't work.

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // there is somebody there nnd it's an ennemy
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}

}

func TestRuleSkillFailTargetNoApplicableTarget_Ally(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = def.MakeTargetTypeProperty(def.TargetTypeFriendOnly)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // there is nobody there, so can't target the ennemy.
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	// and test this time against an ennemy, shouldn't work.

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // there is somebody there nnd it's an ennemy
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}

	// move an ally (entity 2) to the target position.
	ally := gs.Entities[fake.Entity2]
	ally.Position = position.Position{X: 4, Y: 5, Z: 3}
	gs.Entities[fake.Entity2] = ally
	gs.Grid.MoveEntity(position.New(1, 0, 3), position.New(4, 5, 3), fake.Entity2)

	msg = message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       position.New(4, 5, 3), // there is an ally there
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ = gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if reply.HasError {
		t.Errorf("Didn't expect and error but got one: %s", reply.ErrorKey)
	}

}

func TestRuleSkillFailTargetNoApplicableTarget_Cell(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = def.MakeTargetTypeProperty(def.TargetTypeTile)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition, // there is somebody there, so can't target the cell itself.
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.target.none" {
		t.Errorf("Expected error 'skill.target.none', got '%s'", reply.ErrorKey)
	}
}

// Check ability to pay for the skill usage.

func TestRuleSkillFailCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Cooldown = 1
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cooldown" {
		t.Errorf("Expected error 'skill.cooldown', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 11 /*Default MP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.mp" {
		t.Errorf("Expected error 'skill.cost.mp', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailSP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.SPLeech, 11 /*Default HP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.sp" {
		t.Errorf("Expected error 'skill.cost.sp', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 11 /*Default Mvt is 3*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.mvt" {
		t.Errorf("Expected error 'skill.cost.mvt', got '%s'", reply.ErrorKey)
	}
}

func TestRuleSkillFailHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 11 /*Default HP is 10*/, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	if !reply.HasError {
		t.Errorf("Expected error, got none.")
	}

	if reply.ErrorKey != "skill.cost.hp" {
		t.Errorf("Expected error 'skill.cost.hp', got '%s'", reply.ErrorKey)
	}
}

// Check skill cost

func TestRuleSkillSetCooldown(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	// default cooldown cost is 3
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// didn't expect error
	if reply.HasError {
		t.Errorf("Expected no error, got '%s'", reply.ErrorKey)
	}

	// get skill from gs and check cooldown.
	skill := gs.Entities[fake.Attacker].Skills[fake.SkillID]
	if skill.Cooldown != 3 {
		t.Errorf("Expected cooldown to be 3(in GameState), got %d", skill.Cooldown)
	}

	// get skill from reply message and check cooldown.
	skill = reply.Content.(rulermethods.ControllerUseSkillReply).Entity.Skills[fake.SkillID]
	if skill.Cooldown != 3 {
		t.Errorf("Expected cooldown to be 3, got %d", skill.Cooldown)
	}
}

func TestRuleSkillDeduceMP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	// default cooldown cost is 3
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// get Entity from gamestate and check MP
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.MP)
	if prop.GetValue() != 9 { // from 10
		t.Errorf("Expected MP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// get Entity from reply message and check MP.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Entity.GetPropertyC(property.MP)
	if prop.GetValue() != 9 { // from 10
		t.Errorf("Expected MP to be 9, got %d", prop.GetValue())
	}
}

func TestRuleSkillDeduceSP(t *testing.T) {

	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.SPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	// default cooldown cost is 3
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// get Entity from gamestate and check MP
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.SP)
	if prop.GetValue() != 9 { // from 10
		t.Errorf("Expected SP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// get Entity from reply message and check MP.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Entity.GetPropertyC(property.SP)
	if prop.GetValue() != 9 { // from 10
		t.Errorf("Expected SP to be 9, got %d", prop.GetValue())
	}
}

func TestRuleSkillDeduceMvt(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.MvtCost, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	// default cooldown cost is 3
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// get Entity from gamestate and check MP
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.Movement)
	if prop.GetValue() != 2 { // from 3
		t.Errorf("Expected Movement(in GameState) to be 2, got %d", prop.GetValue())
	}

	// get Entity from reply message and check MP.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Entity.GetPropertyC(property.Movement)
	if prop.GetValue() != 2 { // from 3
		t.Errorf("Expected Movement to be 2, got %d", prop.GetValue())
	}
}

func TestRuleSkillDeduceHP(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()

	fake.Skill.Targeting[property.TargetType.String()] = defaultproperty.MakeIntProperty(property.HPLeech, 1, property.FriendlyController, property.Skill)
	gs.addSkillToEntity(fake.Attacker, fake.Skill) // update skill.

	// default cooldown cost is 3
	msg := message.Create(nil,
		rulermethods.ControllerUseSkill{
			EntityID:     fake.Attacker,
			ControllerID: fake.AttackerControllerID,
			Target:       fake.FoePosition,
			SkillID:      fake.SkillID,
		}, nil)

	reply, _, _ := gs.UseSkill(msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	// get Entity from gamestate and check MP
	prop := gs.Entities[fake.Attacker].GetPropertyC(property.HP)
	if prop.GetValue() != 9 { // from 10
		t.Errorf("Expected HP(in GameState) to be 9, got %d", prop.GetValue())
	}

	// get Entity from reply message and check MP.
	prop = reply.Content.(rulermethods.ControllerUseSkillReply).Entity.GetPropertyC(property.HP)
	if prop.GetValue() != 9 { // from 10
		t.Errorf("Expected HP to be 9, got %d", prop.GetValue())
	}
}
