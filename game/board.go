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

func (b Board) valid(m Move) error {
	if m.X >= len(b) ||
		m.X < 0 ||
		m.Y >= len(b) ||
		m.Y < 0 {
		return ErrOutOfBounds
	}
	return nil
}

func (b Board) apply(m Move) error {
	if err := b.valid(m); err != nil {
		return err
	}
	if err := b.intersectionEmpty(m.Position); err != nil {
		return err
	}

	stone := intersection(m.Player)
	b.set(m.Position, stone)

	captured := false
	for _, p := range m.adjacent() {
		switch {
		case !b.rangeCheck(p):
		case b.get(p) == stone:
		case b.clearBounded(p) > 0:
			captured = true
		}
	}

	if !captured && b.bounded(m.Position) {
		b.set(m.Position, empty)
		return fmt.Errorf("Bounded")
	}
	return nil
}

func (b Board) capture(p Position) bool {
	return false
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
