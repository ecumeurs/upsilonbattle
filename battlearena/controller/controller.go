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
	Uuid           uuid.UUID
	Assigned       bool
	ControllerName string
	Validator      controllervalidator.ControllerValidator
	Queue          *messagequeue.MessageQueue
}

// New
func New() *Controller {
	ctrl := &Controller{
		Uuid:     uuid.New(),
		Assigned: false,
	}
	ctrl.act.SetReceiveMessageHandler(ctrl.handleMessage)
	ctrl.act.SetReplyMessageHandler(ctrl.handleReply)
	return ctrl
}

func (c *Controller) handleMessage(msg message.Message) {
	switch msg.TargetMethod.(type) {
	case controllermethods.SetValidatorAndQueue:
		c.SetValidatorAndQueue(msg, msg.TargetMethod.(controllermethods.SetValidatorAndQueue))
	case controllermethods.Send:
		c.Send(msg)
	case controllermethods.ReceiveAPIMessage:
		c.ReceiveAPIMessage(msg)
	}
}

func (c *Controller) handleReply(msg message.Message) {
	switch msg.TargetMethod.(type) {
	case controllermethods.Send:
		c.Send(msg)
	}
}

func (c *Controller) SetValidatorAndQueue(msg message.Message, method controllermethods.SetValidatorAndQueue) {
	c.Validator = method.Validator
	c.Queue = method.Queue
	c.act.NoReply(msg)
}

func (c *Controller) Send(msg message.Message) {
	// this code depend on the type of controller and it's linked API Client
}

func (c *Controller) ReceiveAPIMessage(msg message.Message) {
	c.act.NoReply(msg)
}
