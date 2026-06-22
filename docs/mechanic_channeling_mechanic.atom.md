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
  - [[upsilontypes:mechanic_temporary_entity_system]]
dependents: []
---

# Channeling Mechanic

## INTENT

To implement channeling mechanics where skills have a pre-execution delay (casting time) during which the caster is vulnerable and the effect is delayed until the casting completes.

## THE RULE / LOGIC
**Channeling Mechanic:**

**Core Principle:**
Channeling represents a specialized skill execution phase where a character commits to a high-impact action that requires significant preparation time. This creates a window of vulnerability and tactical risk for the caster in exchange for more efficient Skill Weight (SW) balancing.

**Channeling Cost (Pre-Execution):**
- **Property:** Channeling is measured in delay units (e.g. Fireball with Channeling 400 = 400 delay before the effect resolves).
- **Risk Premium / Vulnerability Premium:** Channeling costs **-15 SW per 10 delay** units (vs -10 SW for normal delay), making channeled skills more powerful for their cost than immediate skills.
- **Sunk Costs:** All pre-execution resources (SP, MP, HP, Channeling delay) are deducted upfront at the start of the channeling phase. These are NOT refunded if the channeling is interrupted or fails.
- **Post-Execution Costs:** Recovery delay is added to the caster's timeline after the effect completes.

**Channeling Process Lifecycle:**
1. **Initiation:** The player selects a channeled skill and target.
2. **Entity Generation:** The system spawns a hidden `TimeBased` temporary channeling entity at the target/caster location to track the channeling duration independently of the caster's own timeline.
3. **Timeline Integration:** The channeling entity is inserted into the global turn queue (Turner) with a delay equal to the skill's defined channeling time; the caster is added to the `IsCasting` state.
4. **Active Phase:** The caster is marked as "Channeling," typically restricting further movement or action and remaining vulnerable to interruption until resolution.
5. **Resolution:** When the channeling entity's turn arrives, the stored skill effect resolves against the original target, post-execution recovery delay is applied, the temporary channeling entity dies, and the caster is released.

**Interruption Mechanics:**
- **Interruption Property:** 0-100, fills when the caster takes damage while casting.
- **Interruption Formula:** Damage-based accumulation — **1 damage = 10 interruption points**.
- **Failure Threshold:** When Interruption **≥ 100**, channeling fails immediately, all pre-execution resources are wasted (not refunded), and the caster is released into a neutral "Interrupted"/recovery state with the temporary entity cleaned up.

**System Integration:**
- **Temporal Decoupling:** Using a separate entity to track channeling time ensures the caster's own turn order is correctly managed relative to spell completion.
- **Visual Feedback:** The interface provides indicators (e.g. "Channeling: Fireball") to both player and opponents, highlighting the interruption window.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_channeling_mechanic]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/skill.go`, `upsilonbattle/battlearena/ruler/rules/beginingofturn.go`
- **Integration:** Works with `[[upsilontypes:mechanic_temporary_entity_system]]`, `[[mechanic_expiration_controller]]`

## EXPECTATION
- A channeled skill deducts all pre-execution costs (SP, MP, HP, Channeling delay) at initiation; these are not refunded on interruption or failure.
- A channeled skill spawns a temporary `TimeBased` entity inserted into the Turner with delay equal to the skill's Channeling value (e.g. 400 delay for Channeling 400); the caster enters `IsCasting` state.
- The skill effect resolves only when the channeling entity's turn arrives, after which the entity dies and the caster is released.
- Taking damage while channeling adds interruption points at 10 points per 1 damage; reaching ≥ 100 interruption fails the channel, wastes resources, and releases the caster.
- Channeled skills are balanced at -15 SW per 10 delay units (vs -10 SW for normal delay).
