package swego

// Calendar types defined in swephexp.h.
const (
	Julian    CalType = 0
	Gregorian CalType = 1
)

// Planet, fictional body and asteroid constants defined in swephexp.h.
const (
	Sun          Planet = 0
	Moon         Planet = 1
	Mercury      Planet = 2
	Venus        Planet = 3
	Mars         Planet = 4
	Jupiter      Planet = 5
	Saturn       Planet = 6
	Uranus       Planet = 7
	Neptune      Planet = 8
	Pluto        Planet = 9
	MeanNode     Planet = 10
	TrueNode     Planet = 11
	MeanApogee   Planet = 12
	OscuApogee   Planet = 13
	Earth        Planet = 14
	Chiron       Planet = 15
	Pholus       Planet = 16
	Ceres        Planet = 17
	Pallas       Planet = 18
	Juno         Planet = 19
	Vesta        Planet = 20
	InterApogee  Planet = 21
	InterPerigee Planet = 22

	Varuna Planet = AstOffset + 20000
	Nessus Planet = AstOffset + 7066

	Cupido   Planet = 40
	Hades    Planet = 41
	Zeus     Planet = 42
	Kronos   Planet = 43
	Apollon  Planet = 44
	Admetos  Planet = 45
	Vulkanus Planet = 46
	Poseidon Planet = 47

	Isis             Planet = 48
	Nibiru           Planet = 49
	Harrington       Planet = 50
	NeptuneLeverrier Planet = 51
	NeptuneAdams     Planet = 52
	PlutoLowell      Planet = 53
	PlutoPickering   Planet = 54
	Vulcan           Planet = 55
	WhiteMoon        Planet = 56
	Proserpina       Planet = 57
	Waldemath        Planet = 58

	EclNut Planet = -1

	AstOffset = 10000
)

//go:generate stringer -type=Planet

// Indexes of related house positions defined in swephexp.h.
const (
	Asc    = 0
	MC     = 1
	ARMC   = 2
	Vertex = 3
	EquAsc = 4 // "equatorial ascendant"
	CoAsc1 = 5 // "co-ascendant" (W. Koch)
	CoAsc2 = 6 // "co-ascendant" (M. Munkasey)
	PolAsc = 7 // "polar ascendant" (M. Munkasey)
)

// Ephemerides that are implemented in the C library.
const (
	JPL        Ephemeris = FlagEphJPL
	Swiss      Ephemeris = FlagEphSwiss
	Moshier    Ephemeris = FlagEphMoshier
	DefaultEph Ephemeris = FlagEphDefault
)

// Calculation flags defined in swephexp.h.
const (
	FlagEphJPL       = 1 << 0
	FlagEphSwiss     = 1 << 1
	FlagEphMoshier   = 1 << 2
	FlagHelio        = 1 << 3
	FlagTruePos      = 1 << 4
	FlagJ2000        = 1 << 5
	FlagNoNut        = 1 << 6
	FlagSpeed        = 1 << 8
	FlagNoGDefl      = 1 << 9
	FlagNoAbber      = 1 << 10
	FlagAstrometric  = FlagNoAbber | FlagNoGDefl
	FlagEquatorial   = 1 << 11
	FlagXYZ          = 1 << 12
	FlagRadians      = 1 << 13
	FlagBary         = 1 << 14
	FlagTopo         = 1 << 15
	FlagSidereal     = 1 << 16
	FlagICRS         = 1 << 17
	FlagJPLHor       = 1 << 18
	FlagJPLHorApprox = 1 << 19
	FlagEphDefault   = FlagEphSwiss
)

