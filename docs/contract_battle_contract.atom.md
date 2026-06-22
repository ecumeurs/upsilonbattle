---
id: contract_battle_contract
status: STABLE
type: CONTRACT
version: 1.0
human_name: Upsilon Battle Engine Contract
layer: BUSINESS
priority: 1
tags: [governance, contract, engine]
parents:
  - [[shared:contract_upsilon_contract]]
dependents: []
---

# Upsilon Battle Engine Contract

## INTENT
Establish the implementation standards and concurrency rules for the Upsilon Battle Engine.

## THE RULE / LOGIC
- **Concurrency:** Must utilize the `[[upsilontools:mech_actor_pattern]]` for thread-safe state management. No shared memory between Arenas.
- **Determinism:** All logic must be deterministic given the same seed and input sequence.
- **Validation:** Every action (Move, Attack, Skill) must be validated against the current state before execution.
- **Performance:** Combat step resolution must complete in under **50 ms** (no nanosecond-latency requirement).
- **Error Safety:** Fail fast on illegal state transitions rather than allowing undefined behavior.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[contract_battle_contract]]`
- **Related Atoms:** `[[upsilontools:mech_actor_pattern]]`, `[[shared:contract_upsilon_contract]]`

## EXPECTATION
- Two Arenas run concurrently with no shared mutable state; state is mutated only through the actor message queue.
- Replaying the same seed and identical input sequence yields a bit-identical match outcome.
- An illegal action (e.g. moving onto an occupied cell, acting out of turn) is rejected by validation before any state mutation.
- A single combat step resolves in under 50 ms.
- An illegal state transition fails fast (returns an error / panics in dev) rather than producing undefined behavior.
