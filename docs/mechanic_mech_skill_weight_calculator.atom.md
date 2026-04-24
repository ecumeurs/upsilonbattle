---
id: mechanic_mech_skill_weight_calculator
status: DRAFT
human_name: Skill Weight Calculator Mechanic
layer: IMPLEMENTATION
tags: [skills, balance, mathematics]
parents: []
type: MECHANIC
version: 2.0
priority: 5
dependents: []
---

# New Atom

## INTENT
To provide a mathematical framework for skill balance using the Skill Weight (SW) system where positive effects add weight and execution costs subtract weight.

## THE RULE / LOGIC
The Skill Weight Calculator implements the Net SW = 0 balance principle:

**Benefits Table (Adds SW):**
- Damage Multiplier: +10 SW per 10% (100% damage = +100 SW)
- Critical Chance: +2 SW per 1% (+25% crit = +50 SW)
- Critical Multiplier: +1 SW per 1% (+50% mult = +50 SW)
- Backstab Modifier: +30 SW flat
- Range Extension: +10 SW per cell > 1 (Range 3 = +20 SW)
- Zone/AoE: +50 SW per extra cell (3-cell line = +100 SW)
- Target Anywhere: +40 SW
- Stun Chance: +2 SW per 1% (50% stun = +100 SW)
- Poison Power: +15 SW per 1 dmg/turn (5 poison = +75 SW)
- Duration: +20 SW per extra turn
- Heal: +15 SW per 10% base
- Shield: +10 SW per point

**Payments Table (Subtracts SW):**
- Delay (+100): -100 SW (baseline)
- Extra Delay: -10 SW per +10 delay
- Channeling: -15 SW per 10 delay (risk premium)
- MP Leech: -15 SW per 1 MP
- SP Leech: -10 SW per 1 SP
- HP Leech: -20 SW per 1 HP (high risk)
- Cooldown: -25 SW per turn

**Balance Rule:** Net SW must equal 0 for a skill to be balanced.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_skill_weight_calculator]]`
- **Related Files:** `upsilonbattle/battlearena/entity/skill/skill.go`

## EXPECTATION
