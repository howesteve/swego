package worker

import (
	"io"
	"reflect"
	"testing"
)

func TestStderrWriter(t *testing.T) {
	cases := []struct {
		name string
		in   func(io.Writer)
		want []*Error
	}{
		{
			"Basic/Unbuffered",
			func(w io.Writer) {
				io.WriteString(w, prefixDebug+"debug\n"+prefixError+"error")
			},
			[]*Error{&Error{"error", "debug", true}},
		},
		{
			"Basic/Buffered",
			func(w io.Writer) {
				io.WriteString(w, prefixDebug+"debug")
				io.WriteString(w, prefixError+"error")
			},
			[]*Error{&Error{"error", "debug", true}},
		},
		{
			"Input/ErrErr",
			func(w io.Writer) {
				io.WriteString(w, prefixError+"length of input data expected")
				io.WriteString(w, prefixError+"reading unexpected EOF (body)")
			},
			[]*Error{
				&Error{Msg: "length of input data expected", Panic: true},
				&Error{Msg: "reading unexpected EOF (body)", Panic: true},
			},
		},
		{
			"Input/ErrDebug",
			func(w io.Writer) {
				io.WriteString(w, prefixError+"length of input data expected")
				io.WriteString(w, prefixDebug+"func=2")
				io.WriteString(w, prefixError+"invalid index (function)")
			},
			[]*Error{
				&Error{Msg: "length of input data expected", Panic: true},
				&Error{"invalid index (function)", "func=2", true},
			},
		},
		{
			"Input/DebugErr",
			func(w io.Writer) {
				io.WriteString(w, prefixDebug+"func=2")
				io.WriteString(w, prefixError+"invalid index (function)")
				io.WriteString(w, prefixError+"length of input data expected")
			},
			[]*Error{
				&Error{"invalid index (function)", "func=2", true},
				&Error{Msg: "length of input data expected", Panic: true},
			},
		},
		{
			"Input/DebugDebug",
			func(w io.Writer) {
				io.WriteString(w, prefixDebug+"c='9' c=57")
				io.WriteString(w, prefixError+"reading unexpected close type marker")
				io.WriteString(w, prefixDebug+"c='>' c=62")
				io.WriteString(w, prefixError+"reading unexpected open type marker")
			},
			[]*Error{
				&Error{"reading unexpected close type marker", "c='9' c=57", true},
				&Error{"reading unexpected open type marker", "c='>' c=62", true},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got []*Error
			w := &stderrWriter{report: func(e *Error) {
				got = append(got, e)
			}}

			c.in(w)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got %q, want: %q", got, c.want)
			}
		})
	}
}
