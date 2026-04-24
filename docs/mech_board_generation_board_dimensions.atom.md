---
id: mech_board_generation_board_dimensions
human_name: Board Dimensions Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_board_generation]]
dependents:
  - [[ui_iso_board]]
---
# Board Dimensions Mechanic

## INTENT
Defines the constraints for the width and height of the board.

## THE RULE / LOGIC
Dimensions: The board is a standard grid (rectangle). Its width and height must each be randomly rolled between 5 and 15 tiles inclusive.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_board_generation_board_dimensions]]`
