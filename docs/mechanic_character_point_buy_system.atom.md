---
id: mechanic_character_point_buy_system
human_name: "Character Point Buy System"
status: DRAFT
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents:
  - [[shared:rule_progression]]
dependents: []
---

# Character Point Buy System

## INTENT
To implement the battle-engine side of the V2 Character Point (CP) point-buy system — both initial character creation (a 100 CP allocation replacing random rolls) and ongoing post-win progression spend — strictly enforcing the economy and taxonomy defined by [[shared:rule_progression]].

## THE RULE / LOGIC
**Core Principle:**
Players allocate a Character Point (CP) pool across attributes instead of receiving random stats, giving full agency over a character's mechanical identity. This atom is the engine implementation of [[shared:rule_progression]]; that rule is authoritative for all costs, the CP cap, and the stat taxonomy.

**CP Economy (per +1 increment, mirrors `[[shared:rule_progression]]`):**
- **Class A — CP-upgradable (Standard):** HP 1 CP · MP 1 CP · SP 1 CP · Attack 5 CP · Defense 5 CP · Movement 30 CP (per cell).
- **Class A — CP-upgradable (Exotic):** JumpHeight 15 CP · CritChance (+1%) 10 CP · CritDamage/CritMultiplier (+5%) 5 CP.
- **Class B — NOT CP-upgradable:** AttackRange, Shield, Accuracy, Dodge (Evasion) are granted only by equipment, buffs, and skills — they never appear in the CP allocation UI.

**Initial Allocation (Character Creation):**
- Every new character starts with exactly **100 CP**.
- **Mandatory baselines:** stats cannot be reduced below the V2 floor (HP 30, Attack 10, Defense 5, Movement 3) to refund CP; non-negativity per `[[shared:rule_progression]]`.
- Creation completes only when the full 100 CP pool is spent.

**Progression Spend (Post-Win):**
- Each recorded victory raises the allowed CP cap by **+10 CP** (cap = `100 + total_wins*10`).
- **Unrestricted spend:** there are NO per-stat purchase caps and NO level/win-gating on any attribute (the legacy "once every 5 wins" Movement gate is removed) — each stat self-balances through its CP cost. Effective-value bounds (e.g. CritChance clamped 0–100% at hit resolution) are applied in the combat pipeline, not at purchase time.
- Players may bank unspent CP toward high-cost upgrades (e.g. Movement).

**Validation:**
- **Affordability:** reject any allocation where `spent_cp + cost > 100 + total_wins*10`.
- **Class enforcement:** reject CP spend on any Class B stat.
- **Non-negativity / baseline:** reject reductions below the V2 floor.
- Successful allocation immediately updates the character's persisted attribute block and deducts CP.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_character_point_buy_system]]`
- **Authoritative Rule:** `[[shared:rule_progression]]`
- **Related Files:** character creation logic, character attribute schema
- **API Endpoints:** `POST /api/v1/character/create` (creation-time CP allocation), progression allocation endpoint
- **UI Components:** character-creation + stat-allocation forms with CP spend display

## EXPECTATION
- A new character begins with exactly 100 CP; creation is rejected until the full pool is allocated.
- Allocating 1 Attack costs 5 CP; 1 HP/MP/SP costs 1 CP each; 1 Movement cell costs 30 CP; +1% CritChance costs 10 CP.
- An allocation is accepted only when `spent_cp + cost <= 100 + total_wins*10`; otherwise rejected.
- Attempting to reduce any stat below its V2 baseline (HP 30 / Attack 10 / Defense 5 / Movement 3) is rejected.
- Attempting to spend CP on a Class B stat (AttackRange, Shield, Accuracy, Dodge) is rejected.
- A recorded victory raises the CP cap by exactly 10; no attribute is subject to a per-stat purchase cap or a level/win gate.
