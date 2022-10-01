package grid

import (
	"fmt"
	"math/rand"

	"github.com/ecumeurs/upsilonbattle/battlearena/grid/cell"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position"
	"github.com/ecumeurs/upsilonbattle/battlearena/grid/position/pattern"
	"github.com/ecumeurs/upsilontools/tools"
	"github.com/google/uuid"
)

type Grid struct {
	Width  int
	Length int
	Height int

	Cells map[position.Position]*cell.Cell
}

// NewGrid Makes a flat grid at groundlevel.
func NewGrid(width, length, groundlevel int) *Grid {
	g := &Grid{
		Width:  width,
		Length: length,
		Height: groundlevel + 2,
		Cells:  make(map[position.Position]*cell.Cell),
	}
	for x := 0; x < width; x++ {
		for y := 0; y < length; y++ {
			for z := 0; z < groundlevel+1; z++ {
				if z < groundlevel {
					g.Cells[position.New(x, y, z)] = cell.NewCell(cell.Dirt, position.New(x, y, z))
				} else {
					g.Cells[position.New(x, y, z)] = cell.NewCell(cell.Ground, position.New(x, y, z))
				}
			}
		}
	}
	return g
}

// RandomPosition returns a random valid position in the grid
func (g *Grid) RandomPosition() position.Position {
	for {
		x := rand.Intn(g.Width)
		y := rand.Intn(g.Length)
		z := g.TopMostGroundAt(x, y)
		pos := position.New(x, y, z)
		// check if there is a ground cell
		if _, found := g.CellAt(pos); found {
			return pos
		}
	}
}

// MoveEntity moves an entity from one position to another
func (g *Grid) MoveEntity(from, to position.Position, EntityID uuid.UUID) error {
	if !g.Contains(to) {
		return fmt.Errorf("to position %v is not in the grid", to)
	}
	c, ok := g.CellAt(from)
	if ok {
		c.EntityID = uuid.Nil
	}
	c, ok = g.CellAt(to)
	if ok {
		c.EntityID = EntityID
	} else {
		return fmt.Errorf("to position %v is not in the grid", to)
	}
	return nil
}

// RemoveEntity
func (g *Grid) RemoveEntity(p position.Position) {
	if !g.Contains(p) {
		return
	}
	g.Cells[p].EntityID = uuid.Nil
}

// CellAt
func (g *Grid) CellAt(p position.Position) (*cell.Cell, bool) {
	if g.PositionIsInGrid(p) {
		c, ok := g.Cells[p]
		return c, ok
	}
	return nil, false
}

func (g *Grid) TopMostCellAt(x, y int) int {
	for z := g.Height; z >= 0; z-- {
		pos := position.New(x, y, z)
		_, found := g.Cells[pos]
		if found {
			return z
		}
	}
	return -1
}

func (g *Grid) LowestCellAt(x, y int) int {
	for z := 0; z <= g.Height; z++ {
		pos := position.New(x, y, z)
		if _, found := g.Cells[pos]; found {
			return z
		}
	}
	return -1
}

func (g *Grid) AllCellsAt(x, y int) []int {
	res := []int{}
	for z := 0; z <= g.Height; z++ {
		pos := position.New(x, y, z)
		if _, found := g.Cells[pos]; found {
			res = append(res, z)
		}
	}
	return res
}

func (g *Grid) TopMostGroundAt(x, y int) int {
	for z := g.Height; z >= 0; z-- {
		pos := position.New(x, y, z)
		c, found := g.Cells[pos]
		if found && c.Type == cell.Ground {
			return z
		}
	}
	return -1
}

func (g *Grid) LowestGroundAt(x, y int) int {
	for z := 0; z <= g.Height; z++ {
		pos := position.New(x, y, z)
		if c, found := g.Cells[pos]; found && c.Type == cell.Ground {
			return z
		}
	}
	return -1
}

func (g *Grid) AllGroundAt(x, y int) []int {
	res := []int{}
	for z := 0; z <= g.Height; z++ {
		pos := position.New(x, y, z)
		if c, found := g.Cells[pos]; found && c.Type == cell.Ground {
			res = append(res, z)
		}
	}
	return res
}

// PositionIsInGrid returns true if the position is in the grid
func (g *Grid) PositionIsInGrid(p position.Position) bool {
	return p.X >= 0 && p.X < g.Width && p.Y >= 0 && p.Y < g.Length && p.Z >= 0 && p.Z < g.Height
}

