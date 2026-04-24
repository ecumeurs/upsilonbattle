---
id: entity_character_skill_inventory
status: DRAFT
tags: ["skills", "inventory", "progression"]
parents: []
dependents: []
priority: 5
layer: ARCHITECTURE
version: 2.0
---

# New Atom

## INTENT
To establish a character skill inventory system where players can store, manage, and access all acquired skills regardless of current equipment status. This enables roguelike-style skill collection and swapping.

## THE RULE / LOGIC
**Skill Inventory System:**

**Inventory vs Equipped:**
- **Inventory:** All skills character has acquired (unlimited capacity)
- **Equipped:** Subset of inventory selected for battle (slot-limited)
- Characters can own more skills than their available slots
- Skills in inventory are accessible for equipment/unequipment

**Acquisition Methods:**
- **Character Creation:** 1 selected skill enters both inventory and is auto-equipped
- **Skill Selection (every 10 levels):** New skill added to inventory
- **Shop Purchase:** Bought skills added to inventory
- **Reforging:** Modified skill remains in inventory

**Inventory Properties:**
- Skill ID and reference to skill template
- Acquisition timestamp (for sorting)
- Reforge count
- Current equipped status (boolean)

**Inventory Operations:**
- View all owned skills
- View skill details (properties, effects, cooldown)
- Mark skill as equipped/unequipped
- View skill statistics (usage, damage dealt, credits earned)

**Credit Assignment:**
- Credits tracked per skill
- Skill credits contribute to total character earnings
- Credit earning uses equipped/used skill

## TECHNICAL INTERFACE

## EXPECTATION
