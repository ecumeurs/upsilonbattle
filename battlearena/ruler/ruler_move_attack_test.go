package ruler
// @test-link [[uc_combat_turn]]

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilontypes/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/sirupsen/logrus"
)

// TestRulerControllerCanMoveAttackAndEndTurn simulates a basic player turn sequence, including movement, combat, and turn finalization.
func TestRulerControllerCanMoveAttackAndEndTurn(t *testing.T) {
	ruler := NewCompleteRuler()
	ctrl := NewFake("Fake1")
	ctrl2 := NewFake("Fake2")
	// Manually assign entities
	i := 0
	for id, e := range ruler.GameState.Entities {
		if i == 0 {
			e.ControllerID = ctrl.ID
		} else {
			e.ControllerID = ctrl2.ID
		}
		ruler.GameState.Entities[id] = e
		i++
	}

	// Flatten grid to avoid height-blocking in 4x4 tests
	for x := 0; x < ruler.GameState.Grid.Width; x++ {
		for y := 0; y < ruler.GameState.Grid.Length; y++ {
			z := ruler.GameState.Grid.TopMostGroundAt(x, y)
			if cell, ok := ruler.GameState.Grid.CellAt(position.New(x, y, z)); ok {
				if z != 0 {
					delete(ruler.GameState.Grid.Cells, cell.Position)
					cell.Position.Z = 0
					ruler.GameState.Grid.Cells[cell.Position] = cell
				}
			}
		}
	}
	for id, e := range ruler.GameState.Entities {
		e.Position.Z = 0
		ruler.GameState.Entities[id] = e
	}

	ruler.Start()
	defer ruler.Stop()
	dChan := make(chan *message.Message, 1)
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl, ControllerID: ctrl.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan
	ruler.SendActor(message.Create(nil, rulermethods.AddController{Controller: ctrl2, ControllerID: ctrl2.ID}, rulermethods.AddControllerReply{}), dChan)
	<-dChan

	ctrl.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)
	ctrl2.ExpectMessage(t, rulermethods.BattleStart{}, 5*time.Second)

	replyChan := make(chan *message.Message)

	ruler.SendActor(message.Create(nil, rulermethods.GetGridState{}, rulermethods.GetGridStateReply{}), replyChan)
	msg := <-replyChan
	grd := msg.Content.(rulermethods.GetGridStateReply).Grid

	ruler.SendActor(message.Create(nil, rulermethods.GetEntitiesState{}, rulermethods.GetEntitiesStateReply{}), replyChan)
	msg = <-replyChan

	var entities []entity.Entity
	var foeEntities []entity.Entity
	for _, ent := range msg.Content.(rulermethods.GetEntitiesStateReply).Entities {
		if ent.ControllerID == ctrl.ID {
			entities = append(entities, ent)
		} else {
			foeEntities = append(foeEntities, ent)
		}
	}

	attacker := entities[0]
	target := foeEntities[0]

	done := false
	turnTimeout := time.After(15 * time.Second)

	hasTestedMove := false
	hasTestedAttack := false

	for !done {
		select {
		case msg := <-ctrl.Inbox:
			if msg.TargetMethod != nil {
				switch m := msg.TargetMethod.(type) {
				case rulermethods.ControllerNextTurn:
					if m.Entity.ID == attacker.ID {
						if !hasTestedMove {
							var nextPos position.Position
							found := false

							for dx := -1; dx <= 1; dx++ {
								for dy := -1; dy <= 1; dy++ {
									if dx == 0 && dy == 0 {
										continue
									}
									// Only orthogonal adjacency is valid in many turn-based grids, but just in case, check proper neighbors.
									if dx != 0 && dy != 0 {
										continue
									}

									p := position.New(attacker.Position.X+dx, attacker.Position.Y+dy, 0)
									p.Z = grd.TopMostGroundAt(p.X, p.Y)
									logrus.Debugf("Attacker Position: %s destination %s", attacker.Position.String(), p.String())

									if !p.IsAdjacent(attacker.Position, 2) {
										continue
									}

									zDiff := p.Z - attacker.Position.Z
									if zDiff < 0 {
										zDiff = -zDiff
									}

									if c, ok := grd.CellAt(p); ok && !c.IsOccupied() && zDiff <= 2 {
										nextPos = p
										found = true
										break
									}
								}
								if found {
									break
								}
							}

							movePath := []position.Position{}
							if found {
								movePath = append(movePath, nextPos)
							}

							ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
								EntityID:     attacker.ID,
								Path:         movePath,
								ControllerID: ctrl.ID,
							}, rulermethods.ControllerMoveReply{}), ctrl.GetCallbackChan())
						} else if !hasTestedAttack {
							ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
								EntityID:     attacker.ID,
								Target:       target.Position,
								ControllerID: ctrl.ID,
							}, rulermethods.ControllerAttackReply{}), ctrl.GetCallbackChan())
						} else {
							ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
								EntityID:     m.Entity.ID,
								ControllerID: ctrl.ID,
							}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
							done = true
						}
					} else {
						ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
							EntityID:     m.Entity.ID,
							ControllerID: ctrl.ID,
						}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
					}
				case rulermethods.ControllerMoveReply:
					hasTestedMove = true
					ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
						EntityID:     attacker.ID,
						ControllerID: ctrl.ID,
					}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
				case rulermethods.ControllerAttackReply:
					hasTestedAttack = true
					done = true // short-circuit end test!
					ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
						EntityID:     attacker.ID,
						ControllerID: ctrl.ID,
					}, rulermethods.EndOfTurn{}), ctrl.GetCallbackChan())
				}
			}
		case msg := <-ctrl2.Inbox:
			if msg.TargetMethod != nil {
				switch m := msg.TargetMethod.(type) {
				case rulermethods.ControllerNextTurn:
					ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
						EntityID:     m.Entity.ID,
						ControllerID: ctrl2.ID,
					}, rulermethods.EndOfTurn{}), ctrl2.GetCallbackChan())
				}
			}
		case <-turnTimeout:
			t.Fatal("Timeout in battle simulation")
		}
	}

	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{ControllerID: ctrl.ID}, nil))
	ruler.NotifyActor(message.Create(nil, rulermethods.ControllerQuit{ControllerID: ctrl2.ID}, nil))

	ctrl.Stop()
	ctrl2.Stop()
}
