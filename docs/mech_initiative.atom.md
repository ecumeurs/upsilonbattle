---
id: mech_initiative
human_name: Initiative & Delay Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[module_backend_initiative_evaluation]]
dependents: []
---
# Initiative & Delay Mechanic

## INTENT
To determine turn order from action delay costs and a random starting roll, granting each character its active turn when its initiative ticker reaches zero.

## THE RULE / LOGIC
**Initiative & Delay Mechanic.**

Turn order is driven by per-character delay tickers managed by the Turner:

- **Initial Initiative Roll:** when the match begins, every character rolls an initial initiative value in the range **1 to 500**.
- **Active State:** a character receives its active turn only when its evaluated initiative ticker reaches **0**.
- **Delay Costs:** actions performed during an active turn incur a cumulative numeric Delay Cost (summed across the turn).
- **Requeue Calculation:** at the end of a turn, the character's delay until its next turn is computed from the summed Delay Cost and the character is re-queued accordingly.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_initiative]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/turner/turner.go`, `turner_state.go`, `upsilonbattle/battlearena/ruler/ruler_turn.go`
- **Test Names:** `turner_test.go`, `ruler_dead_entity_next_turn_test.go`
