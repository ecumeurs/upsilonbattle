# .\tools\actor

[Up](../README.md)

# Actor

An actor is a stand alone thread that can only be accessed by using the corresponding Message Queue. 

The actor will expose it's available actions through it and will only respond through it as well.

This ensure data integrity and thread safety.

# Usage

Declare a new struct embedding the Actor struct and set the ReceiveMessageHandler appropriately

```go
type MyActor struct {
    actor.Actor
}

type MyMethod struct {}

func (a *MyActor) Init() {
    a.Actor.Start()
    a.ReceiveMessageHandler =  a.ReceiveMessage
}

func (a *Actor) ReceiveMessage(msg message.Message) {
    select msg.TargetMethod.(type) {
        case MyMethod:
            a.MyMethod(msg)
    }
}

func (a *Actor) MyMethod(msg message.Message) {
    // Do something
}
```

Then you can create a new instance of your actor and send it a message.
