---
id: entity_shop_item
status: DRAFT
type: ENTITY
layer: ARCHITECTURE
priority: 5
version: 2.0
tags: [shop, items, catalog, iss-074]
human_name: Shop Item Catalog Entity
parents:
  - [[mec_credit_spending_shop]]
  - [[upsilonapi:domain_credit_economy]]
dependents:
  - [[api_shop_browse]]
  - [[entity_player_inventory]]
---

# New Atom

## INTENT
To define the shop catalog entity — the immutable list of items players can purchase with credits. One row per purchasable item; the V2.0 seed contains exactly three rows (Basic Armor, Basic Sword, Swift Boots).

## THE RULE / LOGIC
**Shop Item Entity:**

**Core Fields (`shop_items` table):**
- **id:** UUID primary key (deterministic in V2.0 seed for test stability).
- **name:** Varchar display name (e.g. "Basic Armor").
- **type:** Varchar item category — one of `armor`, `weapon`, `movement`, `utility` (legacy categorisation used in catalog browsing UX; redundant with `slot` for V2.0).
- **slot:** Varchar with CHECK constraint `slot IN ('armor','utility','weapon')`. The 3-slot binding target — see `[[mec_three_slot_equipment_system]]`.
- **properties:** JSON map keyed by engine `ItemProperties` enum (e.g. `{"ArmorRating":5}`, `{"WeaponBaseDamage":5,"WeaponType":"One-Handed Melee","WeaponRange":1}`). Materialised at battle init via the engine's `def.ItemProperty()` factory; see `[[mech_item_buff_application]]`.
- **cost:** Integer credit cost per unit. V2.0 fixed values per `[[rule_item_pricing_simple]]`.
- **available:** Boolean (default `true`). When `false`, the item is hidden from `GET /v1/shop/items` but existing inventory rows remain valid.
- **version:** Varchar tagging the catalog generation (default `'2.0'`); enables future migrations to retire / version items.

**V2.0 Seed (deterministic):**
- Basic Armor — slot=armor, properties={ArmorRating:5}, cost=200.
- Basic Sword — slot=weapon, properties={WeaponBaseDamage:5, WeaponType:"One-Handed Melee", WeaponRange:1}, cost=300.
- Swift Boots — slot=utility, properties={Movement:1}, cost=150.

**Lifecycle:**
- Seeded on `php artisan db:seed`; not user-mutable.
- Admin-only mutation deferred to V2.1+ (no admin endpoint in V2.0).

**Privacy:**
- Catalog is fully public to any authenticated user.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[entity_shop_item]]`
- **Laravel Model:** `App\Models\ShopItem`
- **Migration:** `*_create_item_system_tables.php`
- **Seeder:** `Database\Seeders\ShopItemsSeeder`
- **API:** `[[api_shop_browse]]`

## EXPECTATION
- `php artisan migrate:fresh --seed` produces exactly three rows.
- Deterministic UUIDs allow test scenarios to reference items by ID.
- `slot` value is enforced by CHECK; invalid inserts fail.
- `cost` is a positive integer; `properties` deserialises to a valid item property map.