// Sidereal modes (ayanamsas) implemented in the C library.
const (
	SidmFaganBradley       Ayanamsa = 0
	SidmLahiri             Ayanamsa = 1
	SidmDeluce             Ayanamsa = 2
	SidmRaman              Ayanamsa = 3
	SidmUshashashi         Ayanamsa = 4
	SidmKrishnamurti       Ayanamsa = 5
	SidmDjwhalKhul         Ayanamsa = 6
	SidmYukteshwar         Ayanamsa = 7
	SidmJNBhasin           Ayanamsa = 8
	SidmBabylKruger1       Ayanamsa = 9
	SidmBabylKruger2       Ayanamsa = 10
	SidmBabylKruger3       Ayanamsa = 11
	SidmBabylHuber         Ayanamsa = 12
	SidmBabylEtaPiscium    Ayanamsa = 13
	SidmAldebaran15Tau     Ayanamsa = 14
	SidmHipparchos         Ayanamsa = 15
	SidmSassanian          Ayanamsa = 16
	SidmGalCent0Sag        Ayanamsa = 17
	SidmJ2000              Ayanamsa = 18
	SidmJ1900              Ayanamsa = 19
	SidmB1950              Ayanamsa = 20
	SidmSuryasiddhanta     Ayanamsa = 21
	SidmSuryasiddhantaMSun Ayanamsa = 22
	SidmAryabhata          Ayanamsa = 23
	SidmAryabhataMSun      Ayanamsa = 24
	SidmSSRevati           Ayanamsa = 25
	SidmSSCitra            Ayanamsa = 26
	SidmTrueCitra          Ayanamsa = 27
	SidmTrueRevati         Ayanamsa = 28
	SidmTruePushya         Ayanamsa = 29
	SidmGalCentGilBrand    Ayanamsa = 30
	SidmGalAlignMardyks    Ayanamsa = 31
	SidmGalEquIAU1958      Ayanamsa = 32
	SidmGalEquTrue         Ayanamsa = 33
	SidmGalEquMula         Ayanamsa = 34
	SidmGalTrueMula        Ayanamsa = 35
	SidmGalCentMulaWilhelm Ayanamsa = 36
	SidmAryabhata522       Ayanamsa = 37
	SidmBabylBritton       Ayanamsa = 38
	SidmUser               Ayanamsa = 255
)

// Options that augment a sidereal mode (ayanamsa).
const (
	SidbitEclT0    Ayanamsa = 256
	SidbitSSYPlane Ayanamsa = 512
	SidbitUserUT   Ayanamsa = 1024
)

// Nodes and apsides calculation bits defined in swephexp.h.
const (
	NodbitMean     NodApsMethod = 1
	NodbitOscu     NodApsMethod = 2
	NodbitOscuBary NodApsMethod = 4
	NodbitFoPoint  NodApsMethod = 256
)

// File name of JPL data files defined in swephexp.h.
const (
	FnameDE200 = "de200.eph"
	FnameDE406 = "de406.eph"
	FnameDE431 = "de431.eph"
	FnameDft   = FnameDE431
	FnameDft2  = FnameDE406
)

// House systems implemented in the C library.
const (
	Alcabitius           HSys = 'B'
	Campanus             HSys = 'C'
	EqualMC              HSys = 'D' // Equal houses, where cusp 10 = MC
	Equal                HSys = 'E' // also 'A'
	CarterPoliEquatorial HSys = 'F'
	Gauquelin            HSys = 'G'
	Azimuthal            HSys = 'H' // a.k.a Horizontal
	Sunshine             HSys = 'I' // Makransky, solution Treindl
	SunshineAlt          HSys = 'i' // Makransky, solution Makransky
	Koch                 HSys = 'K'
	PullenSD             HSys = 'L'
	Morinus              HSys = 'M'
	EqualAsc             HSys = 'N' // Equal houses, where cusp 1 = 0° Aries
	Porphyrius           HSys = 'O' // a.k.a Porphyry
	Placidus             HSys = 'P'
	PullenSR             HSys = 'Q'
	Regiomontanus        HSys = 'R'
	Sripati              HSys = 'S'
	PolichPage           HSys = 'T' // a.k.a. Topocentric
	KrusinskiPisaGoelzer HSys = 'U'
	VehlowEqual          HSys = 'V' // Equal Vehlow (Asc in middle of house 1)
	WholeSign            HSys = 'W'
	AxialRotation        HSys = 'X' // a.k.a. Meridian
	APCHouses            HSys = 'Y'
)

// SplitDeg() flags.
// See documentation on swe_split_deg() for more details
const (
	SplitDegRoundSec  = 1
	SplitDegRoundMin  = 2
	SplitDegRoundDeg  = 4
	SplitDegZodiacal  = 8    // split into zodiac signs
	SplitDegNakshatra = 1024 // split into nakshatras
	SplitDegKeepSign  = 16   // don't round to next zodiac sign/nakshatra
	SplitDegKeepDeg   = 32   // don't round to next degree
)
