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

	// update entity in gamestate.
	gs.Entities[ent.ID] = ent
}
