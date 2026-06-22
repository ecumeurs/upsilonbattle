---
id: rule_combat_range_validation
human_name: "Combat Range Validation Rule"
type: RULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: [combat, range, grid]
parents:
  - [[shared:us_take_combat_turn]]
dependents: []
---

# Combat Range Validation Rule

## INTENT
To define the mathematical validation for normal attacks across the 3D grid, ensuring intuitive melee and reaching combat.

## THE RULE / LOGIC
- **Base Distance:** The horizontal range check uses **2D Manhattan Distance** (`|Δx| + |Δy|`).
- **Horizontal Constraint:** `2D_Distance <= entity.AttackRange`.
- **Vertical Constraint:** The height difference (`|Δz|`) between the attacker and the target must be less than or equal to `1 + entity.AttackRange`.
- **Rationale:** This ensures that melee units (Range 1) can hit targets on adjacent cells with up to 2 units of height difference, reflecting their ability to jump/reach (default jump is 2).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_combat_range_validation]]`
- **Implementation:** `upsilonbattle/battlearena/ruler/rules/attack.go`

## EXPECTATION (For Testing)
- Attacker at (0,0,0) with Range 1 can hit target at (1,0,1) -> ACCEPTED (2D dist 1, Δz 1).
- Attacker at (0,0,0) with Range 1 can hit target at (1,0,2) -> ACCEPTED (2D dist 1, Δz 2).
- Attacker at (0,0,0) with Range 1 cannot hit target at (1,0,3) -> REJECTED (Δz 3 > 1+1).
- Attacker at (0,0,0) with Range 1 cannot hit target at (1,1,0) -> REJECTED (2D dist 2 > 1).
