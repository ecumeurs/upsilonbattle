---
id: mech_move_validation_move_validation_jump_limitations
human_name: Jump Limitations Rule
type: MECHANIC
layer: IMPLEMENTATION
version: 1.1
status: STABLE
priority: 5
tags: [movement, terrain, 3d, pathfinding]
parents:
  - [[mech_move_validation]]
dependents: []
---
# Jump Limitations Rule

## INTENT
The Z-axis difference between any two adjacent steps in a move path must not exceed the entity's `JumpHeight` property. This is the sole mechanism that turns terrain elevation into navigational cost.

## THE RULE / LOGIC
- **Property:** `JumpHeight` on the entity (see `property.JumpHeight`, int).
- **Default:** `2`. Absence of `JumpHeight` on an entity resolves to `2` via `entity.GetProperty` → `def.DefaultProperty` → `def.JumpHeight()`.
- **Randomized range:** `[2, 4)` via `entitygenerator.propertyRandomizers`. Generated characters (AI spawns, tests) land on 2 or 3.
- **API-started battles:** `upsilonapi/bridge/bridge.go` instantiates characters from caller-provided stats *without* explicitly setting `JumpHeight`; they therefore fall through to the default of `2`. This is sufficient because the Hill generator caps adjacent-cell Z delta at 2.
- **Validation step:** for each pair of consecutive cells `(a, b)` in the requested path, `|b.Z - a.Z| > JumpHeight` → path rejected with `entity.path.notvalid`.
- **Start step:** the entity's current position to the first path cell is checked with the same constraint.
- **A\* pathfinding:** `grid.AStarPath(start, end, jumpHeight, exclude)` prunes neighbours whose Z delta exceeds `jumpHeight`, so AI controllers cannot propose invalid paths.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_move_validation_move_validation_jump_limitations]]`
- **Property definition:** `upsilonbattle/battlearena/property/propertyenum.go` (`JumpHeight`), defaults in `property/def/entity.go`.
- **Randomization:** `upsilonbattle/battlearena/entity/entitygenerator/entitygenerator.go`.
- **Validation:** `upsilonbattle/battlearena/ruler/rules/move.go` (step-delta check) and `upsilonmapdata/grid/grid.go:AStarPath` (A\* pruning).

## EXPECTATION (For Testing)
- A 2-cell path where `|Δz| = JumpHeight` → accepted.
- A 2-cell path where `|Δz| = JumpHeight + 1` → rejected with `entity.path.notvalid`.
- `TestRuleMoveFailNotAdjascentJumpHeight` in `ruler/rules/rules_move_test.go` covers the rejection case.
- On a Hill-generated 10x10 map a character with the default `JumpHeight=2` can always reach every ground cell.
