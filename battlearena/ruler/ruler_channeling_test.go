package ruler

// @test-link [[mechanic_channeling_mechanic]]

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/google/uuid"
)

// channeledDamageSkill builds a single-target channeled damage skill (default
// Entity targeting, range 1) with the given channel delay.
func channeledDamageSkill(damage, delay int) skill.Skill {
	sk := skill.New()
	sk.Name = "Channeled Bolt"
	sk.Effect.Properties = append(sk.Effect.Properties,
		defaultproperty.MakeIntProperty(property.DamageScale, damage, property.FriendlyController, property.Skill))
	ch := defaultproperty.MakeIntCounterProperty(property.Channeling, 0, delay, property.FriendlyController, property.Skill)
	sk.Costs[ch.Name(property.GameMaster)] = ch
	return sk
}

// TestRulerChannelLifecycle verifies the end-to-end channeling wiring through the
// Ruler actor: using a channeled skill auto-ends the caster's turn (no controller
// EndOfTurn), the caster is rescheduled and re-picked to RESOLVE the channel (the
// target takes damage), and the turn then advances to the appropriate next entity.
func TestRulerChannelLifecycle(t *testing.T) {
	r := NewRuler(uuid.New())
	defer r.Stop()
	r.ShotClockDuration = 0 // no wall-clock dependency

	r.NbControllers = 2
	r.NbEntitiesPerController = 1
	r.GameState.Grid = grid.NewGrid(10, 10, 3)

	ctrl1 := NewFake("Caster-Controller")
	ctrl2 := NewFake("Target-Controller")

	// Caster (ctrl1, team 1) — delay 100 → goes first. Channel delay (50) is shorter
	// than the target's remaining delay, so the caster is re-picked next and resolves.
	caster := entity.New()
	caster.ControllerID = ctrl1.ID
	caster.Type = entity.Character
	caster.CurrentDelay = 100
	casterPos := position.Position{X: 5, Y: 5, Z: 3}
	caster.Position = casterPos
	for _, v := range def.PropertiesForCharacter() {
		caster.Properties[v.Name(property.GameMaster)] = v
	}
	caster.RepsertPropertyValue(property.TeamID, 1)
	sk := channeledDamageSkill(50, 50)
	caster.Skills[sk.ID] = sk
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), casterPos, caster.ID)
	r.GameState.Entities[caster.ID] = caster
	r.GameState.Turner.AddEntity(caster.ID, caster.CurrentDelay)

	// Target (ctrl2, team 2) — adjacent, delay 200, high HP so it survives the hit.
	target := entity.New()
	target.ControllerID = ctrl2.ID
	target.Type = entity.Character
	target.CurrentDelay = 200
	targetPos := position.Position{X: 5, Y: 6, Z: 3}
	target.Position = targetPos
	for _, v := range def.PropertiesForCharacter() {
		target.Properties[v.Name(property.GameMaster)] = v
	}
	target.RepsertPropertyValue(property.TeamID, 2)
	target.RepsertPropertyValue(property.HP, 200)
	target.RepsertPropertyValue(property.Defense, 0)
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), targetPos, target.ID)
	r.GameState.Entities[target.ID] = target
	r.GameState.Turner.AddEntity(target.ID, target.CurrentDelay)

	casterID := caster.ID
	targetID := target.ID
	skillID := sk.ID

	r.Start()

	for _, c := range []*FakeController{ctrl1, ctrl2} {
		dChan := make(chan *message.Message, 1)
		r.SendActor(message.Create(nil, rulermethods.AddController{Controller: c, ControllerID: c.ID}, rulermethods.AddControllerReply{}), dChan)
		<-dChan
	}

	ctrl1.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	// First turn → caster.
	timeout := time.After(5 * time.Second)
waitFirst:
	for {
		select {
		case msg := <-ctrl1.Inbox:
			if nxt, ok := msg.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				if nxt.Entity.ID != casterID {
					t.Fatalf("expected first turn for caster, got %s", nxt.Entity.ID.String()[0:8])
				}
				break waitFirst
			}
		case <-ctrl2.Inbox:
		case <-timeout:
			t.Fatal("timeout waiting for caster's first turn")
		}
	}

	drainInbox(ctrl1, ctrl2)

	// Caster uses the channeled skill. We deliberately do NOT send EndOfTurn:
	// the cast must auto-pass, the channel must resolve, and the turn advance.
	useReplyChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.ControllerUseSkill{
		EntityID:     casterID,
		ControllerID: ctrl1.ID,
		Target:       targetPos,
		SkillID:      skillID,
	}, rulermethods.ControllerUseSkillReply{}), useReplyChan)
	useReply := <-useReplyChan
	if useReply.HasError {
		t.Fatalf("channel cast failed: %s", useReply.ErrorKey)
	}

	// The target must eventually receive its turn — proving the channel resolved
	// and the turn advanced without any EndOfTurn from the caster's controller.
	timeout2 := time.After(5 * time.Second)
waitNext:
	for {
		select {
		case msg := <-ctrl2.Inbox:
			if nxt, ok := msg.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				if nxt.Entity.ID != targetID {
					t.Fatalf("expected next turn for target, got %s", nxt.Entity.ID.String()[0:8])
				}
				break waitNext
			}
		case <-ctrl1.Inbox:
		case <-timeout2:
			t.Fatal("timeout: channel never resolved / turn never advanced to target")
		}
	}

	// The target must have taken the channeled damage, and the caster must no longer
	// be channeling.
	ents, _ := testingFetchEntities(r)
	var foundTarget, foundCaster bool
	for _, e := range ents {
		switch e.ID {
		case targetID:
			foundTarget = true
			if e.GetPropertyC(property.HP).GetValue() >= 200 {
				t.Fatalf("expected target to take channeled damage, HP still %d", e.GetPropertyC(property.HP).GetValue())
			}
		case casterID:
			foundCaster = true
			if e.IsChanneling() {
				t.Fatal("caster should be released after channel resolution")
			}
		}
	}
	if !foundTarget || !foundCaster {
		t.Fatalf("expected both entities present (target=%v caster=%v)", foundTarget, foundCaster)
	}

	ctrl1.Stop()
	ctrl2.Stop()
}
