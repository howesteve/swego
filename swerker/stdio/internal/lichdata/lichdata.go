// Package lichdata provides Lich data element read and write functionality.
//
// See https://github.com/rentzsch/lich for more information about Lich.
package lichdata

import (
	"errors"
	"fmt"
	"io"

	"github.com/philhofer/fwd"
)

type byteReader interface {
	io.Reader
	io.ByteReader
}

// ReadFrom reads a single Lich data element from reader r.
func ReadFrom(r io.Reader) ([]byte, error) {
	if br, ok := r.(byteReader); ok {
		return readData(br)
	}

	return readData(fwd.NewReader(r))
}

// ErrNoLength is returned if no ASCII length data is found.
var ErrNoLength = errors.New("lichdata: no length")

// Type marker errors.
var (
	ErrInvalidOpenMarker  = errors.New("lichdata: invalid open marker")
	ErrInvalidCloseMarker = errors.New("lichdata: invalid close marker")
)

// MaxLengthError is returned when the buffer is bigger than the maximum value
// an int can hold.
type MaxLengthError struct {
	N uint64
}

const maxInt = int(^uint(0) >> 1)

func (e *MaxLengthError) Error() string {
	return fmt.Sprintf("lichdata: length is %d, limit %d exceeded", e.N, maxInt)
}

func readData(r byteReader) ([]byte, error) {
	c, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	var size uint64
	for '0' <= c && c <= '9' {
		size *= 10
		size += uint64(c - '0')

		c, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	if size == 0 {
		return nil, ErrNoLength
	}

	if c != '<' {
		return nil, ErrInvalidOpenMarker
	}

	if size > uint64(maxInt) {
		return nil, &MaxLengthError{size}
	}

	buf := make([]byte, int(size))
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	c, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	if c != '>' {
		return nil, ErrInvalidCloseMarker
	}

	return buf[:n], nil
}

// BufWriter represents a buffered writer.
type BufWriter interface {
	io.Writer
	io.ByteWriter
	Flush() error
}

// Writer buffers an io.Writer and can write Lich data elements to it.
type Writer struct {
	W BufWriter
}

// NewWriter returns a new Writer for underlying writer w.
func NewWriter(w io.Writer) *Writer {
	if bw, ok := w.(BufWriter); ok {
		return &Writer{bw}
	}

	return &Writer{fwd.NewWriter(w)}
}

// Write writes byte slice buf as a Lich data element to the underlying writer
// of writer w.
func (w *Writer) Write(buf []byte) (n int, err error) {
	return writeData(w.W, buf)
}

func writeData(w BufWriter, data []byte) (n int, err error) {
	i, err := fmt.Fprintf(w, "%d<", len(data))
	if err != nil {
		return 0, err
	}

	n += i
	i, err = w.Write(data)
	if err != nil {
		return 0, err
	}

	n += i
	if err = w.WriteByte('>'); err != nil {
		return 0, err
	}

	n++
	err = w.Flush()
	if err != nil {
		return 0, err
	}

	return n, nil
}
