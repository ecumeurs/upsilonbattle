// Package micro contains archetype micro-behaviors.
// Each file is one Behavior that proposes to one or more slots of a DecisionDraft.
//
// Convention: check draft.Slot == nil before proposing (first-writer-wins).
// Self-throttle with ctx.Memory().TurnsSince(Name()) when needed.
package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// BackstabSeeker proposes a target whose rear is already exposed to this entity,
// and proposes a flanking move to reach the rear when not already there.
// Used by the Sneak archetype.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type BackstabSeeker struct{}

// Name returns the stable identifier used in memory records and logs.
func (*BackstabSeeker) Name() string             {
	return "backstab_seeker"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*BackstabSeeker) BaseActivation() float64  {
	return 0.8
}

// Propose routes self to a rear-facing position behind a foe to set up backstab targeting.
func (b *BackstabSeeker) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()

	selfTeam := self.GetPropertyI(property.TeamID).I()

	var bestTarget *entity.Entity
	for _, ent := range entities {
		if ent.ID == self.ID {
			continue
		}
		if ent.GetPropertyI(property.TeamID).I() == selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		// Prefer a target we are already backstabbing, or can reach from the rear.
		if self.IsBackstabbing(ent) {
			e := ent
			bestTarget = &e
			break
		}
	}

	if bestTarget == nil {
		return
	}

	if draft.Target == nil {
		draft.Target = &behavior.TargetSlot{EntityID: bestTarget.ID}
	}

	// Propose a flanking move toward the rear of the target.
	if draft.Move == nil && ctx.RemainingMovement() > 0 {
		draft.Move = backstabMoveToRear(self, bestTarget, ctx)
	}
}

// backstabMoveToRear plans an A* path toward the rear of the target within the movement budget.
func backstabMoveToRear(self entity.Entity, target *entity.Entity, ctx behavior.GameContext) *behavior.MoveSlot {
	grd := ctx.Grid()
	if grd == nil {
		return nil
	}
	rearPos := rearPosition(target)
	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	path, found := grd.AStarPath(self.Position, rearPos, jumpHeight, nil)
	if !found || len(path) <= 1 {
		return nil
	}
	end := ctx.RemainingMovement() + 1
	if end > len(path) {
		end = len(path)
	}
	return &behavior.MoveSlot{Path: path[1:end]}
}

// rearPosition returns the grid cell directly behind the target based on its Orientation.
func rearPosition(target *entity.Entity) position.Position {
	tx, ty, tz := target.Position.X, target.Position.Y, target.Position.Z
	switch target.Orientation {
	case entity.Up: // facing Up(N) → rear is Down(S)
		return position.New(tx, ty-1, tz)
	case entity.Right: // facing Right(E) → rear is Left(W)
		return position.New(tx-1, ty, tz)
	case entity.Down: // facing Down(S) → rear is Up(N)
		return position.New(tx, ty+1, tz)
	case entity.Left: // facing Left(W) → rear is Right(E)
		return position.New(tx+1, ty, tz)
	}
	return target.Position
}

var _ behavior.Behavior = (*BackstabSeeker)(nil)
