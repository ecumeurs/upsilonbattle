---
id: rule_forfeit_battle
human_name: "Forfeit Battle Rule"
type: RULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: [combat, forfeit, resolution]
parents:
  - [[rule_team_mechanics]]
  - [[shared:us_take_combat_turn]]
dependents: []
---

# Forfeit Battle Rule

## INTENT
To allow a player to concede a match, resulting in an immediate victory for the opposing side(s) and proper arena closure.

## THE RULE / LOGIC
- A player may declare "FORFEIT" at any point once the arena has reached the **in_progress** state of [[specification_arena_lifecycle]] — i.e. from the moment the match's first engine tick has run until the match concludes. A forfeit request made earlier, while the arena is still in the created or starting state, is rejected (400 game.not.in.progress) as an illegal state transition — this is defined, compliant behavior, not a bug, and it is not a silent no-op or a queued retry.
- **PvE Resolution (Single Player vs AI):**
  - If the human player forfeits, the battle arena is closed immediately.
  - The human player is marked as "DEFEATED".
- **PvP Resolution (Multi-Player):**
  - If a player forfeits, all entities belonging to that player's `TeamID` are considered to have surrendered.
  - The forfeiting team is marked as "DEFEATED".
  - Victory is handed to the remaining team(s) with active entities.
  - The battle ends immediately and the `winner_team_id` is broadcast to all clients via a `BattleEnd` event.
  - In a 2v2 scenario, the forfeit of any single player on a team covers the entire team.

**Prior wording (superseded 2026-08-24, ISS-102):** the first bullet previously read: "A player may declare \"FORFEIT\" at any time during the match." That wording was ambiguous as to whether the pre-tick setup window (arena created but not yet ticking) counted as part of "the match" — [[specification_arena_lifecycle]] now resolves that ambiguity explicitly: forfeit is legal from **in_progress** onward, not from **created**/**starting**.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /api/v1/game/:id/forfeit` (Standalone Route)
- **Engine Call:** `controllerForfeit`
- **Code Tag:** `@spec-link [[rule_forfeit_battle]]`
- **Related Issue:** `#ISS-003`

## EXPECTATION (For Testing)
- Player A (Team 1) forfeits -> System broadcasts `BattleEnd` with `winner_team_id: 2` -> Arena closed.
- In PvE, Player A forfeits -> System broadcasts `BattleEnd` where `winner_team_id` is the computer team (e.g. 2).
