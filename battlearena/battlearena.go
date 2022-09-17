package battlearena

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/google/uuid"
)

type BattleArena struct {
	Uuid        uuid.UUID
	Grid        *grid.Grid
	Controllers map[uuid.UUID]*controller.Controller
}
