// Package game provides means of playing a game
package game

import (
	"encoding/json"
	"errors"
)

type stones struct {
	Remaining int `json:"remaining"`
	Captured  int `json:"captured"`
}

// MoveError is returned if State.Move can't play the piece
type MoveError string

func (m MoveError) Error() string {
	return string(m)
}

// MoveErrors can occur if input is invalid, or if the player is unable to play
const (
	ErrGameOver     = MoveError("Game Over")
	ErrWrongPlayer  = MoveError("Wrong player for move")
	ErrSpotNotEmpty = MoveError("Position filled")
	ErrOutOfBounds  = MoveError("Out of bounds")
	ErrNoStones     = MoveError("Player out of stones")
	ErrRepeatState  = MoveError("Move recreates previous state")
)

type LastMove struct {
	Move
	PiecesRemoved int
}

type State struct {
	current  Board
	previous Board
	player   Color
	over     bool
	size     int
	pieces   int
	stones   map[Color]*stones
	last     LastMove
}

func New(size, pieces int) *State {
	c := newBoard(size)
	return &State{
		current:  c,
		previous: nil,
		player:   Black,
		over:     false,
		size:     size,
		pieces:   pieces,
		stones: map[Color]*stones{
			White: {pieces, 0},
			Black: {pieces, 0},
		},
	}
}

func (s *State) valid(m Move) error {
	switch {
	case s.over:
		return ErrGameOver
	case m.Player != s.player:
		return ErrWrongPlayer
	case s.stones[m.Player].Remaining <= 0:
		if s.stones[m.Player.Opponent()].Remaining <= 0 {
			s.over = true
			return ErrGameOver
		}
		return ErrNoStones
	}
	return nil
}

func (s *State) Pass(player Color) error {
	if s.over {
		return ErrGameOver
	}
	// Legal to pass when out of stones
	if player != s.player {
		return ErrWrongPlayer
	}
	if s.current.equal(s.previous) == nil {
		s.over = true
		return ErrGameOver
	}
	s.previous = s.current
	s.player = player.Opponent()
	return nil
}

func (s *State) Move(m Move) error {
	if err := s.valid(m); err != nil {
		return err
	}
	b := s.current.copy()
	captured, err := b.apply(m)
	if err != nil {
		return err
	}
	if b.equal(s.previous) == nil {
		return ErrRepeatState
	}
	s.previous = s.current
	s.current = b
	s.stones[m.Player].Remaining--
	s.stones[m.Player].Captured += captured
	s.player = m.Player.Opponent()
	s.last = LastMove{m, captured}
	return nil
}

func (s *State) Score() (black, white int) {
	black, white = s.current.score()
	black += s.stones[Black].Captured
	white += s.stones[White].Captured
	return
}

func (s *State) MarshalJSON() ([]byte, error) {
	data := struct {
		Board         Board    `json:"board"`
		CurrentPlayer Color    `json:"currentplayer"`
		Black         stones   `json:"black"`
		White         stones   `json:"white"`
		LastMove      LastMove `json:"lastmove,omitempty"`
	}{
		s.current,
		s.player,
		*s.stones[Black],
		*s.stones[White],
		s.last,
	}
	return json.Marshal(data)
}

func (s *State) UnmarshalJSON(data []byte) error {
	return errors.New("cannot unmarshal state as data is lost in JSON transmission")
}
