---
id: mechanic_cell_attached_effects
status: DRAFT
priority: 5
version: 2.0
parents:
  - [[mech_positional_effects]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [effects, grid, movement]
human_name: Cell Attached Effects
---

# Cell Attached Effects

## INTENT

To implement effects attached to grid cells that trigger on movement events (walking in/out/through) and modify movement costs, supporting environmental hazards like quagmire, poison clouds, and healing zones.

## THE RULE / LOGIC
**Cell Attached Effects Mechanic:**

**Core Principle:**
Enables the association of persistent skill effects with specific grid coordinates. These effects interact with entities through movement-based triggers and turn-based cycles, modifying traversal costs and applying status changes.

**Data Architecture:**
- **Grid Integration:** Each grid cell maintains a list of associated Effect IDs.
- **State Mapping:** The system tracks positional effects through a global mapping in the game state, allowing for efficient lookup and removal.
- **Relational Integrity:** Effects are stored as independent objects with their own properties (e.g., CasterID, Duration) and are linked to one or more cells.

**Operational Triggers:**
- **On-Entry (OnStepIn):** Resolves immediately when an entity moves into an affected cell.
- **On-Exit (OnStepOut):** Triggers when an entity vacates a cell, often used for removing temporary modifiers or resolving "hazard" damage.
- **Traversal (OnStep):** Applied for every step taken through an affected zone, including the final destination.
- **Periodic (OnTurn):** Triggers at the start of an entity's turn if they occupy a cell containing the effect.

**Movement Cost and Constraints:**
- **Variable Cost Modifiers:** Effects can dynamically increase the movement cost of a cell (e.g., difficult terrain adding +2 cost per step).
- **Cost Resolution:** During path execution, the system sums the base cell crossing cost and all active effect modifiers along the path. Movement is deducted from the character's available movement pool.
- **Flow Control:** Specific effects can include "Force Stop" or "Force End Turn" properties, which immediately terminate a character's movement or turn upon triggering.

**Effect Resolution Lifecycle:**
1. **Detection:** During the movement phase, the system identifies all cells in the planned path that contain effect IDs.
2. **Trigger Evaluation:** For each cell, the system matches the movement event (Enter/Pass/Exit) against the effect's defined trigger type.
3. **Application:** If a match occurs, the effect's payload (damage, heal, buff, debuff) is applied to the traversing entity.
4. **Cleanup:** Effects tagged as "Remove on Trigger" (e.g., a one-time trap) are purged from the grid cell after resolution.
5. **Turn-Based Sync:** At the beginning of each turn, the system re-validates the entity's current position and applies any relevant "On-Turn" effects before allowing action selection.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[cell_attached_effects]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/move.go`, `upsilonbattle/battlearena/ruler/rules/beginingofturn.go`
- **Integration:** Works with `mech_positional_effects`, `mech_trigger_system`

## EXPECTATION
