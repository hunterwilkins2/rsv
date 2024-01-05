package rsv

import (
	"bytes"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		Input  []byte
		Output [][]string
		Error  error
	}{
		{Input: []byte{'a', 'b', 'c', '\xFF', '\xFD'}, Output: [][]string{{"abc"}}, Error: nil},
		{Input: []byte{'\xFD'}, Output: [][]string{{}}, Error: nil},
		{Input: []byte{'\xFF', '\xFD'}, Output: [][]string{{""}}, Error: nil},
		{Input: append(bytes.Repeat([]byte{'a'}, 5000), '\xFF', 'b', '\xFF', '\xFD'), Output: [][]string{{strings.Repeat("a", 5000), "b"}}, Error: nil},
		{Input: []byte{
			'H', 'e', 'l', 'l', 'o', '\xFF', '\xF0', '\x9F', '\x8C', '\x8E', '\xFF', '\xFD',
			'\xFD',
			'\xFF', 'a', 'b', 'c', '\xFF', '\xFD',
		},
			Output: [][]string{
				{"Hello", "🌎"},
				{},
				{"", "abc"},
			}, Error: nil},
		{Input: []byte{}, Output: [][]string{}, Error: nil},
		{Input: []byte{'a'}, Output: [][]string{}, Error: ErrUnterminatedRow},
		{Input: []byte{'a', '\xFF'}, Output: [][]string{}, Error: ErrUnterminatedRow},
		{Input: []byte{'a', '\xFF', '\xFD', 'b', '\xFF'}, Output: [][]string{}, Error: ErrUnterminatedRow},
		{Input: []byte{'a', '\xFD'}, Output: [][]string{}, Error: ErrUnterminatedField},
		{Input: []byte{'a', '\xFF', 'b', '\xFD'}, Output: [][]string{}, Error: ErrUnterminatedField},
		{Input: []byte{'a', '\xFF', '\xFD', 'b', '\xFD'}, Output: [][]string{}, Error: ErrUnterminatedField},
	}

	for n, tt := range testCases {
		r := NewReader(bytes.NewReader(tt.Input))
		out, err := r.ReadAll()
		if err != tt.Error {
			t.Errorf("#%d Got error %s; want %s", n, err, tt.Error)
			return
		}

		if len(out) != len(tt.Output) {
			t.Errorf("#%d Got %v; want %v", n, out, tt.Output)
			return
		}
		for i := 0; i < len(out); i++ {
			if len(out[i]) != len(tt.Output[i]) {
				t.Errorf("#%d Got %v; want %v", n, out, tt.Output)
				return
			}
			for j := 0; j < len(out[i]); j++ {
				if out[i][j] != tt.Output[i][j] {

					t.Errorf("#%d Got %v; want %v", n, out, tt.Output)
				}
			}
		}
	}
}
