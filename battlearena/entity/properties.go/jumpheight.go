package properties

// JumpHeight is the delta height to the next cell (for now)
type JumpHeight struct {
	JumpHeight int
}

func (a *JumpHeight) I() int {
	return a.JumpHeight
}

func (a *JumpHeight) SetI(f int) {
	a.JumpHeight = f
}

func (a *JumpHeight) Name(i InformationLevel) string {
	if i >= FriendlyController {
		return "JumpHeight"
	}
	return ""
}

func (a *JumpHeight) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= FriendlyController {
		return a.JumpHeight
	}
	return nil
}

func (a *JumpHeight) Get() interface{} {
	return a.JumpHeight
}

func (a *JumpHeight) Set(p interface{}) {
	a.JumpHeight = p.(int)
}

func (a *JumpHeight) Increase() {
	// this will be used later on when we implement leveling.
}

func (a *JumpHeight) GetType() PropertyType {
	return Character
}
