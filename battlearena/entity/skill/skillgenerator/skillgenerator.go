package skillgenerator

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity/skill"
	"github.com/ecumeurs/upsilontools/tools"
)

var propertiesCostRandomizers = map[string]tools.IntRange{
	"JumpHeight": tools.IntRange{Start: 1, End: 3},
}
var propertiesTargetingRandomizers = map[string]tools.IntRange{
	"JumpHeight": tools.IntRange{Start: 1, End: 3},
}
var propertiesEffectRandomizers = map[string]tools.IntRange{
	"JumpHeight": tools.IntRange{Start: 1, End: 3},
}

func GenerateRandomSkill() *skill.Skill {
	return nil
}
