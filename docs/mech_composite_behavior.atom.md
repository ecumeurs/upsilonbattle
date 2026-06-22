---
id: mech_composite_behavior
status: DRAFT
priority: 5
version: 2.0
parents:
  - [[mech_behavior_system]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [ai, behavior, composite]
human_name: Composite Behavior
---

# Composite Behavior

## INTENT

To enable complex AI behavior by combining multiple behaviors using OR (first-match) or AND (all-execute) composition modes.

## THE RULE / LOGIC

**Composite Modes:**

**OR Mode (First Match Wins):**
- Behaviors are tried in priority order
- First behavior to return a non-NoDecision wins
- Later behaviors are never consulted
- Use case: "If fleeing, don't attack. If not fleeing, attack."

**AND Mode (All Execute):**
- All behaviors execute in priority order
- Later decisions override earlier decisions
- All behaviors get a chance to influence the outcome
- Use case: "Apply status effects, then execute aggressive behavior."

**Composite Behavior Structure:**

A composite behavior contains:
- Mode (OR or AND)
- List of behavior wrappers, each with:
  - Behavior instance
  - Priority (higher = checked/executed first)

**OR Mode Logic:**

In OR mode:
1. Sort behaviors by priority (highest first)
2. For each behavior:
   - Call `OnTurn()`
   - If decision is not NoDecision, return it
   - Otherwise, continue to next behavior
3. If no behavior made a decision, return Pass

**AND Mode Logic:**

In AND mode:
1. Sort behaviors by priority (highest first)
2. Initialize final decision as NoDecision
3. For each behavior:
   - Call `OnTurn()`
   - If decision is not NoDecision, override final decision
4. Return final decision

**OnDamaged Propagation:**

In both OR and AND mode:
- Call `OnDamaged()` on all behaviors
- All behaviors get notified of damage events
- Order: highest priority first

**Priority System:**

- Higher priority = checked/executed first
- In OR mode: first priority to make a decision wins
- In AND mode: higher priority decisions may be overridden by lower priority depending on sorting
- Use consistent priority values (10 = high, 5 = medium, 1 = low)

**Example - Spooked + Aggressive (OR Mode):**

Composite with two behaviors:
1. SpookedBehavior (priority 10): Flee if HP below threshold
2. AggressiveBehavior (priority 5): Attack nearest foe

Flow:
- Spooked checks if HP < threshold
- If yes: return Flee decision (wins, Aggressive never runs)
- If no: return NoDecision
- Aggressive runs: find nearest foe, attack or move toward

**Example - Expiration + Aggressive (AND Mode):**

Composite with two behaviors:
1. ExpirationBehavior (priority 5): Check expiration (usually no decision)
2. AggressiveBehavior (priority 5): Execute aggressive behavior

Flow:
- Expiration runs: checks duration, returns NoDecision
- Aggressive runs: finds target, attacks
- Final decision is Attack

**Edge Cases:**

- **All behaviors return NoDecision**: Composite returns Pass
- **One behavior returns Pass**: In OR mode, Pass wins. In AND mode, later behaviors can override.
- **Conflicting decisions**: In AND mode, last decision wins (by priority order)

**Behavior Wrapper:**

Each behavior is wrapped with:
- Priority value
- Reference to behavior instance

Wrappers allow the same behavior instance to be used in multiple composites with different priorities.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_composite_behavior]]`
- **Related Files:** `upsilonbattle/battlearena/controller/behavior/composite.go` (new)
- **Integration:** Depends on `mech_behavior_system`

## EXPECTATION
- In OR mode, behaviors are evaluated highest-priority first and the first non-NoDecision result is returned; lower-priority behaviors are not consulted.
- In OR mode where no behavior decides, the composite returns Pass.
- In AND mode, all behaviors run highest-priority first and each non-NoDecision result overrides the prior final decision; the last (lowest-priority) decision wins.
- `OnDamaged()` is propagated to every wrapped behavior in both modes, highest priority first.
- Spooked(priority 10) + Aggressive(priority 5) in OR mode: below-threshold HP yields a Flee decision and Aggressive never runs; otherwise Aggressive attacks.
- Expiration(5) + Aggressive(5) in AND mode: Expiration returns NoDecision, Aggressive's Attack becomes the final decision.
- The same behavior instance can be wrapped in multiple composites with different priorities.
