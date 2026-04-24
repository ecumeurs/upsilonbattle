---
id: mechanic_backstab_detection_algorithm
status: STABLE
dependents:
  - [[mec_backstabbing_mechanic]]
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents:
  - [[domain_upsilon_engine]]
---

# New Atom

## INTENT
To implement backstab detection algorithm that determines when an attack originates from behind the target based on entity orientation, enabling 150% damage multiplier and 50% armor penetration bonuses.

## THE RULE / LOGIC
**Backstab Detection Algorithm:**

**Core Principle:**
A backstab occurs when an offensive action originates from behind the target entity's current facing direction. This position of tactical advantage bypasses defenses and maximizes impact.

**Orientation and Coordinate System:**
- **Facing Directions:** Entities are oriented in one of four cardinal directions: **Up**, **Right**, **Down**, or **Left**.
- **Vector Calculation:** The system determines the "Attack Vector" by calculating the angle between the attacker's coordinates and the target's coordinates.

**Detection Logic Hierarchy:**
1. **Target Orientation Check:** Retrieve the target's current facing direction.
2. **Back Definition:** The "Back" of the target is defined as the direction exactly 180° opposite to its facing.
3. **Angle Comparison:** The algorithm calculates if the Attack Vector falls within a **±45° tolerance** of the target's perfect back angle.
    - **Up Facing:** Backstab range is 135° to 225° (Down).
    - **Right Facing:** Backstab range is 225° to 315° (Left).
    - **Down Facing:** Backstab range is 315° to 45° (Up).
    - **Left Facing:** Backstab range is 45° to 135° (Right).

**Prerequisites for Resolution:**
- **Weapon Range:** Backstab bonuses only apply if the attacker is within their weapon's defined effective range (typically melee or close-range).
- **Line of Sight (LoS):** The path from the attacker to the target's back must be clear of opaque obstacles or walls.
- **Action Type:** In the current version (V2), backstab bonuses apply specifically to **Weapon-as-Skill** attacks; standard magical or area-of-effect skills do not trigger backstab detection.

**Mechanical Impact:**
- **Damage Multiplier:** Valid backstabs apply a **150%** multiplier to the final damage output.
- **Armor Bypass:** Backstabs ignore **50%** of the target's Armor Rating, making them highly effective against heavily armored "tanks."
- **Feedback:** Triggers specialized visual and auditory indicators, such as a distinct "Backstab!" combat log entry and unique damage animations.

**Tactical and AI Dynamics:**
- **Orientation Persistence:** Some entities may auto-face attackers when struck, preventing consecutive backstabs without further repositioning.
- **AI Awareness:** "Sneak" archetypes actively calculate paths to enter the target's back-angle range, while defensive archetypes attempt to position themselves to protect their rear arc.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[backstab_detection_algorithm]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entity.go` (orientation), `upsilonbattle/battlearena/ruler/rules/attack.go`
- **Integration:** Works with `mec_backstabbing_mechanic`, `armor_penetration_system`

## EXPECTATION
