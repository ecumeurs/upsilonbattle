---
id: mech_board_generation
human_name: Board Generation Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[module_backend_board_generation]]
dependents:
  - [[battleui:ui_iso_board]]
---
# Board Generation Mechanic

## INTENT
To generate a combat board by rolling its dimensions within bounds, enforcing a minimum area, and designating random impassable obstacle tiles.

## THE RULE / LOGIC
**Board Generation Mechanic.**

- **Dimensions:** the board is a rectangular grid; its width and height are each randomly rolled between **5 and 15 tiles inclusive**.
- **Minimum Area Constraint:** the total area (width × height) of the rolled board must be **≥ 50 tiles**; rolls below this threshold are rejected/re-rolled.
- **Terrain Obstacles:** after dimension generation, randomly selected tiles are designated as impassable "obstacles."

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_board_generation]]`
