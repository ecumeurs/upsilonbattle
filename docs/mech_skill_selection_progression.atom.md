---
id: mech_skill_selection_progression
human_name: Skill Selection Progression Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [skills, progression, selection]
parents: []
dependents:
  - [[mech_skill_reforging_mechanic]]
---

# Skill Selection Progression Mechanic

## INTENT
To implement the skill selection system where players choose 1 of 3 skills at character creation and every 10 levels, with available skill grades increasing as characters progress.

## THE RULE / LOGIC
**Skill Selection Timeline:**
- **Character Creation:** Choose 1 of 3 random skills (Grade I-II)
- **Every 10 Levels:** Choose 1 of 3 random skills (higher grades available)
- **Skill Reforging:** Every 5 levels, can modify existing skill properties

**Skill Grade Progression:**
- **Level 1-9:** Grade I-II skills offered
- **Level 10-19:** Grade II-III skills offered
- **Level 20-29:** Grade III-IV skills offered
- **Level 30+:** Grade IV-V skills offered

**Selection Process:**
1. System generates 3 random skills from appropriate grade pool
2. Player reviews skill properties, costs, and effects
3. Player selects 1 skill to learn
4. Skill added to character's skill list
5. Unselected skills are discarded

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
- **API Endpoints:** `POST /api/v1/character/{id}/skill-select`, `POST /api/v1/character/{id}/skill-reforge`
