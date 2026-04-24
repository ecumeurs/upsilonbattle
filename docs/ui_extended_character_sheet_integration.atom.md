---
id: ui_extended_character_sheet_integration
status: DRAFT
version: 1.0
parents: []
dependents: []
type: UI
layer: CUSTOMER
---

# New Atom

## INTENT
To implement extended character sheet UI integration that displays all V2 properties, skills, equipment, and progression information, enabling players to understand and manage their complete character capabilities.

## THE RULE / LOGIC
**Extended Character Sheet Integration:**

**Core Principle:**
The character sheet serves as the centralized dashboard for player identity, providing a comprehensive overview of V2 statistics, equipment configurations, and progression milestones.

**Information Architecture:**
- **Identity & Economy:** Displays the character's name, current level, total wins, and real-time credit balance. Includes the total Character Point (CP) pool (base 100 + win-based accumulation).
- **Core Attribute Matrix:**
    - **Health (HP):** Current/Max health bar with percentage markers.
    - **Offense (Attack):** Aggregated value from base stats and weapon bonuses.
    - **Mitigation (Defense):** Combined rating from base stats and armor pieces.
    - **Mobility (Movement):** Maximum movement tiles per turn.
- **Exotic Attributes:** Secondary stats including Critical Chance, Critical Multiplier, Dodge/Evasion, Accuracy, and Jump Height.

**Combat and Tactical Breakdown:**
- **Capabilities Section:** Details offensive modifiers like Attack Range, Backstab Multipliers, and Armor Penetration values.
- **Skill Repository:** Lists all owned skills with their associated Grade, Skill Weight (SW), resource costs (MP/SP), and current cooldown timers.

**Equipment Management:**
- **Active Slot Interface:** A visual 3-slot layout for **Armor**, **Utility**, and **Weapon** categories.
- **Rarity Visuals:** Equipment items are color-coded based on their tier (e.g., Common, Uncommon, Rare).
- **Inventory Integration:** Displays unequipped items available for swapping, with stat comparison tooltips showing projected changes.

**Progression and Milestones:**
- **CP Allocation:** Interactive interface for spending available points on base or exotic stats according to the cost table.
- **XP Visualization:** Progress bar indicating advancement toward the next level and upcoming rewards (e.g., skill selection unlocks).
- **Historical Tracking:** Logs recent stat changes and transaction history for transparency.

**Interactive Features and UI Logic:**
- **Comparison Engine:** High-contrast indicators (green/red) for stat differences when comparing currently equipped gear with inventory candidates.
- **Detailed Modals:** Expanded views for skills and equipment providing precise mathematical formulas and targeting previews on the grid.
- **Performance:** Implements progressive data loading and client-side caching to ensure the interface remains responsive during high-frequency attribute updates.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[extended_character_sheet_integration]]`
- **Related Files:** Character sheet UI components, stat calculation APIs, progression tracking systems
- **Integration:** Works with `character_point_buy_system`, `entity_equipment_system`, `api_equipment_management`

## EXPECTATION
