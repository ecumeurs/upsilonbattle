---
id: mechanic_equipment_shop_inventory_system
status: DRAFT
version: 1.0
parents: []
dependents:
  - [[mechanic_mech_daily_random_shop_roll]]
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To implement equipment shop inventory system where players can browse, filter, and purchase equipment items based on character level, with prices determined by equipment tier and rarity.

## THE RULE / LOGIC
**Equipment Shop Inventory System:**

**Core Principle:**
The shop serves as the primary gateway for character power progression, offering filtered access to equipment based on character level, skill prerequisites, and credit balance.

**Equipment Classifications:**
- **Armor (Defensive):** Head, Body, Hands, Legs, and Feet slots providing Armor Rating and Defense bonuses.
- **Utility (Special):** Neck, Rings, and Belt slots providing stat enhancements, resource pools (MP/SP), and unique effects.
- **Weapons (Offensive):** One-Handed and Two-Handed variants for Melee and Ranged combat, defining attack power and range.

**Inventory Tier System:**
- **Tier 1 (Common):** Entry-level items (e.g., Leather, Iron) with basic stat bonuses.
- **Tier 2 (Uncommon):** Mid-tier items (e.g., Chain, Steel) with enhanced stats and minor MP bonuses.
- **Tier 3 (Rare):** High-tier items (e.g., Plate, Mithral) with significant bonuses and specialized properties.
- **Higher Tiers (Epic/Legendary):** Late-game items with massive multipliers and unique mechanics.

**Shop Logic and Filtering:**
- **Prerequisite Checks:** Items are visible only if the character meets the minimum level requirement and possesses any mandatory precursor skills.
- **Dynamic Filtering:** Users can filter by category (Armor/Weapon/Utility), stat focus (Attack/Defense/HP), or affordability.
- **Sorting:** Inventory can be sorted by price, power rating (SW), tier, or recency.

**Economic Framework:**
- **Pricing Formula:** Final Cost = (Equipment Power Rating) × (Tier Multiplier).
- **Power Rating Calculation:** Derived from the sum of weighted property values (e.g., Damage is weighted higher than HP).
- **Tier Multipliers:** Range from 1.0x for Common to 5.0x for Legendary items.

**Transaction and Inventory Management:**
- **Purchase Workflow:** Validates credit availability and slot compatibility before transferring the item to the character's unequipped inventory.
- **Equipment Slots:** Characters possess three primary active slots (Armor, Utility, Weapon).
- **Recalculation:** Equipping or unequipping an item triggers an immediate recalculation of the character's base and modified statistics.

**UI and Feedback:**
- **Comparison System:** Provides a side-by-side comparison of currently equipped items versus shop candidates, highlighting projected stat changes.
- **History:** Maintains a transaction log of all purchases, including timestamps, item IDs, and final costs.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[equipment_shop_inventory_system]]`
- **Related Files:** Shop management UI, equipment database, character inventory logic
- **Integration:** Works with `mec_credit_spending_shop`, `entity_equipment_system`, `api_equipment_management`

## EXPECTATION
