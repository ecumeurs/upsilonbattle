---
id: mech_character_reroll
human_name: Character Reroll Mechanic
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[mech_character_reroll_availability]]
  - [[mech_character_reroll_effect]]
  - [[mech_character_reroll_limit]]
---
# Character Reroll Mechanic

## INTENT
To aggregate the constituent rules of Character Reroll Mechanic.

## THE RULE / LOGIC
Automatically re-randomize player character roster stats during account creation.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_character_reroll]]`
