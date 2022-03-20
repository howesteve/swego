//go:build (linux && cgo) || (darwin && cgo)
// +build linux,cgo darwin,cgo

package swecgo

import "github.com/howesteve/swego"

// acquire locks the wrapper for exclusive library access.
// release unlocks the wrapper from exclusive library access.

var _ Library = (*wrapper)(nil) // assert interface

func (w *wrapper) Version() (string, error) {
	return Version, nil
}

func (w *wrapper) SetPath(ephepath string) {
	w.acquire()
	setEphePath(ephepath)
	w.release()
}

func (w *wrapper) Close() {
	w.acquire()
	closeEphemeris()
	w.release()
}

const resetDeltaT = -1e-10

func setDeltaT(dt *float64) {
	var f float64
	if dt == nil {
		f = resetDeltaT
	} else {
		f = *dt
	}

	setDeltaTUserDef(f)
}

func setCalcFlagsState(fl *swego.CalcFlags) int32 {
	if fl == nil {
		setDeltaT(nil)
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

	if fl.JPLFile != "" {
		fl.JPLFile = swego.FnameDft
	}

	setJPLFile(fl.JPLFile)
	setDeltaT(fl.DeltaT)
	return fl.Flags
}

func (w *wrapper) PlanetName(pl swego.Planet) (string, error) {
	w.acquire()
	name := planetName(pl)
	w.release()
	return name, nil
}

func (w *wrapper) Calc(et float64, pl swego.Planet, fl *swego.CalcFlags) ([]float64, int, error) {
	w.acquire()
	flags := setCalcFlagsState(fl)
	xx, cfl, err := calc(et, pl, flags)
	w.release()
	return xx, cfl, err
}

func (w *wrapper) CalcUT(ut float64, pl swego.Planet, fl *swego.CalcFlags) ([]float64, int, error) {
	w.acquire()
	flags := setCalcFlagsState(fl)
	xx, cfl, err := calcUT(ut, pl, flags)
	w.release()
	return xx, cfl, err
}

func (w *wrapper) NodAps(et float64, pl swego.Planet, fl *swego.CalcFlags, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	w.acquire()
	flags := setCalcFlagsState(fl)
	nasc, ndsc, peri, aphe, err = nodAps(et, pl, flags, m)
	w.release()
	return
}

func (w *wrapper) NodApsUT(ut float64, pl swego.Planet, fl *swego.CalcFlags, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	w.acquire()
	flags := setCalcFlagsState(fl)
	nasc, ndsc, peri, aphe, err = nodApsUT(ut, pl, flags, m)
	w.release()
	return
}

func (w *wrapper) GetAyanamsaEx(et float64, fl *swego.AyanamsaExFlags) (float64, error) {
	w.acquire()
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	f, err := getAyanamsaEx(et, fl.Flags)
	w.release()
	return f, err
}

func (w *wrapper) GetAyanamsaExUT(ut float64, fl *swego.AyanamsaExFlags) (float64, error) {
	w.acquire()
	setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
	f, err := getAyanamsaExUT(ut, fl.Flags)
	w.release()
	return f, err
}

func (w *wrapper) GetAyanamsaName(ayan swego.Ayanamsa) (string, error) {
	w.acquire()
	name := getAyanamsaName(ayan)
	w.release()
	return name, nil
}

func (w *wrapper) JulDay(y, m, d int, h float64, ct swego.CalType) (float64, error) {
	jd := julDay(y, m, d, h, int(ct))
	return jd, nil
}

func (w *wrapper) RevJul(jd float64, ct swego.CalType) (y, m, d int, h float64, err error) {
	y, m, d, h = revJul(jd, int(ct))
	return y, m, d, h, nil
}

func (w *wrapper) UTCToJD(y, m, d, h, i int, s float64, fl *swego.DateConvertFlags) (et, ut float64, err error) {
	w.acquire()
	setDeltaT(fl.DeltaT)
	et, ut, err = utcToJD(y, m, d, h, i, s, int(fl.Calendar))
	w.release()
	return
}

func (w *wrapper) JdETToUTC(et float64, fl *swego.DateConvertFlags) (y, m, d, h, i int, s float64, err error) {
	w.acquire()
	setDeltaT(fl.DeltaT)
	y, m, d, h, i, s = jdETToUTC(et, int(fl.Calendar))
	w.release()
	return y, m, d, h, i, s, nil
}

func (w *wrapper) JdUT1ToUTC(ut1 float64, fl *swego.DateConvertFlags) (y, m, d, h, i int, s float64, err error) {
	w.acquire()
	setDeltaT(fl.DeltaT)
	y, m, d, h, i, s = jdUT1ToUTC(ut1, int(fl.Calendar))
	w.release()
	return y, m, d, h, i, s, nil
}

func (w *wrapper) HousesEx(ut float64, fl *swego.HousesExFlags, geolat, geolon float64, hsys swego.HSys) ([]float64, []float64, error) {
	w.acquire()
	var flags int32
	if fl != nil {
		flags = fl.Flags
		if (flags & flgSidereal) == flgSidereal {
			setSidMode(fl.SidMode.Mode, fl.SidMode.T0, fl.SidMode.AyanT0)
		}

		setDeltaT(fl.DeltaT)
	} else {
		setDeltaT(nil)
	}

	cusps, ascmc, err := housesEx(ut, flags, geolat, geolon, hsys)
	w.release()
	return cusps, ascmc, err
}

func (w *wrapper) HousesARMC(armc, geolat, eps float64, hsys swego.HSys) ([]float64, []float64, error) {
	w.acquire()
	cusps, ascmc, err := housesARMC(armc, geolat, eps, hsys)
	w.release()
	return cusps, ascmc, err
}

func (w *wrapper) HousePos(armc, geolat, eps float64, hsys swego.HSys, pllng, pllat float64) (float64, error) {
	w.acquire()
	pos, err := housePos(armc, geolat, eps, hsys, pllng, pllat)
	w.release()
	return pos, err
}

func (w *wrapper) HouseName(hsys swego.HSys) (string, error) {
	w.acquire()
	name := houseName(hsys)
	w.release()
	return name, nil
}

func (w *wrapper) DeltaTEx(jd float64, eph swego.Ephemeris) (float64, error) {
	w.acquire()
	dt, err := deltaTEx(jd, int32(eph))
	w.release()
	return dt, err
}

func setTimeEquDeltaT(fl *swego.TimeEquFlags) {
	if fl == nil {
		setDeltaT(nil)
	} else {
		setDeltaT(fl.DeltaT)
	}
}

func (w *wrapper) TimeEqu(jd float64, fl *swego.TimeEquFlags) (float64, error) {
	w.acquire()
	setTimeEquDeltaT(fl)
	f, err := timeEqu(jd)
	w.release()
	return f, err
}

func (w *wrapper) LMTToLAT(lmt, geolon float64, fl *swego.TimeEquFlags) (float64, error) {
	w.acquire()
	setTimeEquDeltaT(fl)
	lat, err := lmtToLAT(lmt, geolon)
	w.release()
	return lat, err
}

func (w *wrapper) LATToLMT(lat, geolon float64, fl *swego.TimeEquFlags) (float64, error) {
	w.acquire()
	setTimeEquDeltaT(fl)
	lmt, err := latToLMT(lat, geolon)
	w.release()
	return lmt, err
}

func setSidTimeDeltaT(fl *swego.SidTimeFlags) {
	if fl == nil {
		setDeltaT(nil)
	} else {
		setDeltaT(fl.DeltaT)
	}
}

func (w *wrapper) SidTime0(ut, eps, nut float64, fl *swego.SidTimeFlags) (float64, error) {
	w.acquire()
	setSidTimeDeltaT(fl)
	f := sidTime0(ut, eps, nut)
	w.release()
	return f, nil
}

func (w *wrapper) SidTime(ut float64, fl *swego.SidTimeFlags) (float64, error) {
	w.acquire()
	setSidTimeDeltaT(fl)
	f := sidTime(ut)
	w.release()
	return f, nil
}
