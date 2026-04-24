---
id: mech_controller_communication_sequence
human_name: Controller-Ruler Communication Sequence
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: [sequence, protocol, controller, ruler]
parents:
  - [[api_controller_methods]]
  - [[api_ruler_methods]]
  - [[mech_controller_handshake]]
dependents: []
---
# Controller-Ruler Communication Sequence

## INTENT
To define the chronological order of messages exchanged between the Ruler and Controllers during a battle session.

## THE RULE / LOGIC
The exchange follows this sequence:

### 1. Initialization Phase (Handshake)
- **Ruler** receives `AddController`.
- **Ruler** notifies **Controller** with `SetQueue`.
- **Controller** requests `GetGridState` and `GetEntitiesState`.
- **Controller** notifies **Ruler** with `ControllerBattleReady`.

### 2. Battle Start
- **Ruler** broadcasts `BattleStart` to all Controllers once all are `BattleReady`.

### 3. Combat Loop (Per Turn)
- **Ruler** determines next turn and notifies active **Controller** with `ControllerNextTurn`.
- **Controller** issues actions (e.g., `ControllerMove`, `ControllerAttack`).
- **Ruler** replies to actions with Result/State updates.
- **Ruler** broadcasts `EntitiesStateChanged` to all Controllers.
- **Controller** signifies end of turn with `EndOfTurn`.

### 4. Conclusion
- **Ruler** broadcasts `BattleEnd` with winner details.
- **Controller** may send `ControllerQuit`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_controller_communication_sequence]]`
- **Related Issue:** `#None`

## EXPECTATION (For Testing)
- Messages must follow this order. Out-of-order actions (e.g., `ControllerMove` before `BattleStart`) should be rejected by the Ruler.
