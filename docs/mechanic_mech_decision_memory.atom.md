---
id: mechanic_mech_decision_memory
status: STABLE
parents:
  - [[mech_behavior_system]]
dependents: []
human_name: AI Decision Memory
type: MECHANIC
layer: IMPLEMENTATION
version: 1.1
priority: 3
tags: [ai, memory, behavior, state]
---

# New Atom

## INTENT
To maintain a fixed-size ring buffer of per-layer decision records across ticks and turns, enabling behavior layers to self-throttle and maintain sticky targeting.

## THE RULE / LOGIC
**Ring buffer:** `DecisionMemory` holds the last 20 `DecisionRecord` entries (constant `memorySize = 20`). Oldest entries are overwritten when the buffer is full.

**Record contents:** `LayerName`, `Skipped` (bool), `Turn`, `Tick`, `TargetID` (uuid.UUID).

**Queries:**
- `TurnsSince(layer string) int` — turns elapsed since that layer last wrote a non-skipped record. Returns a large sentinel if never recorded.
- `CountInLastN(layer string, n int) int` — number of non-skipped activations for a layer in the last N turns.
- `CurrentTarget() uuid.UUID` — the TargetID from the most recent non-skipped record with a non-nil target; sticky across ticks within a turn.
- `ClearTarget()` — clears the sticky target (called on turn advance).
- `AdvanceTick()` / `AdvanceTurn()` — increment internal turn/tick counters.

**Cooldown pattern:** Layers call `memory.TurnsSince("layer_name") < threshold` to self-throttle.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_decision_memory]]`
- **Key type:** `DecisionMemory`
- **Package:** `upsilonbattle/battlearena/controller/behavior/`
- **Parents:** `[[mech_behavior_layered]]`

## EXPECTATION
- `TurnsSince` returns a large sentinel for a layer that has never fired.
- `CurrentTarget` is nil-UUID after `ClearTarget` or when no target has been recorded.
- Ring buffer does not grow beyond 20 entries.
