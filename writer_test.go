package rsv

import (
	"bytes"
	"encoding/csv"
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
	f.Write([]string{"abc"})
	f.Flush()
	err = f.Error()
	if err == nil {
		t.Error("Error should not be nil")
	}
}

var benchmarkWriteData = [][]string{
	{"abc", "def", "12356", "1234567890987654311234432141542132"},
	{"abc", "def", "12356", "1234567890987654311234432141542132"},
	{"abc", "def", "12356", "1234567890987654311234432141542132"},
}

func BenchmarkWrite(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := NewWriter(&bytes.Buffer{})
		err := w.WriteAll(benchmarkWriteData)
		if err != nil {
			b.Fatal(err)
		}
		w.Flush()
	}
}

func BenchmarkCSVWrite(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := csv.NewWriter(&bytes.Buffer{})
		err := w.WriteAll(benchmarkWriteData)
		if err != nil {
			b.Fatal(err)
		}
		w.Flush()
	}
}
