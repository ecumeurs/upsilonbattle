package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
)

// HealAlly proposes using a heal skill on the ally with the lowest HP percentage
// when that ally is below 50% HP. Sets both Target and Action slots.
// Used by the Support archetype.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type HealAlly struct{}

// Name returns the stable identifier used in memory records and logs.
func (*HealAlly) Name() string             {
	return "heal_ally"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*HealAlly) BaseActivation() float64  {
	return 0.85
}

// Propose uses a heal skill on the allied entity below the critical HP threshold.
func (b *HealAlly) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if ctx.HasActed() {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	// Find a heal skill.
	healSkillID := hasSkillWithTag(self, "heal")
	if healSkillID == [16]byte{} {
		return
	}

	ally := lowestHPAlly(self, selfTeam, entities)
	if ally == nil || hpPercent(*ally) >= 50 {
		return
	}

	if draft.Target == nil {
		draft.Target = &behavior.TargetSlot{EntityID: ally.ID}
	}
	if draft.Action == nil {
		draft.Action = &behavior.ActionSlot{
			Type:    behavior.ActionSkill,
			Target:  ally.Position,
			SkillID: healSkillID,
		}
	}
}

var _ behavior.Behavior = (*HealAlly)(nil)
