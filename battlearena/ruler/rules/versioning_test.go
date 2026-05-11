package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"testing"
	"github.com/google/uuid"
)

// @test-link [[api_standard_envelope]]
// @test-link [[rule_credit_action_communication_layer]]

// TestVersioningBitPacking verifies the internal bit-packing logic for game state versioning.
// It ensures that turn indices and action indices are correctly encoded into a single 64-bit version integer.
func TestVersioningBitPacking(t *testing.T) {
	gs := gamestate.New(uuid.New())

	// Initial State
	if gs.Version != 0 {
		t.Errorf("Expected version 0, got %d", gs.Version)
	}

	// First Turn
	gs.IncTurn()
	if gs.TurnIndex != 1 || gs.ActionIndex != 0 {
		t.Errorf("Expected Turn 1, Action 0. Got Turn %d, Action %d", gs.TurnIndex, gs.ActionIndex)
	}
	expectedV1 := int64(1) << 32
	if gs.Version != expectedV1 {
		t.Errorf("Expected version %d, got %d", expectedV1, gs.Version)
	}

	// First Action
	gs.IncAction()
	if gs.TurnIndex != 1 || gs.ActionIndex != 1 {
		t.Errorf("Expected Turn 1, Action 1. Got Turn %d, Action %d", gs.TurnIndex, gs.ActionIndex)
	}
	expectedV11 := (int64(1) << 32) | 1
	if gs.Version != expectedV11 {
		t.Errorf("Expected version %d, got %d", expectedV11, gs.Version)
	}

	// Getters
	if gs.GetTurn() != 1 {
		t.Errorf("Expected GetTurn() 1, got %d", gs.GetTurn())
	}
	if gs.GetAction() != 1 {
		t.Errorf("Expected GetAction() 1, got %d", gs.GetAction())
	}

	// Second Turn
	gs.IncTurn()
	if gs.TurnIndex != 2 || gs.ActionIndex != 0 {
		t.Errorf("Expected Turn 2, Action 0. Got Turn %d, Action %d", gs.TurnIndex, gs.ActionIndex)
	}
	expectedV2 := int64(2) << 32
	if gs.Version != expectedV2 {
		t.Errorf("Expected version %d, got %d", expectedV2, gs.Version)
	}
}
