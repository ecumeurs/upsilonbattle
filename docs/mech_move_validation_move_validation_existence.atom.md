---
id: mech_move_validation_move_validation_existence
human_name: Entity Existence Rule
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_move_validation]]
dependents: []
---
# Entity Existence Rule

## INTENT
The entity must exist within the game state.

## THE RULE / LOGIC
entity.notfound

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation_move_validation_existence]]`
