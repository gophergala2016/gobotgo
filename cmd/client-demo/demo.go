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
var competitive = flag.Bool("competitive", false, "Use a slightly more aggressive algorithm")

func init() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if *competitive {
		log.Println("Playing competitively")
	}
	log.Println("Connecting...")
	c, err := client.New(*url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Player", c.ID())

	if err := play(c); err != game.ErrGameOver {
		log.Fatal(err)
	}

	printWinner(c)
}

func score(c *client.Client) (black, white int) {
	black, white = c.State().Score()
	black += c.BlackStones().Captured
	white += c.WhiteStones().Captured
	return
}

func printWinner(c *client.Client) {
	winner := game.None
	black, white := score(c)
	switch {
	case black > white:
		winner = game.Black
	case black < white:
		winner = game.White
	}
	log.Printf("Score: b: %d, w: %d", black, white)
	log.Println("Game won by", winner)
}

type decision int

func play(c *client.Client) error {
	act := action{Client: c}
	for {
		act.Choose()
		if err := act.Act(); err != nil {
			switch err {
			case game.ErrRepeatState:
				act.Pass()

			case game.ErrNoStones:
				fallthrough
			case game.ErrSpotNotEmpty:
				fallthrough
			case game.ErrSelfCapture:
				log.Printf("(%d): invalid move %+v: '%s'", act.ID(), act.Position, err.Error())
				// Make sure we update the state after this
				act.choice = wait

			default:
				return err
			}
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

// Relative advantage for c with given board
func advantage(c *client.Client, b game.Board) int {
	black, white := b.Score()
	diff := black - white
	if c.Color() == game.White {
		diff = -diff
	}
	return diff
}

// Picks most competitive immediate move, or random.
//   - Passes if everything is worse/indifferent (won't fill territory)
//   - Picks random competitive move
func (act *action) competitive() {
	choices := map[int][]game.Position{}
	b := act.State()
	nothing := advantage(act.Client, b)
	for x, rows := range b {
		for y, cols := range rows {
			if cols != game.None {
				continue
			}
			c := b.Copy()
			pos := game.Position{x, y}
			taken, err := c.Apply(game.Move{act.Color(), pos})
			if err != nil {
				continue
			}
			score := advantage(act.Client, c) + taken
			choices[score] = append(choices[score], pos)
		}
	}
	// we could have a negative advantage
	max := nothing
	for k := range choices {
		if k > max {
			max = k
		}
	}
	if max <= nothing {
		act.choice = pass
	} else {
		act.choice = move
		options := choices[max]
		act.Position = options[rand.Intn(len(options))]
	}
}

// Fair random
func (act *action) rand() {
	// count remaining pieces
	act.choice = move
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
					act.X, act.Y = x, y
					return
				}
				empty++
			}
		}
	}
	log.Println(position, empty)
	log.Fatal("Didn't find an empty space for rant!")
}

func (act *action) Choose() {
	stones := act.WhiteStones()
	if act.Color() == game.Black {
		stones = act.BlackStones()
	}
	switch {
	case act.Opponent() == act.CurrentPlayer():
		act.choice = wait
	case stones.Remaining <= 0:
		act.choice = pass
	case *competitive:
		act.competitive()
	default:
		act.rand()
	}
}
