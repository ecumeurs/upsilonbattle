---
id: mech_skill_validation
human_name: Skill Execution & Validation Mechanic
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[mech_skill_validation_action_state_verification]]
  - [[mech_skill_validation_economic_cost_verification_cooldown_check]]
  - [[mech_skill_validation_economic_cost_verification_stat_leech]]
  - [[mech_skill_validation_entity_targeting_rules_verification]]
  - [[mech_skill_validation_existence_verification]]
  - [[mech_skill_validation_grid_boundaries_verification]]
  - [[mech_skill_validation_range_limit_verification]]
  - [[mech_skill_validation_turn_controller_identity_verification]]
---
# Skill Execution & Validation Mechanic

## INTENT
To aggregate the constituent rules of Skill Execution & Validation Mechanic.

## THE RULE / LOGIC
Mechanic enforcing skills execution and validation constraints

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_validation]]`
