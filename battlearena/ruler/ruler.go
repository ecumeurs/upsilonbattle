package ruler

import (
	"fmt"
	"time"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilontypes/entity/entitygenerator"
	"github.com/ecumeurs/upsilontypes/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rules"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilonmapmaker/gridgenerator"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/ecumeurs/upsilontools/tools/actor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ArenaState int

const (
	WaitingForControllers ArenaState = 1
	InProgress            ArenaState = 2
	Finished              ArenaState = 3
)

// String returns the string representation of the ArenaState.
func (g ArenaState) String() string {
	switch g {
	case WaitingForControllers:
		return "WaitingForControllers"
	case InProgress:
		return "InProgress"
	case Finished:
		return "Finished"
	default:
		return fmt.Sprintf("Unknown GameState %d", g)
	}
}

// Ruler is the main actor managing a battle arena instance.
// It orchestrates turn management, rule enforcement, and controller communication.
type Ruler struct {
	ID uuid.UUID
	*actor.Actor
	GameState    *rules.GameState
	CurrentState ArenaState
	logger       *logrus.Entry

	NbControllers           int
	NbEntitiesPerController int

	ControllerBattleReady map[uuid.UUID]bool
	SetQueueAcks          map[uuid.UUID]bool

	shotClock         *time.Timer
	shotClockVersion  int64
	ShotClockDuration time.Duration
	firstTurnSent     bool
}

// NewCompleteRuler is a factory function that constructs a fully initialized Ruler actor.
// This constructor is designed for rapid deployment of a standard battle arena, pre-populating
// the game state with a default 4x4x5 flat grid and 2 controllers, each assigned one randomly
// generated character entity.
//
// The initialization process involves several distinct steps:
// 1. Unique ID Generation: A new UUID is assigned to the Ruler for global identification.
// 2. Actor Setup: The underlying actor system is initialized with the name "Ruler".
// 3. Game State Creation: A fresh rules.GameState is instantiated and linked to the Ruler.
// 4. Grid Generation: A flat terrain is generated using the gridgenerator package.
// 5. Entity Population: For each expected controller, a random character is generated,
//    assigned to a team (Team 1 or Team 2), and placed at a random position on the grid.
// 6. Turner Initialization: Each entity is added to the initiative turner with a random initial delay.
// 7. Handler Registration: The init() method is called to bind all network and internal message handlers.
//
// WARNING: This function DOES NOT automatically start the Ruler's internal message loop.
// The caller is responsible for invoking the Start() method on the returned Ruler pointer.
// Failure to do so will leave the actor in a dormant state, unable to process registrations or turns.
//
// @spec-link [[api_go_battle_start]]
func NewCompleteRuler() *Ruler {
	id := uuid.New()
	r := Ruler{
		ID:                    id,
		Actor:                 actor.New("Ruler"),
		CurrentState:          WaitingForControllers,
		ControllerBattleReady: make(map[uuid.UUID]bool),
		SetQueueAcks:          make(map[uuid.UUID]bool),
		ShotClockDuration:     30 * time.Second,
	}
	r.GameState = rules.New(r.ID)
	r.logger = logrus.WithFields(logrus.Fields{
		"component": "Ruler",
		"name":      r.Name()})

	gg := gridgenerator.GridGenerator{
		Width:                tools.NewIntRange(4, 4),
		Length:               tools.NewIntRange(4, 4),
		Height:               tools.NewIntRange(5, 5),
		GenerateObstrcution:  false,
		Type:                 gridgenerator.Flat,
		ObstructionRate:      tools.NewIntRange(0, 0),
	}

	r.GameState.Grid = gg.Generate()
	r.NbControllers = 2
	r.NbEntitiesPerController = 1
	nbEntities := r.NbEntitiesPerController * r.NbControllers

	for i := 0; i < nbEntities; i++ {
		e := entitygenerator.GenerateRandomEntity()
		e.Type = entity.Character
		e.Name = fmt.Sprintf("Entity %d", i)
		teamID := 1
		if i >= r.NbEntitiesPerController {
			teamID = 2
		}
		e.RepsertPropertyValue(property.TeamID, teamID)
		e.CurrentDelay = tools.NewIntRange(1000, 2000).Random()
		e.Position = r.GameState.Grid.RandomPosition()
		r.GameState.Grid.MoveEntity(position.New(0, 0, 0), e.Position, e.ID)
		r.GameState.Entities[e.ID] = e
		r.GameState.Turner.AddEntity(e.ID, e.CurrentDelay)
	}

	r.init()
	return &r
}

// NewRuler provides a minimalist constructor for a Ruler actor, allowing for custom arena configuration.
// Unlike NewCompleteRuler, this function returns a Ruler with an empty grid and no pre-populated entities.
// It is the preferred method for production environments where the map and characters are retrieved
// from a database or specified via a configuration file.
//
// Key features of this constructor:
// - Deterministic Seeding: Calls tools.Seed() to ensure consistent random number generation.
// - Clean Slate: Initializes only the essential actor and state containers.
// - Flexibility: Allows the caller to manually invoke SetGrid() and AddEntity() to build the match state.
//
// Lifecycle Requirements:
// 1. Grid Association: A valid grid MUST be provided via SetGrid() before any controllers are added.
// 2. Entity Loading: Participants should be added via AddEntity() to ensure proper registration in the turner.
// 3. Actor Activation: Like all actor-based components, Start() must be called to begin message processing.
//
// @spec-link [[api_go_battle_start]]
func NewRuler(id uuid.UUID) *Ruler {
	tools.Seed()
	r := Ruler{
		ID:                    id,
		Actor:                 actor.New("Ruler"),
		CurrentState:          WaitingForControllers,
		ControllerBattleReady: make(map[uuid.UUID]bool),
		SetQueueAcks:          make(map[uuid.UUID]bool),
		ShotClockDuration:     30 * time.Second,
	}
	r.GameState = rules.New(r.ID)
	r.logger = logrus.WithFields(logrus.Fields{
		"component": "Ruler",
		"name":      r.Name()})

	r.init()
	return &r
}

// SetNbControllers sets the number of expected controllers for the match.
func (r *Ruler) SetNbControllers(nb int) {
	r.NbControllers = nb
}

// SetGrid associates a map grid with the battle arena.
func (r *Ruler) SetGrid(g *grid.Grid) {
	r.GameState.Grid = g
}

// AddEntity adds a new character or object to the battle.
func (r *Ruler) AddEntity(e entity.Entity) {
	e.CurrentDelay = tools.NewIntRange(1000, 1500).Random()

	if r.GameState.Grid == nil {
		r.logger.WithFields(logrus.Fields{
			"entityID": e.ID.String()[0:8]}).Error("Cannot add entity: Grid is not initialized.")
		return
	}

	if e.Position.Equals(position.Position{}) {
		e.Position = r.GameState.Grid.RandomPosition()
	}
	r.GameState.Grid.MoveEntity(position.New(0, 0, 0), e.Position, e.ID)
	r.GameState.Entities[e.ID] = e
	r.GameState.Turner.AddEntity(e.ID, e.CurrentDelay)
}

func (r *Ruler) init() {
	r.AddCallHandler(rulermethods.AddController{}, r.addController, nil)
	r.AddCallHandler(rulermethods.GetState{}, r.getState, nil)
	r.AddCallHandler(rulermethods.GetGridState{}, r.getGridState, nil)
	r.AddCallHandler(rulermethods.GetEntitiesState{}, r.getEntitiesState, nil)
	r.AddCallHandler(rulermethods.GetBoardState{}, r.getBoardState, nil)
	r.AddCallHandler(rulermethods.ControllerMove{}, r.controllerMove, nil)
	r.AddCallHandler(rulermethods.ControllerAttack{}, r.controllerAttack, nil)
	r.AddCallHandler(rulermethods.ControllerUseSkill{}, r.controllerUseSkill, nil)
	r.AddNotificationHandler(rulermethods.NotifyController{}, r.notifyController, nil)
	r.AddCallHandler(rulermethods.EndOfTurn{}, r.endOfTurn, nil)
	r.AddNotificationHandler(rulermethods.ControllerQuit{}, r.controllerQuit, nil)
	r.AddNotificationHandler(rulermethods.BattleStart{}, r.battleStart, nil)
	r.AddNotificationHandler(rulermethods.ControllerBattleReady{}, r.controllerBattleReady, nil)
	r.AddNotificationHandler(rulermethods.InternalTriggerFirstTurn{}, r.processTriggerFirstTurn, nil)
	r.AddNotificationHandler(rulermethods.ControllerTurnReady{}, r.controllerTurnReady, nil)
	r.AddNotificationHandler(rulermethods.ControllerPassed{}, r.controllerPassed, nil)
	r.AddCallHandler(rulermethods.ControllerForfeit{}, r.controllerForfeit, nil)
	r.AddNotificationHandler(rulermethods.Timeout{}, r.timeout, nil)
	r.AddNotificationHandler(actor.ActorAboutToStop{}, r.actorAboutToStop, nil)
	r.AddNotificationHandler(rulermethods.Resurrect{}, r.resurrect, nil)

	r.AddNotificationHandler(rulermethods.TestingDeleteEntity{}, r.testingDeleteEntity, nil)
	r.AddCallHandler(rulermethods.TestingGetState{}, r.testingGetState, nil)
}

// PrintStack prints the current message queue stack for debugging.
func (r *Ruler) PrintStack() {
	r.GetQueue().PrintStack()
}
