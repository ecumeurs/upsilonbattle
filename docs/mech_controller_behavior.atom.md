---
id: mech_controller_behavior
human_name: "Controller Behavior System"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 3
tags: [ai, controller, actor]
parents:
  - [[mech_behavior_system]]
dependents: []
---

# Controller Behavior System

# INTENT
To enable complex, pluggable AI decision-making for actor-based controllers decoupled from engine plumbing.

## THE RULE / LOGIC
- **Asynchronous Execution:** Behaviors run within a `Controller` actor's goroutine, safe from engine loop blocking.
- **Semantic Decisions:** Returns a `Decision` struct (Move, Attack, Skill, Pass) rather than raw messages.
- **Pluggable Architecture:** Controllers hold a `Behavior` interface, allowing hot-swapping or composite logic.
- **Composite Support:** Supports `CompositeBehavior` (AND/OR modes) to combine multiple atomic tactical rules.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_controller_behavior]]`
- **Location:** `upsilonbattle/battlearena/controller/behavior/`

## EXPECTATION
AI characters can perform complex tactical analysis and execute multi-step plans without direct engine-level implementation of the AI logic.
