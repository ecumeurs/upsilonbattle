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
Same goes for defence, accuracy, and so on. 

## Entity Creation

[See Generator](entitygenerator/README.md)

## Entity Evolution

This is for later ;)

## Entity Storage

This is for later ;)


