# .\battlearena\entity\properties\effect\effectapplicator

[Up](../README.md)

# Effect Applicator

This module allow and effect to be enacted on a battle arena.

## Damaging Effect

That's an effect that intent to harm the target.

=> [See Properties Enum](../../propertyenum.go)

At time of writing it consist on the positive presence of:

* `Damage`
* `PoisonPower`
* `StunPower`

And the negative presence of:

* `ShieldPower`


### Hit!

First step is to determine what hits or not. At the moment, everything can be dodged, later expect to have different kind of dodges for different kind of attacks ( probably with the introduction of a `AttackType` property, or something similar)

The rule here is simple:

Roll 1d100, if result is lower than attackers's `Accuracy` - target's `Dodge` then it's a hit.

### Damage Calculation

There are several properties coming into play here:

| Property             | Description                              | Use                  | Def |
| -------------------- | ---------------------------------------- | -------------------- | --- |
| `Attack`             | Attacker's base attack                   | Flat damage          | 1   |
| `Damage`             | Attack multiplier (in %)                 | Damage Rate          | 100 |
| `PoisonPower`        | The base poison power of the attack      | Flat damage          | 0   |
| `StunPower`          | The base stun power of the attack        | Flat damage          | 0   |
| `ShieldPower`        | Direct damage to Shield                  | Flat damage          | 0   |
| `PoisonChance`       | The chance to poison the target          | Chance to poison     | 0   |
| `StunChance`         | The chance to stun the target            | Chance to stun       | 0   |
| `CriticalChance`     | The chance to deal critical damage       | Multiplier           | 0   |
| `CriticalMultiplier` | The multiplier to apply on critic (in %) | Global Multiplier    | 100 |
| `Armor`              | The armor of the target                  | Damage reduction     | 0   |
| `Shield`             | The consummable shield of the target     | Damage reduction     | 0   |
| `Defense`            | The natural defense of the target        | Damage reduction     | 0   |
| `HP`                 | The current HP of the target             | Ressource            | 10  |
| `Poison`             | The current poison level of the target   | Ressource            | 0   |
| `Stun`               | The current stun level of the target     | Ressource            | 0   |
| -------------------- | ---------------------------------------- | -------------------- | --- |

Expect Item to have an `Attack` Property at some point or something that add flat damage to the attack.

Poison true damage are computed thus: `PoisonPower` - `Defense`: It bypasses armor.
Stun true damage are computed thus: `StunPower` - `Armor`: It bypasses defense.

Critical Chance is computed easily: Roll 1d100, if result is lower than `CriticalChance` then it's a critical hit. Expect a counter to be added to the effect at some point. 

True damage are computed thus: (`Attack` * `Damage`) - `Armor` - `Defense` + True Poison + True Stun 

If a critical hit has been delivered, then the damage is multiplied by `CriticalMultiplier`, rounded down to int.

`ShieldPower` is a flat direct damage to the shield, it bypasses all defense and armor. It is applied before the damage application. (Expected to be negative here, for damaging effects)

Poison & Stun chance are simply rolled against 1d100. If the result is lower than the chance, then the effect(True Poison & True Stun) is added to the target's `Poison` & `Stun` properties.

[See Effects of poisonned state & Stunned state](../../../entity/README.md#Status-Effect-and-End-of-Turn)

### Damage Reduction

Damage reduction aims first at discounting the damage from the `Shield` property if any is available. If the shield is consumed, then the damage applied to the `HP` property.

## Healing Effect

That's an effect that intent to restore the target.

=> [See Properties Enum](../../propertyenum.go)

At time of writing it consist on the positive presence of:

* `Heal`
* `ShieldPower`

And the negative presence of:
* `PoisonPower`
* `StunPower`

As negative `PoisonPower` & `StunPower` are used to heal the target.

Notes: 

* Can't overheal
* Can Overshield up to MaxPV *2
* Can't reduce poison below 0 (no preventive cure)
* Can't reduce stun below 0 (no preventive cure)
