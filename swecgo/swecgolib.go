package swecgo

import (
	"errors"
	"unsafe"
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

*/
import "C"

func supportsTLS() bool {
	return C.swex_supports_tls() == 1
}

func withError(fn func(err *C.char) bool) error {
	var _err [C.AS_MAXCH]C.char

	if fn(&_err[0]) {
		return errors.New("swecgo: " + C.GoString(&_err[0]))
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

func setSidMode(mode int32, t0, ayanT0 float64) {
	_mode := C.int32(mode)
	_t0 := C.double(t0)
	_ayanT0 := C.double(ayanT0)
	C.swex_set_sid_mode(_mode, _t0, _ayanT0)
}

func setFileNameJPL(name string) {
	_name := C.CString(name)
	C.swe_set_jpl_file(_name)
	C.free(unsafe.Pointer(_name))
}

func close() {
	C.swe_close()
}

const (
	flgTopo     = C.SEFLG_TOPOCTR
	flgSidereal = C.SEFLG_SIDEREAL
)

func calc(et float64, pl int, fl int32) (xx [6]float64, cfl int, err error) {
	_jd := C.double(et)
	_pl := C.int(pl)
	_fl := C.int32(fl)
	_xx := (*C.double)(unsafe.Pointer(&xx[0]))

	err = withError(func(err *C.char) bool {
		cfl = int(C.swe_calc(_jd, _pl, _fl, _xx, err))
		return cfl == C.ERR
	})

	return
}

func calcUT(ut float64, pl int, fl int32) (xx [6]float64, cfl int, err error) {
	_jd := C.double(ut)
	_pl := C.int32(pl)
	_fl := C.int32(fl)
	_xx := (*C.double)(unsafe.Pointer(&xx[0]))

	err = withError(func(err *C.char) bool {
		cfl = int(C.swe_calc_ut(_jd, _pl, _fl, _xx, err))
		return cfl == C.ERR
	})

	return
}

func planetName(pl int) string {
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

func getAyanamsaEx(et float64, fl int32) (aya float64, err error) {
	_jd := C.double(et)
	_fl := C.int32(fl)
	_aya := (*C.double)(unsafe.Pointer(&aya))

	err = withError(func(err *C.char) bool {
		rc := int(C.swe_get_ayanamsa_ex(_jd, _fl, _aya, err))
		return rc == C.ERR
	})

	return
}

func getAyanamsaExUT(ut float64, fl int32) (aya float64, err error) {
	_jd := C.double(ut)
	_fl := C.int32(fl)
	_aya := (*C.double)(unsafe.Pointer(&aya))

	err = withError(func(err *C.char) bool {
		rc := int(C.swe_get_ayanamsa_ex_ut(_jd, _fl, _aya, err))
		return rc == C.ERR
	})

	return
}

func getAyanamsaName(sidmode int32) string {
	return C.GoString(C.swe_get_ayanamsa_name(C.int32(sidmode)))
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

func jdETToUTC(et float64, gf int) (y, m, d, h, i int, s float64) {
	_jd := C.double(et)
	_gf := C.int32(gf)
	_y := (*C.int32)(unsafe.Pointer(&y))
	_m := (*C.int32)(unsafe.Pointer(&m))
	_d := (*C.int32)(unsafe.Pointer(&d))
	_h := (*C.int32)(unsafe.Pointer(&h))
	_i := (*C.int32)(unsafe.Pointer(&i))
	_s := (*C.double)(unsafe.Pointer(&s))
	C.swe_jdet_to_utc(_jd, _gf, _y, _m, _d, _h, _i, _s)
	return
}

func jdUT1ToUTC(ut float64, gf int) (y, m, d, h, i int, s float64) {
	_jd := C.double(ut)
	_gf := C.int32(gf)
	_y := (*C.int32)(unsafe.Pointer(&y))
	_m := (*C.int32)(unsafe.Pointer(&m))
	_d := (*C.int32)(unsafe.Pointer(&d))
	_h := (*C.int32)(unsafe.Pointer(&h))
	_i := (*C.int32)(unsafe.Pointer(&i))
	_s := (*C.double)(unsafe.Pointer(&s))
	C.swe_jdut1_to_utc(_jd, _gf, _y, _m, _d, _h, _i, _s)
	return
}

func houseName(hsys int) string {
	return C.GoString(C.swe_house_name(C.int(hsys)))
}

func deltaT(jd float64) float64 {
	return float64(C.swe_deltat(C.double(jd)))
}

func deltaTEx(jd float64, eph int32) (float64, error) {
	var deltaT float64

	err := withError(func(err *C.char) bool {
		deltaT = float64(C.swe_deltat_ex(C.double(jd), C.int32(eph), err))
		return *err != '\000'
	})

	return deltaT, err
}

func timeEqu(jd float64) (E float64, err error) {
	_E := (*C.double)(unsafe.Pointer(&E))

	err = withError(func(err *C.char) bool {
		rc := int(C.swe_time_equ(C.double(jd), _E, err))
		return rc == C.ERR
	})

	return
}

func lmtToLAT(jdLMT, geolon float64) (jdLAT float64, err error) {
	_lmt := C.double(jdLMT)
	_geolon := C.double(geolon)
	_lat := (*C.double)(unsafe.Pointer(&jdLAT))

	err = withError(func(err *C.char) bool {
		rc := int(C.swe_lmt_to_lat(_lmt, _geolon, _lat, err))
		return rc == C.ERR
	})

	return
}

func latToLMT(jdLAT, geolon float64) (jdLMT float64, err error) {
	_lat := C.double(jdLAT)
	_geolon := C.double(geolon)
	_lmt := (*C.double)(unsafe.Pointer(&jdLMT))

	err = withError(func(err *C.char) bool {
		rc := int(C.swe_lat_to_lmt(_lat, _geolon, _lmt, err))
		return rc == C.ERR
	})

	return
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
