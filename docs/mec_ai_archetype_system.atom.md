---
id: mec_ai_archetype_system
human_name: AI Archetype System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.1
status: STABLE
priority: 5
tags: [ai,archetypes,behavior,pipeline]
parents:
  - [[mech_behavior_system]]
dependents: []
---

# AI Archetype System Mechanic

## INTENT
To implement four distinct AI archetypes (Fighter, Ranger, Support, Sneak) that follow player progression rules with archetype-specific stat allocation and skill selection.

## THE RULE / LOGIC
**Four archetypes** implemented as `Archetype` interface in `upsilonbattle/battlearena/controller/archetype/`:

| Archetype | Behavior stack (top → bottom) | Stat emphasis |
|-----------|-------------------------------|---------------|
| Fighter   | BattleFocus → ChargeIn → Baseline | Attack, Defense |
| Ranger    | KiteAway → MaintainRange → FocusBackline → Baseline | Attack, Movement |
| Support   | HealAlly → ShieldAlly → StayBehindFront → Baseline | SP, MP, Defense |
| Sneak     | BackstabSeeker → FocusWeakest → Flank → Baseline | Movement, Attack |

**Archetype interface:**
```go
type Archetype interface {
    Slug() string
    Behavior() *behavior.LayeredBehavior   // stack with baseline auto-appended
    StatWeights() StatWeights              // drives CP allocation in entitygenerator
    BuildSkillBundle(grade string, count int) []skill.Skill
}
```

**Registry** (`archetype.Get(slug)`, `archetype.RandomFor(existing []string)`): `RandomFor` respects the team composition constraint by excluding already-capped archetypes from the candidate pool.

**Behavior pipeline:** Each archetype's `Behavior()` returns a `LayeredBehavior` where the baseline `AggressiveBehavior` (BaseActivation=1.0) is always the last layer, guaranteeing a decision every tick. See `[[mechanic_mech_behavior_layered]]`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_ai_archetype_system]]`
- **Key package:** `upsilonbattle/battlearena/controller/archetype/`
- **Files:** `archetype.go`, `fighter.go`, `ranger.go`, `support.go`, `sneak.go`
- **Micro-behaviors:** `upsilonbattle/battlearena/controller/behavior/micro/`
- **Tests:** `archetype_test.go` — registry coverage, stack ordering, stat weights, skill bundle count, RandomFor constraints
