# .\battlearena\ruler\turner

[Up](../README.md)

# Turner 

This module is responsible for the turn logic of the game. It will be in charge of determining the order of the entities playing, and the amount of time they will have to wait before being able to play again.

Order of turn is based on delay remaining. The entity with the lowest delay remaining will be the first to play. Upon playing, all entities will have their delay reduced by the amount of time they had to wait.

The entity with the lowest delay remaining will be the first to play again, and so on.

## Delay Credit

Per movement: +200
Per Attack: +500

More later

## Timeout

It's also responsible for timeout detection and handling!

Timeout: + 1000 

ATM: expect a timeout of 60s.

