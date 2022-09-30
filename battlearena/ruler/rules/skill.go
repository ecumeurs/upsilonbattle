package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type localSkillCtx struct {
	*GameState
	log              *logrus.Entry
	targetedTiles    []position.Position
	targetedEntities []entity.Entity
}

func (gs *GameState) UseSkill(msg message.Message, req rulermethods.ControllerUseSkill) (reply message.Message, damaged []rulermethods.ControllerAttacked, affected []rulermethods.ControllerSkillUsed) {
	ctx := localSkillCtx{
		GameState: gs,
		log: gs.Logger.WithFields(logrus.Fields{
			"RequestID":    msg.RequestId.String()[0:8],
			"controllerID": req.ControllerID.String()[0:8],
			"entityID":     req.EntityID.String()[0:8],
			"skillID":      req.SkillID.String()[0:8],
			"target":       req.Target,
			"rule":         "skill",
		}),
	}
	ctx.log.Debug("Controller attack request")

	ok, reply := ctx.preSkillChecks(msg, req)
	damaged = make([]rulermethods.ControllerAttacked, 0)
	affected = make([]rulermethods.ControllerSkillUsed, 0)
	if !ok {
		return
	}

	ent := ctx.Entities[req.EntityID]
	sk := ent.Skills[req.SkillID]

	// now we have a target identifed! yata!

	// Hit test!
	if sk.Effect.IsDamaging() {
		if len(ctx.targetedEntities) > 0 {
			newTargets := make([]entity.Entity, 0)
			accuracy := ent.GetPropertyI(property.Accuracy).I()

			for _, target := range ctx.targetedEntities {
				dodge := ent.GetPropertyI(property.Dodge).I()
				if tools.RandomInt(0, 100) < accuracy-dodge {
					newTargets = append(newTargets, target)
				}
			}
			ctx.targetedEntities = newTargets
		}
	}

	// Apply skill effect
	if sk.Effect.IsDamaging() {
		damage := sk.GetPropertyI(property.Damage).I()
		stunpwr := sk.GetPropertyI(property.StunPower).I()
		stunchance := sk.GetPropertyI(property.StunChance).I()
		poisonpwr := sk.GetPropertyI(property.PoisonPower).I()
		poisonchance := sk.GetPropertyI(property.PoisonChance).I()

		updatedEntities := make([]entity.Entity, 0)

		for _, target := range ctx.targetedEntities {
			hp := target.GetPropertyC(property.HP).GetValue()
			defense := target.GetPropertyI(property.Defense).I()
			shield := target.GetPropertyC(property.Shield).GetValue()
			armor := target.GetPropertyI(property.ArmorRating).I()

			truepoison := poisonpwr - defense
			truestun := stunpwr - armor

			truedmg := damage - defense - armor + truepoison + truestun

			// check poisoning & stunning.
			if truepoison > 0 {
				if tools.RandomInt(0, 100) < poisonchance {
					// target is poisoned
					// update poison value.
					poison := target.GetPropertyI(property.Poison).I()
					target.GetPropertyI(property.Poison).Set(poison + truepoison)
				}
			}
			if truestun > 0 {
				if tools.RandomInt(0, 100) < stunchance {
					// target is stunned
					// update stun value.
					stun := target.GetPropertyI(property.Stun).I()
					target.GetPropertyI(property.Stun).Set(stun + truestun)
				}
			}
			if truedmg > 0 {
				// first killoff shield.
				if shield > 0 {
					if shield > truedmg {
						shield -= truedmg
						truedmg = 0
					} else {
						truedmg -= shield
						shield = 0
					}
					//update shield
					target.GetPropertyC(property.Shield).SetValue(shield)
				}
				// then kill off hp.
				if hp > 0 {
					if hp > truedmg {
						hp -= truedmg
						truedmg = 0
					} else {
						truedmg -= hp
						hp = 0
					}
					//update hp
					target.GetPropertyC(property.HP).SetValue(hp)
				}
			}
			updatedEntities = append(updatedEntities, target)
		}

		ctx.targetedEntities = updatedEntities

		for _, tar := range updatedEntities {
			damaged = append(damaged, rulermethods.ControllerAttacked{
				ControllerID:         tar.ControllerID,
				Entity:               tar,
				SkillID:              sk.ID,
				AttackerControllerID: ent.ControllerID,
				Attacker:             ent,
			})
		}
	}

	if sk.Effect.IsHealing() {
		heal := sk.GetPropertyI(property.Heal).I()
		updatedEntities := make([]entity.Entity, 0)
		for _, target := range ctx.targetedEntities {
			hp := target.GetPropertyC(property.HP).GetValue()
			maxhp := target.GetPropertyC(property.HP).GetMaxValue()
			if hp < maxhp {
				if hp+heal > maxhp {
					hp = maxhp
				} else {
					hp += heal
				}
				target.GetPropertyC(property.HP).SetValue(hp)
			}
			updatedEntities = append(updatedEntities, target)
		}
		ctx.targetedEntities = updatedEntities

		for _, tar := range updatedEntities {
			affected = append(affected, rulermethods.ControllerSkillUsed{
				ControllerID:        tar.ControllerID,
				Entity:              tar,
				SkillID:             sk.ID,
				EmitterControllerID: ent.ControllerID,
				Emitter:             ent,
			})
		}
	}

	// for the moment ignore buffs and dots.

	// update entities in global context.
	for _, tar := range ctx.targetedEntities {
		ctx.Entities[tar.ID] = tar
	}

	// now user pay the cost
	ent, sk = ctx.paySkillCost(ent, sk)

	// update the skill in entity
	ent.Skills[sk.ID] = sk

	// Update the entity in the game state
	ctx.Entities[req.EntityID] = ent

	// reply to user
	reply = msg.Reply()
	reply.TargetMethod = rulermethods.ControllerUseSkillReply{
		Entity: ent,
	}

	return
}

