package game

import (
	"fmt"
	"strings"
)

type Board [][]intersection

type intersection int

// Player types
const (
	empty = intersection(iota)
	black
	white
)

func (i intersection) String() string {
	switch i {
	case empty:
		return "."
	case black:
		return "b"
	case white:
		return "w"
	default:
		panic("(invalid intersection print)")
	}
}

func (b Board) String() string {
	rows := []string{}
	for _, row := range b {
		cols := []string{}
		for _, c := range row {
			cols = append(cols, c.String())
		}
		rows = append(rows, strings.Join(cols, " "))
	}
	rows = append(rows, "\n") // add trailing newline
	return strings.Join(rows, "\n")
}

func (b Board) valid(m Move) error {
	if m.X >= len(b) ||
		m.X < 0 ||
		m.Y >= len(b) ||
		m.Y < 0 {
		return ErrOutOfBounds
	}
	return nil
}

func (b Board) apply(m Move) (int, error) {
	if err := b.valid(m); err != nil {
		return 0, err
	}
	if err := b.intersectionEmpty(m.Position); err != nil {
		return 0, err
	}

	stone := intersection(m.Player)
	b.set(m.Position, stone)

	points := 0
	for _, p := range m.adjacent() {
		switch {
		case !b.rangeCheck(p):
		case b.get(p) == stone:
		default:
			if count := b.clearBounded(p); count > 0 {
				points += count
			}
		}
	}

	if points == 0 && b.bounded(m.Position) {
		b.set(m.Position, empty)
		return 0, fmt.Errorf("Bounded")
	}
	return points, nil
}

func (b Board) equal(c Board) error {
	switch {
	case b == nil && c == nil:
		return nil
	case c == nil:
		fallthrough
	case b == nil:
		return fmt.Errorf("board is nil")
	}
	d := b.slice()
	e := c.slice()

	for i := range d {
		if d[i] != e[i] {
			return fmt.Errorf("Board state not equal at %d", i)
		}
	}
	return nil
}

func (b Board) intersectionEmpty(p Position) error {
	i := b[p.X][p.Y]
	if i != empty {
		return fmt.Errorf("Intersection %d-%d is not empty", p.X, p.Y)
	}
	return nil
}

func (b Board) set(p Position, i intersection) Board {
	b[p.X][p.Y] = i
	return b
}

func (b Board) get(p Position) intersection {
	return b[p.X][p.Y]
}

func (b Board) rangeCheck(p Position) bool {
	return p.X >= 0 && p.X < len(b) && p.Y >= 0 && p.Y < len(b)
}

func (b Board) bounded(start Position) bool {
	return b.boundedMask(start) != nil
}

// Returns a mask of the bounded positions, or nil if none are bounded
func (b Board) boundedMask(start Position) Board {
	color := b.get(start)

	if color == empty {
		return nil
	}
	mask := newBoard(len(b)).set(start, color)
	// We're probably going to allocate somewhat initially, so lets allocate a bit
	frontier := make([]Position, 0, 64)
	frontier = append(frontier, start)

	// Walk the frontier
	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		// Check canditates up, down, left, right
		// Look for a connected empty. That means we're not bounded
		for _, adj := range current.adjacent() {
			switch {
			case !b.rangeCheck(adj):
				// Don't add out of range positions to the frontier
			case mask.get(adj) == color:
				// Don't add previously checked positions to the frontier
			case b.get(adj) == empty:
				// Not bounded empty connected
				return nil
			case b.get(adj) == color:
				mask.set(adj, color)
				frontier = append(frontier, adj)
			default:
				// Dont' add the opponent's space to the frontier
			}
		}
	}

	// if we've exhausted the frontier, this is empty
	return mask
}

// Counts bounded pieces at p and clears them
func (b Board) clearBounded(start Position) int {
	mask := b.boundedMask(start)
	if mask == nil {
		return 0
	}
	count := 0
	sliced := b.slice()
	for i, state := range mask.slice() {
		if state != empty {
			sliced[i] = empty
			count++
		}
	}
	if count == 0 {
		panic("Mask was returned that was entirely empty")
	}
	return count
}

func (b Board) score() (blackPoints, whitePoints int) {
	points := b.copy()
	mask := newBoard(len(b))
	for x := 0; x < len(b); x++ {
		for y := 0; y < len(b); y++ {
			switch {
			case points[x][y] != empty:
			case mask[x][y] != empty:
			default:
				points.explore(Position{x, y}, mask)
			}
		}
	}
	for _, p := range points.slice() {
		switch p {
		case black:
			blackPoints++
		case white:
			whitePoints++
		}
	}
	return
}

func (b Board) explore(start Position, mask Board) {
	if b.get(start) != empty {
		return
	}

	explored := intersection(-1)
	previous := intersection(-2)
	unbounded := intersection(-3)

	// We're probably going to allocate somewhat initially, so lets allocate a bit
	frontier := make([]Position, 0, 64)
	frontier = append(frontier, start)

	color := empty
	mask.set(start, explored)
	// Walk the frontier
	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		// Check canditates up, down, left, right
		// Look for a connected "other". That means we're not bounded
		for _, adj := range current.adjacent() {
			switch {
			case !b.rangeCheck(adj):
				// Don't add out of range positions to the frontier
			case mask.get(adj) == explored:
				// Don't add previously checked positions to the frontier
			case mask.get(adj) == previous:
				panic("Unexpectedly got to previous!")

			case b.get(adj) == empty:
				mask.set(adj, explored)
				frontier = append(frontier, adj)

			// We've found a color, what do we do?
			case b.get(adj) == color:
			case color == empty:
				color = b.get(adj)
			default:
				color = unbounded
			}
		}
	}

	if color == white || color == black {
		for i, c := range mask.slice() {
			if c == explored {
				b.slice()[i] = color
			}
		}
	}
	// Clear found positions
	m := mask.slice()
	for i, c := range m {
		if c == explored {
			m[i] = previous
		}
	}
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
