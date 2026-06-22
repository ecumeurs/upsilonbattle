---
id: mechanic_character_creation_integration
human_name: "Character Creation Integration"
status: DRAFT
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents:
  - [[module_game]]
dependents: []
type: MECHANIC
---

# Character Creation Integration

## INTENT
To implement character creation and progression integration that bridges V1 legacy characters to V2 systems, handles data migration, and establishes new character creation flow with 100 CP point-buy allocation.

## THE RULE / LOGIC
**Character Creation and Progression Integration:**

**Core Principle:**
This mechanic establishes the standardized lifecycle for character data, encompassing the transition from legacy V1 structures to the robust V2 point-buy system and the subsequent progression through level-based milestones.

**V1 to V2 Migration Logic:**
- **Automatic Rebalancing:** Legacy characters are automatically migrated to the V2 baseline (HP 30-50, Attack 10, Defense 5, Movement 3).
- **CP Compensation:** Players receive a Character Point (CP) grant based on their legacy progression. The pool is reset to a baseline of 100 CP plus a calculated "migration bonus" to ensure no loss of perceived power.
- **Legacy Rewards:** Skill and equipment compensation is awarded based on the legacy character's level and win history, populating the new V2 inventories with appropriate-tier items.

**Character Creation Workflow:**
1. **Identity & Aesthetics:** Define character name (unique) and visual customization (avatar/colors).
2. **CP Allocation (Point-Buy):**
    - **Initial Pool:** Exactly 100 CP available.
    - **Interactive Spending:** Points are spent on core stats (HP, Attack, Defense, Movement) or exotic stats (Crit, Dodge) according to the V2 cost table.
    - **Validation:** Creation is only finalized if exactly 100 CP are allocated and all minimum attribute thresholds are met.
3. **Skill Selection:** The system offers a choice of one skill from a randomized pool of three Grade I-II candidates.
4. **Starter Equipment:** Players may select a basic set of equipment (Armor, Utility, Weapon) to fill their active slots immediately.

**Progression and Scaling:**
- **Level Increments:** Each level-up awards +10 CP for further attribute enhancement.
- **Temporal Milestones:**
    - **Every 5 Levels:** Access to skill reforging (modifying existing skill properties).
    - **Every 10 Levels:** Choice of one new skill from a Grade-appropriate randomized pool.
- **Exotic Stat Unlocking:** Higher-tier exotic stats (e.g., Accuracy, Jump Height) may be locked behind level prerequisites to prevent over-specialization in the early game.

**Data Consistency and Persistence:**
- **The V2 Schema:** Characters are represented by an expanded data structure that includes base/modified stats, available CP, exotic attribute maps, 3-slot equipment references, and credit balances.
- **Migration Traceability:** The system maintains a reference to the original legacy character ID for auditing and migration bonus verification.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_character_creation_integration]]`
- **Related Files:** Character creation API, migration scripts, database schema
- **API Endpoints:** `POST /api/v1/character/create`, `GET /api/v1/character/migrate`
- **UI Components:** Character creation wizard, migration notification interface

## EXPECTATION
