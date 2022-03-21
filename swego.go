// Package swego defines an interface for interfacing with the Swiss Ephemeris.
package swego

// Error represents an error reported by the Swiss Ephemeris library.
type Error string

func (e Error) Error() string {
	return "swisseph: " + string(e)
}

// Planet is the type of planet constants.
type Planet int

// Ayanamsa is the type of sidereal mode constants.
type Ayanamsa int32

// SidMode represents library state changed by swe_set_sid_mode.
type SidMode struct {
	Mode   Ayanamsa
	T0     float64
	AyanT0 float64
}

// GeoLoc represents a geographic location.
type GeoLoc struct {
	Long float64
	Lat  float64
	Alt  float64
}

// CalcFlags represents the library state of swe_calc and swe_calc_ut.
type CalcFlags struct {
	Flags   int32
	TopoLoc *GeoLoc  // Arguments to swe_set_topo
	SidMode *SidMode // Arguments to swe_set_sid_mode
	JPLFile string   // Argument to swe_set_jpl_file
	DeltaT  *float64 // Argument to swe_set_delta_t_userdef, nil resets it.
}

// Copy returns a copy of the calculation flags fl.
func (fl *CalcFlags) Copy() *CalcFlags {
	copy := new(CalcFlags)
	*copy = *fl
	return copy
}

// Ephemeris represents an ephemeris implemented in the C library.
type Ephemeris int32

// SetEphemeris sets the ephemeris flag in fl.
func (fl *CalcFlags) SetEphemeris(eph Ephemeris) { fl.Flags |= int32(eph) }

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (fl *CalcFlags) SetDeltaT(f float64) { fl.DeltaT = &f }

// NodApsMethod is the type of Nodbit constants.
type NodApsMethod int32

