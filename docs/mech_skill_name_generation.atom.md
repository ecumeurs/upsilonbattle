---
id: mech_skill_name_generation
human_name: Skill Name Generation Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 4
tags: [skills, names, generation, aesthetic]
parents:
  - [[upsilonbattleui:req_ui_look_and_feel]]
  - [[mech_ai_name_generation]]
  - [[upsilontypes:module_skill_generator]]
dependents: []
---

# Skill Name Generation Mechanic

## INTENT
To produce diegetic, "Neon in the Dust" flavoured names for procedurally generated skills by combining a tag-driven modifier prefix, a primary-tag subject, and a grade-coloured suffix. Names must never equal a raw property key (`Damage`, `Heal`, `Shield`).

## THE RULE / LOGIC
Names follow the template: `[Modifier] [Subject] [Suffix]` with a hard cap of **24 characters**.

**Modifier prefix** — driven by the first secondary tag that matches:

| Secondary tag | Prefix pool |
|---|---|
| `dot` | Cinder, Sludge, Rot\_, Verm\_ |
| `crit` | Razor, Fang\_, Spike\_ |
| `aoe` | Flux, Cascade\_, Wave\_ |
| `channeled` | Drift\_, Bleed\_, Slow\_ |
| `buff` / `debuff` | Echo\_, Static, Hex\_ |
| none | ∅, Null\_, Void\_, Ghost\_ |

**Subject** — driven by the primary tag:

| Primary tag | Subject pool |
|---|---|
| `melee` | Strike, Bash, Cleaver, Smash |
| `ranged` | Bolt, Lance, Pulse, Tracer |
| `aoe` | Burst, Field, Storm, Bloom |
| `heal` | Mend, Patch, Pulse, Suture |
| `shield` | Bulwark, Aegis, Plate, Shell |
| `trap` | Mine, Snare, Tripwire, Hex |
| `counter` | Riposte, Ricochet, Rebuke |
| `reaction` | Reflex, Backlash |
| `passive` | Aura, Cycle, Drift |
| `stun` | Jolt, Stutter, Lockdown |
| `mobility` | Sprint, Phase, Vector |

**Suffix** — 30% grade-flavored (`_I` … `_V`), 70% haxxor pool (`_X`, `v2`, `_Z`, `_Bot`, `_666`, `_Alpha`).

The full segment is truncated to 24 chars if the combination exceeds the limit.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_skill_name_generation]]`
- **File:** `upsilontypes/entity/skill/skillgenerator/namegen.go`
- **Signature:** `func Name(primaryTag string, secondaryTags []string, grade string) string`
- **Called by:** `skillgenerator.Generate()` after producing the skill, before return.

## EXPECTATION
- Generated name does not equal any raw `SkillProperties` key.
- Length ≤ 24 characters.
- Subject segment is always non-empty (fallback: `"Skill"`).
- 100% of names match pattern: at least one alphabetic word followed by an optional suffix token.
