# .\battlearena\ruler

[Up](../README.md)

# Ruler

This module is the main game logic. It's the one that decides who wins and who loses, and it's the one that validate any input from the controllers.
It also handle the turn logic, and the game state.

## Arena State 

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

[See GameState](./rules/README.md) for internal state.

## Turn Logic

Every entity joining the fight will be assigned a delay credit. This delay credit is the amount of time that the entity will have to wait before being able to act.

The first entity to play it the one with the lowest credit. Upon begining its turn, all remaining entities will have their credit reduced by as much.

The entity with the lowest credit will be the first to play again, and so on.

Per turn, each entity will be able to perform a move and an action.
Both of these will have an impact on the delay credit of the entity, determining their next action opportunity. Controllers will have to take this into account when sending commands to the entities.

[See Turner](turner/README.md) for more information.

## Game Rule Sets

Note: Most rules are set in the rules module! [Check it out](rules/README.md)

(TODO)
These modules will handle the application of the game rules. They will be responsible to alter entities and map states.

These will also influence delay credits.
ATM: 
* Per movement: +200
* Per attack: +500
* Per timeout: +1000 // not implemented
* End of turn: +500

Expect attacks to be variable in the future. Expect movement to be affected by the terrain. Expect delay of other entities to be affected by entities techniques.


## TODO

* Start Battle is too abrupt: need to add a preceding state informing all controller that the game is ready to start, they need to ack this message before the game indeed starts ( allow them to prepare their UI for example(fetch grid and entities) )
  * Note: grid and entities are already available before battle start, but entities won't be marked with their controller id (except for current one)
    * Note2: controllers shouldn't be able to request foreign entities data ... (privacy is a thing)
* Same with End of Turn and Turn Begin. 
* Many other thing ...

## Tests


Most usefull test right now is in ruler_fullgame_test.go. The fake controller is a simple seek and destroy controller that will try to kill the closest entity. It will move towards it, and attack it if it's in range. The game runs for two controllers (both are the same type of controller, so it's a matter turn order to decide who wins)   


