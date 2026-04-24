---
id: mech_version_bit_packing
human_name: "Major.Minor Version Encoding"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 3
tags: [performance, versioning, bitmasking]
parents:
  - [[mech_game_state_versioning]]
dependents: []
---

# Major.Minor Version Encoding

## INTENT
To represent a 2D state progression index (Character Turn . Action Count) within a single 64-bit monotonic integer for compatibility with existing network DTOs and database schemas.

## THE RULE / LOGIC
- **Encoding Strategy:** The 64-bit integer (`Version`) is split into two 32-bit segments:
    - **Major (High 32 bits):** `TurnIndex` — Increments whenever a turn concludes (Pass, Timeout, Forfeit).
    - **Minor (Low 32 bits):** `ActionIndex` — Increments on every state-changing action within a turn (Move, Attack, Skill).
- **Formula:** `Version = (int64(TurnIndex) << 32) | int64(ActionIndex)`
- **Resets:** The `ActionIndex` is reset to `0` every time `TurnIndex` increments.
- **Monotonicity:** Because TurnIndex always increases or stays the same, and ActionIndex resets only when TurnIndex increases, the resulting 64-bit integer remains strictly monotonic over the course of a match.

## TECHNICAL INTERFACE (The Bridge)
- **Go Helper Methods:**
    - `gs.IncTurn()`: `TurnIndex++, ActionIndex = 0`
    - `gs.IncAction()`: `ActionIndex++`
    - `gs.GetTurn()`: `(gs.Version >> 32)`
    - `gs.GetAction()`: `(gs.Version & 0xFFFFFFFF)`
- **Code Tag:** `@spec-link [[mech_version_bit_packing]]`
- **Location:** `battlearena/ruler/rules/gamestate.go`

## EXPECTATION (For Testing)
- Version `1.0` (Turn 1, Action 0) -> `4294967296`
- Version `1.1` (Turn 1, Action 1) -> `4294967297`
- Version `2.0` (Turn 2, Action 0) -> `8589934592`
- Comparison `2.0 > 1.1` must always be true.
- Bitwise extraction of major/minor from a stored `int64` yields the correct integers.
