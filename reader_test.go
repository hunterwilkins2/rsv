package rsv

import (
	"bytes"
	"encoding/csv"
	"io"
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

type nTimes struct {
	s   string
	n   int
	off int
}

func (r *nTimes) Read(p []byte) (n int, err error) {
	for {
		if r.n <= 0 || r.s == "" {
			return n, io.EOF
		}
		n0 := copy(p, r.s[r.off:])
		p = p[n0:]
		n += n0
		r.off += n0
		if r.off == len(r.s) {
			r.off = 0
			r.n--
		}
		if len(p) == 0 {
			return
		}
	}
}

const benchmarkReadData = "x\xFFy\xFFz\xFF\xFF\xFDx\xFFy\xFF\xFF\xFF\xFDx\xFF\xFF\xFF\xFF\xFD\xFF\xFF\xFF\xFF\xFDx\xFFy\xFFz\xFFw\xFF\xFDx\xFFy\xFFz\xFF\xFF\xFDx\xFFy\xFF\xFF\xFF\xFDx\xFF\xFF\xFF\xFF\xFD\xFF\xFF\xFF\xFF\xFD"

func BenchmarkRead(b *testing.B) {
	b.ReportAllocs()
	r := NewReader(&nTimes{s: benchmarkReadData, n: b.N})
	for {
		_, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			b.Fatal(err)
		}
	}
}

const benchmarkCSVData = `x,y,z,w
x,y,z,
x,y,,
x,,,
,,,
"x","y","z","w"
"x","y","z",""
"x","y","",""
"x","","",""
"","","",""
`

func BenchmarkCSVRead(b *testing.B) {
	b.ReportAllocs()
	r := csv.NewReader(&nTimes{s: benchmarkCSVData, n: b.N})
	for {
		_, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			b.Fatal(err)
		}
	}
}
