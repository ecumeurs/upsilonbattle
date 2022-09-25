package position

import (
	"fmt"
	"math"

	"github.com/ecumeurs/upsilontools/tools"
)

// Position is a position in the grid
type Position struct {
	X int
	Y int
	Z int
}

// String
func (p Position) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.X, p.Y, p.Z)
}

// New creates a new position
func New(x, y, z int) Position {
	return Position{X: x, Y: y, Z: z}
}

func (p Position) Equals(p2 Position) bool {
	return p.X == p2.X && p.Y == p2.Y && p.Z == p2.Z
}

// SameXY
func (p Position) SameXY(p2 Position) bool {
	return p.X == p2.X && p.Y == p2.Y
}

// Distance
func (p Position) Distance(p2 Position) int {
	return tools.Abs(p.X-p2.X) + tools.Abs(p.Y-p2.Y) + tools.Abs(p.Z-p2.Z)
}

// RandomPosition
func RandomPosition(width, length, height int) Position {
	return New(tools.RandomInt(0, width), tools.RandomInt(0, length), tools.RandomInt(0, height))
}

// RandomPositionInArea
func RandomPositionInArea(width, length, height, x, y, z int) Position {
	return New(tools.RandomInt(x, x+width), tools.RandomInt(y, y+length), tools.RandomInt(z, z+height))
}

func RandomBorderPosition(width, length, height int) Position {
	r := tools.RandomInt(0, 4)
	if r == 0 {
		return New(tools.RandomInt(0, width), 0, tools.RandomInt(0, height))
	} else if r == 1 {
		return New(tools.RandomInt(0, width), length-1, tools.RandomInt(0, height))
	} else if r == 2 {
		return New(0, tools.RandomInt(0, length), tools.RandomInt(0, height))
	}
	return New(width-1, tools.RandomInt(0, length), tools.RandomInt(0, height))
}

// Add adds a position to another
func (p Position) Add(p2 Position) Position {
	return New(p.X+p2.X, p.Y+p2.Y, p.Z+p2.Z)
}

func ContainedIn(p Position, positions []Position) bool {
	for _, pos := range positions {
		if p.Equals(pos) {
			return true
		}
	}
	return false
}

// Substract
func (p Position) Substract(p2 Position) Position {
	return New(p.X-p2.X, p.Y-p2.Y, p.Z-p2.Z)
}

// IsAdjacent tell wether a position is adjascent on x,y axes and allows height difference as provided by parameter
func (p Position) IsAdjacent(p2 Position, allowedJump int) bool {
	return (p.X == p2.X && tools.Abs(p.Y-p2.Y) == 1 && tools.Abs(p.Z-p2.Z) <= allowedJump) || (p.Y == p2.Y && tools.Abs(p.X-p2.X) == 1 && tools.Abs(p.Z-p2.Z) <= allowedJump)
}

// Return the angle to position in degree. Expect 0 to be north, 90 to be east, 180 to be south and 270 to be west
func (p Position) AngleTo(p2 Position) int {
	dx := p2.X - p.X
	dy := p2.Y - p.Y
	return int(math.Atan2(float64(dy), float64(dx)) * 180 / 3.14)
}
