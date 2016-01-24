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
		b[m.X][m.Y] = intersection(m.Player)
	}
	c := b.copy()
	if cap(c[0]) != 64 {
		t.Errorf("Capacity needed to be 25")
	}
	for _, m := range moves {
		if c[m.X][m.Y] != intersection(m.Player) {
			t.Errorf("Failed to copy %d,%d", m.X, m.Y)
		}
	}
	c[1][1] = intersection(Black)
	b[2][2] = intersection(White)
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
		b       []intersection
		bounded []bool
	}{
		{
			"empty",
			3,
			[]intersection{
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
			[]intersection{
				black, white, empty,
				white, empty, empty,
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
			[]intersection{
				empty, black, empty,
				black, white, black,
				empty, black, empty,
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
			[]intersection{
				empty, white, white,
				white, white, white,
				white, white, white,
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
			[]intersection{
				black, white, white,
				white, white, white,
				white, white, white,
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
			[]intersection{
				empty, empty, empty, empty, empty,
				empty, empty, black, black, empty,
				empty, black, white, white, black,
				empty, empty, black, black, empty,
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
			[]intersection{
				white, white, white, white, white,
				white, black, black, black, white,
				white, black, empty, black, white,
				white, black, black, black, white,
				white, white, white, white, white,
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
			[]intersection{
				white, white, white, white, white,
				white, black, black, black, white,
				white, black, white, black, white,
				white, black, black, black, white,
				white, white, white, white, white,
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
			[]intersection{
				empty, empty, empty, empty, empty, empty, empty,
				empty, white, white, white, white, white, empty,
				empty, white, black, black, black, white, empty,
				empty, white, black, white, black, white, empty,
				empty, white, black, black, black, white, empty,
				empty, white, white, white, white, white, empty,
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
		board []intersection
		p     Position
		mask  []intersection // nil if move is invalid
	}{
		{
			"empty", 2,
			[]intersection{
				empty, empty,
				empty, empty,
			},
			Position{0, 0},
			nil,
		},
		{
			"open corner", 2,
			[]intersection{
				empty, black,
				black, empty,
			},
			Position{1, 0},
			nil,
		},
		{
			"closed corner", 2,
			[]intersection{
				white, black,
				empty, white,
			},
			Position{0, 1},
			[]intersection{
				empty, black,
				empty, empty,
			},
		},
		{
			"shaped bound", 4,
			[]intersection{
				empty, white, white, empty,
				white, black, black, white,
				empty, white, black, white,
				black, black, white, empty,
			},
			Position{2, 2},
			[]intersection{
				empty, empty, empty, empty,
				empty, black, black, empty,
				empty, empty, black, empty,
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

func coalesce(lhs, rhs []intersection) []intersection {
	if lhs != nil {
		return lhs
	}
	return rhs
}

func TestClearBounded(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		initial []intersection
		p       Position
		final   []intersection // nil if unchanged
		removed int
	}{
		{
			"empty", 2,
			[]intersection{
				empty, empty,
				empty, empty,
			},
			Position{0, 0},
			nil,
			0,
		},
		{
			"open corner", 2,
			[]intersection{
				empty, black,
				black, empty,
			},
			Position{1, 0},
			nil,
			0,
		},
		{
			"closed corner", 2,
			[]intersection{
				white, black,
				empty, white,
			},
			Position{0, 1},
			[]intersection{
				white, empty,
				empty, white,
			},
			1,
		},
		{
			"shaped bound", 4,
			[]intersection{
				empty, white, white, empty,
				white, black, black, white,
				empty, white, black, white,
				black, black, white, empty,
			},
			Position{2, 2},
			[]intersection{
				empty, white, white, empty,
				white, empty, empty, white,
				empty, white, empty, white,
				black, black, white, empty,
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
	a[5][5] = intersection(Black)
	a[6][6] = intersection(White)
	b := a.copy()

	err := a.equal(b)

	if err != nil {
		t.Error("Boards not equivalent")
	}
}

func TestBoardNotEqual(t *testing.T) {
	size := 19

	a := newBoard(size)
	a[5][5] = intersection(Black)
	a[6][6] = intersection(White)
	b := newBoard(size)
	b[5][5] = intersection(Black)
	b[6][6] = intersection(Black)

	err := a.equal(b)

	if err == nil {
		t.Error("Boards were equal")
	}
}

func TestApplyMove(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		initial []intersection
		m       Move
		final   []intersection // nil if move is invalid
	}{
		{
			"empty", 2,
			[]intersection{
				empty, empty,
				empty, empty,
			},
			Move{Black, Position{0, 0}},
			[]intersection{
				black, empty,
				empty, empty,
			},
		},
		{
			"capture corner", 2,
			[]intersection{
				empty, black,
				empty, white,
			},
			Move{White, Position{0, 0}},
			[]intersection{
				white, empty,
				empty, white,
			},
		},
		{
			"corner eye, player", 3,
			[]intersection{
				empty, black, empty,
				black, empty, empty,
				empty, empty, empty,
			},
			Move{Black, Position{0, 0}},
			[]intersection{
				black, black, empty,
				black, empty, empty,
				empty, empty, empty,
			},
		},
		{
			"power corner", 2,
			[]intersection{
				empty, black,
				black, black,
			},
			Move{White, Position{0, 0}},
			[]intersection{
				white, empty,
				empty, empty,
			},
		},
		// Self capture tests
		{
			"simple corner", 2,
			[]intersection{
				empty, black,
				black, empty,
			},
			Move{White, Position{0, 0}},
			nil,
		},
		{
			"corner eye, opponent", 3,
			[]intersection{
				empty, black, empty,
				black, empty, empty,
				empty, empty, empty,
			},
			Move{White, Position{0, 0}},
			nil,
		},
		{
			"overlapping capture", 4,
			[]intersection{
				empty, black, white, empty,
				black, empty, black, white,
				empty, black, white, empty,
				empty, empty, empty, empty,
			},
			Move{White, Position{1, 1}},
			[]intersection{
				empty, black, white, empty,
				black, white, empty, white,
				empty, black, white, empty,
				empty, empty, empty, empty,
			},
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
		if err := b.apply(test.m); (test.final == nil) == (err == nil) {
			t.Errorf("movability of '%s' was unexpectedly %v: %s", test.name, (test.final == nil), err.Error())
		}
		if err := b.equal(after); err != nil {
			t.Errorf("Move '%s' result not equal to expected: %s", test.name, err.Error())
		}
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		size         int
		board        []intersection
		black, white int
	}{
		{
			2,
			[]intersection{
				empty, empty,
				empty, empty,
			},
			0, 0,
		},
		{
			2,
			[]intersection{
				empty, black,
				black, empty,
			},
			4, 0,
		},
		{
			2,
			[]intersection{
				white, black,
				empty, white,
			},
			1, 3,
		},
		{
			4,
			[]intersection{
				empty, white, white, empty,
				white, black, black, white,
				empty, white, black, white,
				black, black, white, empty,
			},
			5, 10,
		},
		{
			4,
			[]intersection{
				empty, white, white, empty,
				white, empty, empty, white,
				empty, white, empty, white,
				black, black, white, empty,
			},
			2, 13,
		},
		{
			5,
			[]intersection{
				empty, white, empty, white, empty,
				black, empty, empty, empty, white,
				empty, black, empty, empty, white,
				black, empty, empty, white, empty,
				empty, black, empty, white, empty,
			},
			6, 9,
		},
	}

	for i, test := range tests {
		board := sliceBoard(test.board, test.size)
		b, w := board.score()
		if b != test.black || w != test.white {
			t.Errorf("for %d expected %d-%d, got %d-%d\n%s", i, test.black, test.white, b, w, board)
		}
	}
}

func ExampleIntersection_String() {
	fmt.Println(white, black, empty)
	// Output: w b .
}

func ExampleBoard_String() {
	i := []intersection{
		white, black, empty,
		empty, white, black,
		black, white, empty,
	}
	fmt.Println(sliceBoard(i, 3))
	// Output:
	// w b .
	// . w b
	// b w .
}
