// Package game provides means of playing a game
package game

type stones struct {
	remaining int
	captured  int
}

// MoveError is returned if State.Move can't play the piece
type MoveError string

func (m MoveError) Error() string {
	return string(m)
}

// MoveErrors can occur if input is invalid, or if the player is unable to play
const (
	ErrWrongPlayer  = MoveError("Wrong player for move")
	ErrSpotNotEmpty = MoveError("Position filled")
	ErrOutOfBounds  = MoveError("Out of bounds")
	ErrNoStones     = MoveError("Player out of stones")
)

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
		return ErrWrongPlayer
	case s.stones[m.Player].remaining <= 0:
		return ErrNoStones
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
