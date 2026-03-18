package controllermethods

import (
	"github.com/ecumeurs/upsilontools/tools/actor"
)

// @spec-link [[api_controller_methods]]
type SetQueue struct {
	actor.NoReply
	Ruler     actor.Communication
}

// and forget
type Send struct {
	actor.NoReply
}

type ReceiveAPIMessage struct {
	actor.NoReply
}
