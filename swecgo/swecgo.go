// go:build (linux && cgo) || (darwin && cgo)
//go:build (linux && cgo) || (darwin && cgo)
// +build linux,cgo darwin,cgo

// #cgo LDFLAGS: -ldl -lm
// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: m

// Package swecgo embeds the Swiss Ephemeris library using cgo.
package swecgo

import (
	"sync"

	"github.com/howesteve/swego"
)

// Library extends the main library interface by exposing C library
// life-cycle methods.
type Library interface {
	// The following methods will always return nil as error:
	//  Version
	//  PlanetName
	//  GetAyanamsaName
	//  JulDay
	//  RevJul
	//  JdETToUTC
	//  JdUT1ToUTC
	//  HouseName
	//  SidTime
	//  SidTime0
	swego.Interface

	// SetPath opens the ephemeris and sets the data path.
	SetPath(path string)

	// Close closes the Swiss Ephemeris library.
	// The ephemeris can be reopened by calling SetPath.
	Close()

	// used for locking and prevent other interface implementations
	acquire()
	release()
}

// Open initializes the Swiss Ephemeris C library with DefaultPath as
// ephemeris path. The returned object is safe for concurrent use.
func Open() Library { return OpenWithPath(DefaultPath) }

// OpenWithPath initializes the Swiss Ephemeris C library and calls
// swe_set_ephe_path with ephePath as argument afterwards.
// The returned object is safe for concurrent use.
func OpenWithPath(ephePath string) Library {
	swe := Interface()
	swe.SetPath(ephePath)
	return swe
}

// Interface returns an object that calls the Swiss Ephemeris C library.
// The returned object is safe for concurrent use.
func Interface() Library {
	winit.Do(func() {
		checkLibrary()
		wrap = &wrapper{locker: new(sync.Mutex)}
	})

	return wrap
}

var winit sync.Once
var wrap Library

// wrapper interfaces between swego.Interface and the library functions.
// It protect stateful library functions with a mutex. When the wrapper is
// exclusively locked, the mutex is temporary replaced by a no-op lock.
type wrapper struct {
	locker sync.Locker
}

func (w *wrapper) acquire() { w.locker.Lock() }
func (w *wrapper) release() { w.locker.Unlock() }

type unlocked struct{}

func (unlocked) Lock()   {}
func (unlocked) Unlock() {}

var unlockedWrapper = &wrapper{locker: unlocked{}}

type exclLocked struct {
	*wrapper
	locker sync.Locker
}

func (el exclLocked) ExclusiveUnlock() { el.locker.Unlock() }
func (w *wrapper) ExclusiveLock() swego.LockedInterface {
	w.locker.Lock()
	return exclLocked{
		wrapper: unlockedWrapper, // wrapper with no-op lock
		locker:  w.locker,        // the actual wrapper mutex
	}
}

// Locked exclusively locks the library, disable per function locking and
// exposes the locked library to the callback function. Per function locking
// is restored when execution is returned to the caller.
// If either argument is nil, it panics.
func Locked(swe Library, callback func(swe Library)) {
	if swe == nil {
		panic("swe is nil")
	}

	if callback == nil {
		panic("callback is nil")
	}

	w := swe.(*wrapper)
	l := w.ExclusiveLock().(exclLocked)
	callback(l)
	l.ExclusiveUnlock()
}
