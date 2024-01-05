package rsv

import (
	"bufio"
	"errors"
	"io"
)

// These are the errors that can be returned from [Reader.Read]
var (
	ErrUnterminatedField = errors.New("field is not terminated")
	ErrUnterminatedRow   = errors.New("row is not terminated")
)

type Reader struct {
	r *bufio.Reader
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

// Read reads one record from r.
// If parsing a field and the terminate row byte is found before the terminate value byte,
// the error [rsv.ErrUnterminatedField] and an empty slice of fields is returned.
// If the row is not empty and the terminate row byte is not the last character,
// the error [rsv.ErrUnterminatedRow] and an empty slice of fields is returned.
func (r *Reader) Read() ([]string, error) {
	b, err := r.r.ReadBytes(terminateRow)
	if len(b) == 0 && err == io.EOF {
		return nil, io.EOF
	}
	if b[len(b)-1] != terminateRow {
		return nil, ErrUnterminatedRow
	}

	s := 0 // starting search index
	startedReading := false
	var record []string
	for i := 0; i < len(b)-1; i++ {
		if b[i] == terminateValue {
			length := i - s
			startedReading = false
			record = append(record, string(b[s:s+length]))
			s = i + 1
			continue
		}
		startedReading = true
	}
	if startedReading {
		return nil, ErrUnterminatedField
	}
	return record, nil
}

// ReadAll reads all the remaining records from r.
// Each record is a slice of fields.
func (r *Reader) ReadAll() ([][]string, error) {
	var records [][]string
	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}
