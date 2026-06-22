package battletest

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/gamestate"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rules"
	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/entity/skill"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilontypes/property/def"
	"github.com/ecumeurs/upsilontypes/property/effect"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

// Scenario is a fluent builder over a GameState. It wires controllers, entities,
// skills, and positional effects, then drives the rules layer (UseSkill / Move /
// turns) and exposes inspectors for assertions.
type Scenario struct {
	GS *gamestate.GameState

	t           testing.TB
	controllers map[int]*FakeController // teamID -> controller
	nextDelay   int
}

// Actor is a lightweight handle to an entity placed in the scenario.
type Actor struct {
	ID uuid.UUID
	s  *Scenario
}

// New creates a Scenario backed by a grid of the given dimensions.
func New(t testing.TB, x, y, z int) *Scenario {
	t.Helper()
	gs := gamestate.New(uuid.New())
	gs.Grid = grid.NewGrid(x, y, z)
	return &Scenario{
		GS:          gs,
		t:           t,
		controllers: make(map[int]*FakeController),
		nextDelay:   100,
	}
}

// Controller returns the FakeController owning the given team, creating it on first use.
func (s *Scenario) Controller(team int) *FakeController {
	if c, ok := s.controllers[team]; ok {
		return c
	}
	c := newFakeController()
	s.controllers[team] = c
	s.GS.Controllers[c.ID] = c
	return c
}

// PlaceOption mutates an entity during placement (HP, movement, combat stats, …).
type PlaceOption func(*entity.Entity)

// WithHP sets the entity's current and max HP.
func WithHP(v int) PlaceOption {
	return func(e *entity.Entity) {
		e.RepsertPropertyCMaxValue(property.HP, v)
		e.RepsertPropertyCValue(property.HP, v)
	}
}

// WithMovement sets the entity's current and max movement points.
func WithMovement(v int) PlaceOption {
	return func(e *entity.Entity) {
		e.RepsertPropertyCMaxValue(property.Movement, v)
		e.RepsertPropertyCValue(property.Movement, v)
	}
}

// WithMP sets the entity's current and max mana points.
func WithMP(v int) PlaceOption {
	return func(e *entity.Entity) {
		e.RepsertPropertyCMaxValue(property.MP, v)
		e.RepsertPropertyCValue(property.MP, v)
	}
}

// WithSP sets the entity's current and max skill/stamina points.
func WithSP(v int) PlaceOption {
	return func(e *entity.Entity) {
		e.RepsertPropertyCMaxValue(property.SP, v)
		e.RepsertPropertyCValue(property.SP, v)
	}
}

// WithStat sets an integer stat property (Attack, Defense, ArmorRating, JumpHeight, …).
func WithStat(p interface{}, v int) PlaceOption {
	return func(e *entity.Entity) { e.RepsertPropertyValue(p, v) }
}

// Place creates a Character entity for the given team at pos and registers it on the grid and turner.
func (s *Scenario) Place(team int, pos position.Position, opts ...PlaceOption) *Actor {
	ctrl := s.Controller(team)

	ent := entity.New()
	ent.ControllerID = ctrl.ID
	ent.Position = pos
	ent.Type = entity.Character
	ent.CurrentDelay = s.nextDelay
	s.nextDelay += 50

	for _, v := range def.PropertiesForCharacter() {
		ent.Properties[v.Name(property.GameMaster)] = v
	}
	ent.RepsertPropertyValue(property.TeamID, team)

	for _, o := range opts {
		o(&ent)
	}

	if err := s.GS.Grid.MoveEntity(pos, pos, ent.ID); err != nil {
		s.t.Fatalf("battletest: failed to place entity at %s: %v", pos, err)
	}
	s.GS.Entities[ent.ID] = ent
	s.GS.Turner.AddEntity(ent.ID, ent.CurrentDelay)

	return &Actor{ID: ent.ID, s: s}
}

// Give registers a skill on the actor and returns the skill ID for later use.
func (a *Actor) Give(sk skill.Skill) uuid.UUID {
	ent := a.s.GS.Entities[a.ID]
	ent.Skills[sk.ID] = sk
	a.s.GS.Entities[a.ID] = ent
	return sk.ID
}

// Trap registers a positional effect at pos and returns its effect ID.
func (s *Scenario) Trap(pos position.Position, eff effect.Effect) uuid.UUID {
	id := uuid.New()
	s.GS.Effects[id] = eff
	s.GS.PositionalEffects[pos] = append(s.GS.PositionalEffects[pos], id)
	return id
}

// Turn forces the turn to the actor and clears its acted/moved state (fresh turn).
func (s *Scenario) Turn(a *Actor) {
	s.GS.Turner.ForceTurn(a.ID)
	ent := s.GS.Entities[a.ID]
	ent.CurrentDelay = 0
	acted := ent.GetProperty(property.HasActed)
	acted.Set(false)
	ent.UpdateProperty(acted)
	moved := ent.GetProperty(property.HasMoved)
	moved.Set(false)
	ent.UpdateProperty(moved)
	s.GS.Entities[a.ID] = ent
}

// UseSkill drives rules.UseSkill for the actor against a target position.
func (s *Scenario) UseSkill(caster *Actor, skillID uuid.UUID, target position.Position) (*message.Message, []rulermethods.ControllerAttacked, []rulermethods.ControllerSkillUsed) {
	ent := s.GS.Entities[caster.ID]
	msg := message.Create(nil, rulermethods.ControllerUseSkill{
		ControllerID: ent.ControllerID,
		EntityID:     caster.ID,
		SkillID:      skillID,
		Target:       target,
	}, nil)
	return rules.UseSkill(s.GS, msg, msg.TargetMethod.(rulermethods.ControllerUseSkill))
}

// Move drives rules.Move for the actor along the given path.
func (s *Scenario) Move(a *Actor, path ...position.Position) *message.Message {
	ent := s.GS.Entities[a.ID]
	msg := message.Create(nil, rulermethods.ControllerMove{
		ControllerID: ent.ControllerID,
		EntityID:     a.ID,
		Path:         path,
	}, nil)
	return rules.Move(s.GS, msg, msg.TargetMethod.(rulermethods.ControllerMove))
}

// BeginTurn fires the start-of-turn pipeline (including OnTurn positional effects) for the actor.
func (s *Scenario) BeginTurn(a *Actor) {
	rules.BeginingOfTurn(s.GS, s.GS.Entities[a.ID])
}
