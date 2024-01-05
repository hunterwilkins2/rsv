# Row of String Value (RSV) 

[![Go Reference](https://pkg.go.dev/badge/github.com/hunterwilkins2/rsv/slug.svg)](https://pkg.go.dev/github.com/hunterwilkins2/rsv)
![Unit tests](https://github.com/hunterwilkins2/rsv/actions/workflows/test.yaml/badge.svg)

Package `rsv` reads and writes Row of String Value (rsv) files as described
by [Stenway's video](https://www.youtube.com/watch?v=tb_70o6ohMA).
See original repository for more details https://github.com/Stenway/RSV-Challenge

A rsv file contains zero or more records of one or more fields per record.
Each field is seperated by the terminate value byte (0xFF).
Each record is seperated by the terminate row byte (0xFD).

rsv files have an advantage over other data encoding formats, such as csv, in that
an rsv file doesn't require a header for its fields, supports UTF-8 encodings,
avoid delimiter collisions, and multiple rsv files can be simply concatenated.
## Installation
```
go get github.com/hunterwilkins2/rsv
```

## Example
```go
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
```