// AyanamsaExFlags represents the library state of swe_get_ayanamsa_ex and
// swe_get_ayanamsa_ex_ut.
type AyanamsaExFlags struct {
	Flags   int32
	SidMode *SidMode // Argument to swe_set_sid_mode
	DeltaT  *float64 // Argument to swe_set_delta_t_userdef, nil resets it.
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (fl *AyanamsaExFlags) SetDeltaT(f float64) { fl.DeltaT = &f }

// CalType represents the calendar type used in julian date conversions.
type CalType int

// DateConvertFlags represents the library state of swe_utc_to_jd,
// swe_jdet_to_utc and swe_jdut1_to_utc.
type DateConvertFlags struct {
	Calendar CalType // clearifies the input year, Julian or Gregorian
	DeltaT   *float64
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (fl *DateConvertFlags) SetDeltaT(f float64) { fl.DeltaT = &f }

// HousesExFlags represents library state of swe_houses_ex in a stateless way.
type HousesExFlags struct {
	Flags   int32
	SidMode *SidMode // Argument to swe_set_sid_mode
	DeltaT  *float64 // Argument to swe_set_delta_t_userdef, nil resets it.
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (fl *HousesExFlags) SetDeltaT(f float64) { fl.DeltaT = &f }

// HSys represents house system identifiers used in the C library.
type HSys byte

// NewHSys validates the input and returns a HSys value if valid.
func NewHSys(c byte) (hsys HSys, ok bool) {
	if c == 'i' {
		return HSys(c), true
	}

	// It's trivial to convert lower case to upper case in ASCII.
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}

	switch c {
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y':
		return HSys(c), true
	default:
		return 0, false
	}
}

// TimeEquFlags represents the library state of swe_time_equ, swe_lmt_to_lat
// and swe_lat_to_lmt.
type TimeEquFlags struct {
	DeltaT *float64
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (fl *TimeEquFlags) SetDeltaT(f float64) { fl.DeltaT = &f }

// SidTimeFlags represents the library state of swe_sidtime0 and swe_sidtime.
type SidTimeFlags struct {
	DeltaT *float64
}

// SetDeltaT sets f as delta T in flags object fl.
// Set fl.DeltaT to nil to reset the value within the Swiss Ephemeris.
func (fl *SidTimeFlags) SetDeltaT(f float64) { fl.DeltaT = &f }

// Interface defines a standardized way for interfacing with the Swiss
// Ephemeris library from Go.
type Interface interface {
	// Version returns the version of the Swiss Ephemeris.
	Version() (string, error)

	// PlanetName returns the name of planet pl.
	PlanetName(pl Planet) (string, error)

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

	// GetAyanamsaEx returns the ayanamsa for Julian Date (in Ephemeris Time) et.
	// It is equal to GetAyanamsa but uses the ΔT consistent with the ephemeris
	// passed in fl.Flags.
	GetAyanamsaEx(et float64, fl *AyanamsaExFlags) (float64, error)
	// GetAyanamsaExUT returns the ayanamsa for Julian Date (in Universal Time) ut.
	// It is equal to GetAyanamsaUT but uses the ΔT consistent with the ephemeris
	// passed in fl.Flags.
	GetAyanamsaExUT(ut float64, fl *AyanamsaExFlags) (float64, error)
	// GetAyanamsaName returns the name of sidmode.
	GetAyanamsaName(ayan Ayanamsa) (string, error)

	// JulDay returns the corresponding Julian Date for the given date.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	JulDay(y, m, d int, h float64, ct CalType) (float64, error)
	// RevJul returns the corresponding calendar date for the given Julian Date.
	// Calendar type ct is used to clearify the year y, Julian or Gregorian.
	RevJul(jd float64, ct CalType) (y, m, d int, h float64, err error)
	// UTCToJD returns the corresponding Julian Date in Ephemeris and Universal
	// Time for the given date and accounts for leap seconds in the conversion.
	UTCToJD(y, m, d, h, i int, s float64, fl *DateConvertFlags) (et, ut float64, err error)
	// JdETToUTC returns the corresponding calendar date for the given Julian
	// Date in Ephemeris Time and accounts for leap seconds in the conversion.
	JdETToUTC(et float64, fl *DateConvertFlags) (y, m, d, h, i int, s float64, err error)
	// JdETToUTC returns the corresponding calendar date for the given Julian
	// Date in Universal Time and accounts for leap seconds in the conversion.
	JdUT1ToUTC(ut1 float64, fl *DateConvertFlags) (y, m, d, h, i int, s float64, err error)

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
	HouseName(hsys HSys) (string, error)

	// DeltaTEx returns the ΔT for the Julian Date jd.
	DeltaTEx(jd float64, eph Ephemeris) (float64, error)

	// TimeEqu returns the difference between local apparent and local mean time
	// in days for the given Julian Date (in Universal Time).
	TimeEqu(jd float64, fl *TimeEquFlags) (float64, error)
	// LMTToLAT returns the local apparent time for the given Julian Date (in
	// Local Mean Time) and the geographic longitude.
	LMTToLAT(jdLMT, geolon float64, fl *TimeEquFlags) (float64, error)
	// LATToLMT returns the local mean time for the given Julian Date (in Local
	// Apparent Time) and the geographic longitude.
	LATToLMT(jdLAT, geolon float64, fl *TimeEquFlags) (float64, error)

	// SidTime0 returns the sidereal time for Julian Date jd, ecliptic obliquity
	// eps and nutation nut at the Greenwich medidian, measured in hours.
	SidTime0(ut, eps, nut float64, fl *SidTimeFlags) (float64, error)
	// SidTime returns the sidereal time for Julian Date jd at the Greenwich
	// medidian, measured in hours.
	SidTime(ut float64, fl *SidTimeFlags) (float64, error)

	// SplitDeg takes a decimal degree number as input and provides sign or
	// nakshatra, degree, minutes, seconds and fraction of second.
	// Internally it calls swe_split_deg() and it's reference should be used.
	SplitDeg(ddeg float64, roundflag int) (ideg int32, imin int32, isec int32, dsecfr float64, isgn int32)
}

// Locked tries to exclusively lock the library handle, disable per function
// locking and exposes the locked interface to the callback function. Per
// function locking is restored when execution is returned to the caller. The
// input object is directly passed to the callback function if it does not
// implement ExclusiveLocker. If either argument is nil, it panics.
func Locked(swe Interface, callback func(swe Interface)) {
	if swe == nil {
		panic("swe is nil")
	}

	if callback == nil {
		panic("callback is nil")
	}

	l, ok := swe.(ExclusiveLocker)
	if !ok {
		callback(swe)
		return
	}

	li := l.ExclusiveLock()
	callback(li)
	li.ExclusiveUnlock()
}

// An ExclusiveLocker is a library handle that can exclusively lock itself.
type ExclusiveLocker interface {
	ExclusiveLock() LockedInterface
}

// An LockedInterface is a library handle that is exclusively locked.
type LockedInterface interface {
	Interface
	ExclusiveUnlock()
}
