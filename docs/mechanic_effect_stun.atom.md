---
id: mechanic_effect_stun
status: DRAFT
human_name: Stun Effect Application
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
tags: [combat, status-effect, stun, iss-095]
parents:
  - [[mech_combat_attack_computation]]
priority: 3
dependents: []
---

# Stun Effect Application

## INTENT
To define the conditions under which a stun status effect is applied to a target entity during combat damage resolution.

## THE RULE / LOGIC
**Stun Application Guard:**
1. Compute `trueStun = max(Skill.StunPower - Target.ArmorRating, 0)`.
2. If `trueStun > 0` AND `RandomInt(0, 100) < Skill.StunChance`, apply stun.
3. **Both StunPower AND StunChance must be non-zero** for stun to ever proc (ISS-095 pairing rule).
4. Stun stacks additively: `Target.Stun += trueStun`.
5. StunPower also contributes to the "Stun Tunnel" of damage (adds to TrueDamage).

**Generator Requirement:** Any producer or secondary layer that adds StunChance MUST also add StunPower (and vice-versa). Missing either renders the skill mechanically dead.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_effect_stun]]`
- **Location:** `effectapplicator.go:159-163`
- **Test Names:** `TestStunApplied_WhenBothPowerAndChance`, `TestStunNotApplied_WhenChanceZero`, `TestStunNotApplied_WhenPowerZero`, `TestStunStacking`

## EXPECTATION
- StunPower=5, StunChance=100, no armor → stun=5 applied.
- StunPower=5, StunChance=0 → no stun applied.
- StunPower=0, StunChance=100 → no stun applied (IsDamaging=false).
- Two applications of StunPower=2, StunChance=100 → stun stacks to 4.
