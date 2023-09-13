package battlearena

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/google/uuid"
)

type BattleArena struct {
	*actor.Actor
	Uuid        uuid.UUID
	Controllers map[uuid.UUID]*controller.Controller
	Ruler       ruler.Ruler
}

func NewBattleArena() BattleArena {
	ba := BattleArena{
		Actor:       actor.New("BattleArena"),
		Uuid:        uuid.New(),
		Controllers: make(map[uuid.UUID]*controller.Controller),
		Ruler:       ruler.NewRuler(),
	}
	return ba
}
