---
id: entity_users
human_name: Users Database Entity
type: ENTITY
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[data_persistence]]
  - [[entity_player]]
dependents:
  - [[entity_match_participants]]
---
# Users Database Entity

## INTENT
To define the comprehensive database structure for user accounts in the PostgreSQL schema, encompassing authentication, personal data, game statistics, and role-based access control.

## THE RULE / LOGIC
Defines the `users` table, serving as the primary persistent storage for player identity and metrics.

### Core Authentication & Identity
- `id` (UUID, Primary Key) - Internal unique identifier, not exposed to frontend per [[requirement_customer_user_id_privacy]]
- `account_name` (Varchar, Unique, Not Null) - Public unique display name
- `email` (Varchar, Unique, Not Null) - Contact email for account management
- `password_hash` (Varchar, Not Null) - Bcrypt/hashed password storage (never plaintext)

### Role-Based Access Control
- `role` (Varchar, Default 'Player') - Access level: 'Player' or 'Admin'
  - Enforced by [[rule_admin_access_restriction]]

### Game Statistics & Progression
- `total_wins` (Int, Default 0) - Lifetime match victories
- `total_losses` (Int, Default 0) - Lifetime match defeats
- `ratio` (Numeric, Default 0) - Calculated win/loss ratio for [[ui_leaderboard]]
- `reroll_count` (Int, Default 0) - Tracks character stat rerolls per [[rule_character_reroll_limits]]

### WebSocket Communication
- `ws_channel_key` (UUID, Unique, Nullable) - Secure channel identifier for Reverb WebSocket subscriptions
  - Generated during [[uc_player_registration]]
  - Required for [[module_frontend_matchmaking_orchestration]]

### Personal Data (GDPR Protected)
- `full_address` (Text, Nullable) - Residential address, private
- `birth_date` (Date, Nullable) - Date of birth, private
- Protected by [[rule_gdpr_compliance]] deletion protocols

### Account Lifecycle
- `remember_token` (Varchar, Nullable) - Laravel "remember me" functionality
- `created_at` (Timestamp) - Account creation timestamp
- `updated_at` (Timestamp) - Last modification timestamp
- `deleted_at` (Timestamp, Nullable) - Soft delete timestamp for GDPR compliance

### Database Constraints & Indexes
- **Unique Constraints**: `account_name`, `email`, `ws_channel_key`
- **Indexes**: `account_name`, `email`, `ws_channel_key`, `updated_at`
- **Role Constraint**: Must be 'Player' or 'Admin'

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_users]]`
- **Laravel Model:** `App\Models\User`
- **Foreign Key References:**
  - `characters.player_id` → `users.id` (CASCADE DELETE)
  - `match_participants.player_id` → `users.id` (CASCADE DELETE)

## INTEGRATION NOTES

### Laravel Session Management
- Integrates with Laravel's native session and password reset systems
- Supports [[uc_auth_login]] and [[uc_player_registration]]

### GDPR Compliance
- Soft delete via `deleted_at` allows account deactivation
- Full deletion requires anonymization per [[rule_gdpr_compliance]]

### WebSocket Integration
- `ws_channel_key` enables secure private channel subscriptions
- Required for real-time matchmaking notifications in [[module_frontend_matchmaking_orchestration]]

## TESTING CONSIDERATIONS
- Verify role constraints (Player/Admin only)
- Test cascade delete behavior with characters and match participants
- Validate ws_channel_key uniqueness and nullable behavior
- Ensure GDPR soft delete doesn't break foreign key constraints
