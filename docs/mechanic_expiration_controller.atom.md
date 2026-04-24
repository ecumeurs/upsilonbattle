---
id: mechanic_expiration_controller
status: DRAFT
priority: 5
version: 2.0
parent:
  - [[mech_entity_expiration]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
parents: []
---

# Expiration Controller

## INTENT

To manage the lifecycle of temporary entities, cleaning up expired effects and removing entities when their duration ends naturally.

## THE RULE / LOGIC
**Expiration Controller Mechanic:**

**Core Principle:**
To ensure system stability and prevent resource leaks, every temporary entity is governed by a centralized expiration logic. This system automatically manages the removal of entities that have completed their functional duration or reached their temporal limit.

**Centralized Expiration Logic:**
- **Turn-Based Decrement:** During the `End of Turn` resolution phase, the system audits all active entities.
- **Duration Tracking:** If an entity possesses a `Duration` property, the system decrements its value by 1.
- **Automatic Termination:** When the `Duration` counter reaches zero, the entity is marked for immediate removal.

**Removal and State Cleanup (The "Surgical Removal" Process):**
1. **Registry Purge:** The entity is removed from the primary Game State entity map and the turn order queue (Turner).
2. **Spatial Cleanup:** The entity's coordinates are cleared from the grid, ensuring the cell becomes traversable (if it was an obstacle).
3. **Dependency Cleanup:** The system performs a recursive scan for any `Positional Effects` (e.g., area-of-effect clouds) that were anchored to the expiring entity. If these effects are tagged to "Expire with Caster," they are purged simultaneously.
4. **Attribution Preservation:** Any final resolution events (e.g., a last tick of poison damage) are processed with credit attribution to the original caster before the entity's data is fully destroyed.

**Category-Specific Expiration Patterns:**
- **Persistent Zones:** Invisible anchor entities track the duration. When the anchor expires, the entire associated map zone (Poison Fog, Healing Well) is removed.
- **Trap Mechanics:** Traps expire either upon triggering (one-time use) or after a maximum turn count (timeout), whichever occurs first.
- **Channeling Structures:** These entities are unique in that they expire upon successful execution of their stored payload or upon being forcibly interrupted by an external action.

**System Integrity:**
- **Leak Prevention:** By centralizing expiration in the core ruleset, the system ensures that "orphaned" entities cannot remain in the game state indefinitely.
- **Performance:** Surgical removal minimizes the overhead of managing large numbers of temporary effects, keeping the grid efficient for pathfinding and collision checks.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[expiration_controller]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/endofturn.go`, `upsilonbattle/battlearena/ruler/rules/gamestate.go`
- **Integration:** Works with `mech_entity_expiration`, `mech_positional_effects`

## EXPECTATION
