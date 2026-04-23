package skillgenerator

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill/skillweight"
)

// @test-link [[ISS-065_skill_generation_balance]]
func TestGenerateRandomSkill(t *testing.T) {
	// Generate 100 random skills
	for i := 0; i < 100; i++ {
		sk := GenerateRandomSkill()
		
		pSW, nSW, netSW := skillweight.Calculate(sk)
		
		// 1. Verify net SW = 0
		if netSW != 0 {
			t.Errorf("Skill %d (%s) has non-zero net SW: %d (Positive: %d, Negative: %d)", i, sk.Name, netSW, pSW, nSW)
		}
		
		// 2. Verify grade distribution follows expected ranges
		grade := skillweight.GetGrade(pSW)
		if grade == "" {
			t.Errorf("Skill %d has empty grade", i)
		}
		
		// 3. Verify credit costs = positiveSW * 2
		cost := skillweight.GetCreditCost(pSW)
		if cost != pSW*2 {
			t.Errorf("Skill %d cost mismatch: expected %d, got %d", i, pSW*2, cost)
		}
	}
}