func (ctx *localSkillCtx) preSkillChecks(msg message.Message, req rulermethods.ControllerUseSkill) (ok bool, reply message.Message) {

	// Check if the controller is allowed to use the entity
	if !ctx.CheckControllerForEntity(req.ControllerID, req.EntityID) {
		ctx.log.Error("Controller is not allowed to use this entity")
		return false, msg.ReplyWithError("Controller is not allowed to use this entity", "entity.controller.missmatch")
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

	// seek skill in entity
	skill, found := ent.Skills[req.SkillID]
	if !found {
		ctx.log.Error("Skill not found")
		return false, msg.ReplyWithError("Skill not found", "skill.notfound")
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

func (ctx *localSkillCtx) checkSkillTarget(msg message.Message, user entity.Entity, target position.Position, sk skill.Skill) (bool, message.Message) {
	mech := sk.GetProperty(property.TargetingMechanics).(*def.TargetingMechanicsProperty)
	targetype := sk.GetProperty(property.TargetType).(*def.TargetTypeProperty)
	zone := sk.GetProperty(property.Zone).(*def.ZoneProperty)
	rng := sk.GetProperty(property.Range).(*def.RangeProperty)

	selectedZone := ctx.Grid.SelectPositionsByPattern(target, zone.ZonePattern)
	if len(selectedZone) == 0 {
		ctx.log.Error("Target is not in grid")
		return false, msg.ReplyWithError("Target is not in grid", "skill.target.outofgrid")
	}

	// Check target within range
	dist := user.Position.Distance(target)
	if dist > rng.MaxRange || dist < rng.MinRange {
		ctx.log.Error("Target is not in range")
		return false, msg.ReplyWithError("Target is not in range", "skill.target.range")
	}

	// considere only TargetingMechanicsAnywhere for now ... because it's the only one implemented :p
	if mech.TargetingMechanics == def.TargetingMechanicsAnywhere || mech.TargetingMechanics == def.TargetingMechanicsLOS {
		// only okay.
	}

	switch targetype.TargetType {
	case def.TargetTypeSelf:
		if !user.Position.Equals(target) {
			ctx.log.Error("Target is not self")
			return false, msg.ReplyWithError("Target is not self", "skill.target.self")
		}
		ctx.targetedEntities = append(ctx.targetedEntities, user)
	case def.TargetTypeTile:
		ctx.targetedTiles = selectedZone
	case def.TargetTypeEntity:
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos) // should be ok because it has been veted before.
			if c.EntityID != uuid.Nil {
				ctx.targetedEntities = append(ctx.targetedEntities, ctx.Entities[c.EntityID])
			}
		}
	case def.TargetTypeFriendOnly:
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos) // should be ok because it has been veted before.
			if c.EntityID != uuid.Nil {
				if ctx.Entities[c.EntityID].ControllerID == user.ControllerID {
					ctx.targetedEntities = append(ctx.targetedEntities, ctx.Entities[c.EntityID])
				}
			}
		}
	case def.TargetTypeEnemyOnly:
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos) // should be ok because it has been veted before.
			if c.EntityID != uuid.Nil {
				if ctx.Entities[c.EntityID].ControllerID != user.ControllerID {
					ctx.targetedEntities = append(ctx.targetedEntities, ctx.Entities[c.EntityID])
				}
			}
		}
	case def.TargetTypeEntityOrTile:
		ctx.targetedTiles = make([]position.Position, 0)
		ctx.targetedEntities = make([]entity.Entity, 0)
		for _, pos := range selectedZone {
			c, _ := ctx.Grid.CellAt(pos) // should be ok because it has been veted before.
			if c.EntityID != uuid.Nil {
				ctx.targetedEntities = append(ctx.targetedEntities, ctx.Entities[c.EntityID])
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

func (ctx *localSkillCtx) checkSkillCost(msg message.Message, user entity.Entity, sk skill.Skill) (bool, message.Message) {
	// for all costs in skill check if the user has enough.
	if sk.Cooldown > 0 {
		ctx.log.Error("Skill is on cooldown")
		return false, msg.ReplyWithError("Skill is on cooldown", "skill.cooldown")
	}

	// delay and cooldown don't remove anything.

	if skp := sk.GetPropertyI(property.HPLeech); skp != nil {

		hp := user.GetPropertyC(property.HP)
		if hp.GetValue() < skp.I() {
			ctx.log.Error("Not enough HP")
			return false, msg.ReplyWithError("Not enough HP", "skill.cost.hp")
		}
	}

	if skp := sk.GetPropertyI(property.MPLeech); skp != nil {

		hp := user.GetPropertyC(property.MP)
		if hp.GetValue() < skp.I() {
			ctx.log.Error("Not enough HP")
			return false, msg.ReplyWithError("Not enough HP", "skill.cost.hp")
		}
	}

	if skp := sk.GetPropertyI(property.SPLeech); skp != nil {

		hp := user.GetPropertyC(property.SP)
		if hp.GetValue() < skp.I() {
			ctx.log.Error("Not enough HP")
			return false, msg.ReplyWithError("Not enough HP", "skill.cost.hp")
		}
	}

	if skp := sk.GetPropertyI(property.MvtCost); skp != nil {

		hp := user.GetPropertyC(property.Movement)
		if hp.GetValue() < skp.I() {
			ctx.log.Error("Not enough HP")
			return false, msg.ReplyWithError("Not enough HP", "skill.cost.hp")
		}
	}

	return true, msg
}

func (ctx *localSkillCtx) paySkillCost(user entity.Entity, sk skill.Skill) (entity.Entity, skill.Skill) {
	if skp := sk.GetPropertyC(property.Cooldown); skp != nil {
		sk.Cooldown = skp.GetMaxValue()
	}

	// delay and cooldown don't remove anything.

	if skp := sk.GetPropertyI(property.HPLeech); skp != nil {

		hp := user.GetPropertyC(property.HP)
		hp.SetValue(hp.GetValue() - skp.I())
		user.UpdateProperty(hp)

	}

	if skp := sk.GetPropertyI(property.MPLeech); skp != nil {

		mp := user.GetPropertyC(property.MP)
		mp.SetValue(mp.GetValue() - skp.I())
		user.UpdateProperty(mp)
	}

	if skp := sk.GetPropertyI(property.SPLeech); skp != nil {

		sp := user.GetPropertyC(property.SP)
		sp.SetValue(sp.GetValue() - skp.I())
		user.UpdateProperty(sp)
	}

	if skp := sk.GetPropertyI(property.MvtCost); skp != nil {

		mvt := user.GetPropertyC(property.Movement)
		mvt.SetValue(mvt.GetValue() - skp.I())
		user.UpdateProperty(mvt)
	}

	return user, sk
}
