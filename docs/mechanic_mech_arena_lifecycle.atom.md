---
id: mechanic_mech_arena_lifecycle
status: DRAFT
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 4
version: 1.0
human_name: Arena Lifecycle Management
tags: [performance, lifecycle]
---

# New Atom

## INTENT
Manages the safe shutdown and resource reclamation of battle resources once a match is resolved.

## THE RULE / LOGIC
1. Resolution: Match completion triggers a call to `ArenaBridge.DestroyArena`.
2. Map Removal: Match record is removed from `ArenaBridge.arenas`.
3. Actor Shutdown: `Ruler` actor is signaled to stop.
4. Cascading Shutdown: `Ruler` stops its associated `Controllers` before terminating.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_arena_lifecycle]]`
- **Method:** `ArenaBridge.DestroyArena(matchID)`
- **Signals:** `ActorStop` sent to Ruler and Controllers.

## EXPECTATION
- Call to DestroyArena removes match from bridge.
- Ruler actor Stop signal results in all controllers stopping.
- Goroutine count decreases after battle end.
