package server

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gophergala2016/gobotgo/game"
)

type testWriter struct {
	content []byte
}

func (tw *testWriter) Write(b []byte) (int, error) {
	tw.content = b
	return 0, nil
}

func (tw testWriter) WriteHeader(int) {}

func (tw testWriter) Header() http.Header {
	return http.Header{}
}

var wg sync.WaitGroup

func TestStartHandler(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	tests := []struct {
		input    *http.Request
		expected map[GameID]*Game
		reason   string
	}{
		{
			r,
			map[GameID]*Game{
				1: &Game{
					state:   game.New(19, 180),
					players: map[GameID]game.Color{1: game.Black},
				},
			},
			"One person has joined",
		},
		{
			r,
			map[GameID]*Game{
				1: &Game{
					state:   game.New(19, 180),
					players: map[GameID]game.Color{1: game.Black, 2: game.White},
				},
				2: &Game{
					state:   game.New(19, 180),
					players: map[GameID]game.Color{1: game.Black, 2: game.White},
				},
			},
			"Two people have joined",
		},
	}
	for _, test := range tests {
		startHandler(&testWriter{}, test.input)
		if !gameEqual(test.expected, gameMap) {
			t.Errorf("%s not equal:\nexpected:%v\nactual:%v", test.reason, test.expected, gameMap)
		}
	}

}

func TestWaitHandler(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w3 := testWriter{}
	w4 := testWriter{}
	startHandler(&w3, r)
	if `{"ID":3,"color":"Black"}` != string(w3.content) {
		t.Errorf("Wait handler test %s not equal to expected id 3", string(w3.content))
	}
	startHandler(&w4, r)
	if `{"ID":4,"color":"White"}` != string(w4.content) {
		t.Errorf("Wait handler test %s not equal to expected id 4", string(w4.content))
	}
	wg.Add(1)
	go gameWait(&w4, 4)
	time.Sleep(1 * time.Second)
	playMove(&w3, 3, "[1,1]")
	wg.Wait()
	if `"valid"` != string(w3.content) {
		t.Errorf("Wait handler test move 3 not valid: %s", string(w3.content))
	}
	if `"go bot go"` != string(w4.content) {
		t.Errorf("Wait handler test wait 4 out: %s", string(w4.content))
	}

	wg.Add(1)
	go gameWait(&w3, 3)
	time.Sleep(1 * time.Second)
	playMove(&w4, 4, "[2,2]")
	wg.Wait()
	if `"valid"` != string(w4.content) {
		t.Errorf("Wait handler test move 4 not valid: %s", string(w4.content))
	}
	if `"go bot go"` != string(w3.content) {
		t.Errorf("Wait handler test wait 3 out: %s", string(w3.content))
	}

}

func gameEqual(lhs, rhs map[GameID]*Game) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i, l := range lhs {
		r, ok := rhs[i]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(l.state, r.state) {
			return false
		}
		if !reflect.DeepEqual(l.players, r.players) {
			return false
		}
	}
	return true
}

func gameWait(w http.ResponseWriter, id GameID) {
	path := fmt.Sprintf("/%d/wait/", id)
	r, _ := http.NewRequest("GET", path, nil)
	playHandler(w, r)
	wg.Done()
}

func playMove(w http.ResponseWriter, id GameID, move string) {
	path := fmt.Sprintf("/%d/move/", id)
	r, _ := http.NewRequest("POST", path, bytes.NewBufferString(move))
	playHandler(w, r)
}
