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
