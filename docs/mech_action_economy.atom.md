---
id: mech_action_economy
human_name: Turn Action Economy Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[shared:req_tech_debt_backlog]]
dependents:
  - [[upsilonbattleui:ui_action_panel]]
  - [[upsilonbattleui:ui_initiative_timeline]]
---
# Turn Action Economy Mechanic

## INTENT
To define the delay costs of each action and the temporal constraints (shot clock and timeout penalty) governing a character's active turn.

## THE RULE / LOGIC
**Turn Action Economy Mechanic.**

**Action Delay Costs:**
- **Move:** +20 delay cost per tile moved.
- **Attack:** +100 delay cost.
- **Pass:** +300 delay cost (base).

**Time Constraints:**
- Turn duration is strictly capped at **30 seconds**.
- **Enforcement:** the Go engine uses a `ShotClock` timer that fires an internal `Timeout` notification.
- **Race Prevention:** the `Timeout` handler validates that the Turn Index on the message matches the current Game State version before applying skip logic, so stale timeouts from previous turns are safely ignored (see `[[mech_game_state_versioning]]`).

**Timeout Penalty:**
- If a turn reaches the 30-second cap without completion, an automatic "Pass" is triggered and a strict penalty of **+100 delay cost** is added on top of the base Pass cost, for a total of **+400 delay**.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_action_economy]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/ruler_actions.go`, `ruler_turn.go`, `upsilonbattle/battlearena/ruler/rules/attack.go`, `move.go`, `endofturn.go`
