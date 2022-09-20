# .\battlearena\ruler

[Up](../README.md)

# Ruler

This module is the main game logic. It's the one that decides who wins and who loses, and it's the one that validate any input from the controllers.
It also handle the turn logic, and the game state.

## Game State 

The game state contains a simple enum that can be one of the following:

* `WaitingForControllers`: the game is waiting for players to join the game.
* `InProgress`: the game is still in progress, no one has won yet.
* `Finished`: the game is finished, someone has won.

The first state `WaitingForControllers` is the default state, and it's the one that the game starts with. The game will automatically transition to the `InProgress` state when all controllers have been added to the game. The game will automatically transition to the `Finished` state when the game is over.

During the first state, the game will only accept new controller joining the game, and refuse any command received.
During the second state, the game will accept command received from controllers, and will process them. Although it will only accept commands on the behalf of the current entity whose turn it is.
During the third state, the game will refuse any further command received, except to retrieve the winner of the game and the endgame state.

The Ruler also is in possession of the Grid data and Entities living within.

It will be the responsibility of the master controller to remove this module from the memory.

## Turn Logic

Every entity joining the fight will be assigned a delay credit. This delay credit is the amount of time that the entity will have to wait before being able to act.

The first entity to play it the one with the lowest credit. Upon begining its turn, all remaining entities will have their credit reduced by as much.

The entity with the lowest credit will be the first to play again, and so on.

Per turn, each entity will be able to perform a move and an action.
Both of these will have an impact on the delay credit of the entity, determining their next action opportunity. Controllers will have to take this into account when sending commands to the entities.

## Game Rule Sets

(TODO)
These modules will handle the application of the game rules. They will be responsible to alter entities and map states.

These will also influence delay credits.
ATM: 
* Per movement: +200
* Per attack: +500
* Per timeout: +1000
* End of turn: +500

Expect attacks to be variable in the future. Expect movement to be affected by the terrain. Expect delay of other entities to be affected by entities techniques.

## Controller Validators

These modules will handle the validation of commands received from controllers. They ensure the message received are valid and can be processed by the game.




