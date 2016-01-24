package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gophergala2016/gobotgo/game"
	"github.com/gophergala2016/gobotgo/server"
)

type state struct {
	player   game.Color
	board    game.Board
	lastMove game.Move
}

// Client connects to a URL and begins a game
type Client struct {
	client http.Client
	url    string
	id     server.GameID
	player game.Color
	state
}

func (c *Client) playURL(s string) string {
	return fmt.Sprintf("%s/api/v1/game/play/%d/%s", c.url, c.id, s)
}

func (c *Client) retrieve(s string, v interface{}) error {
	resp, err := c.client.Get(c.url + "/api/v1/game/" + s)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	return dec.Decode(v)
}

func New(url string) (*Client, error) {
	c := Client{
		client: http.Client{Timeout: time.Minute * 10},
		url:    url,
	}
	v := struct {
		ID    server.GameID
		Color game.Color
	}{}
	if err := c.retrieve("start/", &v); err != nil {
		return nil, err
	}
	c.id = v.ID
	c.player = v.Color
	return &c, nil
}

func (c *Client) ID() server.GameID {
	return c.id
}

func (c *Client) State() game.Board {
	return nil
	// return c.board.copy()
}

func (c *Client) move(m []int) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(
		c.playURL("move"),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var response string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	switch response {
	case "valid":
		return nil
	case game.ErrSpotNotEmpty.Error():
		return game.ErrSpotNotEmpty
	case game.ErrOutOfBounds.Error():
		return game.ErrOutOfBounds
	case game.ErrWrongPlayer.Error():
		return game.ErrWrongPlayer
	case game.ErrRepeatState.Error():
		return game.ErrRepeatState
	case game.ErrGameOver.Error():
		return game.ErrGameOver
	default:
		return fmt.Errorf("Bad request: %s", response)
	}
}

func (c *Client) Move(p game.Position) error {
	return c.move([]int{p.X, p.Y})
}

func (c *Client) Pass() error {
	return c.move([]int{})
}

func (c *Client) Color() game.Color {
	return c.player
}

func (c *Client) Opponent() game.Color {
	return c.player.Opponent()
}

func (c *Client) Wait() error {
	_, err := c.client.Get(c.playURL("wait"))
	return err
}

func (c *Client) CurrentPlayer() game.Color {
	return c.state.player
}
