package rules

import (
	"testing"

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/google/uuid"
)

// @test-link [[rule_forfeit_battle]]
// @test-link [[rule_team_mechanics]]

func setTeam(gs *GameState, entityID uuid.UUID, teamID int) {
	ent := gs.Entities[entityID]
	ent.RepsertPropertyValue(property.TeamID, teamID)
	gs.Entities[entityID] = ent
}

func TestRuleForfeitPvP(t *testing.T) {
	gs, fake := makeGameStateForTwo()

	// Assign teams
	// Team 1: Controller 1
	// Team 2: Controller 2
	setTeam(gs, fake.Entity1, 1)
	setTeam(gs, fake.Entity2, 1)
	setTeam(gs, fake.Entity3, 2)
	setTeam(gs, fake.Entity4, 2)

	_, winnerTeamID, finished := gs.Forfeit(fake.Controller1)

	if !finished {
		t.Errorf("Expected battle to be finished")
	}

	if winnerTeamID != 2 {
		t.Errorf("Expected winning team to be 2, got %d", winnerTeamID)
	}

	// Team 1 entities should be removed
	if _, ok := gs.Entities[fake.Entity1]; ok {
		t.Errorf("Expected Entity 1 to be removed")
	}
	if _, ok := gs.Entities[fake.Entity2]; ok {
		t.Errorf("Expected Entity 2 to be removed")
	}

	// Team 2 entities should remain
	if _, ok := gs.Entities[fake.Entity3]; !ok {
		t.Errorf("Expected Entity 3 to remain")
	}
	if _, ok := gs.Entities[fake.Entity4]; !ok {
		t.Errorf("Expected Entity 4 to remain")
	}
}


func TestRuleForfeitNoEntities(t *testing.T) {
	gs, _ := makeGameStateForTwo()
	
	// Try forfeiting for a controller with no entities
	_, _, finished := gs.Forfeit(uuid.New())
	
	if finished {
		t.Errorf("Expected battle to NOT be finished")
	}
}
