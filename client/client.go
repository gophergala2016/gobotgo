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

// Client connects to a URL and begins a game
type Client struct {
	client http.Client
	url    string
	id     server.GameID
	player game.Color
	state  game.PublicState
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

	if err := c.loadState(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) ID() server.GameID {
	return c.id
}

func (c *Client) loadState() error {
	ps := game.PublicState{}
	resp, err := c.client.Get(c.playURL("state"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&ps); err != nil {
		return err
	}
	c.state = ps
	return nil
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

	err = c.loadState()

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
	case game.ErrSelfCapture.Error():
		return game.ErrSelfCapture
	case game.ErrGameOver.Error():
		return game.ErrGameOver
	default:
		return fmt.Errorf("Bad request: %s", response)
	}

	if err != nil {
		return err
	}
	return nil
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
	err2 := c.loadState()
	switch {
	case err != nil:
		return err
	case err2 != nil:
		return err2
	default:
		return nil
	}
}

func (c *Client) CurrentPlayer() game.Color {
	return c.state.CurrentPlayer
}

// State returns a copy of the board state
func (c *Client) State() game.Board {
	return c.state.Board.Copy()
}
