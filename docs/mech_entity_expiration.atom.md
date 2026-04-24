---
id: mech_entity_expiration
status: DRAFT
priority: 5
version: 2.0
parent: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [entities, time-based, expiration]
human_name: Entity Expiration
---

# Entity Expiration

## INTENT

To enable entities with limited lifespans (turrets, delayed effects, walls) by tracking Duration and removing entities when expired.

## THE RULE / LOGIC

**Duration Property:**

Entities can have a Duration property that indicates how many turns they exist:
- Duration > 0: Entity lives for this many turns
- Duration = 0: Entity never expires (unless removed by other means)
- Duration decrements each EndOfTurn

**Expiration Check:**

In the EndOfTurn rule, after an entity passes turn:

1. Check if entity has Duration property
2. If Duration > 0, decrement by 1
3. If new Duration <= 0, remove entity

**RemoveEntity Function:**

Centralized cleanup function that:

1. Removes entity from Turner (no more turns)
2. Removes entity from Grid (frees cell)
3. Removes entity from Entities map
4. Cleans up PositionalEffects owned by this entity

**Cleanup of PositionalEffects:**

When RemoveEntity is called:

1. Iterate all PositionalEffects
2. For each effect at each position:
   - If effect CasterID matches dying entity and ExpiresWithCaster is true:
     - Remove effect from Effects map
   - Otherwise, keep effect
3. Update PositionalEffects to reflect removals

**Duration Initialization:**

Duration can be initialized in two ways:
- As a counter: current value and max value
- As a single value: decrements to 0

**Entity Types with Duration:**

- **TimeBased**: Delayed effects, channeling entities
- **Trap**: Traps that expire after N turns if not triggered
- **AreaEffect**: Zone effects that expire after N turns
- **Obstacle**: Walls/barriers that degrade over time
- **Turret**: Summons that last N turns

**Expiration Flow:**

1. Entity created with Duration = N
2. Each EndOfTurn, Duration decrements
3. When Duration = 0, RemoveEntity called
4. RemoveEntity:
   - Removes from Turner
   - Removes from Grid
   - Removes from Entities map
   - Cleans up PositionalEffects owned by this entity

**Edge Cases:**

- **Duration = 0**: Never expires (unless removed by other means)
- **Negative Duration**: Invalid, should not occur
- **Expiration during turn**: Entity completes current action, then removed at EndOfTurn

**Interaction with Turner:**

Entities that expire naturally:
- Are added to Turner with their initial delay
- Get turns until Duration reaches 0
- When Duration = 0, removed from Turner before next turn

Entities without controllers (like walls):
- May not be added to Turner at all
- Duration checked in EndOfTurn for all entities
- Removed when Duration = 0

**Zone Entity Expiration:**

For zone effects created by an anchor entity:
- Anchor entity has Duration
- All zone effects have CasterID = anchor entity ID
- All zone effects have ExpiresWithCaster = true
- When anchor expires (Duration = 0), RemoveEntity cleanup removes all zone effects

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_entity_expiration]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/endofturn.go`, `upsilonbattle/battlearena/ruler/rules/gamestate.go`

## EXPECTATION
