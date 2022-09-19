# .\battlearena

[Up](../README.md)

# Battle Arena

This is the main entry point to the battle arena. It contains the main game access point, and pointers to all critical game objects.

* The Controllers: they are the main inputs into the game, they are the relay from which the commands are sent to the entities. They come into two kinds: the player controller, and the AI controller.
* The Ruler: it's the main game logic, it's the one that decides who wins and who loses, and it's the one that validate any input from the controllers.

