---
id: entity_player
human_name: Player Account Entity
type: ENTITY
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[shared:requirement_customer_user_id_privacy]]
dependents:
  - [[upsilonapi:infra_seed_admin]]
  - [[upsilonapi:infra_seed_test_account]]
  - [[upsilonapi:rule_admin_access_restriction]]
  - [[upsilonapi:uc_admin_user_management]]
  - [[entity_users]]
---
# Player Account Entity

## INTENT
To define the player account entity: its identity attributes, registration and initial-setup rules, character roster rules, and lifetime statistics tracking.

## THE RULE / LOGIC
**Player Account Entity.**

**Core identity attributes:**
- `account_name` (Public/Unique)
- `full_address` (Private)
- `birth_date` (Private)
- `role` (Admin, Player)
- `id` (Internal UUID, NOT exposed to frontend per [[shared:requirement_customer_user_id_privacy]])
- **Deletion Protocol:** supports soft-deletion and anonymization as specified in [[upsilonapi:rule_gdpr_compliance]].

**Registration:** every player must connect through a logged-in account to play.

**Initial Setup:** upon account creation the player is automatically granted exactly **3 characters**.

**Character Rules Apply:** those granted characters must have their attributes rolled according to the rules defined in `[[upsilontypes:entity_character]]`.

**Statistics Tracking:** the player entity tracks the absolute number of game wins and losses (see `[[entity_users]]` for the persisted columns).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_player]]`
