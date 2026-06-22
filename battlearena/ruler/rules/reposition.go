package rules

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// Movement skills displace a subject (the caster or the target) along the casting ray.
// The defining trait is fly-over: tiles traversed between origin and landing do NOT fire
// positional-effect triggers — only the landing tile fires (OnEnter). This lets a dash jump
// a trap, or a kick push a target over one.
//
// @spec-link [[mech_movement_reposition]]

// repositionDistanceOf returns the signed tile displacement carried by the skill's effect.
// 0 (or absence) means the skill does not reposition.
func repositionDistanceOf(sk skill.Skill) int {
	p := sk.Effect.GetProperty(property.RepositionDistance)
	if p == nil {
		return 0
	}
	return p.(property.IntProperty).I()
}

// repositionSubjectOf returns who the skill displaces. Absence defaults to Self.
func repositionSubjectOf(sk skill.Skill) def.RepositionSubjectType {
	p := sk.Effect.GetProperty(property.RepositionSubject)
	if p == nil {
		return def.RepositionSubjectSelf
	}
	return def.RepositionSubjectType(p.Get().(string))
}

// isReposition reports whether the skill repositions a subject.
func isReposition(sk skill.Skill) bool {
	return repositionDistanceOf(sk) != 0
}

// sign returns -1, 0, or 1 matching the sign of x.
func sign(x int) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

// unitStep returns the per-axis unit direction from→to on the XY plane (Z ignored).
// This yields both cardinal and diagonal directions.
func unitStep(from, to position.Position) position.Position {
	return position.New(sign(to.X-from.X), sign(to.Y-from.Y), 0)
}

// repositionLanding computes the landing tile from a starting position, a forward unit
// direction, and a signed distance. Negative distance reverses the direction (pull). The
// landing keeps the subject's original height.
func repositionLanding(from, forward position.Position, dist int) position.Position {
	if dist < 0 {
		forward = position.New(-forward.X, -forward.Y, 0)
		dist = -dist
	}
	return position.New(from.X+forward.X*dist, from.Y+forward.Y*dist, from.Z)
}

// repositionVector resolves the (origin, forward) for a subject. For Self the caster is the
// subject and the direction is caster→target; for Target the subject is the targeted entity
// and the direction is caster→subject.
func repositionVector(subject def.RepositionSubjectType, casterPos, subjectPos, target position.Position) (origin, forward position.Position) {
	if subject == def.RepositionSubjectSelf {
		return casterPos, unitStep(casterPos, target)
	}
	return subjectPos, unitStep(casterPos, subjectPos)
}

// validLanding reports whether a tile is a legal landing target: in-grid, walkable, and not
// occupied by a blocking entity other than the subject itself.
func (ctx *localSkillCtx) validLanding(pos position.Position, selfID uuid.UUID) (bool, string) {
	c, ok := ctx.Grid.CellAt(pos)
	if !ok {
		return false, "skill.reposition.outofgrid"
	}
	if c.Type != cell.Ground && c.Type != cell.Dirt {
		return false, "skill.reposition.blocked"
	}
	if ctx.HasBlockingEntity(c.EntityIDs, selfID) {
		return false, "skill.reposition.blocked"
	}
	return true, ""
}

// checkReposition validates that every subject of a reposition skill has a legal landing tile.
// It runs in preSkillChecks (before any effect is applied) so a blocked reposition fails cleanly.
func (ctx *localSkillCtx) checkReposition(msg *message.Message, caster entity.Entity, sk skill.Skill, target position.Position) (bool, *message.Message) {
	if !isReposition(sk) {
		return true, msg
	}
	dist := repositionDistanceOf(sk)
	subject := repositionSubjectOf(sk)

	subjects := []entity.Entity{caster}
	if subject == def.RepositionSubjectTarget {
		subjects = ctx.targetedEntities
	}

	for _, s := range subjects {
		origin, forward := repositionVector(subject, caster.Position, s.Position, target)
		if forward.X == 0 && forward.Y == 0 {
			return false, msg.ReplyWithError("Reposition needs a direction", "skill.reposition.nodirection")
		}
		if ok, errkey := ctx.validLanding(repositionLanding(origin, forward, dist), s.ID); !ok {
			return false, msg.ReplyWithError("Reposition landing is blocked", errkey)
		}
	}
	return true, msg
}

// applyReposition moves the skill's subject(s) after effects have been applied, firing only the
// landing tile's OnEnter triggers. It returns the (possibly moved) caster.
func (ctx *localSkillCtx) applyReposition(caster entity.Entity, sk skill.Skill, target position.Position) entity.Entity {
	if !isReposition(sk) {
		return caster
	}
	dist := repositionDistanceOf(sk)
	subject := repositionSubjectOf(sk)

	if subject == def.RepositionSubjectSelf {
		_, forward := repositionVector(subject, caster.Position, caster.Position, target)
		ctx.moveSubject(&caster, repositionLanding(caster.Position, forward, dist))
		if live, ok := ctx.Entities[caster.ID]; ok {
			caster = live
		}
		return caster
	}

	for _, te := range ctx.targetedEntities {
		live, ok := ctx.Entities[te.ID]
		if !ok {
			continue // subject died from the effect
		}
		_, forward := repositionVector(subject, caster.Position, live.Position, target)
		ctx.moveSubject(&live, repositionLanding(live.Position, forward, dist))
	}
	return caster
}

// moveSubject relocates an entity on the grid to the landing tile and fires the landing tile's
// OnEnter positional effects. Intermediate (flown-over) tiles are intentionally skipped.
func (ctx *localSkillCtx) moveSubject(ent *entity.Entity, landing position.Position) {
	from := ent.Position
	if from.Equals(landing) {
		return
	}
	if err := ctx.Grid.MoveEntity(from, landing, ent.ID); err != nil {
		ctx.log.WithError(err).Error("reposition: grid move failed")
		return
	}
	ent.Position = landing
	ctx.Entities[ent.ID] = *ent

	// Fly-over semantics: only the landing tile fires (OnEnter); no OnExit on origin,
	// no OnEnter/OnStep on traversed tiles. This is the divergence from Move().
	ProcessPositionalEffects(ctx.GameState, *ent, landing, property.TriggerOnEnter)
}
