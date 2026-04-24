---
id: mech_controller_handshake
human_name: Controller Handshake Protocol
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: [handshake, initialization, controller]
parents:
  - [[api_controller_methods]]
dependents:
  - [[mech_controller_communication_sequence]]
---
# Controller Handshake Protocol

## INTENT
To ensure bidirectional communication between the Ruler and a Controller is established during the battle registration phase.

## THE RULE / LOGIC
1.  **Registration**: A Controller is added to the Ruler via `AddController`.
2.  **Handshake Notification**: The Ruler immediately sends a `SetQueue` message back to the Controller.
3.  **Reference Storage**: The Controller receives `SetQueue` and saves the `Ruler`'s actor reference.
4.  **Verification**: The Controller typically responds by requesting the initial game state (`GetGridState`, `GetEntitiesState`) using the saved Ruler reference.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_controller_handshake]]`
- **Related Issue:** `#None`
- **Test Names:** `TestRulerBattleBegin`

## EXPECTATION (For Testing)
- After `AddController`, the controller MUST have a non-nil `ruler` reference.
