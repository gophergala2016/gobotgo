package server

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gophergala2016/gobotgo/game"
)

type testWriter struct{}

func (tw testWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (tw testWriter) WriteHeader(int) {}

func (tw testWriter) Header() http.Header {
	return nil
}

func TestStateHandler(t *testing.T) {
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
		startHandler(testWriter{}, test.input)
		if !gameEqual(test.expected, gameMap) {
			t.Errorf("%s not equal:\nexpected:%v\nactual:%v", test.reason, test.expected, gameMap)
		}
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
