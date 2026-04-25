---
id: mec_three_slot_equipment_system
human_name: Three-Slot Equipment System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [equipment, slots, inventory]
parents:
  - [[entity_equipment_system]]
dependents:
  - [[upsilonapi:api_equipment_management]]
---

# Three-Slot Equipment System Mechanic

## INTENT
To implement the simplified 3-slot equipment system with exactly 1 armor slot, 1 utility slot, and 1 weapon slot per character, providing focused equipment progression without inventory management complexity.

## THE RULE / LOGIC
**Three-Slot Equipment System Mechanic:**

**Core Principle:**
To maximize tactical focus and minimize inventory management overhead, every character is restricted to a streamlined three-slot equipment architecture.

**The Three Active Slots:**
1. **Armor Slot:** Accommodates one piece of defensive gear (e.g., Helmets, Chestplates, or Greaves).
2. **Utility Slot:** Reserved for one specialized accessory (e.g., Rings, Amulets, or Belts) that provides unique bonuses or resource enhancements.
3. **Weapon Slot:** Houses the primary offensive tool, supporting One-Handed or Two-Handed Melee and Ranged variants.

**Equipping Rules and Constraints:**
- **Exclusivity:** A character may have exactly one active item per slot. Equipping a new item into an occupied slot automatically moves the previous item back to the character's unequipped inventory.
- **Categorical Integrity:** Items are tagged with specific types (Armor, Utility, Weapon) and can only be equipped in their corresponding slot.
- **Two-Handed Special Case:** Large weapons occupy the Weapon slot but may impose restrictions on the Utility slot to balance their higher base damage output.
- **Instant Stat Resolution:** Swapping equipment triggers an immediate recalculation of the character's modified attributes (e.g., updated Defense or Attack power).

**Inventory and Storage Logic:**
- **Unequipped State:** Items not currently in an active slot are stored in a personal inventory.
- **Capacity:** The unequipped inventory is unrestricted by weight or volume, allowing players to collect diverse equipment without penalty.
- **Accessibility:** Swapping between active slots and the inventory is permitted between matches and does not incur any economic or temporal cost.

**Design Philosophy:**
- **Clarity:** The three-slot system provides a highly readable interface for players to assess their own power and that of their opponents.
- **Progression:** Simplifies the upgrade path by ensuring every purchase decision is a direct replacement for an existing capability, making "stat-up" choices more impactful.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_three_slot_equipment_system]]`
- **Related Files:** `upsilonapi/api/input.go` (Character structure), equipment database schema
