---
id: mechanic_expiration_controller
human_name: "Expiration Controller"
status: DRAFT
priority: 5
version: 2.0
parents:
  - [[upsilontypes:mechanic_temporary_entity_system]]
dependents:
  - [[mech_constructed_obstacle]]
  - [[mech_turret_entity]]
type: MECHANIC
layer: IMPLEMENTATION
---

# Expiration Controller

## INTENT

To manage the lifecycle of temporary entities, cleaning up expired effects and removing entities when their duration ends naturally.

## THE RULE / LOGIC
**Expiration Controller Mechanic:**

**Core Principle:**
To ensure system stability and prevent resource leaks, every temporary entity is governed by a centralized expiration logic. This automation layer manages removal of entities that have completed their functional duration, executed their payload, or reached their temporal limit.

**Duration Property:**
- Duration > 0: entity lives for that many turns.
- Duration = 0: entity never expires (unless removed by other means). Negative Duration is invalid and should not occur.
- Duration may be initialized as a counter (current + max value) or as a single value that decrements to 0; it is sometimes represented as HP.

**Turn-Based Expiration Check (EndOfTurn rule):**
1. After an entity passes its turn, check whether it has a `Duration` property.
2. If Duration > 0, decrement by 1.
3. If the new Duration ≤ 0, call `RemoveEntity`.
*An entity completes its current action during the turn, then is removed at EndOfTurn.*

**Controller Execution Patterns:**
- **Standard Self-Termination:** when the temporary entity's turn arrives, the controller executes the associated skill/effect and immediately initiates the death sequence.
- **Duration-Based (Multi-Turn):** each turn the controller executes the effect (e.g. Poisonous Fog) and decrements the hidden Duration counter; deletion occurs when it reaches zero.
- **One-Time Channeling:** the controller manages a pending-skill entity; once the delay condition is met the skill fires, the casting character is released from its "Casting" state, and the entity is destroyed (or destroyed on forced interruption).
- **Trigger-Based (Traps):** the controller stays dormant; termination is triggered externally (e.g. a "Step-In" movement event) and after the payload resolves the trap entity is removed. A trap expires on trigger (one-time) or after a max turn count (timeout), whichever first.

**RemoveEntity — Surgical Removal & State Cleanup:**
1. **Notification:** the controller sends an "End of Turn"/"Death" message to the ruler with the target Entity ID.
2. **On-Death Resolution:** final logic/visual effects (e.g. a last tick of poison) are processed with credit attribution to the original caster before data is destroyed.
3. **Registry Purge:** the entity is removed from the Turner (no more turns), the Grid (frees the cell, making it traversable if it was an obstacle), and the primary Entities map.
4. **Caster Synchronization:** if the entity was linked to a character (e.g. a channeled spell), the character's properties are updated so it is no longer occupied by that task.
5. **Dependency / PositionalEffects Cleanup:** iterate all `PositionalEffects`; for each effect whose `CasterID` matches the dying entity AND `ExpiresWithCaster = true`, remove it from the Effects map; otherwise keep it.

**Entity Types with Duration:** TimeBased (delayed/channeling entities), Trap, AreaEffect/zone, Obstacle (walls/barriers), Turret (summons).

**Category-Specific Patterns:**
- **Persistent / Zone Effects:** an invisible anchor entity tracks duration; all zone effects carry `CasterID = anchor ID` and `ExpiresWithCaster = true`, so when the anchor expires `RemoveEntity` purges the entire zone (Poison Fog, Healing Well).
- **Controllerless entities (e.g. walls):** may not be added to the Turner at all; their Duration is still checked in EndOfTurn for all entities and they are removed when Duration = 0.

**System Integrity:**
- **Leak Prevention:** centralizing expiration in the core ruleset guarantees no "orphaned" entities remain in the game state.
- **Performance:** surgical removal keeps the grid efficient for pathfinding and collision checks.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_expiration_controller]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/endofturn.go`, `upsilonbattle/battlearena/ruler/gamestate/gamestate_logic.go`, `upsilonbattle/battlearena/ruler/behavior/behavior.go`, `upsilonbattle/battlearena/controller/controller.go`
- **Test Names:** `upsilonbattle/battlearena/ruler/rules/rules_iss066_test.go`
- **Integration:** Works with `[[upsilontypes:mech_positional_effects]]`

## EXPECTATION
- An entity with Duration > 0 has its Duration decremented by 1 at each EndOfTurn and is removed when Duration reaches ≤ 0.
- An entity with Duration = 0 never auto-expires.
- `RemoveEntity` removes the entity from the Turner, the Grid (freeing its cell), and the Entities map in a single pass.
- On removal, PositionalEffects whose `CasterID` matches the dying entity and whose `ExpiresWithCaster = true` are purged; all other effects are retained.
- A zone anchor's expiration removes every zone effect tagged with its `CasterID`.
- A channeling/TimeBased entity self-terminates after firing its payload and releases the linked caster's Casting state.
- A trap entity expires on trigger (Step-In) or after its timeout turn count, whichever comes first.
- Final on-death resolution (e.g. last poison tick) is attributed to the original caster before the entity is destroyed.
- Verified by `upsilonbattle/battlearena/ruler/rules/rules_iss066_test.go`.
