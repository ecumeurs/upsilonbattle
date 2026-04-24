---
id: mech_behavior_system
status: DRAFT
priority: 5
version: 2.0
parent: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [ai, behavior, controller]
human_name: Behavior System
---

# Behavior System

## INTENT

To separate AI decision-making logic from controller plumbing, enabling pluggable, testable, and composable entity behaviors.

## THE RULE / LOGIC

**Separation of Concerns:**

Controllers are responsible for:
- Entity ownership (ControllerID)
- Communication with Ruler (actor pattern)
- Message routing and replies

Behaviors are responsible for:
- Decision-making (what to do)
- State tracking (last damage, preferences)
- Tactical logic (targeting, positioning)

**Behavior Interface:**

The behavior interface defines how AI decisions are made:

- `OnTurn()`: Called when entity's turn begins, returns a decision
- `OnDamaged()`: Called when entity takes damage, for reactive behavior
- `GetName()`: Returns behavior name for debugging

**Game Context:**

Behaviors receive a game context that provides:
- Access to entity by ID
- Access to all entities
- Access to grid state
- Current turn information

**Decision Types:**

Behaviors return decisions that indicate what to do:

- `NoDecision`: Behavior has no opinion, let next behavior try
- `Move`: Move to a position along a path
- `Attack`: Attack a target position
- `Skill`: Use a skill at a target position
- `Pass`: End turn without action
- `Flee`: Move away from threats

**Decision Structure:**

A decision contains:
- Type of action to take
- Target position (for Attack, Skill, Flee)
- Skill ID (for Skill)
- Movement path (for Move, Flee)
- Priority (for AND mode composites)

**Controller Integration:**

Controllers are extended to hold a Behavior instance. When the controller receives a turn notification:

1. Verify the entity belongs to this controller
2. Create game context from known state
3. Call `Behavior.OnTurn()` to get decision
4. Execute the decision by sending appropriate request to Ruler
5. Handle the response

**Example Behavior - Aggressive:**

The aggressive behavior implements standard attack AI:

- Find nearest enemy entity
- If in attack range, attack
- If not in range, move toward target
- If no valid target, pass turn
- On damage, remember the attacker for future targeting

**Behavior State:**

Behaviors can maintain state between calls:
- Preferred target (after being damaged)
- Last known positions
- Tactical preferences

**Testing:**

Behaviors are testable independently of controllers:
- Provide mock game context
- Call `OnTurn()` with test entity
- Verify returned decision

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_behavior_system]]`
- **Related Files:** `upsilonbattle/battlearena/controller/behavior/behavior.go` (new), `upsilonbattle/battlearena/controller/controller.go`
- **Integration:** Works with `mech_composite_behavior` for complex behavior combinations

## EXPECTATION
