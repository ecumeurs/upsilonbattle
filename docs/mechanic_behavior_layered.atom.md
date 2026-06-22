---
id: mechanic_behavior_layered
status: STABLE
parents:
  - [[mech_behavior_system]]
dependents: []
human_name: Layered Behavior Pipeline
type: MECHANIC
layer: IMPLEMENTATION
priority: 4
version: 1.1
tags: [ai, behavior, pipeline, archetypes]
---

# Layered Behavior Pipeline

## INTENT
To evaluate an ordered stack of Behavior layers per tick, producing a single EngineCommand via first-writer-wins slot resolution and grade-scaled activation rolls.

## THE RULE / LOGIC
**Pipeline model:** `LayeredBehavior.Tick(ctx) → EngineCommand`

Layers are evaluated top to bottom. Each layer rolls its activation rate against a grade-scaled threshold before proposing. If the roll fails, the layer is skipped.

**Activation curve:** `effectiveActivation(base, gradeIdx) = clamp(base * (0.4 + gradeIdx * 0.075), 0, 1)`. Base ≥ 1.0 bypasses the roll entirely (always activates). The baseline AggressiveBehavior uses `BaseActivation = 1.0`.

**Slot model:** Each tick has a shared `DecisionDraft` with three slots: `Target`, `Move`, `Action`. First layer to write a slot wins; subsequent layers skip already-set slots.

**Resolution priority:** `Action` (if not acted) → `Move` (if movement budget remains) → `EndOfTurn`.

**Baseline guarantee:** The bottom layer of every archetype stack is `AggressiveBehavior` with `BaseActivation = 1.0`, ensuring a decision is always reached.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_behavior_layered]]`
- **Key types:** `LayeredBehavior`, `Behavior`, `DecisionDraft`, `GameContext`, `EngineCommand`
- **Package:** `upsilonbattle/battlearena/controller/behavior/`
- **Parents:** `[[mechanic_ai_controller_archetypes]]`

## EXPECTATION
- A grade-I entity activates non-baseline layers ~40% of the time; grade-V activates them 100%.
- The baseline layer always activates regardless of grade.
- If all non-baseline layers skip, the baseline fills all three slots.
- Resolution emits exactly one EngineCommand per tick.
