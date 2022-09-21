package controllermethods

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/controllervalidator"
	"github.com/ecumeurs/upsilontools/tools/actor"
)

type SetValidatorAndQueue struct {
	actor.NoReply
	Validator controllervalidator.ControllerValidator
	Ruler     actor.Communication
}

// and forget
type Send struct {
	actor.NoReply
}

type ReceiveAPIMessage struct {
	actor.NoReply
}
