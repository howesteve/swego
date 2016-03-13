// Package swecgo embeds the Swiss Ephemeris library using cgo.
package swecgo

import "github.com/dwlnetnl/swego"

// Library represents the Swiss Ephemeris C library.
type Library struct{}

// Open returns a handle to the Swiss Ephemeris. If ephepath is empty, the
// default (SE_EPHE_PATH) is used.
func Open(ephepath string) *Library {
	if ephepath == "" {
		ephepath = DefaultPath
	}

	withLock(func() {
		setEphePath(ephepath)
	})

	return &Library{}
}

// UsesTLS returns true if Thread Local Storage is enabled and used.
func (l *Library) UsesTLS() bool { return supportsTLS() }

// Close closes the library.
func (l *Library) Close() error {
	withLock(func() {
		close()
	})

	return nil
}

// Version implements swego.Interface.
func (l *Library) Version() string { return Version }

// Calc implements swego.Interface.
func (l *Library) Calc(et float64, pl int, fl swego.CalcFlags) (xx [6]float64, cfl int, err error) {
	withLock(func() {
		setFlagState(fl)
		xx, cfl, err = calc(et, pl, fl.Flags)
	})

	return
}

// CalcUT implements swego.Interface.
func (l *Library) CalcUT(ut float64, pl int, fl swego.CalcFlags) (xx [6]float64, cfl int, err error) {
	withLock(func() {
		setFlagState(fl)
		xx, cfl, err = calcUT(ut, pl, fl.Flags)
	})

	return
}

func setFlagState(fl swego.CalcFlags) {
	if (fl.Flags & flgTopo) == flgTopo {
		setTopo(fl.TopoLoc.Long, fl.TopoLoc.Lat, fl.TopoLoc.Alt)
	}

	if (fl.Flags & flgSidereal) == flgSidereal {
		setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	}
}

// PlanetName implements swego.Interface.
func (l *Library) PlanetName(pl int) (name string) {
	withLock(func() {
		name = planetName(pl)
	})

	return
}
