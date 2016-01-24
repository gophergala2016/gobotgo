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
