---
id: mech_constructed_obstacle
status: DRAFT
priority: 5
version: 2.0
parent:
  - [[mech_entity_expiration]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [entities, obstacles, walls]
human_name: Constructed Obstacle

## INTENT

To implement player-created obstacles (walls, barriers) that block movement, have HP, and may degrade over time.

## THE RULE / LOGIC

**Obstacle Characteristics:**

- **EntityType**: Obstacle
- **Movement**: 0 (cannot move)
- **WalkThrough**: false (blocks movement)
- **HP**: Has health, can be destroyed
- **Duration**: Optional (for decay over time)
- **Controller**: uuid.Nil (no controller) OR ExpirationBehavior only
- **Turner**: May or may not be in Turner (no controller = no turns needed)

**Obstacle Creation:**

When creating an obstacle:

1. Create entity with Obstacle type
2. Set WalkThrough = false (blocks movement)
3. Set HP to desired value
4. Set Duration if obstacle should decay
5. Add entity to Grid at position
6. Optionally add to Turner if it has a controller

**Movement Blocking:**

When validating movement:

1. For each cell in movement path:
   - Check all entities in cell
   - Skip self
   - If any entity has WalkThrough = false:
     - Movement blocked, return error

**Pathfinding Update:**

Pathfinding (A*) must exclude cells with blocking entities:

1. When generating path, check each candidate cell
2. Get all entities in cell
3. If any entity has WalkThrough = false:
   - Exclude cell from path

**Obstacle Damage:**

Obstacles can be targeted by attacks:

1. Check if target is an obstacle entity
2. Calculate damage as normal
3. Apply damage to obstacle's HP
4. If HP <= 0:
   - Remove obstacle from game
5. Credits awarded to attacker

**Obstacle Decay:**

Using Duration property:

1. In EndOfTurn, check all entities
2. If entity has Duration > 0:
   - Decrement Duration
   - If Duration = 0:
     - Remove entity

**Obstacle with Controller:**

For obstacles that need behavior (e.g., force field that pulses):

1. Create obstacle with controller
2. Assign ExpirationBehavior to handle duration
3. Add to Turner for expiration tracking
4. Obstacle gets turns but behavior handles them

**Obstacle Variants:**

- **Stone Wall**: HP = 50, Duration = 0 (permanent until destroyed)
- **Ice Wall**: HP = 20, Duration = 10, WalkThrough = false
- **Force Field**: HP = 30, Duration = 5, WalkThrough = true, pulses damage
- **Spike Trap**: HP = 10, Duration = 20, damages on entry
- **Barrier**: HP = 40, WalkThrough = false, blocks LoS

**Turner Consideration:**

Obstacles without controllers don't need Turner entries.
EndOfTurn iterates over all entities to check Duration, so obstacles will expire regardless.

**Obstacle vs Turret:**

Obstacles differ from turrets:
- Obstacles don't act (no attacks or skills)
- Obstacles block movement (WalkThrough = false)
- Turrets act and attack (have controller and behavior)
- Turrets don't block (WalkThrough = true)

**Zone Control:**

Obstacles enable tactical zone control:
- Block chokepoints
- Create safe zones
- Force specific paths
- Provide cover

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_constructed_obstacle]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/move.go`, `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonmapdata/grid/grid.go`

## EXPECTATION
