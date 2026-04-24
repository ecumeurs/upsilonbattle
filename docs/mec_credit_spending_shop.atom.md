---
id: mec_credit_spending_shop
human_name: Credit Spending Shop Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [economy, shop, progression]
parents:
  - [[domain_credit_economy]]
dependents:
  - [[mec_shop_inventory_system]]
---

# Credit Spending Shop Mechanic

## INTENT
To implement the shop system where players spend credits to purchase skills and equipment, with prices determined by Skill Weight system and equipment tiers.

## THE RULE / LOGIC
**Shop Pricing Formula:**
- **Skill Cost:** Total Positive SW × 2 credits
- **Equipment Cost:** Base cost × tier multiplier

**Skill Availability:**
- Shop offers skills based on character level
- Level 1-9: Grade I-II skills available
- Level 10-19: Grade II-III skills available
- Level 20-29: Grade III-IV skills available
- Level 30+: Grade IV-V skills available

**Equipment Categories:**
- **Armor:** Defensive items with ArmorRating
- **Utility:** Items with special effects and buffs
- **Weapon:** Items that transform basic attacks into skill-based attacks

**Purchase Mechanics:**
- Credits deducted from character balance
- Skill/equipment added to character inventory
- One-time purchases (no consumables in V2.1)
- Permanent character upgrades

**Shop Features:**
- Filter by category and grade
- Search by skill properties
- Preview skill effects and costs
- Affordability indicators based on current credits

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_credit_spending_shop]]`
- **API Endpoints:** `GET /api/v1/shop/inventory`, `POST /api/v1/shop/purchase`
