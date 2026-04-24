package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect/effectapplicator"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ProcessPositionalEffects fires all positional effects at a given position that match the trigger.
// It is called from Move (OnEnter/OnExit) and BeginingOfTurn (OnTurn).
//
// @spec-link [[mech_trigger_system]]
// @spec-link [[mech_positional_effects]]
func (gs *GameState) ProcessPositionalEffects(ent entity.Entity, pos position.Position, trigger property.TriggerTypeValue) {
	effectIDs, ok := gs.PositionalEffects[pos]
	if !ok || len(effectIDs) == 0 {
		return
	}

	log := gs.Logger.WithFields(logrus.Fields{
		"entityID": ent.ID.String()[0:8],
		"pos":      pos,
		"trigger":  trigger,
	})

	// Collect IDs to process (copy slice — we may remove during iteration)
	toProcess := make([]uuid.UUID, len(effectIDs))
	copy(toProcess, effectIDs)

	for _, effectID := range toProcess {
		eff, exists := gs.Effects[effectID]
		if !exists {
			continue
		}

		// Check trigger type matches
		triggerProp := eff.GetProperty(property.TriggerType)
		if triggerProp == nil {
			continue
		}
		effectTrigger := property.TriggerTypeValue(triggerProp.Get().(string))
		if effectTrigger != trigger {
			continue
		}

		log.WithField("effectID", effectID.String()[0:8]).Info("Positional effect triggered")

		// Check TriggerCount (0 = unlimited)
		removeAfter := false
		triggerCountProp := eff.GetProperty(property.TriggerCount)
		if triggerCountProp != nil {
			remaining := triggerCountProp.Get().(int)
			if remaining > 0 {
				remaining--
				triggerCountProp.Set(remaining)
				gs.Effects[effectID] = eff
				if remaining <= 0 {
					removeAfter = true
				}
			}
		}

		// Apply the effect
		gs.applyPositionalEffect(log, ent, eff)

		// RemoveOnTrigger: consume after one application
		removeOnTriggerProp := eff.GetProperty(property.RemoveOnTrigger)
		if removeAfter || (removeOnTriggerProp != nil && removeOnTriggerProp.Get().(bool)) {
			gs.RemovePositionalEffect(effectID, pos)
		}
	}
}

// applyPositionalEffect applies a single positional effect's properties directly to an entity.
// Positional effects differ from skill effects in that they have a single predetermined target
// (the entity that stepped on/occupies the cell). We use effectapplicator in single-target mode.
//
// @spec-link [[mech_positional_effects]]
func (gs *GameState) applyPositionalEffect(log *logrus.Entry, target entity.Entity, eff effect.Effect) {
	log.WithFields(logrus.Fields{
		"targetID": target.ID.String()[0:8],
		"effect":   eff.Name,
	}).Debug("Applying positional effect to entity")

	// ApplyDirectEffect expects a target position and targeted entities list.
	// Since positional effects are pre-targeted (they hit whoever is on the cell),
	// we pass the entity's current position and the entity itself as the single target.
	damaged, affected, _, _, _ := effectapplicator.ApplyDirectEffect(
		log,
		&target,
		eff,
		target.Position,
		[]position.Position{target.Position},
		gs.Grid,
		[]entity.Entity{target},
	)

	// Write back the modified entities
	gs.Entities[target.ID] = target
	for _, d := range damaged {
		gs.Entities[d.ID] = d
	}
	for _, a := range affected {
		gs.Entities[a.ID] = a
	}

	gs.IncVersion()
}
