---
id: req_skill_generation
human_name: Skill Generation System Requirement
type: REQUIREMENT
layer: BUSINESS
version: 1.0
status: DRAFT
priority: 4
tags: [skills, generation, tags, grades, names, icons]
parents:
  - [[shared:req_trpg_game_definition]]
dependents:
  - [[upsilonbattleui:ui_action_panel]]
  - [[upsilonbattleui:ui_skill_icon]]
  - [[mech_skill_equipment_system]]
  - [[mech_skill_selection_progression]]
  - [[mech_skill_validation]]
  - [[upsilontypes:module_skill_generator]]
---

# Skill Generation System Requirement

## INTENT
To specify that the skill roll pipeline must produce grade-aware, category-tagged skills with diegetic names, so a player understands what a skill does from its tag icons and name without opening a tooltip.

## THE RULE / LOGIC
Every procedurally generated skill must satisfy four constraints:

**1. Grade targeting** — the generator accepts a target grade (I–V) and produces a skill whose positive skill weight falls within the grade's PSW band (see `[[rule_skill_grading_system]]`). Grade I is the default.

**2. Tag classification** — after generation, skills carry an ordered `tags []string` derived from their final properties (not authored). Tags drive icon composition, name prefix selection, and action-panel rendering. The 17-tag vocabulary covers behavior (`trap`, `counter`, `reaction`, `passive`), effect family (`heal`, `shield`, `dot`, `stun`, `buff`, `debuff`), delivery (`melee`, `ranged`, `aoe`), and modifiers (`crit`, `channeled`, `instant`, `mobility`).

**3. Diegetic naming** — skill names follow a `[Modifier] [Subject] [Suffix]` template capped at 24 characters. No name may equal a raw property key (`Damage`, `Heal`, `Shield`).

**4. Arena action panel** — equipped skills with `behavior ∈ {Direct, Reaction, Counter, Trap}` appear as clickable buttons in the Battle Arena action panel; skills with `behavior == Passive` render in a read-only passive rail below the active row.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[req_skill_generation]]`
- **Related Issue:** `ISS-065_skill_generation_balance`
- **Test:** `upsiloncli/tests/scenarios/e2e_skill_roll_naming.js`


## EXPECTATION
- Every rolled skill has `instance_data.name` that does not equal a raw property key.
- Every rolled skill has `instance_data.tags` as a non-empty array drawn from the 17-tag vocabulary.
- Grade I skills have PSW ∈ [60, 150].
- `ActionPanel.vue` renders active-skill buttons and a passive rail; passives have no click handler.
