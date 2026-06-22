---
id: mech_move_validation
human_name: Entity Move Validation Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[shared:req_tech_debt_backlog]]
dependents:
  - [[upsiloncli:api_plan_travel_toward]]
---
# Entity Move Validation Mechanic

## INTENT
To validate an entity move request against all path and ownership preconditions before applying it to the game state, rejecting it with a specific error code if any check fails.

## THE RULE / LOGIC
**Entity Move Validation Mechanic.**

A move request must pass every check below before it is applied; the first failing check rejects the request with the indicated error code:

1. **Entity Existence** — the entity must exist within the game state (`entity.notfound`).
2. **Turn Mismatch** — the move command must target the currently active entity in the turn sequence (`entity.turn.mismatch`).
3. **Controller Mismatch** — the requesting Controller ID must own the entity (`entity.controller.mismatch`).
4. **Already Moved** — the entity must not have its `HasMoved` flag set (`entity.movement.already`).
5. **Path Adjacency** — each coordinate in the requested path array must be strictly adjacent to the previous one (`entity.path.notadjascent`, `entity.path.notvalid`).
6. **Path Length / Credits** — the path length must not exceed the entity's remaining Movement credits (`entity.path.too.long`, `entity.movement.nocredits`, `entity.movement.credits`).
7. **Obstacle Collision** — no node in the path may be an Obstacle (`entity.path.obstacle`).
8. **Entity Collision** — the final destination node must not be occupied by another entity (`entity.path.occupied`).
9. **Jump Limitations** — the Z-axis difference between consecutive path cells (and from the entity's current position to the first cell) must not exceed the entity's `JumpHeight`; this is the sole mechanism turning terrain elevation into navigational cost:
   - **Property:** `property.JumpHeight` (int). **Default 2** (absence resolves to 2 via `entity.GetProperty` → `def.DefaultProperty` → `def.JumpHeight()`).
   - **Randomized range:** `[2, 4)` via `entitygenerator.propertyRandomizers`; generated characters land on 2 or 3.
   - **API-started battles:** `upsilonapi/bridge/bridge.go` instantiates characters without explicitly setting `JumpHeight`, so they fall through to the default of 2; sufficient because the Hill generator caps adjacent-cell Z delta at 2. The bridge also projects incoming `X, Y` to `TopMostCellAt(x, y)` so clients may omit `Z`.
   - **Validation step:** for each consecutive pair `(a, b)`, `|b.Z - a.Z| > JumpHeight` → rejected with `entity.path.notvalid`.
   - **A* pathfinding:** `grid.AStarPath(start, end, jumpHeight, exclude)` prunes neighbours whose Z delta exceeds `jumpHeight`, so AI controllers cannot propose invalid paths.
   - **Walkable surfaces:** both `cell.Ground` and `cell.Dirt` are walkable/targetable.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/move.go`, `upsilonmapdata/grid/grid.go` (`AStarPath`), `upsilonbattle/battlearena/property/propertyenum.go`, `upsilonbattle/battlearena/property/def/entity.go`, `upsilonbattle/battlearena/entity/entitygenerator/entitygenerator.go`, `upsilonapi/bridge/bridge.go`
- **Test Names:** `rules_move_test.go`, `rules_move_extended_test.go`, `TestRuleMoveFailNotAdjascentJumpHeight`, `edge_movement_jump_limitations` (EC-09)
