package client

import (
	"fmt"
	"net/http/httptest"
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
		t.Errorf("unexpected error '%s'", err.Error())
	}
}

func move(t *testing.T, c *Client, x, y int) func() {
	return func() { testError(t, c.Move(game.Position{x, y})) }
}

func wait(t *testing.T, c *Client) func() {
	return func() { testError(t, c.Wait()) }
}

func pass(t *testing.T, c *Client) func() {
	return func() { testError(t, c.Pass()) }
}

// assumes gobot/server is well formed
func TestBasic(t *testing.T) {
	ts := httptest.NewServer(server.MuxerAPIv1())
	p1 := newClient(t, ts.URL, 1, game.Black)
	p2 := newClient(t, ts.URL, 2, game.White)

	v := validator(t)
	v.After(2, wait(t, p2))
	v.Before(1, move(t, p1, 0, 2))
	v.Verify(time.Second)
	v.After(4, wait(t, p1))
	// Throw in an invalid move before scheduling the real one
	err := p2.Move(game.Position{0, 2})
	if err == nil {
		t.Errorf("move [0,2] expected error")
	} else {
		fmt.Print(err.Error())
	}
	v.Before(3, move(t, p2, 1, 0))
	v.Verify(time.Second)
	v.After(6, wait(t, p2))
	v.Before(5, pass(t, p1))
	v.Verify(time.Second)
	v.After(8, wait(t, p1))
	v.Before(7, pass(t, p2))
	v.Verify(time.Second)
}