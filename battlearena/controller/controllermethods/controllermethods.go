package controllermethods

import (
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/google/uuid"
)

// @spec-link [[api_controller_methods]]
type SetQueue struct {
	ControllerID uuid.UUID
	Ruler        actor.Communication
}

type SetQueueReply struct {
	ControllerID uuid.UUID
}

// and forget
type Send struct {
	actor.NoReply
}

type ReceiveAPIMessage struct {
	actor.NoReply
}
