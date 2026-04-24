---
id: mechanic_script_lifecycle
human_name: "Script Lifecycle and Teardown"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 2
tags: [scripting, farm, lifecycle]
parents:
  - [[script_farm]]
dependents: []
---

# Script Lifecycle and Teardown

## INTENT
Ensure that every agent script executes a guaranteed cleanup phase regardless of success or failure.

## THE RULE / LOGIC
1.  **Setup Phase**: The agent is initialized, WebSocket connection established.
2.  **Execution Phase**: The main JS script runs.
3.  **Teardown Phase**: 
    -   Must run even if `Execution Phase` throws an exception, assertion fails, or the process receives an INTERRUPT signal (SIGINT/SIGTERM).
    -   Triggered by Go's `defer` mechanism and context cancellation.
    -   Executes the function registered via `upsilon.onTeardown(callback)`.
    -   Gracefully stops the WebSocket listener.

## TECHNICAL INTERFACE (The Bridge)
-   **Method:** `upsilon.onTeardown(callback: function)`
-   **Code Tag:** `@spec-link [[mechanic_script_lifecycle]]`
-   **Implementation:** `internal/script/coordinator.go` (defer block), `internal/script/bridge.go` (`jsOnTeardown`)

## EXPECTATION (For Testing)
1.  Register a teardown hook that logs a message.
2.  Force a script failure (e.g., `upsilon.assert(false)`).
3.  Verify that the special teardown log message appears in the agent's output.
