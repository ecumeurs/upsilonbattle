---
id: mech_board_generation_terrain_obstacles
human_name: Terrain Obstacle Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_board_generation]]
dependents: []
---
# Terrain Obstacle Mechanic

## INTENT
Designates impassable tiles on the generated board.

## THE RULE / LOGIC
Terrain Obstacles: After dimension generation, randomly selected tiles will be designated as impassable "obstacles."

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_board_generation_terrain_obstacles]]`
