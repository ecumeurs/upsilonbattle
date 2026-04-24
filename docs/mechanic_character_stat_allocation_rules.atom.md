---
id: mechanic_character_stat_allocation_rules
status: DRAFT
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents: []
dependents: []
---

# New Atom

## INTENT
To establish character stat allocation rules defining point costs, maximum caps, and level restrictions for both standard and exotic attributes in V2 progression system.

## THE RULE / LOGIC
**Character Stat Allocation Rules:**

**Core Principle:**
Attribute growth in the V2 system is governed by a weighted point-buy economy. Costs reflect the tactical potency of each stat, while level-based restrictions and caps prevent extreme specialization and ensure balanced progression.

**Standard Attribute Configuration:**
- **Health (HP):**
    - **Cost:** +1 HP per 1 CP.
    - **Constraints:** No hard maximum cap; designed for linear scaling to encourage diverse durability builds.
- **Offense (Attack):**
    - **Cost:** +1 Attack per 5 CP.
    - **Impact:** High cost reflects multiplicative damage scaling.
- **Mitigation (Defense):**
    - **Cost:** +1 Defense per 5 CP.
    - **Constraints:** Subject to a soft cap relative to the character's Attack value to prevent over-specialization in passivity.
- **Mobility (Movement):**
    - **Cost:** +1 cell per 30 CP.
    - **Restrictions:** Maximum of 6 cells. Increases are restricted to once every **5 levels** to preserve tactical grid balance.

**Exotic Attribute Configuration:**
- **Critical Chance:**
    - **Cost:** +1% per 10 CP.
    - **Restrictions:** Hard cap of 30%. Increases restricted to once every **5 levels**.
- **Critical Multiplier:**
    - **Cost:** +5% per 5 CP.
    - **Restrictions:** Maximum of 200%. No level-based frequency restriction.
- **Accuracy:**
    - **Cost:** +2% per 3 CP.
    - **Restrictions:** Maximum of 120%. Limited to once every **3 levels**.
- **Dodge (Evasion):**
    - **Cost:** +1% per 5 CP.
    - **Restrictions:** Hard cap of 30%. Limited to once every **3 levels**.
- **Jump Height:**
    - **Cost:** +1 height increment per 15 CP.
    - **Restrictions:** Maximum height of 5. Limited to once every **5 levels**.

**Progression and Validation Logic:**
- **Financial Validation:** The system must verify the character has an unspent CP balance greater than or equal to the total cost of the requested increase.
- **Temporal Validation:** For restricted stats (Movement, Crit, Jump, etc.), the system verifies that the required number of level-ups have occurred since the last increase of that specific attribute.
- **Integrity Validation:** Stat purchases that would exceed a hard cap are automatically rejected.
- **Immediate Resolution:** Successful allocation results in an immediate update to the character's modified stat block and a permanent deduction from their available CP pool.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[character_stat_allocation_rules]]`
- **Related Files:** Character progression logic, stat validation system
- **UI Components:** Stat allocation interface with cost displays and restriction warnings

## EXPECTATION
