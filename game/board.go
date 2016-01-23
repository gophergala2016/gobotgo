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

func (b Board) equal(other Board) error {
	for i := range b {
		if b[i] != other[i] {
			return Errorf("Board state not equal at %d", i)
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

func (b Board) copy() Board {
	l := len(b)
	a := make([]intersection, l*l)
	copy(a, b[0][:l*l])
	return sliceBoard(a, l)
}
