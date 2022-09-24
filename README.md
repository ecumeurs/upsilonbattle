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

[Controller Methods)](https://github.com/ecumeurs/upsilonbattle/blob/main/battlearena/controller/controllermethods/controllermethods.go)


## Tests

[Ruler Tests](battlearena/ruler/ruler_fullgame_test.go) This is the most appropriate test as of now.

[Lesser Ruler Tests](battlearena/ruler/ruler_test.go) This is probably broken.
