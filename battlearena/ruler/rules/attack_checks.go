package rules

import (

	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

// preAttackChecks performs validations before an attack can be executed.
// It verifies entity existence, turn order, range, line of sight, and team-based restrictions (friendly fire).
// Intent: Enforce combat rules and prevent invalid attack requests.
// @spec-link [[mech_action_economy]]
func (ctx *localAttackCtx) preAttackChecks(msg *message.Message, req rulermethods.ControllerAttack) (ok bool, reply *message.Message) {

	ent, found := ctx.Entities[req.EntityID]
	if !found {
		ctx.log.Error("Entity not found")
		return false, msg.ReplyWithError("Entity not found", "entity.notfound")
	}

	// Check if the controller is allowed to use the entity. Same shared
	// gamestate.CheckControllerForEntity gate used by move/skill/pass; no
	// dedicated attack-validation atom exists yet (see mech_action_economy's
	// gap note), so this reuses the canonical "Controller Mismatch" rule
	// (check #3) already documented for the identical entity.controller.mismatch
	// key on the move path.
	// @spec-link [[mech_move_validation]]
	if !ctx.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		ctx.log.Error("Controller is not allowed to use this entity")
		return false, msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.mismatch")
	}

	if ctx.Turner.CurrentEntityTurn != req.EntityID {
		ctx.log.Error("It is not this entity turn")
		return false, msg.ReplyWithError("It is not this entity turn", "entity.turn.mismatch")
	}

	// Check if the attack is valid.
	// Target must resolve to a real cell inside the grid's [0,Width) x [0,Length) x
	// [0,Height) bounds (Grid.CellAt -> PositionIsInGrid); anything outside is rejected
	// before any occupancy/range/team logic runs.
	// @spec-link [[entity_grid]]
	target, found := ctx.Grid.CellAt(req.Target)
	if !found {
		ctx.log.Error("Target is not found")
		return false, msg.ReplyWithError("Invalid target", "entity.attack.target.invalid")
	}

	if target.Type != cell.Ground && target.Type != cell.Dirt {
		ctx.log.Error("Target is not valid")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.celltype")
	}

	// A target cell must hold an entity (grid-cell EntityIDs collection) and, among
	// its occupants, an attackable Character/Monster -- not just a WalkThrough entity.
	// @spec-link [[mechanic_multi_entity_cell_system]]
	if !target.IsOccupied() {
		ctx.log.Error("Target has no entities")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.noentity")
	}

	// Use FindCharacterInCell to ensure we only target attackable entities.
	targetEntity, found := ctx.FindCharacterInCell(target.EntityIDs)
	if !found {
		ctx.log.Error("Target cell has no attackable character")
		return false, msg.ReplyWithError("No attackable entity", "entity.attack.noentity")
	}

	// range check.
	// @spec-link [[rule_combat_range_validation]]
	attackerRange := ent.GetPropertyI(property.AttackRange)
	weaponRange := ent.GetPropertyI(property.WeaponRange)
	effectiveRange := tools.Max(attackerRange.I(), weaponRange.I())

	// Base distance is 2D Manhattan distance (ISS-098)
	distance2D := tools.Abs(ent.Position.X-target.Position.X) + tools.Abs(ent.Position.Y-target.Position.Y)
	// Vertical distance check
	zDiff := tools.Abs(ent.Position.Z - target.Position.Z)

	if effectiveRange < distance2D || zDiff > (effectiveRange+1) {
		ctx.log.WithFields(logrus.Fields{
			"attackrange":    attackerRange.I(),
			"weaponrange":    weaponRange.I(),
			"effectiveRange": effectiveRange,
			"distance2D":     distance2D,
			"zDiff":       zDiff,
			"attacker":    ent.Position,
			"target":      target.Position,
		}).Error("Target is out of range")
		return false, msg.ReplyWithError("Invalid attack", "entity.attack.outofrange")
	}

	// @spec-link [[rule_friendly_fire_team_validation]]
	attackerTeam := ent.GetPropertyI(property.TeamID)
	targetTeam := targetEntity.GetPropertyI(property.TeamID)

	if attackerTeam.I() == targetTeam.I() {
		ctx.log.WithFields(logrus.Fields{
			"attackerTeam": attackerTeam.I(),
			"targetTeam":   targetTeam.I(),
		}).Error("Friendly fire is not allowed")
		return false, msg.ReplyWithError("Friendly fire is not allowed", "entity.attack.friendlyfire")
	}

	propHasActed := ent.GetProperty(property.HasActed).Get().(bool)
	if propHasActed {
		ctx.log.Error("Entity has already acted")
		return false, msg.ReplyWithError("Entity has already acted", "entity.hasacted")
	}

	return true, reply
}
