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
  - [[shared:requirement_customer_user_id_privacy]]
dependents:
  - [[shared:uc_admin_user_management]]
  - [[upsilonapi:infra_seed_admin]]
  - [[upsilonapi:rule_admin_access_restriction]]
  - [[entity_player_entity_character_rules_apply]]
  - [[entity_player_entity_player_initial_setup]]
  - [[entity_player_entity_player_registration]]
  - [[entity_player_entity_player_stats_tracking]]
  - [[entity_users]]
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
