package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/google/uuid"
)

// StayBehindFront proposes a move that keeps the entity behind the ally frontline
// (i.e. the ally closest to the enemy). If the entity is already behind a front-liner,
// it abstains. Used by the Support archetype to stay protected.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type StayBehindFront struct{}

// Name returns the stable identifier used in memory records and logs.
func (*StayBehindFront) Name() string             {
	return "stay_behind_front"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*StayBehindFront) BaseActivation() float64  {
	return 0.7
}

// Propose moves toward allies when isolated, keeping the entity within support range.
func (b *StayBehindFront) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Move != nil || ctx.RemainingMovement() <= 0 {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	// Find the ally nearest to any enemy (the "front-liner").
	var frontLiner *entity.Entity
	bestAllyEnemyDist := int(^uint(0) >> 1)
	for _, ent := range entities {
		if ent.ID == self.ID || ent.GetPropertyI(property.TeamID).I() != selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		// Distance from this ally to nearest enemy.
		minEnemyDist := int(^uint(0) >> 1)
		for _, other := range entities {
			if other.GetPropertyI(property.TeamID).I() == selfTeam {
				continue
			}
			if other.GetPropertyI(property.HP).I() <= 0 {
				continue
			}
			d := tools.Distance(ent.Position.X, ent.Position.Y, other.Position.X, other.Position.Y)
			if d < minEnemyDist {
				minEnemyDist = d
			}
		}
		if minEnemyDist < bestAllyEnemyDist {
			bestAllyEnemyDist = minEnemyDist
			e := ent
			frontLiner = &e
		}
	}

	if frontLiner == nil {
		return
	}

	// If we are already behind the front-liner (farther from enemies than they are),
	// stay put.
	selfEnemyDist := minDistToEnemy(self, selfTeam, entities)
	if selfEnemyDist >= bestAllyEnemyDist {
		return
	}

	// Move to be 2 cells behind the front-liner (away from enemies).
	grd := ctx.Grid()
	if grd == nil {
		return
	}
	jumpHeight := self.GetPropertyI(property.JumpHeight).I()
	path, found := grd.AStarPath(self.Position, frontLiner.Position, jumpHeight, nil)
	if !found || len(path) <= 2 {
		return
	}
	// Stop 2 cells short of the front-liner.
	limit := len(path) - 2
	if limit > ctx.RemainingMovement() {
		limit = ctx.RemainingMovement()
	}
	if limit < 2 {
		return
	}
	draft.Move = &behavior.MoveSlot{Path: path[1:limit]}
}

func minDistToEnemy(self entity.Entity, selfTeam int, entities map[uuid.UUID]entity.Entity) int {
	best := int(^uint(0) >> 1)
	for _, ent := range entities {
		if ent.GetPropertyI(property.TeamID).I() == selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		d := tools.Distance(self.Position.X, self.Position.Y, ent.Position.X, ent.Position.Y)
		if d < best {
			best = d
		}
	}
	return best
}

var _ behavior.Behavior = (*StayBehindFront)(nil)
