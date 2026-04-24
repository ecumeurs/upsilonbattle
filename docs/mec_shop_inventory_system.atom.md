---
id: mec_shop_inventory_system
human_name: Shop Inventory System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [shop, economy, inventory]
parents:
  - [[mec_credit_spending_shop]]
dependents: []
---

# Shop Inventory System Mechanic

## INTENT
To implement the shop inventory system that manages available skills and equipment for purchase, filtering by character level and grade, with affordability checking and purchase management.

## THE RULE / LOGIC
**Inventory Categories:**
- **Skills:** All available skills organized by grade and properties
- **Equipment:** Armor, Utility, Weapon items organized by tier and type
- **Filters:** Property-based filtering, grade-based sorting, cost-based sorting

**Skill Inventory Management:**
- **Level-Based Availability:** Filter skills by character level access
- **Grade Organization:** Skills grouped by I-V grades
- **Property Search:** Find skills by specific properties (damage, healing, range, etc.)
- **Random Selection:** Generate 3 random skills from appropriate pool for selection

**Equipment Inventory Management:**
- **Slot-Based Filtering:** Show only equipment for available slots (armor, utility, weapon)
- **Tier Organization:** Equipment grouped by power tiers
- **Stat Property Search:** Find equipment by specific stat bonuses
- **Affordability Checking:** Highlight items affordable with current credits

**Purchase Management:**
- **Credit Deduction:** Deduct cost from character credit balance
- **Item Acquisition:** Add purchased item/skill to character inventory
- **Inventory Limits:** Enforce equipment slot limits and skill capacity
- **Purchase History:** Track all purchases for audit and player reference

**Shop Features:**
- **Filters:** Filter by category, grade, property, price range
- **Sorting:** Sort by name, price, grade, relevance
- **Preview:** Show detailed skill/item properties before purchase
- **Search:** Text search for skills/items by name or property

**Inventory Database Schema:**
- **shop_skills:** Available skills with properties, grades, costs
- **shop_equipment:** Available equipment with properties, tiers, costs
- **character_purchases:** Purchase history and owned items

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_shop_inventory_system]]`
- **API Endpoints:** `GET /api/v1/shop/skills`, `GET /api/v1/shop/equipment`, `POST /api/v1/shop/purchase`
