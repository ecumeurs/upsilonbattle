---
id: mech_skill_validation_entity_targeting_rules_verification
human_name: Entity Targeting Rules Verification
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
# Entity Targeting Rules Verification

## INTENT
Validate the Entity Targeting Rule based on the explicit TargetType attached to the skill.

## THE RULE / LOGIC
skill.target.none, skill.target.self, skill.target.enemyonly, skill.target.friendonly, skill.target.tile

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_validation_entity_targeting_rules_verification]]`
