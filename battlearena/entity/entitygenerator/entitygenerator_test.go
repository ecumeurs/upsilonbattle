package entitygenerator

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
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

	prop := ent.GetPropertyI("AttackRange", def.AttackRange())

	if prop.Name(property.GameMaster) != "AttackRange" {
		t.Error("Wrong property name")
	}

	if prop.Get().(int) == 0 {
		t.Error("Wrong property value")
	}

	t.Log(property.PrettyPrint(prop, property.GameMaster))

}
