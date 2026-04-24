---
id: mech_action_economy_time_constraint_rules
human_name: Time Constraint Rules Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.1
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_action_economy]]
dependents: []
---
# Time Constraint Rules Mechanic

## INTENT
Defines the temporal constraints for a turn.

## THE RULE / LOGIC
- Turn duration is strictly capped at 30 seconds.
- **Enforcement:** The Go Engine uses a timer (`ShotClock`) that triggers an internal `Timeout` notification.
- **Race Prevention:** The `Timeout` handler validates that the Turn Index on the message matches the current Game State version before applying the skip logic. This ensures stale timeouts (from previous turns) are safely ignored. [[mech_game_state_versioning]]

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_action_economy_time_constraint_rules]]`
