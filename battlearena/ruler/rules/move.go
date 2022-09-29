package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/cell"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

type localMoveCtx struct {
	*GameState
	log *logrus.Entry
}

func (gs *GameState) Move(msg message.Message, req rulermethods.ControllerMove) (reply message.Message) {
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

	// Move the entity
	ctx.Grid.MoveEntity(ent.Position, req.Path[len(req.Path)-1], ent.ID)

	ctx.log.WithFields(logrus.Fields{
		"entityID": req.EntityID.String()[0:8],
		"from":     ent.Position,
		"to":       req.Path[len(req.Path)-1]}).Debug("Entity moved")

	// update entity position
	ent.Position = req.Path[len(req.Path)-1]

	// Compute the new delay
	ent.CurrentDelay = ent.CurrentDelay + len(req.Path)*200

	// Update the entity
	ctx.Entities[req.EntityID] = ent

	// reply with the new entities state (opaque to the client)
	reply = msg.Reply()

	reply.Content = rulermethods.ControllerMoveReply{
		Entity: ent,
	}

	return reply
}

func (ctx *localMoveCtx) preMoveChecks(msg message.Message, req rulermethods.ControllerMove) (ok bool, reply message.Message) {
	// Check if the entity exists.

	// Check if the controller is allowed to move the entity
	if !ctx.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		ctx.log.Error("Controller is not allowed to move this entity")
		return false, msg.ReplyWithError("Controller is not allowed to move this entity", "entity.controller.missmatch")
	}
	if ctx.Turner.CurrentEntityTurn != req.EntityID {
		ctx.log.Error("It is not this entity turn")
		return false, msg.ReplyWithError("It is not this entity turn", "entity.turn.missmatch")
	}

	ent, found := ctx.Entities[req.EntityID]
	if !found {
		ctx.log.Error("Entity not found")
		return false, msg.ReplyWithError("Entity not found", "entity.notfound")
	}

	// fetch movement distance
	mvt := ent.GetProperty("Movement", entity.DefaultMovement())
	movementDistance := mvt.(*property.DefaultIntCounterProperty).Value

	// check path length
	if len(req.Path) > movementDistance {
		ctx.log.WithFields(logrus.Fields{
			"pathLength": len(req.Path),
			"max":        movementDistance,
		}).Error("Path is too long")
		return false, msg.ReplyWithError("Path is too long", "entity.path.too.long")
	}

	// fetch jumpheight
	jumpHeight := ent.GetPropertyI("JumpHeight", entity.DefaultJumpHeight()).I()

	// Check if the path is valid
	cells := ctx.Grid.CellsForPositions(req.Path)
	// a valid path is a path that contains only walkable cells and all cells must be adjascent
	for i, c := range cells {
		if c.Type == cell.Ground {
			if i > 0 && !cells[i-1].Position.IsAdjacent(c.Position, jumpHeight) {
				ctx.log.WithFields(logrus.Fields{
					"jumpHeight": jumpHeight,
				}).Error("Path is not valid")
				return false, msg.ReplyWithError("Invalid path", "entity.path.invalid")
			}
		} else {
			ctx.log.WithFields(logrus.Fields{
				"cellType": c.Type,
				"position": c.Position,
			}).Error("Path is not valid")
			return false, msg.ReplyWithError("Invalid path(wrong type)", "entity.path.invalid")
		}
	}

	return true, reply
}
