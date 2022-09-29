package skillgenerator

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
	"github.com/ecumeurs/upsilontools/tools"
)

func TestGenerateRandomSkill(t *testing.T) {

	skp := def.Attack()
	t.Log(skp)

	skp = defaultproperty.MakeIntProperty(property.Attack, tools.RandomInt(1, 3), property.Public, property.Skill)
	t.Log(skp)

	skp = defaultproperty.MakeIntProperty(property.Damage, tools.RandomInt(1, 3), property.Public, property.Skill)
	t.Log(skp)

	sk := GenerateRandomSkill()
	t.Log(sk)

}
