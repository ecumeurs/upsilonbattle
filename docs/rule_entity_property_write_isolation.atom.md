---
id: rule_entity_property_write_isolation
status: REVIEW
tags: [entity, properties, buffs, write-back, iss-144, iss-145]
layer: ARCHITECTURE
priority: 5
version: 1.0
human_name: Entity Property Write Isolation Rule
dependents: []
type: RULE
parents:
  - [[shared:rule_stat_taxonomy]]
---

# Entity Property Write Isolation Rule

## INTENT
To require that any write deriving a new value from Entity.GetProperty's buff-composed read result strips the buff contribution before that value is persisted, keeping an entity's base property state and its buff-composed read value permanently distinct.

## THE RULE / LOGIC
**Base vs. composed property state:**
- Base state is the value stored directly in `Entity.Properties`, independent of any buff. It is what a property equals when no buff for that property is currently registered.
- Composed state is what `Entity.GetProperty` (and its typed variants `GetPropertyI`/`GetPropertyF`/`GetPropertyC`) returns: the base value with every currently-registered buff for that property (`Entity.Buffs`, applied via `Property.ApplyBuff`) layered on top. Composed state is a read-time projection, not persisted state — it changes automatically when a buff is added or removed, with no write required.
- Base and composed state are the same value only when zero buffs affect that property. Whenever at least one buff is registered, they are numerically distinct, and that distinction must never be collapsed.

**The write isolation invariant:**
- A write to an entity property is a read-modify-write cycle when it first reads a property's current value (to compute a delta, apply damage/cost/restore, or otherwise derive a new value from the existing one) and then persists a new value back onto the entity.
- Any such read-modify-write MUST read and write base state at both ends. Reading through `GetProperty` for the sole purpose of consuming/displaying the value (with no subsequent write) is unaffected by this rule.
- If a read-modify-write instead reads through `GetProperty` (composed) and persists the result via `RepsertPropertyValue`, `RepsertPropertyCMaxValue`, `UpdatePropertyValue`, `UpdateProperty`, or any equivalent path that writes into `Entity.Properties`, the buff's contribution MUST be subtracted from the read value before it is persisted, so only the base component reaches `Entity.Properties`.
- The buff itself is never persisted into `Entity.Properties`. It remains registered in `Entity.Buffs` and is re-applied by the next `GetProperty` call. A correctly-isolated write leaves the buff's contribution to reappear on the next read via composition, not via having been folded into the stored base value.
- This invariant is a property of the `Entity` type's own read/write contract. It applies to every property capable of carrying an active buff, regardless of the buff's origin (equipped item, skill effect, shield effect, or any future buff source) — the invariant does not vary by who registered the buff.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[rule_entity_property_write_isolation]]` (same-submodule form, used inside `upsilonbattle`) / `@spec-link [[upsilonbattle:rule_entity_property_write_isolation]]` (cross-submodule form, used from `upsilontypes`). Now placed atop every accessor and call site listed below, per the ISS-144/ISS-145 fix.
- **Governs (`upsilontypes/entity/entity.go`):** the composed reader `GetProperty` (entity.go:172) — tagged as the source this invariant protects against, not as compliant itself — plus every base-state accessor: `GetBaseProperty` (205), `GetBasePropertyC` (211), `RepsertPropertyValue` (301), `RepsertPropertyCMaxValue` (310), `RepsertPropertyCValue` (319), `UpdatePropertyValue` (328), `UpdatePropertyCMaxValue` (337), `UpdatePropertyCValue` (346), and `AdjustPropertyCValue` (365) — the delta-based write-isolation primitive introduced by the fix and now the primary implementation path for any buffed-property read-modify-write.
- **Corrected call sites (implementation, not violation, as of ISS-144/ISS-145):** `upsilonbattle/battlearena/ruler/rules/endofturn.go` (poison-tick HP deduction and floor, and Movement restore-to-base-max, via `AdjustPropertyCValue`/`GetBasePropertyC`), `upsilonbattle/battlearena/ruler/rules/skill_validation.go` (`paySkillCost`: HP/MP/SP/Movement cost deduction via `AdjustPropertyCValue`), `upsilonbattle/battlearena/ruler/rules/move.go` (move-cost Movement deduction via `AdjustPropertyCValue` — added to this rule's scope during implementation; it carried the identical Movement-escalation defect but was not in the original enumeration), `upsilonbattle/battlearena/ruler/rules/attack.go` (Shield absorption and HP damage via `AdjustPropertyCValue`), `upsilonbattle/battlearena/property/effect/effectapplicator/effectapplicator.go` (Shield deplete/absorb, HP damage, HP heal, and Shield overshield, all via `AdjustPropertyCValue`).
- **Companion atoms:** `[[mechanic_item_buff_application]]` documents the buff-registration/composition mechanism this rule's write side complements (still read-side only in its own documented scope). `[[mechanic_equipment_stat_bonuses]]` documents the equipment stat-bonus policy that is the buff system's only current source.

## EXPECTATION
- After any write derived from a composed read on a buffed property, the value stored in `Entity.Properties` for that property equals the pre-write base value plus only the intended base-level delta — the active buff's contribution is absent from the persisted value.
- Repeated read-modify-write cycles that apply no new base-level delta (e.g. successive turn-end restores against an unchanged buff) leave the base value unchanged across cycles; they do not escalate it.
- Removing every buff affecting a property (unequip, buff expiry) leaves `GetProperty` returning exactly the last correctly-isolated base value — no buff contribution remains folded into it.
