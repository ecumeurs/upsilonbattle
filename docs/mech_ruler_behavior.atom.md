---
id: mech_ruler_behavior
human_name: "Ruler Behavior System"
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: OBSOLETE
priority: 3
tags: [engine, automation, ruler]
parents:
  - [[mech_behavior_system]]
dependents: []
---

# Ruler Behavior System

## INTENT
DEPRECATED: this records the now-excised in-Ruler synchronous AIBehavior path; it survives only as the anchor for the ISS-101 controller-ownership invariant.

## THE RULE / LOGIC
**OBSOLETE — the in-Ruler `AIBehavior` automation path has been removed.**

Previously, `handTurn`/`triggerFirstTurn` read the `property.AIBehavior` slug and, for a non-`"none"` value, diverted the entity's turn into `upsilonbattle/battlearena/ruler/behavior/` (`AggressiveBehavior`/`ExpirationBehavior`), self-dispatching the resulting action and bypassing the controller stack. That mechanism was vestigial dead weight in production: no API DTO carried `AIBehavior` and no service ever set it. It has been excised:
- The `ruler/behavior/` package is deleted.
- The `property.AIBehavior` definition is removed from `upsilontypes`.
- The `behaviorSlug != "none"` branch in `handTurn` is gone.

**What remains under this atom (ISS-101 controller-ownership invariant):** every entity must have a non-nil `ControllerID`. `AddEntity` rejects `uuid.Nil` at registration and `handTurn` panics defensively if it ever sees one. All automation now flows exclusively through the asynchronous controller path (`[[mech_controller_behavior]]` / `controllers.AIController`).

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_ruler_behavior]]`
- **Surviving anchor:** `upsilonbattle/battlearena/ruler/ruler.go` (`AddEntity` Nil-ControllerID guard, ISS-101)
- **Tests:** `TestRulerAddEntityRejectsNilControllerID`, `TestRulerHandTurnRejectsNilControllerID` (`ruler_iss101_test.go`)
- **Removed:** `upsilonbattle/battlearena/ruler/behavior/` (deleted), `property.AIBehavior` (removed from upsilontypes)

## EXPECTATION
No entity is ever automated via an in-Ruler behavior slug. Every entity has a non-nil ControllerID (enforced at AddEntity and defended in handTurn); reaching either with uuid.Nil panics. All AI automation is driven by the asynchronous controller stack only.
