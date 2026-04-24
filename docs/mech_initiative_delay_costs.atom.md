---
id: mech_initiative_delay_costs
human_name: Delay Costs Calculation
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[mech_initiative]]
dependents: []
---
# Delay Costs Calculation

## INTENT
Calculates cumulative numeric Delay Cost for actions during a turn.

## THE RULE / LOGIC
Actions performed during an active turn incur a cumulative numeric Delay Cost.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_initiative_delay_costs]]`
