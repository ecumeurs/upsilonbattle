---
id: mechanic_mec_pass_facing_choice
status: DRAFT
type: MECHANIC
priority: 3
human_name: Pass Orientation Choice
layer: IMPLEMENTATION
version: 1.0
tags: turn,combat,defense
parents: []
dependents: []
---

# New Atom

## INTENT
To allow players to tactically choose their orientation when passing a turn, mitigating the risk of being backstabbed.

## THE RULE / LOGIC
**Pass Facing Choice Logic:**

1. **Trigger:** Player selects 'Pass' action.
2. **Interaction:** UI displays cardinal direction selector (Up, Right, Down, Left).
3. **Execution:**
    - If direction selected: Update entity `Orientation` to selected value, then trigger `EndTurn`.
    - If no selection (timeout or cancel): Maintain current `Orientation`, trigger `EndTurn`.
4. **Validation:** Selected orientation must be a valid `EntityOrientation` enum value.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mec_pass_facing_choice]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/pass.go`, `battleui` ActionPanel
- **API Endpoint:** `POST /v1/battle/pass` with `orientation` parameter.

## EXPECTATION
- Choose 'Pass' -> Prompt for orientation -> Select 'Left' -> Entity faces Left and turn ends.
- Choose 'Pass' -> Select 'None' -> Entity maintains current orientation and turn ends.
