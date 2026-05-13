package behavior

import (
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// TickOutcome reports what happened when the previous EngineCommand was processed.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type TickOutcome int

const (
	// TickNone means no previous command has been sent this turn.
	TickNone TickOutcome = iota
	// TickSuccess means the previous command completed without issue.
	TickSuccess
	// TickBlocked means the move was blocked by an obstacle or entity.
	TickBlocked
	// TickTrap means the entity stepped on a trap during the move.
	TickTrap
)

// TargetSlot holds the entity to engage this tick.
// Set by targeting layers (e.g. FocusWeakest, BackstabSeeker, baseline).
type TargetSlot struct {
	EntityID uuid.UUID
}

// MoveSlot holds a movement path toward a destination.
// Set by movement layers (e.g. KiteAway, MaintainRange, baseline).
type MoveSlot struct {
	Path []position.Position
}

// ActionType identifies what kind of action to take.
type ActionType int

const (
	// ActionAttack is a basic weapon attack.
	ActionAttack ActionType = iota
	// ActionSkill uses an equipped skill.
	ActionSkill
)

// ActionSlot holds the action to execute after movement (or in place of movement).
// Set by action layers (e.g. HealAlly, ShieldAlly, baseline).
type ActionSlot struct {
	Type    ActionType
	Target  position.Position
	SkillID uuid.UUID // populated for ActionSkill
}

// DecisionDraft is the shared mutable state built up across a single pipeline tick.
// Layers write to unset slots (first-writer-wins convention).
//
// Notes is a free-form cross-layer commentary map, designed to support future
// "partial proposal" semantics (e.g. "(move, where?)" or "(target, A or B?)") without
// changing the interface — out of scope for v1.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type DecisionDraft struct {
	Target *TargetSlot
	Move   *MoveSlot
	Action *ActionSlot
	Notes  map[string]any
}

// EngineCommandType enumerates what the controller should send to the ruler.
type EngineCommandType int

const (
	// CmdEndOfTurn signals that the entity's turn is over.
	CmdEndOfTurn EngineCommandType = iota
	// CmdMove requests a movement along Path.
	CmdMove
	// CmdAttack requests a basic attack at Target.
	CmdAttack
	// CmdSkill requests skill use at Target.
	CmdSkill
)

// EngineCommand is the final output of LayeredBehavior.Tick.
// The controller translates this into a ruler message (ControllerMove, ControllerAttack, etc.).
// EntityID and ControllerID are injected by the controller, not the behavior.
//
// @spec-link [[mechanic_mech_behavior_layered]]
type EngineCommand struct {
	Type    EngineCommandType
	Path    []position.Position // CmdMove
	Target  position.Position   // CmdAttack / CmdSkill
	SkillID uuid.UUID           // CmdSkill
}

// Resolve picks one EngineCommand from the draft following the resolution rule:
//
//	Action (if has-not-acted) → Move (if movement remaining) → EndOfTurn
//
// @spec-link [[mechanic_mech_behavior_layered]]
func (d *DecisionDraft) Resolve(ctx GameContext) EngineCommand {
	if d.Action != nil && !ctx.HasActed() {
		switch d.Action.Type {
		case ActionAttack:
			return EngineCommand{Type: CmdAttack, Target: d.Action.Target}
		case ActionSkill:
			return EngineCommand{Type: CmdSkill, Target: d.Action.Target, SkillID: d.Action.SkillID}
		}
	}
	if d.Move != nil && ctx.RemainingMovement() > 0 && len(d.Move.Path) > 0 {
		return EngineCommand{Type: CmdMove, Path: d.Move.Path}
	}
	return EngineCommand{Type: CmdEndOfTurn}
}
