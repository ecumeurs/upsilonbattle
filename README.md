# Upsilon Battle

T RPG

# Version 0.0.3 ( WIP )

Skill usage, limiters (1 action per turn, n move, etc), skill cooldown.
Buffs and Curses 

## Objectives 

* At each turn, mvt,attack counter must be reset. (but not HP)
* Skills:
  * An action that have its own properties. 
  * Attacks using skill will use both set of properties. 
  * Skill have cooldown. 
* Buff and curses are applied to properties for a certain number of rounds.

# Version 0.0.2 

Basic Entity Stats, and rules.

## Objectives

* Entities have stats:
  * HP: takes more hits to be killed.
  * Defence: reduce incoming damages
  * Attack: increase dealt damages
  * Attack Number: Number of attack allowed per turn. (default 1)
  * MVT: Number of cell/turn that can be moved through
  * Jump: Maximum height difference between two adjascent cell for move (defaulted to 2 up to now) 
  * Attack Range (min/max): attack range (right now we are still at very basic attacks)
 

* Ruler must ensure an entity can't move further than expected. (or higher)

* damage computation: attack - defence (simple)

* ruler must ensure the attack range is respected: for example: attack range set to 4-8 won't be able to attack at cqc

* Of course when HP drop to 0: entity is dead and removed from the game (as previously)

* Generated Entities will have random values for their starting stats! (classlike stuff come later)


## Docs

[Properties](battlearena/entity/properties/README.md)

[Rules](battlearena/ruler/rules/README.md)



# Version 0.0.1

## Objectives

* Grid generation
Generate a randomish map based on specs. (Flat lands, Hills, River)

* Battle Start
2 controller need to register themselves

* Turn based game
Each controllers is granted a number of entities on the board. Each entity will be called in turn based on a random ranking selection. 
Each entity is granted two options: move and attack. (at the moment, one entity can do 10 move actions and 10 attack in the same round, no checks is being done ...)
Lastly, an entity (through the controller) can declare its end of turn. The ruler will then declare the next entity's turn

* Battle End
When only one controller's entities remain, the battle is considered done, a winner is declared.
All controller can quit the battle arena then.

## Docs

[Ruler](battlearena/ruler/README.md)

[MessageQueue](https://github.com/ecumeurs/upsilontools/tree/main/tools/messagequeue)

[Actor](https://github.com/ecumeurs/upsilontools/tree/main/tools/actor)

[Ruler Methods](https://github.com/ecumeurs/upsilonbattle/blob/main/battlearena/ruler/rulermethods/rulermethods.go) 
message the ruler expect to receive and struct he will respond with. It also can broadcast some of these as well.

[Controller Methods](https://github.com/ecumeurs/upsilonbattle/blob/main/battlearena/controller/controllermethods/controllermethods.go)


## Tests

[Ruler Tests](battlearena/ruler/ruler_fullgame_test.go) This is the most appropriate test as of now. (valid v0.0.2)

[Lesser Ruler Tests](battlearena/ruler/ruler_test.go) This is probably broken.
