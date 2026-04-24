---
id: mec_cell_attached_effects
human_name: Cell Attached Effects Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [grid, effects, movement]
parents:
  - [[mechanic_mech_temporary_entity_system]]
dependents: []
---

# Cell Attached Effects Mechanic

## INTENT
To implement effects attached to grid cells that trigger on movement events (walking in/out) and modify movement costs, supporting environmental hazards like quagmire without creating multiple entities.

## THE RULE / LOGIC
**Cell Attached Effects Mechanic:**

**Core Principle:**
Enables the attachment of persistent effects to specific grid cells that interact with entity movement and turn-based logic without requiring individual entities for every affected tile.

**Master Entity Architecture:**
- **Centralized Control:** A single "Master Entity" governs an entire area effect (e.g., a 3x3 fog cloud).
- **Lifecycle Management:** The Master Entity tracks its own duration, typically by losing 1 HP per game turn until expiration.
- **Area Definition:** Uses standard skill zone patterns (Cross, Circle, Square) to map which grid cells are currently affected.

**Movement Event Triggers:**
- **On-Entry (OnStepIn):** Triggers an effect or modifier the moment an entity enters a tagged cell.
- **On-Exit (OnStepOut):** Triggers when an entity moves out of a tagged cell.
- **Bidirectional (OnStepBoth):** Triggers on both entry and exit events for maximum interaction.

**Modifier Rules:**
- **Movement Costs:** Affected cells can modify the base movement cost (e.g., Quagmire increasing cost from 1 to 2).
- **Pathfinding Integration:** The movement cost modifications must be exposed to the pathfinding algorithm to ensure AI and player movement previews accurately reflect the cost of traversal.
- **Persistence:** These modifiers are applied dynamically during the movement phase as the entity traverses the grid.

**Interaction Logic:**
- **Periodic Effects (OnTurn):** During the start or end of a turn, the Master Entity applies its core effect (e.g., Poison, Healing) to all entities currently occupying cells within its defined zone.
- **Granularity Choice:**
    - Use individual entities per cell when each tile requires unique state or independent timing.
    - Use a Master Entity when a uniform effect covers multiple tiles and shares a single expiration timer.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_cell_attached_effects]]`
- **Related Files:** `upsilonmapdata/grid/grid.go`, `upsilonbattle/battlearena/ruler/rules/move.go`
