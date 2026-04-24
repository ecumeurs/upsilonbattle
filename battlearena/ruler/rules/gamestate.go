package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect"
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
	// @spec-link [[mech_version_bit_packing]]
	Version     int64
	TurnIndex   uint32
	ActionIndex uint32

	// @spec-link [[mech_positional_effects]]
	// PositionalEffects maps grid positions to the effect IDs active at that cell.
	// Actual effect data is stored in Effects for single-source-of-truth.
	PositionalEffects map[position.Position][]uuid.UUID
	// @spec-link [[mech_positional_effects]]
	// Effects is the central store of all positional effect data, keyed by effect ID.
	Effects map[uuid.UUID]effect.Effect
}

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

// FindCharacterInCell returns the first Character-type entity found in a list of entity IDs.
// This is the canonical way to resolve "who to attack" in a multi-entity cell, since only
// Character/Monster entities are valid attack targets — WalkThrough entities (traps, zones) are not.
//
// @spec-link [[mechanic_multi_entity_cell_system]]
func (gs *GameState) FindCharacterInCell(entityIDs []uuid.UUID) (entity.Entity, bool) {
	for _, id := range entityIDs {
		if ent, ok := gs.Entities[id]; ok {
			if ent.Type == entity.Character || ent.Type == entity.Monster {
				return ent, true
			}
		}
	}
	return entity.Entity{}, false
}

// HasBlockingEntity returns true if any entity in the given ID list has WalkThrough=false.
// A cell with at least one blocking entity cannot be entered by another entity.
//
// @spec-link [[mechanic_multi_entity_cell_system]]
func (gs *GameState) HasBlockingEntity(entityIDs []uuid.UUID, selfID uuid.UUID) bool {
	for _, id := range entityIDs {
		if id == selfID {
			continue
		}
		if ent, ok := gs.Entities[id]; ok {
			walkThroughProp := ent.GetProperty(property.WalkThrough)
			// Absent WalkThrough property = regular entity = blocking (WalkThrough defaults to false).
			if walkThroughProp == nil {
				return true
			}
			isWalkThrough, ok := walkThroughProp.Get().(bool)
			if !ok || !isWalkThrough {
				return true
			}
		}
	}
	return false
}

// FindEffectsByCaster returns all effect IDs (and their Effects) that belong to the given caster.
// Useful when cleaning up after a caster dies or is removed mid-combat.
//
// @spec-link [[mechanic_effect_caster_tracking]]
func (gs *GameState) FindEffectsByCaster(casterID uuid.UUID) map[uuid.UUID]effect.Effect {
	result := make(map[uuid.UUID]effect.Effect)
	for id, eff := range gs.Effects {
		if eff.CasterID == casterID {
			result[id] = eff
		}
	}
	return result
}

// RemoveEntity removes a temporary entity from all game state structures atomically.
// It removes the entity from the Turner, the Grid, the Entities map, and cleans up
// any positional effects owned by this entity (where ExpiresWithCaster is true).
//
// @spec-link [[mech_entity_expiration]]
// @spec-link [[mechanic_effect_caster_tracking]]
func (gs *GameState) RemoveEntity(entityID uuid.UUID) {
	gs.Logger.WithFields(logrus.Fields{
		"entityID": entityID.String()[0:8],
	}).Info("Removing entity from game state")

	// Remove from Turner (handles the case where it is the current turn entity)
	gs.Turner.RemoveEntity(entityID)

	// Remove from Grid
	if ent, exists := gs.Entities[entityID]; exists {
		gs.Grid.RemoveEntity(ent.Position, entityID)
	}

	// Cleanup PositionalEffects owned by this entity (ExpiresWithCaster)
	// @spec-link [[mechanic_effect_caster_tracking]]
	for pos, effectIDs := range gs.PositionalEffects {
		surviving := effectIDs[:0]
		for _, effectID := range effectIDs {
			if eff, exists := gs.Effects[effectID]; exists {
				if eff.CasterID == entityID {
					// Check ExpiresWithCaster flag on the effect
					if eff.HasProperty(property.ExpiresWithCaster) {
						delete(gs.Effects, effectID)
						continue // drop this effectID
					}
				}
			}
			surviving = append(surviving, effectID)
		}
		if len(surviving) == 0 {
			delete(gs.PositionalEffects, pos)
		} else {
			gs.PositionalEffects[pos] = surviving
		}
	}

	// Remove from Entities map
	delete(gs.Entities, entityID)
}

// RemovePositionalEffect removes a single positional effect from a cell and the Effects store.
//
// @spec-link [[mech_positional_effects]]
func (gs *GameState) RemovePositionalEffect(effectID uuid.UUID, pos position.Position) {
	delete(gs.Effects, effectID)

	if ids, ok := gs.PositionalEffects[pos]; ok {
		updated := ids[:0]
		for _, id := range ids {
			if id != effectID {
				updated = append(updated, id)
			}
		}
		if len(updated) == 0 {
			delete(gs.PositionalEffects, pos)
		} else {
			gs.PositionalEffects[pos] = updated
		}
	}

	// Also remove from cell's EffectIDs
	if c, ok := gs.Grid.CellAt(pos); ok {
		c.RemoveEffect(effectID)
	}
}

func (gs *GameState) GetEntities() map[uuid.UUID]entity.Entity {
	return gs.Entities
}

func (gs *GameState) GetGrid() *grid.Grid {
	return gs.Grid
}

func (gs *GameState) GetLogger() *logrus.Entry {
	return gs.Logger
}
