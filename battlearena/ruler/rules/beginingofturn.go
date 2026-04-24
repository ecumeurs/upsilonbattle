package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/sirupsen/logrus"
)

func (gs *GameState) BeginingOfTurn(ent entity.Entity) {
	loclog := gs.Logger.WithFields(logrus.Fields{
		"entityID": ent.ID.String()[0:8],
		"rule":     "BeginingOfTurn",
	})
	loclog.Debug("Begining of turn")

	// check stunned!

	stunvalue := ent.GetPropertyI(property.Stun).I()
	if stunvalue > ent.GetPropertyC(property.HP).GetMaxValue()/2 {
		// stunned!
		loclog.WithFields(logrus.Fields{
			"entityID": ent.ID.String()[0:8],
			"stun":     stunvalue,
		}).Debug("Entity is stunned")
		ent.UpdatePropertyValue(property.Stun, 0)
		ent.UpdatePropertyValue(property.HasActed, true)
		ent.UpdatePropertyValue(property.HasMoved, true)
		stunvalue = 0
	}
	if stunvalue > 0 {
		// Not stunned yet.
		stunvalue = stunvalue / 2
		if stunvalue <= 1 {
			stunvalue = 0
		}
		ent.UpdatePropertyValue(property.Stun, stunvalue)
		loclog.WithFields(logrus.Fields{
			"entityID": ent.ID.String()[0:8],
			"stun":     stunvalue,
		}).Debug("Entity is not stunned yet")

	}

	// @spec-link [[mech_trigger_system]]
	// Fire OnTurn positional effects for the entity's current cell.
	// This runs after stun resolution so the entity's final state is used.
	gs.ProcessPositionalEffects(ent, ent.Position, property.TriggerOnTurn)

	// Reload entity in case positional effects modified it
	ent = gs.Entities[ent.ID]

	// update entity in gamestate.
	gs.Entities[ent.ID] = ent
}
