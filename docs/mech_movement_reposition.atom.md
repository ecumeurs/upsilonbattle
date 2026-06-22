---
id: mech_movement_reposition
human_name: Movement Skill Reposition Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 4
tags: [skills, movement, traps, positions]
parents:
  - [[contract_battle_contract]]
dependents: []
---

# Movement Skill Reposition Mechanic

## INTENT

To let skills displace a subject — the caster (dash/teleport) or the target (push/pull/kick) —
along the casting ray, with the defining trait that **tiles flown over do not fire positional
effects; only the landing tile does**. This enables skills to jump a target over a trap, or a
caster to escape over one, and composes with other effect properties (e.g. "retreat with a
defensive buff").

## THE RULE / LOGIC

**Model (upsilontypes):**

A reposition is expressed by two effect properties carried in `Skill.Effect.Properties`:

- `RepositionSubject` — `"Self"` (caster moves) or `"Target"` (targeted entity moves). Absence defaults to Self.
- `RepositionDistance` — signed tile count along the casting ray. Positive moves along the ray
  (dash forward / push away from caster); negative moves against it (pull toward caster).
  `0` (or absence) means the skill does not reposition.

A skill is a movement skill iff `RepositionDistance != 0`. The classifier tags any such skill
`mobility` (whenever a position changes).

**Direction (the casting ray):**

- **Self**: forward = unit step from caster toward the aim target (`unitStep(casterPos, target)`).
  The aim target (TargetType `Tile`/`EntityOrTile`) indicates direction only; the caster jumps
  `RepositionDistance` tiles that way. Self skills apply their co-located effects (buffs) to the
  caster itself.
- **Target**: forward = unit step from caster toward the subject (`unitStep(casterPos, subjectPos)`).
  Each targeted entity is displaced `RepositionDistance` tiles (negative = pulled toward caster).

`unitStep` takes the per-axis sign on the XY plane, so both cardinal and diagonal directions are
supported. Height (Z) is preserved.

**Validation (pre-application, in `preSkillChecks` → `checkReposition`):**

For every subject, the landing tile (`origin + forward * distance`) must be:
1. In-grid — else `skill.reposition.outofgrid`.
2. Walkable (Ground/Dirt) and not occupied by a blocking entity — else `skill.reposition.blocked`.
3. Have a non-zero direction — else `skill.reposition.nodirection`.

A failed reposition aborts the skill cleanly with no state change.

**Execution (after effects resolve, in `UseSkill` → `applyReposition`):**

1. Effects (damage/heal/buff) are applied first via `ApplyDirectEffect`.
2. Each surviving subject is relocated on the grid to its landing tile.
3. **Only the landing tile fires** `ProcessPositionalEffects(..., TriggerOnEnter)`. There is no
   `OnExit` on the origin and no `OnEnter`/`OnStep` on traversed tiles. This is the deliberate
   divergence from `Move()`, which fires triggers per step.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[mech_movement_reposition]]`
- **Model Files:** `upsilontypes/property/propertyenum.go`, `upsilontypes/property/def/skill.go`,
  `upsilontypes/entity/skill/skillgenerator/classifier.go`
- **Execution Files:** `upsilonbattle/battlearena/ruler/rules/reposition.go`,
  `upsilonbattle/battlearena/ruler/rules/skill.go`,
  `upsilonbattle/battlearena/ruler/rules/skill_validation.go`
- **Integration:** Reuses `mech_positional_effects` / `mech_trigger_system` for landing triggers.
- **Tests:** `@test-link [[mech_movement_reposition]]` — see the reposition trap matrix in
  `upsilonbattle/battlearena/battletest/reposition_test.go`.

## EXPECTATION

- Dash/push over an intermediate trap → trap does not fire; landing on a trap → it fires.
- Pull moves a target toward the caster; kicks (`O X`) push an adjacent target over one tile.
- Movement composes with co-located effects (e.g. retreat + shield) and consumes `MvtCost`.
- Out of scope: re-enabling the generator's `mobility` producer; non-landing trigger types.
