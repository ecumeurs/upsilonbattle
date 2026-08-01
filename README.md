# Upsilon Battle

Deterministic, actor-based tactical RPG combat engine — the `BattleArena` — powering turn-based battles for the Upsilon platform.

## What this is

`upsilonbattle` is the core combat simulation library of the Upsilon platform. It has no
network layer and no persistence of its own: it is a pure Go engine embedded by
[`upsilonapi`](https://github.com/ecumeurs/upsilonapi) (the "Bridge" service), which imports
`github.com/ecumeurs/upsilonbattle/battlearena` and exposes it over JSON APIs to real
clients (the web/battle UI, bots, CLI tooling).

A match is represented by a `BattleArena`: a self-contained actor (`Uuid`, `Controllers`,
`Ruler`) created via `battlearena.NewBattleArena`. Each arena owns exactly one `Ruler`, the
actor responsible for enforcing the rules and driving the match to completion.

### Core concepts

* **Ruler** (`battlearena/ruler`) — the match's authoritative actor. It owns the `GameState`
  (grid + entities) and is the only thing allowed to mutate it once the match has started; all
  interaction happens through its actor message queue, never through direct pointer access.
  It tracks the arena's lifecycle (`WaitingForControllers` → `InProgress` → `Finished`),
  validates every incoming action, and declares the winner.
* **Controllers** (`battlearena/controller`) — the seam between a Ruler and a player (human or
  AI). Each controller registers with the Ruler, is granted a squad of entities, and issues
  move/attack/skill commands on their behalf. AI archetypes and layered behaviors live under
  `battlearena/controller/archetype` and `battlearena/controller/behavior`.
* **Dynamic initiative** — rather than a fixed turn order, every entity accrues a delay credit
  (movement, attacks, and ending a turn all add credit); the entity with the lowest credit acts
  next. This produces a non-linear, tactically meaningful pace instead of a strict round-robin.
* **Grid-based tactics** (`battlearena/ruler/gamestate`, via `upsilonmapdata`/`upsilonmapmaker`)
  — matches take place on a procedurally generated board with terrain, movement range, jump
  height, and attack range/line validation.
* **Rules & mechanics** (`battlearena/ruler/rules`, `battlearena/property`) — stat scaling,
  damage computation (attack vs. defense, armor penetration, backstab detection), skills with
  cooldowns, buffs/curses/effects (poison, stun, shield, heal), and equipment-driven stat
  bonuses.

### Design constraints

Per the engine's contract (`docs/contract_battle_contract.atom.md`):

* **Determinism** — the same seed and input sequence must always produce a bit-identical
  outcome.
* **No shared memory** — two arenas run fully concurrently; all state mutation happens only
  through the actor message queue.
* **Fail fast** — illegal actions (moving onto an occupied cell, acting out of turn, invalid
  skill use) are rejected by validation *before* any state mutation; illegal state transitions
  panic/error rather than silently producing undefined behavior.
* **Performance** — a single combat step resolves in well under 50 ms.

## Repository layout

```
battlearena/
  battlearena.go          entry point: BattleArena / NewBattleArena
  ruler/                  match actor: lifecycle, turn order, GameState, rule enforcement
    gamestate/            grid + entity state owned by the Ruler
    rulermethods/          message types the Ruler accepts / replies with
    rules/                 individual game rule implementations
    turner/                initiative/delay-credit turn scheduling
  controller/             player-facing seam: registration, commands, AI archetypes/behaviors
    controllermethods/     message types a Controller accepts / replies with
  property/               stats, effects (buffs/curses), damage & modifier computation
  battletest/             deterministic scenario sandbox (see below)
docs/                     ATD atoms (contract, vision, mechanics, rules) for this engine
```

Deeper per-package documentation lives alongside the code (e.g.
`battlearena/ruler/README.md`, `battlearena/property/effect/README.md`,
`battlearena/battletest/README.md`).

## Testing

The engine is tested at two altitudes:

* **Integration** — the `ruler` package tests (`NewRuler` / `Start` / `SendActor`) drive the full
  actor/message protocol, turn & shot clock, match lifecycle, win detection, and races.
* **Mechanics** — the [`battletest` scenario sandbox](battlearena/battletest/README.md) is a
  deterministic, service-free harness over `GameState` for testing skills, movement, effects, and
  traps directly against the rules layer (no actor, no clock, no network). Use it for new
  mechanic tests. See its README for the full manual.

Run the full suite with:

```sh
go test ./...
```
