---
id: mech_combat_shielding
human_name: "Combat Shielding Mechanic"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 5
tags: [combat, shield, defense]
parents:
  - [[mech_combat_attack_computation]]
dependents: []
---

# Combat Shielding Mechanic

## INTENT
To define Shielding as a secondary health buffer that absorbs damage after personal mitigation.

## THE RULE / LOGIC
1.  **Post-Mitigation**: Unlike the "Bubble" concept, these shields absorb `TrueDamage` *after* the target's Armor and Defense have already reduced the incoming hit.
2.  **Absorption**: Absorbs damage at a 1:1 ratio.
3.  **Persistence**: Shields persist between turns until depleted.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_combat_shielding]]`
- **Related Property:** `property.Shield`

## EXPECTATION (For Testing)
- If a character has 10 Defense and 10 Shield, and takes 15 raw damage:
    - 5 damage proceeds past Defense.
    - 5 damage is absorbed by the Shield (Shield becomes 5).
    - HP remains unchanged.
