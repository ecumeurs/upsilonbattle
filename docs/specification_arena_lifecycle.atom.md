---
id: specification_arena_lifecycle
status: REVIEW
type: SPECIFICATION
layer: ARCHITECTURE
priority: 5
tags: combat,lifecycle,arena
parents:
  - [[shared:us_take_combat_turn]]
human_name: Arena Lifecycle State Machine
version: 1.1
dependents: []
---

# Arena Lifecycle State Machine

## INTENT
To make the arena's lifecycle explicit as a four-state machine (created, starting, in_progress, concluded) so that gateway code, the engine, and clients can all reason about which player actions are legal at any given moment, instead of that legality being an undocumented side effect of internal engine checks.

## THE RULE / LOGIC
1. **created** — the hub has written the match record (CreateMatchRecord) and is about to call, or has just called, the engine's arena-start endpoint. No engine-side arena exists yet. No player action is valid in this state; at most the client has received a match.found notification.
2. **starting** — the engine has registered the arena (its Ruler actor exists, with ArenaState == WaitingForControllers) but has not yet begun its game loop tick. Every player-initiated action — move, attack, skill use, pass, and forfeit alike — is rejected with 400 game.not.in.progress. This is the pre-tick window ISS-102 observed; the rejection is not a bug, it is this state's defined behavior, consistent with the fail-fast-on-illegal-state-transitions guarantee.
3. **in_progress** — the engine's game loop has started (ArenaState == InProgress, signaled by the game.started event). The full action set is legal: move, attack, skill use, pass, and forfeit, each subject to its own governing rule (forfeit specifically: [[upsilonbattle:rule_forfeit_battle]]).
4. **concluded** — the engine has resolved a win or forfeit condition (ArenaState == Finished). No further player action is accepted from this point. Only teardown proceeds: arena removal from the bridge and actor shutdown, per [[upsilonbattle:mechanic_arena_lifecycle]].

Transitions are forward and one-way along the ordinary path: created -> starting -> in_progress -> concluded, with concluded being terminal (teardown, not a fifth state). Exactly one exception to no state is skipped is permitted:

5. **created -> concluded (engine-start-failure exception, ISS-106)** — if the hub's call to the engine's arena-start endpoint never results in an acknowledged arena, the match is concluded immediately from created, skipping starting and in_progress entirely. This covers two cases: a transport error reaching the endpoint (the engine was unreachable), and a well-formed response that explicitly rejects the start (a Success:false envelope — the engine reached its rule checks and refused, rather than failing to respond). In both cases no engine-side arena was ever registered, so there is no starting/in_progress state for the match to have occupied; leaving the match record in created forever (a zombie match no client can ever progress or requeue) would be worse than concluding it directly. This mirrors the equivalent exception already established at poll-time: when a status poll finds the arena missing or unreachable from the engine after the match had reached in_progress, the match is likewise concluded immediately rather than left dangling (upsilonhub/internal/games/battle/matchmaking_status.go, verifyActiveMatch, the ISS-054 precedent) — the same engine-has-no-arena-for-this-match logic, applied here at creation-time instead of poll-time.

Outside this one named exception, transitions remain strictly forward and no other skip or reordering is legal: a match may not, for example, jump from starting to concluded without having reached in_progress, nor may concluded be re-entered once reached.

## TECHNICAL INTERFACE
- **Hub side (created):** upsilonhub/internal/games/battle/matchmaking.go, CreateMatch — writes the match record, then calls the engine's arena-start endpoint (POST /v1/arena/start, [[upsilonapi:api_go_battle_start]]).
- **Hub side (created -> concluded exception, ISS-106):** upsilonhub/internal/games/battle/matchmaking.go, CreateMatch's two failure branches on the StartArena call result — a transport error (ErrEngineUnreachable) and a well-formed Success:false rule-rejection envelope — both call Battle.ConcludeNow(ctx, matchID) to conclude the match immediately instead of leaving it in created. Mirrors the poll-time precedent in matchmaking_status.go's verifyActiveMatch (ISS-054).
- **Engine side (starting / in_progress / concluded):** upsilonbattle/battlearena/ruler/ruler.go defines ArenaState (WaitingForControllers=1, InProgress=2, Finished=3) on the Ruler actor's CurrentState field.
- **Action guards:** upsilonbattle/battlearena/ruler/ruler_actions.go and ruler_turn.go reject any action with r.CurrentState != InProgress, replying game.not.in.progress — this is the enforcement point for the starting state's no-actions rule.
- **Concluded teardown:** governed separately by [[upsilonbattle:mechanic_arena_lifecycle]] (ArenaBridge.DestroyArena), reached identically whether concluded was arrived at via the ordinary in_progress -> concluded path or via the created -> concluded exception above.

## EXPECTATION
- A move, attack, skill-use, pass, or forfeit request arriving while the arena is created or starting returns 400 game.not.in.progress, and this is compliant behavior, not a defect.
- The same request arriving once the arena is in_progress succeeds, subject to its own governing rule.
- No player action is accepted once the arena is concluded; only [[upsilonbattle:mechanic_arena_lifecycle]] teardown proceeds.
- Along the ordinary path, the four states are reached strictly in order (created -> starting -> in_progress -> concluded); no state is skipped and none is re-entered.
- Exception: if CreateMatch's call to start the arena fails — transport error or a Success:false rejection — the match transitions directly from created to concluded, never passing through starting or in_progress; this is the sole permitted skip.
