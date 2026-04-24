---
id: mech_board_generation_min_area_constraint
human_name: Minimum Area Constraint Mechanic
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
# Minimum Area Constraint Mechanic

## INTENT
Ensures the total area of the board meets or exceeds a minimum threshold.

## THE RULE / LOGIC
Minimum Size Constraint: The total area (width × height) of the rolled board must be greater than or equal to 50 tiles.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_board_generation_min_area_constraint]]`
