---
id: mec_backstabbing_mechanic
human_name: Backstabbing Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: STABLE
priority: 5
tags: [combat, positioning, damage]
parents:
  - [[mechanic_backstab_detection_algorithm]]
dependents: []
---

# Backstabbing Mechanic

## INTENT
To implement backstabbing combat mechanic with 150% damage multiplier and 50% armor penetration for attacks from behind, applicable to weapon attacks only.

## THE RULE / LOGIC
**Backstab Detection:**
- **Condition:** Attack from behind target (opposite orientation)
- **Detection Algorithm:** Attacker must be within 45° of target's back angle
- **Orientation Check:** Uses existing EntityOrientation system (Up, Right, Down, Left)

**Damage Calculation:**
- **Base Multiplier:** 150% damage (1.5× base damage)
- **Armor Penetration:** Ignores 50% armor rating
- **Shield Application:** Shield still applies fully (not penetrated)
- **Formula:** `(BaseDamage × 1.5) - (ArmorRating × 0.5) - Shield`

**Scope Limitations:**
- **Weapon Attacks Only:** Skills do not benefit from backstabbing in V2.1
- **All Weapons:** Bows, pistols, melee weapons all support backstabbing
- **No Skill Synergy:** Future skills may have "Backstab Enabled" property

**Visual Feedback:**
- Backstab indicators when positioning behind enemy
- Damage multiplier display in combat log
- Critical backstab highlighting for special effects

**AI Integration:**
- AI attempts to avoid backstabs when possible
- Sneak archetype AI prioritizes backstabbing opportunities
- Support AI positions to protect vulnerable allies

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_backstabbing_mechanic]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entity.go`, `upsilonbattle/battlearena/ruler/rules/attack.go`
