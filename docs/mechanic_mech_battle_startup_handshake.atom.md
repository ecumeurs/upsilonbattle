---
id: mechanic_mech_battle_startup_handshake
status: STABLE
type: MECHANIC
parents:
  - [[uc_combat_turn]]
version: 1.0
dependents: []
human_name: Battle Startup Handshake
layer: ARCHITECTURE
priority: 5
---

# New Atom

## INTENT
To ensure the transition from Turn 0 (Match Initialization) to Turn 1 (Tactical Combat) is handled asynchronously, preventing delivery races and state corruption in external controllers (bots).

## THE RULE / LOGIC
- When the Ruler receives a `BattleStart` notification or signals combined readiness (`isBattleReadyToExecute`), it broadcasts the `BattleStart` meta-event (Version 0).
- Instead of triggering the first turn synchronously, it schedules a `SelfNotifyDelayed` message for `InternalTriggerFirstTurn`.
- This ensures that the initialization state is fully dispatched and processed by the Bridge/Webhook layer before Turn 1 (Active combat) begins.

## TECHNICAL INTERFACE
- **Internal Message:** `rulermethods.InternalTriggerFirstTurn{}`
- **Code Tag:** `@spec-link [[mech_battle_startup_handshake]]`
- **Delay Constraint:** 100ms minimum handover delay.

## EXPECTATION
- Request for Battle Start arrives -> BattleStart broadcast with Version 0 sent immediately.
- 100ms pause -> First tactical turn (Version 4294967296) triggered.
- Result: Webhook consumers receive the game initialization state before receiving the first combat turn.
