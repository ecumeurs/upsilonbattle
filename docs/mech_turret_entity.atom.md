---
id: mech_turret_entity
status: DRAFT
priority: 5
version: 2.0
parents:
  - [[mech_behavior_system]]
  - [[mechanic_expiration_controller]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [entities, turrets, ai]
human_name: Turret Entity
---

# Turret Entity

## INTENT

To implement stationary entities that act on their own turn, using ranged attacks, and may have limited lifespans.

## THE RULE / LOGIC

**Turret Characteristics:**

- **EntityType**: TimeBased
- **Movement**: 0 (cannot move)
- **AttackRange**: > 0 (ranged attack)
- **Duration**: Optional (limited lifespan)
- **Controller**: Assigned controller (usually Aggressive)
- **Behavior**: AggressiveBehavior (finds nearest foe, attacks if in range)
- **WalkThrough**: true (other entities can move through)

**Turret Creation:**

When creating a turret:

1. Create entity with TimeBased type
2. Set Movement = 0, AttackRange = desired range
3. Set Duration if turret has limited lifespan
4. Set WalkThrough = true
5. Add attack skill(s) to entity
6. Add entity to Turner with 0 delay (immediate first turn)
7. Add entity to Grid at position
8. Assign controller with AggressiveBehavior

**Turret Behavior:**

Using AggressiveBehavior with Movement = 0:

1. Movement = 0, so skip pathfinding
2. Find nearest foe
3. Check if foe is in attack range
4. If in range: attack, add Delay, requeue
5. If not in range: pass turn (cannot move)
6. EndOfTurn: decrement Duration, remove if 0

**Turret Turn Flow:**

1. Turner gives turret a turn
2. AggressiveBehavior.OnTurn() called
3. Movement = 0, so skip movement logic
4. Find nearest foe in attack range
5. If in range: attack, add Delay, requeue
6. If not in range: pass turn
7. EndOfTurn: decrement Duration, remove if 0

**Composite Turret (Expiration + Aggressive):**

For turrets that need both expiration tracking and aggressive behavior:

1. Create CompositeBehavior with AND mode
2. Add ExpirationBehavior (checks duration, usually no decision)
3. Add AggressiveBehavior (finds and attacks targets)
4. Assign composite to turret's controller

**Turret Variants:**

- **Basic Turret**: Duration = 0 (infinite), AttackRange = 3
- **Sniper Turret**: Duration = 5, AttackRange = 8, Attack = 5
- **Rapid Fire Turret**: Duration = 3, AttackRange = 2, Attack = 2, fires multiple times
- **Healing Turret**: Uses healing skill instead of attack
- **Buffing Turret**: Applies buffs to allies in range

**Turret Credits:**

All damage dealt by turret awards credits to:
- The turret's controller (player who placed it)
- Tracked via effect.CasterID

**Turret vs Mobile Entity:**

Turrets differ from mobile entities:
- Cannot move (Movement = 0)
- Usually have limited Duration
- Often stationary strategic assets
- Can be targeted and destroyed

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_turret_entity]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entitygenerator/`, `upsilonbattle/battlearena/controller/behavior/aggressive.go`

## EXPECTATION
