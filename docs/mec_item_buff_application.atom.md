---
id: mec_item_buff_application
status: DRAFT
tags: [items, buffs, engine, iss-074]
parents:
  - [[entity_character_equipment]]
  - [[mec_equipment_stat_bonuses]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
human_name: Item Buff Application Mechanic
version: 2.0
---

# New Atom

## INTENT
To define how equipped items are projected onto the engine's `Entity` at battle init via the existing buff infrastructure (`Forever=true` `TemporaryProperties`), and how unequip is mirrored at runtime via origin-based removal.

## THE RULE / LOGIC
**Item-as-Buff Projection (battle init):**

For each entity in the incoming `ArenaStartRequest`, the bridge processes `entity.equipped_items[]` (populated by `UpsilonEntityResource` from `character_equipment` joined with `player_inventory.shop_item.properties`):

```go
for _, item := range entity.EquippedItems {
    buff := property.TemporaryProperties{
        Forever:        true,                     // never ticked down
        OriginEntityID: uuid.MustParse(item.ItemID),
        Properties:     map[string]property.Property{},
    }
    for key, raw := range item.Properties {
        buff.Properties[key] = def.ItemProperty(property.ItemProperties(key)).WithValue(raw)
    }
    entity.RegisterBuff(buff)

    // Weapon-as-skill: items can carry an Effect (skill ID).
    if effect, ok := item.Properties["Effect"]; ok && effect != nil {
        skill := registry.LookupSkill(effect)
        entity.RegisterSkill(skill)
    }
}
```

**Why `Forever=true`:**
- Equipment buffs must not decay with `BuffTickDown`. They persist for the entire match unless the item is unequipped.
- The existing `BuffTickDown` already exempts `Forever` buffs.

**Unequip at runtime (in-match equip changes are NOT supported in V2.0; this is for save-state restore and future runtime swap):**
```go
func (e *Entity) RemoveBuffsByOrigin(originID uuid.UUID) {
    kept := make([]property.TemporaryProperties, 0, len(e.Buffs))
    for _, b := range e.Buffs {
        if b.OriginEntityID != originID {
            kept = append(kept, b)
        }
    }
    e.Buffs = kept
}
```

**Property resolution:**
- `entity.GetProperty(name)` already iterates buffs and applies them via `prop.ApplyBuff(buff)` (see `entity.go:155-165`).
- Stacking across multiple buffs of the same property is additive by default; item-only properties (AttackRange, Shield) become first-class on the entity once any buff supplies them.

**Serialization contract:**
- `equipped_items` is part of the additive `ArenaStartRequest` extension. Old engine versions tolerate the unknown field per Go JSON unmarshal default. See `[[shared:api_standard_envelope]]` rules and the additive-contract guidance in CLAUDE.md.

**Skill placeholder:**
- The sibling `equipped_skills []string` field on the `ArenaStartRequest` Player.Entity payload is reserved for ISS-073. The bridge currently ignores it. Wiring is documented in ISS-073's preparation block.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mec_item_buff_application]]`
- **Files:** `upsilonapi/bridge/bridge.go` (insertion at the entity-bootstrap loop), `upsilonbattle/battlearena/entity/entity.go` (new `RemoveBuffsByOrigin`), `upsilonbattle/battlearena/property/buff.go` (existing `Forever` flag), `upsilonbattle/battlearena/property/def/item.go` (existing factories).
- **Test Names:** `TestArenaInit_EquippedItemsBecomeBuffs`, `TestRemoveBuffsByOrigin`, `TestItemEffectRegistersSkill`.

## EXPECTATION
- An entity with one armor + one sword equipped boots with two `Forever=true` buffs whose `OriginEntityID` matches the inventory item IDs.
- `entity.GetProperty("Armor").CValue` reflects the armor's ArmorRating contribution.
- `RemoveBuffsByOrigin(armor_id)` strips exactly the armor buff, leaving the sword buff intact.
- Items carrying an `Effect` register the corresponding skill on the entity.
