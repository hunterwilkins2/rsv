package rsv

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	testCases := []struct {
		Input  [][]string
		Output string
	}{
		{Input: [][]string{{"abc"}}, Output: "abc\xff\xfd"},
		{Input: [][]string{{`a"a`}}, Output: "a\"a\xff\xfd"},
		{Input: [][]string{}, Output: ""},
		{Input: [][]string{{}}, Output: "\xfd"},
		{Input: [][]string{{""}}, Output: "\xff\xfd"},
		{Input: [][]string{{""}, {""}}, Output: "\xff\xfd\xff\xfd"},
		{Input: [][]string{{""}, {""}, {""}}, Output: "\xff\xfd\xff\xfd\xff\xfd"},
		{Input: [][]string{{""}, {""}, {"a"}}, Output: "\xff\xfd\xff\xfda\xff\xfd"},
		{Input: [][]string{{""}, {"a"}, {""}}, Output: "\xff\xfda\xff\xfd\xff\xfd"},
		{Input: [][]string{{"a"}, {""}, {""}}, Output: "a\xff\xfd\xff\xfd\xff\xfd"},
		{Input: [][]string{{""}, {"a"}, {"a"}}, Output: "\xff\xfda\xff\xfda\xff\xfd"},
		{Input: [][]string{{"a"}, {"a"}, {""}}, Output: "a\xff\xfda\xff\xfd\xff\xfd"},
		{Input: [][]string{{"a"}, {"a"}, {"a"}}, Output: "a\xff\xfda\xff\xfda\xff\xfd"},
		{Input: [][]string{{"Hello", "🌎"}}, Output: "Hello\xff🌎\xff\xfd"},
		{Input: [][]string{{"Hello", "🌎"}, {}, {"", "abc"}}, Output: "Hello\xff🌎\xff\xfd\xfd\xffabc\xff\xfd"},
	}

	for n, tt := range testCases {
		b := &strings.Builder{}
		f := NewWriter(b)
		err := f.WriteAll(tt.Input)
		if err != nil {
			t.Errorf("Unexpected error: %s\n", err)
		}

		out := b.String()
		if out != tt.Output {
			t.Errorf("#%d: out=%q; want=%q", n, out, tt.Output)
		}
	}
}

type errorWriter struct{}

func (e errorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("mock error")
}

func TestError(t *testing.T) {
	b := &bytes.Buffer{}
	f := NewWriter(b)
	f.Write([]string{"abc"})
	f.Flush()
	err := f.Error()

	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	f = NewWriter(errorWriter{})
	f.WriteAll([][]string{{"abc"}})
	f.Flush()
	err = f.Error()
	if err == nil {
		t.Error("Error should not be nil")
	}
}
