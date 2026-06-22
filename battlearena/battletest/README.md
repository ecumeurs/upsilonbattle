# `battletest` — Skill Scenario Sandbox

A deterministic, service-free harness for driving the battle engine in Go unit tests.

> **Spec atom:** [[module_skill_sandbox]] · **Related:** [[mech_movement_reposition]]

---

## 1. Why this exists

Asserting skill, movement, and positional-effect mechanics through a *full* battle is painful:
spinning up a `Ruler` actor, registering controllers over the message protocol, waiting on the
turn/shot clock, and reading state back asynchronously — all to check that "a kick pushes the
target two tiles." That machinery is the right thing to test in the `ruler` package, but it is
the wrong tool for verifying a single mechanic.

`battletest` formalizes the ad-hoc per-package helpers (`makeGameStateForTwo*`, hand-rolled
`FakeController`s) into one fluent builder over `gamestate.GameState`. It drives the **rules
layer directly** — `rules.UseSkill`, `rules.Move`, `rules.BeginingOfTurn`, positional effects —
with no actor, no async dispatch, no network, and no AI.

### Test altitudes

| Layer | What it verifies | Tool |
|-------|------------------|------|
| **Integration** | actor lifecycle, message routing, turn/shot clock, match end, win detection, races | `ruler` package tests (`NewRuler`/`Start`/`SendActor`) |
| **Mechanics** | damage, targeting, costs, effects, traps, reposition | **`battletest`** |
| **Pure unit** | a single function in isolation (e.g. damage formula) | `effectapplicator` tests |

`battletest` is for the middle row. It is **not** a replacement for the integration suite: it
deliberately omits the orchestration layer those tests exist to cover.

---

## 2. Quick start

```go
package mypkg_test

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/battletest"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
)

func TestStrike(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)

	attacker := s.Place(1, position.New(5, 5, 3), battletest.WithStat(property.Attack, 5))
	foe := s.Place(2, position.New(5, 6, 3), battletest.WithHP(20))

	strike := attacker.Give(
		battletest.NewSkill("Strike").TargetType("Entity").Range(1, 1).Damage(100).Build(),
	)

	s.Turn(attacker)
	reply, damaged, _ := s.UseSkill(attacker, strike, s.Pos(foe))

	if reply.HasError {
		t.Fatalf("unexpected error: %s", reply.ErrorKey)
	}
	if s.HP(foe) != 15 {
		t.Errorf("expected foe HP 15, got %d", s.HP(foe))
	}
	_ = damaged
}
```

The shape is always **setup → drive → assert**.

---

## 3. Constraint: the import cycle

`battletest` imports the `rules` package. Therefore **a test in `package rules` cannot import
`battletest`** — that would create an import cycle in the test binary and Go will reject it.

- ✅ New tests live in their **own package** (e.g. `battletest_test`, or any `*_test` external
  package, or a feature package's tests).
- ✅ Tests in the `ruler` package or other downstream packages may import `battletest`.
- ❌ In-package `rules` tests (`package rules`) must keep using the local `makeGameStateForTwo*`
  helpers. To move one onto `battletest`, first convert it to an **external** `package rules_test`.

---

## 4. API reference

### 4.1 Constructing a scenario

```go
func New(t testing.TB, x, y, z int) *Scenario
```
Creates a `GameState` backed by a grid of `x × y × z` cells. Pass the test's `*testing.T`; the
sandbox uses it for `Fatalf` on setup failures.

The underlying state is exposed as `s.GS` (`*gamestate.GameState`) as an escape hatch for
assertions or setups the helpers don't cover.

### 4.2 Controllers & entities

```go
func (s *Scenario) Controller(team int) *FakeController
func (s *Scenario) Place(team int, pos position.Position, opts ...PlaceOption) *Actor
```
`Place` builds a `Character` entity for `team`, loads the default character properties, registers
it on the grid and the turner, and assigns it the team's controller (created on first use, one
per team). It returns an `*Actor` handle (a thin wrapper around the entity UUID).

Each `Place` gets a unique, increasing `CurrentDelay`, so turn order is stable.

#### Place options

| Option | Effect |
|--------|--------|
| `WithHP(v)` | sets current **and** max HP |
| `WithMovement(v)` | sets current and max movement points |
| `WithMP(v)` / `WithSP(v)` | sets mana / skill points |
| `WithStat(p, v)` | sets any integer property (`property.Attack`, `property.Defense`, `property.ArmorRating`, `property.JumpHeight`, …) |

```go
caster := s.Place(1, position.New(5, 5, 3),
	battletest.WithHP(30),
	battletest.WithMovement(4),
	battletest.WithStat(property.Attack, 8),
)
```

### 4.3 Skills

```go
func (a *Actor) Give(sk skill.Skill) uuid.UUID   // registers a skill, returns its ID
```

Build skills fluently with `SkillSpec` (starts from `skill.New()`: Direct behavior, range 1,
Single zone, Entity target):

```go
func NewSkill(name string) *SkillSpec

func (b *SkillSpec) Behavior(bt def.BehaviorType) *SkillSpec   // Direct, Trap, Passive, …
func (b *SkillSpec) TargetType(tt def.TargetTypes) *SkillSpec  // Entity, Self, Tile, EnemyOnly, …
func (b *SkillSpec) Range(min, max int) *SkillSpec
func (b *SkillSpec) Zone(pattern string) *SkillSpec            // "Neighbours", "Line:3", "Circle:2", …

func (b *SkillSpec) Cost(p property.Property) *SkillSpec       // arbitrary cost
func (b *SkillSpec) MvtCost(v int) *SkillSpec

func (b *SkillSpec) Effect(p property.Property) *SkillSpec     // arbitrary effect property
func (b *SkillSpec) Damage(v int) *SkillSpec
func (b *SkillSpec) Heal(v int) *SkillSpec
func (b *SkillSpec) ShieldPower(v int) *SkillSpec
func (b *SkillSpec) Reposition(subject def.RepositionSubjectType, dist int) *SkillSpec

func (b *SkillSpec) Build() skill.Skill
```

