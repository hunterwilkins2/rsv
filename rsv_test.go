package rsv_test

import (
	"os"
	"testing"

	"github.com/hunterwilkins2/rsv"
)

func TestEmptyFile(t *testing.T) {
	f, err := os.Create("/tmp/empty.rsv")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		err := os.Remove("/tmp/empty.rsv")
		if err != nil {
			t.Fatal(err)
		}
	}()

	r := rsv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		t.Errorf("Unexpected Error: %s", err)
		return
	}

	if len(data) != 0 {
		t.Errorf("Expected data to be empty")
	}
}

func TestReadWriteFile(t *testing.T) {
	data := [][]string{
		{"abc", "def", "12356", "1234567890987654311234432141542132"},
		{"abc", "def", "12356", "1234567890987654311234432141542132"},
		{"abc", "def", "12356", "1234567890987654311234432141542132"},
	}

	f, err := os.Create("/tmp/test.rsv")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove("/tmp/test.rsv")
		if err != nil {
			t.Fatal(err)
		}
	}()

	w := rsv.NewWriter(f)
	err = w.WriteAll(data)
	f.Close()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	f, err = os.Open("/tmp/test.rsv")
	if err != nil {
		t.Fatal(err)
	}

	r := rsv.NewReader(f)
	out, err := r.ReadAll()
	f.Close()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(out) != len(data) {
		t.Errorf("Got %v; want %v", out, data)
		return
	}
	for i := 0; i < len(out); i++ {
		if len(out[i]) != len(data[i]) {
			t.Errorf("Got %v; want %v", out, data)
			return
		}
		for j := 0; j < len(out[i]); j++ {
			if out[i][j] != data[i][j] {

				t.Errorf("Got %v; want %v", out, data)
			}
		}
	}
}

func TestAppendFiles(t *testing.T) {
	data := [][]string{
		{"abc", "def", "12356", "1234567890987654311234432141542132"},
		{"abc", "def", "12356", "1234567890987654311234432141542132"},
		{"abc", "def", "12356", "1234567890987654311234432141542132"},
	}

	f, err := os.Create("/tmp/append.rsv")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove("/tmp/append.rsv")
		if err != nil {
			t.Fatal(err)
		}
	}()

	w := rsv.NewWriter(f)
	err = w.WriteAll(data)
	f.Close()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	f, err = os.OpenFile("/tmp/append.rsv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}

	w = rsv.NewWriter(f)
	err = w.WriteAll(data)
	f.Close()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	f, err = os.Open("/tmp/append.rsv")
	if err != nil {
		t.Fatal(err)
	}

	r := rsv.NewReader(f)
	out, err := r.ReadAll()
	f.Close()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	concatData := append(data, data...)
	if len(out) != len(concatData) {
		t.Errorf("Got %v; want %v", out, concatData)
		return
	}
	for i := 0; i < len(out); i++ {
		if len(out[i]) != len(concatData[i]) {
			t.Errorf("Got %v; want %v", out, concatData)
			return
		}
		for j := 0; j < len(out[i]); j++ {
			if out[i][j] != concatData[i][j] {

				t.Errorf("Got %v; want %v", out, concatData)
			}
		}
	}
}
