---
id: rule_ai_team_composition_rules
status: DRAFT
version: 1.0
parents: []
dependents: []
type: RULE
layer: CUSTOMER
---

# New Atom

## INTENT
To implement AI team composition rules that enforce maximum limits on archetype variety per AI team, ensuring balanced team composition with appropriate support and specialist limitations.

## THE RULE / LOGIC
**AI Team Composition Rules:**

**Core Principle:**
AI team structures are governed by strict archetype constraints to ensure tactical variety and prevent unbalanced combat encounters.

**Archetype Constraints:**
- **Maximum Limits:** Each AI team is restricted to a maximum of **one Support** and **one Sneak**.
- **Flexible Slots:** No maximum limit is enforced for **Fighter** or **Ranger** archetypes.
- **Team Size Cap:** An AI team cannot exceed a total of 5 members.

**Team Generation Hierarchy:**
The system populates teams sequentially based on total size:
1. **1v1:** Any archetype allowed (no constraints).
2. **2v2:** Prioritizes 1 Fighter and 1 Ranger.
3. **3v3:** Adds 1 Support to the base 2v2 composition.
4. **4v4:** Adds 1 Sneak to the base 3v3 composition.
5. **5v5:** Fills the final slot with an additional Fighter or Ranger.

**Composition Validation Logic:**
- **Mandatory Check:** Before match start, the team must be validated against the archetype maximums.
- **Fallback:** If a generated team violates limits, the selection algorithm must re-run until a valid composition is achieved.

**Difficulty and Scaling Integration:**
- **Level Sync:** All AI team members are scaled to the average level of the player team.
- **Dynamic Aggression:** AI behavior (targeting precision, risk tolerance) scales based on the player's recent win rate.
- **Archetype Synergy:** Advanced generation algorithms prioritize templates that maximize combat synergy (e.g., combining a frontline Fighter with a protective Support).

**Matchmaking and Substituted Roles:**
- **Strategic Variety:** Matchmaking uses multiple pre-defined templates (Balanced, Aggressive, Tactical, Ranged Focus) to keep encounters unpredictable.
- **Loss Adaptation:** If a specialized AI (Support/Sneak) is eliminated, the remaining team members adjust their tactical priorities rather than attempting to fill the vacant role.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ai_team_composition_rules]]`
- **Related Files:** `upsilonbattle/battlearena/controller/teamcomposition.go`, matchmaking logic
- **Integration:** Works with `ai_controller_archetypes`, `mec_ai_archetype_system`, `ai_progression_matching`

## EXPECTATION
