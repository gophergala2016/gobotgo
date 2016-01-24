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
	s.current = sliceBoard([]Color{
		empty, Black, White, empty,
		Black, empty, empty, White,
		empty, Black, White, empty,
		empty, empty, empty, empty,
	}, size)
	for i, test := range tests {
		if err := s.Move(test.Move); err != test.err {
			t.Errorf("for move %d expected error '%s' got '%s'", i, test.err, err)
		}
	}
}

func TestDoublePassEndsGame(t *testing.T) {
	size := 4
	s := New(size, 30)
	s.current = sliceBoard([]Color{
		empty, Black, White, empty,
		Black, empty, empty, White,
		empty, Black, White, empty,
		empty, empty, empty, empty,
	}, size)
	moves := []Move{
		{Black, Position{1, 2}},
		{White, Position{1, 1}},
		{Black, Position{3, 3}},
		{White, Position{2, 3}},
		{Black, Position{1, 2}},
	}
	for i, m := range moves {
		if err := s.Move(m); err != nil {
			t.Fatalf("did not expect move #%d:%v to end in '%s'", i, m, err)
		}
	}
	if err := s.Pass(Black); err != ErrWrongPlayer {
		t.Fatalf("Expected '%s' when black tried to pass, got '%s'", ErrWrongPlayer, err)
	}
	if err := s.Pass(White); err != nil {
		t.Fatalf("Expected white to be able to pass, got '%s'", err)
	}
	if err := s.Move(Move{White, Position{0, 0}}); err != ErrWrongPlayer {
		t.Fatalf("Expected '%s' when white went after passing, got '%s'", ErrWrongPlayer, err)
	}
	if err := s.Move(Move{Black, Position{0, 0}}); err != nil {
		t.Fatalf("Expected black to go after white passed, got '%s'", err)
	}
	passes := []struct {
		c   Color
		err error
	}{
		{White, nil},
		{White, ErrWrongPlayer},
		{Black, ErrGameOver},
		{Black, ErrGameOver},
		{White, ErrGameOver},
		{White, ErrGameOver},
	}
	for i, pass := range passes {
		if err := s.Pass(pass.c); err != pass.err {
			t.Fatalf("Expected Pass(%s) at %d to be '%s', got '%s'", pass.c, i, pass.err, err)
		}
	}
	moves = []Move{
		{Black, Position{0, 4}},
		{White, Position{0, 4}},
		{Black, Position{1, 1}},
		{White, Position{1, 1}},
	}
	for i, m := range moves {
		if err := s.Move(m); err != ErrGameOver {
			t.Fatalf("Expected Move(%v) at %d to be '%s', got '%s'", m, i, ErrGameOver, err)
		}
	}
}

func TestPassBlackDoesNotEntGame(t *testing.T) {
	s := New(5, 20)
	if err := s.Pass(Black); err != nil {
		t.Fatalf("Did not expect to be unable to pass black, '%s'", err)
	}
}

func TestOutOfStonesEnds(t *testing.T) {
	stoneCount := 4
	size := 4
	s := New(size, stoneCount)
	for i := 0; i < stoneCount; i++ {
		if err := s.Pass(Black); err != nil {
			t.Errorf("Unexepected error for %d:Pass(Black), got '%s'", i, err.Error())
		}
		m := Move{White, Position{i, 0}}
		if err := s.Move(m); err != nil {
			t.Errorf("Unexepected error for move %d:%v, got '%s'", i, m, err.Error())
		}
	}
	for i := 0; i < stoneCount-1; i++ {
		m := Move{Black, Position{i, 1}}
		if err := s.Move(m); err != nil {
			t.Errorf("Unexepected for move %d:%v, got '%s'", i, m, err.Error())
		}
		m = Move{White, Position{i, 2}}
		if err := s.Move(m); err != ErrNoStones {
			t.Errorf("Expected %d:%v to be '%s', got '%s'", i, m, ErrNoStones, err.Error())
		}
		if err := s.Pass(White); err != nil {
			t.Errorf("Unexepected error for %d:Pass(White), got '%s'", i, err.Error())
		}
	}
	m := Move{Black, Position{3, 1}}
	if err := s.Move(m); err != nil {
		t.Errorf("Unexepected for move %v, got '%s'", m, err.Error())
	}
	tests := []struct {
		Move
		err error
	}{
		{Move{White, Position{3, 3}}, ErrGameOver},
		{Move{Black, Position{3, 3}}, ErrGameOver},
	}

	for i, test := range tests {
		if err := s.Move(test.Move); err != test.err {
			t.Fatalf("Expected %d:Move(%v) to be '%s', got '%s'", i, test.Move, test.err, err)
		}
	}
}

func TestScoring(t *testing.T) {
	score := struct {
		black, white int
	}{5, 2}

	// assume alternating moves
	moves := []Position{
		// black, white
		{1, 1}, {0, 1},
		{0, 2}, {0, 3},
		{0, 0}, {0, 4},
	}
	p := Black
	s := New(5, 100)
	for i, m := range moves {
		move := Move{p, m}
		if err := s.Move(move); err != nil {
			t.Fatalf("failed move %d:%v, got '%s'", i, move, err.Error())
		}
		p = p.opponent()
	}

	t.Log(s.stones[Black], s.stones[White])
	b, w := s.Score()
	if score.black != b || score.white != w {
		t.Errorf("unmatched score, expected %d-%d, got %d-%d", score.black, score.white, b, w)
	}
}
