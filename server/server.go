package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gophergala2016/gobotgo/game"
)

type GameID uint64

type gameIDChan chan GameID

type Game struct {
	state   *game.State
	players map[GameID]game.Color
	turn    chan game.Color
}

var nextGame = make(chan *Game, 1)
var masterID = make(gameIDChan, 1)
var gameMap = map[GameID]*Game{}

const notSet GameID = 0

func init() {
	nextGame <- &Game{}
	masterID <- 1
}

func main() {
	port := ":8100"
	root := "/api/v1"
	mux := http.NewServeMux()
	mux.HandleFunc(filepath.Join(root, "game/start")+"/", startHandler)
	play := filepath.Join(root, "game/play") + "/"
	mux.Handle(play, http.StripPrefix(play, http.HandlerFunc(playHandler)))
	log.Fatal(http.ListenAndServe(port, mux))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// name := r.FormValue("name")
	size := parseSize(r)
	// _, poll := r.Form["poll"]
	g := <-nextGame
	if g.state == nil {
		g.state = game.New(size, 180)
		g.players = map[GameID]game.Color{}
		g.turn = make(chan game.Color, 1)
		g.turn <- game.Black
	}
	id := masterID.next()
	var c game.Color
	switch {
	case len(g.players) == 0:
		c = game.Black
		nextGame <- g
	case len(g.players) == 1:
		c = game.White
		nextGame <- &Game{}
	}
	gameMap[id] = g
	g.players[id] = c
	s := struct {
		ID    GameID     `json: "id"`
		Color game.Color `json:"color"`
	}{
		id, c,
	}
	b, err := json.Marshal(&s)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("JSON marshal error for %v: %s", s, err.Error()))
	}
	w.Write(b)
}

func parseSize(r *http.Request) int {
	size, _ := strconv.Atoi(r.FormValue("size"))
	if size == 0 {
		return 19
	}
	return size
}

func (c gameIDChan) next() GameID {
	id := <-c
	c <- id + 1
	return id
}

func (id GameID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func (g Game) stateHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(g.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "JSON error")
		log.Println(err)
		return
	}
	w.Write(b)
}

func (g *Game) moveHandler(w http.ResponseWriter, r *http.Request, id GameID) {
	t := <-g.turn
	p, ok := g.players[id]
	if !ok {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("No player for id %d", id))
		g.turn <- t
		return
	}
	m, err := g.parseMove(r, p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		g.turn <- t
		return
	}
	if err := g.state.Move(m); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		g.turn <- t
		return
	}
	w.Write([]byte("valid"))
	g.changeTurn(t)
}

func (g *Game) waitHandler(w http.ResponseWriter, r *http.Request, id GameID) {
	p, ok := g.players[id]
	if !ok {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("No player for id %d", id))
		return
	}
	for t := <-g.turn; t != p; t = <-g.turn {
	}
	w.Write([]byte("go bot go"))
}

func (g Game) parseMove(r *http.Request, c game.Color) (game.Move, error) {
	d := json.NewDecoder(r.Body)
	var move [2]int
	if err := d.Decode(&move); err != nil {
		return game.Move{}, fmt.Errorf("Decode move error: %s", err.Error())
	}
	return game.Move{
		Player:   c,
		Position: game.Position{move[0], move[1]},
	}, nil
}

func (g Game) changeTurn(t game.Color) {
	switch t {
	case game.Black:
		g.turn <- game.White
	default:
		g.turn <- game.Black
	}
}
