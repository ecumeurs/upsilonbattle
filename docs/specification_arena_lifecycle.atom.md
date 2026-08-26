---
id: specification_arena_lifecycle
status: DRAFT
type: SPECIFICATION
layer: ARCHITECTURE
priority: 5
tags: combat,lifecycle,arena
parents:
  - [[shared:us_take_combat_turn]]
human_name: Arena Lifecycle State Machine
version: 1.0
dependents: []
---

# New Atom

## INTENT
To make the arena's lifecycle explicit as a four-state machine (created, starting, in_progress, concluded) so that gateway code, the engine, and clients can all reason about which player actions are legal at any given moment, instead of that legality being an undocumented side effect of internal engine checks.

## THE RULE / LOGIC
1. **created** — the hub has written the match record (CreateMatchRecord) and is about to call, or has just called, the engine's arena-start endpoint. No engine-side arena exists yet. No player action is valid in this state; at most the client has received a match.found notification.
2. **starting** — the engine has registered the arena (its Ruler actor exists, with ArenaState == WaitingForControllers) but has not yet begun its game loop tick. Every player-initiated action — move, attack, skill use, pass, and forfeit alike — is rejected with 400 game.not.in.progress. This is the pre-tick window ISS-102 observed; the rejection is not a bug, it is this state's defined behavior, consistent with the fail-fast-on-illegal-state-transitions guarantee.
3. **in_progress** — the engine's game loop has started (ArenaState == InProgress, signaled by the game.started event). The full action set is legal: move, attack, skill use, pass, and forfeit, each subject to its own governing rule (forfeit specifically: [[upsilonbattle:rule_forfeit_battle]]).
4. **concluded** — the engine has resolved a win or forfeit condition (ArenaState == Finished). No further player action is accepted from this point. Only teardown proceeds: arena removal from the bridge and actor shutdown, per [[upsilonbattle:mechanic_arena_lifecycle]].
Transitions are strictly forward and one-way: created -> starting -> in_progress -> concluded, with concluded being terminal (teardown, not a fifth state).

## TECHNICAL INTERFACE
- **Hub side (created):** upsilonhub/internal/games/battle/matchmaking.go, CreateMatch — writes the match record, then calls the engine's arena-start endpoint (POST /v1/arena/start, [[upsilonapi:api_go_battle_start]]).
- **Engine side (starting / in_progress / concluded):** upsilonbattle/battlearena/ruler/ruler.go defines ArenaState (WaitingForControllers=1, InProgress=2, Finished=3) on the Ruler actor's CurrentState field.
- **Action guards:** upsilonbattle/battlearena/ruler/ruler_actions.go and ruler_turn.go reject any action with r.CurrentState != InProgress, replying game.not.in.progress — this is the enforcement point for the starting state's 'no actions' rule.
- **Concluded teardown:** governed separately by [[upsilonbattle:mechanic_arena_lifecycle]] (ArenaBridge.DestroyArena).

## EXPECTATION
- A move, attack, skill-use, pass, or forfeit request arriving while the arena is created or starting returns 400 game.not.in.progress, and this is compliant behavior, not a defect.
- The same request arriving once the arena is in_progress succeeds, subject to its own governing rule.
- No player action is accepted once the arena is concluded; only [[upsilonbattle:mechanic_arena_lifecycle]] teardown proceeds.
- The four states are reached strictly in order (created -> starting -> in_progress -> concluded); no state is skipped and none is re-entered.
