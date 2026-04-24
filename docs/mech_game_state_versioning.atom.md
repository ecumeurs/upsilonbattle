---
id: mech_game_state_versioning
human_name: "Game State Versioning & Deduplication"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 1
tags: [performance, networking, deduplication]
parents:
  - [[api_go_webhook_callback]]
dependents:
  - [[mech_version_bit_packing]]
---

# Game State Versioning & Deduplication

## INTENT
Ensure only the latest state is broadcasted and processed by deduplicating monotonic state updates to optimize network and CPU usage.

## THE RULE / LOGIC
- **Monotonic Tick:** Every state-changing action in the Go Engine must increment the `GameState.Version` field (int64).
- **Major.Minor Encoding:** The version is split into `TurnIndex` (major) and `ActionIndex` (minor) using bit-packing to ensure monotonicity while tracking tactical progression. See [[mech_version_bit_packing]] for bitwise details.
- **Single Source of Truth:** The `version` field is the definitive marker of match progression. Legacy `turn` counters are strictly mapped to this version to ensure consistency across the stack.
- **Engine-Side Deduplication:**
    - The engine tracks the `lastSentWebhookVersion` for each active match ID.
    - Before forwarding any event (Move, Attack, TurnStart, etc.) to the Laravel webhook, the engine checks if `currentVersion > lastSentVersion`.
    - If it is not greater, the webhook is dropped to prevent redundant updates.
- **Gateway-Side Validation:**
    - Laravel strictly validates incoming webhooks via `WebhookRequest`, requiring `data.match_id` and `data.version`.
    - Incoming webhooks with `version <= current_version` are ignored, preventing race conditions or out-of-order processing.
    - Version `0` is treated as a special initialization case and is always accepted if the match is new.
- **Broadcast Efficiency:**
    - Laravel only broadcasts `BoardUpdated` events to WebSockets for successfully validated new versions.

## TECHNICAL INTERFACE (The Bridge)
- **Field Name:** `sequence` (in BoardState DTO), `version` (in ArenaEvent DTO and Database).
- **Code Tag:** `@spec-link [[mech_game_state_versioning]]`
- **Location (Go):** `battlearena/ruler/rules/gamestate.go`, `upsilonapi/bridge/bridge.go`
- **Location (Laravel):** `app/Http/Controllers/API/WebhookController.php`

## EXPECTATION (For Testing)
- Multiple identical state updates from the engine result in only ONE database write and ONE broadcast.
- Out-of-order webhooks (lower version than stored) are silently dropped.
- The state version is visible in the frontend/CLI to ensure clients can detect missing or stale updates.
