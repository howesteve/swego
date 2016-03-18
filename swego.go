// Package swego defines an interface for interfacing with the Swiss Ephemeris.
package swego

// CalcFlags represents the flags argument of swe_calc and swe_calc_ut in a
// stateless way.
type CalcFlags struct {
	Flags   int32
	TopoLoc TopoLoc
	SidMode SidMode

	// FileNameJPL represents the argument to swe_set_jpl_file.
	FileNameJPL string
}

// TopoLoc represents the arguments to swe_set_topo.
type TopoLoc struct {
	Lat  float64
	Long float64
	Alt  float64
}

// AyanamsaExFlags represents the flags argument of swe_get_ayanamsa_ex and
// swe_get_ayanamsa_ex_ut in a stateless way.
type AyanamsaExFlags struct {
	Flags   int32
	SidMode SidMode
}

// HousesExFlags represents the flags argument of swe_houses_ex in a stateless
// way.
type HousesExFlags struct {
	Flags   int32
	SidMode SidMode
}

// SidMode represents the arguments of swe_set_sid_mode.
type SidMode struct {
	Mode   int32
	T0     float64
	AyanT0 float64
}

// CalType represents the calendar type used in julian date conversion.
type CalType byte

// Calendar types.
const (
	Julian    CalType = 'j'
	Gregorian CalType = 'g'
)

// Interface defines a standardized way for interfacing with the Swiss
// Ephemeris library from Go.
type Interface interface {
	Version() string

	SetPath(ephepath string)
	Close()

	Calc(et float64, pl int, fl CalcFlags) (xx [6]float64, cfl int, err error)
	CalcUT(ut float64, pl int, fl CalcFlags) (xx [6]float64, cfl int, err error)

	PlanetName(pl int) string

	GetAyanamsa(et float64) float64
	GetAyanamsaUT(ut float64) float64
	GetAyanamsaEx(et float64, fl AyanamsaExFlags) (float64, error)
	GetAyanamsaExUT(ut float64, fl AyanamsaExFlags) (float64, error)
	GetAyanamsaName(sidmode int32) string

	JulDay(y, m, d int, h float64, ct CalType) float64
	RevJul(jd float64, c byte) (y, m, d int, h float64)
	UTCToJD(y, m, d int, h float64, ct CalType) float64
	JdETToUTC(et float64, c byte) (y, m, d, h, i int, s float64)
	JdUT1ToUTC(ut1 float64, c byte) (y, m, d, h, i int, s float64)

	Houses(ut, geolat, geolon float64, hsys int) ([]float64, [10]float64)
	HousesEx(ut float64, fl HousesExFlags, geolat, geolon float64, hsys int) ([]float64, [10]float64)
	HousesArmc(armc, geolat, eps float64, hsys int) ([]float64, [10]float64)
	HousePos(armc, geolat, eps float64, hsys int, xpin [2]float64) (float64, error)
	HouseName(hsys int) string

	DeltaT(jd float64) float64
	DeltaTEx(jd float64, fl int32) (float64, error)

	TimeEqu(jd float64) (float64, error)
	LMTToLAT(jdLMT, geolon float64) (float64, error)
	LATToLMT(jdLAT, geolon float64) (float64, error)

	SidTime0(ut, eps, nut float64) float64
	SidTime(ut float64) float64
}
