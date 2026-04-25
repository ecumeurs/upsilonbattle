---
id: entity_character
human_name: Character Entity
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[entity_character_distribute_remaining_points]]
  - [[rule_character_create_character]]
  - [[rule_progression]]
---
# Character Entity

## INTENT
To aggregate the constituent rules of Character Entity.

## THE RULE / LOGIC
Defines the baseline stat block and attributes of a playable character unit.

Attributes:
* HP (consumable, on the board, in game only). When a character is eliminated, HP is set to 0.
* Max HP (attribute)
* Dead (boolean, in game only): True if the character has been eliminated in the current session.
* Attack (attribute)
* Defense (attribute)
* Move (consumable, on the board, in game only)
* Max Move (attribute)
* Position (on the board, in game only): {x,y}. Note: Deceased characters have their position cleared from the grid.
* Name
* ID
* Player ID (the player that owns this character, UUID assigned to that player for a game, in game only)

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_character]]`
