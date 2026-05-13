package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/property"
)

// FocusBackline proposes an enemy marked as "support" or "ranger" as the primary target.
// Let the next layer or baseline fill other slots.
// Used by the Ranger and Sneak archetypes.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
// @spec-link [[mechanic_ai_controller_archetypes]]
type FocusBackline struct{}

// Name returns the stable identifier used in memory records and logs.
func (*FocusBackline) Name() string             {
	return "focus_backline"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*FocusBackline) BaseActivation() float64  {
	return 0.7
}

// Propose targets the nearest enemy with a Support or Ranger archetype label.
func (b *FocusBackline) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if draft.Target != nil {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	var backlineTarget *entity.Entity
	for _, ent := range entities {
		if ent.ID == self.ID || ent.GetPropertyI(property.TeamID).I() == selfTeam {
			continue
		}
		if ent.GetPropertyI(property.HP).I() <= 0 {
			continue
		}
		if isBacklineArchetype(ent) {
			e := ent
			backlineTarget = &e
			break
		}
	}

	if backlineTarget == nil {
		return
	}
	draft.Target = &behavior.TargetSlot{EntityID: backlineTarget.ID}
}

// isBacklineArchetype returns true if the entity is marked as a support or ranger archetype.
func isBacklineArchetype(ent entity.Entity) bool {
	if !ent.HasProperty(property.AIArchetype) {
		return false
	}
	archetype, ok := ent.GetProperty(property.AIArchetype).Get().(string)
	if !ok {
		return false
	}
	return archetype == "support" || archetype == "ranger"
}

var _ behavior.Behavior = (*FocusBackline)(nil)
