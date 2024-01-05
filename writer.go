// Package rsv reads and writes Row of String Value (rsv) files as described
// by Stenway's video https://www.youtube.com/watch?v=tb_70o6ohMA.
// See original repository for more details https://github.com/Stenway/RSV-Challenge
//
// A rsv file contains zero or more records of one or more fields per record.
// Each field is seperated by the terminate value byte (0xFF).
// Each record is seperated by the terminate row byte (0xFD).
//
// rsv files have an advantage over other data encoding formats, such as csv, in that
// an rsv file doesn't require a header for its fields, supports UTF-8 encodings,
// avoid delimiter collisions, and multiple rsv files can be simply concatenated.
package rsv

import (
	"bufio"
	"io"
)

const (
	terminateValue = 0xFF
	terminateRow   = 0xFD
)

// A Writer writes records using rsv encoding
type Writer struct {
	w *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// Error returns any error that might have occurred during
// [Writer.Write] or [Writer.Flush]
func (w *Writer) Error() error {
	_, err := w.w.Write(nil)
	return err
}

// Flush writes any buffered data to the underlying [io.Writer]
func (w *Writer) Flush() {
	w.w.Flush()
}

// Write writes a single rsv record to w
// A record is a slice of strings.
// Writes are buffered, so [Writer.Flush] must eventually be called
// to ensure that the record is written to the underlying [io.Writer]
func (w *Writer) Write(record []string) error {
	for _, field := range record {
		_, err := w.w.Write([]byte(field))
		if err != nil {
			return err
		}
		err = w.w.WriteByte(terminateValue)
		if err != nil {
			return err
		}
	}
	return w.w.WriteByte(terminateRow)
}

// WriteAll writes multiple records to w using [Writer.Write]
// and then calls [Writer.Flush], returning any errors from the Flush
func (w *Writer) WriteAll(records [][]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	return w.w.Flush()
}
