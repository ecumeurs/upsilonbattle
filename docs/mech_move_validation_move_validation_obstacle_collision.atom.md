---
id: mech_move_validation_move_validation_obstacle_collision
human_name: Obstacle Collision Rule
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
# Obstacle Collision Rule

## INTENT
No node within the path can be classified as an Obstacle.

## THE RULE / LOGIC
entity.path.obstacle

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation_move_validation_obstacle_collision]]`
