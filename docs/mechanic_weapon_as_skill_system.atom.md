---
id: mechanic_weapon_as_skill_system
human_name: Weapon-as-Skill System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority=5
tags=[equipment, weapons, combat]
parents:
  - [[upsilonapi:entity_equipment_system]]
dependents: []
---

# Weapon-as-Skill System Mechanic

## INTENT
To implement the weapon-as-skill system where equipped weapons transform basic attacks into skill-based attacks with properties like range, damage, critical chance, enabling weapon variety through the skill system.

## THE RULE / LOGIC
**Weapon-as-Skill System Mechanic:**

**Core Principle:**
This system unifies basic attacks and specialized abilities by treating equipped weapons as "permanent skills." This allows the combat engine to use a single, robust calculation pathway for all offensive actions.

**Weapon-to-Skill Mapping:**
Equipping a weapon dynamically generates a virtual skill entry using the following property translations:
- **Range:** The weapon's `WeaponRange` property dictates the skill's effective distance and targeting constraints.
- **Offense:** `WeaponBaseDamage` is mapped directly to the skill's primary damage property.
- **Precision:** `CritChance` and `CritMultiplier` from the weapon are integrated into the skill's critical hit resolution logic.
- **Specialized Behavior:** Attributes like "Backstab Bonus" or "Shield Penetration" are treated as skill-level effects.

**Combat Execution Flow:**
1. **Trigger:** A player or AI initiates a "Basic Attack."
2. **Detection:** The system checks the character's active **Weapon Slot**.
3. **Branching Logic:**
    - **Weapon Equipped:** The system constructs a temporary skill block using the weapon's current attributes. This "Weapon Skill" is then passed through the standard **Skill Resolution Pipeline**, inheriting all global modifiers, evasion checks, and damage mitigation rules.
    - **Unarmed (Fallback):** If the slot is empty, the system defaults to a baseline "Unarmed Attack" with minimal range and fixed low damage.
4. **Resolution:** The final damage value is calculated based on the weapon skill's properties versus the target's defensive attributes (Defense, Armor, Shields).

**Weapon Variety through Skill Properties:**
- **Melee Weapons (Swords/Axes):** Generate skills with short range (1-2) and high base damage.
- **Ranged Weapons (Bows/Crossbows):** Generate skills with high range (4+) and varying accuracy/crit profiles.
- **Utility Weapons (Daggers/Staves):** Generate skills with unique modifiers like enhanced backstab penetration or secondary resource bonuses.

**System Benefits:**
- **Code Consistency:** Eliminates separate logic branches for "Basic Attacks" versus "Special Skills."
- **Scalability:** New weapon types can be introduced by simply defining new property sets without modifying the core attack logic.
- **Transparency:** Players can view their "Weapon Skill" in the character sheet to understand exactly how their equipment influences their combat performance.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mechanic_weapon_as_skill_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonbattle/battlearena/property/def/item.go`