```txt
time="2022-09-24T09:15:45+02:00" level=info msg="Actor started" component=actor name=Ruler
time="2022-09-24T09:15:45+02:00" level=info msg="Starting message queue" component=messagequeue name=Ruler
time="2022-09-24T09:15:45+02:00" level=info msg="Received message" RequestID=bb552ee7 component=Ruler message_type=actor.ActorStarted
time="2022-09-24T09:15:45+02:00" level=warning msg="Unhandled message" component=actor message="[R bb552ee7]" message_type=actor.ActorStarted name=Ruler
time="2022-09-24T09:15:45+02:00" level=info msg="Actor started" component=actor name=Fake1
time="2022-09-24T09:15:45+02:00" level=info msg="Starting message queue" component=messagequeue name=Fake1
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received message: actor.ActorStarted" Controller=Fake1 controllerID=fe47d14e message_type=actor.ActorStarted
time="2022-09-24T09:15:45+02:00" level=info msg="Actor started" component=actor name=Fake2
time="2022-09-24T09:15:45+02:00" level=warning msg="Unhandled message" component=actor message="[R e66d54c2]" message_type=actor.ActorStarted name=Fake1
time="2022-09-24T09:15:45+02:00" level=info msg="Starting message queue" component=messagequeue name=Fake2
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received message: actor.ActorStarted" Controller=Fake2 controllerID=a77e24dd message_type=actor.ActorStarted
time="2022-09-24T09:15:45+02:00" level=warning msg="Unhandled message" component=actor message="[R 16d7a1ab]" message_type=actor.ActorStarted name=Fake2
time="2022-09-24T09:15:45+02:00" level=info msg="Received message" RequestID=2ddd49b0 component=Ruler message_type=rulermethods.AddController
time="2022-09-24T09:15:45+02:00" level=info msg=AddController ControllerID=fe47d14e RequestID=2ddd49b0 component=Ruler name=Ruler
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received message: controllermethods.SetValidatorAndQueue" Controller=Fake1 controllerID=fe47d14e message_type=controllermethods.SetValidatorAndQueue
time="2022-09-24T09:15:45+02:00" level=info msg="Received message" RequestID=072387d8 component=Ruler message_type=rulermethods.AddController
time="2022-09-24T09:15:45+02:00" level=info msg=AddController ControllerID=a77e24dd RequestID=072387d8 component=Ruler name=Ruler
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received message: controllermethods.SetValidatorAndQueue" Controller=Fake2 controllerID=a77e24dd message_type=controllermethods.SetValidatorAndQueue
time="2022-09-24T09:15:45+02:00" level=info msg="Received message" RequestID=a85e1e4d component=Ruler message_type=rulermethods.GetGridState
time="2022-09-24T09:15:45+02:00" level=info msg="Received message" RequestID=2da7c4b0 component=Ruler message_type=rulermethods.GetEntitiesState
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received reply: rulermethods.GetGridStateReply" Controller=Fake1
time="2022-09-24T09:15:45+02:00" level=info msg="Received message" RequestID=d1a020f5 component=Ruler message_type=rulermethods.GetGridState
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received reply: rulermethods.GetEntitiesStateReply" Controller=Fake1
time="2022-09-24T09:15:45+02:00" level=info msg="Received message" RequestID=def3b917 component=Ruler message_type=rulermethods.GetEntitiesState
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received reply: rulermethods.GetGridStateReply" Controller=Fake2
time="2022-09-24T09:15:45+02:00" level=info msg="Controller received reply: rulermethods.GetEntitiesStateReply" Controller=Fake2
time="2022-09-24T09:15:45+02:00" level=info msg="New Turn Received" Controller=Fake1 RequestId=2da7c4b0 Turn="00000000 0,69c770f5 1102,07263540 1638,e55d435a 1688,406a4354 1873,"
time="2022-09-24T09:15:45+02:00" level=info msg="New Turn Received" Controller=Fake2 RequestId=def3b917 Turn="00000000 0,69c770f5 1102,07263540 1638,e55d435a 1688,406a4354 1873,"
time="2022-09-24T09:15:47+02:00" level=info msg="Received message" RequestID=1f37b7ed component=Ruler message_type=rulermethods.BattleStart
time="2022-09-24T09:15:47+02:00" level=info msg="Game started" RequestID=1f37b7ed component=Ruler name=Ruler
time="2022-09-24T09:15:47+02:00" level=info msg="First entity to play" RequestID=1f37b7ed component=Ruler entityID=69c770f5 name=Ruler
time="2022-09-24T09:15:47+02:00" level=info msg="Controller received message: rulermethods.BattleStart" Controller=Fake2 controllerID=a77e24dd message_type=rulermethods.BattleStart
time="2022-09-24T09:15:47+02:00" level=info msg="##### BattleStart #####"
time="2022-09-24T09:15:47+02:00" level=info msg="Controller received message: rulermethods.BattleStart" Controller=Fake1 controllerID=fe47d14e message_type=rulermethods.BattleStart
time="2022-09-24T09:15:47+02:00" level=info msg="Received message" RequestID=154acc41 component=Ruler message_type=rulermethods.GetGridState
time="2022-09-24T09:15:47+02:00" level=info msg="##### BattleStart #####"
time="2022-09-24T09:15:47+02:00" level=info msg="Controller received reply: rulermethods.GetGridStateReply" Controller=Fake2
time="2022-09-24T09:15:47+02:00" level=info msg="Received message" RequestID=5eaae410 component=Ruler message_type=rulermethods.GetEntitiesState
time="2022-09-24T09:15:47+02:00" level=info msg="Received message" RequestID=5934d800 component=Ruler message_type=rulermethods.GetGridState
time="2022-09-24T09:15:47+02:00" level=info msg="Controller received reply: rulermethods.GetEntitiesStateReply" Controller=Fake2
time="2022-09-24T09:15:47+02:00" level=info msg="New Turn Received" Controller=Fake2 RequestId=5eaae410 Turn="69c770f5 0,07263540 536,e55d435a 586,406a4354 771,"
time="2022-09-24T09:15:47+02:00" level=info msg="Controller received reply: rulermethods.GetGridStateReply" Controller=Fake1
time="2022-09-24T09:15:47+02:00" level=info msg="Received message" RequestID=c4b79a49 component=Ruler message_type=rulermethods.GetEntitiesState
time="2022-09-24T09:15:47+02:00" level=info msg="Controller received reply: rulermethods.GetEntitiesStateReply" Controller=Fake1
time="2022-09-24T09:15:47+02:00" level=info msg="New Turn Received" Controller=Fake1 RequestId=c4b79a49 Turn="69c770f5 0,07263540 536,e55d435a 586,406a4354 771,"
time="2022-09-24T09:15:49+02:00" level=info msg="Controller received message: rulermethods.ControllerNextTurn" Controller=Fake1 controllerID=fe47d14e message_type=rulermethods.ControllerNextTurn
time="2022-09-24T09:15:49+02:00" level=info msg="##### Turn BEGIN #####" Controller=fe47d14e ControllerName=Fake1 EntityID=69c770f5 RequestID=c0075e32 Turn="69c770f5 0,07263540 536,e55d435a 586,406a4354 771,"
time="2022-09-24T09:15:49+02:00" level=info msg=selectNearestFoe controller=fe47d14e entity=69c770f5 pos="(19, 39, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg=candidate candidate_controller=a77e24dd candidate_entity=e55d435a candidate_pos="(23, 36, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg=selected selected_controller=a77e24dd selected_entity=e55d435a selected_pos="(23, 36, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg=candidate candidate_controller=a77e24dd candidate_entity=07263540 candidate_pos="(26, 14, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg=nearest nearest_controller=a77e24dd nearest_entity=e55d435a nearest_pos="(23, 36, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg="Moving attacker" EntityID=69c770f5 Expected="(22, 36, 11)" Position="(19, 39, 11)" TargetEntity=e55d435a TargetPosition="(23, 36, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg="Received message" RequestID=3d491c5c component=Ruler message_type=rulermethods.ControllerMove
time="2022-09-24T09:15:49+02:00" level=info msg="Controller received reply: rulermethods.ControllerMoveReply" Controller=Fake1
time="2022-09-24T09:15:49+02:00" level=info msg=Attacking EntityID=69c770f5 Position="(22, 36, 11)" TargetEntity=e55d435a TargetPosition="(23, 36, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg="Received message" RequestID=9288c14e component=Ruler message_type=rulermethods.ControllerAttack
time="2022-09-24T09:15:49+02:00" level=info msg="##### Entity removed #####" RequestID=9288c14e component=Ruler entityID=e55d435a name=Ruler position="(23, 36, 11)"
time="2022-09-24T09:15:49+02:00" level=info msg="Controller received reply: rulermethods.ControllerAttackReply" Controller=Fake1
time="2022-09-24T09:15:49+02:00" level=info msg="Attack done, ending turn" Controller=fe47d14e ControllerName=Fake1 Error=false Message=
time="2022-09-24T09:15:49+02:00" level=info msg="Controller received message: rulermethods.ControllerAttacked" Controller=Fake2 controllerID=a77e24dd message_type=rulermethods.ControllerAttacked
time="2022-09-24T09:15:49+02:00" level=info msg="Received message" RequestID=b896d408 component=Ruler message_type=rulermethods.EndOfTurn
time="2022-09-24T09:15:49+02:00" level=info msg="##### END OF TURN #####" RequestID=b896d408 component=Ruler controllerID=fe47d14e entityID=69c770f5 name=Ruler
time="2022-09-24T09:15:49+02:00" level=info msg="Controller received reply: rulermethods.EndOfTurn" Controller=Fake1
time="2022-09-24T09:15:49+02:00" level=info msg="Controller received message: rulermethods.EntitiesStateChanged" Controller=Fake1 controllerID=fe47d14e message_type=rulermethods.EntitiesStateChanged
time="2022-09-24T09:15:49+02:00" level=info msg="New Turn Received" Controller=Fake1 RequestId=2e077ec4 Turn="69c770f5 0,07263540 536,406a4354 235,69c770f5 1664,"
time="2022-09-24T09:15:49+02:00" level=info msg="Controller received message: rulermethods.EntitiesStateChanged" Controller=Fake2 controllerID=a77e24dd message_type=rulermethods.EntitiesStateChanged
time="2022-09-24T09:15:49+02:00" level=info msg="New Turn Received" Controller=Fake2 RequestId=3d36c8da Turn="69c770f5 0,07263540 536,406a4354 235,69c770f5 1664,"
time="2022-09-24T09:15:51+02:00" level=info msg="Controller received message: rulermethods.ControllerNextTurn" Controller=Fake2 controllerID=a77e24dd message_type=rulermethods.ControllerNextTurn
time="2022-09-24T09:15:51+02:00" level=info msg="##### Turn BEGIN #####" Controller=a77e24dd ControllerName=Fake2 EntityID=07263540 RequestID=cd49805b Turn="07263540 0,406a4354 235,69c770f5 1664,"
time="2022-09-24T09:15:51+02:00" level=info msg=selectNearestFoe controller=a77e24dd entity=07263540 pos="(26, 14, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg=candidate candidate_controller=fe47d14e candidate_entity=69c770f5 candidate_pos="(22, 36, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg=selected selected_controller=fe47d14e selected_entity=69c770f5 selected_pos="(22, 36, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg=candidate candidate_controller=fe47d14e candidate_entity=406a4354 candidate_pos="(30, 13, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg=selected selected_controller=fe47d14e selected_entity=406a4354 selected_pos="(30, 13, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg=nearest nearest_controller=fe47d14e nearest_entity=406a4354 nearest_pos="(30, 13, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg="Moving attacker" EntityID=07263540 Expected="(29, 13, 11)" Position="(26, 14, 11)" TargetEntity=406a4354 TargetPosition="(30, 13, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg="Received message" RequestID=b80cb7a1 component=Ruler message_type=rulermethods.ControllerMove
time="2022-09-24T09:15:51+02:00" level=info msg="Controller received reply: rulermethods.ControllerMoveReply" Controller=Fake2
time="2022-09-24T09:15:51+02:00" level=info msg=Attacking EntityID=07263540 Position="(29, 13, 11)" TargetEntity=406a4354 TargetPosition="(30, 13, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg="Received message" RequestID=1f2c3ebd component=Ruler message_type=rulermethods.ControllerAttack
time="2022-09-24T09:15:51+02:00" level=info msg="##### Entity removed #####" RequestID=1f2c3ebd component=Ruler entityID=406a4354 name=Ruler position="(30, 13, 11)"
time="2022-09-24T09:15:51+02:00" level=info msg="Controller received reply: rulermethods.ControllerAttackReply" Controller=Fake2
time="2022-09-24T09:15:51+02:00" level=info msg="Attack done, ending turn" Controller=a77e24dd ControllerName=Fake2 Error=false Message=
time="2022-09-24T09:15:51+02:00" level=info msg="Received message" RequestID=b2769da9 component=Ruler message_type=rulermethods.EndOfTurn
time="2022-09-24T09:15:51+02:00" level=info msg="Controller received message: rulermethods.ControllerAttacked" Controller=Fake1 controllerID=fe47d14e message_type=rulermethods.ControllerAttacked
time="2022-09-24T09:15:51+02:00" level=info msg="##### END OF TURN #####" RequestID=b2769da9 component=Ruler controllerID=a77e24dd entityID=07263540 name=Ruler
time="2022-09-24T09:15:51+02:00" level=info msg="Controller received reply: rulermethods.EndOfTurn" Controller=Fake2
time="2022-09-24T09:15:51+02:00" level=info msg="Controller received message: rulermethods.EntitiesStateChanged" Controller=Fake2 controllerID=a77e24dd message_type=rulermethods.EntitiesStateChanged
time="2022-09-24T09:15:51+02:00" level=info msg="New Turn Received" Controller=Fake2 RequestId=383137d0 Turn="07263540 0,69c770f5 1664,07263540 136,"
time="2022-09-24T09:15:51+02:00" level=info msg="Controller received message: rulermethods.EntitiesStateChanged" Controller=Fake1 controllerID=fe47d14e message_type=rulermethods.EntitiesStateChanged
time="2022-09-24T09:15:51+02:00" level=info msg="New Turn Received" Controller=Fake1 RequestId=731eb043 Turn="07263540 0,69c770f5 1664,07263540 136,"
time="2022-09-24T09:15:53+02:00" level=info msg="Controller received message: rulermethods.ControllerNextTurn" Controller=Fake1 controllerID=fe47d14e message_type=rulermethods.ControllerNextTurn
time="2022-09-24T09:15:53+02:00" level=info msg="##### Turn BEGIN #####" Controller=fe47d14e ControllerName=Fake1 EntityID=69c770f5 RequestID=113bc09b Turn="69c770f5 0,07263540 136,"
time="2022-09-24T09:15:53+02:00" level=info msg=selectNearestFoe controller=fe47d14e entity=69c770f5 pos="(22, 36, 11)"
time="2022-09-24T09:15:53+02:00" level=info msg=candidate candidate_controller=a77e24dd candidate_entity=07263540 candidate_pos="(29, 13, 11)"
time="2022-09-24T09:15:53+02:00" level=info msg=selected selected_controller=a77e24dd selected_entity=07263540 selected_pos="(29, 13, 11)"
time="2022-09-24T09:15:53+02:00" level=info msg=nearest nearest_controller=a77e24dd nearest_entity=07263540 nearest_pos="(29, 13, 11)"
time="2022-09-24T09:15:53+02:00" level=info msg="Moving attacker" EntityID=69c770f5 Expected="(29, 14, 11)" Position="(22, 36, 11)" TargetEntity=07263540 TargetPosition="(29, 13, 11)"
time="2022-09-24T09:15:53+02:00" level=info msg="Received message" RequestID=43347f3a component=Ruler message_type=rulermethods.ControllerMove
time="2022-09-24T09:15:53+02:00" level=info msg="Controller received reply: rulermethods.ControllerMoveReply" Controller=Fake1
time="2022-09-24T09:15:53+02:00" level=info msg=Attacking EntityID=69c770f5 Position="(29, 14, 11)" TargetEntity=07263540 TargetPosition="(29, 13, 11)"
time="2022-09-24T09:15:53+02:00" level=info msg="Received message" RequestID=2680ee5c component=Ruler message_type=rulermethods.ControllerAttack
time="2022-09-24T09:15:53+02:00" level=info msg="##### Entity removed #####" RequestID=2680ee5c component=Ruler entityID=07263540 name=Ruler position="(29, 13, 11)"
time="2022-09-24T09:15:53+02:00" level=info msg="Controller received reply: rulermethods.ControllerAttackReply" Controller=Fake1
time="2022-09-24T09:15:53+02:00" level=info msg="Attack done, ending turn" Controller=fe47d14e ControllerName=Fake1 Error=false Message=
time="2022-09-24T09:15:53+02:00" level=info msg="Received message" RequestID=30d7d39f component=Ruler message_type=rulermethods.EndOfTurn
time="2022-09-24T09:15:53+02:00" level=info msg="##### END OF BATTLE! #####" RequestID=30d7d39f component=Ruler controllerID=fe47d14e entityID=69c770f5 name=Ruler
time="2022-09-24T09:15:53+02:00" level=info msg="Controller received message: rulermethods.ControllerAttacked" Controller=Fake2 controllerID=a77e24dd message_type=rulermethods.ControllerAttacked
time="2022-09-24T09:15:53+02:00" level=info msg="Controller received reply: rulermethods.EndOfTurn" Controller=Fake1
time="2022-09-24T09:15:53+02:00" level=info msg="Controller received message: rulermethods.BattleEnd" Controller=Fake2 controllerID=a77e24dd message_type=rulermethods.BattleEnd
time="2022-09-24T09:15:53+02:00" level=info msg="##### BattleEnd #####"
time="2022-09-24T09:15:53+02:00" level=info msg="Controller received message: rulermethods.BattleEnd" Controller=Fake1 controllerID=fe47d14e message_type=rulermethods.BattleEnd
time="2022-09-24T09:15:53+02:00" level=info msg="##### BattleEnd #####"
time="2022-09-24T09:15:53+02:00" level=info msg="Battle Finished, doing end of game Checks"
time="2022-09-24T09:15:53+02:00" level=info msg="Received message" RequestID=df8b137d component=Ruler message_type=rulermethods.GetGridState
time="2022-09-24T09:15:53+02:00" level=info msg="Received message" RequestID=52b572e5 component=Ruler message_type=rulermethods.GetEntitiesState
time="2022-09-24T09:15:53+02:00" level=info msg="END OF LAST CHECKS"
time="2022-09-24T09:15:53+02:00" level=info msg="Received message" RequestID=9e97873f component=Ruler message_type=rulermethods.ControllerQuit
time="2022-09-24T09:15:53+02:00" level=info msg="ctrl1 stopped"
time="2022-09-24T09:15:53+02:00" level=info msg="Received message" RequestID=44b1ea4b component=Ruler message_type=rulermethods.ControllerQuit
time="2022-09-24T09:15:53+02:00" level=info msg="Stopping message queue" component=messagequeue name=Fake1
time="2022-09-24T09:15:53+02:00" level=info msg="ctrl2 stopped"
time="2022-09-24T09:15:53+02:00" level=info msg="Stopping message queue" component=messagequeue name=Fake2
```


