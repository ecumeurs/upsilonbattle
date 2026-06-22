---
id: mechanic_ai_termination
status: DRAFT
human_name: AI Termination Pattern
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
dependents: []
priority: 4
tags: [concurrency, ai]
parents:
  - [[module_actor_concurrency]]
---

# AI Termination Pattern

## INTENT
Ensures AI controllers terminate without blocking their main actor loop during match resolution.

## THE RULE / LOGIC
**AI Termination Pattern:**

**Core Principle:**
To prevent deadlocks and ensure responsive system shutdown during match resolution, AI controllers must employ non-blocking communication patterns for lifecycle termination signals.

**Termination Logic and Concurrency Rules:**
1. **Buffered Communication:** Channels used to signal "Battle Finished" or similar end-of-lifecycle events must be initialized with a buffer (minimum size 1). This ensures that the sending routine can proceed without waiting for a consumer to be ready.
2. **Non-Blocking Signal Dispatch:** All attempts to send a termination signal into a lifecycle channel must use a non-blocking selection pattern. If the channel is full or no consumer is available, the send operation should be skipped rather than allowing the controller to hang.
3. **Lifecycle Priority:** The processing of termination signals (e.g., "Battle Over") must never impede the controller's ability to receive and process a higher-priority "Actor Stop" message.
4. **Independent Resolution:** Once a termination signal is dispatched, the controller enters a "Graceful Shutdown" state where it ceases tactical calculations but remains receptive to final system cleanup commands.

**System Stability Requirements:**
- **Hang Prevention:** The controller must be guaranteed to exit its primary execution loop within a deterministic timeframe after a match concludes.
- **Async Safety:** The termination pattern must account for scenarios where multiple systems might attempt to signal the end of a battle simultaneously.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_ai_termination]]`
- **Component:** `AggressiveController`
- **Channel:** `BattleFinished (buffered)`

## EXPECTATION
- AI Controller does not hang during BattleEnd.
- Message processing continues until ActorStop is received.
- BattleFinished channel send returns immediately if unconsumed.
