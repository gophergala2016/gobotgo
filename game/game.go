// Package game provides means of playing a game
package game

import "fmt"

type stones struct {
	remaining int
	captured  int
}

type State struct {
	current  Board
	previous Board
	player   color
	size     int
	pieces   int
	stones   map[color]*stones
}

func New(size, pieces int) *State {
	c := newBoard(size)
	p := newBoard(size)
	return &State{
		current:  c,
		previous: p,
		player:   Black,
		size:     size,
		pieces:   pieces,
		stones: map[color]*stones{
			White: {pieces, 0},
			Black: {pieces, 0},
		},
	}
}

func (s *State) valid(m Move) error {
	switch {
	case m.Player != s.player:
		return MoveError(fmt.Sprintf("Not your turn"))
	case s.stones[m.Player].remaining <= 0:
		return MoveError(fmt.Sprintf("Out of pieces"))
	}
	return nil
}

func (s *State) Move(m Move) error {
	if err := s.valid(m); err != nil {
		return err
	}
	b := s.current.copy()
	if err := b.apply(m); err != nil {
		return err
	}
	s.previous = s.current
	s.current = b
	s.stones[m.Player].remaining--
	s.player = m.Player.opponent()
	return nil
}
