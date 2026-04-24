---
id: mech_matchmaking
human_name: Matchmaking Mechanics
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: [matchmaking, logic, logic-flow]
parents:
  - [[req_matchmaking]]
dependents: []
---
# Matchmaking Mechanics

## INTENT
To define the logical conditions for grouping players into game sessions based on match format.

## THE RULE / LOGIC
- **Queue Membership:** A player enters the queue for a specific `game_mode`.
- **Character Selection:** The system automatically pulls the player's first 3 characters (active roster) upon joining.
- **Match Triggers:**
  - **1v1 PVP:** Requires exactly 2 distinct players in the queue.
  - **1v1 PVE:** Requires 1 player; system generates an AI opponent immediately.
  - **2v2 PVP:** Requires exactly 4 distinct players in the queue. Players are grouped into teams (2 vs 2).
  - **2v2 PVE:** Requires 2 distinct players on the same team vs 2 AI-controlled entities.
- **Session Initiation:** Once the count for the requested `game_mode` is met:
  1. A `GameMatch` entity is created.
  2. The `UpsilonApiService` is called to `startArena`.
  3. Players are removed from the `MatchmakingQueue`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_matchmaking]]`
- **Related Issue:** `#ISS-007`

## EXPECTATION (For Testing)
- For 1v1 PVP, two calls to `joinMatch` from different users should trigger exactly one `startArena` call.
- For 2v2 PVP, four calls to `joinMatch` are required before `startArena` is triggered.
