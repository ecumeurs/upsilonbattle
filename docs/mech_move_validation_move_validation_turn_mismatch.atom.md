---
id: mech_move_validation_move_validation_turn_mismatch
human_name: Turn Mismatch Rule
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_move_validation]]
dependents: []
---
# Turn Mismatch Rule

## INTENT
A move command must match the currently active entity in the turn sequence.

## THE RULE / LOGIC
entity.turn.missmatch

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation_move_validation_turn_mismatch]]`
