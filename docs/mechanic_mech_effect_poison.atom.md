---
id: mechanic_mech_effect_poison
status: DRAFT
version: 1.0
parents:
  - [[mech_combat_attack_computation]]
human_name: Poison Effect Application
priority: 3
tags: [combat, status-effect, poison, dot, iss-095]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To define the conditions under which a poison status effect is applied to a target entity during combat damage resolution.

## THE RULE / LOGIC
**Poison Application Guard:**
1. Compute `truePoison = max(Skill.PoisonPower - Target.Defense, 0)`.
2. If `truePoison > 0` AND `RandomInt(0, 100) < Skill.PoisonChance`, apply poison.
3. **Both PoisonPower AND PoisonChance must be non-zero** for poison to ever proc (ISS-095 pairing rule).
4. Poison stacks additively: `Target.Poison += truePoison`.
5. PoisonPower also contributes to the "Poison Tunnel" of damage (adds to TrueDamage).

**Generator Requirement:** Any producer or secondary layer that adds PoisonPower MUST also add PoisonChance (and vice-versa). Missing either renders the skill mechanically dead.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_effect_poison]]`
- **Location:** `effectapplicator.go:154-158`
- **Test Names:** `TestPoisonApplied_WhenBothPowerAndChance`, `TestPoisonNotApplied_WhenChanceZero`, `TestPoisonNotApplied_WhenPowerZero`, `TestPoisonStacking`

## EXPECTATION
- PoisonPower=3, PoisonChance=100, no defense → poison=3 applied.
- PoisonPower=3, PoisonChance=0 → no poison applied.
- PoisonPower=0, PoisonChance=100 → no poison applied (IsDamaging=false).
- Two applications of PoisonPower=3, PoisonChance=100 → poison stacks to 6.
