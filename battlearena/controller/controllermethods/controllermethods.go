package controllermethods

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/controllervalidator"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue"
)

type SetValidatorAndQueue struct {
	actor.NoReply
	Validator controllervalidator.ControllerValidator
	Queue     *messagequeue.MessageQueue
}

type Send struct{}

type ReceiveAPIMessage struct {
	actor.NoReply
}
