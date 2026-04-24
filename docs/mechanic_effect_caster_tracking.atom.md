---
id: mechanic_effect_caster_tracking
status: DRAFT
priority: 5
version: 2.0
parent: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
tags: [effects, caster, credits]
parents: []
---

# Effect Caster Tracking

## INTENT

To implement effect caster tracking system where all combat effects remember their original caster until effect ends, enabling proper credit assignment, interruption mechanics, and support play credit earning.

## THE RULE / LOGIC
**Effect Caster Tracking Mechanic:**

**Core Principle:**
To ensure accurate credit attribution and mechanical integrity, the system maintains a persistent link between every active effect and its original creator. This "CasterID" tracking is fundamental to the game's economic and tactical resolution layers.

**Tracking Architecture:**
- **The CasterID Anchor:** Every effect object (Damage, Healing, Shields, Status Effects) contains a mandatory `CasterID` field that stores the unique identifier of the entity that initiated the action.
- **Relational Integrity:** This link is preserved even if the effect is transferred (e.g., a "Poison Cloud" creates "Poisoned" status effects on multiple targets; each status effect inherits the cloud's original `CasterID`).

**Lifecycle and Expiration Rules:**
- **The "Expires with Caster" Property:** A boolean flag that determines the effect's persistence upon the death or removal of the caster.
    - **True:** The effect (e.g., a summoned barrier or a channeled beam) is immediately purged when the caster leaves the game state.
    - **False:** The effect (e.g., a fired projectile or a lingering debuff) continues to resolve until its own duration or condition is met, regardless of the caster's status.
- **Cleanup Resolution:** Upon entity removal, the system audits all grid cells and active status effect stacks, surgically removing only those orphaned effects tagged with `ExpiresWithCaster = True`.

**Credit and Economic Attribution:**
- **Direct Rewards:** Credits for damage dealt or healing performed are routed directly to the account of the `CasterID`.
- **Post-Mortem Attribution:** If a caster dies but their lingering effects (e.g., a poison DOT) continue to deal damage, the resulting credits are still credited to the original caster's total, acknowledging their contribution to the team's success.
- **Support Play Validation:** Credits for damage mitigation (Shields) or tactical control (Stuns) are awarded to the original caster, providing an economic incentive for non-offensive playstyles.

**System Integration:**
- **Interruption Mechanics:** Tracking allows the system to identify which active effects should be cancelled when a character is stunned or silenced.
- **Auditability:** The combat log and replay system use `CasterID` to clearly communicate to players exactly who was responsible for any given tactical event.

## TECHNICAL INTERFACE

- **Code Tag:** `@spec-link [[effect_caster_tracking]]`
- **Related Files:** `upsilonbattle/battlearena/property/effect/effect.go`, `upsilonbattle/battlearena/property/effect/effectapplicator/effectapplicator.go`

## EXPECTATION
