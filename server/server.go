package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gophergala2016/gobotgo/game"
)

type gameID uint64

type gameIDChan chan gameID

type Game struct {
	state game.State
	black gameID
	white gameID
}

var nextGame = make(chan *Game, 1)
var masterID = make(gameIDChan, 1)
var gameMap = map[gameID]*Game{}

const notSet gameID = 0

func init() {
	nextGame <- &Game{}
	masterID <- 1
}

func main() {
	root := "/api/v1"
	port := ":8100"
	mux := http.NewServeMux()
	mux.HandleFunc(filepath.Join(root, "game/start")+"/", startHandler)
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// name := r.FormValue("name")
	size := parseSize(r)
	// _, poll := r.Form["poll"]
	g := <-nextGame
	if g.state.Empty() {
		g.state = game.New(size)
	}
	switch {
	case g.black == notSet:
		g.black = masterID.next()
		gameMap[g.black] = g
		nextGame <- g
		w.Write([]byte(g.black.String()))
	case g.white == notSet:
		g.white = masterID.next()
		gameMap[g.white] = g
		nextGame <- &Game{}
		w.Write([]byte(g.white.String()))
	}
}

func parseSize(r *http.Request) int {
	size, _ := strconv.Atoi(r.FormValue("size"))
	if size == 0 {
		return 19
	}
	return size
}

func (c gameIDChan) next() gameID {
	id := <-c
	c <- id + 1
	return id
}

func (id gameID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
