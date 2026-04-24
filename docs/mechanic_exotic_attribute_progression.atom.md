---
id: mechanic_exotic_attribute_progression
status: DRAFT
priority: 5
version: 2.0
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To implement exotic attribute progression mechanics for V2 characters, defining how attributes like Critical Chance, Critical Multiplier, Dodge, Accuracy, and Jump Height increase through level-based progression and character point allocation.

## THE RULE / LOGIC
**Exotic Attribute Progression Mechanic:**

**Core Principle:**
Exotic attributes represent specialized combat capabilities that offer tactical depth beyond standard durability and offense. Their progression is gated by a tiered level-restriction system to ensure that high-value tactical advantages are earned through consistent character development.

**Tiered Level Restrictions:**
Attribute increases are permitted according to the following temporal cadence:
- **Phase 1 (Every Level):** Characters can always invest points in **Standard Stats** (HP, Attack, Defense) and **Critical Multiplier**.
- **Phase 2 (Every 3 Levels):** Unlocks the ability to increase precision and avoidance stats: **Accuracy** and **Dodge**.
- **Phase 3 (Every 5 Levels):** Unlocks high-impact tactical stats: **Critical Chance**, **Jump Height**, and **Movement**.

**Attribute Specification and Caps:**
1. **Critical Chance:**
    - **Utility:** Determines the frequency of critical hits.
    - **Economic Constraint:** 10 CP per +1% increment.
    - **Boundaries:** Hard-capped at 30% to prevent guaranteed critical damage.
2. **Critical Multiplier:**
    - **Utility:** Dictates the damage bonus applied when a critical hit occurs.
    - **Economic Constraint:** 5 CP per +5% increment.
    - **Boundaries:** Maximum of 200% (2x damage).
3. **Accuracy & Dodge:**
    - **Utility:** Accuracy ensures skill execution hits the target; Dodge provides a percentage chance to evade incoming skill-based effects.
    - **Economic Constraint:** Accuracy (3 CP per +2%), Dodge (5 CP per +1%).
    - **Boundaries:** Dodge is capped at 30% to prevent "untouchable" character builds.
4. **Jump Height:**
    - **Utility:** Increases the character's ability to traverse vertical terrain obstacles and reach high-ground advantage points.
    - **Economic Constraint:** 15 CP per +1 height unit.
    - **Boundaries:** Capped at 5 units to maintain map integrity.

**Progression Logic:**
- **Financial Validation:** The system confirms sufficient Character Points (CP) are available for the upgrade.
- **Milestone Verification:** The system checks the character's level against the last recorded increase for that specific exotic attribute to ensure the restriction interval (3 or 5 levels) has been respected.
- **Build Archetyping:** By gating these stats, the system encourages distinct character archetypes (e.g., Evasive Scouts, Heavy Hitters, Precision Assassins) as players reach level milestones.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[exotic_attribute_progression]]`
- **Related Files:** Character progression system, stat validation logic
- **UI Components:** Character sheet with exotic attributes display and upgrade buttons

## EXPECTATION
