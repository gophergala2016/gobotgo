package game

import "testing"

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
				if bounded != b.bounded(x, y) {
					t.Errorf("bounded test for '%s' at %d,%d was unexpectedly %v", test.name, x, y, !bounded)
				}
			}
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
