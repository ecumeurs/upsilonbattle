package rules

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

type localSkillCtx struct {
	*gamestate.GameState
	log              *logrus.Entry
	targetedTiles    []position.Position
	targetedEntities []entity.Entity
}

// preSkillChecks performs a sequence of validations before a skill can be used.
// It verifies entity existence, controller authorization, turn order, skill existence, and action state.
// Intent: Prevent illegal skill usage and ensure game state consistency.
// @spec-link [[mech_skill_validation]]
func (ctx *localSkillCtx) preSkillChecks(msg *message.Message, req rulermethods.ControllerUseSkill) (ok bool, reply *message.Message) {

	ent, found := ctx.Entities[req.EntityID]
	if !found {
		ctx.log.Error("Entity not found")
		return false, msg.ReplyWithError("Entity not found", "entity.notfound")
	}

	// Check if the controller is allowed to use the entity
	if !ctx.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		ctx.log.Error("Controller is not allowed to use this entity")
		return false, msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.mismatch")
	}

	if ctx.Turner.CurrentEntityTurn != req.EntityID {
		ctx.log.Error("It is not this entity turn")
		return false, msg.ReplyWithError("It is not this entity turn", "entity.turn.mismatch")
	}

	// seek skill in entity
	skill, found := ent.Skills[req.SkillID]
	if !found {
		ctx.log.Error("Skill not found")
		return false, msg.ReplyWithError("Skill not found", "skill.notfound")
	}

	// entity has already acted this turn.
	if ent.HasActed() {
		ctx.log.Error("Entity has already acted this turn")
		return false, msg.ReplyWithError("Entity has already acted this turn", "entity.alreadyacted")
	}

	// no target for passives!
	if skill.IsPassive() || skill.IsReaction() || skill.IsCounter() {
		return true, reply
	}

	// Skill target is valid.
	ok, reply = ctx.checkSkillTarget(msg, ent, req.Target, skill)
	if !ok {
		return false, reply
	}

	// checks if the user can pay the skill cost
	ok, reply = ctx.checkSkillCost(msg, ent, skill)
	if !ok {
		return false, reply
	}

	return true, reply
}

// checkSkillTarget validates the selected target position and entities against the skill's targeting rules.
// It checks range, grid boundaries, and target type compatibility (Self, Tile, Entity, Friend, Enemy).
// Intent: Enforce tactical constraints and targeting mechanics.
// @spec-link [[mech_skill_validation_entity_targeting_rules_verification]]
func (ctx *localSkillCtx) checkSkillTarget(msg *message.Message, user entity.Entity, target position.Position, sk skill.Skill) (bool, *message.Message) {
	mech := sk.GetProperty(property.TargetingMechanics).(*defaultproperty.DefaultStringProperty)
	targetype := sk.GetProperty(property.TargetType).(*defaultproperty.DefaultStringProperty)
	zone := sk.GetProperty(property.Zone).(*def.ZoneProperty)
	rng := sk.GetProperty(property.Range).(property.IntCounterProperty)

	// Check target within range
	// Use 2D Manhattan distance (ignoring height) 
	dist := tools.Abs(user.Position.X-target.X) + tools.Abs(user.Position.Y-target.Y)
	if dist > rng.GetMaxValue() || dist < rng.GetValue() {
		ctx.log.WithFields(logrus.Fields{
			"dist":     dist,
			"minRange": rng.GetValue(),
			"maxRange": rng.GetMaxValue(),
			"userPos":  user.Position,
			"target":   target,
		}).Error("Target is not in range")
		return false, msg.ReplyWithError("Target is not in range", "skill.target.range")
	}

	selectedZone := ctx.Grid.SelectPositionsByPattern(target, zone.ZonePattern)

	if len(selectedZone) == 0 {
		ctx.log.Error("Target is not in grid")
		return false, msg.ReplyWithError("Target is not in grid", "skill.target.outofgrid")
	}

	// considere only TargetingMechanicsAnywhere for now ... because it's the only one implemented :p
	if mech.Get().(string) == string(def.TargetingMechanicsAnywhere) || mech.Get().(string) == string(def.TargetingMechanicsLOS) {
		// only okay.
	}

	switch def.TargetTypes(targetype.Get().(string)) {
	case def.TargetTypeSelf:
		if !user.Position.Equals(target) {
			ctx.log.Error("Target is not self")
			return false, msg.ReplyWithError("Target is not self", "skill.target.self")
		}
		ctx.targetedEntities = append(ctx.targetedEntities, user)
	case def.TargetTypeTile:
		// TargetTypeTile targets empty cells (no entity present)
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos)
			if !c.IsOccupied() {
				ctx.targetedTiles = append(ctx.targetedTiles, pos)
			}
		}
	case def.TargetTypeEntity:
		// @spec-link [[mechanic_multi_entity_cell_system]]
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos)
			for _, entID := range c.EntityIDs {
				if ent, ok := ctx.Entities[entID]; ok {
					ctx.targetedEntities = append(ctx.targetedEntities, ent)
				}
			}
		}
	case def.TargetTypeFriendOnly:
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos)
			for _, entID := range c.EntityIDs {
				if targetEnt, ok := ctx.Entities[entID]; ok {
					if targetEnt.GetPropertyI(property.TeamID).I() == user.GetPropertyI(property.TeamID).I() {
						ctx.targetedEntities = append(ctx.targetedEntities, targetEnt)
					}
				}
			}
		}
	case def.TargetTypeEnemyOnly:
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos)
			for _, entID := range c.EntityIDs {
				if targetEnt, ok := ctx.Entities[entID]; ok {
					if targetEnt.GetPropertyI(property.TeamID).I() != user.GetPropertyI(property.TeamID).I() {
						ctx.targetedEntities = append(ctx.targetedEntities, targetEnt)
					}
				}
			}
		}
	case def.TargetTypeEntityOrTile:
		ctx.targetedTiles = make([]position.Position, 0)
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos)
			if c.IsOccupied() {
				for _, entID := range c.EntityIDs {
					if ent, ok := ctx.Entities[entID]; ok {
						ctx.targetedEntities = append(ctx.targetedEntities, ent)
					}
				}
			} else {
				ctx.targetedTiles = append(ctx.targetedTiles, pos)
			}
		}
	}


	if len(ctx.targetedTiles) == 0 && len(ctx.targetedEntities) == 0 {
		ctx.log.Error("No Target found")
		return false, msg.ReplyWithError("No Target found", "skill.target.none")
	}

	return true, msg
}

