# .\battlearena\entity

[Up](../README.md)

# Entity

Entities represent characters, monster and other things that get their own turn in the battle.

Entities may or may not all have the same properties available, some may be generic (like HP) and others may be specific to the entity type.

## Entity Properties

The following properties are available for all entities:

* `name` - The name of the entity

Other properties will be handled through a specific module, as they may or may not be added to the targeted entity at creation time.

Note: Most properties comes in PAIR (like, Attack and Defense). While its preferable for both entities in an attack to have both properties, it is not mandatory. If one of the entities is missing a property, the attack will be considered as a "defaulted" attack. It might mean the attack automatically miss, or whatever. This will be up to the rules.Attack to decide. 
Same goes for Defense, accuracy, and so on. 

## Status Effect and End of Turn

Status effects:
* `Poison` - If the poison level is greater than 0, then the poison effect is applied to the entity. The poison effect is applied by subtracting the poison level from the entity's HP. If the HP is reduced to 1, then the poison effect is ignored. **At the end of the turn** of the owner: it is then divided by 2, if the result is equals or lower than 1, it's removed.
* `Stun` - If the stun level reaches MaxHP/2 or more, **at the begining of the turn**, the entity is stunned and can't do anything, and the games skips the entity's turn. If stunned, the stun level is reset to 0. Otherwise it's divided by 2, if it reaches 1 or less, it's removed.

At the end of each entities turn:
* Available movement is reset to it's maximum value (and can be reduced/increased by other effects)
* HasMoved, HasActed flags are reset. 

## Entity Creation

[See Generator](entitygenerator/README.md)

## Entity Evolution

This is for later ;)

## Entity Storage

This is for later ;)