// ForcePositionToGround
func (g *Grid) ForcePositionToGround(p position.Position) (position.Position, bool) {
	if g.PositionIsInGrid(p) {
		return position.New(p.X, p.Y, g.TopMostGroundAt(p.X, p.Y)), true
	}
	return position.New(-1, -1, -1), false
}

// Display the grid
func (g *Grid) Display() {
	for z := 0; z < g.Height; z++ {
		for y := 0; y < g.Length; y++ {
			for x := 0; x < g.Width; x++ {
				pos := position.New(x, y, z)
				c, ok := g.Cells[pos]
				if !ok {
					continue
				}
				switch c.Type {
				case cell.Ground:
					if c.EntityID == uuid.Nil {
						print(".")
					} else {
						print("x")
					}
				case cell.Obstacle:
					print("#")
				}
			}
			println()
		}
		println()
		println()
	}
}

func (g *Grid) RemoveCell(p position.Position) {
	delete(g.Cells, p)
}

// ReplaceCell expects the new Cell to be at the same x,y position as the old one (z can be different)
func (g *Grid) ReplaceCell(p position.Position, c *cell.Cell) {
	if !g.Contains(p) {
		return
	}
	if p.SameXY(c.Position) {
		g.RemoveCell(p)
		g.Cells[c.Position] = c
	}
}

func (g *Grid) ReplaceCellType(p position.Position, t cell.CellType) {
	if !g.Contains(p) {
		return
	}
	g.Cells[p].Type = t
}

func lighterColor(height, maxheight, basecolor, maxcolor int) int {
	return basecolor + (maxcolor-basecolor)*height/maxheight
}

func hexColor(r, g, b int) string {
	return fmt.Sprintf("0x%02x%02x%02x", r, g, b)
}

// GenerateCellAsObeliskCube returns the cell as an obelisk cube
func (g *Grid) generateCellAsObeliskCube(p position.Position) string {
	res := ""
	res += fmt.Sprintf("position = new obelisk.Point3D(%d, %d, %d);\n", p.X*20, p.Y*20, p.Z*20)
	switch g.Cells[p].Type {
	case cell.Ground:
		if g.Cells[p].EntityID != uuid.Nil {
			res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 0, 80), lighterColor(p.Z, g.Height, 0, 80), lighterColor(p.Z, g.Height, 0, 80)) + ");\n"
		} else {
			res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 0, 180), 255, lighterColor(p.Z, g.Height, 0, 180)) + ");\n"
		}
	case cell.Obstacle:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, lighterColor(p.Z, g.Height, 80, 160), lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	case cell.Water:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(lighterColor(p.Z, g.Height, 80, 160), lighterColor(p.Z, g.Height, 50, 130), 255) + ");\n"
	case cell.Dirt:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(150, 155, lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	case cell.Debug:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, lighterColor(p.Z, g.Height, 50, 130), 255) + ");\n"
	case cell.Debug2:
		res += "color = new obelisk.CubeColor().getByHorizontalColor(" + hexColor(255, 255, lighterColor(p.Z, g.Height, 50, 130)) + ");\n"
	}
	res += "cube = new obelisk.Cube(dimension, color, false);\n"
	res += "pixelView.renderObject(cube, position);\n"
	return res
}

func (g *Grid) GenerateHTML() string {
	res := "<html><head><script src=\"https://unpkg.com/obelisk.js@1.2.2/build/obelisk.min.js\"></script></script></head>\n<body><canvas id=\"canvas-demo\" width=\"5000\" height=\"5000\"></canvas>\n"
	res += "<script>var canvas = document.getElementById('canvas-demo');\n"

	res += "var point = new obelisk.Point(1000, 500);var pixelView = new obelisk.PixelView(canvas,point);\n"

	res += "// create dimension instance \n"
	res += "var dimension = new obelisk.CubeDimension(20, 20, 20);\n"
	res += "// create color instance \n"
	res += "var color = new obelisk.CubeColor().getByHorizontalColor(obelisk.ColorPattern.BLUE);\n"
	res += "var cube = new obelisk.Cube(dimension, color, true);\n"
	res += "var position = new obelisk.Point3D(0, 0, 0);\n"

	// render all cubes
	for z := 0; z < g.Height; z++ {
		for y := 0; y < g.Length; y++ {
			for x := 0; x < g.Width; x++ {
				pos := position.New(x, y, z)
				_, ok := g.Cells[pos]
				if !ok {
					continue
				}
				res += g.generateCellAsObeliskCube(pos)
			}
		}
	}

	res += "</script></body></html>"
	return res
}

