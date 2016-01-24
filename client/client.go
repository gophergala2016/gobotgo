package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	board  game.Board
	color  game.Color
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
	c.color = v.Color
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
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch string(data) {
	case "valid":
		return nil
	case "Game Over":
		return nil
	default:
		return fmt.Errorf("Bad request: %s", string(data))
	}
}

func (c *Client) Move(p game.Position) error {
	return c.move([]int{p.X, p.Y})
}

func (c *Client) Pass() error {
	return c.move([]int{})
}

func (c *Client) Color() game.Color {
	return c.color
}

func (c *Client) Opponent() game.Color {
	return c.Opponent()
}

func (c *Client) Wait() error {
	_, err := c.client.Get(c.playURL("wait"))
	return err
}
