---
id: mech_initiative
human_name: Initiative & Delay Mechanic
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[mech_initiative_active_state]]
  - [[mech_initiative_delay_costs]]
  - [[mech_initiative_initiative_roll]]
  - [[mech_initiative_requeue_calculation]]
---
# Initiative & Delay Mechanic

## INTENT
To aggregate the constituent rules of Initiative & Delay Mechanic.

## THE RULE / LOGIC
Determines turn order based on action weight and random starting values.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_initiative]]`
