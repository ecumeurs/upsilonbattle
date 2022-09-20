package cell

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/google/uuid"
)

// CellType describe all the possible types of cell
type CellType int

const (
	// WallCellType is the type of a wall cell
	Obstacle CellType = 0
	// GroundCellType is the type of a ground cell
	Ground CellType = 1
	Water  CellType = 2
	Dirt   CellType = 3
	Debug  CellType = 4
	Debug2 CellType = 5
)

// Cell is a cell in the grid
type Cell struct {
	Type     CellType
	EntityID uuid.UUID
	Position position.Position
}
