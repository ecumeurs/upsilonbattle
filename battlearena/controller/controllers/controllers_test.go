package controllers

// @test-link [[mechanic_mech_behavior_layered]]

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TestFindNextAppropriateTarget verifies that selectNearestFoe returns an enemy entity, not an ally or the entity itself.
func TestFindNextAppropriateTarget(t *testing.T) {
	ctrl := NewRdAggressiveController("TestCtrl")
	ctrl.RequestLogger = ctrl.Logger

	populateMixedEntities(ctrl, 10)

	for _, e := range ctrl.KnownEntities {
		if e.ControllerID == ctrl.ID {
			assertNearestFoe(t, ctrl, e)
		}
	}
}

func populateMixedEntities(ctrl *AIController, count int) {
	for i := 0; i < count; i++ {
		e := entity.Entity{
			ID:       uuid.New(),
			Position: position.Position{X: i * 10, Y: i * 3},
		}
		e.Properties = make(map[string]property.Property)
		if i%2 == 0 {
			e.ControllerID = uuid.New()
			e.RepsertPropertyValue(property.TeamID, 1)
		} else {
			e.ControllerID = ctrl.ID
			e.RepsertPropertyValue(property.TeamID, 2)
		}
		ctrl.KnownEntities[e.ID] = e
	}
}

func assertNearestFoe(t *testing.T, ctrl *AIController, self entity.Entity) {
	t.Helper()
	ent, err := ctrl.selectNearestFoe(self, ctrl.KnownEntities)
	if err != nil {
		t.Error(err)
		return
	}
	logrus.WithFields(logrus.Fields{
		"found_pos":        ent.Position,
		"found_entity":     ent.ID.String()[0:8],
		"found_controller": ent.ControllerID.String()[0:8],
	}).Info("found")
	if ent.ControllerID == ctrl.ID {
		t.Error("Found own entity")
	}
	if ent.ID == self.ID {
		t.Error("Found same entity")
	}
}
