package properties

// Movement is the movement range property of a character (for now)
type Movement struct {
	Movement int
}

func (a *Movement) I() int {
	return a.Movement
}

func (a *Movement) SetI(f int) {
	a.Movement = f
}

func (a *Movement) Name(i InformationLevel) string {
	if i >= FriendlyController {
		return "Movement"
	}
	return ""
}

func (a *Movement) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= FriendlyController {
		return a.Movement
	}
	return nil
}

func (a *Movement) Get() interface{} {
	return a.Movement
}

func (a *Movement) Set(p interface{}) {
	a.Movement = p.(int)
}

func (a *Movement) Increase() {
	// this will be used later on when we implement leveling.
}

func (a *Movement) GetType() PropertyType {
	return Character
}
