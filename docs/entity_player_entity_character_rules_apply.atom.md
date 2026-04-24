---
id: entity_player_entity_character_rules_apply
human_name: Character Rules Apply to Players
type: ENTITY
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[entity_player]]
dependents: []
---
# Character Rules Apply to Players

## INTENT
Apply character rules to player entities upon creation.

## THE RULE / LOGIC
These characters must have their attributes rolled according to the rules defined in `entity_character`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_player_entity_character_rules_apply]]`
