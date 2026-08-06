---
id: module_skill_sandbox
human_name: Skill Scenario Sandbox
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 5
tags: [testing, sandbox, skills]
parents:
  - [[shared:req_tech_debt_backlog]]
dependents: []
---

# Skill Scenario Sandbox

## INTENT

To provide a deterministic, service-free harness for driving the battle engine in unit tests.
Asserting skill, movement, and positional-effect mechanics through a full battle (matchmaking,
AI, network) is slow and flaky. The sandbox formalizes the ad-hoc per-package helpers
(`makeGameStateForTwo*`, `FakeController`) into one fluent builder over `GameState`, so mechanics
are set up and asserted directly against the rules layer.

## THE RULE / LOGIC

The `battletest` package (`upsilonbattle/battlearena/battletest`) exposes a `Scenario` builder:

**Setup**
- `New(t, x, y, z)` — a GameState backed by a grid of the given dimensions.
- `Place(team, pos, opts...)` — a Character entity for a team; options set HP/MP/SP/Movement and
  stats (`WithHP`, `WithMovement`, `WithStat`, …). A `FakeController` is created per team.
- `Actor.Give(skill)` — register a skill (built fluently via `NewSkill(...)`).
- `Trap(pos, effect)` — register a positional effect; `PoisonTrap(...)` is the canonical observable.

**Drive**
- `Turn(actor)` — force the turn and clear acted/moved state.
- `UseSkill(caster, skillID, target)`, `Move(actor, path...)`, `BeginTurn(actor)` — invoke the
  exported rules-layer functions directly.

**Inspect**
- `Pos`, `HP`, `Movement`, `Shield`, `Poison`, `Stat`, `HasActed`, `Alive`, `TrapAt`.

Determinism comes from the default RNG (override with `tools.TesterRand` when needed). Combined
with `TrapAt` and the status inspectors, tests assert *which* tiles fired (e.g. fly-over vs
landing) without any production hooks.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[module_skill_sandbox]]`
- **Files:** `upsilonbattle/battlearena/battletest/{scenario,inspect,builders,fakecontroller}.go`
- **Used by:** the reposition trap matrix (`battletest/reposition_test.go`) and sandbox smoke
  tests (`battletest/scenario_test.go`).

## EXPECTATION

- A skill/trap/movement scenario can be expressed in a handful of fluent calls and asserted
  deterministically, with no services running.
