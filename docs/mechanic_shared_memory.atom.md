---
id: mechanic_shared_memory
human_name: "Agent Shared Memory"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 2
tags: [scripting, farm, coordination]
parents:
  - [[script_farm]]
dependents: []
---

# Agent Shared Memory

## INTENT
Enable coordination between isolated agents in a farm (e.g., sharing a Match ID during PVP).

## THE RULE / LOGIC
1.  Provide a thread-safe `SharedStore` accessible by all agents in a single farm run.
2.  Methods must use read/write mutexes to prevent race conditions during concurrent execution.
3.  Store is volatile and cleared between farm invocations.

## TECHNICAL INTERFACE (The Bridge)
-   **Method:** `upsilon.setShared(key: string, value: any)`
-   **Method:** `upsilon.getShared(key: string): any`
-   **Code Tag:** `@spec-link [[mechanic_shared_memory]]`
-   **Implementation:** `internal/script/store.go`, `internal/script/bridge.go` (`jsSetShared`, `jsGetShared`)

## EXPECTATION (For Testing)
1.  Agent A sets a value: `upsilon.setShared("test_key", "hello")`.
2.  Agent B waits/sleeps and then reads: `val = upsilon.getShared("test_key")`.
3.  Verify `val === "hello"` in Agent B.
