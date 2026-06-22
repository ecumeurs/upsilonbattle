package ruler
// @test-link [[mechanic_arena_lifecycle]]
// @test-link [[uc_combat_turn]]

import (
	"reflect"
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FakeController struct {
	*actor.Actor
	ID               uuid.UUID
	Inbox            chan *message.Message
	History          []*message.Message
	KnownEntities    map[uuid.UUID]entity.Entity
	ruler            actor.Communication
	battleready      bool
	receivedGrid     bool
	receivedEntities bool
}

// NewFake creates a new FakeController with a given name for testing purposes.
func NewFake(name string) *FakeController {
	ctrl := &FakeController{
		Actor:         actor.New(name),
		Inbox:         make(chan *message.Message, 10000),
		ID:            uuid.New(),
		KnownEntities: make(map[uuid.UUID]entity.Entity),
		battleready:   false,
	}

	ctrl.AddNotificationHandler(controllermethods.SetQueue{}, ctrl.SetQueue, nil)
	ctrl.AddNotificationHandler(controllermethods.Send{}, ctrl.Send, nil)
	ctrl.AddNotificationHandler(controllermethods.ReceiveAPIMessage{}, ctrl.ReceiveAPIMessage, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerNextTurn{}, ctrl.ControllerNextTurn, nil)
	ctrl.AddNotificationHandler(rulermethods.BattleStart{}, ctrl.BattleStart, nil)
	ctrl.AddNotificationHandler(rulermethods.BattleEnd{}, ctrl.BattleEnd, nil)
	ctrl.AddNotificationHandler(rulermethods.EntitiesStateChanged{}, ctrl.EntitiesStateChanged, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerAttacked{}, ctrl.ControllerAttacked, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerPassed{}, ctrl.NoOp, nil)
	ctrl.AddNotificationHandler(rulermethods.ControllerMoved{}, ctrl.NoOp, nil)

	ctrl.AddReplyHandler(rulermethods.GetStateReply{}, ctrl.GetStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.GetGridStateReply{}, ctrl.GetGridStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.GetEntitiesStateReply{}, ctrl.GetEntitiesStateReply, nil)
	ctrl.AddReplyHandler(rulermethods.ControllerMoveReply{}, ctrl.ControllerMoveReply, nil)
	ctrl.AddReplyHandler(rulermethods.ControllerAttackReply{}, ctrl.ControllerAttackReply, nil)
	ctrl.AddReplyHandler(rulermethods.EndOfTurn{}, ctrl.EndOfTurnReply, nil)

	ctrl.Start()

	return ctrl
}

// getMessageTypeName returns the string representation of the message's method type.
func getMessageTypeName(msg *message.Message) string {
	if msg.TargetMethod != nil {
		return reflect.TypeOf(msg.TargetMethod).String()
	}
	if msg.Content != nil {
		return reflect.TypeOf(msg.Content).String()
	}
	return ""
}

// ExpectMessage waits for a message of a specific type to arrive in the controller's inbox.
func (c *FakeController) ExpectMessage(t *testing.T, expectedType interface{}, timeout time.Duration) *message.Message {
	t.Helper()
	expectedTypeName := reflect.TypeOf(expectedType).String()

	for i, msg := range c.History {
		if getMessageTypeName(msg) == expectedTypeName {
			c.History = append(c.History[:i], c.History[i+1:]...)
			return msg
		}
	}

	start := time.Now()
	for {
		select {
		case msg := <-c.Inbox:
			if getMessageTypeName(msg) == expectedTypeName {
				return msg
			}
			c.History = append(c.History, msg)
		case <-time.After(timeout - time.Since(start)):
			t.Fatalf("Timeout waiting for message of type %s on %s", expectedTypeName, c.Name())
			return nil
		}
	}
}

// Close is a placeholder for closing the controller.
func (c *FakeController) Close() {
	// Not needed for Inbox approach
}

// triggerStopper is an internal helper to log and enqueue incoming messages.
func (c *FakeController) triggerStopper(msg *message.Message) {
	logrus.WithFields(logrus.Fields{
		"controller": c.Name(),
		"msgType":    getMessageTypeName(msg),
	}).Info("FakeController received message")
	c.Inbox <- msg
}

// SetQueue is a notification handler that assigns the ruler to the controller.
func (c *FakeController) SetQueue(ctx actor.NotificationContext) {
	c.triggerStopper(ctx.Msg)
	m := ctx.Msg.TargetMethod.(controllermethods.SetQueue)
	c.ID = m.ControllerID
	c.ruler = m.Ruler
	c.ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), c.GetCallbackChan())
	c.ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), c.GetCallbackChan())
}

