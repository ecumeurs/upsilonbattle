package rules

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/def"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type FakeController struct {
	ID             uuid.UUID
	NotifyMessages []*message.Message
	SentMessages   []*message.Message // expect only no reply.
}

func newFakeController() *FakeController {
	return &FakeController{
		ID:             uuid.New(),
		NotifyMessages: make([]*message.Message, 0),
		SentMessages:   make([]*message.Message, 0),
	}
}

func (fc *FakeController) NotifyActor(m *message.Message) {
	fc.NotifyMessages = append(fc.NotifyMessages, m)
}

func (fc *FakeController) SendActor(m *message.Message, replyTo chan *message.Message) {
	fc.SentMessages = append(fc.SentMessages, m)
	go func() {
		replyTo <- m.Reply()
	}()
}

type FakeState struct {
	Controller1     uuid.UUID
	Controller2     uuid.UUID
	Entity1         uuid.UUID
	Entity2         uuid.UUID
	Entity3         uuid.UUID
	Entity4         uuid.UUID
	FakeController1 *FakeController
	FakeController2 *FakeController
}

func makeGameStateForTwo() (*GameState, FakeState) {

	gs := New(uuid.New())
	ctrl := newFakeController()
	gs.Controllers[ctrl.ID] = ctrl
	ctrl2 := newFakeController()
	gs.Controllers[ctrl2.ID] = ctrl2

	gs.Grid = grid.NewGrid(10, 10, 3)
	// Generate 2 entities for each controller.

	ent1 := entity.New()
	ent1.ControllerID = ctrl.ID
	ent1.Position = position.Position{X: 0, Y: 0, Z: 3}
	ent1.CurrentDelay = 200
	ent1.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		ent1.Properties[v.Name(property.GameMaster)] = v
	}
	gs.Grid.MoveEntity(ent1.Position, ent1.Position, ent1.ID)
	gs.Entities[ent1.ID] = ent1
	gs.Turner.AddEntity(ent1.ID, ent1.CurrentDelay)

	ent2 := entity.New()
	ent2.ControllerID = ctrl.ID
	ent2.Position = position.Position{X: 1, Y: 0, Z: 3}
	ent2.CurrentDelay = 250
	ent2.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		ent2.Properties[v.Name(property.GameMaster)] = v
	}
	gs.Grid.MoveEntity(ent2.Position, ent2.Position, ent2.ID)
	gs.Entities[ent2.ID] = ent2
	gs.Turner.AddEntity(ent2.ID, ent2.CurrentDelay)

	ent3 := entity.New()
	ent3.ControllerID = ctrl2.ID
	ent3.Position = position.Position{X: 9, Y: 9, Z: 3}
	ent3.CurrentDelay = 300
	ent3.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		ent3.Properties[v.Name(property.GameMaster)] = v
	}
	gs.Grid.MoveEntity(ent3.Position, ent3.Position, ent3.ID)
	gs.Entities[ent3.ID] = ent3
	gs.Turner.AddEntity(ent3.ID, ent3.CurrentDelay)

	ent4 := entity.New()
	ent4.ControllerID = ctrl2.ID
	ent4.Position = position.Position{X: 8, Y: 9, Z: 3}
	ent4.CurrentDelay = 350
	ent4.Type = entity.Character
	for _, v := range def.PropertiesForCharacter() {
		ent4.Properties[v.Name(property.GameMaster)] = v
	}
	gs.Grid.MoveEntity(ent4.Position, ent4.Position, ent4.ID)
	gs.Entities[ent4.ID] = ent4
	gs.Turner.AddEntity(ent4.ID, ent4.CurrentDelay)

	fake := FakeState{
		Controller1:     ctrl.ID,
		Controller2:     ctrl2.ID,
		Entity1:         ent1.ID,
		Entity2:         ent2.ID,
		Entity3:         ent3.ID,
		Entity4:         ent4.ID,
		FakeController1: ctrl,
		FakeController2: ctrl2,
	}

	return gs, fake
}
