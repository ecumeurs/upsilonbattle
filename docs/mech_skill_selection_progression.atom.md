---
id: mech_skill_selection_progression
human_name: Skill Selection Progression Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [skills, progression, selection]
parents:
  - [[req_skill_generation]]
dependents:
  - [[mech_skill_reforging_mechanic]]
---

# Skill Selection Progression Mechanic

## INTENT
To implement the skill selection system where a character receives a single grade-gated skill roll at creation (consuming a one-time roulette_used flag, not a choice of 3), and players choose 1 of 3 skills every 10 levels thereafter, with available skill grades increasing as characters progress.

## THE RULE / LOGIC
**Skill Selection Timeline:**
- **Character Creation:** A single grade-gated skill roll (Grade I-II) — not a choice of 3. The engine generates exactly one skill for the requested grade and it is acquired directly; the character's `roulette_used` flag (default false) is marked consumed by this roll. (Verified from code: the roll handler marks the flag but does not itself reject a subsequent call — whether repeat rolls are blocked elsewhere is not confirmed in this pass.)
- **Every 10 Levels:** Choose 1 of 3 random skills (higher grades available)
- **Skill Reforging:** Every 5 levels, can modify existing skill properties

**Skill Grade Progression:**
- **Level 1-9:** Grade I-II skills offered
- **Level 10-19:** Grade II-III skills offered
- **Level 20-29:** Grade III-IV skills offered
- **Level 30+:** Grade IV-V skills offered

**Selection Process (Every 10 Levels):**
1. System generates 3 random skills from appropriate grade pool
2. Player reviews skill properties, costs, and effects
3. Player selects 1 skill to learn
4. Skill added to character's skill list
5. Unselected skills are discarded

**Selection Process (Character Creation):**
1. System generates exactly 1 random skill from the grade-gated pool
2. The rolled skill is acquired directly onto the character — no review-and-pick-1-of-3 step
3. The character's `roulette_used` flag is marked consumed

**Skill Pool Design:**
- **Base Skills:** Predefined skill templates by grade
- **Procedural Skills:** Generated from property combinations
- **Skill Archetypes:** Offensive, Defensive, Utility, Movement
- **Availability:** Skills filtered by character level and grade

**Reforging Mechanics:**
- **Cost:** Credits based on grade change magnitude
- **Limitations:** Can't increase skill grade beyond current level access
- **Risk:** Reforging may change skill behavior significantly

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_selection_progression]]`
- **API Endpoints:** `POST /api/v1/profile/character/{characterId}/skills/roll` — shipped, drives the character-creation single grade-gated roll (grade selectable via `?grade=`, gated by the account's total wins). The every-10-levels choose-1-of-3 flow and `POST /api/v1/character/{id}/skill-reforge` are not confirmed against current code in this pass — do not treat their endpoint shapes as verified.
