---
id: mechanic_ai_controller_archetypes
status: DRAFT
version: 1.0
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To implement individual AI controller archetypes for Fighter, Ranger, Support, and Sneak behaviors, defining specific decision trees, skill priorities, and tactical patterns for each archetype.

## THE RULE / LOGIC
**AI Controller Archetypes:**

**Core Principle:**
Each AI archetype extends the base controller with specialized behavior patterns, skill selection priorities, and tactical positioning rules.

**Fighter Controller Archetype:**
- **Tactical Intent:** Direct aggression and frontline engagement.
- **Positioning:** Prioritizes cells adjacent to enemies or within immediate melee range.
- **Skill Priority:** High-damage melee skills, self-shields, and gap-closing "charge" abilities.
- **Decision Loop:**
    1. Check for killable enemies within reach and execute a killing blow.
    2. If no immediate kill is possible, use a charge skill to close distance.
    3. Default to attacking the nearest enemy with the most efficient melee skill.
- **Stat Distribution:** Focuses on Attack (60%) and Defense (25%), with minor HP and Movement.

**Ranger Controller Archetype:**
- **Tactical Intent:** Ranged harassment and kiting.
- **Positioning:** Maintains an optimal distance (4-7 cells) from enemies; favors high ground or cover.
- **Skill Priority:** Ranged attacks, area-denial traps, and mobility skills (dashes/teleports).
- **Decision Loop:**
    1. If outside optimal range, move to a safe distance.
    2. Place traps at chokepoints if available.
    3. Attack the most exposed enemy from range.
- **Stat Distribution:** Focuses on Attack (50%) and Accuracy (20%), with moderate Movement.

**Support Controller Archetype:**
- **Tactical Intent:** Ally preservation and combat enhancement.
- **Positioning:** Stays within proximity of teammates, blocking paths to vulnerable allies.
- **Skill Priority:** HP restoration, damage-absorption shields, and offensive/defensive buffs.
- **Decision Loop:**
    1. Heal the most damaged ally if below a critical threshold.
    2. Apply shields to allies threatened by high-damage enemies.
    3. Buff teammates lacking active enhancements.
    4. Move to a position that maximizes support coverage.
- **Stat Distribution:** Focuses on MP/SP pools (40%) and Defense (25%), with zero Attack.

**Sneak Controller Archetype:**
- **Tactical Intent:** Flanking and precision status application.
- **Positioning:** Aggressively seeks cells behind enemies to trigger backstab bonuses.
- **Skill Priority:** High-crit backstab skills, stealth, evasion buffs, and poison/stun effects.
- **Decision Loop:**
    1. Attempt to move behind a vulnerable target for a backstab attack.
    2. Apply status effects (poison/stun) to high-value targets.
    3. Use evasion or stealth if under direct threat.
- **Stat Distribution:** Focuses on Movement (30%) and Attack (25%), with high Dodge and Crit.

**Team Composition Logic:**
- **Constraints:** Maximum of 1 Support and 1 Sneak per team. No limits on Fighters or Rangers.
- **Selection Priority:** Teams always start with 1 Fighter and 1 Ranger. Support is added at size 3, Sneak at size 4. Remaining slots are filled with Fighters.

**Performance Optimization:**
- **Turn-Based Caching:** Decision trees and tactical evaluations are cached per controller archetype for the duration of a turn to reduce redundant calculations.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ai_controller_archetypes]]`
- **Related Files:** `upsilonbattle/battlearena/controller/controllers/fighter.go`, `upsilonbattle/battlearena/controller/controllers/ranger.go`, etc.
- **Integration:** Works with `mec_ai_archetype_system`, `ai_progression_matching`

## EXPECTATION
