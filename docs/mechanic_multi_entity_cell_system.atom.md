---
id: mechanic_multi_entity_cell_system
status: DRAFT
version: 2.0
parent: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
tags: [grid, entities, multi-entity]
parents: []
---

# Multi-Entity Cell System

## INTENT

To implement multi-entity cell system where multiple entities (characters, effects, traps, obstacles) can occupy the same grid cell simultaneously, enabling complex tactical interactions.

## THE RULE / LOGIC
**Multi-Entity Cell System Mechanic:**

**Core Principle:**
This mechanic fundamentally shifts the grid from a "one-entity-per-cell" model to a "multi-occupant" model. This allows for complex tactical overlays where characters, environmental effects, traps, and obstacles can all occupy the same spatial coordinate while maintaining distinct collision and interaction rules.

**Cell Architecture (V2):**
- **Entity Identification:** Each grid cell maintains an `EntityIDs` collection (an array of unique identifiers) rather than a single occupant reference.
- **Effect Identification:** A separate `EffectIDs` collection tracks positional status effects anchored to the cell.
- **Terrain Dynamics:** Cells possess a `CrossingCost` property, allowing for varied terrain types (e.g., roads, rubble, mud) to influence movement efficiency independently of any active effects.

**Collision and Traversal Rules:**
- **The WalkThrough Constraint:** Traversal through a cell is governed by the `WalkThrough` property of every entity present within its `EntityIDs` list.
    - **Blocking:** If any entity in the cell has `WalkThrough = False` (e.g., another Character or a solid Barrier), the cell is considered blocked for movement.
    - **Permissive:** If all entities in the cell have `WalkThrough = True` (e.g., Traps or non-corporeal spirits), the cell remains traversable.
- **Pathfinding Integration (A*):** The navigation algorithm considers a cell "impassable" if any of its current occupants are marked as non-traversable.

**Cumulative Movement Cost Calculation:**
The total cost to enter a cell is the sum of:
1. **Base Step Cost:** The fundamental cost of one step (typically 1).
2. **Terrain Crossing Cost:** The additional cost defined by the cell's `CrossingCost` property.
3. **Active Effect Penalties:** The sum of all `MvtCost` properties from effects currently present in the cell's `EffectIDs` list.
*Note: Movement resolution allows characters to enter a cell even if the cost exceeds their remaining movement pool, effectively ending their movement at that location.*

**Grid Lifecycle Management:**
- **Atomic Movement:** The `MoveEntity` operation surgically removes an ID from its source cell's list and appends it to the destination cell's list, ensuring the total entity count in the game state remains consistent.
- **Occupancy Retrieval:** The system provides specialized lookup methods to return the full collection of entity IDs at any given coordinate, enabling area-of-effect calculations and proximity checks.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[multi_entity_cell_system]]`
- **Related Files:** `upsilonmapdata/grid/cell/cell.go`, `upsilonmapdata/grid/grid.go`, `upsilonbattle/battlearena/ruler/rules/move.go`

## EXPECTATION
