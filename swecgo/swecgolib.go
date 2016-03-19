package swecgo

import (
	"errors"
	"unicode"
	"unsafe"

	"github.com/dwlnetnl/swego"
)

/*

// To aid debugging the Swiss Ephemeris supports tracing. It will write two
// files in the current working directory. Note that the Go tool may change the
// current working directory. The tracing facility has 2 modes:
//
//  TRACE=1 writes swetrace.{c,txt} and appends each session to it.
//  TRACE=2 writes swetrace_<pid>.{c, txt} where pid is the running process id.
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
// tracing can be reopened after it is closed for the first time.
//
// If the internal package is dirty and Go needs to build it every time it's
// very easy to enable the tracing and the trace limit per build using the
// CGO_CFLAGS environment variable. To inspect the Go build process pass the -x
// argument flag.

// #cgo CFLAGS: -DTRACE=1
// #cgo CFLAGS: -DTRACE=2
// #cgo CFLAGS: -DTRACE_COUNT_MAX=10000000
// #cgo CFLAGS: -DTRACE_CLOSE=1

// ----------

// Enable thread local storage in library.
// #cgo CFLAGS: -DTLS_ENABLED=1

int swex_supports_tls() {
#if TLS_ENABLED
	return 1;
#else
	return 0;
#endif
}

// ----------

#cgo CFLAGS: -w

#include <stdlib.h>
#include <string.h>
#include "swephexp.h"
#include "sweph.h"

void swex_set_topo(double geolon, double geolat, double geoalt) {
	if (swed.geopos_is_set == FALSE
		|| swed.topd.geolon != geolon
		|| swed.topd.geolat != geolat
		|| swed.topd.geoalt != geoalt
	)	{
		swe_set_topo(geolon, geolat, geoalt);
	}
}

void swex_set_sid_mode(int32 sid_mode, double t0, double ayan_t0) {
	if (swed.ayana_is_set == FALSE
		|| swed.sidd.sid_mode != sid_mode
		|| swed.sidd.ayan_t0 != ayan_t0
		|| swed.sidd.t0 != t0
	) {
		swe_set_sid_mode(sid_mode, t0, ayan_t0);
	}
}

void swex_set_jpl_file(const char *fname) {
	if (strncmp(fname, swed.jplfnam, strlen(fname)) != 0) {
		swe_set_jpl_file(fname);
	}
}

*/
import "C"

func supportsTLS() bool {
	return C.swex_supports_tls() == 1
}

const errPrefix = "swecgo: "

func withError(fn func(err *C.char) bool) error {
	var _err [C.AS_MAXCH]C.char

	if fn(&_err[0]) {
		return errors.New(errPrefix + C.GoString(&_err[0]))
	}

	return nil
}

// Version constant contains the current Swiss Ephemeris version.
const Version = C.SE_VERSION

// DefaultPath is the default ephemeris path defined by the library.
const DefaultPath = C.SE_EPHE_PATH

func setEphePath(path string) {
	_path := C.CString(path)
	C.swe_set_ephe_path(_path)
	C.free(unsafe.Pointer(_path))
}

func setTopo(lng, lat, alt float64) {
	_lng := C.double(lng)
	_lat := C.double(lat)
	_alt := C.double(alt)
	C.swex_set_topo(_lng, _lat, _alt)
}

func setSidMode(mode swego.Ayanamsa, t0, ayanT0 float64) {
	_mode := C.int32(mode)
	_t0 := C.double(t0)
	_ayanT0 := C.double(ayanT0)
	C.swex_set_sid_mode(_mode, _t0, _ayanT0)
}

func setFileNameJPL(name string) {
	_name := C.CString(name)
	C.swex_set_jpl_file(_name)
	C.free(unsafe.Pointer(_name))
}

func close() {
	C.swe_close()
}

const (
	flgTopo     = C.SEFLG_TOPOCTR
	flgSidereal = C.SEFLG_SIDEREAL
)

type _calcFunc func(jd C.double, fl C.int32, xx *C.double, err *C.char) C.int32

func _calc(jd float64, fl int32, fn _calcFunc) (xx [6]float64, cfl int, err error) {
	_jd := C.double(jd)
	_fl := C.int32(fl)
	_xx := (*C.double)(unsafe.Pointer(&xx[0]))

	err = withError(func(err *C.char) bool {
		cfl = int(fn(_jd, _fl, _xx, err))
		return cfl == C.ERR
	})

	return
}

func calc(et float64, pl swego.Planet, fl int32) ([6]float64, int, error) {
	return _calc(et, fl, func(jd C.double, fl C.int32, xx *C.double, err *C.char) C.int32 {
		return C.swe_calc(jd, C.int(pl), fl, xx, err)
	})
}

func calcUT(ut float64, pl swego.Planet, fl int32) ([6]float64, int, error) {
	return _calc(ut, fl, func(jd C.double, fl C.int32, xx *C.double, err *C.char) C.int32 {
		return C.swe_calc_ut(jd, C.int32(pl), fl, xx, err)
	})
}

func planetName(pl swego.Planet) string {
	var _name [C.AS_MAXCH]C.char
	C.swe_get_planet_name(C.int(pl), &_name[0])
	return C.GoString(&_name[0])
}