func (g *Grid) FindLowestLevel() int {
	lowest := 0
	for _, c := range g.Cells {
		if c.Position.Z < lowest {
			lowest = c.Position.Z
		}
	}
	return lowest
}

// FindNearestCellMatchingPredicate returns the nearest cell matching the predicate
func (g *Grid) FindNearestCellMatchingPredicate(p position.Position, predicate func(*cell.Cell) bool) (*cell.Cell, bool) {
	p.Z = g.LowestCellAt(p.X, p.Y)
	if !g.Contains(p) {
		return nil, false
	}

	if predicate(g.Cells[p]) {
		return g.Cells[p], true
	}

	visited := map[position.Position]bool{p: true}
	queue := []position.Position{p}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		for _, n := range g.SelectPositionsByPattern2D(pos, pattern.Neighbours2D()) {
			if visited[n] {
				continue
			}
			visited[n] = true
			if predicate(g.Cells[n]) {
				return g.Cells[n], true
			}
			queue = append(queue, n)
		}
	}
	return nil, false
}

func (g *Grid) Contains(p position.Position) bool {
	_, ok := g.Cells[p]
	return ok
}

// SelectPositionsByPattern returns the positions matching the pattern
func (g *Grid) SelectPositionsByPattern(origin position.Position, pat pattern.Pattern) []position.Position {
	res := []position.Position{}
	pos := pat.ApplyInArea(origin, g.Width, g.Length, g.Height)
	for _, p := range pos {
		if g.Contains(p) {
			res = append(res, p)
		}
	}
	return res
}

func (g *Grid) SelectPositionsByPattern2D(origin position.Position, pat pattern.Pattern2D) []position.Position {
	res := []position.Position{}
	for _, p := range pat {
		pos := origin.Add(p)
		pos.Z = g.TopMostCellAt(pos.X, pos.Y)
		if g.Contains(pos) {
			res = append(res, pos)
		}
	}
	return res
}

// CellsForPositions returns the cells for the given positions
func (g *Grid) CellsForPositions(pos []position.Position) []*cell.Cell {
	res := []*cell.Cell{}
	for _, p := range pos {
		c, ok := g.Cells[p]
		if ok {
			res = append(res, c)
		}
	}
	return res
}

// AStarPath returns the path from start to end using A*
func (g *Grid) AStarPath(start, end position.Position, jumpHeight int) ([]position.Position, bool) {
	if !g.Contains(start) || !g.Contains(end) {
		return nil, false
	}
	if start == end {
		return []position.Position{start}, true
	}

	visited := map[position.Position]int{}
	queue := []position.Position{start}
	parents := map[position.Position]position.Position{}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		if pos.Equals(end) {
			return g.reconstructPath(visited, start, end, jumpHeight), true
		}
		visited[pos] = visited[parents[pos]] + 1
		for _, n := range g.SelectPositionsByPattern2D(pos, pattern.Neighbours2D()) {
			if _, found := visited[n]; found {
				continue
			}
			if tools.Abs(n.Z-pos.Z) > jumpHeight {
				continue
			}
			if c, found := g.CellAt(n); found && c.Type == cell.Ground {
				queue = append(queue, n)
				parents[n] = pos
			}
		}
	}
	return nil, false
}

func (g *Grid) reconstructPath(visited map[position.Position]int, start, end position.Position, jumpHeight int) []position.Position {
	res := []position.Position{end}
	for res[len(res)-1] != start {
		// find the lowest number within adjascents
		lowest := 999999
		var lowestPos position.Position
		for _, n := range g.SelectPositionsByPattern2D(res[len(res)-1], pattern.Neighbours2D()) {
			if tools.Abs(n.Z-res[len(res)-1].Z) <= jumpHeight {
				if vis, found := visited[n]; found && vis < lowest {
					lowest = vis
					lowestPos = n
				}
			}
		}
		res = append(res, lowestPos)
	}

	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
