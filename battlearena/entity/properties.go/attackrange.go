package properties

// AttackRange is the range of attack of a character(will be moved to skill and item later on)
type AttackRange struct {
	AttackRange int
}

func (a *AttackRange) I() int {
	return a.AttackRange
}

func (a *AttackRange) SetI(f int) {
	a.AttackRange = f
}

func (a *AttackRange) Name(i InformationLevel) string {
	if i >= FriendlyController {
		return "AttackRange"
	}
	return ""
}

func (a *AttackRange) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= FriendlyController {
		return a.AttackRange
	}
	return nil
}

func (a *AttackRange) Get() interface{} {
	return a.AttackRange
}

func (a *AttackRange) Set(p interface{}) {
	a.AttackRange = p.(int)
}

func (a *AttackRange) Increase() {
	// this will be used later on when we implement leveling.
}

func (a *AttackRange) GetType() PropertyType {
	return Character
}
