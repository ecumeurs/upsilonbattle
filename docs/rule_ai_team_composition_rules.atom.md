---
id: rule_ai_team_composition_rules
status: STABLE
version: 1.0
parents: []
dependents: []
type: RULE
layer: BUSINESS
human_name: AI Team Composition Rules
tags: [ai,matchmaking,team,composition]
priority: 4
---

# New Atom

## INTENT
To implement AI team composition rules that enforce maximum limits on archetype variety per AI team, ensuring balanced team composition with appropriate support and specialist limitations.

## THE RULE / LOGIC
**Constraint:** Per AI team — maximum **1 Support**, maximum **1 Sneak**. Fighter and Ranger are uncapped.

**PHP enforcement** (`MatchMakingController::assignAIArchetypes(int $count): array`):
Iterates slots. Once `support` is picked it is removed from the available pool; same for `sneak`. Fighter and Ranger remain available throughout. Result: any 3-entity AI team always satisfies the constraint.

**Go enforcement** (`upsilonapi/bridge/validateTeamComposition(players []api.Player) error`):
Counts support and sneak archetypes per team across all entities in `AutoGen=true` players. Returns an error if either count exceeds 1 on the same team. Called during `StartArena` before entity generation.

**Violation response:** Go returns a 400-envelope error; StartArena is rejected. The PHP layer prevents this from ever being reached through normal matchmaking.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[rule_team_composition]]`
- **PHP:** `battleui/app/Http/Controllers/API/MatchMakingController.php` `assignAIArchetypes()`
- **Go:** `upsilonapi/bridge/bridge_start.go` `validateTeamComposition()`
- **Tests (PHP):** `PVEMatchmakingTest::test_ai_entities_are_auto_gen_with_archetype`
- **Tests (Go):** `bridge_start_archetype_test.go` `TestValidateTeamCompositionRejectsTwoSupports`, `TestValidateTeamCompositionRejectsTwoSneaks`

## EXPECTATION
- A normal PVE match never generates two Supports or two Sneaks on the AI team.
- A request with two Sneaks on the same team returns an error from `validateTeamComposition`.
- The same constrained archetype is allowed on different teams (team 1 support + team 2 support is valid).
