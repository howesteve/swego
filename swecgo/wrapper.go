package swecgo

import (
	"errors"
	"sync"
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

var gMutex sync.Mutex

func withLock(fn func()) {
	if supportsTLS() {
		fn()
		return
	}

	gMutex.Lock()
	fn()
	gMutex.Unlock()
}

// Version constant contains the current Swiss Ephemeris version.
const Version = C.SE_VERSION

// DefaultPath is the default ephemeris path.
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

const errPrefix = "swecgo: "

func calc(et float64, pl int, fl int32) (xx [6]float64, cfl int, err error) {
	_jd := C.double(et)
	_pl := C.int(pl)
	_fl := C.int32(fl)

	_xx := (*C.double)(unsafe.Pointer(&xx[0]))
	var _err [C.AS_MAXCH]C.char

	cfl = int(C.swe_calc(_jd, _pl, _fl, _xx, &_err[0]))
	if cfl == C.ERR {
		err = errors.New(errPrefix + C.GoString(&_err[0]))
		return
	}

	return
}

func calcUT(ut float64, pl int, fl int32) (xx [6]float64, cfl int, err error) {
	_jd := C.double(ut)
	_pl := C.int32(pl)
	_fl := C.int32(fl)

	_xx := (*C.double)(unsafe.Pointer(&xx[0]))
	var _err [C.AS_MAXCH]C.char

	cfl = int(C.swe_calc_ut(_jd, _pl, _fl, _xx, &_err[0]))
	if cfl == C.ERR {
		err = errors.New(errPrefix + C.GoString(&_err[0]))
		return
	}

	return
}

func planetName(pl int) string {
	var _name [C.AS_MAXCH]C.char
	C.swe_get_planet_name(C.int(pl), &_name[0])
	return C.GoString(&_name[0])
}
