// +build linux,cgo darwin,cgo

// Package swecgo embeds the Swiss Ephemeris library using cgo.
package swecgo

import (
	"runtime"
	"sync"

	"github.com/astrotools/swego"
)

// Call calls fn within an initialized execution context. The initialization of
// this context is done by calling initFn. If initFn is nil, the default data
// path is set to DefaultPath. For more information see the Programmer's
// Documentation about swe_set_ephe_path.
func Call(initFn func(swego.Interface), fn func(swego.Interface)) {
	inv := NewInvoker(initFn)
	inv.Invoke(fn)
}

// NewInvoker initializes an execution context and returns it.
// If initFn is nil, the default data path is set to DefaultPath. For more
// information see the Programmer's Documentation about swe_set_ephe_path.
func NewInvoker(initFn func(swego.Interface)) swego.Invoker {
	if initFn == nil {
		initFn = func(swe swego.Interface) {
			swe.SetPath(DefaultPath)
		}
	}

	if supportsTLS() {
		panic("swecgo: Thread Local Storage (TLS) is not supported")
		// inv := tlsInvoker{}
		// inv.Invoke(initFn)
		// return inv
	}

	gInvoker.once.Do(func() {
		gInvoker.inv = newMuInvoker()
		gInvoker.inv.Invoke(initFn)
	})

	return gInvoker.inv
}

var gInvoker struct {
	inv  swego.Invoker
	once sync.Once
}

type muInvoker struct {
	mu sync.Mutex
}

func newMuInvoker() *muInvoker { return &muInvoker{} }

// Invoke implements interface swego.Invoker.
func (inv *muInvoker) Invoke(fn func(swego.Interface)) error {
	inv.mu.Lock()
	fn(gWrapper)
	inv.mu.Unlock()
	return nil
}

type tlsInvoker struct{}

func (tlsInvoker) Invoke(fn func(swego.Interface)) error {
	runtime.LockOSThread()
	fn(gWrapper)
	runtime.UnlockOSThread()
	return nil
}

type chanInvoker struct {
	fnCh chan func()
}

func newChanInvoker() *chanInvoker {
	inv := &chanInvoker{fnCh: make(chan func())}
	go inv.runLoop()
	return inv
}

func (inv *chanInvoker) runLoop() {
	for fn := range inv.fnCh {
		fn()
	}
}

// Invoke implements interface swego.Invoker.
func (inv *chanInvoker) Invoke(fn func(swego.Interface)) error {
	var wg sync.WaitGroup

	wg.Add(1)
	inv.fnCh <- func() {
		fn(gWrapper)
		wg.Done()
	}

	wg.Wait()
	return nil
}
