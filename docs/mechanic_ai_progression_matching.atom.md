---
id: mechanic_ai_progression_matching
status: STABLE
version: 1.0
parents:
  - [[mech_behavior_system]]
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
human_name: AI Grade Progression
tags: [ai,grade,progression,cp]
---

# AI Grade Progression

## INTENT
To implement AI progression matching system where AI characters follow the same point-buy progression rules as players, scaling stats, skills, and level according to player team averages.

## THE RULE / LOGIC
**Grade table** (`upsilontypes/entity/grade/`):

Grades: I, I+, II, II+, III, III+, IV, IV+, V (9 tiers, indices 0–8).

```
GradeFromWins(wins):
  wins < 0 → clamp to 0
  wins ≥ 40 → "V" (hard cap, no V+)
  main = wins / 10           // 0..3
  if wins % 10 ≥ 6 → grades[main*2 + 1]   // the + tier
  else              → grades[main*2]        // base tier

GradeIndex(g) → 0..8  (errors on unknown string)
CPForGrade(g) = 100 + 50 * GradeIndex(g)  // 100..500
```

**Stat allocation** (`upsilontypes/entity/entitygenerator/`):

Each archetype declares `StatWeights` (HP, SP, MP, Attack, Defense, Movement, AttackRange). The generator distributes the CP pool proportionally to those weights, applying a minimum floor per stat. Higher grade → higher CP pool → better stats across the board.

**Skill bundle** (`archetype.BuildSkillBundle(grade, count)`):

Selects skills from the archetype's allowed tag set, filtered to tags appropriate for the grade. Count scales with grade index (Grade I: 2 skills, Grade V: 4 skills).

**Grade → CP mapping** (key boundary values):
| Wins | Grade | CP pool |
|------|-------|---------|
| 0    | I     | 100     |
| 10   | II    | 200     |
| 20   | III   | 300     |
| 30   | IV    | 400     |
| ≥40  | V     | 500     |

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_ai_progression_matching]]`
- **Key package:** `upsilontypes/entity/grade/` — grade table and CP formula
- **Key package:** `upsilontypes/entity/entitygenerator/` — stat generation from weights + CP
- **Key package:** `upsilonbattle/battlearena/controller/archetype/` — archetype stat weights + skill bundle
- **Bridge integration:** `upsilonapi/bridge/bridge_start.go` `generateEntityFromArchetype()`
- **PHP surface:** `battleui/app/Http/Controllers/API/MatchMakingController.php` — passes `total_wins` on AI player payload

## EXPECTATION
- `GradeFromWins(0)` = "I", `GradeFromWins(40)` = "V", `GradeFromWins(100)` = "V".
- `CPForGrade("I")` = 100, `CPForGrade("V")` = 500.
- Grade V entity generated from archetype has strictly higher HP than Grade I entity of the same archetype.
- Skill count scales: Grade I entities receive fewer skills than Grade V entities.
