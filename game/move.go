// Move is used to provide an action
package game

import "fmt"

type Position struct {
	X, Y int
}

type Move struct {
	Player player
	Position
}

type MoveError string

type player int

const (
	Black = player(black)
	White = player(white)
)

func (m MoveError) Error() string {
	return string(m)
}

func (m Move) valid(max int, current player) error {
	switch {
	case m.X >= max:
		return MoveError(fmt.Sprintf("X coordinate %d higher than size %d", m.X, max))
	case m.X < 0:
		return MoveError(fmt.Sprintf("X coordinate %d less than 0", m.X))
	case m.Y >= max:
		return MoveError(fmt.Sprintf("Y coordinate %d higher than size %d", m.Y, max))
	case m.Y < 0:
		return MoveError(fmt.Sprintf("Y coordinate %d less than 0", m.Y))
	case m.Player != current:
		return MoveError(fmt.Sprintf("Not your turn"))
	default:
		return nil
	}
}
