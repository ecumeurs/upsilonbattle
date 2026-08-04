---
id: mech_character_reroll
human_name: Character Reroll Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[module_game]]
dependents: []
---
# Character Reroll Mechanic

## INTENT
To let a player re-randomize their entire character roster while their account has not yet played a match, up to a fixed number of attempts.

## THE RULE / LOGIC
**Character Reroll Mechanic.**

- **Availability:** the reroll is allowed only while the account has not yet played its first match — not tied to any specific screen or one-time flow; it is reachable from the account's profile at any point before that first match resolves.
- **Effect:** a reroll completely discards the current 3-character set and mathematically rolls a brand-new set, so all 3 characters receive new stats.
- **Limit:** a player may reroll at most **3 times** total; a `reroll_count` (default 0, tracked in player_stats) is checked against the limit before each attempt and incremented on use. Once the account's first match resolves (win or loss), reroll is permanently unavailable regardless of remaining count.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_character_reroll]]`
