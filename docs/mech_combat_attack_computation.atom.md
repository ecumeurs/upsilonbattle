---
id: mech_combat_attack_computation
human_name: "Combat Attack Computation"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: REVIEW
priority: 5
tags: [combat, damage, math]
parents:
  - [[module_backend_combat_math]]
dependents:
  - [[mech_combat_shielding]]
  - [[mechanic_effect_critical]]
  - [[mechanic_effect_damage]]
  - [[mechanic_effect_heal]]
  - [[mechanic_effect_poison]]
  - [[mechanic_effect_shield]]
  - [[mechanic_effect_stun]]
---

# Combat Attack Computation

## INTENT
To define the mathematical sequence for damage resolution in the core engine.

## THE RULE / LOGIC
This mechanic defines the mathematical sequence for damage resolution in the core engine, covering both skill-based ("Damaging" tag) and standard (non-skill) attacks.

**A. Standard (non-skill) attack** — linear reduction model with a floor of 1:
1. **Total Attack:** `TotalAttack = Attacker.Attack + Attacker.WeaponBaseDamage`
2. **Effective Defense:** `EffectiveDefense = Target.Defense + Target.ArmorRating`
3. **Multiplier:** `Multiplier = 1.0` for a normal hit. If the attack qualifies as a backstab, `Multiplier` and `EffectiveDefense` are both adjusted instead — the detection rule and the resulting 1.5x damage / 50% armor-penetration values are owned by [[mechanic_backstab_detection_algorithm]] and [[mechanic_armor_penetration_system]] respectively and are not restated here.
4. **Direct Computation:** `Damage = int(TotalAttack * Multiplier) - EffectiveDefense`
5. **Minimum Floor:** `FinalDamage = max(1, Damage)`
6. **Shield Absorption:** `FinalDamage` then passes through the shield-absorption rule owned by [[mech_combat_shielding]] before it reduces `Target.HP` — not restated here.
7. **Resolution:** `Target.HP = Target.HP - FinalDamage` (the post-shield remainder from step 6).

**B. Skill attack (Skill has a "Damaging" tag)** — "Three-Tunnel" model before shielding:
1. **Hit Test (Skills Only):** Accuracy vs Dodge roll.
2. **Mitigation Tunnels:**
    - **Physical:** `Phys = max((Attacker.Attack * Skill.Damage / 100) - Target.Defense - Target.Armor, 1)`
    - **Poison:** `Pois = max(Skill.PoisonPower - Target.Defense, 0)`
    - **Stun:** `Stun = max(Skill.StunPower - Target.Armor, 0)`
3. **Grand Total:** `TrueDamage = Phys + Pois + Stun`
4. **Crit Step:** if the hit crits, `FinalDmg = floor(TrueDamage * CritMultiplier)`
5. **Shield Step:** `Shield` absorbs `FinalDmg` 1:1.
6. **Resolution:** remaining damage reduces `Target.HP`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_combat_attack_computation]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `effectapplicator.go`
- **Test Names:** `rules_attack_test.go`, `rules_attack_failure_test.go`, `e2e_melee_attack_damage.js`
- **Integration:** Backstab detection and its 1.5x damage / 50% armor-penetration modifiers are owned by [[mechanic_backstab_detection_algorithm]] and [[mechanic_armor_penetration_system]]; shield absorption (both the standard and skill attack paths) is owned by [[mech_combat_shielding]].

## EXPECTATION (For Testing)
**Standard attack:**
- Attacker 10 Attack vs Target 5 Defense → 5 damage.
- Attacker 10 Attack vs Target 15 Defense → 1 damage (floor).

**Skill attack:**
- Attacking with 10 Phys and 5 Poison vs 5 Defense → Phys becomes 5, Poison becomes 0, total 5 damage applied to Shield then HP.
- Each mitigation tunnel is floored independently (Physical floor 1, Poison/Stun floor 0) before summing.
- A critical hit multiplies the summed TrueDamage by CritMultiplier (floored) before shield absorption.
- Shield absorbs final damage 1:1; only the overflow reduces HP.
