package game

type Position struct {
	X, Y int
}

// Move is used to provide an action
type Move struct {
	Player player
	Position
}

type player int

// There are only two players
const (
	Black = player(black)
	White = player(white)
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
