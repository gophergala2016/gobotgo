package game

type Position struct {
	X, Y int
}

// Move is used to provide an action
type Move struct {
	Player color
	Position
}

type color int

// There are only two players
const (
	Black = color(black)
	White = color(white)
)

func (p Position) add(q Position) Position {
	return Position{
		p.X + q.X,
		p.Y + q.Y,
	}
}

func (p Position) adjacent() [4]Position {
	var adj [4]Position
	adjacentMoves := [4]Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for i, m := range adjacentMoves {
		adj[i] = p.add(m)
	}
	return adj
}

func (c color) opponent() color {
	if c == White {
		return Black
	}
	return White
}

func (c color) String() string {
	switch c {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "(invalid color)"
	}
}
