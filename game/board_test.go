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
		{White, 0, 0},
		{White, 4, 7},
		{Black, 6, 3},
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
