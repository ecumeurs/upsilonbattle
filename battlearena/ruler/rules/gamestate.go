package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type GameState struct {
	RulerID     uuid.UUID
	Grid        *grid.Grid
	Turner      turner.Turner
	Entities    map[uuid.UUID]entity.Entity
	Controllers map[uuid.UUID]actor.Communication
	Logger      *logrus.Entry
	// @spec-link [[rule_team_mechanics]]
	WinnerID uuid.UUID
	// @spec-link [[mech_game_state_versioning]]
	// @spec-link [[mech_version_bit_packing]]
	Version     int64
	TurnIndex   uint32
	ActionIndex uint32
}

func New(rulerID uuid.UUID) *GameState {
	gs := &GameState{
		RulerID:     rulerID,
		Turner:      turner.NewTurner(),
		Entities:    make(map[uuid.UUID]entity.Entity),
		Controllers: make(map[uuid.UUID]actor.Communication),
	}

	gs.Logger = logrus.WithFields(logrus.Fields{
		"component": "rules",
		"rulerID":   rulerID.String()[0:8],
	})

	return gs
}

func (gs *GameState) CheckControllerForEntity(controllerID uuid.UUID, entityID uuid.UUID) bool {
	for _, e := range gs.Entities {
		if e.ID == entityID {
			return e.ControllerID == controllerID
		}
	}
	return false
}

func (gs *GameState) UpdateVersion() {
	gs.Version = (int64(gs.TurnIndex) << 32) | int64(gs.ActionIndex)
}

func (gs *GameState) IncVersion() {
	gs.IncAction()
}

func (gs *GameState) IncAction() {
	gs.ActionIndex++
	gs.UpdateVersion()
}

func (gs *GameState) IncTurn() {
	gs.TurnIndex++
	gs.ActionIndex = 0
	gs.UpdateVersion()
}

func (gs *GameState) GetTurn() uint32 {
	return uint32(gs.Version >> 32)
}

func (gs *GameState) GetAction() uint32 {
	return uint32(gs.Version & 0xFFFFFFFF)
}
