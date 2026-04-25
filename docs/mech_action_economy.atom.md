---
id: mech_action_economy
human_name: Turn Action Economy Mechanic
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[mech_action_economy_action_cost_rules]]
  - [[mech_action_economy_time_constraint_rules]]
  - [[mech_action_economy_timeout_penalty_rules]]
---
# Turn Action Economy Mechanic

## INTENT
To aggregate the constituent rules of Turn Action Economy Mechanic.

## THE RULE / LOGIC
Defines the allowable actions and temporal constraints for a character's active turn.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_action_economy]]`
