package game

import (
	"testing"
)

const defaultBoardSize = 19

func TestMoveValid(t *testing.T) {
	tests := []struct {
		input    Move
		current  player
		expected error
		reason   string
	}{
		{
			Move{Black, Position{0, 0}},
			Black,
			nil,
			"Successful move",
		},
		{
			Move{Black, Position{20, 12}},
			Black,
			MoveError("X coordinate 20 higher than size 19"),
			"X coordinate out of bounds positive",
		},
		{
			Move{Black, Position{-18, 18}},
			Black,
			MoveError("X coordinate -18 less than 0"),
			"X coordinate out of bounds negative",
		},
		{
			Move{Black, Position{12, 20}},
			Black,
			MoveError("Y coordinate 20 higher than size 19"),
			"Y coordinate out of bounds positive",
		},
		{
			Move{Black, Position{12, -18}},
			Black,
			MoveError("Y coordinate -18 less than 0"),
			"Y coordinate out of bounds negative",
		},
		{
			Move{White, Position{12, 12}},
			Black,
			MoveError("Not your turn"),
			"Wrong player color moved",
		},
	}

	for _, test := range tests {
		state := New(defaultBoardSize)
		state.player = test.current
		err := state.Move(test.input)
		if err != test.expected {
			t.Error(test.reason, err)
		}
	}
}
