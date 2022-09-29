package properties

// Attack is the base damaging property of a character (for now)
type Attack struct {
	Attack int
}

func (a *Attack) I() int {
	return a.Attack
}

func (a *Attack) SetI(f int) {
	a.Attack = f
}

func (a *Attack) Name(i InformationLevel) string {
	if i >= FriendlyController {
		return "Attack"
	}
	return ""
}

func (a *Attack) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= FriendlyController {
		return a.Attack
	}
	return nil
}

func (a *Attack) Get() interface{} {
	return a.Attack
}

func (a *Attack) Set(p interface{}) {
	a.Attack = p.(int)
}

func (a *Attack) Increase() {
	// this will be used later on when we implement leveling.
}

func (a *Attack) GetType() PropertyType {
	return Character // it will be evolved to base damage later on and Attack will be moved to Skill and Damage will be created for item.
}
