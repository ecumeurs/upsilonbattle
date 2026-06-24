package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/defaultproperty"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Channeling: a skill carrying a Channeling delay > 0 commits its caster to a
// deferred, high-impact action. The caster pays all costs upfront, is locked out
// (cannot move/act/retarget), and sits dormant in the Turner until its delayed
// turn arrives — at which point the effect resolves. Taking enough damage while
// channeling interrupts it (costs are sunk, never refunded).
//
// This is a CASTER-RESCHEDULE mechanic: the caster is already a Turner cursor, so
// we just push its next tick out by the channel delay. No separate entity.
//
// @spec-link [[mechanic_channeling_mechanic]]

// isChanneledSkill reports whether a skill is channeled (Channeling delay > 0).
func isChanneledSkill(sk skill.Skill) bool {
	return sk.GetPropertyC(property.Channeling).GetMaxValue() > 0
}

// channelDelay returns the channeling delay (in delay units) of a channeled skill.
func channelDelay(sk skill.Skill) int {
	return sk.GetPropertyC(property.Channeling).GetMaxValue()
}

// isEntityTargetMode reports whether a skill aims at an entity (which the channel
// must follow if it moves) rather than a fixed tile.
func isEntityTargetMode(sk skill.Skill) bool {
	tt := def.TargetTypes(sk.GetProperty(property.TargetType).(*defaultproperty.DefaultStringProperty).Get().(string))
	return tt != def.TargetTypeTile
}

// beginChannel starts a channel: it pays every pre-execution cost upfront (sunk,
// never refunded), captures the target per the skill's targeting mode, marks the
// caster IsCasting, and locks it out. The effect itself is NOT applied here — it
// resolves later via ResolveChannel when the dormant caster is re-picked.
//
// @spec-link [[mechanic_channeling_mechanic]]
func (ctx *localSkillCtx) beginChannel(msg *message.Message, req rulermethods.ControllerUseSkill, ent entity.Entity, sk skill.Skill) (reply *message.Message, damaged []rulermethods.ControllerAttacked, affected []rulermethods.ControllerSkillUsed) {
	damaged = make([]rulermethods.ControllerAttacked, 0)
	affected = make([]rulermethods.ControllerSkillUsed, 0)

	// Pay ALL pre-execution costs upfront. These are sunk: not refunded on
	// interruption, fizzle, or death.
	ent, sk = ctx.paySkillCost(ent, sk)
	ent.Skills[sk.ID] = sk

	// Capture the target the way the skill's targeting mode means it: entity-target
	// channels follow the entity (it can move before resolution); tile/self channels
	// resolve at a fixed coordinate.
	cs := &entity.CastingState{SkillID: sk.ID}
	if isEntityTargetMode(sk) && len(ctx.targetedEntities) > 0 {
		cs.TargetEntity = ctx.targetedEntities[0].ID
	} else {
		pos := req.Target
		cs.TargetPos = &pos
	}
	ent.IsCasting = cs

	// Lock the caster out for the channel's duration.
	hasActed := ent.GetProperty(property.HasActed)
	hasActed.Set(true)
	ent.UpdateProperty(hasActed)
	hasMoved := ent.GetProperty(property.HasMoved)
	hasMoved.Set(true)
	ent.UpdateProperty(hasMoved)

	ctx.Entities[req.EntityID] = ent
	ctx.IncVersion()

	ctx.log.WithField("channelDelay", channelDelay(sk)).Info("Channel started")
	reply = msg.Reply()
	reply.Content = rulermethods.ControllerUseSkillReply{Attacker: ent, Results: []rulermethods.ActionResult{}}
	return
}

// ResolveChannel resolves a caster's in-flight channel: it re-derives the live
// target (entity targets may have moved), re-validates targeting (the target may
// have moved out of range → the channel fizzles), then applies the stored skill's
// effect. The caster is released (IsCasting cleared) either way. The returned
// damaged/affected mirror UseSkill's so the actor can fan out notifications.
//
// fizzled is true when the channel produced no effect (target gone or out of range);
// costs remain sunk and the caster still recovers normally.
//
// @spec-link [[mechanic_channeling_mechanic]]
func ResolveChannel(gs *gamestate.GameState, casterID uuid.UUID) (damaged []rulermethods.ControllerAttacked, affected []rulermethods.ControllerSkillUsed, fizzled bool) {
	ent, ok := gs.Entities[casterID]
	if !ok || ent.IsCasting == nil {
		return nil, nil, false
	}

	ctx := localSkillCtx{
		GameState: gs,
		log: gs.Logger.WithFields(logrus.Fields{
			"entityID": casterID.String()[0:8],
			"skillID":  ent.IsCasting.SkillID.String()[0:8],
			"rule":     "channel.resolve",
		}),
	}

	cs := ent.IsCasting
	sk := ent.Skills[cs.SkillID]

	// Re-derive the live target position. Entity targets follow the entity.
	var targetPos position.Position
	if cs.TargetEntity != uuid.Nil {
		tgt, alive := gs.Entities[cs.TargetEntity]
		if !alive {
			ctx.log.Info("Channel fizzled: target entity gone")
			ent.IsCasting = nil
			gs.Entities[casterID] = ent
			return nil, nil, true
		}
		targetPos = tgt.Position
	} else {
		targetPos = *cs.TargetPos
	}

	// Re-validate targeting against current positions — the target could have moved
	// out of range during the channel.
	msg := message.Create(nil, nil, nil)
	if valid, _ := ctx.checkSkillTarget(msg, ent, targetPos, sk); !valid {
		ctx.log.Info("Channel fizzled: target invalid at resolution")
		ent.IsCasting = nil
		gs.Entities[casterID] = ent
		return nil, nil, true
	}

	// Release the caster before applying so any self-targeted effect sees a
	// non-casting entity.
	ent.IsCasting = nil
	gs.Entities[casterID] = ent

	_, damaged, affected, errStr, errKey := ctx.applyDirectSkillEffect(&ent, sk, targetPos)
	if errStr != "" {
		ctx.log.WithField("errKey", errKey).Error("Channel resolution effect error: " + errStr)
		return nil, nil, true
	}
	gs.Entities[casterID] = ent
	return damaged, affected, false
}

// ApplyInterruption fills a channeling entity's interruption gauge when it takes
// damage (10 points per 1 damage). At >= 100 the channel fails: IsCasting is
// cleared (sunk costs are NOT refunded) and the caster is pulled out of its distant
// dormant slot and re-queued into a near recovery (the skill's Delay cost), NOT the
// channel delay. No-op for non-casting entities. Returns true if it interrupted.
//
// Called at every damage site (skill effects, attacks, positional effects) so any
// source of damage can break a channel.
//
// @spec-link [[mechanic_channeling_mechanic]]
func ApplyInterruption(gs *gamestate.GameState, target *entity.Entity, damage int) bool {
	if target.IsCasting == nil || damage <= 0 {
		return false
	}

	target.IsCasting.Interruption += damage * 10
	if target.IsCasting.Interruption < 100 {
		return false
	}

	// Channel broken. Costs were paid upfront and are not refunded.
	sk := target.Skills[target.IsCasting.SkillID]
	recovery := sk.GetPropertyC(property.Delay).GetMaxValue() // absence defaults to 500

	target.IsCasting = nil
	gs.Turner.RemoveEntity(target.ID)
	target.CurrentDelay = recovery
	gs.Turner.AddEntity(target.ID, recovery)

	gs.Logger.WithFields(logrus.Fields{
		"entityID": target.ID.String()[0:8],
		"recovery": recovery,
	}).Info("Channel interrupted")
	return true
}
