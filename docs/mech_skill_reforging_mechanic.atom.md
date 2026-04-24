---
id: mech_skill_reforging_mechanic
human_name: Skill Reforging Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [skills, progression, modification]
parents:
  - [[mech_skill_selection_progression]]
dependents: []
---

# Skill Reforging Mechanic

## INTENT
To implement the skill reforging system where players can modify existing skills every 5 levels, changing properties and effects for credits, with limitations based on skill grade and character level access.

## THE RULE / LOGIC
**Reforging Triggers:**
- **Frequency:** Available every 5 character levels
- **Cost:** Credits based on magnitude of changes and grade modifications
- **Limitations:** Cannot increase skill grade beyond current level access

**Reforging Operations:**
- **Property Modification:** Change skill property values (increase damage, extend range, etc.)
- **Effect Changes:** Replace or modify skill effects
- **Cost Adjustment:** Modify skill costs (delay, cooldown, resource costs)
- **Grade Changes:** Limited to within current level access range

**Reforging Cost Formula:**
- **Property Changes:** Credit cost = (Property SW difference × 2)
- **Grade Increase:** Credit cost = (New Grade SW - Current Grade SW) × 2
- **Complex Changes:** Additional credit cost for major reforgs

**Reforging Restrictions:**
- **Level 1-9:** Cannot exceed Grade II
- **Level 10-19:** Cannot exceed Grade III
- **Level 20-29:** Cannot exceed Grade IV
- **Level 30+: Cannot exceed Grade V

**Reforging Process:**
1. Player selects skill to reforge
2. System shows available modifications and costs
3. Player confirms changes
4. Credits deducted, skill updated
5. Old skill configuration lost (no undo)

## TECHNICAL INTERFACE (The Bridge
- **Code Tag:** `@spec-link [[mech_skill_reforging_mechanic]]`
- **API Endpoints:** `POST /api/v1/character/{id}/skill-reforge`, `GET /api/v1/skill/{id}/reforge-options`
