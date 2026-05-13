---
id: mech_ruler_behavior
human_name: "Ruler Behavior System"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 3
tags: [engine, automation, ruler]
parents:
  - [[mech_behavior_system]]
dependents: []
---

# Ruler Behavior System

## INTENT
To provide low-overhead, synchronous automation for simple entities directly within the battle engine.

## THE RULE / LOGIC
- **Synchronous Execution:** Behaviors are executed within the Ruler's main goroutine during the entity's turn.
- **Direct Messaging:** The behavior must return a `*message.Message` which is immediately dispatched by the Ruler.
- **Non-Actor Entities:** This system is the default for entities where `ControllerID` is null (Traps, Turrets, Simple Summons).
- **Performance Constraint:** Behavior logic must be non-blocking and computationally inexpensive.
- **Expiration Support:** The `ExpirationBehavior` is the core implementation for time-limited entities that do not perform actions (like delayed effect anchors). It returns a simple `EndOfTurn` message to allow the engine to process duration decrement.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_ruler_behavior]]`
- **Location:** `upsilonbattle/battlearena/ruler/behavior/`

## EXPECTATION
Entities without controllers act according to their assigned behavior slug during their initiative turn without spawning additional goroutines.
