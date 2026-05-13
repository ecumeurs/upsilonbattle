package controller

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/google/uuid"
)

// Controller is the base struct for all entity controllers.
// It handles actor plumbing (message routing, ruler communication) and optionally
// delegates AI decision-making to a pluggable Behavior.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type Controller struct {
	*actor.Actor
	ID             uuid.UUID
	// EntityID is the entity this controller is responsible for.
	// When using the one-controller-per-entity model, this is set at construction.
	EntityID       uuid.UUID
	Assigned       bool
	ControllerName string
	Ruler          actor.Communication
	// Behavior is the optional pluggable AI decision system.
	// When nil, the controller falls back to its own inline logic.
	// @spec-link [[mechanic_mech_behavior_layered]]
	Behavior       behavior.Behavior
}


// New
func New() *Controller {
	id := uuid.New()
	ctrl := &Controller{
		ID:    id,
		Actor: actor.New(fmt.Sprintf("Controller %s", id.String()[0:8])),

		Assigned: false,
	}

	ctrl.AddNotificationHandler(controllermethods.SetQueue{}, ctrl.setQueue, nil)
	ctrl.AddNotificationHandler(controllermethods.Send{}, ctrl.send, nil)
	ctrl.AddNotificationHandler(controllermethods.ReceiveAPIMessage{}, ctrl.receiveAPIMessage, nil)
	ctrl.AddReplyHandler(controllermethods.Send{}, ctrl.sendReply, nil)
	ctrl.Start()

	return ctrl
}

func NewController(id uuid.UUID) *Controller {
	ctrl := &Controller{
		ID:    id,
		Actor: actor.New(fmt.Sprintf("Controller %s", id.String()[0:8])),

		Assigned: false,
	}

	ctrl.AddNotificationHandler(controllermethods.SetQueue{}, ctrl.setQueue, nil)
	ctrl.AddNotificationHandler(controllermethods.Send{}, ctrl.send, nil)
	ctrl.AddNotificationHandler(controllermethods.ReceiveAPIMessage{}, ctrl.receiveAPIMessage, nil)
	ctrl.AddReplyHandler(controllermethods.Send{}, ctrl.sendReply, nil)
	ctrl.Start()

	return ctrl
}

// @spec-link [[mech_controller_handshake]]
func (c *Controller) setQueue(ctx actor.NotificationContext) {
	method := ctx.Msg.TargetMethod.(controllermethods.SetQueue)
	c.Ruler = method.Ruler
	fmt.Println("Controller: ", c.ID, " is now assigned to a player")
}

func (c *Controller) send(ctx actor.NotificationContext) {
	// this code depend on the type of controller and it's linked API Client
	fmt.Println("Controller: ", c.ID, ctx.Msg, " sent a message: ", ctx.Msg.TargetMethod)
}

func (c *Controller) sendReply(ctx actor.ReplyContext) {
	fmt.Println("Controller: ", c.ID, ctx.Msg, " sent a message (from reply): ", ctx.Msg.Content)
}

func (c *Controller) receiveAPIMessage(ctx actor.NotificationContext) {
}
