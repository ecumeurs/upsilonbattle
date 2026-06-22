// Package battletest provides a deterministic, service-free scenario sandbox for
// driving the battle engine in unit tests. It formalizes the ad-hoc per-package
// helpers (makeGameStateForTwo*, FakeController) into a fluent builder so skill,
// movement, and positional-effect mechanics can be set up and asserted directly
// against the rules layer — no matchmaking, AI, or network.
//
// @spec-link [[module_skill_sandbox]]
package battletest

import (
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// FakeController is a no-op actor.Communication implementation that records the
// messages it receives, so tests can assert on notifications and replies.
type FakeController struct {
	ID             uuid.UUID
	NotifyMessages []*message.Message
	SentMessages   []*message.Message
}

// newFakeController creates a FakeController with a fresh ID and empty inboxes.
func newFakeController() *FakeController {
	return &FakeController{
		ID:             uuid.New(),
		NotifyMessages: make([]*message.Message, 0),
		SentMessages:   make([]*message.Message, 0),
	}
}

// NotifyActor records a fire-and-forget notification.
func (fc *FakeController) NotifyActor(m *message.Message) {
	fc.NotifyMessages = append(fc.NotifyMessages, m)
}

// SendActor records a request and replies asynchronously, mirroring the real actor contract.
func (fc *FakeController) SendActor(m *message.Message, replyTo chan *message.Message) {
	fc.SentMessages = append(fc.SentMessages, m)
	go func() {
		replyTo <- m.Reply()
	}()
}
