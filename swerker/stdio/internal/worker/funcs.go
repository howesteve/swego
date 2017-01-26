package worker

//go:generate msgp -file $GOFILE

// Funcs hold the list of function names returned from the worker.
// The index of the function is the element index.
type Funcs []string

// LastIdx retuns the last valid index i found in functions list s.
func (s Funcs) LastIdx() uint8 { return uint8(len(s)) - 1 }

// FuncsMap returns a map from function name to the corresponding index.
func (s Funcs) FuncsMap() FuncsMap {
	m := make(FuncsMap)
	for idx, name := range s {
		m[name] = uint8(idx)
	}

	return m
}

// Lookup does a linear search in s for function fn.
func (s Funcs) Lookup(fn string) (idx uint8, ok bool) {
	for i, name := range s {
		if name == fn {
			return uint8(i), true
		}
	}

	return 0, false
}

//msgp:ignore FuncsMap

// FuncsMap maps a function name to a RPC function index.
type FuncsMap map[string]uint8
