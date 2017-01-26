package worker

import (
	"reflect"
	"testing"
)

var testFuncs = Funcs{
	"rpc_funcs", // required by worker RPC system
	"test_crash",
	"test_error",
	"swe_version",
}

func TestFuncs_LastIdx(t *testing.T) {
	got := testFuncs.LastIdx()
	const want = 3

	if got != want {
		t.Errorf("LastIdx() = %d, want: %d", got, want)
	}
}

func TestFuncs_FuncsMap(t *testing.T) {
	got := testFuncs.FuncsMap()
	want := FuncsMap{
		"rpc_funcs":   0,
		"test_crash":  1,
		"test_error":  2,
		"swe_version": 3,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want: %v", got, want)
	}
}

func TestFuncs_Lookup(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		idx, ok := testFuncs.Lookup("rpc_funcs")
		if idx != 0 {
			t.Errorf("idx = %d, want: 0", idx)
		}

		if !ok {
			t.Error("ok = false, want: true")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		idx, ok := testFuncs.Lookup("not_existent")
		if idx != 0 {
			t.Errorf("idx = %d, want: 0", idx)
		}

		if ok {
			t.Error("ok = true, want: false")
		}
	})
}
