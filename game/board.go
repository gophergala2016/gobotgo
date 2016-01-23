package game

import "fmt"

type Board [][]intersection

type intersection int

// Player types
const (
	empty = intersection(iota)
	black
	white
)

func (b Board) set(m Move) error {
	if err := b.intersectionEmpty(m); err != nil {
		return err
	}
	b[m.X][m.Y] = intersection(m.Player)
	return nil
}

func (b Board) equal(c Board) error {
	d := b.slice()
	e := c.slice()

	for i := range d {
		if d[i] != e[i] {
			return fmt.Errorf("Board state not equal at %d", i)
		}
	}
	return nil
}

func (b Board) intersectionEmpty(m Move) error {
	i := b[m.X][m.Y]
	if i != empty {
		return fmt.Errorf("Intersection %d-%d is not empty", m.X, m.Y)
	}
	return nil
}

func (b Board) bounded(x, y int) bool {
	color := b[x][y]
	if color == empty {
		return false
	}
	mask := newBoard(len(b))
	mask[x][y] = color
	// We're probably going to allocate somewhat initially, so lets allocate a bit
	type pos struct{ x, y int }
	frontier := make([]pos, 0, 64)
	frontier = append(frontier, pos{x, y})

	moves := [4]pos{
		{0, 1},
		{1, 0},
		{0, -1},
		{-1, 0},
	}

	rangeCheck := func(i int) bool {
		return i >= 0 && i < len(b)
	}

	// Walk the frontier
	for len(frontier) > 0 {
		f := frontier[0]
		frontier = frontier[1:]
		// Check canditates up, down, left, right
		// Look for a connected empty. That means we're not bounded
		for _, m := range moves {
			c := pos{f.x + m.x, f.y + m.y}
			switch {
			case !rangeCheck(c.x) || !rangeCheck(c.y):
			case b[c.x][c.y] == empty:
				return false
			case mask[c.x][c.y] == color:
				mask[c.x][c.y] = color
				frontier = append(frontier, c)
			default:
				// Dont' add the opponent's space to the frontier
			}
		}
	}

	// if we've exhausted the frontier, this is empty
	return true
}

func newBoard(size int) Board {
	return sliceBoard(make([]intersection, size*size), size)
}

func sliceBoard(i []intersection, size int) Board {
	if len(i) != size*size {
		panic("intersection list isn't size^2")
	}
	b := make(Board, size)
	// Only allocate once
	for row := range b {
		b[row] = i[:size]
		i = i[size:]
	}

	return b
}

func (b Board) slice() []intersection {
	if cap(b[0]) != len(b)*len(b[0]) {
		panic("board does not have entire allocation at board 0")
	}
	return b[0][:cap(b[0])]
}

func (b Board) copy() Board {
	l := len(b)
	a := make([]intersection, l*l)
	copy(a, b.slice())
	return sliceBoard(a, l)
}
