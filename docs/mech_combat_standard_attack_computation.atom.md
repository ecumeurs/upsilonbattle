---
id: mech_combat_standard_attack_computation
human_name: "Standard Combat Attack Computation"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 5
tags: [combat, damage, math]
parents:
  - [[module_backend_combat_math]]
dependents: []
---

# Standard Combat Attack Computation

## INTENT
To define the simplest mathematical sequence for standard (non-skill) physical attacks.

## THE RULE / LOGIC
Standard attacks follow a linear reduction model with a floor of 1:

1.  **Direct Computation**:
    - `Damage = Attacker.Attack - Target.Defense`
2.  **Minimum Floor**:
    - `FinalDamage = max(1, Damage)`
3.  **Resolution**:
    - `Target.HP = Target.HP - FinalDamage`

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_combat_standard_attack_computation]]`
- **Related Files:** `attack.go`

## EXPECTATION (For Testing)
- If Attacker has 10 Attack and Target has 5 Defense:
    - Result: 5 Damage.
- If Attacker has 10 Attack and Target has 15 Defense:
    - Result: 1 Damage (floor).
