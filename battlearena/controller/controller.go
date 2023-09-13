package controller

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type Controller struct {
	*actor.Actor
	ID             uuid.UUID
	Assigned       bool
	ControllerName string
	Ruler          actor.Communication
}

// New
func New() *Controller {
	id := uuid.New()
	ctrl := &Controller{
		ID:    id,
		Actor: actor.New(fmt.Sprintf("Controller %s", id.String()[0:8])),

		Assigned: false,
	}

	ctrl.AddMethod(controllermethods.SetQueue{}, ctrl.setQueue, nil)
	ctrl.AddMethod(controllermethods.Send{}, ctrl.send, nil)
	ctrl.AddMethod(controllermethods.ReceiveAPIMessage{}, ctrl.receiveAPIMessage, nil)
	ctrl.AddReply(controllermethods.Send{}, ctrl.send, nil)
	ctrl.Start()

	return ctrl
}

func (c *Controller) setQueue(msg *message.Message) bool {
	method := msg.Content.(controllermethods.SetQueue)
	c.Ruler = method.Ruler
	c.NoReply(msg)
	fmt.Println("Controller: ", c.ID, " is now assigned to a player")
	return true
}

func (c *Controller) send(msg *message.Message) bool {
	// this code depend on the type of controller and it's linked API Client
	c.NoReply(msg)
	fmt.Println("Controller: ", c.ID, msg, " sent a message: ", msg.Content)
	return true
}

func (c *Controller) receiveAPIMessage(msg *message.Message) bool {
	c.NoReply(msg)
	return true
}
