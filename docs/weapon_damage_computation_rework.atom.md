---
id: weapon_damage_computation_rework
status: DRAFT
parents: []
dependents: []
version: 1.0
---

# New Atom

## INTENT
To implement weapon-based attack computation system where equipped weapons transform basic attacks into skill-based attacks, utilizing weapon properties as skill attributes for full damage calculation.

## THE RULE / LOGIC
**Weapon Damage Computation Rework:**

**Core Principle:**
Combat resolution is unified by transforming every weapon-based attack into a temporary "Skill" execution. This ensures that basic attacks benefit from the full complexity of the Skill System, including critical hits, accuracy checks, and advanced mitigation formulas.

**Weapon-to-Skill Transformation Logic:**
- **Dynamic Skill Generation:** Upon initiating an attack, the system constructs a virtual skill block using the weapon's current properties.
    - **Offense:** `WeaponBaseDamage` is treated as the skill's base power.
    - **Reach:** `WeaponRange` defines the skill's targeting distance.
    - **Precision:** Weapon `CritChance` and `CritMultiplier` are added to the character's baseline exotic attributes.
    - **Efficiency:** `AttackSpeed` (or Cooldown) from the weapon determines the recovery `Delay` applied after execution.
- **Categorical Behavior:**
    - **Melee Weapons:** Inherit high accuracy (100% baseline) and higher base damage but limited range (1-2).
    - **Ranged Weapons:** Inherit high range (4-7) but lower base damage and variable accuracy (80-95%) to balance their distance advantage.
    - **Heavy Weapons (2H):** Provide massive damage bonuses but restrict the character's utility slot and impose higher recovery delays.

**Resolution Pipeline:**
1. **Initiation:** The engine identifies the weapon in the character's active slot.
2. **Synthesis:** A temporary skill object is generated, combining weapon attributes with the character's total effective stats.
3. **Execution:** The attack enters the standard Skill Resolution Pipeline:
    - **Hit Test:** Accuracy vs Evasion/Dodge.
    - **Crit Test:** Combined Critical Chance check.
    - **Mitigation:** Base Damage vs Target Defense and Armor.
    - **Multipliers:** Positional bonuses (Backstabs) and elemental synergies are applied.
4. **Finalization:** The resulting damage is deducted from the target, and the recovery delay is applied to the attacker.

**System Integration and Benefits:**
- **Unified Logic:** Removes redundant code by processing all combat actions through the same mathematical pathway.
- **Enhanced Customization:** Allows for specialized weapons (e.g., Elemental, Shield-Piercing, High-Crit) to be implemented using existing skill properties.
- **Preview Accuracy:** Enables the UI to provide a highly accurate "Expected Damage" preview by running a dry-run of the synthesized skill.
- **Scalability:** New weapon types can be introduced as pure data definitions without requiring updates to the core combat rules.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[weapon_damage_computation_rework]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonbattle/battlearena/entity/skill/skill.go`
- **Integration:** Works with `mec_weapon_as_skill_system`, `armor_penetration_system`, `backstab_detection_algorithm`

## EXPECTATION
