package property

import (
	"fmt"
)

type PropertyType int

const (
	None      PropertyType = 0
	Character PropertyType = 1
	// @spec-link [[mech_entity_properties_skill_properties]]
	Skill     PropertyType = 2
	Item      PropertyType = 3
)

type InformationLevel int

const (
	Public             InformationLevel = 0
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
	Duplicate() Property
	ApplyBuff(p Property) Property
	UnapplyBuff(p Property) Property
}

// PrettyPrint
func PrettyPrint(p Property, i InformationLevel) string {
	return fmt.Sprintf("%s: %v", p.Name(i), p.UserFriendlyGet(i))
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

type BoolProperty interface {
	Property
	B() bool
	SetB(bool)
}

type IntCounterProperty interface {
	IntProperty
	GetValue() int
	GetMaxValue() int
	SetValue(int)
	SetMaxValue(int)
}
