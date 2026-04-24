---
id: mech_skill_equipment_system
status: DRAFT
layer: IMPLEMENTATION
version: 2.0
tags: ["skills", "equipment", "progression"]
parents: []
dependents: []
priority: 5
---

# New Atom

## INTENT
To implement the skill equipment system where players select which skills from their inventory are active during battle. This is the bridge between skill inventory and actual combat usage, enforcing slot limits and preparation decisions.

## THE RULE / LOGIC
**Equipment System:**

**Preparation Phase (Pre-Battle):**
- Players select skills from inventory to equip before match starts
- Equipped skills appear in battle action panel
- Cannot equip more skills than available slots
- Equipment decisions are strategic (offensive vs defensive, utility)

**Equipment Operations:**
- **Equip Skill:** Move skill from inventory to equipped slot
- **Unequip Skill:** Return skill to inventory, free slot
- **Auto-Equip:** New skills from selection are auto-equipped if slots available

**Slot Allocation:**
- Skills have implicit slot positions (1, 2, 3, 4, 5)
- UI shows equipped skills with clear slot indicators
- Players can swap equipped skills by re-equipping

**Battle Integration:**
- Engine receives equipped skills as part of entity initialization
- Only equipped skills are registered via entity.RegisterSkill()
- Inventory skills are not available during battle

**Equipment Persistence:**
- Equipment state is persistent across matches
- Last equipped skill set is remembered
- Players can change equipment between matches

**Validation:**
- Cannot equip unowned skills
- Cannot equip duplicate instances of same skill
- Must respect slot limit from rule_character_skill_slots

## TECHNICAL INTERFACE

## EXPECTATION
