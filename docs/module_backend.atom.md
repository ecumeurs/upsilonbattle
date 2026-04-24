---
id: module_backend
human_name: UpsilonBattle Backend Component
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[infra_mvp_docker]]
  - [[module_backend_action_economy]]
  - [[module_backend_board_generation]]
  - [[module_backend_combat_math]]
  - [[module_backend_initiative_evaluation]]
---
# UpsilonBattle Backend Component

## INTENT
To aggregate the constituent rules of UpsilonBattle Backend Component.

## THE RULE / LOGIC
UpsilonBattle Backend Component governs battle-related logic in an isolated Go application providing a JSON API.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[module_backend]]`
