---
id: mechanic_mech_effect_critical
status: DRAFT
parents:
  - [[mech_combat_attack_computation]]
dependents: []
human_name: Critical Hit Mechanic
version: 1.0
tags: [combat, critical, damage, iss-095]
type: MECHANIC
layer: IMPLEMENTATION
priority: 3
---

# New Atom

## INTENT
To define the critical hit roll and damage multiplier applied during combat damage resolution.

## THE RULE / LOGIC
**Critical Hit Resolution:**
1. If `Skill.CriticalChance > 0` AND `RandomInt(0, 100) < CriticalChance`:
   - `multiplier = CriticalMultiplier / 100.0` (default CritMultiplier = 100 → 1.0x, i.e., no bonus).
2. `TrueDamage = floor(TrueDamage * multiplier)`.
3. Critical applies AFTER all three damage tunnels (Physical + Poison + Stun) are summed.
4. Critical applies BEFORE shield absorption.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_effect_critical]]`
- **Location:** `effectapplicator.go:135-139`
- **Test Names:** `TestEffectApplicatorCrit`

## EXPECTATION
- CritChance=100, CritMultiplier=200, base damage=3 → final damage = 6.
- CritChance=0 → no multiplier applied (default 1.0x).
