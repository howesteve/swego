// Package swerker provides an interface for interfacing with worker processes.
package swerker

import "github.com/tinylib/msgp/msgp"

// Dispatcher dispatches calls to a backend worker.
type Dispatcher interface {
	// IndexForName look up an index for a function name. If the function name is
	// found the index and true is returned, otherwise 0 and false. The lookup is
	// done in constant time if possible.
	//
	// This mapping is specific for a backend and Swiss Ephemeris version and
	// must not be cached by the client.
	IndexForName(string) (uint8, bool)

	// Dispatch dispatches a call to a backend.
	Dispatch(*Call) (msgp.Raw, error)
}

//go:generate msgp

//msgp:tuple Call
//msgp:tuple CtxCall

// A Call represents a call to the Swiss Ephemeris. It has optionally a
// context, one or more context calls that are executed before the actual call.
type Call struct {
	Ctx  []*CtxCall
	Func uint8
	Args msgp.Raw
}

// A CtxCall is a context call.
type CtxCall struct {
	Func uint8
	Args msgp.Raw
}
