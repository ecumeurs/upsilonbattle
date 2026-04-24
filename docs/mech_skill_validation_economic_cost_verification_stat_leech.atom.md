---
id: mech_skill_validation_economic_cost_verification_stat_leech
human_name: Economic Cost Verification - Stat Leech
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
# Economic Cost Verification - Stat Leech

## INTENT
Check if the entity possesses enough points in the respective stat to pay for the action.

## THE RULE / LOGIC
skill.cost.mp, skill.cost.sp, skill.cost.hp, skill.cost.mvt

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_validation_economic_cost_verification_stat_leech]]`
