package gridgenerator

import (
	"fmt"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/cell"
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/position/pattern"
	"github.com/ecumeurs/upsilonbattle/tools"
)

type GridType int

const (
	Flat     GridType = 0
	Hill     GridType = 1
	River    GridType = 2
	Mountain GridType = 3
)

type GridGenerator struct {
	Width  tools.IntRange
	Length tools.IntRange
	Height tools.IntRange

	GenerateObstrcution bool
	ObstructionRate     tools.IntRange

	Type GridType
}

func fillBellowGroundWithDirt(g *grid.Grid) {
	for _, c := range g.Cells {
		if c.Type == cell.Dirt {
			continue
		}
		for z := c.Position.Z - 1; z >= 0; z-- {
			if z > c.Position.Z-5 {
				_, ok := g.Cells[position.New(c.Position.X, c.Position.Y, z)]
				if !ok {
					d := new(cell.Cell)
					d.Position = position.New(c.Position.X, c.Position.Y, z)
					d.Type = cell.Dirt
					g.Cells[d.Position] = d
				} else {
					g.ReplaceCellType(position.New(c.Position.X, c.Position.Y, z), cell.Dirt)
				}
			} else {
				// remove all cells below the ground
				delete(g.Cells, position.New(c.Position.X, c.Position.Y, z))
			}
		}
	}
	ensureNoDirtIsVisibleOnTop(g)
}

func ensureNoDirtIsVisibleOnTop(g *grid.Grid) {
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Length; y++ {
			z := g.TopMostGroundAt(x, y)
			if c, ok := g.Cells[position.New(x, y, z)]; ok {
				if c.Type == cell.Dirt {
					c.Type = cell.Ground
				}
			}
		}
	}
}

// Generate a Grid. Note that you only need to keep cells where there is a ground or an obstacle
func (g *GridGenerator) Generate() (res *grid.Grid) {

	switch g.Type {
	case Flat:
		res = g.generateFlat()
	case Hill:
		res = g.generateHill()
	case River:
		res = g.generateRiver()
	}

	fillBellowGroundWithDirt(res)

	return res
}

func (g *GridGenerator) generateFlat() *grid.Grid {
	gr := new(grid.Grid)
	gr.Cells = make(map[position.Position]*cell.Cell)
	gr.Width = g.Width.Random()
	gr.Length = g.Length.Random()
	gr.Height = g.Height.Random()
	obstruction := g.ObstructionRate.Random()

	// determine a random height for the ground and create the ground, do not create cell where there is no ground (above or below)
	// add a random number of obstacles based on obstruction and replace some top most ground by obstacles
	// while the ground may be flat, it is expected a slight variation in height
	// the ground is expected to be flat in the middle and to have a slight slope on the sides

	// determine a random height for the ground
	ground_height := tools.RandomInt(2, gr.Height)
	fmt.Println("ground height", ground_height)

	// create the ground, do not create cell where there is no ground (above or below)
	for x := 0; x < gr.Width; x++ {
		for y := 0; y < gr.Length; y++ {
			c := new(cell.Cell)
			c.Position = position.New(x, y, ground_height)
			c.Type = cell.Ground
			gr.Cells[c.Position] = c
		}
	}

	// add minor variations in height
	for x := 0; x < gr.Width; x++ {
		for y := 0; y < gr.Length; y++ {
			z := gr.TopMostGroundAt(x, y)
			rand := tools.RandomInt(0, 100)
			if rand > 90 {
				c := new(cell.Cell)
				c.Position = position.New(x, y, z+1)
				c.Type = cell.Ground
				gr.ReplaceCell(position.New(x, y, z), c)
			} else if rand > 88 {
				c := new(cell.Cell)
				c.Position = position.New(x, y, z-1)
				c.Type = cell.Ground
				gr.ReplaceCell(position.New(x, y, z), c)
			}
		}
	}

	if g.GenerateObstrcution {
		// add a random number of obstacles based on obstruction and replace some top most ground by obstacles
		for i := 0; i < obstruction; i++ {
			x := tools.RandomInt(0, gr.Width-1)
			y := tools.RandomInt(0, gr.Length-1)
			z := gr.TopMostGroundAt(x, y)
			gr.ReplaceCellType(position.New(x, y, z), cell.Obstacle)
		}
	}

	return gr
}

func adjascentHeights(g *grid.Grid, p position.Position) []int {
	adjascent := []int{}
	for _, pt := range pattern.Neighbours() {
		if p.Add(pt).Equals(p) {
			continue
		}
		if c, ok := g.Cells[p.Add(pt)]; ok {
			adjascent = append(adjascent, c.Position.Z)
		}
	}
	return adjascent
}

func adjascentNearestHeight(g *grid.Grid, p position.Position, target_height int) int {
	adjascent := adjascentHeights(g, p)
	min := 999999
	for _, h := range adjascent {
		if h < min {
			min = h
		}
	}
	return min
}

