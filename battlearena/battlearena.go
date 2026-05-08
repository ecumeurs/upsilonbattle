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
	Ruler       *ruler.Ruler
	Metadata    map[string]interface{}
}

// NewBattleArena initializes a new battle arena instance with the specified UUID.
// It creates a new actor with the given UUID and bootstraps the underlying ruler
// and controller maps. This function is the primary entry point for match instantiation
// in the Go engine, orchestrating the creation of the core Ruler component.
// @spec-link [[api_go_battle_start]]
// @spec-link [[mechanic_mech_arena_lifecycle]]
func NewBattleArena(id uuid.UUID) *BattleArena {
	ba := BattleArena{
		Actor:       actor.New("BattleArena"),
		Uuid:        id,
		Controllers: make(map[uuid.UUID]*controller.Controller),
		Ruler:       ruler.NewRuler(id),
		Metadata:    make(map[string]interface{}),
	}
	return &ba
}
