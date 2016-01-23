// Package game provides means of playing a game
package game

import "fmt"

type State struct {
	current  Board
	previous Board
	player   player
	size     int
}

func New(size int) State {
	c := newBoard(size)
	p := newBoard(size)
	return State{
		current:  c,
		previous: p,
		player:   Black,
		size:     size,
	}
}

func (s State) valid(m Move) error {
	if m.Player != s.player {
		return MoveError(fmt.Sprintf("Not your turn"))
	}
	return nil
}

func (s State) Move(m Move) error {
	if err := s.valid(m); err != nil {
		return err
	}
	b := s.current.copy()
	if err := b.apply(m); err != nil {
		return err
	}
	s.previous = s.current
	s.current = b
	return nil
}
