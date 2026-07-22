---
id: mechanic_channeling_mechanic
human_name: Channeling Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 3.0
status: SPECIFIED
priority: 5
tags: [time-based, skills, casting, initiative]
parents:
  - [[mech_initiative]]
  - [[mechanic_pre_post_execution_costs]]
dependents: []
---

# Channeling Mechanic

## INTENT

To implement channeled skills: a skill with a Channeling delay > 0 commits its
caster to a high-impact action that resolves *later*. The caster pays all costs
upfront and is locked out and vulnerable while the channel runs; the effect only
fires when the channel completes. Taking enough damage interrupts the channel.

## THE RULE / LOGIC

**Core Principle — Caster Reschedule (NOT a separate entity):**
The caster is already a cursor in the initiative queue (Turner). Channeling does
NOT spawn a hidden tracking entity; it simply pushes the caster's next tick out by
the channel delay. While dormant the caster stays in the arena (in `Entities`, on
the grid) and is therefore targetable. When its delayed tick arrives, the stored
skill resolves.

> Supersedes v2.0, which modelled channeling as a spawned `TimeBased` temporary
> entity. That dependency on the temporary-entity system is removed.

**Channeling Cost (Pre-Execution):**
- **Property:** `Channeling` is measured in delay units (absence = 0 = not a
  channel). e.g. Fireball with Channeling 400 reschedules the caster +400.
- **Risk Premium:** Channeling is balanced at **-15 SW per 10 delay** units (vs
  -10 SW for normal delay) — more power per cost in exchange for the vulnerability.
- **Sunk Costs:** All pre-execution resources (SP, MP, HP, Mvt) are paid upfront at
  cast initiation and are NOT refunded if the channel is interrupted or fizzles.

**State — `IsCasting`:**
The caster holds an `IsCasting` record `{SkillID, TargetEntity | TargetPos,
Interruption}`. The target is captured per the skill's targeting mode: entity-target
skills store `TargetEntity` (the channel FOLLOWS the entity, which may move before
resolution); tile/self skills store a fixed `TargetPos`.

**Lifecycle:**
1. **Cast init** (`rules/skill.go UseSkill` → `beginChannel`): detect
   `Channeling.MaxValue > 0`; validate + pay ALL costs upfront; set `IsCasting`;
   mark `HasActed`/`HasMoved`; do NOT apply the effect.
2. **Auto-pass** (`ruler_actions.go controllerUseSkill` → `rules/endofturn.go`):
   the caster's turn ends immediately — it does NOT wait for the controller to send
   an EndOfTurn. The channel branch in EndOfTurn reschedules the caster by
   `+channelDelay` (NOT the flat 300 Pass base), keeps it locked (no HasActed/
   HasMoved reset, no movement restore), and re-queues it in the Turner.
3. **Dormant:** the caster sits in the Turner at `+channelDelay`, never picked until
   completion, still targetable on the grid.
4. **Resolution** (`ruler_turn.go handTurn` → `rules/channeling.go ResolveChannel`):
   when the dormant caster is re-picked, it does NOT start a shot clock or receive a
   ControllerNextTurn. The target is re-derived (entity targets follow to their
   current position) and **re-validated** (range/grid/type). If the target is gone
   or out of range the channel **FIZZLES** (no effect; costs stay sunk). Otherwise
   the stored skill's effect resolves. `IsCasting` is cleared, the caster runs a
   normal end-of-turn (recovery delay, flag reset, re-queue), and the turn is handed
   to the appropriate next entity in initiative (which may be the caster again).
5. **Chained / ticking channel:** a channel whose resolved skill is itself a channel
   simply starts a new cast at step 1 — no special case, no separate entity.

**Interruption Mechanics:**
- **Property:** `Interruption` 0-100, accumulates when the caster takes damage.
- **Formula:** **1 damage = 10 interruption points**.
- **Failure Threshold:** at `Interruption >= 100` the channel fails immediately —
  `IsCasting` is cleared, all pre-execution resources are wasted (not refunded), and
  the caster is pulled out of its distant dormant slot and re-queued into a near
  recovery equal to the skill's `Delay` cost (absence = 500), explicitly NOT the
  channel delay (auto-pass to recovery).
- **Damage sites hooked:** skill effects, basic attacks, and positional effects all
  feed the gauge via a single `ApplyInterruption` helper. (Poison cannot interrupt:
  a dormant caster never takes its own end-of-turn while channeling.)

**Death during channel:** normal `RemoveEntity` cleanup; the channel just never
resolves. No special handling.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_channeling_mechanic]]`
- **State:** `entity.CastingState` + `Entity.IsCasting` / `Entity.IsChanneling()`
  (`upsilontypes/entity/entity.go`).
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/channeling.go`
  (`beginChannel`, `ResolveChannel`, `ApplyInterruption`),
  `rules/skill.go` (`UseSkill`, `applyDirectSkillEffect`),
  `rules/endofturn.go` (channel reschedule), `rules/attack.go` &
  `rules/positionaleffect.go` (interruption hooks),
  `ruler/ruler_turn.go` (`handTurn`/`resolveChannel`/`advanceTurn`),
  `ruler/ruler_actions.go` (auto-pass).
- **API Projection:** `upsilonapi/api` `Casting` DTO + `convertCastingState`
  (serializes the "Channeling: X" indicator + interruption gauge).
- **Integration:** Works with `[[mechanic_pre_post_execution_costs]]` and
  `[[mech_initiative]]`.

## EXPECTATION
- A channeled skill deducts all pre-execution costs at initiation; these are not
  refunded on interruption, fizzle, or death.
- Cast initiation ends the caster's turn immediately (no controller EndOfTurn) and
  reschedules it by the skill's Channeling value; the caster enters `IsCasting`.
- The effect resolves only when the dormant caster is re-picked, after the target is
  re-derived and re-validated; an out-of-range or gone target fizzles the channel.
- Entity-target channels resolve against the target's CURRENT position (it may have
  moved); tile-target channels resolve at the fixed cast-time tile.
- Taking damage while channeling adds 10 interruption per 1 damage; reaching >= 100
  fails the channel, wastes the sunk costs, and recovers the caster at the skill's
  Delay cost (not the channel delay).
- Channeled skills are balanced at -15 SW per 10 delay units (vs -10 for normal).
