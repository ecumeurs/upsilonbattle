# .\battlearena\controller

[Up](../README.md)

## Controller Overview

Controllers are responsible for managing entities within the battle arena. They communicate with the central `Ruler` via the actor model and asynchronous message passing. The controller processes notifications (such as `BattleStart`, `ControllerNextTurn`, `EntitiesStateChanged`) and sends commands/requests (like `ControllerMove`, `ControllerAttack`, `EndOfTurn`) to the ruler.

## Testing and "Stoppers"

