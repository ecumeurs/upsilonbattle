---
id: mech_skill_validation_action_state_verification
human_name: Action State Verification
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_skill_validation]]
dependents: []
---
# Action State Verification

## INTENT
Prevent an entity from having already acted for a specific skill.

## THE RULE / LOGIC
entity.alreadyacted

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_validation_action_state_verification]]`
