package rsv_test

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hunterwilkins2/rsv"
)

func ExampleReader() {
	records := [][]string{
		{"Hello", "🌎"},
		{},
		{"", "abc"},
	}

	b := bytes.NewBuffer([]byte{})
	w := rsv.NewWriter(b)

	// Write records into file/buffer in rsv format
	err := w.WriteAll(records)
	if err != nil {
		log.Fatal(err)
	}

	// Read records from file/buffer
	r := rsv.NewReader(b)
	data, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
	// Output:
	// [[Hello 🌎] [] [ abc]]
}
