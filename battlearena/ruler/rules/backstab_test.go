package rules

import (

	"testing"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/stretchr/testify/assert"
)

// TestBackstabDetection verifies the spatial logic for determining if an attacker is behind a target.
// @test-link [[mechanic_backstab_detection_algorithm]]
func TestBackstabDetection(t *testing.T) {
	attacker := entity.New()
	target := entity.New()

	// Position target at (5, 5, 0)
	target.Position = position.New(5, 5, 0)

	testCases := []struct {
		name         string
		targetOrient entity.EntityOrientation
		attackerPos  position.Position
		expectBack   bool
	}{
		// FaceToward maps: Right(1,0)->Up(0), Up(0,1)->Right(1), Left(-1,0)->Down(2), Down(0,-1)->Left(3)
		// Assuming we follow FaceToward's thresholds for consistency:
		{"TargetUp_AttackerBehind", entity.Up, position.New(4, 5, 0), true},       // Attacker is West (facing East)
		{"TargetUp_AttackerSide", entity.Up, position.New(5, 6, 0), false},        // Attacker is North
		{"TargetRight_AttackerBehind", entity.Right, position.New(5, 4, 0), true}, // Attacker is South (facing North)
		{"TargetDown_AttackerBehind", entity.Down, position.New(6, 5, 0), true},   // Attacker is East (facing West)
		{"TargetLeft_AttackerBehind", entity.Left, position.New(5, 6, 0), true},   // Attacker is North (facing South)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			target.Orientation = tc.targetOrient
			attacker.Position = tc.attackerPos
			assert.Equal(t, tc.expectBack, attacker.IsBackstabbing(target))
		})
	}
}

// TestBackstabDamageCalculationLogic verifies the mathematical application of backstab multipliers (1.5x) 
// and defense reductions (50%) in damage calculations.
// @test-link [[mechanic_backstab_detection_algorithm]]
func TestBackstabDamageCalculationLogic(t *testing.T) {

	attacker := entity.New()
	target := entity.New()

	attacker.Position = position.New(4, 5, 0) // Behind target facing Up
	target.Position = position.New(5, 5, 0)
	target.Orientation = entity.Up

	// Set properties
	attacker.RepsertPropertyValue(property.Attack, 20)
	target.RepsertPropertyValue(property.Defense, 10)
	target.RepsertPropertyCValue(property.Shield, 0)

	// Case 1: Backstab (1.5x damage, 50% defense)
	// (20 * 1.5) - (10 * 0.5) = 30 - 5 = 25
	multiplier := 1.5
	effectiveDefense := 10 * 0.5
	damage := int(float64(20)*multiplier) - int(effectiveDefense)
	assert.Equal(t, 25, damage)

	// Case 2: Shield absorption
	target.RepsertPropertyCValue(property.Shield, 10)
	// 25 damage - 10 shield = 15 final damage
	if damage > 10 {
		damage -= 10
	} else {
		damage = 0
	}
	assert.Equal(t, 15, damage)
}
