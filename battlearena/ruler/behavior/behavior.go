package behavior

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilonmapdata/grid/position/pattern"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GameState interface to avoid circular dependency
type GameState interface {
	GetEntities() map[uuid.UUID]entity.Entity
	GetGrid() *grid.Grid
	GetLogger() *logrus.Entry
}

// Behavior defines how an automated entity (Trap, TimeBased, etc.) decides its action.
// @spec-link [[mech_behavior_system]]
type Behavior interface {
	Decide(gs GameState, ent entity.Entity) *message.Message
}

// ExpirationBehavior is the default behavior for entities that just wait to die.
// They always return EndOfTurn.
// @spec-link [[mech_entity_expiration]]
type ExpirationBehavior struct{}

// Decide returns an EndOfTurn message, as expiration behavior simply waits for the entity to expire.
func (b *ExpirationBehavior) Decide(gs GameState, ent entity.Entity) *message.Message {
	return message.Create(nil, rulermethods.EndOfTurn{
		EntityID:     ent.ID,
		ControllerID: ent.ControllerID,
	}, &rulermethods.EndOfTurn{})
}

// AggressiveBehavior finds the nearest enemy and moves toward them.
// @spec-link [[mech_behavior_system]]
type AggressiveBehavior struct{}

// Decide selects an action (Move or Attack) based on the proximity of the nearest enemy entity.
func (b *AggressiveBehavior) Decide(gs GameState, ent entity.Entity) *message.Message {
	entities := gs.GetEntities()
	myTeam := ent.GetPropertyI(property.TeamID).I()

	var nearestEnemy *entity.Entity
	minDist := 999999

	for _, other := range entities {
		if other.ID == ent.ID {
			continue
		}
		if other.GetPropertyI(property.TeamID).I() == myTeam {
			continue
		}
		// Only attack characters/monsters
		if other.Type != entity.Character && other.Type != entity.Monster {
			continue
		}

		dist := ent.Position.Distance(other.Position)
		if dist < minDist {
			minDist = dist
			target := other
			nearestEnemy = &target
		}
	}

	if nearestEnemy == nil {
		return (&ExpirationBehavior{}).Decide(gs, ent)
	}

	// For now, aggressive behavior just moves one step toward the enemy if not adjacent.
	// In a real implementation, it would use pathfinding and check for attack range.
	// Since we want to keep it simple for now, we just move one step.
	
	if ent.Position.IsAdjacent(nearestEnemy.Position, ent.GetPropertyI(property.JumpHeight).I()) {
		// Already adjacent, attack!
		return message.Create(nil, rulermethods.ControllerAttack{
			EntityID:     ent.ID,
			ControllerID: ent.ControllerID,
			Target:       nearestEnemy.Position,
		}, &rulermethods.ControllerAttackReply{})
	}

	// Move one step toward enemy
	path := []position.Position{}
	// Find adjacent position to ent that is closer to nearestEnemy
	bestPos := ent.Position
	for _, adj := range gs.GetGrid().SelectPositionsByPattern2D(ent.Position, pattern.Neighbours2D()) {
		dist := adj.Distance(nearestEnemy.Position)
		if dist < bestPos.Distance(nearestEnemy.Position) {
			// Check if cell is occupied
			if c, ok := gs.GetGrid().CellAt(adj); ok {
				// Simplified occupancy check: just check if any entity is there
				if len(c.EntityIDs) == 0 {
					bestPos = adj
				}
			}
		}
	}

	if !bestPos.Equals(ent.Position) {
		path = append(path, bestPos)
		return message.Create(nil, rulermethods.ControllerMove{
			EntityID:     ent.ID,
			ControllerID: ent.ControllerID,
			Path:         path,
		}, &rulermethods.ControllerMoveReply{})
	}

	return (&ExpirationBehavior{}).Decide(gs, ent)
}

// GetBehavior returns a behavior implementation by its slug.
func GetBehavior(slug string) Behavior {
	switch slug {
	case "aggressive":
		return &AggressiveBehavior{}
	default:
		return &ExpirationBehavior{}
	}
}
