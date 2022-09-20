package battlearena

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type BattleArena struct {
	act         actor.Actor
	Uuid        uuid.UUID
	Controllers map[uuid.UUID]*controller.Controller
	Ruler       ruler.Ruler
}

func NewBattleArena() BattleArena {
	ba := BattleArena{
		Uuid:        uuid.New(),
		Controllers: make(map[uuid.UUID]*controller.Controller),
		Ruler:       ruler.NewRuler(),
	}
	ba.act.SetReceiveMessageHandler(ba.handleMessage)
	return ba
}

func (b *BattleArena) handleMessage(msg message.Message) bool {
	switch msg.TargetMethod.(type) {

	}
	return false
}
