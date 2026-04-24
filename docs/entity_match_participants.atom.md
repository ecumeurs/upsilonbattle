---
id: entity_match_participants
human_name: Match Participants Database Entity
type: ENTITY
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[entity_game_match]]
  - [[entity_users]]
dependents: []
---
# Match Participants Database Entity

## INTENT
To define the database structure linking users (human or AI/bot) to specific game matches, enabling historical tracking, leaderboard calculations, and progression validation.

## THE RULE / LOGIC
Defines the `match_participants` table, serving as the junction table between `users` and `game_matches`.

### Core Match Participation
- `id` (BigInt, Primary Key, Auto-increment) - Internal unique identifier
- `match_id` (UUID, Foreign Key → `game_matches.id`) - Reference to the associated match
- `team` (Int, Not Null) - Team assignment: 1 or 2
- `status` (Varchar, Nullable) - Match outcome: 'WIN' or 'LOSS'
  - Enforced by CHECK constraint
  - Used for [[ui_leaderboard]] win/loss calculations

### Player Identification (Human vs AI/Bot)
- `player_id` (UUID, Foreign Key → `users.id`, **Nullable**) - Reference to user account
  - **Critical:** Nullable to support AI/bot participants in PvE modes
  - When NULL, represents AI-controlled opponents or system-generated entities
  - Required for human participant [[rule_progression]] calculations
  - Optional for AI participants who don't need progression tracking

### Match Metadata
- `created_at` (Timestamp) - When the participant joined the match
- `updated_at` (Timestamp) - Last modification timestamp (e.g., status update)

### Database Constraints & Indexes
- **Foreign Keys:**
  - `match_id` → `game_matches.id` (CASCADE DELETE)
  - `player_id` → `users.id` (CASCADE DELETE, Nullable)
- **Check Constraints:**
  - `status` must be 'WIN' or 'LOSS' when set
- **Indexes:** `match_id`, `player_id`

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_match_participants]]`
- **Laravel Model:** `App\Models\MatchParticipant`
- **Use Cases:**
  - Track human player performance for [[rule_progression]]
  - Support PvE modes with AI/bot opponents
  - Enable [[uc_admin_history_management]] for match audit
  - Power [[ui_leaderboard]] calculations

## INTEGRATION NOTES

### Human vs AI/Bot Distinction
The nullable `player_id` field enables flexible match participation:

**Human Participants (player_id NOT NULL):**
- Standard PvP and PvE matches
- Eligible for [[rule_progression]] rewards
- Contribute to [[ui_leaderboard]] rankings
- Generate match history for [[uc_admin_history_management]]

**AI/Bot Participants (player_id NULL):**
- PvE opponents controlled by game logic
- No progression or leaderboard impact
- Still require proper match completion tracking
- May be referenced in match history but excluded from player metrics

### Match Lifecycle Integration

**Match Start:**
- Participants created when [[api_go_battle_start]] is invoked
- `status` is NULL initially
- `team` assignment based on [[entity_game_match]] game mode

**Match Completion:**
- `status` updated to 'WIN' or 'LOSS' by battle engine
- Triggers [[rule_progression]] calculations for human participants
- Updates [[entity_users]] statistics via cascade operations

**Match History:**
- All participants retained for [[uc_admin_history_management]]
- AI participants included for completeness but excluded from player metrics

### Progression Validation
- Human participants validate against [[rule_progression]] constraints
- Check for duplicate participation (anti-cheat)
- Ensure proper win/loss attribution
- Support [[us_win_progression_win_alloc_point]] logic

## TESTING CONSIDERATIONS
- Verify NULL player_id behavior for AI/bot participants
- Test cascade delete from both matches and users tables
- Validate CHECK constraints for status field
- Ensure progression logic only applies to human participants
- Test team assignment consistency with game mode
- Verify match history includes both human and AI participants

## BUSINESS RULE CONNECTIONS
- Triggers [[rule_progression]] validation on match completion
- Supports [[us_win_progression_win_alloc_point]] for human participants
- Excludes AI participants from [[ui_leaderboard]] calculations
- Enables [[uc_admin_history_management]] for match audit trails
