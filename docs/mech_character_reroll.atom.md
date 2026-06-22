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
To let a player re-randomize their entire starting character roster during account creation, up to a fixed number of attempts.

## THE RULE / LOGIC
**Character Reroll Mechanic.**

- **Availability:** the reroll is allowed only while the account is in the creation flow, after the initial 3 characters have been generated.
- **Effect:** a reroll completely discards the current 3-character set and mathematically rolls a brand-new set, so all 3 characters receive new stats.
- **Limit:** a player may reroll at most **3 times** per account creation; a `reroll_count` (default 0, see `[[entity_users]]`) is checked against the limit before each attempt and incremented on use.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_character_reroll]]`
