---
id: rule_credit_earning
human_name: Combat Credit Earning Rule
type: RULE
layer: ARCHITECTURE
version: 2.0
status: STABLE
priority: 5
tags: [economy, credits, combat, status-effects, support]
parents:
  - [[upsilonapi:domain_credit_economy]]
dependents: []
---

# Combat Credit Earning Rule

## INTENT
Define how characters earn credits during combat — from damage dealt, healing performed, damage mitigated, and status effects applied — with all credits assigned to the caster of the effect.

## THE RULE / LOGIC
**Base Rule:** 1 HP of absolute effect = 1 credit.

**Damage Credits:**
- Credits equal the actual HP lost by the target (after shields absorb).
- Example: 15 damage vs 0 shield = 15 credits; 15 damage vs 10 shield = 5 credits.
- Applies to direct attacks and damaging skills.

**Healing Credits:**
- Credits equal the actual HP recovered by the target (capped by MaxHP).
- Example: heal 15 HP to a target missing 5 HP = 5 credits.

**Damage Mitigation Credits (Support):**
- 1 HP mitigated (shields, damage reduction) = 1 credit, awarded to the mitigation's caster.

**Status Effect Credits (Flat Rate):**
- Credits Earned: `SkillWeight / 10` per application (poison, stun, buff, debuff).
- Example: 100 SW = 10 credits; 200 SW = 20 credits; 50 SW = 5 credits.
- Awarded once at the moment of application — no ongoing per-turn credits.
- Each application/stack earns separately; multi-effect skills apply the formula per component.

**Credit Assignment & Tracking:**
- Credits always go to the effect's caster (attacker, healer, or shielder), never the target.
- Effects must track `CasterID` until the effect ends; shields keep earning for the caster even after the caster dies.
- Tracked per character per match; balance synced to the central user account via webhooks.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_credit_earning]]`
- **Test Names:** `TestDamageCreditEarning`, `TestHealingCreditEarning`, `TestShieldCreditEarning`, `TestPoisonCreditEarning`, `TestStunCreditEarning`, `TestBuffCreditEarning`
- **Assignment Logic:**
```go
// When shield blocks damage
if shield.BlocksDamage > 0 {
    credits.Earned += shield.BlocksDamage
    credits.AssignedTo = shield.CasterID  // Original caster
}
```

## EXPECTATION (For Testing)
- Damage, healing, and mitigation each award credits equal to the absolute HP affected, assigned to the caster.
- Status effects award `SkillWeight / 10` once per application, with no per-turn credits.
- Shields continue crediting their caster after the caster's death; credits are never assigned to the target.
