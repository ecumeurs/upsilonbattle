---
id: entity_game_match
human_name: Game Match Entity
type: ENTITY
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[data_persistence]]
dependents:
  - [[entity_match_participants]]
  - [[uc_admin_history_management]]
---
# Game Match Entity

## INTENT
To define the data structure and caching metadata for an active or completed gameplay match in the BattleUI schema.

## THE RULE / LOGIC
Defines the `game_matches` table, ensuring match metadata and asynchronous events can be queried independently.

Attributes:
* ID (UUID)
* Game State Cache (JSON of BoardStateDTO)
* Grid Cache (JSON mapping of Grid elements)
* Turn (Integer, current turn number)
* Started At (Timestamp)
* Concluded At (Timestamp, nullable)
* Winning Team ID (Integer, nullable)
* Game Mode (String, e.g., 'pvp', 'pve')

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_game_match]]`
