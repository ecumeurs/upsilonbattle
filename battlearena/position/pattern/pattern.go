package pattern

import (
	"github.com/ecumeurs/upsilonbattle/battlearena/position"
	"github.com/ecumeurs/upsilontools/tools"
)

type Pattern []position.Position

func (p Pattern) Contains(origin, pos position.Position) bool {
	for _, p := range p {
		if origin.Add(p).Equals(pos) {
			return true
		}
	}
	return false
}

func (p Pattern) ContainsAny(origin position.Position, positions []position.Position) bool {
	for _, pos := range positions {
		if p.Contains(origin, pos) {
			return true
		}
	}
	return false
}

func (p Pattern) ContainsAll(origin position.Position, positions []position.Position) bool {
	for _, pos := range positions {
		if !p.Contains(origin, pos) {
			return false
		}
	}
	return true
}

// Circle creates a circle pattern
func Circle(radius int) Pattern {
	p := Pattern{}
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			for z := -radius; z <= radius; z++ {
				if x*x+y*y+z*z <= radius*radius {
					p = append(p, position.New(x, y, z))
				}
			}
		}
	}
	return p
}

// Square creates a square pattern
func Square(width, length, height int) Pattern {
	p := Pattern{}
	for x := -width; x <= width; x++ {
		for y := -length; y <= length; y++ {
			for z := -height; z <= height; z++ {
				p = append(p, position.New(x, y, z))
			}
		}
	}
	return p
}

// Line creates a line pattern
func Line(length int) Pattern {
	p := Pattern{}
	for i := 0; i < length; i++ {
		p = append(p, position.New(i, 0, 0))
	}
	return p
}

// Neighbour creates a neighbour pattern
func Neighbours() Pattern {
	return Square(1, 1, 1)
}

func AdjacentDirections() Pattern {
	return Pattern{
		position.New(0, 0, 1),
		position.New(0, 0, -1),
		position.New(0, 1, 0),
		position.New(0, -1, 0),
		position.New(1, 0, 0),
		position.New(-1, 0, 0),
	}
}

// Apply applies a pattern to a position
func (p Pattern) Apply(origin position.Position) []position.Position {
	result := []position.Position{}
	for _, pos := range p {
		result = append(result, origin.Add(pos))
	}
	return result
}

// ApplyInArea applies a pattern to a position in an area
func (p Pattern) ApplyInArea(origin position.Position, width, length, height int) []position.Position {
	result := []position.Position{}
	for _, pos := range p {
		pos = origin.Add(pos)
		if pos.X >= 0 && pos.X < width && pos.Y >= 0 && pos.Y < length && pos.Z >= 0 && pos.Z < height {
			result = append(result, pos)
		}
	}
	return result
}

// PathTo Generate a path from 0,0,0 to another Position as directly as possible.
func PathTo(pos position.Position) Pattern {
	p := Pattern{}
	for {
		if pos.X == 0 && pos.Y == 0 && pos.Z == 0 {
			break
		}
		p = append(p, pos)
		if pos.X > 0 {
			pos.X--
		} else if pos.X < 0 {
			pos.X++
		}
		if pos.Y > 0 {
			pos.Y--
		} else if pos.Y < 0 {
			pos.Y++
		}
		if pos.Z > 0 {
			pos.Z--
		} else if pos.Z < 0 {
			pos.Z++
		}
	}
	return p
}

// Enlarge makes a pattern bigger by adding a layer of cells around it
func (p Pattern) Enlarge(radius int) Pattern {
	result := Pattern{}
	for _, pos := range p {
		for x := -radius; x <= radius; x++ {
			for y := -radius; y <= radius; y++ {
				for z := -radius; z <= radius; z++ {
					if !position.ContainedIn(pos.Add(position.New(x, y, z)), result) {
						result = append(result, pos.Add(position.New(x, y, z)))
					}
				}
			}
		}
	}
	return result
}

// Enlarge makes a pattern bigger by adding a layer of cells around it
func (p Pattern) EnlargeVarying(r tools.IntRange) Pattern {
	result := Pattern{}
	for _, pos := range p {
		radius := r.Random()
		for x := -radius; x <= radius; x++ {
			for y := -radius; y <= radius; y++ {
				for z := -radius; z <= radius; z++ {
					if !position.ContainedIn(pos.Add(position.New(x, y, z)), result) {
						result = append(result, pos.Add(position.New(x, y, z)))
					}
				}
			}
		}
	}
	return result
}

type Pattern2D []position.Position

func Square2D(width, length int) Pattern2D {
	p := Pattern2D{}
	for x := -width; x <= width; x++ {
		for y := -length; y <= length; y++ {
			p = append(p, position.New(x, y, 0))
		}
	}
	return p
}

func Circle2D(radius int) Pattern2D {
	p := Pattern2D{}
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			if x*x+y*y <= radius*radius {
				p = append(p, position.New(x, y, 0))
			}
		}
	}
	return p
}

func Line2D(length int) Pattern2D {
	p := Pattern2D{}
	for i := 0; i < length; i++ {
		p = append(p, position.New(i, 0, 0))
	}
	return p
}

func Neighbours2D() Pattern2D {
	return Square2D(1, 1)
}

func PathTo2D(pos position.Position) Pattern2D {
	p := Pattern2D{}
	for {
		if pos.X == 0 && pos.Y == 0 {
			break
		}
		p = append(p, pos)
		if pos.X > 0 {
			pos.X--
		} else if pos.X < 0 {
			pos.X++
		}
		if pos.Y > 0 {
			pos.Y--
		} else if pos.Y < 0 {
			pos.Y++
		}
	}
	return p
}

// Enlarge makes a pattern bigger by adding a layer of cells around it
func (p Pattern2D) Enlarge(radius int) Pattern2D {
	result := Pattern2D{}
	for _, pos := range p {
		for x := -radius; x <= radius; x++ {
			for y := -radius; y <= radius; y++ {
				if !position.ContainedIn(pos.Add(position.New(x, y, 0)), result) {
					result = append(result, pos.Add(position.New(x, y, 0)))
				}
			}
		}
	}
	return result
}

// Enlarge makes a pattern bigger by adding a layer of cells around it
func (p Pattern2D) EnlargeVarying(r tools.IntRange) Pattern2D {
	result := Pattern2D{}
	for _, pos := range p {
		radius := r.Random()
		for x := -radius; x <= radius; x++ {
			for y := -radius; y <= radius; y++ {
				if !position.ContainedIn(pos.Add(position.New(x, y, 0)), result) {
					result = append(result, pos.Add(position.New(x, y, 0)))
				}
			}
		}
	}
	return result
}
