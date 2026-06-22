---
id: mech_behavior_system
status: STABLE
priority: 5
version: 2.0
dependents:
  - [[mech_composite_behavior]]
  - [[mech_controller_behavior]]
  - [[mech_ruler_behavior]]
  - [[mech_turret_entity]]
  - [[mechanic_ai_controller_archetypes]]
  - [[mechanic_ai_progression_matching]]
  - [[mechanic_behavior_layered]]
  - [[mechanic_decision_memory]]
type: MODULE
layer: ARCHITECTURE
tags: [ai, behavior, controller]
human_name: Behavior System
parents:
  - [[shared:req_tech_debt_backlog]]
---

# Behavior System

## INTENT
To provide a unified framework for entity automation across both synchronous engine logic and asynchronous actor controllers.

## THE RULE / LOGIC
The Upsilon Battle engine supports two distinct automation paths depending on entity complexity and execution context:

1. **Ruler Behavior (`[[mech_ruler_behavior]]`):** Low-overhead, synchronous automation for simple effects (traps, turrets).
2. **Controller Behavior (`[[mech_controller_behavior]]`):** High-level, asynchronous AI for complex actors (heroes, bosses).

This parent atom governs the high-level separation of concerns between these two subsystems.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_behavior_system]]`
- **Related Files:** `upsilonbattle/battlearena/controller/behavior/behavior.go` (new), `upsilonbattle/battlearena/controller/controller.go`
- **Integration:** Works with `mech_composite_behavior` for complex behavior combinations

## EXPECTATION
- An entity routed through Ruler Behavior is resolved synchronously within the engine tick (no actor message round-trip) for simple effects (traps, turrets).
- An entity routed through Controller Behavior is driven by an asynchronous actor controller for complex actors (heroes, bosses).
- Every automated entity resolves through exactly one of the two paths; neither path leaves an entity without a decision.
- Composite behaviors can be assembled on top of either path (see `[[mech_composite_behavior]]`).
