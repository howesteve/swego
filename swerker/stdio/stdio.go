// Package stdio implements a dispatcher that runs multiple swerker-stdio
// processes.
package stdio

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/astrotools/swego/swerker"
	"github.com/astrotools/swego/swerker/stdio/internal/worker"

	"github.com/tinylib/msgp/msgp"
)

// Dispatcher runs a set of swerker-stdio worker processes.
type Dispatcher struct {
	procs     int
	path      string
	data      string
	workers   []worker.Worker
	workersMu sync.RWMutex // protects workers
	queue     chan task
	crashed   chan worker.Worker
	workDone  chan struct{}
	closed    chan struct{}
	onNewErr  func(error)
	onExitErr func(error)
	funcs     worker.FuncsMap
	lastIdx   uint8
}

type task struct {
	call   *swerker.Call
	result chan result // TODO: benchmark impact of pooling
}

type result struct {
	data msgp.Raw
	err  error
}

// An Option configures an optional Dispatcher parameter.
type Option func(*Dispatcher)

// NumWorkers configures the number of processes are started by Dispatcher.
// If num is 0, the number of logical processors usable by the current process
// is used.
func NumWorkers(num int) Option {
	return func(d *Dispatcher) {
		d.procs = num
	}
}

// OnNewError configures a Dispatcher to call fn when a worker could not be
// restarted.
func OnNewError(fn func(err error)) Option {
	return func(d *Dispatcher) {
		d.onNewErr = fn
	}
}

// OnExitError configures a Dispatcher to call fn with the exit error of a
// crashed worker.
func OnExitError(fn func(err error)) Option {
	return func(d *Dispatcher) {
		d.onExitErr = fn
	}
}

var newWorker = worker.New // for testing

// New returns a Dispatcher that interfaces via swerker-stdio with the
// Swiss Ephemeris. As it takes the file system path to the binary and the
// number of instances of the program as arguments. By default the number of
// logical processors usable by the current process is used.
func New(path string, opts ...Option) (d *Dispatcher, err error) {
	d = &Dispatcher{
		path:     path,
		queue:    make(chan task),
		crashed:  make(chan worker.Worker),
		workDone: make(chan struct{}),
		closed:   make(chan struct{}),
	}

	for _, opt := range opts {
		opt(d)
	}

	if d.procs == 0 {
		d.procs = runtime.NumCPU()
	}

	defer func() {
		if err != nil {
			d.Close()
		}
	}()

	d.workers = make([]worker.Worker, d.procs)
	for i := 0; i < d.procs; i++ {
		w, funcs, err := d.newWorker()
		if err != nil {
			return nil, err
		}

		if i == 0 {
			d.funcs = funcs.FuncsMap()
			d.lastIdx = funcs.LastIdx()
		}

		d.workers[i] = w
	}

	go d.restartWorkers()
	return d, nil
}

func (d *Dispatcher) newWorker() (worker.Worker, worker.Funcs, error) {
	w, funcs, err := newWorker(d.path)
	if err != nil {
		return nil, nil, err
	}

	if d.data != "" {
		if idx, ok := d.IndexForName("swe_set_ephe_path"); ok {
			var args []byte
			args = msgp.AppendArrayHeader(args, 1)
			args = msgp.AppendString(args, d.data)
			w.Call(&swerker.Call{Func: idx, Args: args})
		}
	}

	go d.runWorker(w)
	return w, funcs, nil
}

func (d *Dispatcher) runWorker(w worker.Worker) {
	for t := range d.queue {
		d.workersMu.RLock()

		data, crashed, err := w.Call(t.call)
		t.result <- result{data, err}
		if crashed {
			err := w.Exit()
			if err != nil && d.onExitErr != nil {
				d.onExitErr(err)
			}

			d.crashed <- w
			d.workersMu.RUnlock()
			return
		}

		d.workersMu.RUnlock()
	}

	if idx, ok := d.IndexForName("swe_close"); ok {
		w.Call(&swerker.Call{Func: idx})
	}

	w.Exit()
}

func (d *Dispatcher) restartWorkers() {
	for {
		select {
		case cw := <-d.crashed:
			d.workersMu.Lock()

			for i, w := range d.workers {
				if w == cw {
					w, _, err := d.newWorker()
					if err != nil && d.onNewErr != nil {
						d.onNewErr(err)
					} else {
						d.workers[i] = w
					}

					break
				}
			}

			d.workersMu.Unlock()
		case <-d.workDone:
			d.closed <- struct{}{}
			return
		}
	}
}

// Close terminates and cleans the worker processes.
func (d *Dispatcher) Close() error {
	close(d.queue)
	close(d.workDone)
	<-d.closed
	return nil
}

// Path returns the path of the swerker-stdio binary used by dispatcher d.
func (d *Dispatcher) Path() string { return d.path }

// DataPath configures one or more ephemeris data paths. These are passed to
// each worker by calling swe_set_ephe_path. The paths are combined to a list
// of separated paths.
func DataPath(paths ...string) Option {
	return func(d *Dispatcher) {
		d.data = combineDataPaths(paths)
	}
}

const sep = string(filepath.Separator)
const lsep = string(filepath.ListSeparator)

func combineDataPaths(paths []string) string {
	var s []string

	for _, path := range paths {
		if !strings.HasSuffix(path, sep) {
			path += sep
		}

		s = append(s, path)
	}

	return strings.Join(s, lsep)
}

// DataPath returns the list of ephemeris data paths send to the workers.
func (d *Dispatcher) DataPath() string { return d.data }

// DataPaths returns a slice of ephemeris data paths send to the workers.
func (d *Dispatcher) DataPaths() (s []string) {
	paths := strings.Split(d.DataPath(), lsep)
	for _, path := range paths {
		if strings.HasSuffix(path, sep) {
			path = path[:len(path)-len(sep)]
		}

		s = append(s, path)
	}

	return s
}

// IndexForName implements swerker.Dispatcher interface.
func (d *Dispatcher) IndexForName(name string) (uint8, bool) {
	idx, ok := d.funcs[name]
	return idx, ok
}

// UnimplementedError is returned if the requested function is not implemented
// by the worker process.
type UnimplementedError struct {
	Func uint8
}

func (e *UnimplementedError) Error() string {
	return fmt.Sprintf("stdio: unimplemented function %d", e.Func)
}

// Dispatch implements swerker.Dispatcher interface.
func (d *Dispatcher) Dispatch(c *swerker.Call) (msgp.Raw, error) {
	if c.Func > d.lastIdx {
		return nil, &UnimplementedError{c.Func}
	}

	t := task{c, make(chan result)}
	d.queue <- t
	r := <-t.result
	return r.data, r.err
}

// Version returns the Swiss Ephemeris version linked by the swerker-stdio
// binary.
func Version(path string) (v string, err error) {
	w, funcs, err := newWorker(path)
	if err != nil {
		return "", err
	}

	// terminate worker subprocess on function return
	defer func() {
		exitErr := w.Exit()
		if err == nil && exitErr != nil {
			err = exitErr
		}
	}()

	const name = "swe_version"
	fn, ok := funcs.Lookup(name)
	if !ok {
		return "", fmt.Errorf("stdio: function %q not found", name)
	}

	data, crashed, err := w.Call(&swerker.Call{Func: fn})
	if crashed {
		return "", err
	}

	size, data, err := msgp.ReadArrayHeaderBytes(data)
	if err != nil {
		return "", errors.New("stdio: unexpected return value type")
	}

	if size != 1 {
		return "", errors.New("stdio: unexpected return value length")
	}

	v, _, err = msgp.ReadStringBytes(data)
	return v, err
}
