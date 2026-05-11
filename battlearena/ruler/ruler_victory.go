package ruler

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/ecumeurs/upsilontypes/property"
)

// evaluateVictory checks if a win condition is met and transitions the state if necessary.
func (r *Ruler) evaluateVictory(nextTurnEnt uuid.UUID) {
	remainingTeams := make(map[int]bool)
	winningTeamID := 0
	for _, ent := range r.GameState.Entities {
		remainingTeams[ent.GetPropertyI(property.TeamID).I()] = true
		winningTeamID = ent.GetPropertyI(property.TeamID).I()
	}

	if len(remainingTeams) <= 1 || nextTurnEnt == uuid.Nil {
		// @spec-link [[rule_team_mechanics]]
		r.RequestLogger.Info("##### END OF BATTLE! #####")
		r.CurrentState = Finished
		r.GameState.WinnerTeamID = winningTeamID
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerTeamID: winningTeamID,
				Version:      r.GameState.Version,
			}, nil))
		}
	}
}
