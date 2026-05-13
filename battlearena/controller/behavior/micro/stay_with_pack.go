package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontools/tools"
)

const packRadius = 3

// StayWithPack proposes moving toward the nearest ally when the entity is isolated
// (no ally within packRadius cells). Helps keep the team together.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type StayWithPack struct{}

// Name returns the stable identifier used in memory records and logs.
func (*StayWithPack) Name() string             {
	return "stay_with_pack"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*StayWithPack) BaseActivation() float64  {
	return 0.6
}

// Propose moves toward the nearest ally when self is isolated beyond the pack threshold.
func (b *StayWithPack) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Move != nil || ctx.RemainingMovement() <= 0 {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	// Check isolation: is any ally within packRadius cells?
	for _, ent := range entities {
		if ent.ID == self.ID || ent.GetPropertyI(property.TeamID).I() != selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		if tools.Distance(self.Position.X, self.Position.Y, ent.Position.X, ent.Position.Y) <= packRadius {
			return // not isolated
		}
	}

	// Isolated: find nearest ally and move toward them.
	var nearestAlly *entity.Entity
	bestDist := int(^uint(0) >> 1)
	for _, ent := range entities {
		if ent.ID == self.ID || ent.GetPropertyI(property.TeamID).I() != selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		d := tools.Distance(self.Position.X, self.Position.Y, ent.Position.X, ent.Position.Y)
		if d < bestDist {
			bestDist = d
			e := ent
			nearestAlly = &e
		}
	}
	if nearestAlly == nil {
		return
	}

	grd := ctx.Grid()
	if grd == nil {
		return
	}
	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	path, found := grd.AStarPath(self.Position, nearestAlly.Position, jumpHeight, nil)
	if !found || len(path) <= 1 {
		return
	}
	budget := ctx.RemainingMovement()
	// Stop one cell short of the ally.
	limit := len(path) - 1
	if limit > budget {
		limit = budget
	}
	if limit < 2 {
		return
	}
	draft.Move = &behavior.MoveSlot{Path: path[1:limit]}
}

var _ behavior.Behavior = (*StayWithPack)(nil)
