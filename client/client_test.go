package client

import (
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/gophergala2016/gobotgo/game"
	"github.com/gophergala2016/gobotgo/server"
)

func newClient(t *testing.T, URL string, ID server.GameID, color game.Color) *Client {
	c, err := New(URL)
	if err != nil {
		t.Fatal("failed to initialize client: '%s'", err)
	}
	if c.ID() != ID {
		t.Fatalf("expected ID %d, got %d", ID, c.ID())
	}
	if c.Color() != color {
		t.Fatalf("Expected color %s, got %s", color, c.Color())
	}
	return c
}

func testError(t *testing.T, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(2)
		t.Errorf("%s:%d: unexpected error '%s'", file, line, err.Error())
	}
}

func testOver(t *testing.T, err error) {
	if err != game.ErrGameOver {
		_, file, line, _ := runtime.Caller(2)
		t.Errorf("%s:%d: unexpected error, expected '%s', got '%s'", file, line, game.ErrGameOver, err)
	}
}

func move(t *testing.T, c *Client, x, y int) func() {
	return func() { testError(t, c.Move(game.Position{x, y})) }
}

func wait(t *testing.T, c *Client) func() {
	return func() { testError(t, c.Wait()) }
}

func passOver(t *testing.T, c *Client) func() {
	return func() { testOver(t, c.Pass()) }
}

func pass(t *testing.T, c *Client) func() {
	return func() { testError(t, c.Pass()) }
}

func gameOver(t *testing.T, c *Client) func() {
	return func() {
		testOver(t, c.Pass())
		testOver(t, c.Move(game.Position{0, 2}))
		testOver(t, c.Move(game.Position{0, 2}))
		testError(t, c.Wait())
	}
}

// assumes gobot/server is well formed
func TestBasic(t *testing.T) {
	ts := httptest.NewServer(server.MuxerAPIv1())
	p1 := newClient(t, ts.URL, 1, game.Black)
	p2 := newClient(t, ts.URL, 2, game.White)

	// white can't play yet
	if err := p2.Move(game.Position{0, 2}); err != game.ErrWrongPlayer {
		t.Errorf("move [0,2] expected error '%s', got '%s'", game.ErrWrongPlayer, err)
	}

	v := validator(t)
	v.After(2, wait(t, p2))
	v.Before(1, move(t, p1, 0, 2))
	v.Verify(time.Second)
	v.After(4, wait(t, p1))
	// Throw in an invalid move before scheduling the real one
	// This test is slightly broken
	if err := p2.Move(game.Position{0, 2}); err != game.ErrSpotNotEmpty {
		t.Errorf("move [0,2] expected error '%s', got '%s'", game.ErrSpotNotEmpty, err)
	}
	v.Before(3, move(t, p2, 1, 0))
	v.Verify(time.Second)
	v.After(6, wait(t, p2))
	v.Before(5, pass(t, p1))
	v.Verify(time.Second)
	v.After(8, wait(t, p1))
	v.Before(7, passOver(t, p2))
	v.Verify(time.Second)
	v.After(9, gameOver(t, p1))
	v.After(9, gameOver(t, p2))
	v.After(9, gameOver(t, p1))
	v.After(9, gameOver(t, p2))
	v.Verify(time.Second)

	testPosition(t, p1, game.White, 1, 0)
	testPosition(t, p2, game.White, 1, 0)
}

func testPosition(t *testing.T, c *Client, p game.Color, x, y int) {
	found := c.state.Board[x][y]
	if p != found {
		t.Errorf("at %d-%d expected %s, found %s", p, found)
	}
}
