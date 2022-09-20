package controller

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/controllervalidator"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type Controller struct {
	act            *actor.Actor
	ID             uuid.UUID
	Assigned       bool
	ControllerName string
	Validator      controllervalidator.ControllerValidator
	Queue          *messagequeue.MessageQueue
}

// New
func New() *Controller {
	ctrl := &Controller{
		ID: uuid.New(),

		Assigned: false,
	}

	ctrl.act = actor.New(fmt.Sprintf("Controller %s", ctrl.ID.String()[0:8]))
	ctrl.act.SetReceiveMessageHandler(ctrl.handleMessage)
	ctrl.act.SetReplyMessageHandler(ctrl.handleReply)
	ctrl.act.Start()

	return ctrl
}

func (c *Controller) handleMessage(msg message.Message) bool {
	switch msg.TargetMethod.(type) {
	case controllermethods.SetValidatorAndQueue:
		c.setValidatorAndQueue(msg, msg.TargetMethod.(controllermethods.SetValidatorAndQueue))
		return true
	case controllermethods.Send:
		c.send(msg)
		return true
	case controllermethods.SendExpectReply:
		c.sendAndExpectReply(msg)
		return true
	case controllermethods.ReceiveAPIMessage:
		c.receiveAPIMessage(msg)
		return true
	}
	return false
}

func (c *Controller) handleReply(msg message.Message) bool {
	switch msg.TargetMethod.(type) {
	case controllermethods.Send:
		c.send(msg)
		return true
	}
	return false
}

func (c *Controller) setValidatorAndQueue(msg message.Message, method controllermethods.SetValidatorAndQueue) {
	c.Validator = method.Validator
	c.Queue = method.Queue
	c.act.NoReply(msg.Reply())
	fmt.Println("Controller: ", c.ID, " is now assigned to a player")
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
	c.act.NoReply(msg.Reply())
	fmt.Println("Controller: ", c.ID, msg, " sent a message: ", msg.Content)
}

func (c *Controller) sendAndExpectReply(msg message.Message) {
	// this code depend on the type of controller and it's linked API Client
	c.act.NoReply(msg.Reply()) // TODO: change this: this should expect some reply ...
}

func (c *Controller) receiveAPIMessage(msg message.Message) {
	c.act.NoReply(msg.Reply())
}
