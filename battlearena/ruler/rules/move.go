package rules

import (
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)


type localMoveCtx struct {
	*GameState
	log *logrus.Entry
}

// @spec-link [[mech_action_economy_action_cost_rules]]
func (gs *GameState) Move(msg *message.Message, req rulermethods.ControllerMove) (reply *message.Message) {
	ctx := localMoveCtx{
		GameState: gs,
		log: gs.Logger.WithFields(logrus.Fields{
			"RequestID":    msg.RequestId.String()[0:8],
			"controllerID": req.ControllerID.String()[0:8],
			"entityID":     req.EntityID.String()[0:8],
			"path":         req.Path,
			"rule":         "move",
		}),
	}
	ctx.log.Debug("Controller move request")

	ok, reply := ctx.preMoveChecks(msg, req)
	if !ok {
		return reply
	}

	ent := ctx.Entities[req.EntityID]
	fromPos := ent.Position
	destPos := req.Path[len(req.Path)-1]

	// @spec-link [[mech_trigger_system]]
	// Fire OnExit for the cell the entity is leaving.
	ctx.ProcessPositionalEffects(ent, fromPos, property.TriggerOnExit)

	// Move the entity
	ctx.Grid.MoveEntity(fromPos, destPos, ent.ID)

	ctx.log.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8],
		"from":     fromPos,
		"to":       destPos}).Debug("Entity moved")

	// Update entity position in GameState so it's correct for effect processing and reply.
	ent.Position = destPos
	ctx.Entities[req.EntityID] = ent

	// @spec-link [[mech_trigger_system]]
	// Fire OnEnter for the cell the entity just entered.
	ctx.ProcessPositionalEffects(ent, destPos, property.TriggerOnEnter)

	// Reload entity in case OnEnter positional effects modified it (e.g. applied poison).
	ent = ctx.Entities[req.EntityID]

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + len(req.Path)*20

	// Pay the move cost ! (1/tile)
	prop := ent.GetPropertyC(property.Movement)
	prop.SetValue(prop.GetValue() - len(req.Path))

	ent.UpdateProperty(prop)

	// Update the entity
	ctx.Entities[req.EntityID] = ent

	ctx.IncVersion()

	// notify all controllers of the movement.
	for _, ctrl := range ctx.Controllers {
		ctrl.NotifyActor(message.Create(nil, rulermethods.ControllerMoved{
			EntityID:     req.EntityID,
			Path:         req.Path,
			ControllerID: req.ControllerID,
			Version:      gs.Version,
		}, nil))
	}

	// reply with the new entities state (opaque to the client)
	reply = msg.Reply()

	reply.Content = rulermethods.ControllerMoveReply{
		Entity: gs.Entities[req.EntityID],
	}

	return reply
}

func (ctx *localMoveCtx) preMoveChecks(msg *message.Message, req rulermethods.ControllerMove) (ok bool, reply *message.Message) {
	// Check if the entity exists.

	ent, found := ctx.Entities[req.EntityID]
	if !found {
		ctx.log.Error("Entity not found")
		return false, msg.ReplyWithError("Entity not found", "entity.notfound")
	}

	// Check if the controller is allowed to move the entity
	if !ctx.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		ctx.log.Error("Controller is not allowed to move this entity")
		return false, msg.ReplyWithError("Controller is not allowed to move this entity", "entity.controller.missmatch")
	}
	if ctx.Turner.CurrentEntityTurn != req.EntityID {
		ctx.log.Error("It is not this entity turn")
		return false, msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch")
	}

	// fetch movement distance
	mvt := ent.GetProperty(property.Movement)
	movementDistance := mvt.(*defaultproperty.DefaultIntCounterProperty).MaxValue

	// check path length
	if len(req.Path) > movementDistance {
		ctx.log.WithFields(logrus.Fields{
			"pathLength": len(req.Path),
			"max":        movementDistance,
		}).Error("Path is too long")
		return false, msg.ReplyWithError("Path is too long", "entity.path.too.long")
	}

	// fetch jumpheight
	jumpHeight := ent.GetPropertyI(property.JumpHeight).I()

	// Check if the path is valid
	cells := ctx.Grid.CellsForPositions(req.Path)
	// a valid path is a path that contains only walkable cells and all cells must be adjascent
	for i, c := range cells {
		if c.Type == cell.Ground || c.Type == cell.Dirt {
			// Multi-entity cells: a cell is blocked if it contains a non-WalkThrough entity (other than self).
			// @spec-link [[mechanic_multi_entity_cell_system]]
			if ctx.HasBlockingEntity(c.EntityIDs, req.EntityID) {
				ctx.log.WithFields(logrus.Fields{
					"position": c.Position,
				}).Error("Path contains a blocking entity")
				return false, msg.ReplyWithError("Path contains an occupied cell", "entity.path.occupied")
			}
			if i > 0 && !cells[i-1].Position.IsAdjacent(c.Position, jumpHeight) {
				ctx.log.WithFields(logrus.Fields{
					"jumpHeight": jumpHeight,
				}).Error("Path is not valid")
				return false, msg.ReplyWithError("Invalid path", "entity.path.notvalid")
			}
		} else {
			ctx.log.WithFields(logrus.Fields{
				"cellType": c.Type,
				"position": c.Position,
			}).Error("Path is not valid")
			return false, msg.ReplyWithError("Invalid path(wrong type)", "entity.path.obstacle")
		}

	}


	// ensure entity has movement credits to perform the action.
	prop := ent.GetPropertyC(property.Movement)
	if prop.GetValue() <= 0 {
		ctx.log.Error("Entity has no movement credits")
		return false, msg.ReplyWithError("Entity has no movement credits", "entity.movement.nocredits")
	}

	if prop.GetValue() < len(req.Path) {
		ctx.log.Error("Entity has not enough movement credits")
		return false, msg.ReplyWithError("Entity has not enough movement credits", "entity.movement.credits")
	}

	// ensure entity is adjacent to the first move.
	if !ent.Position.IsAdjacent(req.Path[0], jumpHeight) {
		ctx.log.Error("Entity is not adjacent to the first move")
		return false, msg.ReplyWithError("Entity is not adjacent to the first move", "entity.path.notadjacent")
	}

	// can't move if has already moved this turn.
	propMoved := ent.GetProperty(property.HasMoved)
	if propMoved.Get().(bool) {
		ctx.log.Error("Entity has already moved this turn")
		return false, msg.ReplyWithError("Entity has already moved this turn", "entity.movement.already")
	}

	return true, reply
}
