---
id: mechanic_armor_penetration_system
status: DRAFT
version: 1.0
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To implement armor penetration system where certain attacks (like backstabs) can bypass partial armor protection, reducing damage mitigation effectiveness of target's defensive bonuses.

## THE RULE / LOGIC
**Armor Penetration System:**

**Core Principle:**
Armor penetration allows specific attacks to ignore a portion of the target's Armor Rating, facilitating damage delivery against high-defense builds.

**Damage Calculation Framework:**
- **Standard Formula:** Damage = Attack Power - (Defense + Armor Rating) - Shield Value.
- **Penetration Formula:** Damage = Attack Power - (Defense + (Armor Rating * (1 - Penetration %))) - Shield Value.
- **Clamping:** All damage calculations are clamped to a minimum of 1 damage.

**Backstab Penetration Mechanic:**
- **Penetration Amount:** Backstabs automatically ignore **50%** of the target's Armor Rating.
- **Synergy:** This effect stacks with the standard backstab damage multiplier (e.g., 1.5x) to maximize effectiveness against "tanks."

**Penetration Classifications:**
- **Percentage-Based:** Reduces armor by a relative fraction (e.g., Backstabs, certain high-tier skills).
- **Fixed/Flat:** Ignores a specific integer amount of armor (e.g., Piercing ammunition, basic weapon properties).
- **Shield Interaction:** Standard armor penetration does **not** bypass temporary Shields. Shields must be depleted or bypassed using specialized "Shield Penetration" effects.

**Stacking and Priority Rules:**
- **Non-Additive:** Percentage-based penetration effects do not stack additively. The system applies only the highest available penetration percentage.
- **Resolution Order:** Defense is always calculated before Armor Rating; penetration only modifies the Armor Rating component of the mitigation.

**Equipment and Attributes:**
- **Piercing Weapons:** Inherit a base penetration percentage (e.g., 25%).
- **Daggers:** Grant additional penetration bonuses specifically for backstab maneuvers.
- **Temporary Buffs:** Skills or potions can grant short-term penetration attributes.

**Balance and Dynamics:**
- **Anti-Tank Utility:** Penetration prevents high-armor characters from becoming functionally invincible to low-damage attacks.
- **Strategic Layer:** Creates a counter-play dynamic where heavy armor is effective against standard volleys but vulnerable to flanking and specialized piercing attacks.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[armor_penetration_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonbattle/battlearena/entity/entity.go`
- **Integration:** Works with `backstab_detection_algorithm`, `mec_backstabbing_mechanic`, `weapon_damage_computation_rework`

## EXPECTATION
