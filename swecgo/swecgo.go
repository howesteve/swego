// +build linux,cgo darwin,cgo

// Package swecgo embeds the Swiss Ephemeris library using cgo.
package swecgo

import (
	"sync"

	"github.com/dwlnetnl/swego"
)

// Call calls fn within an initialized execution context. The initialization of
// this context is done by calling initFn. If initFn is nil, the default data
// path is set to DefaultPath. For more information see the Programmer's
// Documentation about swe_set_ephe_path.
//
// In non-TLS-mode only a single fn can be executed at any point time. This is
// achieved by sending fn over a channel to a separate goroutine that executes
// all closures it receives and Call blocks waiting until fn is done executing.
//
// TLS-mode (Thread Local Storage) is currently not implemented.
func Call(initFn func(swego.Interface), fn func(swego.Interface)) {
	initGWrapper(initFn)
	gWrapper.execute(fn)
	// In TLS-mode fn would be called on a single thread in a pool of locked OS
	// threads. For more information about this see runtime.LockOSThread.
}

// NewInvoker initializes an execution context and returns it.
// If initFn is nil, the default data path is set to DefaultPath. For more
// information see the Programmer's Documentation about swe_set_ephe_path.
func NewInvoker(initFn func(swego.Interface)) swego.Invoker {
	initGWrapper(initFn)
	return gWrapper
}

// Invoke implements interface swego.Invoker.
func (w *wrapper) Invoke(fn func(swego.Interface)) error {
	w.execute(fn)
	return nil
}

// ----------

var gWrapper *wrapper
var gWrapperOnce sync.Once

func initGWrapper(initFn func(swego.Interface)) {
	gWrapperOnce.Do(func() {
		if supportsTLS() {
			panic("Swiss Ephemeris library with Thread Local Storage enabled " +
				"is not supported")
		}

		gWrapper = newWrapper()

		if initFn == nil {
			initFn = func(swe swego.Interface) {
				swe.SetPath(DefaultPath)
			}
		}

		gWrapper.execute(initFn)
	})
}

// ----------

type wrapper struct {
	fnCh chan func()
}

func newWrapper() *wrapper {
	w := &wrapper{fnCh: make(chan func())}
	go w.runLoop()
	return w
}

func (w *wrapper) runLoop() {
	// Run loop is a separate goroutine.
	// The OS thread is always locked during a cgo call.
	// This code might change if Thread Local Storage (TLS) is used.

	for fn := range w.fnCh {
		fn()
	}
}

func (w *wrapper) execute(fn func(swego.Interface)) {
	var wg sync.WaitGroup

	wg.Add(1)
	w.fnCh <- func() {
		fn(w)
		wg.Done()
	}

	wg.Wait()
}

// ----------

var _ swego.Interface = (*wrapper)(nil) // assert interface

// Version implements swego.Interface.
func (w *wrapper) Version() string { return Version }

// GetLibraryPath implements swego.Interface.
func (w *wrapper) GetLibraryPath() string { return getLibraryPath() }

// SetPath implements swego.Interface.
func (w *wrapper) SetPath(ephepath string) { setEphePath(ephepath) }

// Close implements swego.Interface.
func (w *wrapper) Close() {
	closeEphemeris()
	// w.fnCh should be closed, but w is gWrapper and is global we must not.
}

func setCalcFlagsState(fl *swego.CalcFlags) {
	if (fl.Flags & flgTopo) == flgTopo {
		setTopo(fl.TopoLoc.Long, fl.TopoLoc.Lat, fl.TopoLoc.Alt)
	}

	if (fl.Flags & flgSidereal) == flgSidereal {
		setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	}

	if fl.FileNameJPL != "" {
		fl.FileNameJPL = swego.FnameDft
	}

	setFileNameJPL(fl.FileNameJPL)
}

// PlanetName implements swego.Interface.
func (w *wrapper) PlanetName(pl swego.Planet) string { return planetName(pl) }

// Calc implements swego.Interface.
func (w *wrapper) Calc(et float64, pl swego.Planet, fl *swego.CalcFlags) ([]float64, int, error) {
	setCalcFlagsState(fl)
	return calc(et, pl, fl.Flags)
}

// CalcUT implements swego.Interface.
func (w *wrapper) CalcUT(ut float64, pl swego.Planet, fl *swego.CalcFlags) ([]float64, int, error) {
	setCalcFlagsState(fl)
	return calcUT(ut, pl, fl.Flags)
}

// NodAps implements swego.Interface.
func (w *wrapper) NodAps(et float64, pl swego.Planet, fl *swego.CalcFlags, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	setCalcFlagsState(fl)
	return nodAps(et, pl, fl.Flags, m)
}

// NodApsUT implements swego.Interface.
func (w *wrapper) NodApsUT(ut float64, pl swego.Planet, fl *swego.CalcFlags, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	setCalcFlagsState(fl)
	return nodApsUT(ut, pl, fl.Flags, m)
}

