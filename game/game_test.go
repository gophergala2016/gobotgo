package game

import "testing"

func TestTurnOrder(t *testing.T) {
	s := New(4)
	switch {
	case s.player != Black:
		t.Error("Expected first play to be Black")
	case s.Move(Move{White, Position{0, 0}}) == nil:
		t.Error("White should not be able to move")
	case s.player != Black:
		t.Error("Player should still be Black")
	case s.Move(Move{Black, Position{0, 0}}) != nil:
		t.Error("Black should have moved")
	case s.player != White:
		t.Error("Player should now be White")
	case s.Move(Move{Black, Position{1, 1}}) == nil:
		t.Error("Black should not be able to move")
	case s.player != White:
		t.Error("Player should still be White")
	case s.Move(Move{White, Position{1, 1}}) != nil:
		t.Error("White should have moved")
	case s.player != Black:
		t.Error("Player should now be Black")
	}
}
