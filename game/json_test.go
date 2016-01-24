package game

import (
	"encoding/json"
	"testing"
)

func TestColorUnmarshall(t *testing.T) {
	tests := []struct {
		Color
		result string
	}{
		{Black, `"Black"`},
		{White, `"White"`},
		{Color(5), `"None"`},
	}

	for _, test := range tests {
		data, err := json.Marshal(&test.Color)
		if err != nil {
			t.Errorf("unable to unmarshal %s: '%s'", test.Color, err.Error())
		}
		if string(data) != test.result {
			t.Errorf("expecte %s to marshal to %s, got %s", test.Color, test.result, string(data))
		}
	}
}

func TestColorMarshal(t *testing.T) {
	tests := []struct {
		json string
		Color
	}{
		{`"Black"`, Black},
		{`"White"`, White},
		{`"None"`, None},
		// Be accepting of others poorly formed JSON
		{`"bLaCK"`, Black},
		{`"wHITE"`, White},
		{`"noice is'''' none"`, None},
	}

	for _, test := range tests {
		c := Color(5)
		if err := json.Unmarshal([]byte(test.json), &c); err != nil {
			t.Errorf("unexpected error unmarshalling %s, got '%s'", test.json, err.Error())
			continue
		}
		if c != test.Color {
			t.Errorf("expected %s to unmarshal to %s, got %s", test.Color, test.Color, c)
		}
	}
}

// Just confirm we have to do nothing
func TestBoardMarshal(t *testing.T) {
	tests := []struct {
		size  int
		board []Color
		json  string
	}{
		{
			3,
			make([]Color, 9),
			`[["None","None","None"],["None","None","None"],["None","None","None"]]`,
		},
		{
			3,
			[]Color{
				White, Black, empty,
				empty, White, Black,
				Black, White, empty,
			},
			`[["White","Black","None"],["None","White","Black"],["Black","White","None"]]`,
		},
	}

	for _, test := range tests {
		b := sliceBoard(test.board, test.size)
		data, err := json.Marshal(b)
		if err != nil {
			t.Logf("board:\n%s", b)
			t.Errorf("couldn't unmarshal board becaus '%s'", err.Error())
			continue
		}
		if string(data) != test.json {
			t.Errorf("data didn't unmarshal:\nexp: %s\ngot: %s", test.json, string(data))
		}
	}
}
