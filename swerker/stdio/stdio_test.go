package stdio

import (
	"bytes"
	"errors"
	"os/exec"
	"reflect"
	"testing"

	"github.com/howesteve/swego/swerker"
	"github.com/howesteve/swego/swerker/stdio/internal/worker"

	"github.com/tinylib/msgp/msgp"
)

type testWorker struct {
	path string
	call callFunc
	exit exitFunc
}

type newFunc func(string) (worker.Worker, worker.Funcs, error)
type callFunc func(*swerker.Call) (msgp.Raw, bool, error)
type exitFunc func() error

func (w *testWorker) Exit() error { return w.exit() }
func (w *testWorker) Call(c *swerker.Call) (msgp.Raw, bool, error) {
	return w.call(c)
}

func newTestWorker(funcs worker.Funcs, call callFunc, exit exitFunc) newFunc {
	return func(path string) (worker.Worker, worker.Funcs, error) {
		return &testWorker{path, call, exit}, funcs, nil
	}
}

const workerPath = "/path/to/swerker-stdio"

func TestNew(t *testing.T) {
	dataPaths := []string{"/path/to/longfiles", "/path/to/files"}
	funcs := worker.Funcs{"rpc_funcs", "swe_set_ephe_path"}

	defer func() { newWorker = worker.New }()
	newWorker = newTestWorker(funcs, func(c *swerker.Call) (msgp.Raw, bool, error) {
		const idx = 1
		if c.Func != idx {
			t.Errorf("swe_set_ephe_path func = %d, want: %d", c.Func, idx)
		}

		var args []byte
		args = msgp.AppendArrayHeader(args, 1)
		args = msgp.AppendString(args, combineDataPaths(dataPaths))

		if !bytes.Equal([]byte(c.Args), args) {
			t.Errorf("swe_set_ephe_path args =\n\t[% x]\nwant:\n\t[% x]", c.Args, args)
		}

		if c.Ctx != nil {
			t.Errorf("swe_set_ephe_path ctx = %#v, want: nil", c.Ctx)
		}

		return msgp.Raw{0x90}, false, nil
	}, func() error {
		return nil
	})

	d, err := New(workerPath, NumWorkers(2), DataPath(dataPaths...))
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if d == nil {
		t.Fatal("d == nil")
	}

	if got := d.Path(); got != workerPath {
		t.Errorf("Path() = %q, want: %q", got, workerPath)
	}

	if got := d.DataPaths(); !reflect.DeepEqual(got, dataPaths) {
		t.Errorf("data path = %q, want: %q", got, dataPaths)
	}
}

func TestClose(t *testing.T) {
	funcs := worker.Funcs{"rpc_funcs", "swe_close"}

	defer func() { newWorker = worker.New }()
	newWorker = newTestWorker(funcs, func(c *swerker.Call) (msgp.Raw, bool, error) {
		const idx = 1
		if c.Func != idx {
			t.Errorf("swe_close func = %d, want: %d", c.Func, idx)
		}

		if !bytes.Equal([]byte(c.Args), nil) {
			t.Errorf("swe_close args =\n\t[% x]\nwant: []", c.Args)
		}

		if c.Ctx != nil {
			t.Errorf("swe_close ctx = %#v, want: nil", c.Ctx)
		}

		return msgp.Raw{0x90}, false, nil
	}, func() error {
		return nil
	})

	d, err := New(workerPath, NumWorkers(1))
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if d == nil {
		t.Fatal("d == nil")
	}

	if err := d.Close(); err != nil {
		t.Errorf("err = %v, want: nil", err)
	}
}

func TestIndexForName(t *testing.T) {
	funcs := worker.Funcs{"rpc_funcs"}

	d := &Dispatcher{funcs: funcs.FuncsMap()}
	idx, ok := d.IndexForName("rpc_funcs")
	if !ok {
		t.Error("ok = false, want: true")
	}

	if idx != 0 {
		t.Errorf("idx = %d, want: 0", idx)
	}
}

