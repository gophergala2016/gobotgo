package game

import (
	"encoding/json"
	"strings"
)

type Position struct {
	X, Y int
}

// Move is used to provide an action
type Move struct {
	Player Color
	Position
}

type Color int

// There are only two players
const (
	None = Color(iota)
	Black
	White
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

func (c Color) Opponent() Color {
	switch c {
	case White:
		return Black
	case Black:
		return White
	default:
		return None
	}
}

func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "None"
	}
}

func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Color) UnmarshalJSON(data []byte) error {
	s := strings.ToLower(string(data))
	switch {
	case s == `"black"`:
		*c = Black
	case s == `"white"`:
		*c = White
	default:
		*c = None
	}
	return nil
}
