---
id: mech_character_reroll_limit
human_name: Reroll Limit
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_character_reroll]]
dependents: []
---
# Reroll Limit

## INTENT
Ensure the player can only re-roll up to a maximum of exactly 3 times per account creation.

## THE RULE / LOGIC
Is reroll counter at or below limit before allowing another attempt?

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_character_reroll_limit]]`
