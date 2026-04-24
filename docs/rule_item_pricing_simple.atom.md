---
id: rule_item_pricing_simple
status: DRAFT
priority: 5
layer: ARCHITECTURE
version: 2.0
tags: ["shop", "pricing", "items", "economy"]
parents: []
dependents: []
---

# New Atom

## INTENT
To establish a simple pricing model for shop items in V2 testing phase, using fixed credit costs instead of the full Skill Weight system. Prices are based on item type and stat bonuses, providing a predictable economy for initial testing.

## THE RULE / LOGIC
**Simple Pricing Formula:**
- **Armor Items:** Fixed 200 credits
- **Weapon Items:** Fixed 300 credits
- **Movement Items:** Fixed 150 credits
- **Utility Items:** Fixed 100 credits (future)

**Pricing Rationale:**
- Weapons provide damage bonus (highest value) = highest cost
- Armor provides defensive bonus = medium cost
- Movement provides mobility = lower cost
- No SW calculation yet (for V2.1 full procedural system)

**Property Value Mapping:**
- **ArmorRating (+1):** Not included in pricing (baseline stat boost)
- **WeaponRating (+5):** Core property, included in base cost
- **Movement (+1):** Core stat boost, included in base cost

**Price Adjustments:**
- Fixed costs allow predictable credit economy
- No complex math for V2 initial rollout
- Future: SW-based pricing for procedural item generation

**Purchase Validation:**
- Cannot purchase if player credits < item cost
- Cannot purchase duplicates of unique items
- Inventory capacity check (if applicable)

**Future Roadmap:**
- V2.1: SW-based procedural item pricing
- V2.2: Item rarity tiers (Common, Rare, Epic)
- V2.3: Item upgrade and enchantment system

## TECHNICAL INTERFACE

## EXPECTATION
