// @test-link [[rule_entity_property_write_isolation]]
package rules

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/google/uuid"
)

// ISS-144 regression coverage: entity property writes must isolate base state
// from buff-composed reads. See issues/ISS-144_20260827_buff_writeback_folds_into_base_state.md
// and upsilonbattle/docs/rule_entity_property_write_isolation.atom.md.

// registerForeverBuff attaches a Forever (non-expiring) buff for the given
// property to entityID, with the supplied Value/MaxValue deltas.
func registerForeverBuff(gs *gamestate.GameState, entityID uuid.UUID, prop property.Key, valueDelta, maxDelta int) {
	ent := gs.Entities[entityID]
	buff := property.MakeTemporaryProperties(0)
	buff.Forever = true
	buff.Properties[property.PropertyToString(prop)] = defaultproperty.MakeIntCounterProperty(prop, valueDelta, maxDelta, property.Public, property.Character)
	ent.RegisterBuff(buff)
	gs.Entities[entityID] = ent
}

// elapseTurnFor forces entityID's turn and runs EndOfTurn on it, failing the
// test immediately if the rule itself rejects the request.
func elapseTurnFor(t *testing.T, gs *gamestate.GameState, entityID, controllerID uuid.UUID) {
	t.Helper()
	gs.Turner.ForceTurn(entityID)
	msg := message.Create(nil, rulermethods.EndOfTurn{EntityID: entityID, ControllerID: controllerID}, nil)
	ok, reply := EndOfTurn(gs, msg, msg.TargetMethod.(rulermethods.EndOfTurn), gs.Entities[entityID])
	if !ok {
		t.Fatalf("EndOfTurn failed unexpectedly: %s", reply.ErrorKey)
	}
}

// TestEndOfTurn_MovementRestore_PairedBuffDoesNotEscalate reproduces ISS-144's
// primary example: a +10/+10 Movement item buff must not compound every time
// end-of-turn restores Movement to its max. Default character Movement is 3/3
// (see def.PropertiesForCharacter), so with the buff active composed Movement
// must stay pinned at 13/13 forever, not climb 13 -> 23 -> 33 -> 43 as it does
// pre-fix (composed max leaking into persisted base on every restore).
func TestEndOfTurn_MovementRestore_PairedBuffDoesNotEscalate(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	registerForeverBuff(gs, fake.Entity1, property.Movement, 10, 10)

	for turn := 1; turn <= 4; turn++ {
		elapseTurnFor(t, gs, fake.Entity1, fake.Controller1)

		mvt := gs.Entities[fake.Entity1].GetPropertyC(property.Movement)
		if mvt.GetValue() != 13 || mvt.GetMaxValue() != 13 {
			t.Fatalf("turn %d: expected composed Movement 13/13 (no escalation), got %d/%d",
				turn, mvt.GetValue(), mvt.GetMaxValue())
		}
	}
}

// TestEndOfTurn_MovementRestore_MaxOnlyBuffDoesNotEscalate reproduces ISS-144's
// max-only variant: a buff that grants extra Movement capacity but no extra
// current value (+0/+10) must not escalate either. Composed Movement must stay
// pinned at 3/13 forever, not climb 13/23 -> 23/33 -> 33/43 as it does pre-fix.
func TestEndOfTurn_MovementRestore_MaxOnlyBuffDoesNotEscalate(t *testing.T) {
	gs, fake := makeGameStateForTwo()
	registerForeverBuff(gs, fake.Entity1, property.Movement, 0, 10)

	for turn := 1; turn <= 4; turn++ {
		elapseTurnFor(t, gs, fake.Entity1, fake.Controller1)

		mvt := gs.Entities[fake.Entity1].GetPropertyC(property.Movement)
		if mvt.GetValue() != 3 || mvt.GetMaxValue() != 13 {
			t.Fatalf("turn %d: expected composed Movement 3/13 (no escalation), got %d/%d",
				turn, mvt.GetValue(), mvt.GetMaxValue())
		}
	}
}