// Send is a notification handler that captures outgoing messages.
func (c *FakeController) Send(ctx actor.NotificationContext) {
	c.triggerStopper(ctx.Msg)
}

// ReceiveAPIMessage is a notification handler for incoming API messages.
func (c *FakeController) ReceiveAPIMessage(ctx actor.NotificationContext) {
	c.triggerStopper(ctx.Msg)
}

// ControllerNextTurn is a notification handler for the turn transition signal.
func (c *FakeController) ControllerNextTurn(ctx actor.NotificationContext) {
	m := ctx.Msg.TargetMethod.(rulermethods.ControllerNextTurn)
	if m.Entity.ControllerID == c.ID {
		c.triggerStopper(ctx.Msg)
	}
}

// BattleStart is a notification handler for the battle initiation signal.
func (c *FakeController) BattleStart(ctx actor.NotificationContext) {
	c.triggerStopper(ctx.Msg)
}

// BattleEnd is a notification handler for the battle termination signal.
func (c *FakeController) BattleEnd(ctx actor.NotificationContext) {
	c.triggerStopper(ctx.Msg)
}

// ControllerAttacked is a notification handler for when the controller's entity is attacked.
func (c *FakeController) ControllerAttacked(ctx actor.NotificationContext) {
	c.triggerStopper(ctx.Msg)
}

// EntitiesStateChanged is a notification handler for bulk entity state updates.
func (c *FakeController) EntitiesStateChanged(ctx actor.NotificationContext) {
	c.triggerStopper(ctx.Msg)
	c.RequestLogger.WithFields(logrus.Fields{
		"Turn": ctx.Msg.TargetMethod.(rulermethods.EntitiesStateChanged).Turn.String()}).Info("New Turn Received")
	c.KnownEntities = make(map[uuid.UUID]entity.Entity)
	for _, e := range ctx.Msg.TargetMethod.(rulermethods.EntitiesStateChanged).Entities {
		c.KnownEntities[e.ID] = e
	}
}

// GetStateReply is a reply handler for the full game state.
func (c *FakeController) GetStateReply(ctx actor.ReplyContext) {
	c.triggerStopper(ctx.Msg)
}

// checkReadiness verifies if the controller has received all necessary data to start.
func (c *FakeController) checkReadiness() {
	if c.receivedGrid && c.receivedEntities && !c.battleready && c.ruler != nil {
		c.battleready = true
		c.ruler.NotifyActor(message.Create(nil, rulermethods.ControllerBattleReady{
			ControllerID: c.ID,
		}, nil))
	}
}

// GetGridStateReply is a reply handler for the grid state.
func (c *FakeController) GetGridStateReply(ctx actor.ReplyContext) {
	c.triggerStopper(ctx.Msg)
	c.receivedGrid = true
	c.checkReadiness()
}

// GetEntitiesStateReply is a reply handler for the entities state.
func (c *FakeController) GetEntitiesStateReply(ctx actor.ReplyContext) {
	c.triggerStopper(ctx.Msg)
	c.receivedEntities = true
	c.checkReadiness()
}

// ControllerMoveReply is a reply handler for a movement request.
func (c *FakeController) ControllerMoveReply(ctx actor.ReplyContext) {
	c.triggerStopper(ctx.Msg)
}

// ControllerAttackReply is a reply handler for an attack request.
func (c *FakeController) ControllerAttackReply(ctx actor.ReplyContext) {
	c.triggerStopper(ctx.Msg)
}

// EndOfTurnReply is a reply handler for a turn end request.
func (c *FakeController) EndOfTurnReply(ctx actor.ReplyContext) {
	c.triggerStopper(ctx.Msg)
}

// NoOp is an empty notification handler.
func (c *FakeController) NoOp(ctx actor.NotificationContext) {}
