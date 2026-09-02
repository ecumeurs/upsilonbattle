package rules

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// @test-link [[mechanic_channeling_mechanic]]

// withChannel turns a base skill into a channeled damage skill: a Damage effect
// plus a Channeling cost of the given delay (Delay cost is left absent → 500).
func withChannel(sk skill.Skill, damage, delay int) skill.Skill {
	sk.Effect.Properties = append(sk.Effect.Properties,
		defaultproperty.MakeIntProperty(property.DamageScale, damage, property.FriendlyController, property.Skill))
	ch := defaultproperty.MakeIntCounterProperty(property.Channeling, 0, delay, property.FriendlyController, property.Skill)
	sk.Costs[ch.Name(property.GameMaster)] = ch
	return sk
}

// moveEntity relocates an entity on the grid and updates its stored position.
func moveEntity(gs *gamestate.GameState, id uuid.UUID, to position.Position) {
	e := gs.Entities[id]
	gs.Grid.MoveEntity(e.Position, to, id)
	e.Position = to
	gs.Entities[id] = e
}

// TestChannel_CastInit_PaysAndLocksWithoutEffect: using a channeled skill enters
// IsCasting (capturing the entity target), locks the caster, and applies NO effect.
func TestChannel_CastInit_PaysAndLocksWithoutEffect(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	sk := withChannel(fake.Skill, 50, 400)
	addSkillToEntity(gs, fake.Attacker, sk)

	msg := message.Create(nil, rulermethods.ControllerUseSkill{
		EntityID:     fake.Attacker,
		ControllerID: fake.AttackerControllerID,
		Target:       fake.FoePosition,
		SkillID:      sk.ID,
	}, nil)
	reply, damaged, affected := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	assert.False(t, reply.HasError, "channel cast should succeed")
	assert.Empty(t, damaged, "no damage at cast init")
	assert.Empty(t, affected, "no status at cast init")

	caster := gs.Entities[fake.Attacker]
	assert.True(t, caster.IsChanneling(), "caster should be channeling")
	assert.Equal(t, sk.ID, caster.IsCasting.SkillID)
	assert.Equal(t, fake.Foe, caster.IsCasting.TargetEntity, "entity-target captured for follow")
	assert.Nil(t, caster.IsCasting.TargetPos, "entity-target channel stores no tile")
	assert.True(t, caster.HasActed(), "caster locked out")

	assert.Equal(t, 10, gs.Entities[fake.Foe].GetPropertyC(property.HP).GetValue(),
		"foe must be untouched until resolution")
}

// TestChannel_CastInit_TileTargetStoresPosition: a tile-targeted channel stores a
// fixed TargetPos rather than an entity.
func TestChannel_CastInit_TileTargetStoresPosition(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	sk := withChannel(fake.Skill, 50, 400)
	tt := defaultproperty.MakeValidatedStringProperty(property.TargetType, "Tile", property.FriendlyController, property.Skill, []string{"Tile"})
	sk.Targeting[tt.Name(property.GameMaster)] = tt
	addSkillToEntity(gs, fake.Attacker, sk)

	// An empty in-range tile next to the caster (5,5,3).
	tile := position.New(6, 5, 3)
	msg := message.Create(nil, rulermethods.ControllerUseSkill{
		EntityID:     fake.Attacker,
		ControllerID: fake.AttackerControllerID,
		Target:       tile,
		SkillID:      sk.ID,
	}, nil)
	reply, _, _ := UseSkill(gs, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))

	assert.False(t, reply.HasError)
	caster := gs.Entities[fake.Attacker]
	assert.True(t, caster.IsChanneling())
	assert.Equal(t, uuid.Nil, caster.IsCasting.TargetEntity, "tile channel stores no entity")
	if assert.NotNil(t, caster.IsCasting.TargetPos) {
		assert.True(t, caster.IsCasting.TargetPos.Equals(tile))
	}
}

// TestChannel_EndOfTurn_ReschedulesByChannelDelay: ending a casting entity's turn
// reschedules it by the channel delay (not the flat 300 Pass) and keeps it locked.
func TestChannel_EndOfTurn_ReschedulesByChannelDelay(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	sk := withChannel(fake.Skill, 50, 400)
	addSkillToEntity(gs, fake.Attacker, sk)

	a := gs.Entities[fake.Attacker]
	a.IsCasting = &entity.CastingState{SkillID: sk.ID, TargetEntity: fake.Foe}
	a.CurrentDelay = 0
	acted := a.GetProperty(property.HasActed)
	acted.Set(true)
	a.UpdateProperty(acted)
	gs.Entities[fake.Attacker] = a
	gs.Turner.ForceTurn(fake.Attacker)

	msg := message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     fake.Attacker,
		ControllerID: fake.AttackerControllerID,
	}, nil)
	ok, _ := EndOfTurn(gs, msg, msg.TargetMethod.(rulermethods.EndOfTurn), gs.Entities[fake.Attacker])

	assert.True(t, ok)
	caster := gs.Entities[fake.Attacker]
	assert.Equal(t, 400, caster.CurrentDelay, "rescheduled by channel delay, not 300")
	assert.True(t, caster.HasActed(), "lockout retained (flags not reset)")

	d, err := gs.Turner.GetEntityDelay(fake.Attacker)
	assert.NoError(t, err)
	assert.Equal(t, 400, d, "dormant in Turner at +channelDelay")
}

