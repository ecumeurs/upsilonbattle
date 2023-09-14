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
