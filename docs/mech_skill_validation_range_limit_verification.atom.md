---
id: mech_skill_validation_range_limit_verification
human_name: Range Limit Verification
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
# Range Limit Verification

## INTENT
Check if the mathematical distance to the target falls between the skill's MinRange and MaxRange.

## THE RULE / LOGIC
skill.target.range

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_validation_range_limit_verification]]`
