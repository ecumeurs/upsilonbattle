package micro

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/controller/behavior"
	"github.com/ecumeurs/upsilontypes/property"
)

// ShieldAlly proposes applying a shield skill to the most-damaged ally.
// Only fires when the ally is below 70% HP and a shield skill is available.
// Used by the Support archetype.
//
// @spec-link [[mechanic_ai_controller_archetypes]]
type ShieldAlly struct{}

// Name returns the stable identifier used in memory records and logs.
func (*ShieldAlly) Name() string             {
	return "shield_ally"
}
// BaseActivation returns this layer's declared activation probability (0–1, grade-scaled by the pipeline).
func (*ShieldAlly) BaseActivation() float64  {
	return 0.75
}

// Propose applies a shield skill to the most-damaged living ally.
func (b *ShieldAlly) Propose(ctx behavior.GameContext, draft *behavior.DecisionDraft) {
	if ctx.HasActed() || draft.Action != nil {
		return
	}
	self := ctx.SelfEntity()
	entities := ctx.KnownEntities()
	selfTeam := self.GetPropertyI(property.TeamID).I()

	shieldSkillID := hasSkillWithTag(self, "shield")
	if shieldSkillID == [16]byte{} {
		return
	}

	ally := lowestHPAlly(self, selfTeam, entities)
	if ally == nil || hpPercent(*ally) >= 70 {
		return
	}

	if draft.Target == nil {
		draft.Target = &behavior.TargetSlot{EntityID: ally.ID}
	}
	draft.Action = &behavior.ActionSlot{
		Type:    behavior.ActionSkill,
		Target:  ally.Position,
		SkillID: shieldSkillID,
	}
}

var _ behavior.Behavior = (*ShieldAlly)(nil)
