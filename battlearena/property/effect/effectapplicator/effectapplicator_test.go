package effectapplicator

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

type effectapplicatorContext struct {
	Grid      *grid.Grid
	Entities  []entity.Entity
	Pos       []position.Position
	TargetPos position.Position

	Caster entity.Entity
	Target entity.Entity
	Other  entity.Entity
}

func makeTestingEnvironment() effectapplicatorContext {
	res := effectapplicatorContext{
		Grid:     grid.NewGrid(10, 10, 3),
		Entities: []entity.Entity{},
		Pos:      []position.Position{},
	}

	res.Caster = entity.New()
	res.Caster.ControllerID = uuid.New()
	res.Caster.Position = position.Position{X: 0, Y: 0, Z: 3}
	res.Caster.CurrentDelay = 200
	res.Caster.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		res.Caster.Properties[v.Name(property.GameMaster)] = v
	}
	res.Grid.MoveEntity(res.Caster.Position, res.Caster.Position, res.Caster.ID)

	res.Target = entity.New()
	res.Target.ControllerID = uuid.New()
	res.Target.Position = position.Position{X: 1, Y: 0, Z: 3}
	res.TargetPos = res.Target.Position
	res.Target.CurrentDelay = 250
	res.Target.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		res.Target.Properties[v.Name(property.GameMaster)] = v
	}
	res.Grid.MoveEntity(res.Target.Position, res.Target.Position, res.Target.ID)
	res.Entities = append(res.Entities, res.Target)

	res.Other = entity.New()
	res.Other.ControllerID = uuid.New()
	res.Other.Position = position.Position{X: 2, Y: 0, Z: 3}
	res.Other.CurrentDelay = 300
	res.Other.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		res.Other.Properties[v.Name(property.GameMaster)] = v
	}
	res.Grid.MoveEntity(res.Other.Position, res.Other.Position, res.Other.ID)

	return res
}
