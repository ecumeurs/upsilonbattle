package gamestate

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// FindCharacterInCell returns the first Character-type entity found in a list of entity IDs.
// This is the canonical way to resolve "who to attack" in a multi-entity cell, since only
// Character/Monster entities are valid attack targets — WalkThrough entities (traps, zones) are not.
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
// @spec-link [[mechanic_expiration_controller]]
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
	for pos, effectIDs := range gs.PositionalEffects {
		surviving := effectIDs[:0]
		for _, effectID := range effectIDs {
			if gs.shouldEffectExpireWithCaster(effectID, entityID) {
				delete(gs.Effects, effectID)
				continue
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

// shouldEffectExpireWithCaster determines if an effect should be removed because its caster was deleted.
func (gs *GameState) shouldEffectExpireWithCaster(effectID uuid.UUID, entityID uuid.UUID) bool {
	eff, exists := gs.Effects[effectID]
	return exists && eff.CasterID == entityID && eff.HasProperty(property.ExpiresWithCaster)
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
