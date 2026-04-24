---
id: mech_initiative_active_state
human_name: Active State Logic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_initiative]]
dependents: []
---
# Active State Logic

## INTENT
Determines when a character receives their active turn.

## THE RULE / LOGIC
A character receives their active turn only when their evaluated initiative ticker reaches `0`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_initiative_active_state]]`
