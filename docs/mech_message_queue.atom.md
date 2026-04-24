---
id: mech_message_queue
human_name: "Message Queue Mechanic"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: [messaging, async, queue, concurrency]
parents:
  - [[module_actor_concurrency]]
dependents: []
---

# Message Queue Mechanic

## INTENT
Provide a single-threaded asynchronous execution queue for messages to ensure strictly sequential processing of logic.

## THE RULE / LOGIC
- **Unified Stimuli Intake**: The queue serves as the single intake point for all actor stimulae, including new messages and replies (callbacks).
- **FIFO Ordering**: Messages must be processed in the exact order they are received via the `inputChan`.
- **Sequential Execution**: Only one stimulus (message or reply) can be "in-flight" (executing) at any given time.
- **Acknowledge Required**: The queue remains blocked for the next stimulus until an acknowledgment (ACK) is received via the `executorReplyChan`.
- **Internal Non-Blocking Buffer**: The queue uses an internal slice to store incoming stimulae immediately, ensuring the sender (caller or replier) is never blocked.
- **Termination Lifecycle**: When `Stop` or `PrepareStop` is called, the queue must complete the current stimulus (if any) and then terminate processing.
- **Resilience**: The queue must not crash if an acknowledgment is received while the message list is empty (phantom ACK).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag**: `@spec-link [[mech_message_queue]]`
- **Primary Methods**:
  - `Send(msg)`: Enqueue a message.
  - `Start()`: Begin the internal processing loop.
  - `GetExecutorChan()`: Channel for the executor to receive messages.
  - `GetExecutorReplyChan()`: Channel for the executor to signal completion.

## EXPECTATION (For Testing)
- **Sequentiality**: Multiple enqueued messages are dispatched one by one.
- **Blocking**: No new message is dispatched until the previous one is ACKed.
- **Robustness**: Unexpected ACKs must be ignored or logged, but NEVER cause a panic.
- **Concurrency**: Access to the internal message list must be thread-safe.
