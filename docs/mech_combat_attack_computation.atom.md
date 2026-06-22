---
id: mech_combat_attack_computation
human_name: "Combat Attack Computation"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
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
1. **Direct Computation:** `Damage = Attacker.Attack - Target.Defense`
2. **Minimum Floor:** `FinalDamage = max(1, Damage)`
3. **Resolution:** `Target.HP = Target.HP - FinalDamage`

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
- **Test Names:** `rules_attack_test.go`, `rules_attack_failure_test.go`

## EXPECTATION (For Testing)
**Standard attack:**
- Attacker 10 Attack vs Target 5 Defense → 5 damage.
- Attacker 10 Attack vs Target 15 Defense → 1 damage (floor).

**Skill attack:**
- Attacking with 10 Phys and 5 Poison vs 5 Defense → Phys becomes 5, Poison becomes 0, total 5 damage applied to Shield then HP.
- Each mitigation tunnel is floored independently (Physical floor 1, Poison/Stun floor 0) before summing.
- A critical hit multiplies the summed TrueDamage by CritMultiplier (floored) before shield absorption.
- Shield absorbs final damage 1:1; only the overflow reduces HP.
