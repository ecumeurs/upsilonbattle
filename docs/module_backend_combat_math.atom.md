---
id: module_backend_combat_math
human_name: Combat Math Logic
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[module_backend]]
dependents:
  - [[mech_combat_attack_computation]]
  - [[mech_combat_standard_attack_computation]]
---
# Combat Math Logic

## INTENT
Resolve combat-related HP calculations

## THE RULE / LOGIC
Must expose operations via a JSON API.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[module_backend_combat_math]]`
