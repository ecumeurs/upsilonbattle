---
id: mec_ai_archetype_system
human_name: AI Archetype System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [ai, archetypes, progression]
parents: []
dependents: []
---

# AI Archetype System Mechanic

## INTENT
To implement four distinct AI archetypes (Fighter, Ranger, Support, Sneak) that follow player progression rules with archetype-specific stat allocation and skill selection.

## THE RULE / LOGIC
**Four AI Archetypes:**

**Fighter Controller:**
- **Skills:** High damage melee skills, defensive skills, charge abilities
- **Stat Priority:** Attack > Defense > Movement > HP
- **Tactics:** Direct approach, aggressive positioning, focus on damage
- **Suitable For:** Frontline combat, tank roles, aggressive play

**Ranger Controller:**
- **Skills:** Ranged damage skills, trap placement, movement/positioning skills
- **Stat Priority:** Attack > Movement > Accuracy > Defense
- **Tactics:** Kiting behavior, maintain range, use terrain advantages
- **Suitable For:** Ranged combat, positioning, tactical play

**Support Controller:**
- **Skills:** Healing, shielding, buff/debuff application, ally positioning
- **Stat Priority:** MP/SP > Defense > HP > Attack
- **Tactics:** Stay near allies, protect weak team members, prioritize healing
- **Suitable For:** Team support, defensive play, ally protection

**Sneak Controller:**
- **Skills:** Backstabbing bonuses, movement skills, poison/stun application, stealth
- **Stat Priority:** Movement > Attack > Dodge > CritChance
- **Tactics:** Flanking, target weak enemies, avoid direct confrontation
- **Suitable For:** Assassin play, backstabbing, tactical positioning

**Progression System:**
- **Same Rules:** +10 CP per win, max 100 + (total_wins × 10)
- **Skill Selection:** Choose from archetype-appropriate skill pools
- **Level Matching:** AI level matches average player level
- **Stat Allocation:** Archetype-specific priorities for point distribution

**Team Composition Rules:**
- Maximum 1 Support per AI team
- Maximum 1 Sneak per AI team
- Ensures balanced team composition

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_ai_archetype_system]]`
- **Controllers:** `FighterController`, `RangerController`, `SupportController`, `SneakController`
