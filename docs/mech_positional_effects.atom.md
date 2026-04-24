---
id: mech_positional_effects
status: DRAFT
priority: 5
version: 2.0
parent: []
dependents:
  - [[mechanic_cell_attached_effects]]
type: MECHANIC
layer: IMPLEMENTATION
tags: [effects, grid, positions]
human_name: Positional Effects
---

# Positional Effects

## INTENT

To attach effects to grid positions instead of entities, enabling zone effects, traps, and terrain modifiers that persist independent of any character.

## THE RULE / LOGIC

**Data Structure:**

Positional effects are stored in three places:

- **Cell.EffectIDs**: Array of effect IDs attached to this cell
- **GameState.PositionalEffects**: Map from position to effect IDs for quick lookup
- **GameState.Effects**: Central storage of all effect data by ID

**Effect Ownership:**

Each effect has a CasterID:
- For zone effects: CasterID is the anchor entity ID
- For traps: CasterID is the trap entity ID
- For player-created effects: CasterID is the player entity ID

**Effect Properties for Positional Effects:**

New properties for positional effects:
- TriggerType: When the effect fires
- RemoveOnTrigger: Whether to remove effect after firing
- TriggerCount: How many times effect can trigger
- ExpiresWithCaster: Whether to remove effect when caster dies

**Creating a Positional Effect:**

When adding a positional effect:

1. Set effect's CasterID
2. Generate unique effect ID
3. Store effect in GameState.Effects map
4. Add effect ID to cell's EffectIDs array
5. Update GameState.PositionalEffects map

**Removing a Positional Effect:**

When removing a positional effect:

1. Remove effect from GameState.Effects map
2. Remove effect ID from cell's EffectIDs array
3. Update GameState.PositionalEffects map
4. If cell has no more effects, remove position from map

**Cleanup on Caster Death:**

When an entity dies, all effects with ExpiresWithCaster = true and matching CasterID are removed:

1. Iterate all PositionalEffects
2. For each effect at each position:
   - If effect CasterID matches dying entity and ExpiresWithCaster is true:
     - Remove effect from Effects map
   - Otherwise, keep effect
3. Update PositionalEffects to reflect removals

**Zone Entity Pattern:**

For effects that span multiple cells (poisonous fog, healing zone):

1. Create invisible anchor entity with Duration
2. Set effect CasterID = anchor entity ID
3. Set effect ExpiresWithCaster = true
4. Attach effects to cells in zone pattern
5. When anchor entity dies (Duration = 0), RemoveEntity cleanup removes all zone effects

**Anchor Entity Characteristics:**

- Type: TimeBased
- ControllerID: Nil (no controller) or ExpirationBehavior only
- Position: Center of zone
- Properties: Duration, Invisible = true
- Does not block movement (WalkThrough = true)
- Has no collision

**Zone Patterns:**

Effects can be attached to multiple cells using patterns:
- Single cell: Effect affects 1x1 area
- Radius: Circular area of N cells radius
- Pattern: Specific shape (line, cone, cross, etc.)

**Effect Lifecycle:**

1. **Creation**: Effect created and attached to cells
2. **Active**: Effect triggers based on trigger type
3. **Expiration**: Effect removed when:
   - Trigger count reaches 0
   - RemoveOnTrigger is true and effect fires
   - Caster dies and ExpiresWithCaster is true
   - Anchor entity dies (for zone effects)

**Credit Assignment:**

Credits for effects are awarded to the original caster (CasterID), not the entity that happens to be at the position when the effect triggers.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_positional_effects]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/gamestate.go`, `upsilonmapdata/grid/cell/cell.go`
- **Integration:** Works with `mech_trigger_system` for effect execution

## EXPECTATION
