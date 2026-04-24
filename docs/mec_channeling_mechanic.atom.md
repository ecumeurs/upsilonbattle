---
id: mec_channeling_mechanic
human_name: Channeling Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [time-based, skills, casting]
parents:
  - [[mechanic_mech_temporary_entity_system]]
dependents: []
---

# Channeling Mechanic

## INTENT
To implement channeling mechanics where skills have a pre-execution delay (casting time) during which the caster is vulnerable and the effect is delayed until the casting completes.

## THE RULE / LOGIC
Channeling represents skills that require time to cast before the effect takes place:

**Channeling Cost (Pre-Execution):**
- **Property:** Channeling (measured in delay units)
- **Example:** Fireball with Channeling 400 = 400 delay before effect
- **Risk Premium:** Channeling costs -15 SW per 10 delay (vs -10 SW for normal delay)

**Channeling Process:**
1. Player selects channeling skill
2. Create temporary channeling entity at target location
3. Add caster to IsCasting state
4. Channeling entity added to Turner with casting delay
5. During channeling, caster is vulnerable to interruption
6. When channeling entity's turn arrives, effect executes
7. Channeling entity dies, caster is released

**Interruption Mechanics:**
- **Interruption Property:** 0-100, fills when caster takes damage while casting
- **Interruption Formula:** Damage-based accumulation (1 damage = 10 interruption points)
- **When Interruption ≥ 100:** Channeling fails, resources wasted, caster released

**Cost Timing:**
- **Pre-Execution Costs:** SP, MP, Channeling delay paid upfront
- **Post-Execution Costs:** Delay added after effect completes
- **Resource Risk:** Interrupted channeling still consumes pre-execution costs

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_channeling_mechanic]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/skill.go`