func (g *GridGenerator) generateHill() *grid.Grid {
	ob := g.GenerateObstrcution
	g.GenerateObstrcution = false
	gr := g.generateFlat()
	g.GenerateObstrcution = ob
	for it := 0; it < tools.RandomInt(1, 5); it++ {
		x := tools.RandomInt(0, gr.Width)
		y := tools.RandomInt(0, gr.Length)
		z := gr.TopMostGroundAt(x, y)
		hillsize := tools.RandomInt(8, 12)
		hillheight := z + tools.RandomInt(3, 8)

		if hillheight >= gr.Height {
			hillheight = gr.Height - 1
		}

		// update the ground cells height to create the hill. We do not create new cells, we only update the height of the existing ones
		// ensure the slope is never higher that 2 in height from adjascent cells, apply a curve to the slope
		// we also remove all cells above the new height
		// be sure that we only edit cells that are in the grid
		for x2 := x - hillsize; x2 < x+hillsize; x2++ {
			for y2 := y - hillsize; y2 < y+hillsize; y2++ {
				if x2 >= 0 && x2 < gr.Width && y2 >= 0 && y2 < gr.Length {
					// calculate the distance from the center of the hill
					dist := tools.FloatDistance(float64(x2), float64(y2), float64(x), float64(y))
					// calculate the height of the ground at this distance from the center of the hill
					height := hillheight + int(float64(hillheight)*(1.0-dist/float64(hillsize)))
					nearestHeight := adjascentNearestHeight(gr, position.New(x2, y2, height), height)
					if nearestHeight < height-2 {
						height = nearestHeight + 2
					}
					if height >= gr.Height {
						height = gr.Height - 1
					}

					c := cell.Cell{}
					c.Position = position.New(x2, y2, height)
					c.Type = cell.Ground
					gr.Cells[c.Position] = &c
					gr.ReplaceCell(position.New(x2, y2, gr.TopMostGroundAt(x2, y2)), &c)
				}
			}
		}
	}

	if g.GenerateObstrcution {
		// add a random number of obstacles based on obstruction and replace some top most ground by obstacles
		obstruction := g.ObstructionRate.Random()
		for i := 0; i < obstruction; i++ {
			x := tools.RandomInt(0, gr.Width)
			y := tools.RandomInt(0, gr.Length)
			z := gr.TopMostGroundAt(x, y)
			gr.ReplaceCellType(position.New(x, y, z), cell.Obstacle)
		}
	}
	return gr
}

// GenerateRiver generates a river from a flat grid. A river is a path of water cells that are connected to the ground. Height of the river may be lower than the ground.
// A river must run through the whole grid
// There might be a branching in the river, but no more than 2.
// The river width may not be constant, but it is expected to be between 3 and 5 cells wide
// The river is expected to be flat in the middle and to have a slight slope on the sides

func (g *GridGenerator) generateRiver() *grid.Grid {
	gr := g.generateFlat()

	// determine river origin point
	origin := position.RandomBorderPosition(gr.Width, gr.Length, gr.Height)
	origin, _ = gr.ForcePositionToGround(origin)
	endpoints := []position.Position{}
	for i := 0; i < tools.RandomInt(1, 3); i++ {
		// determine river end point
		end := position.RandomBorderPosition(gr.Width, gr.Length, gr.Height)
		end, _ = gr.ForcePositionToGround(end)
		endpoints = append(endpoints, end)
	}

	// determine river width
	width := tools.RandomInt(2, 4)

	// determine river depth
	waterlevel := tools.Min(tools.RandomInt(1, 3), tools.Max(0, gr.FindLowestLevel()-1))

	waypoints := []position.Position{}
	for i := 0; i < tools.RandomInt(1, 3); i++ {
		// determine river waypoints
		waypoint := position.RandomPosition(gr.Width, gr.Length, gr.Height)
		waypoint, _ = gr.ForcePositionToGround(waypoint)
		waypoints = append(waypoints, waypoint)
	}

	currentRiverPosition := origin
	for _, waypoint := range waypoints {
		// generate river from current position to waypoint
		currentRiverPosition = g.generateRiverSegment(gr, currentRiverPosition, waypoint, width, waterlevel)
	}

	for _, endpoint := range endpoints {
		// seek nearest river position to endpoint
		nearestRiverCell, found := gr.FindNearestCellMatchingPredicate(endpoint, func(c *cell.Cell) bool {
			return c.Type == cell.Water
		})

		if !found {
			fmt.Printf("No river found near endpoint %v\n", endpoint)
		} else if nearestRiverCell.Position.Distance(endpoint) > 1 {
			fmt.Printf("Generating River from %v to %v\n", currentRiverPosition, nearestRiverCell.Position)

			// generate river from current position to endpoint
			currentRiverPosition = g.generateRiverSegment(gr, nearestRiverCell.Position, endpoint, tools.RandomInt(2, 3), waterlevel)
		} else {
			fmt.Printf("River already present at endpoint %v\n", endpoint)
		}

	}

	if g.GenerateObstrcution {
		obstruction := g.ObstructionRate.Random()
		// add a random number of obstacles based on obstruction and replace some top most ground by obstacles
		for i := 0; i < obstruction; i++ {

			x := tools.RandomInt(0, gr.Width-1)
			y := tools.RandomInt(0, gr.Length-1)
			z := gr.TopMostGroundAt(x, y)
			if c, found := gr.CellAt(position.New(x, y, z)); found && c.Type == cell.Water {
				i--
			} else {
				gr.ReplaceCellType(position.New(x, y, z), cell.Obstacle)
			}
		}
	}
	return gr
}

// generateRiverSegment generates a river segment from a start position to an end position.
func (g *GridGenerator) generateRiverSegment(gr *grid.Grid, start position.Position, end position.Position, width int, waterlevel int) position.Position {
	// find all position between start and end
	for _, p := range pattern.PathTo2D(end.Substract(start)).Enlarge(width) {
		// we also remove all cells above the new height
		for z := waterlevel; z < gr.Height; z++ {
			delete(gr.Cells, position.New(start.X+p.X, start.Y+p.Y, z))
		}
		// be sure that we only edit cells that are in the grid
		c := cell.Cell{}
		c.Position = position.New(start.X+p.X, start.Y+p.Y, waterlevel)
		c.Type = cell.Water
		gr.Cells[c.Position] = &c
	}
	return end
}