// checkSkillCost verifies if the entity has sufficient resources (HP, MP, SP, Movement) to pay the skill's cost.
// It also checks if the skill is currently on cooldown.
// Intent: Maintain economic balance and prevent spamming of powerful abilities.
// @spec-link [[mech_skill_validation_economic_cost_verification_stat_leech]]
// @spec-link [[mech_skill_validation_economic_cost_verification_cooldown_check]]
func (ctx *localSkillCtx) checkSkillCost(msg *message.Message, user entity.Entity, sk skill.Skill) (bool, *message.Message) {
	// for all costs in skill check if the user has enough.
	if sk.Cooldown > 0 {
		ctx.log.Error("Skill is on cooldown")
		return false, msg.ReplyWithError("Skill is on cooldown", "skill.cooldown")
	}

	// delay and cooldown don't remove anything.

	skp := sk.GetPropertyI(property.HPLeech)
	hp := user.GetPropertyC(property.HP)
	if hp.GetValue() < skp.I() {
		ctx.log.Error("Not enough HP")
		return false, msg.ReplyWithError("Not enough HP", "skill.cost.hp")
	}

	skp = sk.GetPropertyI(property.MPLeech)
	hp = user.GetPropertyC(property.MP)
	if hp.GetValue() < skp.I() {
		ctx.log.Error("Not enough MP")
		return false, msg.ReplyWithError("Not enough MP", "skill.cost.mp")
	}

	skp = sk.GetPropertyI(property.SPLeech)
	hp = user.GetPropertyC(property.SP)
	if hp.GetValue() < skp.I() {
		ctx.log.Error("Not enough SP")
		return false, msg.ReplyWithError("Not enough SP", "skill.cost.sp")
	}

	skp = sk.GetPropertyI(property.MvtCost)
	hp = user.GetPropertyC(property.Movement)
	if hp.GetValue() < skp.I() {
		ctx.log.Error("Not enough Mvt")
		return false, msg.ReplyWithError("Not enough Mvt", "skill.cost.mvt")
	}

	return true, msg
}

// paySkillCost deducts the required resources from the entity and sets the skill on cooldown.
// Intent: Finalize the transaction of using a skill.
func (ctx *localSkillCtx) paySkillCost(user entity.Entity, sk skill.Skill) (entity.Entity, skill.Skill) {
	skpc := sk.GetPropertyC(property.Cooldown)
	sk.Cooldown = skpc.GetMaxValue()

	skp := sk.GetPropertyI(property.HPLeech)
	hp := user.GetPropertyC(property.HP)
	hp.SetValue(hp.GetValue() - skp.I())
	user.UpdateProperty(hp)

	skp = sk.GetPropertyI(property.MPLeech)
	mp := user.GetPropertyC(property.MP)
	mp.SetValue(mp.GetValue() - skp.I())
	user.UpdateProperty(mp)

	skp = sk.GetPropertyI(property.SPLeech)
	sp := user.GetPropertyC(property.SP)
	sp.SetValue(sp.GetValue() - skp.I())
	user.UpdateProperty(sp)

	skp = sk.GetPropertyI(property.MvtCost)
	mvt := user.GetPropertyC(property.Movement)
	mvt.SetValue(mvt.GetValue() - skp.I())
	user.UpdateProperty(mvt)

	return user, sk
}
