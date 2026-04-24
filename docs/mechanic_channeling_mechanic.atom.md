---
id: mechanic_channeling_mechanic
human_name: Channeling Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [time-based, skills, casting]
parents:
  - [[mech_temporary_entity_system]]
dependents: []
---

# Channeling Mechanic

## INTENT

To implement channeling mechanics where skills have a pre-execution delay (casting time) during which the caster is vulnerable and the effect is delayed until the casting completes.

## THE RULE / LOGIC
**Channeling Mechanic:**

**Core Principle:**
Channeling represents a specialized skill execution phase where a character commits to a high-impact action that requires significant preparation time. This creates a window of vulnerability and tactical risk for the caster in exchange for more efficient Skill Weight (SW) balancing.

**Resource and Risk Economy:**
- **Vulnerability Premium:** Channeled skills provide more powerful effects for their cost compared to immediate skills (balancing at -15 SW per 10 delay units).
- **Sunk Costs:** All pre-execution resources (SP, MP, HP) are deducted at the start of the channeling phase. These are not refunded if the channeling is interrupted or fails.
- **Interruption Threshold:** Taking damage while channeling accumulates "Interruption Points." If this total exceeds the character's stability threshold, the channeling fails immediately.

**Channeling Process Lifecycle:**
1. **Initiation:** The player selects a channeled skill and target.
2. **Entity Generation:** The system spawns a hidden `TimeBased` temporary entity at the caster's location to track the channeling duration independently of the caster's own timeline.
3. **Timeline Integration:** The channeling entity is inserted into the global turn queue (Turner) with a delay equal to the skill's defined channeling time.
4. **Active Phase:** The caster is marked as "Channeling," typically restricting further movement or action until resolution.
5. **Resolution:** When the channeling entity's turn arrives:
    - The stored skill effect is resolved against the original target.
    - Post-execution recovery delay is applied to the caster's internal timeline.
    - The temporary channeling entity is removed from the grid and state.

**System Integration:**
- **Temporal Decoupling:** By using a separate entity to track the channeling time, the system ensures that the caster's own turn order is correctly managed relative to the completion of the spell.
- **Visual Feedback:** The interface provides indicators (e.g., "Channeling: Fireball") to both the player and opponents, highlighting the tactical window available for interruption.
- **Interruption Resolution:** A failed channeling event triggers a specific "Interrupted" state, immediately cleaning up the temporary entity and returning the caster to a neutral recovery state.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mec_channeling_mechanic]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/skill.go`, `upsilonbattle/battlearena/ruler/rules/beginingofturn.go`
- **Integration:** Works with `mech_temporary_entity_system`, `mech_entity_expiration`

## EXPECTATION
