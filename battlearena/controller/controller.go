package controller

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/controllervalidator"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type Controller struct {
	act            actor.Actor
	ID             uuid.UUID
	Assigned       bool
	ControllerName string
	Validator      controllervalidator.ControllerValidator
	Queue          *messagequeue.MessageQueue
}

// New
func New() *Controller {
	ctrl := &Controller{
		ID:       uuid.New(),
		Assigned: false,
	}
	ctrl.act.SetReceiveMessageHandler(ctrl.handleMessage)
	ctrl.act.SetReplyMessageHandler(ctrl.handleReply)
	return ctrl
}

func (c *Controller) handleMessage(msg message.Message) {
	switch msg.TargetMethod.(type) {
	case controllermethods.SetValidatorAndQueue:
		c.setValidatorAndQueue(msg, msg.TargetMethod.(controllermethods.SetValidatorAndQueue))
	case controllermethods.Send:
		c.send(msg)
	case controllermethods.SendExpectReply:
		c.sendAndExpectReply(msg)
	case controllermethods.ReceiveAPIMessage:
		c.receiveAPIMessage(msg)
	}
}

func (c *Controller) handleReply(msg message.Message) {
	switch msg.TargetMethod.(type) {
	case controllermethods.Send:
		c.send(msg)
	}
}

func (c *Controller) setValidatorAndQueue(msg message.Message, method controllermethods.SetValidatorAndQueue) {
	c.Validator = method.Validator
	c.Queue = method.Queue
	c.act.NoReply(msg)
}

// Notify Actor
func (c *Controller) NotifyActor(msg message.Message) {
	c.act.SendMessageAndForget(msg)
}

func (c *Controller) SendActor(msg message.Message, callback chan message.Message) {
	c.act.SendMessage(msg, callback)
}

func (c *Controller) send(msg message.Message) {
	// this code depend on the type of controller and it's linked API Client
	c.act.NoReply(msg)
}

func (c *Controller) sendAndExpectReply(msg message.Message) {
	// this code depend on the type of controller and it's linked API Client
	c.act.NoReply(msg) // TODO: change this: this should expect some reply ...
}

func (c *Controller) receiveAPIMessage(msg message.Message) {
	c.act.NoReply(msg)
}
