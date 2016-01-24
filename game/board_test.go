package game

import (
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	b := newBoard(5)
	// We use b[0]
	if cap(b[0]) != 25 {
		t.Errorf("Capacity needed to be 25")
	}
}

func TestCopy(t *testing.T) {
	moves := []Move{
		{White, Position{0, 0}},
		{White, Position{4, 7}},
		{Black, Position{6, 3}},
	}
	b := newBoard(8)
	for _, m := range moves {
		b[m.X][m.Y] = m.Player
	}
	c := b.copy()
	if cap(c[0]) != 64 {
		t.Errorf("Capacity needed to be 25")
	}
	for _, m := range moves {
		if c[m.X][m.Y] != m.Player {
			t.Errorf("Failed to copy %d,%d", m.X, m.Y)
		}
	}
	c[1][1] = Black
	b[2][2] = White
	if b[1][1] != empty {
		t.Error("Copied from copy to orignal")
	}
	if c[2][2] != empty {
		t.Error("Copied from orignal to copy")
	}
}

func TestBounded(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		b       []Color
		bounded []bool
	}{
		{
			"empty",
			3,
			[]Color{
				empty, empty, empty,
				empty, empty, empty,
				empty, empty, empty,
			},
			[]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
		},
		{
			"corner",
			3,
			[]Color{
				Black, White, empty,
				White, empty, empty,
				empty, empty, empty,
			},
			[]bool{
				true, false, false,
				false, false, false,
				false, false, false,
			},
		},
		{
			"diamond",
			3,
			[]Color{
				empty, Black, empty,
				Black, White, Black,
				empty, Black, empty,
			},
			[]bool{
				false, false, false,
				false, true, false,
				false, false, false,
			},
		},
		{
			"empty corner",
			3,
			[]Color{
				empty, White, White,
				White, White, White,
				White, White, White,
			},
			[]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
		},
		{
			"trapped corner",
			3,
			[]Color{
				Black, White, White,
				White, White, White,
				White, White, White,
			},
			[]bool{
				true, true, true,
				true, true, true,
				true, true, true,
			},
		},
		{
			"filled box in space",
			5,
			[]Color{
				empty, empty, empty, empty, empty,
				empty, empty, Black, Black, empty,
				empty, Black, White, White, Black,
				empty, empty, Black, Black, empty,
				empty, empty, empty, empty, empty,
			},
			[]bool{
				false, false, false, false, false,
				false, false, false, false, false,
				false, false, true, true, false,
				false, false, false, false, false,
				false, false, false, false, false,
			},
		},
		{
			"surrounded eye safe",
			5,
			[]Color{
				White, White, White, White, White,
				White, Black, Black, Black, White,
				White, Black, empty, Black, White,
				White, Black, Black, Black, White,
				White, White, White, White, White,
			},
			[]bool{
				true, true, true, true, true,
				true, false, false, false, true,
				true, false, false, false, true,
				true, false, false, false, true,
				true, true, true, true, true,
			},
		},
		{
			"surrounded eye fully bounded",
			5,
			[]Color{
				White, White, White, White, White,
				White, Black, Black, Black, White,
				White, Black, White, Black, White,
				White, Black, Black, Black, White,
				White, White, White, White, White,
			},
			[]bool{
				true, true, true, true, true,
				true, true, true, true, true,
				true, true, true, true, true,
				true, true, true, true, true,
				true, true, true, true, true,
			},
		},
		{
			"surrounded eye not bounded",
			7,
			[]Color{
				empty, empty, empty, empty, empty, empty, empty,
				empty, White, White, White, White, White, empty,
				empty, White, Black, Black, Black, White, empty,
				empty, White, Black, White, Black, White, empty,
				empty, White, Black, Black, Black, White, empty,
				empty, White, White, White, White, White, empty,
				empty, empty, empty, empty, empty, empty, empty,
			},
			[]bool{
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
				false, false, true, true, true, false, false,
				false, false, true, true, true, false, false,
				false, false, true, true, true, false, false,
				false, false, false, false, false, false, false,
				false, false, false, false, false, false, false,
			},
		},
	}
	for _, test := range tests {
		b := sliceBoard(test.b, test.size)
		for x, row := range b {
			for y := range row {
				bounded := test.bounded[x*test.size+y]
				if bounded != b.bounded(Position{x, y}) {
					t.Errorf("bounded test for '%s' at %d,%d was unexpectedly %v", test.name, x, y, !bounded)
				}
			}
		}
	}
}

