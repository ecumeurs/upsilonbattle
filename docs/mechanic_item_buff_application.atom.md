---
id: mechanic_item_buff_application
status: DRAFT
tags: [items, buffs, engine, iss-074]
parents:
  - [[mechanic_equipment_stat_bonuses]]
  - [[upsilonapi:entity_character_equipment]]
dependents:
  - [[upsilontypes:mechanic_buff_attribution_accessor]]
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
human_name: Item Buff Application Mechanic
version: 2.0
---

# Item Buff Application Mechanic

## INTENT
To define how equipped items are projected onto the engine's `Entity` at battle init via the existing buff infrastructure (`Forever=true` `TemporaryProperties`), and how unequip is mirrored at runtime via origin-based removal.

## THE RULE / LOGIC
**Item-as-Buff Projection (battle init):**

For each entity in the incoming `ArenaStartRequest`, the bridge processes `entity.equipped_items[]` (populated by `UpsilonEntityResource` from `character_equipment` joined with `player_inventory.shop_item.properties`). For each equipped item, the bridge builds a `Forever=true` `property.TemporaryProperties` buff whose `OriginEntityID` is the item's ID, resolves each of the item's properties into a buff property, and registers the buff on the entity via `RegisterBuff`. Items that carry an `Effect` (a skill ID) also register the corresponding skill on the entity via `RegisterSkill`.

**Why `Forever=true`:**
- Equipment buffs must not decay with `BuffTickDown`. They persist for the entire match unless the item is unequipped.
- The existing `BuffTickDown` already exempts `Forever` buffs.

**Unequip at runtime (in-match equip changes are NOT supported in V2.0; this is for save-state restore and future runtime swap):**
`RemoveBuffsByOrigin(originID)` strips every buff on the entity whose `OriginEntityID` matches the given origin, keeping the rest.

**Property resolution:**
- `entity.GetProperty(name)` already iterates buffs and applies them via `prop.ApplyBuff(buff)` (see `entity.go:155-165`).
- Stacking across multiple buffs of the same property is additive by default; item-only properties (AttackRange, Shield) become first-class on the entity once any buff supplies them.

**Serialization contract:**
- `equipped_items` is part of the additive `ArenaStartRequest` extension. Old engine versions tolerate the unknown field per Go JSON unmarshal default. See `[[upsilonapi:api_standard_envelope]]` rules and the additive-contract guidance in CLAUDE.md.

**Skill placeholder:**
- The sibling `equipped_skills []string` field on the `ArenaStartRequest` Player.Entity payload is reserved for ISS-073. The bridge currently ignores it. Wiring is documented in ISS-073's preparation block.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_item_buff_application]]`
- **Files:** `upsilonapi/bridge/bridge.go` (insertion at the entity-bootstrap loop), `upsilontypes/entity/entity.go` (`RemoveBuffsByOrigin`; entity moved here from the now-removed `upsilonbattle/battlearena/entity/entity.go` package), `upsilonbattle/battlearena/property/buff.go` (existing `Forever` flag), `upsilonbattle/battlearena/property/def/item.go` (existing factories).
- **Test Names:** `TestArenaInit_EquippedItemsBecomeBuffs`, `TestRemoveBuffsByOrigin`, `TestItemEffectRegistersSkill`.
- **Companion atom:** `[[rule_entity_property_write_isolation]]` documents the write side of the same `Entity` property system: this atom's buff registration/composition (`RegisterBuff`, `GetProperty`) is the read side, and any code path that reads a buff-composed value and persists a new value back onto the entity must follow that rule's base-state isolation invariant instead of writing the composed value directly.

## EXPECTATION
- An entity with one armor + one sword equipped boots with two `Forever=true` buffs whose `OriginEntityID` matches the inventory item IDs.
- `entity.GetProperty("Armor").CValue` reflects the armor's ArmorRating contribution.
- `RemoveBuffsByOrigin(armor_id)` strips exactly the armor buff, leaving the sword buff intact.
- Items carrying an `Effect` register the corresponding skill on the entity.
