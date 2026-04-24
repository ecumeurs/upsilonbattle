---
id: mec_multi_entity_cell_system
human_name: Multi-Entity Cell System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [grid, entities, movement]
parents: []
dependents: []
---

# Multi-Entity Cell System Mechanic

## INTENT
To enable multiple entities to occupy the same grid cell, supporting character + effects co-location while maintaining character-vs-character collision rules.

## THE RULE / LOGIC
**Multi-Entity Cell System Mechanic:**

**Core Principle:**
Extends the grid architecture to allow multiple entities of varying types to occupy the same spatial coordinate simultaneously, facilitating complex interactions between characters and environmental effects while maintaining strict collision boundaries.

**Occupancy Constraints:**
- **Character Exclusive Slot:** Each grid cell supports a maximum of **one** Character entity. Attempts to move a character into a cell already occupied by another character will result in a collision failure.
- **Effect/Temporary Stack:** Each grid cell can support **multiple** Effect or Temporary entities (e.g., a trap, a fog cloud, and a channeling marker can all exist in the same cell).
- **Hybrid Occupation:** A cell may simultaneously contain one Character and multiple Temporary entities.

**Collision and Traversal Logic:**
- **The WalkThrough Property:** Every entity possesses a boolean `WalkThrough` attribute that dictates its physical presence on the grid.
    - **WalkThrough = True:** The entity allows other entities to pass through or occupy its cell (e.g., beneficial clouds, hidden traps, channeling indicators).
    - **WalkThrough = False:** The entity acts as a solid obstacle, blocking all movement into that cell (e.g., Characters, summoned barriers, static walls).
- **Movement Validation Hierarchy:**
    1. **Character Check:** If the cell contains a character other than the mover, movement is denied.
    2. **Obstacle Check:** The system iterates through all non-character entities in the cell. If any entity has `WalkThrough = False`, movement is denied.
    3. **Resolution:** If neither check fails, movement into the cell is permitted.

**Entity-Specific Defaults:**
- **Characters:** Inherently `WalkThrough = False`. They always occupy the exclusive character slot.
- **Traps:** Typically `WalkThrough = True`. They must allow entities to enter the cell to trigger the "Step-In" event.
- **Environmental Hazards:** Usually `WalkThrough = True` (e.g., Poison Gas) but can be `False` (e.g., Ice Wall) depending on the desired mechanical impact.
- **Markers/Indicators:** Always `WalkThrough = True` as they are purely informational or represent pending actions.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_multi_entity_cell_system]]`
- **Related Files:** `upsilonmapdata/grid/grid.go`, `upsilonbattle/battlearena/ruler/rules/move.go`
