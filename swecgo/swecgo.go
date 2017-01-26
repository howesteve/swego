// +build linux,cgo darwin,cgo

// Package swecgo embeds the Swiss Ephemeris library using cgo.
package swecgo

import (
	"sync"

	"github.com/astrotools/swego"
)

var invoker swego.Invoker

func init() {
	if supportsTLS() {
		panic("swecgo: Thread Local Storage (TLS) is not supported")
	}

	invoker = new(wrapper)
}

type wrapper struct {
	mu sync.Mutex
}

// Invoke implements interface swego.Invoker.
func (w *wrapper) Invoke(fn func(swe swego.Interface)) error {
	w.mu.Lock()
	fn(w)
	w.mu.Unlock()
	return nil
}

// Invoker returns an initialized invoker.
func Invoker() swego.Invoker { return invoker }

// Call calls fn within an initialized execution context.
func Call(fn func(swe swego.Interface)) { Invoker().Invoke(fn) }

// Version implements swego.Interface.
func (*wrapper) Version() string { return Version }

// SetPath implements swego.Interface.
func (*wrapper) SetPath(ephepath string) { setEphePath(ephepath) }

// Close implements swego.Interface.
func (*wrapper) Close() { closeEphemeris() }

func setCalcFlagsState(fl *swego.CalcFlags) int32 {
	if fl == nil {
		return 0
	}

	if (fl.Flags & flgTopo) == flgTopo {
		var lng, lat, alt float64

		if fl.TopoLoc != nil {
			lng = fl.TopoLoc.Long
			lat = fl.TopoLoc.Lat
			alt = fl.TopoLoc.Alt
		}

		setTopo(lng, lat, alt)
	}

	if (fl.Flags & flgSidereal) == flgSidereal {
		var mode swego.Ayanamsa
		var t0, ayanT0 float64

		if fl.SidMode != nil {
			mode = fl.SidMode.Mode
			t0 = fl.SidMode.T0
			ayanT0 = fl.SidMode.AyanT0
		}

		setSidMode(mode, t0, ayanT0)
	}

	if fl.FileNameJPL != "" {
		fl.FileNameJPL = swego.FnameDft
	}

	setFileNameJPL(fl.FileNameJPL)
	return fl.Flags
}

// PlanetName implements swego.Interface.
func (*wrapper) PlanetName(pl swego.Planet) string { return planetName(pl) }

// Calc implements swego.Interface.
func (*wrapper) Calc(et float64, pl swego.Planet, fl *swego.CalcFlags) ([]float64, int, error) {
	flags := setCalcFlagsState(fl)
	return calc(et, pl, flags)
}

// CalcUT implements swego.Interface.
func (*wrapper) CalcUT(ut float64, pl swego.Planet, fl *swego.CalcFlags) ([]float64, int, error) {
	flags := setCalcFlagsState(fl)
	return calcUT(ut, pl, flags)
}

// NodAps implements swego.Interface.
func (*wrapper) NodAps(et float64, pl swego.Planet, fl *swego.CalcFlags, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	flags := setCalcFlagsState(fl)
	return nodAps(et, pl, flags, m)
}

// NodApsUT implements swego.Interface.
func (*wrapper) NodApsUT(ut float64, pl swego.Planet, fl *swego.CalcFlags, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	flags := setCalcFlagsState(fl)
	return nodApsUT(ut, pl, flags, m)
}

// GetAyanamsa implements swego.Interface.
func (*wrapper) GetAyanamsa(et float64, sidmode *swego.SidMode) float64 {
	setSidMode(sidmode.Mode, sidmode.T0, sidmode.AyanT0)
	return getAyanamsa(et)
}

// GetAyanamsaUT implements swego.Interface.
func (*wrapper) GetAyanamsaUT(ut float64, sidmode *swego.SidMode) float64 {
	setSidMode(sidmode.Mode, sidmode.T0, sidmode.AyanT0)
	return getAyanamsaUT(ut)
}

// GetAyanamsaEx implements swego.Interface.
func (*wrapper) GetAyanamsaEx(et float64, fl *swego.AyanamsaExFlags) (float64, error) {
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	return getAyanamsaEx(et, fl.Flags)
}

