---
id: mech_ai_name_generation
human_name: "AI Entity Name Generation"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 3
tags: [ai, names, aesthetic, generation]
parents:
  - [[req_ui_look_and_feel]]
dependents: []
---

# AI Entity Name Generation

## INTENT
To dynamically generate "haxxor" style, gritty, and sci-fi names for AI players and their combat entities to enhance atmospheric immersion.

## THE RULE / LOGIC
Names are generated using a "Segmented Concatenation" pattern:
- **Pattern A (Haxxor):** `[Prefix][Subject][Suffix]`
- **Pattern B (Industrial):** `[Technical_Noun]-[ID]`
- **Pattern C (Abstract):** `[Abstract_Noun]_[Hex_Code]`

**Dictionaries:**
- **Prefixes:** `Null_`, `Void_`, `DeathX`, `Rust_`, `Cyber`, `Neon`, `Ghost_`, `Cinder`
- **Subjects:** `Vermin`, `Proxy`, `Ghost`, `Core`, `Code`, `Glitch`, `Zero`, `One`
- **Suffixes:** `_X`, `v2`, `_Bot`, `_666`, `_Alpha`, `_Z`
- **Technical Nouns:** `Scrap`, `Static`, `Sludge`, `Terminal`, `Node`, `Array`
- **Abstract Nouns:** `Fracture`, `Desolation`, `Entropy`, `Echo`

**Exclusions:**
- Do not use generic terms like "Bot 1" or "Computer".
- Names must not exceed 20 characters.

## TECHNICAL INTERFACE (The Bridge)
- **Logic Location:** `\App\Models\Character::generateAIName()` and `MatchMakingController`
- **Code Tag:** `@spec-link [[mech_ai_name_generation]]`
- **Input:** Seed or Entity Type
- **Output:** String (The Name)

## EXPECTATION (For Testing)
- Generated names should be unique within a single match.
- Names should strictly follow the "Neon in the Dust" aesthetic defined in [[req_ui_look_and_feel]].
- 100% of generated names must match one of the three patterns (Haxxor, Industrial, Abstract).
