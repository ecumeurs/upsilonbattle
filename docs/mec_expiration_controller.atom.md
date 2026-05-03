---
id: mec_expiration_controller
human_name: Expiration Controller Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [controllers, entities, time-based]
parents:
  - [[upsilontypes:mechanic_mech_temporary_entity_system]]
dependents: []
---

# Expiration Controller Mechanic

## INTENT
To implement the ExpirationController that manages the lifecycle of temporary entities, handling their death and cleanup when their duration expires or their effect completes.

## THE RULE / LOGIC
**Expiration Controller Mechanic:**

**Core Principle:**
The Expiration Controller is a specialized automation layer that manages the lifecycle of temporary entities, ensuring they execute their intended effects and are cleanly removed from the game state upon completion or duration expiration.

**Controller Execution Patterns:**
- **Standard Self-Termination:** Upon the arrival of the temporary entity's turn, the controller executes the associated skill/effect and immediately initiates the entity's death sequence.
- **Duration-Based (Multi-Turn):**
    - The entity possesses a hidden "Duration" counter (often represented as HP).
    - Each turn, the controller executes the effect (e.g., Poisonous Fog) and decrements the counter by 1.
    - The entity is marked for deletion only when the counter reaches zero.
- **One-Time Channeling:**
    - The controller manages an entity that represents a pending skill (channeling).
    - Once the delay condition is met, the skill is fired, the casting character is released from their "Casting" state, and the temporary entity is destroyed.
- **Trigger-Based (Traps):**
    - The controller remains dormant during turn cycles.
    - Termination is triggered externally (e.g., by a "Step-In" movement event).
    - After the trap triggers and resolves its payload, the controller ensures the trap entity is removed from the grid.

**Termination and Cleanup Sequence:**
1. **Notification:** The controller sends an "End of Turn" or "Death" message to the system ruler, specifying the target Entity ID.
2. **On-Death Resolution:** The system executes any final logic or visual effects associated with the entity's removal.
3. **Registry Removal:** The entity is surgically removed from the grid coordinates, the turn order (Turner), and the primary game state entity map.
4. **Caster Synchronization:** If the temporary entity was linked to a character (e.g., a channeled spell), the system updates the character's properties to reflect that they are no longer occupied by that specific channeling task.

**System Integration:**
- **Reliability:** By decoupling lifecycle management from the primary player/AI controllers, the Expiration Controller ensures that temporary entities are never "orphaned" in the game state, preventing memory leaks and grid clutter.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_expiration_controller]]`
- **Related Files:** `upsilonbattle/battlearena/controller/controller.go`, `upsilonbattle/battlearena/ruler/rules/endofturn.go`
