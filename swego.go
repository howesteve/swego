// Package swego defines an interface for interfacing with the Swiss Ephemeris.
package swego

import "unicode"

// CalType represents the calendar type used in julian date conversion.
type CalType int

// Planet is the type of planet constants.
type Planet int

// Ayanamsa is the type of sidereal mode constants.
type Ayanamsa int32

// SidMode represents the arguments of swe_set_sid_mode.
type SidMode struct {
	Mode   Ayanamsa
	T0     float64
	AyanT0 float64
}

// TopoLoc represents the arguments to swe_set_topo.
type TopoLoc struct {
	Long float64
	Lat  float64
	Alt  float64
}

// CalcFlags represents the flags argument of swe_calc and swe_calc_ut in a
// stateless way.
type CalcFlags struct {
	Flags   int32
	TopoLoc *TopoLoc
	SidMode *SidMode

	// FileNameJPL represents the argument to swe_set_jpl_file.
	FileNameJPL string
}

// Ephemeris represents an ephemeris implemented in the C library.
type Ephemeris int32

// SetEphemeris sets the ephemeris flag in fl.
func (fl *CalcFlags) SetEphemeris(eph Ephemeris) { fl.Flags |= int32(eph) }

// NodApsMethod is the type of Nodbit constants.
type NodApsMethod int32

// AyanamsaExFlags represents the flags argument of swe_get_ayanamsa_ex and
// swe_get_ayanamsa_ex_ut in a stateless way.
type AyanamsaExFlags struct {
	Flags   int32
	SidMode *SidMode
}

// HousesExFlags represents the flags argument of swe_houses_ex in a stateless
// way.
type HousesExFlags struct {
	Flags   int32
	SidMode *SidMode
}

// HSys represents house system identifiers used in the C library.
// It's assumed to be an upper case ASCII letter.
type HSys byte

// NewHSys validates the input and returns a HSys value if valid.
func NewHSys(char byte) (hsys HSys, ok bool) {
	if char == 'i' {
		return HSys(char), true
	}

	// package unicode operates on runes
	char = byte(unicode.ToUpper(rune(char)))

	switch char {
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y':
		return HSys(char), true
	default:
		return '\000', false
	}
}