// GetAyanamsaExUT implements swego.Interface.
func (*wrapper) GetAyanamsaExUT(ut float64, fl *swego.AyanamsaExFlags) (float64, error) {
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	return getAyanamsaExUT(ut, fl.Flags)
}

// GetAyanamsaName implements swego.Interface.
func (*wrapper) GetAyanamsaName(ayan swego.Ayanamsa) string {
	return getAyanamsaName(ayan)
}

// JulDay implements swego.Interface.
func (*wrapper) JulDay(y, m, d int, h float64, ct swego.CalType) float64 {
	return julDay(y, m, d, h, int(ct))
}

// RevJul implements swego.Interface.
func (*wrapper) RevJul(jd float64, ct swego.CalType) (y, m, d int, h float64) {
	return revJul(jd, int(ct))
}

// UTCToJD implements swego.Interface.
func (*wrapper) UTCToJD(y, m, d, h, i int, s float64, ct swego.CalType) (et, ut float64, err error) {
	return utcToJD(y, m, d, h, i, s, int(ct))
}

// JdETToUTC implements swego.Interface.
func (*wrapper) JdETToUTC(et float64, ct swego.CalType) (y, m, d, h, i int, s float64) {
	return jdETToUTC(et, int(ct))
}

// JdUT1ToUTC implements swego.Interface.
func (*wrapper) JdUT1ToUTC(ut1 float64, ct swego.CalType) (y, m, d, h, i int, s float64) {
	return jdUT1ToUTC(ut1, int(ct))
}

// Houses implements swego.Interface.
func (*wrapper) Houses(ut, geolat, geolon float64, hsys swego.HSys) ([]float64, []float64, error) {
	return houses(ut, geolat, geolon, hsys)
}

// HousesEx implements swego.Interface.
func (*wrapper) HousesEx(ut float64, fl *swego.HousesExFlags, geolat, geolon float64, hsys swego.HSys) ([]float64, []float64, error) {
	var flags int32
	if fl != nil {
		flags = fl.Flags
		if (flags & flgSidereal) == flgSidereal {
			setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
		}
	}

	return housesEx(ut, flags, geolat, geolon, hsys)
}

// HousesARMC implements swego.Interface.
func (*wrapper) HousesARMC(armc, geolat, eps float64, hsys swego.HSys) ([]float64, []float64, error) {
	return housesARMC(armc, geolat, eps, hsys)
}

// HousePos implements swego.Interface.
func (*wrapper) HousePos(armc, geolat, eps float64, hsys swego.HSys, pllng, pllat float64) (float64, error) {
	return housePos(armc, geolat, eps, hsys, pllng, pllat)
}

// HouseName implements swego.Interface.
func (*wrapper) HouseName(hsys swego.HSys) string {
	return houseName(hsys)
}

// DeltaT implements swego.Interface.
func (*wrapper) DeltaT(jd float64) float64 { return deltaT(jd) }

// DeltaTEx implements swego.Interface.
func (*wrapper) DeltaTEx(jd float64, eph swego.Ephemeris) (float64, error) {
	return deltaTEx(jd, int32(eph))
}

// SetDeltaTUserDef implements swego.Interface.
func (*wrapper) SetDeltaTUserDef(v float64) { setDeltaTUserDef(v) }

// TimeEqu implements swego.Interface.
func (*wrapper) TimeEqu(jd float64) (float64, error) { return timeEqu(jd) }

// LMTToLAT implements swego.Interface.
func (*wrapper) LMTToLAT(jdLMT, geolon float64) (float64, error) {
	return lmtToLAT(jdLMT, geolon)
}

// LATToLMT implements swego.Interface.
func (*wrapper) LATToLMT(jdLAT, geolon float64) (float64, error) {
	return latToLMT(jdLAT, geolon)
}

// SidTime0 implements swego.Interface.
func (*wrapper) SidTime0(ut, eps, nut float64) float64 {
	return sidTime0(ut, eps, nut)
}

// SidTime implements swego.Interface.
func (*wrapper) SidTime(ut float64) float64 { return sidTime(ut) }