func TestDispatch(t *testing.T) {
	funcs := worker.Funcs{"rpc_funcs", "test_func"}

	defer func() { newWorker = worker.New }()
	newWorker = newTestWorker(funcs, func(c *swerker.Call) (msgp.Raw, bool, error) {
		const idx = 1
		if c.Func != idx {
			t.Errorf("test_func func = %d, want: %d", c.Func, idx)
		}

		if !bytes.Equal([]byte(c.Args), nil) {
			t.Errorf("test_func args =\n\t[% x]\nwant: []", c.Args)
		}

		if c.Ctx != nil {
			t.Errorf("test_func ctx = %#v, want: nil", c.Ctx)
		}

		return msgp.Raw{0x90}, false, nil
	}, func() error {
		return nil
	})

	d, err := New(workerPath, NumWorkers(1))
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if d == nil {
		t.Fatal("d == nil")
	}

	t.Run("Unimplemented", func(t *testing.T) {
		fn := funcs.LastIdx() + 1

		data, err := d.Dispatch(&swerker.Call{Func: fn})
		if data != nil {
			t.Errorf("data = [% x], want: nil", data)
		}

		want := &UnimplementedError{Func: fn}
		if !reflect.DeepEqual(err, want) {
			t.Errorf("err = %v, want: %v", err, want)
		}
	})

	t.Run("Call", func(t *testing.T) {
		data, err := d.Dispatch(&swerker.Call{Func: funcs.LastIdx()})
		if err != nil {
			t.Errorf("err = %v, want: nil", err)
		}

		if !bytes.Equal(data, msgp.Raw{0x90}) {
			t.Errorf("data = [% x], want: [90]", data)
		}
	})

	if err := d.Close(); err != nil {
		t.Errorf("err = %v, want: nil", err)
	}
}

func TestCrash(t *testing.T) {
	funcs := worker.Funcs{"rpc_funcs", "test_crash"}
	exitErr := &exec.Error{Name: workerPath, Err: errors.New("some exit error")}

	defer func() { newWorker = worker.New }()
	newWorker = newTestWorker(funcs, func(c *swerker.Call) (msgp.Raw, bool, error) {
		const idx = 1
		if c.Func != idx {
			t.Errorf("test_func func = %d, want: %d", c.Func, idx)
		}

		if !bytes.Equal([]byte(c.Args), nil) {
			t.Errorf("test_func args =\n\t[% x]\nwant: []", c.Args)
		}

		if c.Ctx != nil {
			t.Errorf("test_func ctx = %#v, want: nil", c.Ctx)
		}

		return nil, true, &worker.Error{
			Msg:   "test_crash called",
			Debug: "func=test_crash",
			Panic: true,
		}
	}, func() error {
		return exitErr
	})

	d, err := New(workerPath, NumWorkers(2), OnExitError(func(err error) {
		if err != exitErr {
			t.Errorf("err = %v, want: %#v", err, exitErr)
		}
	}), OnNewError(func(err error) {
		if err != nil {
			t.Errorf("err = %v, want: nil", err)
		}
	}))

	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if d == nil {
		t.Fatal("d == nil")
	}

	data, err := d.Dispatch(&swerker.Call{Func: funcs.LastIdx()})
	if data != nil {
		t.Errorf("data = [% x], want: nil", data)
	}

	want := &worker.Error{
		Msg:   "test_crash called",
		Debug: "func=test_crash",
		Panic: true,
	}

	if !reflect.DeepEqual(err, want) {
		t.Errorf("err = %v, want: %#v", err, want)
	}

	if err := d.Close(); err != nil {
		t.Errorf("err = %v, want: nil", err)
	}
}

func TestVersion(t *testing.T) {
	funcs := worker.Funcs{"rpc_funcs", "swe_version"}
	const version = "2.00"

	defer func() { newWorker = worker.New }()
	newWorker = newTestWorker(funcs, func(c *swerker.Call) (msgp.Raw, bool, error) {
		const idx = 1
		if c.Func != idx {
			t.Errorf("swe_version func = %d, want: %d", c.Func, idx)
		}

		if c.Args != nil {
			t.Errorf("swe_version args = [% x], want: nil", c.Args)
		}

		if c.Ctx != nil {
			t.Errorf("swe_version ctx = %#v, want: nil", c.Ctx)
		}

		return msgp.Raw("\x91\xa4" + version), false, nil
	}, func() error {
		return nil
	})

	v, err := Version(workerPath)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if v != version {
		t.Errorf("v = %q, want: %q", v, version)
	}
}
