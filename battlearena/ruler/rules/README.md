# .\battlearena\ruler\rules

[Up](../README.md)

# Rules

This module is responsible for specific point of rules used during the battle.
The two main rules begins with movement and attack. They will further fork as the game logic evolves.

Rules heavily rely on properties and will use them to compute their effects.
Note: if a property is missing on an entity, a default value will be used instead (often to the detriment of the entity missing it.)

[Properties](../../properties/README.md)

## GameState

Structs and helper functions around the game state.

## Use Skill 

Allow one to activate a skill.  It's job is mostly to handle targeting validation, and cost of the skill.
Effect application is handled by the [EffectApplicator](../../property/effect/effectapplicator/README.md)


## Movement

Movement is the action entities uses to move from tile to tile.

### Local Ctx

this local context is used to store the current state of the movement.
It's mostly useless right now but expect further use in the future, especially when traps and triggered actions will come into play.

```go
type localMoveCtx struct {
	*GameState
	log *logrus.Entry
}
```

### v0.0.3

Added checks for HasMoved and HasActed. Can't move after having acted.
=> HasMoved is only set after having acted, or having been stunned (begining the turn with more than Half MaxPV in Stun applied)
Expect HasMoved to be set in reaction to other effects in the future. (like traps, reactions, and such.)

### v0.0.2

Move now makes use of the new properties system. It will use the following properties:

* `Movement`: define the reach of the entity.
* `JumpHeight`: define the maximum height the entity can jump.

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

### v0.0.2

Move now makes use of the new properties system. It will use the following properties:

* `Attack`: damage that will be inflicted to the target.
* `AttackRange`: define the reach of the entity.
* `Defense`: damage reduction applied to the entity.
* `HP`: Maximum number of hit points the entity can take.

### v0.0.1

Attack are death sentences: the attacked entity will be removed from the game.
Attack reach is 1 tile away from the attacker.

## Test 	

`rules_test.go` contains structures and functions to assist testings.

See dedicated test files for rules testing.