// GetAyanamsa implements swego.Interface.
func (w *wrapper) GetAyanamsa(et float64, sidmode *swego.SidMode) float64 {
	setSidMode(sidmode.Mode, sidmode.T0, sidmode.AyanT0)
	return getAyanamsa(et)
}

// GetAyanamsaUT implements swego.Interface.
func (w *wrapper) GetAyanamsaUT(ut float64, sidmode *swego.SidMode) float64 {
	setSidMode(sidmode.Mode, sidmode.T0, sidmode.AyanT0)
	return getAyanamsaUT(ut)
}

// GetAyanamsaEx implements swego.Interface.
func (w *wrapper) GetAyanamsaEx(et float64, fl *swego.AyanamsaExFlags) (float64, error) {
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	return getAyanamsaEx(et, fl.Flags)
}

// GetAyanamsaExUT implements swego.Interface.
func (w *wrapper) GetAyanamsaExUT(ut float64, fl *swego.AyanamsaExFlags) (float64, error) {
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	return getAyanamsaExUT(ut, fl.Flags)
}

// GetAyanamsaName implements swego.Interface.
func (w *wrapper) GetAyanamsaName(ayan swego.Ayanamsa) string {
	return getAyanamsaName(ayan)
}

// JulDay implements swego.Interface.
func (w *wrapper) JulDay(y, m, d int, h float64, ct swego.CalType) float64 {
	return julDay(y, m, d, h, int(ct))
}

// RevJul implements swego.Interface.
func (w *wrapper) RevJul(jd float64, ct swego.CalType) (y, m, d int, h float64) {
	return revJul(jd, int(ct))
}

// UTCToJD implements swego.Interface.
func (w *wrapper) UTCToJD(y, m, d, h, i int, s float64, ct swego.CalType) (et, ut float64, err error) {
	return utcToJD(y, m, d, h, i, s, int(ct))
}

// JdETToUTC implements swego.Interface.
func (w *wrapper) JdETToUTC(et float64, ct swego.CalType) (y, m, d, h, i int, s float64) {
	return jdETToUTC(et, int(ct))
}

// JdUT1ToUTC implements swego.Interface.
func (w *wrapper) JdUT1ToUTC(ut1 float64, ct swego.CalType) (y, m, d, h, i int, s float64) {
	return jdUT1ToUTC(ut1, int(ct))
}

// Houses implements swego.Interface.
func (w *wrapper) Houses(ut, geolat, geolon float64, hsys swego.HSys) ([]float64, []float64, error) {
	return houses(ut, geolat, geolon, hsys)
}

// HousesEx implements swego.Interface.
func (w *wrapper) HousesEx(ut float64, fl *swego.HousesExFlags, geolat, geolon float64, hsys swego.HSys) ([]float64, []float64, error) {
	if (fl.Flags & flgSidereal) == flgSidereal {
		setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	}

	return housesEx(ut, fl.Flags, geolat, geolon, hsys)
}

// HousesARMC implements swego.Interface.
func (w *wrapper) HousesARMC(armc, geolat, eps float64, hsys swego.HSys) ([]float64, []float64, error) {
	return housesARMC(armc, geolat, eps, hsys)
}

// HousePos implements swego.Interface.
func (w *wrapper) HousePos(armc, geolat, eps float64, hsys swego.HSys, pllng, pllat float64) (float64, error) {
	return housePos(armc, geolat, eps, hsys, pllng, pllat)
}

// HouseName implements swego.Interface.
func (w *wrapper) HouseName(hsys swego.HSys) string {
	return houseName(hsys)
}

// DeltaT implements swego.Interface.
func (w *wrapper) DeltaT(jd float64) float64 { return deltaT(jd) }

// DeltaTEx implements swego.Interface.
func (w *wrapper) DeltaTEx(jd float64, eph swego.Ephemeris) (float64, error) {
	return deltaTEx(jd, int32(eph))
}

// SetDeltaTUserDef implements swego.Interface.
func (w *wrapper) SetDeltaTUserDef(v float64) { setDeltaTUserDef(v) }

// TimeEqu implements swego.Interface.
func (w *wrapper) TimeEqu(jd float64) (float64, error) { return timeEqu(jd) }

// LMTToLAT implements swego.Interface.
func (w *wrapper) LMTToLAT(jdLMT, geolon float64) (float64, error) {
	return lmtToLAT(jdLMT, geolon)
}

// LATToLMT implements swego.Interface.
func (w *wrapper) LATToLMT(jdLAT, geolon float64) (float64, error) {
	return latToLMT(jdLAT, geolon)
}

// SidTime0 implements swego.Interface.
func (w *wrapper) SidTime0(ut, eps, nut float64) float64 {
	return sidTime0(ut, eps, nut)
}

// SidTime implements swego.Interface.
func (w *wrapper) SidTime(ut float64) float64 { return sidTime(ut) }