// TestChannel_Interruption_AccumulatesAndBreaks: damage fills the gauge at 10 per
// 1 damage; reaching 100 clears IsCasting and re-queues at the skill Delay (500),
// not the channel delay.
func TestChannel_Interruption_AccumulatesAndBreaks(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	sk := withChannel(fake.Skill, 50, 400)
	addSkillToEntity(gs, fake.Foe, sk)

	foe := gs.Entities[fake.Foe]
	foe.IsCasting = &entity.CastingState{SkillID: sk.ID, TargetEntity: fake.Attacker}
	gs.Entities[fake.Foe] = foe
	gs.Turner.RemoveEntity(fake.Foe)
	gs.Turner.AddEntity(fake.Foe, 400) // dormant at channel delay

	f := gs.Entities[fake.Foe]
	broke := ApplyInterruption(gs, &f, 5)
	assert.False(t, broke, "50 < 100 must not break")
	assert.Equal(t, 50, f.IsCasting.Interruption)

	broke = ApplyInterruption(gs, &f, 5)
	assert.True(t, broke, "reaching 100 breaks the channel")
	assert.Nil(t, f.IsCasting, "IsCasting cleared on interrupt")

	d, err := gs.Turner.GetEntityDelay(fake.Foe)
	assert.NoError(t, err)
	assert.Equal(t, 500, d, "recovers at skill Delay (500), not the channel delay")
}

// TestChannel_Interruption_NoopWhenNotCasting: the helper is a no-op for a
// non-casting entity.
func TestChannel_Interruption_NoopWhenNotCasting(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	f := gs.Entities[fake.Foe]
	assert.False(t, ApplyInterruption(gs, &f, 100))
	assert.Nil(t, f.IsCasting)
}

// TestChannel_Resolve_EntityTargetFollowsAndHits: resolution re-derives the target
// entity's CURRENT position (it moved) and still hits it.
func TestChannel_Resolve_EntityTargetFollowsAndHits(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	sk := withChannel(fake.Skill, 50, 400)
	addSkillToEntity(gs, fake.Attacker, sk)

	a := gs.Entities[fake.Attacker]
	a.IsCasting = &entity.CastingState{SkillID: sk.ID, TargetEntity: fake.Foe}
	gs.Entities[fake.Attacker] = a

	// Foe moves from (5,6,3) to another tile still within range 1 of the caster (5,5,3).
	moveEntity(gs, fake.Foe, position.New(6, 5, 3))

	damaged, _, fizzled := ResolveChannel(gs, fake.Attacker)

	assert.False(t, fizzled, "in-range target must resolve")
	if assert.Len(t, damaged, 1) {
		assert.Greater(t, damaged[0].Damage, 0)
		assert.Equal(t, fake.Foe, damaged[0].Entity.ID)
	}
	assert.False(t, gs.Entities[fake.Attacker].IsChanneling(), "caster released after resolution")
}

// TestChannel_Resolve_FizzleWhenTargetOutOfRange: if the target moved out of range
// during the channel, resolution fizzles — no damage, caster still released.
func TestChannel_Resolve_FizzleWhenTargetOutOfRange(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	sk := withChannel(fake.Skill, 50, 400)
	addSkillToEntity(gs, fake.Attacker, sk)

	a := gs.Entities[fake.Attacker]
	a.IsCasting = &entity.CastingState{SkillID: sk.ID, TargetEntity: fake.Foe}
	gs.Entities[fake.Attacker] = a

	moveEntity(gs, fake.Foe, position.New(5, 9, 3)) // distance 4 > range 1

	damaged, affected, fizzled := ResolveChannel(gs, fake.Attacker)

	assert.True(t, fizzled, "out-of-range target fizzles")
	assert.Empty(t, damaged)
	assert.Empty(t, affected)
	assert.Equal(t, 10, gs.Entities[fake.Foe].GetPropertyC(property.HP).GetValue(), "foe untouched")
	assert.False(t, gs.Entities[fake.Attacker].IsChanneling(), "caster released even on fizzle")
}

// TestChannel_Resolve_FizzleOnDeadTarget: an entity-target channel whose target is
// gone resolves to nothing.
func TestChannel_Resolve_FizzleOnDeadTarget(t *testing.T) {
	gs, fake := makeGameStateForTwoSkill()
	sk := withChannel(fake.Skill, 50, 400)
	addSkillToEntity(gs, fake.Attacker, sk)

	a := gs.Entities[fake.Attacker]
	a.IsCasting = &entity.CastingState{SkillID: sk.ID, TargetEntity: uuid.New()} // never existed
	gs.Entities[fake.Attacker] = a

	damaged, affected, fizzled := ResolveChannel(gs, fake.Attacker)

	assert.True(t, fizzled)
	assert.Empty(t, damaged)
	assert.Empty(t, affected)
	assert.False(t, gs.Entities[fake.Attacker].IsChanneling())
}
