# .\battlearena\ruler\rules

[Up](../README.md)

# Rules

This module is responsible for specific point of rules used during the battle.
The two main rules begins with movement and attack. They will further fork as the game logic evolves.

## GameState

Structs and helper functions around the game state.

## Movement

Movement is the action entities uses to move from tile to tile.

### Local Ctx

this local context is used to store the current state of the movement.
It's mostly useless right now but expect further use in the future, especially when traps and triggered actions will come into play.

```go
type localMovekCtx struct {
	*GameState
	log *logrus.Entry
}
```

### v0.0.1

There are but few restriction:

* Height difference between two tiles must be less than 3.
* Entities can only move on `Ground` tiles.

## Attack

Attack is the action entities uses to attack other entities.

### Local Ctx

Local ctx is used to transit information during an attack round. At the moment, it's mostly empty but expect it to carry more information in the future, especially when attacks will trigger combos and movements.

```go
type localAttackCtx struct {
	*GameState
	log *logrus.Entry
}
```


### v0.0.1

Attack are death sentences: the attacked entity will be removed from the game.
Attack reach is 1 tile away from the attacker.