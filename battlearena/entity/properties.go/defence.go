package properties

// Defence is the base damage reduction property of a character (for now)
type Defence struct {
	Defence int
}

func (a *Defence) I() int {
	return a.Defence
}

func (a *Defence) SetI(f int) {
	a.Defence = f
}

func (a *Defence) Name(i InformationLevel) string {
	if i >= FriendlyController {
		return "Defence"
	}
	return ""
}

func (a *Defence) UserFriendlyGet(i InformationLevel) interface{} {
	if i >= FriendlyController {
		return a.Defence
	}
	return nil
}

func (a *Defence) Get() interface{} {
	return a.Defence
}

func (a *Defence) Set(p interface{}) {
	a.Defence = p.(int)
}

func (a *Defence) Increase() {
	// this will be used later on when we implement leveling.
}

func (a *Defence) GetType() PropertyType {
	return Character
}
