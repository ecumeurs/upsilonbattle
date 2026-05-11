package ruler

import (
	"fmt"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller/controllermethods"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rules"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

// addController handles the addition of a controller to the battle.
func (r *Ruler) addController(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.AddController)
	r.RequestLogger.WithFields(logrus.Fields{
		"ControllerID": req.ControllerID.String()[0:8]}).Info("AddController")

	if _, ok := r.GameState.Controllers[req.ControllerID]; ok {
		r.RequestLogger.Warn("Controller already registered")
		ctx.Reply(ctx.Msg.ReplyWithError(fmt.Sprintf("Controller %s already registered", req.ControllerID), "controller.already.registered"))
		return
	}

	if r.GameState.Grid == nil {
		r.RequestLogger.Error("Cannot add controller: Grid is not initialized")
		ctx.Reply(ctx.Msg.ReplyWithError("Grid not initialized", "arena.not_ready.no_grid"))
		return
	}

	// @spec-link [[mech_controller_handshake]]
	// @spec-link [[mech_controller_communication_sequence]]
	req.Controller.NotifyActor(message.Create(nil, controllermethods.SetQueue{
		ControllerID: req.ControllerID,
		Ruler:        r,
	}, nil))

	r.GameState.Controllers[req.ControllerID] = req.Controller
	r.ControllerBattleReady[req.ControllerID] = false

	reply := ctx.Msg.Reply()
	ent := make([]entity.Entity, 0)
	for _, e := range r.GameState.Entities {
		ent = append(ent, e)
	}

	reply.Content = rulermethods.AddControllerReply{
		ControllerID: req.ControllerID,
		Grid:         r.GameState.Grid,
		TurnState:    r.GameState.Turner.GetTurnState(),
		Entities:     ent,
	}

	ctx.Reply(reply)

	if len(r.GameState.Controllers) == r.NbControllers {
		if r.CurrentState == InProgress {
			r.RequestLogger.Info("All controllers re-registered after resurrection, skipping BattleStart")
		} else {
			r.RequestLogger.Info("All controllers registered, scheduling BattleStart")
			r.SelfNotifyDelayed(rulermethods.BattleStart{}, 20*time.Millisecond)
		}
	}
}

// controllerMove handles a movement request from a controller.
// @spec-link [[mech_action_economy_action_cost_rules]]
func (r *Ruler) controllerMove(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerMove)
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	reply := rules.Move(r.GameState, ctx.Msg, req)
	ctx.Reply(reply)

	if !reply.HasError {
		ent := make([]entity.Entity, 0, len(r.GameState.Entities))
		for _, e := range r.GameState.Entities {
			ent = append(ent, e)
		}
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.GameState.Turner.GetTurnState(),
				Version:  r.GameState.Version,
			}, nil))
		}
	}
}

// controllerAttack handles an attack request from a controller.
// @spec-link [[mech_action_economy_action_cost_rules]]
func (r *Ruler) controllerAttack(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerAttack)
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	reply := rules.Attack(r.GameState, ctx.Msg, req)
	ctx.Reply(reply)

	if !reply.HasError {
		ent := make([]entity.Entity, 0, len(r.GameState.Entities))
		for _, e := range r.GameState.Entities {
			ent = append(ent, e)
		}
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.GameState.Turner.GetTurnState(),
				Version:  r.GameState.Version,
			}, nil))
		}
	}
}

// controllerUseSkill handles a skill usage request from a controller.
func (r *Ruler) controllerUseSkill(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerUseSkill)
	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	reply, damaged, affected := rules.UseSkill(r.GameState, ctx.Msg, req)
	ctx.Reply(reply)

	for _, d := range damaged {
		foectrlid := d.ControllerID
		foectrl, found := r.GameState.Controllers[foectrlid]
		if !found {
			r.RequestLogger.WithFields(logrus.Fields{
				"foeControllerID": foectrlid.String()[0:8]}).Error("Foe controller not found")
		} else {
			foectrl.NotifyActor(message.Create(nil, d, nil))
		}
	}

	for _, d := range affected {
		targetctrlid := d.ControllerID
		targetctrl, found := r.GameState.Controllers[targetctrlid]
		if !found {
			r.RequestLogger.WithFields(logrus.Fields{
				"targetControllerID": targetctrlid.String()[0:8]}).Error("target controller not found")
		} else {
			targetctrl.NotifyActor(message.Create(nil, d, nil))
		}
	}

	if !reply.HasError {
		ent := make([]entity.Entity, 0, len(r.GameState.Entities))
		for _, e := range r.GameState.Entities {
			ent = append(ent, e)
		}
		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.EntitiesStateChanged{
				Entities: ent,
				Turn:     r.GameState.Turner.GetTurnState(),
				Version:  r.GameState.Version,
			}, nil))
		}
	}
}

// controllerForfeit handles the forfeiture of a controller.
// @spec-link [[rule_forfeit_battle]]
func (r *Ruler) controllerForfeit(ctx actor.CallContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerForfeit)
	r.RequestLogger.WithFields(logrus.Fields{
		"controllerID": req.ControllerID.String()[0:8],
		"entityID":     req.EntityID.String()[0:8]}).Info("ControllerForfeit")

	if r.CurrentState != InProgress {
		r.RequestLogger.Error("Game is not in progress")
		ctx.Reply(ctx.Msg.ReplyWithError("Game is not in progress", "game.not.in.progress"))
		return
	}

	_, winnerTeamID, finished := rules.Forfeit(r.GameState, req.ControllerID)

	if finished {
		r.CurrentState = Finished
		r.GameState.WinnerTeamID = winnerTeamID
		r.RequestLogger.Info("##### END OF BATTLE! #####")

		for _, ctrl := range r.GameState.Controllers {
			ctrl.NotifyActor(message.Create(nil, rulermethods.BattleEnd{
				WinnerTeamID: winnerTeamID,
				Version:      r.GameState.Version,
			}, nil))
		}
	}

	ctx.Reply(ctx.Msg.Reply())
}

// controllerQuit handles a controller disconnected or quitting.
func (r *Ruler) controllerQuit(ctx actor.NotificationContext) {
	req := ctx.Msg.TargetMethod.(rulermethods.ControllerQuit)
	r.RequestLogger.Debug("Controller quit notification")

	_, found := r.GameState.Controllers[req.ControllerID]
	if found {
		delete(r.GameState.Controllers, req.ControllerID)
		r.RequestLogger.Info("Controller removed from match")

		for id, ent := range r.GameState.Entities {
			if ent.ControllerID == req.ControllerID {
				r.GameState.Grid.RemoveEntity(ent.Position, id)
				delete(r.GameState.Entities, id)
				r.GameState.Turner.RemoveEntity(id)
			}
		}

		if r.CurrentState != Finished {
			r.evaluateVictory(r.GameState.Turner.CurrentEntityTurn)
		}
	}
}

// notifyController is a notification handler for controller-specific events.
func (r *Ruler) notifyController(ctx actor.NotificationContext) {
}

// controllerPassed is a notification handler for when a controller skips its turn.
func (r *Ruler) controllerPassed(ctx actor.NotificationContext) {
	r.RequestLogger.Debug("Controller passed notification ignored by ruler")
}
