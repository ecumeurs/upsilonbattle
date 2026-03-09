package controllers

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestFindNextAppropriateTarget(t *testing.T) {
	ctrl := NewRdAggressiveController("TestCtrl")
	ctrl.RequestLogger = ctrl.Logger
	for i := 0; i < 10; i++ {
		e := entity.Entity{
			ID:       uuid.New(),
			Position: position.Position{X: i * 10, Y: i * 3},
		}
		if i%2 == 0 {
			e.ControllerID = uuid.New()
		} else {
			e.ControllerID = ctrl.ID
		}
		ctrl.KnownEntities[e.ID] = e
	}

	// find first entity of this controller
	for _, e := range ctrl.KnownEntities {
		if e.ControllerID == ctrl.ID {
			ent, err := ctrl.selectNearestFoe(e, ctrl.KnownEntities)

			logrus.WithFields(logrus.Fields{
				"found_pos":        ent.Position,
				"found_entity":     ent.ID.String()[0:8],
				"found_controller": ent.ControllerID.String()[0:8]}).Info("found")
			if err != nil {
				t.Error(err)
			} else {
				if ent.ControllerID == ctrl.ID {
					t.Error("Found own entity")
				}
				if ent.ID == e.ID {
					t.Error("Found same entity")
				}
			}
		}
	}

}
