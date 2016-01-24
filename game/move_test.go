package game

import (
	"testing"
)

const defaultBoardSize = 19

func TestMoveValid(t *testing.T) {
	tests := []struct {
		input    Move
		current  Color
		expected error
	}{
		{
			Move{Black, Position{0, 0}},
			Black,
			nil,
		},
		// Out of bounds
		{
			Move{Black, Position{20, 12}},
			Black,
			ErrOutOfBounds,
		},
		{
			Move{White, Position{-18, 18}},
			White,
			ErrOutOfBounds,
		},
		{
			Move{White, Position{12, 20}},
			White,
			ErrOutOfBounds,
		},
		{
			Move{Black, Position{12, -18}},
			Black,
			ErrOutOfBounds,
		},
		// Wrong color
		{
			Move{White, Position{12, 12}},
			Black,
			ErrWrongPlayer,
		},
		{
			Move{Black, Position{12, 12}},
			White,
			ErrWrongPlayer,
		},
	}

	for _, test := range tests {
		state := New(defaultBoardSize, 10)
		state.player = test.current
		err := state.Move(test.input)
		if err != test.expected {
			t.Errorf("for %v, expected '%s', got '%s'", test.input, test.expected, err)
		}
	}
}

func TestOpponent(t *testing.T) {
	if Black.Opponent() != White {
		t.Errorf("Black opponent should be White, was %s", Black.Opponent())
	}
	if White.Opponent() != Black {
		t.Errorf("White opponent should be Black, was %s", White.Opponent())
	}
}