For anything not covered by a named setter, use `.Effect(...)` / `.Cost(...)` with a
`defaultproperty.Make*` value.

### 4.4 Positional effects (traps)

```go
func (s *Scenario) Trap(pos position.Position, eff effect.Effect) uuid.UUID

// canonical observable: guaranteed poison, optional self-consumption
func PoisonTrap(power int, trigger property.TriggerTypeValue, removeOnTrigger bool) effect.Effect
```
`PoisonTrap` is the recommended observable for trigger assertions because it is deterministic
(100% chance) and visible via `s.Poison(...)`. For other effects, build an `effect.Effect`
directly and ensure it carries a `def.TriggerType(...)` property (a positional effect without a
`TriggerType` will **panic** by design — Fail-Fast).

### 4.5 Driving the engine

```go
func (s *Scenario) Turn(a *Actor)                  // force turn to actor, clear acted/moved
func (s *Scenario) UseSkill(caster *Actor, skillID uuid.UUID, target position.Position)
	(*message.Message, []rulermethods.ControllerAttacked, []rulermethods.ControllerSkillUsed)
func (s *Scenario) Move(a *Actor, path ...position.Position) *message.Message
func (s *Scenario) BeginTurn(a *Actor)             // fires start-of-turn (incl. OnTurn traps)
```
`UseSkill`/`Move` build the request message and call the real `rules` functions. Inspect the
returned `*message.Message`'s `HasError` / `ErrorKey` for validation outcomes.

### 4.6 Inspectors

```go
func (s *Scenario) Entity(a *Actor) entity.Entity   // raw snapshot
func (s *Scenario) Alive(a *Actor) bool             // still in the game state?
func (s *Scenario) Pos(a *Actor) position.Position
func (s *Scenario) HP(a *Actor) int
func (s *Scenario) Movement(a *Actor) int
func (s *Scenario) Shield(a *Actor) int
func (s *Scenario) Poison(a *Actor) int
func (s *Scenario) Stat(a *Actor, p interface{}) int
func (s *Scenario) HasActed(a *Actor) bool
func (s *Scenario) TrapAt(pos position.Position) bool   // any positional effect at pos?
```

---

## 5. Determinism

Hit rolls, crit rolls, and status-chance rolls go through `tools.RandomInt`. For repeatable
results:

- Use **100% / 0%** chances where you want certainty (e.g. `PoisonTrap` is 100%; default skill
  accuracy is 100, dodge 0 → always hits; crit chance 0 → never crits). Most mechanic tests need
  no seeding at all.
- When you must control a roll, install a deterministic generator with
  `tools.TesterRand(func(n int) int { return ... })` before the action (and restore it after).

---

## 6. Worked example — the fly-over trap rule

Movement skills reposition a subject and fire **only the landing tile's** triggers; flown-over
tiles are skipped. The sandbox makes this directly assertable:

```go
func TestDashOverTrap(t *testing.T) {
	s := battletest.New(t, 10, 10, 3)
	caster := s.Place(1, position.New(5, 5, 3), battletest.WithMovement(3))

	// non-removing traps on the flown-over tiles
	s.Trap(position.New(5, 6, 3), battletest.PoisonTrap(10, property.TriggerOnEnter, false))
	s.Trap(position.New(5, 7, 3), battletest.PoisonTrap(10, property.TriggerOnEnter, false))

	dash := caster.Give(battletest.NewSkill("Dash").
		TargetType(def.TargetTypeTile).Range(1, 3).
		Reposition(def.RepositionSubjectSelf, 3).Build())

	s.Turn(caster)
	s.UseSkill(caster, dash, position.New(5, 8, 3)) // aim = direction

	if !s.Pos(caster).Equals(position.New(5, 8, 3)) { t.Error("should have landed on the 3rd tile") }
	if s.Poison(caster) != 0 { t.Error("should have flown over the traps unharmed") }
}
```

See `reposition_test.go` for the full trap matrix (dash/push over & onto, pull, kick, retreat +
buff, MvtCost, blocked/out-of-grid landings) and `scenario_test.go` for basic skill/trap smoke
tests.

---

## 7. What this sandbox does **not** do

- **No orchestration.** No `Ruler` actor, message dispatch, turn/shot clock, match-end, or
  win-condition detection. Test those in the `ruler` package.
- **No AI / controllers-as-players.** Controllers are inert `FakeController`s that only record
  the notifications they receive.
- **No network / serialization.** State is manipulated in-process.

If a test needs any of the above, it belongs in the integration suite, not here.

---

## 8. Files

| File | Contents |
|------|----------|
| `scenario.go` | `Scenario`, `New`, `Place` + options, `Trap`, `Turn`, drivers |
| `inspect.go` | state inspectors |
| `builders.go` | `SkillSpec` fluent skill builder, `PoisonTrap` |
| `fakecontroller.go` | recording `actor.Communication` stub |
| `scenario_test.go` | sandbox smoke tests |
| `reposition_test.go` | movement-skill trap matrix |
