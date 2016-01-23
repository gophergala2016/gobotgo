package game

import (
	"testing"
)

func TestMoveValid(t *testing.T) {
	size := 19
	tests := []struct {
		input    Move
		current  player
		expected error
		reason   string
	}{
		{
			Move{Black, 0, 0},
			Black,
			nil,
			"Successful move",
		},
		{
			Move{Black, 20, 12},
			Black,
			MoveError("X coordinate 20 higher than size 19"),
			"X coordinate out of bounds positive",
		},
		{
			Move{Black, -18, 18},
			Black,
			MoveError("X coordinate -18 less than 0"),
			"X coordinate out of bounds negative",
		},
		{
			Move{Black, 12, 20},
			Black,
			MoveError("Y coordinate 20 higher than size 19"),
			"Y coordinate out of bounds positive",
		},
		{
			Move{Black, 12, -18},
			Black,
			MoveError("Y coordinate -18 less than 0"),
			"Y coordinate out of bounds negative",
		},
		{
			Move{White, 12, 12},
			Black,
			MoveError("Not your turn"),
			"Wrong player color moved",
		},
	}

	for _, test := range tests {
		err := test.input.valid(size, test.current)
		if err != test.expected {
			t.Error(test.reason, err)
		}
	}
}
