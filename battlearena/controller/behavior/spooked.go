package behavior

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

// SpookedBehavior drives an entity to flee from the nearest enemy rather than attack.
// Used for non-combat temporary entities that should keep their distance.
//
// @spec-link [[mech_behavior_system]]
type SpookedBehavior struct{}

func (b *SpookedBehavior) OnTurn(ctx GameContext) Decision {
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()

	selfTeam := self.GetPropertyI(property.TeamID).I()
	selfHP := self.GetPropertyI(property.HP).I()
	if selfHP <= 0 {
		return Decision{Type: DecisionPass}
	}

	// Find nearest living enemy
	var nearestEnemy *position.Position
	bestDist := int(^uint(0) >> 1)

	for _, ent := range entities {
		if ent.ID == self.ID {
			continue
		}
		entHP := ent.GetPropertyI(property.HP).I()
		entTeam := ent.GetPropertyI(property.TeamID).I()
		if entHP <= 0 || entTeam == selfTeam {
			continue
		}
		dist := self.Position.Distance(ent.Position)
		if dist < bestDist {
			bestDist = dist
			pos := ent.Position
			nearestEnemy = &pos
		}
	}

	if nearestEnemy == nil {
		return Decision{Type: DecisionPass}
	}

	// Move in the opposite direction (flee logic delegated to controller pathfinding)
	// The controller uses this decision and computes the actual flee path.
	// We encode flee direction as a negative-target convention (invert target from self).
	fleeTarget := position.New(
		2*self.Position.X-nearestEnemy.X,
		2*self.Position.Y-nearestEnemy.Y,
		self.Position.Z,
	)
	return Decision{Type: DecisionMove, Target: fleeTarget}
}
