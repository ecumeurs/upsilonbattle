---
id: rule_character_skill_slots
status: DRAFT
parents: []
dependents: []
priority: 5
layer: ARCHITECTURE
version: 2.0
tags: ["skills", "progression", "equipment"]
---

# New Atom

## INTENT
To establish the character skill slot system where characters gain a new skill slot every 10 levels, starting with 1 base slot at level 1. This defines the maximum number of skills a character can have equipped for battle.

## THE RULE / LOGIC
**Skill Slot Formula:** Slots = 1 + (CharacterLevel / 10)

**Slot Progression:**
- **Level 1-9:** 1 skill slot (base)
- **Level 10-19:** 2 skill slots
- **Level 20-29:** 3 skill slots
- **Level 30-39:** 4 skill slots
- **Level 40+:** 5 skill slots (soft cap)

**Slot Validation Rules:**
- Characters cannot equip more skills than their available slots
- Equipped skills must be a subset of learned skills
- Removing an equipped skill frees a slot
- Adding an equipped skill requires an available slot

**Slot vs Learned Distinction:**
- **Learned Skills:** All skills the character has acquired (unlimited)
- **Equipped Skills:** Skills ready for use in battle (slot-limited)
- Character can learn more skills than they have slots
- Only equipped skills appear in battle action panel

## TECHNICAL INTERFACE
- Code Tag: @spec-link [[rule_character_skill_slots]]
- API Endpoints: POST /api/v1/character/{id}/equip-skill, DELETE /api/v1/character/{id}/unequip-skill
- Test Names: TestSkillSlotCalculation, TestEquipSkillExceedsSlots, TestUnequipSkillFreesSlot

## EXPECTATION
