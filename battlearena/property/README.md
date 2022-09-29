# .\battlearena\entity\properties

[Up](../README.md)

# Properties

This module generalise properties as a structured item.

## Creation

Creating new properties isn't really complex, just add either a new files under the `properties` folder, or complete the `properties` file with a new struct.
This struct will be stored along the entity and be used by rules to compute the result of an attack, or skill use or whatever.

Properties come in multiple kinds, but mostly two (at the moment)

* IntProperty: a property stored as an integer
* FloatProperty: a property stored as a float

The Property interface allow one to access informations about the property with a visibility level:

* `Public` : The property as visible by everyone  
* `ArenaObserver` : The property as visible by everyone in battle
* `ForeignController` : The property as visible by any foe
* `FriendlyController` : The property as visible by any friend
* `OwnController` : The property as visible by the owner of the entity
* `Analyser` : The property as visible by the analyser
* `ExpertAnalyst` : same, but another level of analysis
* `SpecialistAnalyst` : same, but another level of analysis
* `MasterAnalyst` : same, but another level of analysis
* `GameMaster` : The truest value of the property, as used by the rules.

Mostly, you will want to set value for name/get for Friendly and above: this will allow the controller to see the value of the property, and prevent foreigner to see it.

Lastly but not least there are two further methods associated with properties that aren't used as of now:

* Type: (None, Character, Skill, Item) Some properties won't be available for everything and will be restricted:
  * For example: Base Damage will be restricted to Character, Damage to Item, Attack to Skills. End damage will be computed by the rules and will probably look like this: `(BaseDamage + Damage) * Attack` see [Rules](../rules/README.md)
* Increase() will be set later on to handle property increase (level up)
  
## Defined properties

Character properties:
* Int:
  * Attack
  * Defense
  * Movement
  * AttackRange
  * JumpHeight
  * HP

## TODO


* At some point one should probably add temporary value to such properties(like a buff or a debuff) and a way to handle them.
* HP should be able to store a max value and a current value!