// Interface defines a standardized way for interfacing with the Swiss
// Ephemeris library from Go.
type Interface interface {
	// Version returns the version of the Swiss Ephemeris.
	Version() string
	// GetLibraryPath returns the file system path of the executable.
	GetLibraryPath() string

	// SetPath opens the ephemeris and sets the data path.
	SetPath(path string)
	// Close closes the Swiss Ephemeris library.
	// The ephemeris can be reopened by calling SetPath.
	Close()

	// PlanetName returns the name of planet pl.
	PlanetName(pl Planet) string

	// Calc computes the position and optionally the speed of planet pl at Julian
	// Date (in Ephemeris Time) et with calculation flags fl.
	Calc(et float64, pl Planet, fl *CalcFlags) (xx []float64, cfl int, err error)
	// CalcUT computes the position and optionally the speed of planet pl at
	// Julian Date (in Universal Time) ut with calculation flags fl. Within the C
	// library swe_deltat is called to convert Universal Time to Ephemeris Time.
	CalcUT(ut float64, pl Planet, fl *CalcFlags) (xx []float64, cfl int, err error)

	// NodAps computes the positions of planetary nodes and apsides (perihelia,
	// aphelia, second focal points of the orbital ellipses) for planet pl at
	// Julian Date (in Ephemeris Time) et with calculation flags fl using method
	// m.
	NodAps(et float64, pl Planet, fl *CalcFlags, m NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error)
	// NodApsUT computes the positions of planetary nodes and apsides (perihelia,
	// aphelia, second focal points of the orbital ellipses) for planet pl at
	// Julian Date (in Ephemeris Time) et with calculation flags fl using method
	// m. Within the C library swe_deltat is called to convert Universal Time to
	// Ephemeris Time.
	NodApsUT(ut float64, pl Planet, fl *CalcFlags, m NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error)

	// GetAyanamsa returns the ayanamsa for Julian Date (in Ephemeris Time) et.
	// You should use GetAyanamsaEx, see the Programmer's Documentation.
	GetAyanamsa(et float64, sidmode *SidMode) float64
	// GetAyanamsaUT returns the ayanamsa for Julian Date (in Universal Time) ut.
	// You should use GetAyanamsaExUT, see the Programmer's Documentation.
	GetAyanamsaUT(ut float64, sidmode *SidMode) float64
	// GetAyanamsaEx returns the ayanamsa for Julian Date (in Ephemeris Time) et.
	// It is equal to GetAyanamsa but uses the ΔT consistent with the ephemeris
	// passed in fl.Flags.
	GetAyanamsaEx(et float64, fl *AyanamsaExFlags) (float64, error)
	// GetAyanamsaExUT returns the ayanamsa for Julian Date (in Universal Time) ut.
	// It is equal to GetAyanamsaUT but uses the ΔT consistent with the ephemeris
	// passed in fl.Flags.
	GetAyanamsaExUT(ut float64, fl *AyanamsaExFlags) (float64, error)
	// GetAyanamsaName returns the name of sidmode.
	GetAyanamsaName(ayan Ayanamsa) string

	// JulDay returns the corresponding Julian Date for the given date. Calendar
	// type ct is used to clearify the year y, Julian or Gregorian.
	JulDay(y, m, d int, h float64, ct CalType) float64
	// RevJul returns the corresponding calendar date for the given Julian Date.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	RevJul(jd float64, ct CalType) (y, m, d int, h float64)
	// UTCToJD returns the corresponding Julian Date in Ephemeris and Universal
	// Time for the given date and accounts for leap seconds in the conversion.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	UTCToJD(y, m, d, h, i int, s float64, ct CalType) (et, ut float64, err error)
	// JdETToUTC returns the corresponding calendar date for the given Julian
	// Date in Ephemeris Time and accounts for leap seconds in the conversion.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	JdETToUTC(et float64, ct CalType) (y, m, d, h, i int, s float64)
	// JdETToUTC returns the corresponding calendar date for the given Julian
	// Date in Universal Time and accounts for leap seconds in the conversion.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	JdUT1ToUTC(ut1 float64, ct CalType) (y, m, d, h, i int, s float64)

	// Houses returns the house cusps and related positions for the given
	// geographic location using the given house system. The return values may
	// contain data in case of an error. Geolat and geolon are in degrees.
	Houses(ut, geolat, geolon float64, hsys HSys) ([]float64, []float64, error)
	// HousesEx returns the house cusps and related positions for the given
	// geographic location using the given house system and the provided flags
	// (reference frame). The return values may contain data in case of an error.
	// Geolat and geolon are in degrees.
	HousesEx(ut float64, fl *HousesExFlags, geolat, geolon float64, hsys HSys) ([]float64, []float64, error)
	// HousesArmc returns the house cusps and related positions for the given
	// geographic location using the given house system, ecliptic obliquity and
	// ARMC (also known as RAMC). The return values may contain data in case of
	// an error. ARMC, geolat, geolon and eps are in degrees.
	HousesARMC(armc, geolat, eps float64, hsys HSys) ([]float64, []float64, error)
	// HousePos returns the house position for the ecliptic longitude and
	// latitude of a planet for a given ARMC (also known as RAMC) and geocentric
	// latitude using the given house system. ARMC, geolat, eps, pllng and pllat
	// are in degrees.
	// Before calling HousePos either Houses, HousesEx or HousesARMC should be
	// called first.
	HousePos(armc, geolat, eps float64, hsys HSys, pllng, pllat float64) (float64, error)
	// HouseName returns the name of the house system.
	HouseName(hsys HSys) string

	// DeltaT returns the ΔT for the Julian Date jd.
	// You should use DeltaTEx, see the Programmer's Documentation.
	DeltaT(jd float64) float64
	// DeltaTEx returns the ΔT for the Julian Date jd.
	DeltaTEx(jd float64, eph Ephemeris) (float64, error)
	// SetDeltaTUserDef sets a user supplied ΔT value to be returned by the
	// DeltaT and DeltaTEx methods. Use ResetDeltaT to reset the value.
	SetDeltaTUserDef(v float64)

	// TimeEqu returns the difference between local apparent and local mean time
	// in days for the given Julian Date (in Universal Time).
	TimeEqu(jd float64) (float64, error)
	// LMTToLAT returns the local apparent time for the given Julian Date (in
	// Local Mean Time) and the geographic longitude.
	LMTToLAT(jdLMT, geolon float64) (float64, error)
	// LATToLMT returns the local mean time for the given Julian Date (in Local
	// Apparent Time) and the geographic longitude.
	LATToLMT(jdLAT, geolon float64) (float64, error)

	// SidTime0 returns the sidereal time for Julian Date jd, ecliptic obliquity
	// eps and nutation nut at the Greenwich medidian, measured in hours.
	SidTime0(ut, eps, nut float64) float64
	// SidTime returns the sidereal time for Julian Date jd at the Greenwich
	// medidian, measured in hours.
	SidTime(ut float64) float64
}

// An Invoker invokes a function in an initialized execution context.
type Invoker interface {
	Invoke(fn func(Interface)) error
}