func getAyanamsa(et float64) float64 {
	return float64(C.swe_get_ayanamsa(C.double(et)))
}

func getAyanamsaUT(ut float64) float64 {
	return float64(C.swe_get_ayanamsa_ut(C.double(ut)))
}

type _getAyanamsaExFunc func(jd C.double, fl C.int32, aya *C.double, err *C.char) C.int32

func _getAyanamsaEx(jd float64, fl int32, fn _getAyanamsaExFunc) (aya float64, err error) {
	_jd := C.double(jd)
	_fl := C.int32(fl)
	_aya := (*C.double)(unsafe.Pointer(&aya))

	err = withError(func(err *C.char) bool {
		rc := int(fn(_jd, _fl, _aya, err))
		// rc := int(C.swe_get_ayanamsa_ex(_jd, _fl, _aya, err))
		return rc == C.ERR
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
	_y := (*C.int)(unsafe.Pointer(&y))
	_m := (*C.int)(unsafe.Pointer(&m))
	_d := (*C.int)(unsafe.Pointer(&d))
	_h := (*C.double)(unsafe.Pointer(&h))
	C.swe_revjul(_jd, _gf, _y, _m, _d, _h)
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
		rc := int(C.swe_utc_to_jd(_y, _m, _d, _h, _i, _s, _gf, &jds[0], err))
		return rc == C.ERR
	})

	et = float64(jds[0])
	ut = float64(jds[1])
	return
}

type _jdToUTCFunc func(jd C.double, gf C.int32, y, m, d, h, i *C.int32, s *C.double)

func _jdToUTC(jd float64, gf int, fn _jdToUTCFunc) (y, m, d, h, i int, s float64) {
	_jd := C.double(jd)
	_gf := C.int32(gf)
	_y := (*C.int32)(unsafe.Pointer(&y))
	_m := (*C.int32)(unsafe.Pointer(&m))
	_d := (*C.int32)(unsafe.Pointer(&d))
	_h := (*C.int32)(unsafe.Pointer(&h))
	_i := (*C.int32)(unsafe.Pointer(&i))
	_s := (*C.double)(unsafe.Pointer(&s))
	fn(_jd, _gf, _y, _m, _d, _h, _i, _s)
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

func _houses(lat float64, hsys swego.HSys, fn _housesFunc) (_ []float64, ascmc [10]float64, err error) {
	_lat := C.double(lat)
	_hsys := C.int(hsys)

	var cusps [37]float64
	_cusps := (*C.double)(unsafe.Pointer(&cusps[0]))
	_ascmc := (*C.double)(unsafe.Pointer(&ascmc[0]))

	if C.ERR == fn(_lat, _hsys, _cusps, _ascmc) {
		err = errors.New(errPrefix + "error calculating houses")
	}

	n := 13
	if swego.HSys(unicode.ToUpper(rune(hsys))) == swego.Gauquelin {
		n = 37
	}

	return cusps[:n:n], ascmc, err
}

func houses(ut, lat, lng float64, hsys swego.HSys) ([]float64, [10]float64, error) {
	return _houses(lat, hsys, func(lat C.double, hsys C.int, cusps, ascmc *C.double) C.int {
		_jd := C.double(ut)
		_lng := C.double(lng)
		return C.swe_houses(_jd, lat, _lng, hsys, cusps, ascmc)
	})
}

func housesEx(ut float64, fl int32, lat, lng float64, hsys swego.HSys) ([]float64, [10]float64, error) {
	return _houses(lat, hsys, func(lat C.double, hsys C.int, cusps, ascmc *C.double) C.int {
		_jd := C.double(ut)
		_fl := C.int32(fl)
		_lng := C.double(lng)
		return C.swe_houses_ex(_jd, _fl, lat, _lng, hsys, cusps, ascmc)
	})
}

func housesArmc(armc, lat, eps float64, hsys swego.HSys) ([]float64, [10]float64, error) {
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

func deltaT(jd float64) float64 {
	return float64(C.swe_deltat(C.double(jd)))
}

func deltaTEx(jd float64, eph int32) (deltaT float64, err error) {
	err = withError(func(err *C.char) bool {
		deltaT = float64(C.swe_deltat_ex(C.double(jd), C.int32(eph), err))
		return *err != '\000'
	})

	return
}

func timeEqu(jd float64) (E float64, err error) {
	_E := (*C.double)(unsafe.Pointer(&E))

	err = withError(func(err *C.char) bool {
		rc := int(C.swe_time_equ(C.double(jd), _E, err))
		return rc == C.ERR
	})

	return
}

type _convertLMTLATFunc func(from, lng C.double, to *C.double, err *C.char) C.int32

func _convertLMTLAT(from, geolon float64, fn _convertLMTLATFunc) (to float64, err error) {
	_from := C.double(from)
	_lng := C.double(geolon)
	_to := (*C.double)(unsafe.Pointer(&to))

	err = withError(func(err *C.char) bool {
		rc := int(fn(_from, _lng, _to, err))
		return rc == C.ERR
	})

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
