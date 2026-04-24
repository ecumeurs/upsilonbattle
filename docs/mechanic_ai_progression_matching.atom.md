---
id: mechanic_ai_progression_matching
status: DRAFT
version: 1.0
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To implement AI progression matching system where AI characters follow the same point-buy progression rules as players, scaling stats, skills, and level according to player team averages.

## THE RULE / LOGIC
**AI Progression Matching:**

**Core Principle:**
AI characters adhere to the same progression mechanics as players to ensure competitive fairness and tactical transparency.

**Synchronization Rules:**
- **Baseline Stats:** AI starts with V2 baseline values (HP 30-50, ATK 10, DEF 5, MOV 3).
- **Point Economy:** AI receives 100 Character Points (CP) at creation and accumulates +10 CP per victory, identical to player rewards.
- **CP Costs:** Purchasing stats (HP, Attack, Defense, Movement) and exotic attributes (Crit, Dodge, Accuracy) uses the standard player-facing cost table.

**Level Matching Logic:**
- **Average Calculation:** The system determines the average level of the player team as the primary baseline for AI generation.
- **Distribution Patterns:**
    - **Exact Match:** AI levels precisely match the player average (standard for 1v1).
    - **Variance:** AI levels may fluctuate by ±1 level relative to the average.
    - **Progressive Spread:** In multi-character matches, AI levels are distributed around the average to provide a mix of "elites" and "grunts."

**Skill Grade and Selection:**
- **Grade Access:** AI skill availability is restricted by level-based grade thresholds (e.g., Grade III becomes available at Level 10).
- **Selection Algorithm:** Every 10 levels, the AI selects one skill from a pool of three random candidates filtered by its specific archetype.
- **Reforging:** AI can perform skill reforging every 5 levels, following player-layer logic.

**Difficulty Scaling and Balancing:**
- **Dynamic Adjustment:** The AI's CP pool can be adjusted by ±2 CP based on the player's recent win rate to maintain a consistent challenge.
- **Personality Matrix:** Difficulty levels (Easy, Normal, Hard) define the AI's personality parameters, including Aggression, Risk Tolerance, Adaptiveness, and Team Coordination.

**Fair Play Constraints:**
- **Parity:** No hidden bonuses are granted to AI outside of the CP and Level matching system.
- **Transparency:** Players should be able to deduce AI builds based on visible performance and level parity.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ai_progression_matching]]`
- **Related Files:** AI creation logic, progression tracking, matchmaking integration
- **Integration:** Works with `mec_ai_archetype_system`, `ai_controller_archetypes`, `character_stat_allocation_rules`

## EXPECTATION
