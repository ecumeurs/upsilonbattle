package rules

import (
	"github.com/ecumeurs/upsilontypes/entity/skill/skillweight"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/effect/effectapplicator"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

// @spec-link [[mechanic_mec_skill_payload_resolution]]
// @spec-link [[rule_skill_grading_system]]



// UseSkill handles the execution of an active skill by an entity.
// It performs pre-skill checks, calculates effects via the effect applicator,
// applies damage/status effects, awards credits, and handles skill costs and action state updates.
// Intent: Provide a unified entry point for all skill-based interactions in the battle arena.
// Inputs: msg (*message.Message) - the incoming request message, req (rulermethods.ControllerUseSkill) - the skill usage parameters.
// Outputs: reply (*message.Message) - the outcome message, damaged ([]rulermethods.ControllerAttacked) - list of entities that took damage, affected ([]rulermethods.ControllerSkillUsed) - list of entities affected by status changes.
func (gs *GameState) UseSkill(msg *message.Message, req rulermethods.ControllerUseSkill) (reply *message.Message, damaged []rulermethods.ControllerAttacked, affected []rulermethods.ControllerSkillUsed) {
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

	// now we have a target identifed! yata! For now only work on direct effect ... later :)
	dds, aff, credits, err, errkey := effectapplicator.ApplyDirectEffect(ctx.log, &ent, sk.Effect, req.Target, ctx.targetedTiles, ctx.Grid, ctx.targetedEntities)
	if err != "" {
		ctx.log.Error(err)
		return msg.ReplyWithError(err, errkey), damaged, affected
	}

	// Status Effect Credits (Flat Rate)
	// rule: SkillWeight / 10 credits per application
	if len(aff) > 0 && sk.Effect.IsOverTime() {
		pos, _, _ := skillweight.Calculate(&sk)
		statusCredits := pos / 10
		if statusCredits > 0 {
			credits = append(credits, rulermethods.CreditAward{
				PlayerID: ent.ControllerID,
				Amount:   statusCredits * len(aff), // reward per target affected
				Source:   "status",
			})
		}
	}

	ctx.IncVersion()

	// update entities in global context.
	results := make([]rulermethods.ActionResult, 0)
	for _, res := range dds {
		tar := res.Target
		if tar.GetPropertyC(property.HP).GetValue() <= 0 {
			ctx.log.WithField("targetID", tar.ID.String()[0:8]).Info("Entity killed by skill effect")
			gs.RemoveEntity(tar.ID)
		} else {
			ctx.Entities[tar.ID] = tar
		}
		
		res.CreditAwards = credits
		results = append(results, res)

		damaged = append(damaged, rulermethods.ControllerAttacked{
			ControllerID:         tar.ControllerID,
			Entity:               tar,
			SkillID:              sk.ID,
			AttackerControllerID: ent.ControllerID,
			Attacker:             ent,
			Damage:               res.Damage,
			PrevHP:               res.PrevHP,
			NewHP:                res.NewHP,
			Dead:                 tar.GetPropertyC(property.HP).GetValue() <= 0,
			CreditAwards:         credits,
			Version:              gs.Version,
		})
	}

	for _, res := range aff {
		tar := res.Target
		if tar.GetPropertyC(property.HP).GetValue() <= 0 {
			ctx.log.WithField("targetID", tar.ID.String()[0:8]).Info("Entity killed by skill effect")
			gs.RemoveEntity(tar.ID)
		} else {
			ctx.Entities[tar.ID] = tar
		}

		res.CreditAwards = credits
		results = append(results, res)

		affected = append(affected, rulermethods.ControllerSkillUsed{
			ControllerID:        tar.ControllerID,
			Entity:              tar,
			SkillID:             sk.ID,
			EmitterControllerID: ent.ControllerID,
			Emitter:             ent,
			CreditAwards:        credits,
			Version:             gs.Version,
		})
	}

	// now user pay the cost
	ent, sk = ctx.paySkillCost(ent, sk)

	// update the skill in entity
	ent.Skills[sk.ID] = sk

	// Update the entity in the game state
	// Set that unit acted
	hasActed := ent.GetProperty(property.HasActed)
	hasActed.Set(true)
	ent.UpdateProperty(hasActed)

	// Set that unit moved
	hasMoved := ent.GetProperty(property.HasMoved)
	hasMoved.Set(true)
	ent.UpdateProperty(hasMoved)

	ctx.Entities[req.EntityID] = ent


	// reply to user
	reply = msg.Reply()
	reply.Content = rulermethods.ControllerUseSkillReply{
		Attacker: ent,
		Results:  results,
	}

	return
}