func TestBoundedMask(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		board []Color
		p     Position
		mask  []Color // nil if move is invalid
	}{
		{
			"empty", 2,
			[]Color{
				empty, empty,
				empty, empty,
			},
			Position{0, 0},
			nil,
		},
		{
			"open corner", 2,
			[]Color{
				empty, Black,
				Black, empty,
			},
			Position{1, 0},
			nil,
		},
		{
			"closed corner", 2,
			[]Color{
				White, Black,
				empty, White,
			},
			Position{0, 1},
			[]Color{
				empty, Black,
				empty, empty,
			},
		},
		{
			"shaped bound", 4,
			[]Color{
				empty, White, White, empty,
				White, Black, Black, White,
				empty, White, Black, White,
				Black, Black, White, empty,
			},
			Position{2, 2},
			[]Color{
				empty, empty, empty, empty,
				empty, Black, Black, empty,
				empty, empty, Black, empty,
				empty, empty, empty, empty,
			},
		},
	}

	for _, test := range tests {
		b := sliceBoard(test.board, test.size)
		mask := b.boundedMask(test.p)
		if mask == nil {
			if test.mask != nil {
				t.Errorf("mask '%s' unexpectedly empty", test.name)
			}
		} else {
			if test.mask == nil {
				t.Errorf("mask '%s' unexpectedly exists", test.name)
			} else if err := mask.equal(sliceBoard(test.mask, test.size)); err != nil {
				t.Errorf("mask '%s' result not equal to expected: %s", test.name, err.Error())
			}
		}
	}
}

func coalesce(lhs, rhs []Color) []Color {
	if lhs != nil {
		return lhs
	}
	return rhs
}

func TestClearBounded(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		initial []Color
		p       Position
		final   []Color // nil if unchanged
		removed int
	}{
		{
			"empty", 2,
			[]Color{
				empty, empty,
				empty, empty,
			},
			Position{0, 0},
			nil,
			0,
		},
		{
			"open corner", 2,
			[]Color{
				empty, Black,
				Black, empty,
			},
			Position{1, 0},
			nil,
			0,
		},
		{
			"closed corner", 2,
			[]Color{
				White, Black,
				empty, White,
			},
			Position{0, 1},
			[]Color{
				White, empty,
				empty, White,
			},
			1,
		},
		{
			"shaped bound", 4,
			[]Color{
				empty, White, White, empty,
				White, Black, Black, White,
				empty, White, Black, White,
				Black, Black, White, empty,
			},
			Position{2, 2},
			[]Color{
				empty, White, White, empty,
				White, empty, empty, White,
				empty, White, empty, White,
				Black, Black, White, empty,
			},
			3,
		},
	}

	for _, test := range tests {
		before := sliceBoard(test.initial, test.size)
		after := sliceBoard(coalesce(test.final, test.initial), test.size)
		b := before.copy()
		removed := b.clearBounded(test.p)
		if removed != test.removed {
			t.Errorf("clear '%s' expected to remove %d pieces, removed %d pieces", test.name, test.removed, removed)
		}
		if err := b.equal(after); err != nil {
			t.Errorf("clear '%s' result not equal to expected: %s", test.name, err.Error())
		}
	}
}

func TestBoardEqual(t *testing.T) {
	size := 19

	a := newBoard(size)
	a[5][5] = Black
	a[6][6] = White
	b := a.copy()

	err := a.equal(b)

	if err != nil {
		t.Error("Boards not equivalent")
	}
}

func TestBoardNotEqual(t *testing.T) {
	size := 19

	a := newBoard(size)
	a[5][5] = Black
	a[6][6] = White
	b := newBoard(size)
	b[5][5] = Black
	b[6][6] = Black

	err := a.equal(b)

	if err == nil {
		t.Error("Boards were equal")
	}
}

