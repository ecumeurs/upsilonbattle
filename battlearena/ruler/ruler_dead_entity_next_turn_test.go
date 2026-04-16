package ruler

// @test-link [[mech_initiative_active_state]]

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TestRulerNextTurnSkipsDeadEntity is a regression test for ISS-046.
//
// Scenario: Entity A (ctrl1) kills Entity B (ctrl2) during its turn.
// Entity B happened to be the very next entity in the turn queue.
// After the kill, A calls EndOfTurn.
// The ruler must skip the now-dead B and hand the next turn to the living
// Entity C (also ctrl2), NOT hang waiting for a dead entity to respond.
//
// Without the dead-entity skip loop in ruler.endOfTurn, the game permanently
// stalls because ControllerNextTurn is sent to uuid.Nil (zero-value
// ControllerID of the deleted entity), no controller receives it, and the
// shot clock reissues the EndOfTurn against the same dead entity again.
func TestRulerNextTurnSkipsDeadEntity(t *testing.T) {
	logrus.SetLevel(logrus.WarnLevel) // suppress noise; flip to DebugLevel to diagnose failures

	// ── Ruler setup ────────────────────────────────────────────────────────

	r := NewRuler(uuid.New())
	defer r.Stop()
	// Disable the shot clock so the test does not depend on wall-clock timing.
	r.ShotClockDuration = 0

	r.NbControllers = 2
	r.NbEntitiesPerController = 1 // 1 entity per side so the queue is deterministic

	r.GameState.Grid = grid.NewGrid(10, 10, 3)

	ctrl1 := NewFake("Attacker-Controller")
	ctrl2 := NewFake("Victim-Controller")

	// ── Entity A: the attacker (ctrl1, Team 1) ─────────────────────────────
	//   Delay = 100 → goes first.
	entA := entity.New()
	entA.ControllerID = ctrl1.ID
	entA.Type = entity.Character
	entA.CurrentDelay = 100
	attackerPos := position.Position{X: 5, Y: 5, Z: 3}
	entA.Position = attackerPos
	for _, v := range def.PropertiesForCharacter() {
		entA.Properties[v.Name(property.GameMaster)] = v
	}
	entA.RepsertPropertyValue(property.TeamID, 1)
	// Give the attacker a very high Attack so the one-hit kill is guaranteed.
	// Default Defense for characters is typically 0–5, so setting Attack=9999
	// ensures computedDamage = Max(1, 9999-defense) >> any default HP value.
	entA.RepsertPropertyValue(property.Attack, 9999)
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), attackerPos, entA.ID)
	r.GameState.Entities[entA.ID] = entA
	r.GameState.Turner.AddEntity(entA.ID, entA.CurrentDelay) // delay 100 → first

	// ── Entity B: the victim (ctrl2, Team 2) ──────────────────────────────
	//   Delay = 200 → second in queue (will be next after A ends its turn).
	entB := entity.New()
	entB.ControllerID = ctrl2.ID
	entB.Type = entity.Character
	entB.CurrentDelay = 200
	victimPos := position.Position{X: 5, Y: 6, Z: 3} // adjacent to attacker
	entB.Position = victimPos
	for _, v := range def.PropertiesForCharacter() {
		entB.Properties[v.Name(property.GameMaster)] = v
	}
	entB.RepsertPropertyValue(property.TeamID, 2)
	// Force HP to 1 so a single attack kills it regardless of default stats.
	entB.RepsertPropertyValue(property.HP, 1)
	entB.RepsertPropertyValue(property.Defense, 0)
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), victimPos, entB.ID)
	r.GameState.Entities[entB.ID] = entB
	r.GameState.Turner.AddEntity(entB.ID, entB.CurrentDelay) // delay 200 → second

	// ── Entity C: the survivor (ctrl2, Team 2) ────────────────────────────
	//   Delay = 300 → third in queue. This is the entity that SHOULD receive
	//   the ControllerNextTurn after B is killed and A ends its turn.
	entC := entity.New()
	entC.ControllerID = ctrl2.ID
	entC.Type = entity.Character
	entC.CurrentDelay = 300
	survivorPos := position.Position{X: 3, Y: 3, Z: 3}
	entC.Position = survivorPos
	for _, v := range def.PropertiesForCharacter() {
		entC.Properties[v.Name(property.GameMaster)] = v
	}
	entC.RepsertPropertyValue(property.TeamID, 2)
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), survivorPos, entC.ID)
	r.GameState.Entities[entC.ID] = entC
	r.GameState.Turner.AddEntity(entC.ID, entC.CurrentDelay) // delay 300 → third

	// Snapshot IDs for assertions (entities are value types, so capture now).
	attackerID := entA.ID
	victimID := entB.ID
	survivorID := entC.ID

	// ── Register controllers ───────────────────────────────────────────────
	{
		dChan := make(chan *message.Message, 1)
		r.SendActor(message.Create(nil, rulermethods.AddController{
			Controller:   ctrl1,
			ControllerID: ctrl1.ID,
		}, nil), dChan)
		<-dChan
	}
	{
		dChan := make(chan *message.Message, 1)
		r.SendActor(message.Create(nil, rulermethods.AddController{
			Controller:   ctrl2,
			ControllerID: ctrl2.ID,
		}, nil), dChan)
		<-dChan
	}

	// ── Wait for BattleStart ───────────────────────────────────────────────
	ctrl1.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	// ── Wait for the first ControllerNextTurn (should go to ctrl1 / entA) ─
	var firstTurnMsg *message.Message
	timeout := time.After(5 * time.Second)
