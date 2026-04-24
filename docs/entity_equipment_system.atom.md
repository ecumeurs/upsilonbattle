---
id: entity_equipment_system
human_name: Equipment System Entity
type: ENTITY
layer: ARCHITECTURE
version: 2.0
status: DRAFT
priority: 5
tags: [equipment, inventory, progression]
parents: []
dependents:
  - [[api_equipment_management]]
  - [[mec_equipment_stat_bonuses]]
  - [[mec_three_slot_equipment_system]]
  - [[mec_weapon_as_skill_system]]
---

# Equipment System Entity

## INTENT
To define the equipment system structure with 3-slot inventory (1 armor, 1 utility, 1 weapon) and weapon-as-skill mechanics.

## THE RULE / LOGIC
**Equipment Slots:**
- **Armor Slot:** One piece of armor (provides ArmorRating)
- **Utility Slot:** One utility item (special effects, buffs)
- **Weapon Slot:** One weapon (transforms basic attack into skill-based attack)

**Equipment Properties:**
- **Armor:** ArmorRating (defense bonus), Durability
- **Weapons:** WeaponRange, WeaponBaseDamage, WeaponType, CritChance
- **Utility:** Special effects, passive buffs, consumable items

**Weapon-as-Skill System:**
- Equipped weapon transforms basic attack into skill-based attack
- Weapon properties become skill properties (Range, Damage, Crit)
- Basic attack uses skill computation instead of simple formula
- Enables weapon variety without separate attack systems

**Equipment Effects:**
- **Stat Bonuses:** Direct attribute modifications (Attack +2, Defense +1)
- **Skill Augmentation:** Weapons add properties to basic attack
- **Passive Effects:** Utility items grant ongoing benefits

**Database Schema:**
- `character_equipment`: Tracks equipped items per slot
- `equipment_library`: Available items with properties and costs
- `character_inventory`: Owned but unequipped items

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_equipment_system]]`
- **Laravel Model:** `App\Models\Character`, `App\Models\Equipment`
