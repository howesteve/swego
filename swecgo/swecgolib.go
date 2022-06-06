//go:build (linux && cgo) || (darwin && cgo)
// +build linux,cgo darwin,cgo

// #cgo LDFLAGS: -ldl -lm
// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: m

package swecgo

import (
	"unsafe"

	"github.com/howesteve/swego"
)

/*
// To aid debugging the Swiss Ephemeris supports tracing. It will write two
// files in the current working directory. Note that the Go tool may change the
// current working directory. The tracing facility has 2 modes:
//
//  TRACE=1 writes swetrace.{c,txt} and appends each session to it.
//  TRACE=2 writes swetrace_<pid>.{c,txt} where pid is the running process id.
//
// To build a version of the Swiss Ephemeris that has tracing has enabled
// uncomment one of the #cgo CFLAGS lines below or overwrite the cgo CFLAGS by
// setting the CGO_CFLAGS environment variable to "-DTRACE=1" or "-DTRACE=2".
//
// The tracing facility is limited to TRACE_COUNT_MAX entries. You can change
// this limit by overwriting this value by uncommenting the line below or using
// the CGO_CFLAGS environment variable ("-DTRACE_COUNT_MAX=value").
//
// When swe_close() is called it may possible to also close the tracing system.
// To enable this uncomment the TRACE_CLOSE=1 line or set the CGO_CFLAGS
// environment variable ("-DTRACE_CLOSE=1"). There is no guarantee that the
// tracing can be reopened after it is closed the first time.
//
// It's very easy to enable the tracing and the trace limit per build using the
// CGO_CFLAGS environment variable. It can be helpful to inspect the Go build
// process pass the -x argument flag.

// #cgo CFLAGS: -DTRACE=1
// #cgo CFLAGS: -DTRACE=2
// #cgo CFLAGS: -DTRACE_COUNT_MAX=10000000
// #cgo CFLAGS: -DTRACE_CLOSE=1

// ----------

// Disable thread local storage in library.
// #cgo CFLAGS: -DTLSOFF=1

// ----------

#cgo CFLAGS: -DTLSOFF=1
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS:-ldl -lm

#include "swephexp.h"
#include "sweph.h"
#include "swex.h"
#include "sweversion.h"

double swecgo_deltat_automatic() {
	return SE_DELTAT_AUTOMATIC;
}
*/
import "C"

func checkLibrary() {
	if C.swex_supports_tls() {
		panic("swecgo: Thread Local Storage (TLS) is not supported")
	}

	if resetDeltaT != C.swecgo_deltat_automatic() {
		panic("swecgo: SE_DELTAT_AUTOMATIC mismatch")
	}
}

// withError calls fn with a pre allocated error variable that can passed to a
// function in the C library. The code block must return true if an error is
// returned from the C call.
func withError(fn func(err *C.char) bool) error {
	var _err [C.AS_MAXCH]C.char

	if fn(&_err[0]) {
		return swego.Error(C.GoString(&_err[0]))
	}

	return nil
}

// Swiss Ephemeris version constants.
const (
	Version      = C.SE_VERSION
	VersionMajor = C.SWEX_VERSION_MAJOR
	VersionMinor = C.SWEX_VERSION_MINOR
	VersionPatch = C.SWEX_VERSION_PATCH
)

// DefaultPath is the default ephemeris path defined by the library.
const DefaultPath = C.SE_EPHE_PATH

const (
	flgTopo     = C.SEFLG_TOPOCTR
	flgSidereal = C.SEFLG_SIDEREAL
)

func setEphePath(path string) {
	_path := C.CString(path)
	C.swe_set_ephe_path(_path)
	C.free(unsafe.Pointer(_path))
}

func setJPLFile(name string) {
	_name := C.CString(name)
	C.swex_set_jpl_file(_name)
	C.free(unsafe.Pointer(_name))
}

func setTopo(lng, lat, alt float64) {
	C.swex_set_topo(C.double(lng), C.double(lat), C.double(alt))
}

func setSidMode(mode swego.Ayanamsa, t0, ayanT0 float64) {
	C.swex_set_sid_mode(C.int32_t(mode), C.double(t0), C.double(ayanT0))
}

func closeEphemeris() {
	C.swe_close()
}

