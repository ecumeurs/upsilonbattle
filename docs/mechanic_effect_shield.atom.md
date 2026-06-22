---
id: mechanic_effect_shield
status: DRAFT
tags: [combat, shield, defense, iss-095]
dependents: []
human_name: Shield Effect Application
type: MECHANIC
priority: 3
parents:
  - [[mech_combat_attack_computation]]
layer: IMPLEMENTATION
version: 1.0
---

# Shield Effect Application

## INTENT
To define how shield (overshield) is granted to a target entity and how it caps at 2× max HP.

## THE RULE / LOGIC
**Shield Application:**
1. If `Skill.ShieldPower > 0` (healing branch):
   - `Target.Shield = min(Target.Shield + ShieldPower, Target.MaxHP * 2)` — capped at 2× max HP.
2. If `Skill.ShieldPower < 0` (damage branch):
   - Shield is depleted: `Target.Shield = max(Target.Shield + ShieldPower, 0)`.

**Shield Absorption (Damage Branch):**
1. After computing TrueDamage, shield absorbs 1:1.
2. `absorbed = min(Shield, TrueDamage)`.
3. `Shield -= absorbed`, `TrueDamage -= absorbed`.
4. Remaining TrueDamage reduces HP.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_effect_shield]]`
- **Location:** `effectapplicator.go:148-152, 167-188, 248-250`
- **Test Names:** `TestShieldOvershield_CappedAt2xMaxHP`, `TestEffectApplicatorDamage1TargetShield`, `TestEffectApplicatorShielding`

## EXPECTATION
- ShieldPower=100 on entity with MaxHP=10 → Shield capped at 20.
- Shield=5 absorbs 6 damage → 1 HP damage, Shield=0.
- ShieldPower=5 on target with Shield=0 → Shield=5.
