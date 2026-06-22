---
id: mechanic_effect_heal
status: DRAFT
type: MECHANIC
priority: 3
human_name: Heal Effect Application
layer: IMPLEMENTATION
version: 1.0
tags: [combat, heal, iss-095]
parents:
  - [[mech_combat_attack_computation]]
dependents: []
---

# Heal Effect Application

## INTENT
To define how HP restoration is applied to a target entity, including overheal prevention.

## THE RULE / LOGIC
**Heal Resolution:**
1. Read `Skill.Heal` from the effect.
2. `actualHeal = min(Heal, Target.MaxHP - Target.HP)` — prevents overheal.
3. `Target.HP += actualHeal`.
4. Credits awarded = `actualHeal`.

**Note:** An effect is classified as healing when `IsHealing()` returns true (requires `Heal > 0`, `ShieldPower > 0`, or negative status power). Healing and damaging branches are NOT mutually exclusive — both can run on the same effect.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_effect_heal]]`
- **Location:** `effectapplicator.go:220-261`
- **Test Names:** `TestEffectApplicatorHeal`, `TestEffectApplicatorOverheal`, `TestHealAndShield_Combined`

## EXPECTATION
- Heal=5 on target with 5/10 HP → HP becomes 10, actualHeal=5.
- Heal=15 on target with 5/10 HP → HP becomes 10, actualHeal=5 (capped).
- Heal=5 + ShieldPower=3 → HP restored AND shield applied.
