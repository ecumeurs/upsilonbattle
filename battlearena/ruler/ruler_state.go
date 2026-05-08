package ruler

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

// @spec-link [[api_go_battle_start]]
// @spec-link [[mechanic_mech_arena_lifecycle]]

// getState returns the high-level status of the battle arena.
func (r *Ruler) getState(ctx actor.CallContext) {
	r.RequestLogger.Debug("GetState")
	reply := ctx.Msg.Reply()
	reply.Content = rulermethods.GetStateReply{
		GameState:               r.CurrentState.String(),
		NbControllers:           len(r.GameState.Controllers),
		NbControllersExpected:   r.NbControllers,
		NbEntitiesPerController: r.NbEntitiesPerController,
		CurrentEntityTurn:       r.GameState.Turner.CurrentEntityTurn,
	}

	ctx.Reply(reply)
}

// getBoardState returns the complete spatial and entity state of the arena.
func (r *Ruler) getBoardState(ctx actor.CallContext) {
	r.RequestLogger.Debug("GetBoardState")
	req := ctx.Msg.TargetMethod.(rulermethods.GetBoardState)

	entities := make([]entity.Entity, 0, len(r.GameState.Entities))
	for _, e := range r.GameState.Entities {
		entities = append(entities, e)
	}

	reply := ctx.Msg.Reply()
	reply.Content = rulermethods.GetBoardStateReply{
		Grid:          r.GameState.Grid,
		Entities:      entities,
		TurnState:     r.GameState.Turner.GetTurnState(),
		WinnerTeamID:  r.GameState.WinnerTeamID,
		Version:       r.GameState.Version,
		ActionContext: req.ActionContext,
	}

	ctx.Reply(reply)
}

// getGridState returns only the map data.
func (r *Ruler) getGridState(ctx actor.CallContext) {
	r.RequestLogger.Debug("GetGridState")
	reply := ctx.Msg.Reply()
	reply.Content = rulermethods.GetGridStateReply{
		Grid: r.GameState.Grid,
	}

	ctx.Reply(reply)
}

// getEntitiesState returns the state of all entities in the arena.
func (r *Ruler) getEntitiesState(ctx actor.CallContext) {
	r.RequestLogger.Debug("GetEntitiesState")

	reply := ctx.Msg.Reply()
	ent := make([]entity.Entity, 0)
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	reply.Content = rulermethods.GetEntitiesStateReply{
		Entities: ent,
		Turn:     r.GameState.Turner.GetTurnState(),
	}

	ctx.Reply(reply)
}

func (r *Ruler) testingDeleteEntity(ctx actor.NotificationContext) {
	msg := ctx.Msg.TargetMethod.(rulermethods.TestingDeleteEntity)
	r.logger.WithField("entityID", msg.EntityID.String()[0:8]).Info("Testing-only removal of entity")
	delete(r.GameState.Entities, msg.EntityID)
	r.GameState.Turner.RemoveEntity(msg.EntityID)
}

func (r *Ruler) testingGetState(ctx actor.CallContext) {
	ctx.Reply(message.Create(nil, rulermethods.TestingGetStateReply{
		CurrentEntityTurn: r.GameState.Turner.CurrentEntityTurn,
		CurrentState:      r.CurrentState.String(),
		WinnerTeamID:      r.GameState.WinnerTeamID,
	}, nil))
}