func planetName(pl swego.Planet) string {
	var _name [C.AS_MAXCH]C.char
	C.swe_get_planet_name(C.int(pl), &_name[0])
	return C.GoString(&_name[0])
}

type _calcFunc func(jd C.double, fl C.int32, xx *C.double, err *C.char) C.int32

func _calc(jd float64, fl int32, fn _calcFunc) (_ []float64, cfl int, err error) {
	_jd := C.double(jd)
	_fl := C.int32(fl)

	// Both the float64 and C.double types are defined as an IEEE-754 64-bit
	// floating-point number. This means the representation in memory of a
	// float64 array is equivalent to that of a C.double array. It means it is
	// possible to cast between the two types. In Go land such operation is
	// considered unsafe, hence the use of the unsafe package.
	var xx [6]float64
	_xx := (*C.double)(unsafe.Pointer(&xx[0]))

	err = withError(func(err *C.char) bool {
		cfl = int(fn(_jd, _fl, _xx, err))
		return cfl == C.ERR
	})

	return xx[:], cfl, err
}

func calc(et float64, pl swego.Planet, fl int32) ([]float64, int, error) {
	return _calc(et, fl, func(jd C.double, fl C.int32, xx *C.double, err *C.char) C.int32 {
		return C.swe_calc(jd, C.int(pl), fl, xx, err)
	})
}

func calcUT(ut float64, pl swego.Planet, fl int32) ([]float64, int, error) {
	return _calc(ut, fl, func(jd C.double, fl C.int32, xx *C.double, err *C.char) C.int32 {
		return C.swe_calc_ut(jd, C.int32(pl), fl, xx, err)
	})
}

type _nodApsFunc func(jd C.double, pl, fl, m C.int32, nasc, ndsc, peri, aphe *C.double, err *C.char) C.int32

func _nodAps(jd float64, pl swego.Planet, fl int32, m swego.NodApsMethod, fn _nodApsFunc) (_, _, _, _ []float64, err error) {
	_jd := C.double(jd)
	_pl := C.int32(pl)
	_fl := C.int32(fl)
	_m := C.int32(m)

	// Both the float64 and C.double types are defined as an IEEE-754 64-bit
	// floating-point number. This means the representation in memory of a
	// float64 array is equivalent to that of a C.double array. It means it is
	// possible to cast between the two types. In Go land such operation is
	// considered unsafe, hence the use of the unsafe package.
	var nasc, ndsc, peri, aphe [6]float64
	_nasc := (*C.double)(unsafe.Pointer(&nasc[0]))
	_ndsc := (*C.double)(unsafe.Pointer(&ndsc[0]))
	_peri := (*C.double)(unsafe.Pointer(&peri[0]))
	_aphe := (*C.double)(unsafe.Pointer(&aphe[0]))

	err = withError(func(err *C.char) bool {
		return C.ERR == fn(_jd, _pl, _fl, _m, _nasc, _ndsc, _peri, _aphe, err)
	})

	return nasc[:], ndsc[:], peri[:], aphe[:], err
}

func nodAps(et float64, pl swego.Planet, fl int32, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	return _nodAps(et, pl, fl, m, func(jd C.double, pl, fl, m C.int32, nasc, ndsc, peri, aphe *C.double, err *C.char) C.int32 {
		return C.swe_nod_aps(jd, pl, fl, m, nasc, ndsc, peri, aphe, err)
	})
}

func nodApsUT(ut float64, pl swego.Planet, fl int32, m swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error) {
	return _nodAps(ut, pl, fl, m, func(jd C.double, pl, fl, m C.int32, nasc, ndsc, peri, aphe *C.double, err *C.char) C.int32 {
		return C.swe_nod_aps_ut(jd, pl, fl, m, nasc, ndsc, peri, aphe, err)
	})
}

type _getAyanamsaExFunc func(jd C.double, fl C.int32, aya *C.double, err *C.char) C.int32

func _getAyanamsaEx(jd float64, fl int32, fn _getAyanamsaExFunc) (aya float64, err error) {
	_jd := C.double(jd)
	_fl := C.int32(fl)
	_aya := (*C.double)(unsafe.Pointer(&aya))

	err = withError(func(err *C.char) bool {
		return C.ERR == fn(_jd, _fl, _aya, err)
	})

	return
}

