---
id: mechanic_mech_temporary_entity_system
status: DRAFT
parent: []
dependents:
  - [[mec_cell_attached_effects]]
  - [[mec_channeling_mechanic]]
  - [[mec_effect_caster_tracking]]
  - [[mec_expiration_controller]]
human_name: Temporary Entity System
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
priority: 5
tags: [time-based, entities, effects]
---

# Temporary Entity System

## INTENT

To implement the temporary entity system where time-based mechanics (channeling, traps, area effects) are represented as entities with controlled lifespans and trigger behaviors.

## THE RULE / LOGIC

The Temporary Entity System provides unified infrastructure for all time-based game mechanics:

**Core Principle:** 1 skill effect = 1 entity (simplified approach)

**Entity Types** (v2.0 additions):
- **Character**: Player-controlled entities
- **Monster**: Enemy entities
- **TimeBased**: Channeling skills, delayed effects, turrets
- **Trap**: Triggers on OnStep, has Duration timeout
- **AreaEffect**: Affects zone each turn (anchor for positional effects)
- **Obstacle**: Walls/barriers (don't act, block movement)
- **Others**: Miscellaneous entities

**New Entity Properties:**
- **Duration**: How many turns to live
- **ExpiresWithCaster**: Remove when caster dies
- **WalkThrough**: Can walk through this entity?
- **Invisible**: Not visible to clients

**Trigger System:**
- **OnTurn**: Execute when entity's turn arrives (channeling, area effects)
- **OnStep**: Execute when entity stepped on (traps, quagmire)
- **OnDeath**: Execute when entity dies (explosions, cleanup)
- **OnEnter**: Execute when entity enters cell (poison fog)
- **OnExit**: Execute when entity leaves cell (remove debuffs)

**Zone Entity Pattern (v2.0):**

For effects spanning multiple cells (poisonous fog, healing zone):

- Create invisible anchor entity (TimeBased, has Duration)
- Set effect CasterID = anchor entity ID for all zone effects
- Attach PositionalEffects to cells in zone
- When anchor dies (Duration=0), RemoveEntity cleanup removes all zone effects

**Zone Entity Creation:**

When creating a zone effect:

1. Create invisible anchor entity with Duration and Invisible properties
2. Add anchor to Entities and Turner (for expiration tracking)
3. Calculate zone pattern (e.g., 3x3 cross)
4. For each cell in pattern:
   - Create effect with appropriate properties (poison, healing, etc.)
   - Set effect CasterID = anchor entity ID
   - Set ExpiresWithCaster = true
   - Attach effect to cell
5. Store effect in Effects map

**Expiration** (simplified v2.0):
- No ExpirationController needed
- Duration check in EndOfTurn for all entities
- When Duration = 0, RemoveEntity called
- RemoveEntity cleans up Turner, Grid, Entities, and PositionalEffects

**Effect Caster Tracking:**
- All effects remember caster via CasterID
- On entity death, effects with ExpiresWithCaster = true are removed
- Credits awarded to original caster for all damage/effects

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_temporary_entity_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/endofturn.go`, `upsilonbattle/battlearena/ruler/rules/gamestate.go`
- **Integration:** Works with `mech_positional_effects`, `mech_entity_expiration`

## EXPECTATION
