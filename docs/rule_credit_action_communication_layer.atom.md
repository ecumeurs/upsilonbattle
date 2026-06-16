---
id: rule_credit_action_communication_layer
status: STABLE
priority: 5
layer: ARCHITECTURE
version: 2.0
tags: ["credits", "communication", "api"]
parents:
  - [[req_tech_debt_backlog]]
dependents: []
type: RULE
---

# New Atom

## INTENT
To establish the communication layer protocol for all combat-related actions, ensuring traceability through request IDs, version control for compatibility, effect feedback for result confirmation, and credit tracking association with player IDs.

## THE RULE / LOGIC
**Credit Action Communication Layer Protocol:**

**Core Principle:**
This protocol defines the standardized bidirectional communication between the game engine and its clients (API, UI, CLI), ensuring that every combat action is traceable, versioned, and correctly attributed for economic rewards.

**Message and Response Envelopes:**
- **Action Request Envelope:**
    - **Traceability:** Every request must carry a unique `RequestID` (UUIDv7) for end-to-end tracking.
    - **Versioning:** A strict version string (e.g., "v2.0.0") is required to ensure compatibility between engine and client.
    - **Identity:** Clearly identifies the source `EntityID` and, where applicable, the `TargetID`.
- **Resolution Response Envelope:**
    - **Correlation:** The response must echo the original `RequestID` to allow client-side synchronization.
    - **Outcome:** A boolean `Success` flag paired with a descriptive `Error` string in case of failure.
    - **State Delta:** Includes a `Modified` block containing the specific game state changes (HP updates, position shifts) resulting from the action.
    - **Economic Attribution:** Explicitly returns the `Credits` earned by the action and the `PlayerID` of the recipient.

**Credit Attribution Standards:**
- **Mandatory Attribution:** Every action capable of generating value (damage, healing, utility) must resolve with a credit award, even if the value is zero.
- **Persistence:** Credits are assigned to the `PlayerID` linked to the original action creator. The system ensures that support players (e.g., those providing shields) receive their earned credits even if their character is removed from the grid before the credit-triggering event resolves.
- **Immediate Settlement:** Credits returned in the `ActionResponse` are considered finalized and are immediately reflected in the player's account balance.

**Traceability and Logging Protocol:**
- **Request Identification:** The use of UUIDv7 ensures global uniqueness and temporal sorting for log analysis.
- **Ref-ID Convention:** System logs must prefix all action-related entries with an 8-character `ref_id` derived from the Request ID for easy cross-referencing.
- **Version Guarding:** The engine performs a handshake validation on the `Version` field; mismatching versions trigger an automatic rejection to prevent state corruption.

**Transactional Lifecycle:**
1. **Initiation:** Client dispatches an `ActionMessage`.
2. **Validation:** The engine verifies state integrity, resource availability, and protocol version.
3. **Execution:** The engine resolves the mechanic and calculates game state deltas and credit awards.
4. **Synchronization:** The engine dispatches an `ActionResponse` to the requester and broadcasts state updates to all other observers.
5. **Persistence:** The calculated credits are committed to the player's permanent database record.

## TECHNICAL INTERFACE

## EXPECTATION