func getAyanamsaEx(et float64, fl int32) (float64, error) {
	return _getAyanamsaEx(et, fl, func(jd C.double, fl C.int32, aya *C.double, err *C.char) C.int32 {
		return C.swe_get_ayanamsa_ex(jd, fl, aya, err)
	})
}

func getAyanamsaExUT(ut float64, fl int32) (float64, error) {
	return _getAyanamsaEx(ut, fl, func(jd C.double, fl C.int32, aya *C.double, err *C.char) C.int32 {
		return C.swe_get_ayanamsa_ex_ut(jd, fl, aya, err)
	})
}

func getAyanamsaName(ayan swego.Ayanamsa) string {
	return C.GoString(C.swe_get_ayanamsa_name(C.int32(ayan)))
}

func julDay(y, m, d int, h float64, gf int) float64 {
	_y := C.int(y)
	_m := C.int(m)
	_d := C.int(d)
	_h := C.double(h)
	_gf := C.int(gf)
	return float64(C.swe_julday(_y, _m, _d, _h, _gf))
}

func revJul(jd float64, gf int) (y, m, d int, h float64) {
	_jd := C.double(jd)
	_gf := C.int(gf)
	var _y, _m, _d C.int
	var _h C.double

	C.swe_revjul(_jd, _gf, &_y, &_m, &_d, &_h)

	y = int(_y)
	m = int(_m)
	d = int(_d)
	h = float64(_h)
	return
}

func utcToJD(y, m, d, h, i int, s float64, gf int) (et, ut float64, err error) {
	_y := C.int32(y)
	_m := C.int32(m)
	_d := C.int32(d)
	_h := C.int32(h)
	_i := C.int32(i)
	_s := C.double(s)
	_gf := C.int32(gf)
	var jds [2]C.double

	err = withError(func(err *C.char) bool {
		return C.ERR == C.swe_utc_to_jd(_y, _m, _d, _h, _i, _s, _gf, &jds[0], err)
	})

	et = float64(jds[0])
	ut = float64(jds[1])
	return
}

type _jdToUTCFunc func(jd C.double, gf C.int32, y, m, d, h, i *C.int32, s *C.double)

func _jdToUTC(jd float64, gf int, fn _jdToUTCFunc) (y, m, d, h, i int, s float64) {
	_jd := C.double(jd)
	_gf := C.int32(gf)
	var _y, _m, _d, _h, _i C.int32
	var _s C.double

	fn(_jd, _gf, &_y, &_m, &_d, &_h, &_i, &_s)

	y = int(_y)
	m = int(_m)
	d = int(_d)
	h = int(_h)
	i = int(_i)
	s = float64(_s)
	return
}

func jdETToUTC(et float64, gf int) (y, m, d, h, i int, s float64) {
	return _jdToUTC(et, gf, func(jd C.double, gf C.int32, y, m, d, h, i *C.int32, s *C.double) {
		C.swe_jdet_to_utc(jd, gf, y, m, d, h, i, s)
	})
}

func jdUT1ToUTC(ut float64, gf int) (y, m, d, h, i int, s float64) {
	return _jdToUTC(ut, gf, func(jd C.double, gf C.int32, y, m, d, h, i *C.int32, s *C.double) {
		C.swe_jdut1_to_utc(jd, gf, y, m, d, h, i, s)
	})
}

type _housesFunc func(lat C.double, hsys C.int, cusps, ascmc *C.double) C.int

func _houses(lat float64, hsys swego.HSys, fn _housesFunc) (_, _ []float64, err error) {
	_lat := C.double(lat)
	_hsys := C.int(hsys)

	// Both the float64 and C.double types are defined as an IEEE-754 64-bit
	// floating-point number. This means the representation in memory of a
	// float64 array is equivalent to that of a C.double array. It means it is
	// possible to cast between the two types. In Go land such operation is
	// considered unsafe, hence the use of the unsafe package.
	var cusps [37]float64
	var ascmc [10]float64
	_cusps := (*C.double)(unsafe.Pointer(&cusps[0]))
	_ascmc := (*C.double)(unsafe.Pointer(&ascmc[0]))

	if C.ERR == fn(_lat, _hsys, _cusps, _ascmc) {
		err = swego.Error("swe_house() error")
	}

	// The house system letters are practically constants. If those are changed,
	// it is done via a new version of the Swiss Ephemeris anyway. Also this is
	// already fairly low level code, so a check for the Gauquelin 'house system'
	// letter is no problem.
	n := 13
	if hsys == 'G' || hsys == 'g' {
		n = 37
	}

	return cusps[:n:n], ascmc[:], err
}

