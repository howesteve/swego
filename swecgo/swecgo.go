// Package swecgo embeds the Swiss Ephemeris library using cgo.
package swecgo

import (
	"runtime"
	"sync"

	"github.com/dwlnetnl/swego"
)

// Call calls fn within an initialized execution context. The initialization of
// this context is done by calling init. If init is nil, the default data path
// is set, the value of DefaultPath. For more information see the/ Programmer's
// documentation about swe_set_ephe_path.
//
// In non-TLS-mode only a single fn can be executed at a single time. This is
// achieved by a single mutex that is locked before a fn is called and unlocked
// after fn has been executed.
//
// TLS-mode (Thread Local Storage) is currently not implemented.
func Call(init func(swe swego.Interface), fn func(swe swego.Interface)) {
	// In TLS-mode fn would be called on a single thread in a pool of locked OS
	// threads. For more information about this see runtime.LockOSThread.

	gWrapper.once.Do(func() {
		if supportsTLS() {
			panic("Swiss Ephemeris library with Thread Local Storage enabled " +
				"is not supported")
		}

		if init == nil {
			init = func(swe swego.Interface) {
				swe.SetPath(DefaultPath)
			}
		}

		gWrapper.fnCh = make(chan func())
		go gWrapper.runLoop()
		gWrapper.execute(init)
	})

	gWrapper.execute(fn)
}

var gWrapper = new(wrapper)

type wrapper struct {
	once sync.Once
	fnCh chan func()
}

func (w *wrapper) runLoop() {
	runtime.LockOSThread()

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

// Version implements swego.Interface.
func (w *wrapper) Version() string { return Version }

// SetPath implements swego.Interface.
func (w *wrapper) SetPath(ephepath string) { setEphePath(ephepath) }

// Close implements swego.Interface.
func (w *wrapper) Close() { close() }

func setCalcFlagsState(fl swego.CalcFlags) {
	if (fl.Flags & flgTopo) == flgTopo {
		setTopo(fl.TopoLoc.Long, fl.TopoLoc.Lat, fl.TopoLoc.Alt)
	}

	if (fl.Flags & flgSidereal) == flgSidereal {
		setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	}

	if fl.FileNameJPL != "" {
		setFileNameJPL(fl.FileNameJPL)
	}
}

// Calc implements swego.Interface.
func (w *wrapper) Calc(et float64, pl swego.Planet, fl swego.CalcFlags) ([6]float64, int, error) {
	setCalcFlagsState(fl)
	return calc(et, pl, fl.Flags)
}

// CalcUT implements swego.Interface.
func (w *wrapper) CalcUT(ut float64, pl swego.Planet, fl swego.CalcFlags) ([6]float64, int, error) {
	setCalcFlagsState(fl)
	return calcUT(ut, pl, fl.Flags)
}

// PlanetName implements swego.Interface.
func (w *wrapper) PlanetName(pl swego.Planet) string { return planetName(pl) }

// GetAyanamsa implements swego.Interface.
func (w *wrapper) GetAyanamsa(et float64, sidmode swego.SidMode) float64 {
	setSidMode(sidmode.Mode, sidmode.T0, sidmode.AyanT0)
	return getAyanamsa(et)
}

// GetAyanamsaUT implements swego.Interface.
func (w *wrapper) GetAyanamsaUT(ut float64, sidmode swego.SidMode) float64 {
	setSidMode(sidmode.Mode, sidmode.T0, sidmode.AyanT0)
	return getAyanamsaUT(ut)
}

// GetAyanamsaEx implements swego.Interface.
func (w *wrapper) GetAyanamsaEx(et float64, fl swego.AyanamsaExFlags) (float64, error) {
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	return getAyanamsaEx(et, fl.Flags)
}

// GetAyanamsaExUT implements swego.Interface.
func (w *wrapper) GetAyanamsaExUT(ut float64, fl swego.AyanamsaExFlags) (float64, error) {
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	return getAyanamsaExUT(ut, fl.Flags)
}

// GetAyanamsaName implements swego.Interface.
func (w *wrapper) GetAyanamsaName(sidmode int32) string {
	return getAyanamsaName(sidmode)
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
func (w *wrapper) Houses(ut, geolat, geolon float64, hsys swego.HSys) ([]float64, [10]float64, error) {
	return houses(ut, geolat, geolon, hsys)
}

// HousesEx implements swego.Interface.
func (w *wrapper) HousesEx(ut float64, fl swego.HousesExFlags, geolat, geolon float64, hsys swego.HSys) ([]float64, [10]float64, error) {
	if (fl.Flags & flgSidereal) == flgSidereal {
		setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	}

	return housesEx(ut, fl.Flags, geolat, geolon, hsys)
}

// HousesArmc implements swego.Interface.
func (w *wrapper) HousesArmc(armc, geolat, eps float64, hsys swego.HSys) ([]float64, [10]float64, error) {
	return housesArmc(armc, geolat, eps, hsys)
}

// HousePos implements swego.Interface.
func (w *wrapper) HousePos(armc, geolat, eps float64, hsys swego.HSys, pllng, pllat float64) (float64, error) {
	return housePos(armc, geolat, eps, hsys, pllng, pllat)
}

// HouseName implements swego.Interface.
func (w *wrapper) HouseName(hsys swego.HSys) string { return houseName(hsys) }

// DeltaT implements swego.Interface.
func (w *wrapper) DeltaT(jd float64) float64 { return deltaT(jd) }

// DeltaTEx implements swego.Interface.
func (w *wrapper) DeltaTEx(jd float64, fl int32) (float64, error) {
	return deltaTEx(jd, fl)
}

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
