---
id: mech_skill_validation_grid_boundaries_verification
human_name: Grid Boundaries Verification
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
# Grid Boundaries Verification

## INTENT
Ensure that the target coordinate is within the mapped boundaries of the Grid.

## THE RULE / LOGIC
skill.target.outofgrid

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_validation_grid_boundaries_verification]]`