waitFirst:
	for {
		select {
		case msg := <-ctrl1.Inbox:
			if _, ok := msg.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				firstTurnMsg = msg
				break waitFirst
			}
		case <-ctrl2.Inbox:
			// ctrl2 messages during startup (e.g. EntitiesStateChanged) — ignore.
		case <-timeout:
			t.Fatal("Timeout waiting for first ControllerNextTurn (expected Attacker's turn)")
		}
	}

	firstTurn := firstTurnMsg.TargetMethod.(rulermethods.ControllerNextTurn)
	if firstTurn.Entity.ID != attackerID {
		t.Fatalf("Expected first turn to go to Attacker (%s), got %s",
			attackerID.String()[0:8], firstTurn.Entity.ID.String()[0:8])
	}

	// ── Attack: Attacker kills Victim ──────────────────────────────────────
	attackReplyChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.ControllerAttack{
		EntityID:     attackerID,
		ControllerID: ctrl1.ID,
		Target:       victimPos,
	}, rulermethods.ControllerAttackReply{}), attackReplyChan)

	attackReply := <-attackReplyChan
	if attackReply.HasError {
		t.Fatalf("Attack failed unexpectedly: %s", attackReply.ErrorKey)
	}

	// Verify victim is actually gone from game state.
	// We query entities state directly from the ruler to confirm.
	entityStateChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), entityStateChan)
	entityStateMsg := <-entityStateChan
	entityState := entityStateMsg.Content.(rulermethods.GetEntitiesStateReply)
	for _, e := range entityState.Entities {
		if e.ID == victimID {
			t.Fatal("Victim (Entity B) should have been removed from game state after lethal attack, but it still exists")
		}
	}

	// Drain any pending EntitiesStateChanged notifications before EndOfTurn.
	drainInbox(ctrl1, ctrl2)

	// ── End of turn: this is the critical moment ───────────────────────────
	// Without the fix, calling EndOfTurn here causes the ruler to hand the
	// next turn to the dead victimID (which is at the front of the queue),
	// find it absent from gs.Entities, and stall forever.
	eotReplyChan := make(chan *message.Message, 1)
	r.SendActor(message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     attackerID,
		ControllerID: ctrl1.ID,
	}, rulermethods.EndOfTurn{}), eotReplyChan)
	<-eotReplyChan

	// ── Assert: next ControllerNextTurn must go to the living Survivor ─────
	// We give the game 3 seconds to deliver the turn. If it times out, the
	// bug is present (the game hung trying to hand a turn to the dead victim).
	timeout2 := time.After(3 * time.Second)
waitNext:
	for {
		select {
		case msg := <-ctrl2.Inbox:
			if nxt, ok := msg.TargetMethod.(rulermethods.ControllerNextTurn); ok {
				// The dead victim must never receive the turn signal.
				if nxt.Entity.ID == victimID {
					t.Fatalf("REGRESSION ISS-046: Ruler handed the next turn to the dead entity (victimID=%s). The dead-entity skip guard is not working.",
						victimID.String()[0:8])
				}
				// The survivor must receive it.
				if nxt.Entity.ID != survivorID {
					t.Fatalf("Expected next turn to go to Survivor (%s), got %s",
						survivorID.String()[0:8], nxt.Entity.ID.String()[0:8])
				}
				break waitNext // ✅ Correct behaviour confirmed.
			}
		case <-ctrl1.Inbox:
			// ignore other ctrl1 messages (EntitiesStateChanged etc.)
		case <-timeout2:
			t.Fatal("REGRESSION ISS-046: Timeout – no ControllerNextTurn received after killing the next-in-queue entity and ending the turn. The game is hung (dead entity turn assignment).")
		}
	}

	ctrl1.Stop()
	ctrl2.Stop()
}

// drainInbox consumes all messages currently in the Inbox channels of the
// given controllers without blocking. Used to reset state between phases.
func drainInbox(ctrls ...*FakeController) {
	for _, c := range ctrls {
		for {
			select {
			case <-c.Inbox:
			default:
				return // channel empty
			}
		}
	}
}
