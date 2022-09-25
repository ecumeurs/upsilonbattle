package properties

// HP is the hp property of a character (for now)
type HP struct {
	HP int
}

func (a *HP) I() int {
	return a.HP
}

func (a *HP) SetI(f int) {
	a.HP = f
}

func (a *HP) Name(i InformationLevel) string {
	return "HP"
}

func (a *HP) UserFriendlyGet(i InformationLevel) interface{} {
	return a.HP
}

func (a *HP) Get() interface{} {
	return a.HP
}

func (a *HP) Set(p interface{}) {
	a.HP = p.(int)
}

func (a *HP) Increase() {
	// this will be used later on when we implement leveling.
}

func (a *HP) GetType() PropertyType {
	return Character // it will be evolved to base damage later on and HP will be moved to Skill and Damage will be created for item.
}
