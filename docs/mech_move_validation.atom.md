---
id: mech_move_validation
human_name: Entity Move Validation Mechanic
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[api_plan_travel_toward]]
  - [[mech_move_validation_move_validation_already_moved]]
  - [[mech_move_validation_move_validation_controller_mismatch]]
  - [[mech_move_validation_move_validation_entity_collision]]
  - [[mech_move_validation_move_validation_existence]]
  - [[mech_move_validation_move_validation_jump_limitations]]
  - [[mech_move_validation_move_validation_obstacle_collision]]
  - [[mech_move_validation_move_validation_path_adjacency]]
  - [[mech_move_validation_move_validation_path_length_credits]]
  - [[mech_move_validation_move_validation_turn_mismatch]]
---
# Entity Move Validation Mechanic

## INTENT
To aggregate the constituent rules of Entity Move Validation Mechanic.

## THE RULE / LOGIC
Validates the rules governing entity movements before applying them to the game state.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation]]`
