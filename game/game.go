// Package game provides means of playing a game
package game

type intersection int
type player int

// Player types
const (
	empty = intersection(iota)
	black
	white
)

const (
	Black = player(black)
	White = player(white)
)

// Move is used to provide an action
type Move struct {
	Player player
	X, Y   int
}

type Board [][]intersection

type State struct {
	current  Board
	previous Board
	player   player
	size     int
}

func New(size int) State {
	c := make(Board, size)
	p := make(Board, size)
	return State{
		current:  c,
		previous: p,
		player:   White,
		size:     size,
	}
}
