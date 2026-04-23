package skillweight

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
)

// @test-link [[ISS-065_skill_weight_calculator]]
func TestSkillWeightCalculator(t *testing.T) {
	s1 := skill.New()
	
	// Add Damage property (+100 SW)
	s1.Effect.Properties = append(s1.Effect.Properties, defaultproperty.MakeIntProperty(property.Damage, 100, property.Public, property.Skill))
	// Add Critical Chance property (+50 SW)
	s1.Effect.Properties = append(s1.Effect.Properties, defaultproperty.MakeIntProperty(property.CriticalChance, 25, property.Public, property.Skill))
	// Set Range to 3 (+20 SW)
	s1.Targeting[property.Range.String()] = def.MakeRangeProperty(1, 3)

	pSW, nSW, _ := Calculate(s1)

	// Note: skill.New() comes with default TargetAnywhere (+40 SW)
	// So expected Positive SW = 100 (Damage) + 50 (Crit) + 20 (Range) + 40 (TargetAnywhere default) = 210
	expectedPSW := 210
	if pSW != expectedPSW {
		t.Errorf("Expected %d positive SW, got %d", expectedPSW, pSW)
	}

	// Note: skill.New() comes with default Delay (500) (-500 SW) and Cooldown (3 turns) (-75 SW)
	// So expected Negative SW = -500 - 75 = -575
	expectedNSW := -575
	if nSW != expectedNSW {
		t.Errorf("Expected %d negative SW, got %d", expectedNSW, nSW)
	}
}

// @test-link [[ISS-065_skill_grading_algorithm]]
func TestSkillGrading(t *testing.T) {
	testCases := []struct {
		sw    int
		grade string
	}{
		{100, "Grade I"},   // 0-150
		{250, "Grade II"},  // 151-300
		{400, "Grade III"}, // 301-500
		{600, "Grade IV"},  // 501-750
		{800, "Grade V"},   // 750+
	}

	for _, tc := range testCases {
		result := GetGrade(tc.sw)
		if result != tc.grade {
			t.Errorf("Expected %s for SW %d, got %s", tc.grade, tc.sw, result)
		}
	}
}

// @test-link [[ISS-065_credit_cost_formula]]
func TestCreditCostCalculation(t *testing.T) {
	testCases := []struct {
		positiveSW   int
		expectedCost int
	}{
		{100, 200}, // Basic attack
		{250, 500}, // Fireball
		{400, 800}, // Meteor Swarm
	}

	for _, tc := range testCases {
		result := GetCreditCost(tc.positiveSW)
		if result != tc.expectedCost {
			t.Errorf("Expected cost %d for positiveSW %d, got %d", tc.expectedCost, tc.positiveSW, result)
		}
	}
}
