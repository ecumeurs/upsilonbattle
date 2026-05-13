// Package micro shared helpers used across micro-behavior files.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
package micro

import (
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/google/uuid"
)

// skillTag is the set of skill property tags understood by micro-behaviors.
// Maps a string tag to the corresponding SkillProperties key.
var skillTagMap = map[string]property.SkillProperties{
	"heal":   property.Heal,
	"shield": property.ShieldPower,
}

// maxHPOf returns the maximum HP of an entity via the counter property's max value.
func maxHPOf(ent entity.Entity) int {
	return ent.GetPropertyC(property.HP).GetMaxValue()
}

// nearestLivingEnemy returns the closest living enemy by Manhattan distance, or nil.
func nearestLivingEnemy(self entity.Entity, selfTeam int, entities map[uuid.UUID]entity.Entity) *entity.Entity {
	var best *entity.Entity
	bestDist := int(^uint(0) >> 1)
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
		d := tools.Distance(self.Position.X, self.Position.Y, ent.Position.X, ent.Position.Y)
		if d < bestDist {
			bestDist = d
			e := ent
			best = &e
		}
	}
	return best
}

// lowestHPEnemy returns the living enemy with the fewest HP, or nil.
func lowestHPEnemy(self entity.Entity, selfTeam int, entities map[uuid.UUID]entity.Entity) *entity.Entity {
	var best *entity.Entity
	bestHP := int(^uint(0) >> 1)
	for _, ent := range entities {
		if ent.ID == self.ID {
			continue
		}
		if ent.GetPropertyI(property.TeamID).I() == selfTeam {
			continue
		}
		hp := ent.GetPropertyI(property.HP).I()
		if hp <= 0 {
			continue
		}
		if hp < bestHP {
			bestHP = hp
			e := ent
			best = &e
		}
	}
	return best
}

// lowestHPAlly returns the living ally (excluding self) with the fewest current HP, or nil.
func lowestHPAlly(self entity.Entity, selfTeam int, entities map[uuid.UUID]entity.Entity) *entity.Entity {
	var best *entity.Entity
	bestHP := int(^uint(0) >> 1)
	for _, ent := range entities {
		if ent.ID == self.ID {
			continue
		}
		if ent.GetPropertyI(property.TeamID).I() != selfTeam {
			continue
		}
		hp := ent.GetPropertyI(property.HP).I()
		if hp <= 0 {
			continue
		}
		if hp < bestHP {
			bestHP = hp
			e := ent
			best = &e
		}
	}
	return best
}

// hpPercent returns the HP percentage [0, 100] for an entity (0 if maxHP is 0).
func hpPercent(ent entity.Entity) int {
	hp := ent.GetPropertyI(property.HP).I()
	maxHP := maxHPOf(ent)
	if maxHP <= 0 {
		return 0
	}
	return hp * 100 / maxHP
}

// hasSkillWithTag returns the first equipped skill UUID matching a named tag ("heal", "shield"),
// or uuid.Nil when no such skill is equipped.
func hasSkillWithTag(ent entity.Entity, tag string) uuid.UUID {
	prop, ok := skillTagMap[tag]
	if !ok {
		return uuid.Nil
	}
	for id, sk := range ent.Skills {
		if sk.HasProperty(prop) {
			return id
		}
	}
	return uuid.Nil
}
