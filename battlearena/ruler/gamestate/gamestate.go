package gamestate

import (
	"github.com/ecumeurs/upsilontypes/entity"

	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
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
	WinnerTeamID int
	// @spec-link [[mech_game_state_versioning]]
	Version     int64
	TurnIndex   uint32
	ActionIndex uint32

	// PositionalEffects maps grid positions to the effect IDs active at that cell.
	// Actual effect data is stored in Effects for single-source-of-truth.
	PositionalEffects map[position.Position][]uuid.UUID
	// Effects is the central store of all positional effect data, keyed by effect ID.
	Effects map[uuid.UUID]effect.Effect
}

// New initializes a new GameState instance with the given ruler ID.
// It sets up the logger, versioning, and internal maps for entities and effects.
func New(rulerID uuid.UUID) *GameState {
	gs := &GameState{
		RulerID:           rulerID,
		Turner:            turner.NewTurner(),
		Entities:          make(map[uuid.UUID]entity.Entity),
		Controllers:       make(map[uuid.UUID]actor.Communication),
		PositionalEffects: make(map[position.Position][]uuid.UUID),
		Effects:           make(map[uuid.UUID]effect.Effect),
	}

	gs.Logger = logrus.WithFields(logrus.Fields{
		"component": "rules",
		"rulerID":   rulerID.String()[0:8],
	})

	return gs
}

// CheckControllerForEntity verifies if a given controller ID is authorized to command an entity.
// Returns true if authorized, false otherwise.
func (gs *GameState) CheckControllerForEntity(controllerID uuid.UUID, entityID uuid.UUID) bool {
	for _, e := range gs.Entities {
		if e.ID == entityID {
			return e.ControllerID == controllerID
		}
	}
	return false
}

// GetEntities returns the map of all active entities in the game state.
func (gs *GameState) GetEntities() map[uuid.UUID]entity.Entity {
	return gs.Entities
}

// GetGrid returns the underlying grid structure.
func (gs *GameState) GetGrid() *grid.Grid {
	return gs.Grid
}

// GetLogger returns the GameState's logger instance.
func (gs *GameState) GetLogger() *logrus.Entry {
	return gs.Logger
}
