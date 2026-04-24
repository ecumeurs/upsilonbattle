package rules

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// stringProperty is a simple implementation of property.Property for test purposes.
type stringProperty struct {
	property.Property
	name  string
	value string
}

func (s stringProperty) Name(i property.InformationLevel) string { return s.name }
func (s stringProperty) Get() interface{}                        { return s.value }
func (s stringProperty) Set(p interface{})                       { s.value = p.(string) }
func (s stringProperty) GetType() property.PropertyType          { return property.Skill }
func (s stringProperty) UserFriendlyGet(i property.InformationLevel) interface{} {
	return s.value
}
func (s stringProperty) Duplicate() property.Property {
	return stringProperty{name: s.name, value: s.value}
}

// TestRuleTrapTriggerOnEnter verifies that an entity entering a cell with an OnEnter trap triggers it.
// @test-link [[mech_trigger_system]]
func TestRuleTrapTriggerOnEnter(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	trapPos := position.Position{X: 0, Y: 1, Z: 3}

	// Create a poison trap effect
	trapEffect := effect.Effect{
		Name:     "Poison Trap",
		CasterID: uuid.Nil,
		Properties: []property.Property{
			defaultproperty.MakeIntProperty(property.PoisonPower, 10, property.GameMaster, property.Skill),
			defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.GameMaster, property.Skill),
			stringProperty{name: string(property.TriggerType), value: string(property.TriggerOnEnter)},
			defaultproperty.MakeBoolProperty(property.RemoveOnTrigger, true, property.GameMaster, property.Skill),
		},
	}
	effectID := uuid.New()
	gs.Effects[effectID] = trapEffect
	gs.PositionalEffects[trapPos] = append(gs.PositionalEffects[trapPos], effectID)

	// Move entity 1 into the trap
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{
				trapPos,
			},
		}, nil)

	gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	// Verify entity 1 is poisoned
	ent1 := gs.Entities[fake.Entity1]
	poison := ent1.GetPropertyI(property.Poison).I()
	if poison != 10 {
		t.Errorf("Expected poison power 10, got %d", poison)
	}

	// Verify trap is removed (RemoveOnTrigger = true)
	if _, ok := gs.Effects[effectID]; ok {
		t.Errorf("Trap effect should have been removed")
	}
}

// TestRuleEntityExpiration verifies that a temporary entity is removed when its Duration reaches zero.
// @test-link [[mech_entity_expiration]]
func TestRuleEntityExpiration(t *testing.T) {
	gs, fake := makeGameStateForTwo()

	// Create a temporary entity with Duration = 1
	tempEnt := entity.New()
	tempEnt.ID = uuid.New()
	tempEnt.ControllerID = fake.Controller1
	tempEnt.Type = entity.TimeBased
	tempEnt.Position = position.Position{X: 5, Y: 5, Z: 3}
	tempEnt.CurrentDelay = 100
	tempEnt.RepsertPropertyCValue(property.EntityDuration, 1)
	tempEnt.RepsertPropertyCMaxValue(property.EntityDuration, 1)

	gs.Entities[tempEnt.ID] = tempEnt
	gs.Turner.AddEntity(tempEnt.ID, tempEnt.CurrentDelay)
	gs.Grid.MoveEntity(position.New(0, 0, 0), tempEnt.Position, tempEnt.ID)

	gs.Turner.ForceTurn(tempEnt.ID)

	// End turn for the temporary entity
	msg := message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     tempEnt.ID,
		ControllerID: fake.Controller1,
	}, nil)

	gs.EndOfTurn(msg, msg.TargetMethod.(rulermethods.EndOfTurn), tempEnt)

	// Verify entity is removed
	if _, ok := gs.Entities[tempEnt.ID]; ok {
		t.Errorf("Temporary entity should have been removed after expiration")
	}
}

