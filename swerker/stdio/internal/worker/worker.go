// Package worker provides an interface to the swerker-stdio worker binary.
package worker

import (
	"errors"
	"os/exec"

	"github.com/howesteve/swego/swerker"
	"github.com/howesteve/swego/swerker/stdio/internal/lichdata"

	"github.com/tinylib/msgp/msgp"
)

// Worker runs and interacts with a swerker-stdio subprocess.
type Worker interface {
	Call(c *swerker.Call) (data msgp.Raw, crashed bool, err error)
	Exit() error
}

type worker struct {
	path string
	cmd  *exec.Cmd
	in   *lichdata.Writer
	out  chan msgp.Raw
	err  chan *Error

	waitErr error // valid after exited is closed
	waited  chan struct{}
}

// New runs the swerker-stdio binary found at the specified path as process and
// returns the RPC functions it exposes.
func New(path string) (Worker, Funcs, error) {
	w := &worker{
		path:   path,
		out:    make(chan msgp.Raw),
		err:    make(chan *Error),
		waited: make(chan struct{}),
	}

	if err := w.startProcess(); err != nil {
		return nil, nil, err
	}

	funcs, err := w.unmarshalFuncs()
	if err != nil {
		return nil, nil, err
	}

	return w, funcs, nil
}

// for testing
var execCommand = exec.Command
var execCmdArgs []string

func (w *worker) startProcess() error {
	w.cmd = execCommand(w.path, execCmdArgs...)

	in, err := w.cmd.StdinPipe()
	if err != nil {
		return err
	}

	w.in = lichdata.NewWriter(in)
	w.cmd.Stdout = &stdoutWriter{write: func(out msgp.Raw) { w.out <- out }}
	w.cmd.Stderr = &stderrWriter{report: func(err *Error) { w.err <- err }}

	if err = w.cmd.Start(); err != nil {
		return err
	}

	go w.waitForExit()
	return nil
}

func (w *worker) waitForExit() {
	w.waitErr = w.cmd.Wait()
	close(w.waited)
}

// Exit terminates the subprocess. If the process doesn't complete successfully
// the error is of type *exec.ExitError. Other error types may be returned for
// I/O problems.
func (w *worker) Exit() error {
	if !w.exited() {
		w.in.W.WriteByte('\n')
		w.in.W.Flush()
		<-w.waited
	}

	return w.waitErr
}

func (w *worker) exited() bool {
	select {
	default:
		return false
	case <-w.waited:
		return true
	}
}

// NoFuncsError is returned when no initial funcs are returned.
type NoFuncsError struct {
	// Err is the underlying error why no funcs are available.
	Err error
}

func (e *NoFuncsError) Error() string {
	return "worker: no initial funcs"
}

func (w *worker) unmarshalFuncs() (Funcs, error) {
	var data msgp.Raw
	select {
	case data = <-w.out:
	case err := <-w.err:
		<-w.waited
		return nil, &NoFuncsError{err}
	case <-w.waited:
		return nil, &NoFuncsError{w.waitErr}
	}

	var funcs Funcs
	if _, err := funcs.UnmarshalMsg(data); err != nil {
		return nil, &NoFuncsError{err}
	}

	return funcs, nil
}

// ErrProcessExited is returned if a Call is made to a worker that has a
// terminated subprocess.
var ErrProcessExited = errors.New("worker: process has exited")

// UnexpectedExitError is returned when the subprocess is unexpeced exited.
type UnexpectedExitError struct {
	// Err is the underlying error why the process has exited.
	Err error
}

func (e *UnexpectedExitError) Error() string {
	return "worker: unexpected exit"
}

// Call executes function call c in the worker subprocess.
// Value crashed is true if the subprocess is crashed during the call.
func (w *worker) Call(c *swerker.Call) (data msgp.Raw, crashed bool, err error) {
	if w.exited() {
		return nil, false, ErrProcessExited
	}

	data, err = c.MarshalMsg(nil)
	if err != nil {
		return nil, false, err
	}

	if _, err = w.in.Write(data); err != nil {
		return nil, false, err
	}

	select {
	case data = <-w.out:
	case err := <-w.err:
		<-w.waited
		return nil, true, err
	case <-w.waited:
		return nil, true, &UnexpectedExitError{w.waitErr}
	}

	if msgp.NextType(data) == msgp.MapType {
		var em ErrorMap

		_, err := em.UnmarshalMsg(data)
		if err == nil {
			err = em.error()
		}

		return nil, false, err
	}

	return data, false, nil
}
