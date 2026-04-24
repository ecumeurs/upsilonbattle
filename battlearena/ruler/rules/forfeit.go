package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// @spec-link [[rule_forfeit_battle]]
// @spec-link [[rule_team_mechanics]]

// Forfeit handles the forfeiture of a controller.
// It determines the team of the forfeiting controller and declares all team entities defeated.
// It returns the winner controller ID (if any), the winner team ID (if any) and whether the battle is finished.
func (gs *GameState) Forfeit(controllerID uuid.UUID) (uuid.UUID, int, bool) {
	gs.Logger.WithFields(logrus.Fields{
		"controllerID": controllerID.String()[0:8]}).Info("GameState.Forfeit")

	// 1. Find the team of the forfeiting controller
	var forfeitingTeam int
	found := false
	for _, ent := range gs.Entities {
		if ent.ControllerID == controllerID {
			forfeitingTeam = ent.GetPropertyI(property.TeamID).I()
			found = true
			break
		}
	}

	if !found {
		gs.Logger.Error("Forfeiting controller has no entities")
		return uuid.Nil, 0, false
	}

	// 2. Remove all entities belonging to that team
	for id, ent := range gs.Entities {
		if ent.GetPropertyI(property.TeamID).I() == forfeitingTeam {
			gs.Logger.WithFields(logrus.Fields{
				"entityID": id.String()[0:8],
				"teamID":   forfeitingTeam}).Info("Removing forfeiting team entity")
			gs.Grid.RemoveEntity(ent.Position, id)
			delete(gs.Entities, id)
			gs.Turner.RemoveEntity(id)
		}
	}

	// 3. Check for remaining teams
	remainingTeams := make(map[int]uuid.UUID) // TeamID -> A ControllerID from that team
	for _, ent := range gs.Entities {
		remainingTeams[ent.GetPropertyI(property.TeamID).I()] = ent.ControllerID
	}

	if len(remainingTeams) <= 1 {
		winnerControllerID := uuid.Nil
		winnerTeamID := 0
		for teamID, ctrlID := range remainingTeams {
			winnerControllerID = ctrlID
			winnerTeamID = teamID
			break
		}
		return winnerControllerID, winnerTeamID, true
	}

	gs.IncVersion()

	return uuid.Nil, 0, false
}
