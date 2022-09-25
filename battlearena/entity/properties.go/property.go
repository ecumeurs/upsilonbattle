package properties

type PropertyType int

const (
	None      PropertyType = 0
	Character PropertyType = 1
	Skill     PropertyType = 2
	Item      PropertyType = 3
)

type InformationLevel int

const (
	ExternalObserver   InformationLevel = 0
	ArenaObserver      InformationLevel = 1
	ForeignController  InformationLevel = 2
	FriendlyController InformationLevel = 3

	OwnController InformationLevel = 4

	Analyser          InformationLevel = 5
	ExpertAnalyst     InformationLevel = 6
	SpecialistAnalyst InformationLevel = 7
	MasterAnalyst     InformationLevel = 8

	GameMaster InformationLevel = 9
)

type Property interface {
	Name(i InformationLevel) string                 // most will always reply with a name, some might be hidden by restrictions of with a scrambled name.
	UserFriendlyGet(i InformationLevel) interface{} // most will be expected to return an int (float will be frowned upon) but might return a string if appropriate (status for example) may return nil in which case the information won't be displayed.
	Get() interface{}                               // this will be used mostly internally to compute values from rules.
	Set(p interface{})                              // this will be used mostly internally to compute values from rules.
	Increase()                                      // This won't be used in v0.0.2 but later on when we implement leveling.
	GetType() PropertyType
}

// these interface are here for rules ...
type IntProperty interface {
	Property
	I() int
	SetI(int)
}

type FloatProperty interface {
	Property
	F() float64
	SetF(float64)
}

// note: futher properties may be added per entity basis.
func DefaultPropertiesForCharacter() []Property {
	return []Property{
		&HP{HP: 10},
		&Attack{Attack: 3},
		&AttackRange{AttackRange: 1},
		&Defence{Defence: 0},
		&Movement{Movement: 5},
		&JumpHeight{JumpHeight: 5},
	}
}

type DefaultIntProperty int

// implements IntProperty
func (d *DefaultIntProperty) Name(i InformationLevel) string {
	return ""
}

func (d *DefaultIntProperty) UserFriendlyGet(i InformationLevel) interface{} {
	return *d
}

func (d *DefaultIntProperty) Get() interface{} {
	return *d
}

func (d *DefaultIntProperty) Set(p interface{}) {
	// do nothing!
}

func (d *DefaultIntProperty) Increase() {
	// do nothing
}

func (d *DefaultIntProperty) GetType() PropertyType {
	return None
}

func (d *DefaultIntProperty) I() int {
	return int(*d)
}

func (d *DefaultIntProperty) SetI(int) {
	// do nothing
}

type DefaultFloatProperty int

// implements IntProperty
func (d *DefaultFloatProperty) Name(i InformationLevel) string {
	return ""
}

func (d *DefaultFloatProperty) UserFriendlyGet(i InformationLevel) interface{} {
	return *d
}

func (d *DefaultFloatProperty) Get() interface{} {
	return *d
}

func (d *DefaultFloatProperty) Set(p interface{}) {
	// do nothing!
}

func (d *DefaultFloatProperty) Increase() {
	// do nothing
}

func (d *DefaultFloatProperty) GetType() PropertyType {
	return None
}

func (d *DefaultFloatProperty) I() float64 {
	return float64(*d)
}

func (d *DefaultFloatProperty) SetI(float64) {
	// do nothing
}
