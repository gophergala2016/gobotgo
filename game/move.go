// Move is used to provide an action
package game

import "fmt"

type Move struct {
	Player player
	X, Y   int
}

type player int

const (
	Black = player(black)
	White = player(white)
)

func (m Move) valid(max int) error {
	switch {
	case m.X >= max:
		return fmt.Errorf("X cooridant %d higher than size %d", m.X, max)
	case m.X < 0:
		return fmt.Errorf("X cooridant %d less than 0", m.X)
	case m.Y >= max:
		return fmt.Errorf("Y cooridant %d higher than size %d", m.Y, max)
	case m.Y < 0:
		return fmt.Errorf("Y cooridant %d less than 0", m.Y)
	default:
		return nil
	}
}
