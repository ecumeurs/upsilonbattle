---
id: mech_move_validation_move_validation_already_moved
human_name: Already Moved Rule
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
# Already Moved Rule

## INTENT
The entity must not carry the HasMoved flag set to true.

## THE RULE / LOGIC
entity.movement.already

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation_move_validation_already_moved]]`
