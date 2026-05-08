---
id: mechanic_mech_effect_damage
status: DRAFT
parents:
  - [[mech_combat_attack_computation]]
dependents: []
human_name: Damage Effect Resolution
layer: IMPLEMENTATION
priority: 3
tags: [combat, damage, iss-095]
type: MECHANIC
version: 1.0
---

# New Atom

## INTENT
To define the physical damage tunnel computation from skill damage percentage and entity attack stat.

## THE RULE / LOGIC
**Physical Damage Tunnel:**
1. `rawPhysical = Attacker.Attack * Skill.Damage / 100` — Skill.Damage is a percentage multiplier.
2. `physicalDamage = max(rawPhysical - Target.Defense - Target.ArmorRating, 0)`.
3. If `physicalDamage + poisonDamage + stunDamage = 0`, no HP is deducted.

**Damage Default:** Skill.Damage defaults to 100 (100% of Attack). The blueprint explicitly initializes Damage to 0 for non-damaging skills.

**IsDamaging Check:** An effect enters the damage branch only if `IsDamaging()` returns true: Damage > 0, StunPower > 0, PoisonPower > 0, or ShieldPower < 0.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_effect_damage]]`
- **Location:** `effectapplicator.go:118-193`
- **Test Names:** `TestEffectApplicatorDamage1Target`, `TestEffectApplicatorDamage1Defense`, `TestEffectApplicatorDamage1TargetShield`

## EXPECTATION
- Attack=3, Damage=200 (200% multiplier), Defense=0 → rawDmg=6.
- Attack=3, Damage=200, ArmorRating=5 → rawDmg=max(6-5,0)=1.
- Damage=0, no other positive properties → IsDamaging()=false, no damage branch.
