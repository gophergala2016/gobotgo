package game

import "testing"

func TestTurnOrder(t *testing.T) {
	s := New(4, 20)
	switch {
	case s.player != Black:
		t.Error("Expected first play to be Black")
	case s.Move(Move{White, Position{0, 0}}) != ErrWrongPlayer:
		t.Error("White should not be able to move")
	case s.player != Black:
		t.Error("Player should still be Black")
	case s.Move(Move{Black, Position{0, 0}}) != nil:
		t.Error("Black should have moved")
	case s.player != White:
		t.Error("Player should now be White")
	case s.Move(Move{Black, Position{1, 1}}) != ErrWrongPlayer:
		t.Error("Black should not be able to move")
	case s.player != White:
		t.Error("Player should still be White")
	case s.Move(Move{White, Position{1, 1}}) != nil:
		t.Error("White should have moved")
	case s.player != Black:
		t.Error("Player should now be Black")
	}
}

func TestRemainingCountdown(t *testing.T) {
	tests := []struct {
		Move
		black, white int
	}{
		{Move{Black, Position{1, 1}}, 1, 2},
		{Move{White, Position{1, 0}}, 1, 1},
		{Move{Black, Position{2, 1}}, 0, 1},
		{Move{White, Position{2, 0}}, 0, 0},
	}
	s := New(4, 2)
	for i, test := range tests {
		if err := s.Move(test.Move); err != nil {
			t.Fatalf("Unexpected error for move %d:%+v:%s", i, test.Move, err.Error())
		}
		if test.black != s.stones[Black].remaining ||
			test.white != s.stones[White].remaining {
			t.Fatalf("Failed remaing for row %d. Expected %d,%d got %d,%d.", i, test.black, test.white, s.stones[Black].remaining, s.stones[White].remaining)
		}
	}
}

func TestRemainingStop(t *testing.T) {
	tests := []Move{
		{Black, Position{1, 1}},
		{White, Position{1, 2}},
	}
	s := New(4, 0)
	for i, test := range tests {
		s.player = test.Player
		if err := s.Move(test); err != ErrNoStones {
			t.Errorf("Expected error '%s' for move %d:%+v, got '%s'", ErrNoStones, i, test, err.Error())
		}
		if s.stones[Black].remaining != 0 || s.stones[White].remaining != 0 {
			t.Errorf("Remaining unexpectedly not zero: %d,%d", s.stones[Black].remaining, s.stones[White].remaining)
		}
	}
}

func TestCannotLoopGameState(t *testing.T) {
	tests := []struct {
		Move
		err error
	}{
		{Move{Black, Position{1, 2}}, nil},
		{Move{White, Position{1, 1}}, nil},
		{Move{Black, Position{1, 2}}, ErrRepeatState},
		{Move{White, Position{1, 2}}, ErrWrongPlayer},
		{Move{Black, Position{0, 0}}, nil},
	}
	size := 4
	s := New(size, 20)
	s.current = sliceBoard([]intersection{
		empty, black, white, empty,
		black, empty, empty, white,
		empty, black, white, empty,
		empty, empty, empty, empty,
	}, size)
	for i, test := range tests {
		if err := s.Move(test.Move); err != test.err {
			t.Errorf("for move %d expected error '%s' got '%s'", i, test.err, err)
		}
	}
}
