package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gophergala2016/gobotgo/client"
	"github.com/gophergala2016/gobotgo/game"
)

var url = flag.String("url", "http://localhost:8100", "Root URL of gobotgo service")

func init() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Println("Connecting...")
	c, err := client.New(*url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Player", c.ID())

	if err := play(c); err != game.ErrGameOver {
		log.Fatal(err)
	}

	log.Println("Game one by %d", c.Opponent())
}

type decision int

func play(c *client.Client) error {
	act := action{Client: c}
	for {
		act.Choose()
		if err := act.Act(); err != nil {
			return err
		}
	}
}

const (
	wait = decision(iota)
	pass
	move
)

type action struct {
	*client.Client
	choice decision
	game.Position
}

func (act *action) Act() error {
	switch act.choice {
	case wait:
		return act.Wait()
	case pass:
		return act.Pass()
	case move:
		return act.Move(act.Position)
	default:
		return fmt.Errorf("unknown choice (%d)", act.choice)
	}
}

// Fair random
func (act *action) rand() game.Position {
	// count remaining pieces
	b := act.State()
	empty := 0
	for _, rows := range b {
		for _, cols := range rows {
			if cols == game.None {
				empty++
			}
		}
	}
	if empty == 0 {
		log.Println(b)
		log.Fatal("Board out of positions!")
	}

	position := rand.Intn(empty)
	empty = 0
	for x, rows := range b {
		for y, cols := range rows {
			if cols == game.None {
				if empty == position {
					return game.Position{x, y}
				}
				empty++
			}
		}
	}
	log.Println(position, empty)
	log.Fatal("Didn't find an empty space for rant!")
	return game.Position{}
}

func (act *action) Choose() {
	switch {
	case act.Opponent() == act.CurrentPlayer():
		act.choice = wait
	default:
		act.choice = move
		act.Position = act.rand()
	}
}
