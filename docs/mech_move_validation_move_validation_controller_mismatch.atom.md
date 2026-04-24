---
id: mech_move_validation_move_validation_controller_mismatch
human_name: Controller Mismatch Rule
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
# Controller Mismatch Rule

## INTENT
The requested Controller ID must match the owning Controller of the entity.

## THE RULE / LOGIC
entity.controller.missmatch

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation_move_validation_controller_mismatch]]`