// TestRuleExpiresWithCaster verifies that positional effects are removed when their caster dies.
// @test-link [[mechanic_effect_caster_tracking]]
func TestRuleExpiresWithCaster(t *testing.T) {
	gs, fake := makeGameStateForTwo()

	casterID := fake.Entity1
	effectPos := position.Position{X: 2, Y: 2, Z: 3}

	// Create an effect owned by Entity 1
	ownedEffect := effect.Effect{
		Name:     "Caster Owned Effect",
		CasterID: casterID,
		Properties: []property.Property{
			defaultproperty.MakeBoolProperty(property.ExpiresWithCaster, true, property.GameMaster, property.Skill),
			stringProperty{name: string(property.TriggerType), value: string(property.TriggerOnTurn)},
		},
	}
	effectID := uuid.New()
	gs.Effects[effectID] = ownedEffect
	gs.PositionalEffects[effectPos] = append(gs.PositionalEffects[effectPos], effectID)

	// Kill the caster
	gs.RemoveEntity(casterID)

	// Verify effect is removed
	if _, ok := gs.Effects[effectID]; ok {
		t.Errorf("Effect should have been removed when caster died")
	}
}

// TestRuleTrapTriggerOnExit verifies that an entity leaving a cell with an OnExit trap triggers it.
// @test-link [[mech_trigger_system]]
func TestRuleTrapTriggerOnExit(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	gs.Turner.ForceTurn(fake.Entity1)

	startPos := gs.Entities[fake.Entity1].Position
	destPos := position.Position{X: 0, Y: 1, Z: 3}

	// Create a trap that triggers OnExit
	trapEffect := effect.Effect{
		Name:     "Exit Trap",
		CasterID: uuid.Nil,
		Properties: []property.Property{
			defaultproperty.MakeIntProperty(property.Damage, 50, property.GameMaster, property.Skill),
			stringProperty{name: string(property.TriggerType), value: string(property.TriggerOnExit)},
		},
	}
	effectID := uuid.New()
	gs.Effects[effectID] = trapEffect
	gs.PositionalEffects[startPos] = append(gs.PositionalEffects[startPos], effectID)

	// Move away from the start position
	msg := message.Create(nil,
		rulermethods.ControllerMove{
			EntityID:     fake.Entity1,
			ControllerID: fake.Controller1,
			Path: []position.Position{
				destPos,
			},
		}, nil)

	gs.Move(msg, msg.TargetMethod.(rulermethods.ControllerMove))

	// Verify entity 1 took damage (Default HP is 10, but MakeGameStateForTwo uses 10. Damage is 50% of attack? No, ApplyDirectEffect uses (attack*damage/100)-defense)
	// ent1 attack is 3. 3 * 50 / 100 = 1. Defense is 0. 1 - 0 = 1 damage.
	ent1 := gs.Entities[fake.Entity1]
	hp := ent1.GetPropertyC(property.HP).GetValue()
	if hp >= 10 {
		t.Errorf("Expected HP < 10, got %d", hp)
	}
}

// TestRuleTrapTriggerOnTurn verifies that an entity starting its turn in a cell with an OnTurn trap triggers it.
// @test-link [[mech_trigger_system]]
func TestRuleTrapTriggerOnTurn(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	
	ent1 := gs.Entities[fake.Entity1]
	pos := ent1.Position

	// Create a trap that triggers OnTurn
	trapEffect := effect.Effect{
		Name:     "Turn Trap",
		CasterID: uuid.Nil,
		Properties: []property.Property{
			defaultproperty.MakeIntProperty(property.PoisonPower, 5, property.GameMaster, property.Skill),
			defaultproperty.MakeIntProperty(property.PoisonChance, 100, property.GameMaster, property.Skill),
			stringProperty{name: string(property.TriggerType), value: string(property.TriggerOnTurn)},
		},
	}
	effectID := uuid.New()
	gs.Effects[effectID] = trapEffect
	gs.PositionalEffects[pos] = append(gs.PositionalEffects[pos], effectID)

	// Begin turn
	gs.BeginingOfTurn(ent1)

	// Verify entity 1 is poisoned
	ent1 = gs.Entities[fake.Entity1]
	poison := ent1.GetPropertyI(property.Poison).I()
	if poison != 5 {
		t.Errorf("Expected poison power 5, got %d", poison)
	}
}
