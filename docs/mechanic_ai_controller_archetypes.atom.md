---
id: mechanic_ai_controller_archetypes
human_name: "Ai Controller Archetypes"
status: STABLE
version: 2.1
parents:
  - [[mech_behavior_system]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# Ai Controller Archetypes

## INTENT
To implement four distinct AI controller archetypes (Fighter, Ranger, Support, Sneak), each a code-registered `Archetype` with its own layered behavior stack, stat weights, and skill bundle, defining the decision pattern for AI-controlled characters.

## THE RULE / LOGIC
**AI Controller Archetypes:**

**Core Principle:**
Each AI archetype is implemented as an `Archetype` value in `upsilonbattle/battlearena/controller/archetype/`, extending baseline behavior with a specialized layered behavior stack, stat weighting, and skill selection.

**Archetype interface:**
```go
type Archetype interface {
    Slug() string
    Behavior() *behavior.LayeredBehavior   // stack with baseline auto-appended
    StatWeights() StatWeights              // drives CP allocation in entitygenerator
    BuildSkillBundle(grade string, count int) []skill.Skill
}
```

**Behavior stack per archetype (top → bottom) and stat emphasis:**

| Archetype | Behavior stack (top → bottom) | Stat emphasis |
|-----------|-------------------------------|---------------|
| Fighter   | BattleFocus → ChargeIn → Baseline | Attack, Defense |
| Ranger    | KiteAway → MaintainRange → FocusBackline → Baseline | Attack, Movement |
| Support   | HealAlly → ShieldAlly → StayBehindFront → Baseline | SP, MP, Defense |
| Sneak     | BackstabSeeker → FocusWeakest → Flank → Baseline | Movement, Attack |

**Per-archetype tactical intent and decision loop:**
- **Fighter** — Direct aggression / frontline engagement. Positions adjacent to enemies (melee range). Decision loop: (1) check for killable enemies within reach and execute a killing blow; (2) otherwise use a charge skill to close distance; (3) default to attacking the nearest enemy with the most efficient melee skill. Stat distribution emphasizes Attack and Defense.
- **Ranger** — Ranged harassment / kiting. Maintains optimal distance (~4-7 cells), favors high ground or cover. Decision loop: (1) if outside optimal range, move to a safe distance; (2) place traps at chokepoints if available; (3) attack the most exposed enemy from range. Stat distribution emphasizes Attack and Movement/Accuracy.
- **Support** — Ally preservation / combat enhancement. Stays near teammates, blocking paths to vulnerable allies; zero Attack. Decision loop: (1) heal the most damaged ally below a critical threshold; (2) apply shields to allies threatened by high-damage enemies; (3) buff teammates lacking active enhancements; (4) move to maximize support coverage. Stat distribution emphasizes SP/MP pools and Defense.
- **Sneak** — Flanking / precision status application. Aggressively seeks cells behind enemies to trigger backstab bonuses. Decision loop: (1) move behind a vulnerable target for a backstab attack; (2) apply status effects (poison/stun) to high-value targets; (3) use evasion or stealth if under direct threat. Stat distribution emphasizes Movement and Attack (high Dodge/Crit).

**Registry:**
- `archetype.Get(slug)` returns a named archetype.
- `archetype.RandomFor(existing []string)` returns a random archetype, respecting team-composition constraints by excluding already-capped archetypes from the candidate pool (see `[[rule_ai_team_composition_rules]]`).

**Behavior pipeline:** Each archetype's `Behavior()` returns a `LayeredBehavior` where the baseline `AggressiveBehavior` (BaseActivation=1.0) is always the last layer, guaranteeing a decision every tick. See `[[mechanic_behavior_layered]]`.

**Performance:** Decision trees and tactical evaluations are cached per archetype for the duration of a turn to reduce redundant calculation.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_ai_controller_archetypes]]`
- **Key package:** `upsilonbattle/battlearena/controller/archetype/`
- **Files:** `upsilonbattle/battlearena/controller/archetype/archetype.go`, `fighter.go`, `ranger.go`, `support.go`, `sneak.go`
- **Micro-behaviors:** `upsilonbattle/battlearena/controller/behavior/micro/`
- **Tests:** `upsilonbattle/battlearena/controller/archetype/archetype_test.go`
- **Integration:** Works with `[[mechanic_ai_progression_matching]]`, `[[rule_ai_team_composition_rules]]`, `[[mechanic_behavior_layered]]`

## EXPECTATION
- The registry exposes exactly four archetypes (Fighter, Ranger, Support, Sneak); `archetype.Get` returns each (`TestRegistryCoversAllFour`).
- Every archetype's `Behavior()` stack has the baseline `AggressiveBehavior` as its last layer (`TestBehaviorStackHasBaselineAsLastLayer`).
- Every archetype's behavior stack has at least two layers (`TestBehaviorStackHasAtLeastTwoLayers`).
- Every archetype's `StatWeights()` returns strictly positive weights (`TestStatWeightsArePositive`).
- `archetype.RandomFor(existing)` never returns an archetype already at its composition cap and returns a valid archetype for an empty team (`TestRandomForRespectsCompositionConstraints`, `TestRandomForEmptyTeam`).
- `BuildSkillBundle(grade, count)` returns exactly `count` skills (`TestBuildSkillBundleReturnsRequestedCount`).
