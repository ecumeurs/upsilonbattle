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
---

# Combat Attack Computation

## INTENT
To define the mathematical sequence for damage resolution in the core engine.

## THE RULE / LOGIC
This mechanic is only applicable for Skills that have a "Damaging" tag.
Standard attacks (no skill) are not using this mechanic.

Damage resolution follows a "Three-Tunnel" model before applying shielding:

1.  **Hit Test (Skills Only)**: Accuracy vs Dodge roll.
2.  **Mitigation Tunnels**:
    - **Physical**: `Phys = max((Attacker.Attack * Skill.Damage / 100) - Target.Defense - Target.Armor, 1)`
    - **Poison**: `Pois = max(Skill.PoisonPower - Target.Defense, 0)`
    - **Stun**: `Stun = max(Skill.StunPower - Target.Armor, 0)`
3.  **Grand Total**:
    - `TrueDamage = Phys + Pois + Stun`
4.  **Crit Step**:
    - If hit crits: `FinalDmg = floor(TrueDamage * CritMultiplier)`
5.  **Shield Step**:
    - `Shield` absorbs `FinalDmg` 1:1.
6.  **Resolution**:
    - Remaining damage reduces `Target.HP`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_combat_attack_computation]]`
- **Related Files:** `effectapplicator.go`

## EXPECTATION (For Testing)
- If attacking with 10 Phys, 5 Poison vs 5 Defense:
    - Phys becomes 5.
    - Poison becomes 0.
    - Total 5 damage to Shield/HP.
