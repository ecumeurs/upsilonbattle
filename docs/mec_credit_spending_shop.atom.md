---
id: mec_credit_spending_shop
human_name: Credit Spending Shop Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.1
status: DRAFT
priority: 5
tags: [economy, shop, progression]
parents:
  - [[upsilonapi:domain_credit_economy]]
dependents:
  - [[upsilonapi:api_shop_browse]]
  - [[upsilonapi:api_shop_purchase]]
  - [[mec_shop_inventory_system]]
  - [[upsilontypes:entity_shop_item]]
---

# Credit Spending Shop Mechanic

## INTENT
To implement the shop system where players spend credits to purchase skills and equipment, with prices determined by Skill Weight system and equipment tiers.

## THE RULE / LOGIC
**V2.0 Fixed Pricing (this issue, ISS-074):**
- Per `[[rule_item_pricing_simple]]`, V2.0 catalog entries carry an explicit `cost` column populated at seed time. There is no procedural pricing in V2.0.
- Reference catalog (V2.0):
  - Basic Armor — 200 credits — slot=armor — properties={ArmorRating:5}
  - Basic Sword — 300 credits — slot=weapon — properties={WeaponBaseDamage:5, WeaponType:"One-Handed Melee", WeaponRange:1}
  - Swift Boots — 150 credits — slot=utility — properties={Movement:1}

**V2.1+ Procedural Pricing (deferred):**
- Skill cost = Total Positive Skill Weight × 2 credits.
- Equipment cost = Equipment Power Rating × Tier Multiplier (Common 1.0× → Legendary 5.0×).
- Reintroduced when the procedural item generator lands.

**Purchase Mechanics (V2.0):**
- Endpoint: `[[api_shop_purchase]]` — `POST /v1/shop/purchase` body `{ shop_item_id, quantity? }`.
- Service: DB-transactional. Steps:
  1. Lock user row, check `users.credits >= shop_items.cost × quantity`.
  2. Insufficient → 422 with `meta.reason = "insufficient_credits"`.
  3. Debit credits.
  4. Upsert `player_inventory` row (increment quantity if exists, capped at 99 per `[[rule_quantity_cap]]`).
  5. Insert `inventory_transactions` audit row (transaction_type=`purchase`).
  6. Insert `credit_transactions` audit row (source=`shop_purchase`).
- **Crash early:** any service-level failure rolls back the transaction.

**Shop Browse (V2.0):**
- Endpoint: `[[api_shop_browse]]` — `GET /v1/shop/items` returns all `available=true` rows. No filtering / pagination in V2.0 (3-item catalog).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_credit_spending_shop]]`
- **API Endpoints:** `[[api_shop_browse]]`, `[[api_shop_purchase]]`
- **Service:** `App\Services\ShopService::purchase` (Laravel)
- **Migration:** `*_create_item_system_tables.php` (creates `shop_items`, `player_inventory`, `inventory_transactions`)
