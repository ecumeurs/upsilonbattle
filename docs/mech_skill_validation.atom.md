---
id: mech_skill_validation
human_name: Skill Execution & Validation Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[req_skill_generation]]
dependents: []
---
# Skill Execution & Validation Mechanic

## INTENT
To validate a skill-use request against all preconditions before execution, rejecting it with a specific error code if any check fails.

## THE RULE / LOGIC
**Skill Execution & Validation Mechanic.**

A skill-use request must pass every verification below before the skill executes; the first failing check rejects the request with the indicated error code:

1. **Existence Verification** — both the Entity and the requested Skill ID must exist on the roster (`skill.notfound`).
2. **Turn / Controller Identity Verification** — the entity must be on its active turn and owned by the controller issuing the command (`entity.turn.mismatch`, `entity.controller.mismatch`).
3. **Action State Verification** — the entity must not have already acted for this skill (`entity.alreadyacted`).
4. **Economic Cost — Cooldown Check** — the skill's `Cooldown` property must be 0 (`skill.cooldown`).
5. **Economic Cost — Stat Leech** — the entity must hold enough points in each required stat to pay the action cost (`skill.cost.mp`, `skill.cost.sp`, `skill.cost.hp`, `skill.cost.mvt`).
6. **Grid Boundaries Verification** — the target coordinate must be within the mapped Grid (`skill.target.outofgrid`).
7. **Range Limit Verification** — the mathematical distance to the target must fall between the skill's `MinRange` and `MaxRange` (`skill.target.range`).
8. **Entity Targeting Rules Verification** — the target must satisfy the skill's explicit `TargetType` (`skill.target.none`, `skill.target.self`, `skill.target.enemyonly`, `skill.target.friendonly`, `skill.target.tile`).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_skill_validation]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/skill_validation.go`
- **Test Names:** `rules_skill_validation_test.go`, `rules_skill_cooldown_test.go`, `rules_skill_leech_mp_sp_test.go`, `rules_skill_leech_hp_mvt_test.go`, `rules_skill_targeting_test.go`, `aoe_skill_test.go`
