---
id: mechanic_backstab_detection_algorithm
human_name: "Backstab Detection Algorithm"
status: STABLE
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents:
  - [[upsilonapi:domain_upsilon_engine]]
---

# Backstab Detection Algorithm

## INTENT
To implement backstab detection algorithm that determines when an attack originates from behind the target based on entity orientation, enabling 150% damage multiplier and 50% armor penetration bonuses.

## THE RULE / LOGIC
**Backstab Detection Algorithm:**

**Core Principle:**
A backstab occurs when an offensive action originates from behind the target entity's current facing direction. This position of tactical advantage bypasses defenses and maximizes impact.

**Orientation and Coordinate System:**
- **Facing Directions:** Entities are oriented in one of four cardinal directions: **Up**, **Right**, **Down**, or **Left** (existing EntityOrientation system).
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
- **Action Type:** In the current version (V2), backstab bonuses apply specifically to **Weapon-as-Skill** attacks; standard magical or area-of-effect skills do not trigger backstab detection. Future skills may carry a "Backstab Enabled" property.

**Mechanical Impact:**
- **Damage Multiplier:** Valid backstabs apply a **150%** (1.5×) multiplier to the final damage output.
- **Armor Bypass:** Backstabs ignore **50%** of the target's Armor Rating, making them highly effective against heavily armored "tanks."
- **Shield:** Shield still applies fully (it is NOT penetrated by a backstab).
- **Damage Formula:** `(BaseDamage × 1.5) - (ArmorRating × 0.5) - Shield`.
- **Weapon Coverage:** All weapons (bows, pistols, melee) support backstab.
- **Feedback:** Triggers specialized visual and auditory indicators, such as a distinct "Backstab!" combat log entry, damage-multiplier display, critical backstab highlighting, and unique damage animations.

**Tactical and AI Dynamics:**
- **Orientation Persistence:** Some entities may auto-face attackers when struck, preventing consecutive backstabs without further repositioning.
- **AI Awareness:** "Sneak" archetypes actively calculate paths to enter the target's back-angle range and prioritize backstab opportunities; defensive/support archetypes attempt to position themselves to protect their rear arc and protect vulnerable allies.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mechanic_backstab_detection_algorithm]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entity.go` (orientation), `upsilonbattle/battlearena/ruler/rules/attack.go`
- **Test Names:** see `upsilonbattle/battlearena/ruler/rules/backstab_test.go`
- **Integration:** Works with `[[mechanic_armor_penetration_system]]`

## EXPECTATION
- An attack whose Attack Vector falls within the ±45° back-arc of the target's facing triggers backstab: 150% (1.5×) damage and 50% armor penetration, per `(BaseDamage × 1.5) - (ArmorRating × 0.5) - Shield`.
- An attack from outside the back-arc resolves as a normal (non-backstab) attack with no multiplier and no armor bypass.
- Shield value is subtracted in full on a backstab (shield is never penetrated).
- All weapon types (bow, pistol, melee) can trigger backstab; non-weapon skills do NOT trigger backstab in V2.
- Backstab bonuses apply only when the attacker is within weapon range AND line of sight to the target is clear; otherwise no backstab bonus.
- Each cardinal facing maps to the documented back-arc range (e.g. Up facing → backstab arc 135°–225°).
- Verified by `upsilonbattle/battlearena/ruler/rules/backstab_test.go`.
