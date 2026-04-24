---
id: mec_pre_post_execution_costs
human_name: Pre and Post Execution Costs Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [skills, costs, time-based]
parents: []
dependents: []
---

# Pre and Post Execution Costs Mechanic

## INTENT
To implement the dual-cost system for skills with pre-execution costs (SP, MP, Channeling) paid upfront and post-execution costs (Delay) paid after effect completes.

## THE RULE / LOGIC
**Pre and Post Execution Costs Mechanic:**

**Core Principle:**
Skills in Upsilon employ a dual-phase cost structure that separates resource investment (pre-execution) from tactical recovery (post-execution). This ensures that high-impact actions carry both immediate resource risks and subsequent temporal penalties.

**Pre-Execution Costs (Phase 1: Initiation):**
- **Definition:** Costs that must be settled at the exact moment a skill is initiated by the player or AI.
- **Resource Leech:** Includes the immediate deduction of Stamina Points (SP) and Mana Points (MP).
- **Temporal Commitment (Channeling):** Represents a "risk premium" where the character commits to a specific action over a defined period before the effect occurs.
- **Sunk Cost Rule:** Once initiation begins, pre-execution resource costs are considered consumed. If the action is interrupted or the target becomes invalid during the initiation phase, these resources are **not** refunded.

**Execution Phase (Phase 2: Resolution):**
- The skill's primary effect (damage, healing, status application) is resolved.
- For channeled skills, this occurs at the end of the channeling temporary entity's lifecycle.

**Post-Execution Costs (Phase 3: Recovery):**
- **Recovery Delay:** A temporal penalty added to the caster's internal timeline *after* the skill effect has successfully resolved.
- **Timeline Impact:** This delay determines how long the caster must wait before their next turn arrives in the global initiative queue.
- **Example:** A basic attack may have a 0-cost initiation but impose a +100 recovery delay, representing the time taken to return to a combat-ready stance.

**Sequence of Operations:**
1. **Validation:** System verifies the caster has sufficient MP/SP and isn't currently under a restrictive status (e.g., Stun).
2. **Settlement:** Deduct all initiation costs (MP, SP).
3. **Commitment:** If the skill requires channeling, a temporary entity is spawned to track the initiation time.
4. **Resolution:** Upon completion of any commitment phase, the skill's logic executes on the target.
5. **Recovery:** The caster's initiative delay is updated by adding the skill's defined recovery cost.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_pre_post_execution_costs]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/skill.go`