func housesEx(ut float64, fl int32, lat, lng float64, hsys swego.HSys) ([]float64, []float64, error) {
	return _houses(lat, hsys, func(lat C.double, hsys C.int, cusps, ascmc *C.double) C.int {
		_jd := C.double(ut)
		_fl := C.int32(fl)
		_lng := C.double(lng)
		return C.swe_houses_ex(_jd, _fl, lat, _lng, hsys, cusps, ascmc)
	})
}

func housesARMC(armc, lat, eps float64, hsys swego.HSys) ([]float64, []float64, error) {
	return _houses(lat, hsys, func(lat C.double, hsys C.int, cusps, ascmc *C.double) C.int {
		_armc := C.double(armc)
		_eps := C.double(eps)
		return C.swe_houses_armc(_armc, lat, _eps, hsys, cusps, ascmc)
	})
}

func housePos(armc, geolat, eps float64, hsys swego.HSys, pllng, pllat float64) (pos float64, err error) {
	_armc := C.double(armc)
	_lat := C.double(geolat)
	_eps := C.double(eps)
	_hsys := C.int(hsys)
	xpin := [2]C.double{C.double(pllat), C.double(pllng)}

	err = withError(func(err *C.char) bool {
		pos = float64(C.swe_house_pos(_armc, _lat, _eps, _hsys, &xpin[0], err))
		return *err != '\000'
	})

	return
}

func houseName(hsys swego.HSys) string {
	return C.GoString(C.swe_house_name(C.int(hsys)))
}

func deltaTEx(jd float64, eph int32) (deltaT float64, err error) {
	err = withError(func(err *C.char) bool {
		deltaT = float64(C.swe_deltat_ex(C.double(jd), C.int32(eph), err))
		return *err != '\000'
	})

	return
}

func setDeltaTUserDef(v float64) {
	C.swe_set_delta_t_userdef(C.double(v))
}

func timeEqu(jd float64) (E float64, err error) {
	var _E C.double

	err = withError(func(err *C.char) bool {
		return C.ERR == C.swe_time_equ(C.double(jd), &_E, err)
	})

	E = float64(_E)
	return
}

type _convertLMTLATFunc func(from, lng C.double, to *C.double, err *C.char) C.int32

func _convertLMTLAT(from, geolon float64, fn _convertLMTLATFunc) (to float64, err error) {
	_from := C.double(from)
	_lng := C.double(geolon)
	var _to C.double

	err = withError(func(err *C.char) bool {
		return C.ERR == fn(_from, _lng, &_to, err)
	})

	to = float64(_to)
	return
}

func lmtToLAT(jdLMT, geolon float64) (float64, error) {
	return _convertLMTLAT(jdLMT, geolon, func(from, lng C.double, to *C.double, err *C.char) C.int32 {
		return C.swe_lmt_to_lat(from, lng, to, err)
	})
}

func latToLMT(jdLAT, geolon float64) (float64, error) {
	return _convertLMTLAT(jdLAT, geolon, func(from, lng C.double, to *C.double, err *C.char) C.int32 {
		return C.swe_lat_to_lmt(from, lng, to, err)
	})
}

func sidTime0(ut, eps, nut float64) float64 {
	_ut := C.double(ut)
	_eps := C.double(eps)
	_nut := C.double(nut)
	return float64(C.swe_sidtime0(_ut, _eps, _nut))
}

func sidTime(ut float64) float64 {
	return float64(C.swe_sidtime(C.double(ut)))
}

// Returns (ideg, imin, isec, dsecfr, isgn)
func splitDeg(ddeg float64, roundflag int) (int32, int32, int32, float64, int32) {
	var ideg C.int32
	var imin C.int32
	var isec C.int32
	var dsecfr C.double
	var isgn C.int32
	C.swe_split_deg(C.double(ddeg), C.int32(roundflag), &ideg, &imin, &isec, &dsecfr, &isgn)
	return int32(ideg), int32(imin), int32(isec), float64(dsecfr), int32(isgn)
}
