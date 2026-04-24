---
id: mechanic_character_point_buy_system
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
To implement character creation point-buy system where players receive 100 Character Points (CP) to strategically allocate attributes instead of receiving 4 random points on base stats.

## THE RULE / LOGIC
**Character Point-Buy System:**

**Core Principle:**
Replaces randomized stat growth with a strategic Point-Buy system, granting players agency over their character's mechanical identity through the allocation of a starting Character Point (CP) pool.

**Economic Framework (Costs per +1 Increment):**
- **Durability (HP):** 1 CP per point. Linear and affordable, enabling high health pools.
- **Combat Proficiency (Attack/Defense):** 5 CP per point. Reflects the significant impact of raw offensive and defensive scaling.
- **Tactical Mobility (Movement):** 30 CP per cell. High cost acts as a natural inhibitor for excessive movement range.
- **Exotic Attributes:**
    - **Critical Chance (+1%):** 10 CP.
    - **Critical Multiplier (+5%):** 5 CP.
    - **Jump Height (+1 tile):** 15 CP.

**Allocation Rules and Constraints:**
- **Initial Grant:** Every new character begins with exactly **100 CP**.
- **Mandatory Minimums:** Stats cannot be reduced below the V2 baseline (HP 30, Attack 10, Defense 5, Movement 3) to "refund" points.
- **Full Allocation:** The character creation process is only complete when the entire 100 CP pool has been spent.
- **Attribute Caps:** Exotic stats are subject to secondary soft caps (e.g., maximum Critical Chance thresholds) to prevent game-breaking specialization.

**Progression Integration:**
- **Accumulation:** Characters earn **+10 CP** for every recorded victory, allowing for continuous growth beyond the initial 100 CP.
- **Scaling Impacts:** The transition to 100 CP (versus the legacy 4-point system) allows for more granular skill percentage calculations and meaningful equipment interactions.
- **Flexibility:** Players may choose to "save" CP from victories to purchase high-cost upgrades (like Movement) in the future.

**System Validation:**
- **Integrity Check:** The game engine validates every CP expenditure against the character's current level and total accumulated points.
- **Schema Enforcement:** CP allocations are persisted as part of the character's primary attribute block, ensuring consistency across matches and UI displays.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[character_point_buy_system]]`
- **Related Files:** Character creation logic, character database schema
- **API Endpoints:** `POST /api/v1/character/create` (with CP allocation)
- **UI Components:** Character creation form with CP spending display

## EXPECTATION
