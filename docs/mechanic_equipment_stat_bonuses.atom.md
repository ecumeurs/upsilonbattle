---
id: mechanic_equipment_stat_bonuses
human_name: Equipment Stat Bonuses Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [equipment, stats, bonuses]
parents:
  - [[upsilonapi:entity_equipment_system]]
dependents:
  - [[mechanic_item_buff_application]]
---

# Equipment Stat Bonuses Mechanic

## INTENT
To implement equipment stat bonus system where equipped items provide direct attribute modifications to characters, with armor adding defense, weapons adding attack power, and utility items providing special effects.

## THE RULE / LOGIC
**Equipment Stat Bonuses Mechanic:**

**Core Principle:**
Equipped items provide dynamic modifications to a character's base attributes, allowing for specialized builds and strategic scaling through gear progression.

**Bonus Classifications:**
- **Direct Additives:** Flat numerical increases to core stats (e.g., +5 Defense from Plate Armor, +8 Attack from a Steel Sword).
- **Percentage Modifiers:** Relative adjustments to existing base stats (e.g., +10% Maximum HP).
- **Property Injection:** Granting entirely new capabilities or exotic stats that the character may not possess natively (e.g., enabling Critical Chance, providing Jump Height bonuses, or granting elemental resistances).
- **Conditional Modifiers:** Bonuses that activate only under specific battlefield conditions (e.g., +5 Attack when health is below 20%).

**Slot-Specific Influence:**
- **Armor (Defensive):** Primary source of Defense, Armor Rating, and physical/elemental mitigation. Heavy armor may incur movement or accuracy penalties.
- **Weapons (Offensive):** Determines base Attack power and effective Range. Weapons often include modifiers for Critical Hits and Backstab effectiveness.
- **Utility (Specialist):** Focuses on resource pool expansion (MP/SP), movement utility (Jump, Speed), and unique utility properties (Stealth detection, flight).

**Stat Aggregation and Stacking Rules:**
- **Additive Stacking:** Bonuses for the same property from different equipment slots (Armor + Utility + Weapon) are summed together to create the final modifier.
- **Resolution Order:** Base Character Stats are calculated first, followed by flat equipment additives, and finally any percentage multipliers.
- **Property Overrides:** Certain equipment properties (like Weapon Range) may override the character's default base values rather than adding to them.

**Lifecycle and Maintenance:**
- **Equipped Status:** Bonuses are strictly active only while the item occupies one of the character's 3 active equipment slots.
- **Dynamic Recalculation:** Any change in equipment (equipping, unequipping, or item destruction) triggers an immediate update of the character's modified stat block.
- **Persistence:** Permanent bonuses remain as long as the item is equipped, whereas conditional bonuses are toggled by the game engine based on real-time state checks.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mechanic_equipment_stat_bonuses]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entity.go`, `upsilonbattle/battlearena/property/def/item.go`
