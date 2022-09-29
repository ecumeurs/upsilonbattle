package entitygenerator

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/properties"
)

func TestEntityGenerator(t *testing.T) {
	ent := GenerateRandomEntity()
	t.Log(ent.PrettyString())
	ent = GenerateRandomEntity()
	t.Log(ent.PrettyString())
	ent = GenerateRandomEntity()
	t.Log(ent.PrettyString())
	ent = GenerateRandomEntity()
	t.Log(ent.PrettyString())
	ent = GenerateRandomEntity()
	t.Log(ent.PrettyString())
}
func TestEntityGeneratorGetProperty(t *testing.T) {
	ent := GenerateRandomEntity()
	t.Log(ent.PrettyString())

	var defaultAttackRangeProp = properties.DefaultIntProperty(1)
	prop := ent.GetPropertyI("AttackRange", &defaultAttackRangeProp)

	if prop.Name(properties.GameMaster) != "AttackRange" {
		t.Error("Wrong property name")
	}

	if prop.Get().(int) == 0 {
		t.Error("Wrong property value")
	}

	t.Log(properties.PrettyPrint(prop, properties.GameMaster))

}
