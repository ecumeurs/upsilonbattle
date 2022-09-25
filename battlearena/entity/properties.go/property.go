package properties

type PropertyType int

const (
	Character PropertyType = 0
	Skill     PropertyType = 1
	Item      PropertyType = 2
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
	I() int
	SetI(int)
}

type FloatProperty interface {
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
	}
}
