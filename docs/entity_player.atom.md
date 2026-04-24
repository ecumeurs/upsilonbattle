---
id: entity_player
human_name: Player Account Entity
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[requirement_customer_user_id_privacy]]
dependents:
  - [[entity_player_entity_character_rules_apply]]
  - [[entity_player_entity_player_initial_setup]]
  - [[entity_player_entity_player_registration]]
  - [[entity_player_entity_player_stats_tracking]]
  - [[entity_users]]
  - [[infra_seed_admin]]
  - [[rule_admin_access_restriction]]
  - [[uc_admin_user_management]]
---
# Player Account Entity

## INTENT
To aggregate the constituent rules of Player Account Entity.

## THE RULE / LOGIC
Initial setup and registration for player accounts.
Core attributes for identity:
- `account_name` (Public/Unique)
- `full_address` (Private)
- `birth_date` (Private)
- `role` (Admin, Player)
- `id` (Internal UUID, NOT exposed to frontend per [[requirement_customer_user_id_privacy]])
- **Deletion Protocol:** Supports soft-deletion and anonymization as specified in [[rule_gdpr_compliance]].

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_player]]`
