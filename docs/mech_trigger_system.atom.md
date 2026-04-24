---
id: mech_trigger_system
status: DRAFT
priority: 5
version: 2.0
parent:
  - [[mech_positional_effects]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [effects, triggers, movement]
human_name: Trigger System
---

# Trigger System

## INTENT

To define when positional effects fire, enabling effects that trigger on movement events, turn events, and death events.

## THE RULE / LOGIC

**Trigger Types:**

- `OnEnter`: When entity enters cell
- `OnExit`: When entity leaves cell
- `OnStep`: Each step through cell
- `OnTurn`: Each turn while entity is in cell
- `OnDeath`: When entity dies in/on this cell

**Trigger Behavior:**

**OnEnter:**
- Fires exactly once when entity first enters the cell
- Does not fire again while entity remains in cell
- Use case: Apply poison when stepping into poisonous fog

**OnExit:**
- Fires when entity leaves the cell
- Only fires if entity was previously in the cell
- Use case: Remove movement debuff when leaving quagmire

**OnStep:**
- Fires for each cell the entity moves through
- Fires for both entry and subsequent steps within same turn
- Use case: Traps trigger when stepped on

**OnTurn:**
- Fires at the beginning of entity's turn if entity is in the cell
- Fires every turn while entity remains in cell
- Use case: Poison damage applied each turn while standing in poison zone

**OnDeath:**
- Fires when entity dies while in or on this cell
- Use case: Explosion triggers when burning entity dies

**Movement Trigger Execution:**

When an entity moves through a path:

1. For each cell in the path:
   - Get all effects on the cell
   - For each effect:
     - Determine if trigger fires based on position in path
     - If final position: OnEnter and OnStep triggers fire
     - If passing through: OnStep triggers fire
     - Apply effect to entity
     - Handle removal/trigger count

**Turn Trigger Execution:**

At the beginning of each entity's turn:

1. Get cell at entity's current position
2. For each effect on the cell:
   - If trigger type is OnTurn, apply effect
   - Handle removal/trigger count

**Trigger Count:**

Effects can have a trigger count:
- `0`: Unlimited triggers
- `> 0`: Can trigger this many times before being removed

When a trigger fires:
- If RemoveOnTrigger is true, remove effect immediately
- Otherwise, decrement trigger count
- If trigger count reaches 0, remove effect

**Trigger Stacking:**

Multiple effects can trigger on the same cell:
- Order: First Come First Served (based on when effects were added)
- Effects don't prevent each other from triggering
- All valid effects fire

**Effect Properties for Triggers:**

Effects need these properties:
- TriggerType: Which trigger fires this effect
- RemoveOnTrigger: Should effect be removed after firing?
- TriggerCount: How many times can this trigger?

**Movement Cost Interaction:**

Triggers can modify movement cost:
- Effects can have MvtCost property
- Costs are summed from all triggering effects
- Costs are applied during move execution
- Insufficient movement doesn't prevent entry, just consumes all remaining movement

**Force Stop Triggers:**

Effects can force entity to stop acting:
- ForceStopMove: Sets HasMoved flag, prevents further movement
- ForceEndTurn: Sets HasActed flag, ends turn immediately

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_trigger_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/move.go`, `upsilonbattle/battlearena/ruler/rules/beginingofturn.go`
- **Integration:** Works with `mech_positional_effects`

## EXPECTATION