func TestApplyMove(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		initial []Color
		m       Move
		final   []Color // nil if move is invalid
		points  int
	}{
		{
			"empty", 2,
			[]Color{
				empty, empty,
				empty, empty,
			},
			Move{Black, Position{0, 0}},
			[]Color{
				Black, empty,
				empty, empty,
			},
			0,
		},
		{
			"capture corner", 2,
			[]Color{
				empty, Black,
				empty, White,
			},
			Move{White, Position{0, 0}},
			[]Color{
				White, empty,
				empty, White,
			},
			1,
		},
		{
			"corner eye, player", 3,
			[]Color{
				empty, Black, empty,
				Black, empty, empty,
				empty, empty, empty,
			},
			Move{Black, Position{0, 0}},
			[]Color{
				Black, Black, empty,
				Black, empty, empty,
				empty, empty, empty,
			},
			0,
		},
		{
			"power corner", 2,
			[]Color{
				empty, Black,
				Black, Black,
			},
			Move{White, Position{0, 0}},
			[]Color{
				White, empty,
				empty, empty,
			},
			3,
		},
		// Self capture tests
		{
			"simple corner", 2,
			[]Color{
				empty, Black,
				Black, empty,
			},
			Move{White, Position{0, 0}},
			nil,
			0,
		},
		{
			"corner eye, opponent", 3,
			[]Color{
				empty, Black, empty,
				Black, empty, empty,
				empty, empty, empty,
			},
			Move{White, Position{0, 0}},
			nil,
			0,
		},
		{
			"overlapping capture", 4,
			[]Color{
				empty, Black, White, empty,
				Black, empty, Black, White,
				empty, Black, White, empty,
				empty, empty, empty, empty,
			},
			Move{White, Position{1, 1}},
			[]Color{
				empty, Black, White, empty,
				Black, White, empty, White,
				empty, Black, White, empty,
				empty, empty, empty, empty,
			},
			1,
		},
	}

	for _, test := range tests {
		before := sliceBoard(test.initial, test.size)
		after := sliceBoard(coalesce(test.final, test.initial), test.size)
		b := before.copy()
		if nil != b.intersectionEmpty(test.m.Position) {
			t.Errorf("Test %s tried to move to non-empty space %d, %d", test.name, test.m.X, test.m.Y)
			continue
		}
		points, err := b.Apply(test.m)
		if (test.final == nil) == (err == nil) {
			t.Errorf("movability of '%s' was unexpectedly %v: %s", test.name, (test.final == nil), err.Error())
		}
		if points != test.points {
			t.Errorf("expected %d points, got %d", test.points, points)
		}
		if err := b.equal(after); err != nil {
			t.Errorf("Move '%s' result not equal to expected: %s", test.name, err.Error())
		}
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		size         int
		board        []Color
		black, white int
	}{
		{
			2,
			[]Color{
				empty, empty,
				empty, empty,
			},
			0, 0,
		},
		{
			2,
			[]Color{
				empty, Black,
				Black, empty,
			},
			4, 0,
		},
		{
			2,
			[]Color{
				White, Black,
				empty, White,
			},
			1, 3,
		},
		{
			4,
			[]Color{
				empty, White, White, empty,
				White, Black, Black, White,
				empty, White, Black, White,
				Black, Black, White, empty,
			},
			5, 10,
		},
		{
			4,
			[]Color{
				empty, White, White, empty,
				White, empty, empty, White,
				empty, White, empty, White,
				Black, Black, White, empty,
			},
			2, 13,
		},
		{
			5,
			[]Color{
				empty, White, empty, White, empty,
				Black, empty, empty, empty, White,
				empty, Black, empty, empty, White,
				Black, empty, empty, White, empty,
				empty, Black, empty, White, empty,
			},
			6, 9,
		},
	}

	for i, test := range tests {
		board := sliceBoard(test.board, test.size)
		b, w := board.Score()
		if b != test.black || w != test.white {
			t.Errorf("for %d expected %d-%d, got %d-%d\n%s", i, test.black, test.white, b, w, board)
		}
	}
}

func ExampleColor_Dot() {
	fmt.Println(None.Dot(), White.Dot(), Black.Dot(), empty.Dot())
	// Output: . w b .
}

func ExampleBoard_String() {
	i := []Color{
		White, Black, empty,
		empty, White, Black,
		Black, White, empty,
	}
	fmt.Println(sliceBoard(i, 3))
	// Output:
	// w b .
	// . w b
	// b w .
}
