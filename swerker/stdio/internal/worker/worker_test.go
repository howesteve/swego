package worker

import (
	"flag"
	"reflect"
	"testing"

	"github.com/astrotools/swego/swerker"
	"github.com/astrotools/swego/swerker/stdio/internal/lichdata"

	"github.com/tinylib/msgp/msgp"
)

var (
	useWorker  = flag.Bool("worker", false, "Test against actual swedenw-stdio binary.")
	workerPath = flag.String("worker.path", "../cmd/swerker/swerker-stdio", "Path to swerker-stdio binary.")
)

func TestNoFuncs(t *testing.T) {
	defer swizzle("NoFuncs", "-dangerous_no_funcs_on_init")()

	w, funcs, err := New(*workerPath)
	if w != nil {
		t.Errorf("w = %v, want: nil", w)
	}

	if funcs != nil {
		t.Errorf("funcs = %v, want: nil", funcs)
	}

	if _, ok := err.(*NoFuncsError); !ok {
		t.Errorf("err = %#v, want: %T value", err, (*NoFuncsError)(nil))
	}
}

func TestFuncsPanic(t *testing.T) {
	mockOnly(t)
	defer swizzle("FuncsPanic")()

	w, funcs, err := New(*workerPath)
	if w != nil {
		t.Errorf("w = %v, want: nil", w)
	}

	if funcs != nil {
		t.Errorf("funcs = %v, want: nil", funcs)
	}

	want := &NoFuncsError{Err: &Error{Msg: "funcs panic", Panic: true}}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("err = %#v, want: %q", err, want)
	}
}

func TestInvalidFuncsData(t *testing.T) {
	defer swizzle("InvalidFuncsData", "-dangerous_invalid_funcs_on_init")()

	w, funcs, err := New(*workerPath)
	if w != nil {
		t.Error("w != nil")
	}

	if funcs != nil {
		t.Error("funcs != nil")
	}

	if _, ok := err.(*NoFuncsError); !ok {
		t.Errorf("err.(type) = %T, want: %T", err, (*NoFuncsError)(nil))
	}
}

func TestInvalidFuncsType(t *testing.T) {
	defer swizzle("InvalidFuncsType", "-dangerous_invalid_funcs_types_on_init")()

	w, funcs, err := New(*workerPath)
	if w != nil {
		t.Error("w != nil")
	}

	if funcs != nil {
		t.Error("funcs != nil")
	}

	want := &NoFuncsError{Err: msgp.TypeError{
		Method:  msgp.ArrayType,
		Encoded: msgp.NilType,
	}}

	if !reflect.DeepEqual(err, want) {
		t.Errorf("err = %#v, want: %#v", err, want)
	}
}

func TestParseFuncs(t *testing.T) {
	defer swizzle("Call")()

	w, funcs, err := New(*workerPath)
	if err != nil {
		t.Fatal(err)
	}

	want := workerFuncs(t)
	if err := w.Exit(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(funcs, want) {
		t.Errorf("got: %+v, want: %+v", funcs, want)
	}
}

func workerFuncs(t *testing.T) Funcs {
	// New
	cmd := execCommand(*workerPath)

	wc, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}

	rc, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	// (*Worker).unmarshalFuncs
	data, err := lichdata.ReadFrom(rc)
	if err != nil {
		t.Fatal(err)
	}

	var funcs Funcs
	_, err = funcs.UnmarshalMsg(data)
	if err != nil {
		t.Fatal(err)
	}

	// (*Worker).Close
	wc.Write([]byte{'\n'})
	wc.Close()

	// (*Worker).waitForExit
	err = cmd.Wait()
	if err != nil {
		t.Fatal(err)
	}

	return funcs
}

func TestCall(t *testing.T) {
	defer swizzle("Call")()

	w, funcs, err := New(*workerPath)
	if err != nil {
		t.Fatal(err)
	}

	defer w.Exit()

	resp, crashed, err := w.Call(&swerker.Call{Func: 0}) // rpc_funcs
	if crashed {
		t.Fatalf("worker crashed: %v", err)
	}

	if err != nil {
		t.Errorf("err = %#v, want: nil", err)
	}

	var got Funcs
	if _, err := got.UnmarshalMsg(resp); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, funcs) {
		t.Errorf("got %q, want: %q", got, funcs)
	}
}

func TestCallWhenExited(t *testing.T) {
	defer swizzle("Call")()

	w, _, err := New(*workerPath)
	if err != nil {
		t.Fatal(err)
	}

	w.Exit()

	resp, crashed, err := w.Call(&swerker.Call{Func: 0}) // rpc_funcs
	if crashed {
		t.Fatalf("worker crashed: %v", err)
	}

	if resp != nil {
		t.Errorf("resp = [% x], want: nil", resp)
	}

	if err != ErrProcessExited {
		t.Errorf("err = %#v, want: %q", err, ErrProcessExited)
	}
}

func TestCallError(t *testing.T) {
	defer swizzle("CallError")()

	w, funcs, err := New(*workerPath)
	if err != nil {
		t.Fatal(err)
	}

	defer w.Exit()

	testError, ok := funcs.Lookup("test_error")
	if !ok {
		t.Fatal(`worker does not implement "test_error" function`)
	}

	resp, crashed, err := w.Call(&swerker.Call{Func: testError})
	if resp != nil {
		t.Errorf("resp = [% x], want: nil", resp)
	}

	if crashed {
		t.Error("process is crashed!")
	}

	want := &Error{Msg: "test_error called", Debug: "func=test_error"}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("err = %#v, want: %q", err, want)
	}
}

func TestCrashedCall(t *testing.T) {
	defer swizzle("CrashedCall")()

	w, funcs, err := New(*workerPath)
	if err != nil {
		t.Fatal(err)
	}

	defer w.Exit()

	testCrash, ok := funcs.Lookup("test_crash")
	if !ok {
		t.Fatal(`worker does not implement "test_crash" function`)
	}

	resp, crashed, err := w.Call(&swerker.Call{Func: testCrash})
	if !crashed {
		t.Error("process has not crashed")
	}

	if resp != nil {
		t.Errorf("resp = [% x], want: nil", resp)
	}

	want := &Error{"test_crash called", "func=test_crash", true}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("err = %#v, want: %q", err, want)
	}
}

func TestExitedAfterCrashedCall(t *testing.T) {
	defer swizzle("CrashedCall")()

	w, funcs, err := New(*workerPath)
	if err != nil {
		t.Fatal(err)
	}

	defer w.Exit()

	testCrash, ok := funcs.Lookup("test_crash")
	if !ok {
		t.Fatal(`worker does not implement "test_crash" function`)
	}

	_, crashed, _ := w.Call(&swerker.Call{Func: testCrash})
	if !crashed {
		t.Fatal("worker is not crashed")
	}

	if !w.(*worker).exited() {
		t.Error("Call returned before process has exited")
	}
}

func TestUnexpectedExit(t *testing.T) {
	mockOnly(t)
	defer swizzle("UnexpectedExit")()

	w, _, err := New(*workerPath)
	if err != nil {
		t.Fatal(err)
	}

	defer w.Exit()

	resp, crashed, err := w.Call(&swerker.Call{Func: 0}) // rpc_funcs
	if !crashed {
		t.Error("worker is not crashed")
	}

	if resp != nil {
		t.Errorf("resp = [% x], want: nil", resp)
	}

	if _, ok := err.(*UnexpectedExitError); !ok {
		t.Errorf("err.(type) = %T, want: %T", err, (*UnexpectedExitError)(nil))
	}
}
