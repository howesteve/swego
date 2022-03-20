package worker

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/howesteve/swego/swerker/stdio/internal/lichdata"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp -file $GOFILE

//msgp:ignore Error

// Error represent an error in a worker subprocess.
type Error struct {
	Msg   string
	Debug string
	Panic bool
}

func (e *Error) Error() string {
	if e.Debug == "" {
		return e.Msg
	}

	return fmt.Sprintf("%s [%s]", e.Msg, e.Debug)
}

// ErrorMap is returned by the worker process in case of a RPC error.
// It contains always an "err" key describing the error.
// Optionally it contains a "dbg" key with additional (debug) information.
type ErrorMap map[string]string

func (em ErrorMap) error() *Error {
	return &Error{Msg: em["err"], Debug: em["dbg"]}
}

var readerPool = sync.Pool{New: func() interface{} {
	return new(bytes.Reader)
}}

func newReader(buf []byte) *bytes.Reader {
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(buf)
	return r
}

func freeReader(r *bytes.Reader) {
	readerPool.Put(r)
}

type stdoutWriter struct {
	write func(msgp.Raw)
}

func (w *stdoutWriter) Write(data []byte) (int, error) {
	r := newReader(data)
	msg, err := lichdata.ReadFrom(r)
	freeReader(r)

	if err != nil {
		return 0, err
	}

	w.write(msg)
	return len(data), nil
}

type stderrWriter struct {
	report func(*Error)
	debug  string
}

const (
	prefixDebug    = "DEBUG: "
	prefixError    = "ERROR: "
	prefixDebugLen = len(prefixDebug)
	prefixErrorLen = len(prefixError)
)

func (w *stderrWriter) Write(data []byte) (int, error) {
	r := newReader(data)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, prefixDebug):
			w.debug = line[prefixDebugLen:]

		case strings.HasPrefix(line, prefixError):
			w.report(&Error{line[prefixErrorLen:], w.debug, true})
			w.debug = ""
		}
	}

	freeReader(r)
	if err := s.Err(); err != nil {
		return 0, err
	}

	return len(data), nil
}
