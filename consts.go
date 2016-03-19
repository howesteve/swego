package swego

// Calendar types.
const (
	Julian    CalType = 0
	Gregorian CalType = 1
)

// Planet, fictional body and asteroid constants.
const (
	Sun          = 0
	Moon         = 1
	Mercury      = 2
	Venus        = 3
	Mars         = 4
	Jupiter      = 5
	Saturn       = 6
	Uranus       = 7
	Neptune      = 8
	Pluto        = 9
	MeanNode     = 10
	TrueNode     = 11
	MeanApogee   = 12
	OscuApogee   = 13
	Earth        = 14
	Chiron       = 15
	Pholus       = 16
	Ceres        = 17
	Pallas       = 18
	Juno         = 19
	Vesta        = 20
	InterApogee  = 21
	InterPerigee = 22

	Varuna = AstOffset + 20000
	Nessus = AstOffset + 7066

	Cupido   = 40
	Hades    = 41
	Zeus     = 42
	Kronos   = 43
	Apollon  = 44
	Admetos  = 45
	Vulkanus = 46
	Poseidon = 47

	Isis             = 48
	Nibiru           = 49
	Harrington       = 50
	NeptuneLeverrier = 51
	NeptuneAdams     = 52
	PlutoLowell      = 53
	PlutoPickering   = 54
	Vulcan           = 55
	WhiteMoon        = 56
	Proserpina       = 57
	Waldemath        = 58

	AstOffset = 10000
)

// Indexes of related house positions.
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

// Calculation flags.
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
	FlagEquatorial   = 1 << 11
	FlagXYZ          = 1 << 11
	FlagRadians      = 1 << 12
	FlagBary         = 1 << 13
	FlagTopo         = 1 << 14
	FlagSidereal     = 1 << 15
	FlagICRS         = 1 << 16
	FlagJPLHor       = 1 << 17
	FlagJPLHorApprox = 1 << 18
	FlagEphDefault   = FlagEphSwiss
)

// Sidereal modes (ayanamsas).
const (
	SidmFaganBradley = iota
	SidmLahiri
	SidmDeluce
	SidmRaman
	SidmUshashashi
	SidmKrishnamurti
	SidmDjwhalKhul
	SidmYukteshwar
	SidmJNBhasin
	SidmBabylKruger1
	SidmBabylKruger2
	SidmBabylKruger3
	SidmBabylHuber
	SidmBabylEtaPiscium
	SidmAldebaran15Tau
	SidmHipparchos
	SidmSassanian
	SidmGalCent0Sag
	SidmJ2000
	SidmJ1900
	SidmB1950
	SidmSuryasiddhanta
	SidmSuryasiddhantaMSun
	SidmAryabhata
	SidmAryabhataMSun
	SidmSSRevati
	SidmSSCitra
	SidmTrueCitra
	SidmTrueRevati
	SidmTruePushya
	SidmUser    = 255
	SidmDefault = SidmFaganBradley
)

// // Nodes and apsides calculation bits.
// const (
// 	NodBitMean       = 1
// 	NodBitOscu       = 2
// 	NodBitOscuBary   = 4
// 	NodBitFocalPoint = 256
// )

// File name of JPL data files.
const (
	FnameDE200 = "de200.eph"
	FnameDE406 = "de406.eph"
	FnameDE431 = "de431.eph"
	FnameDft   = FnameDE431
	FnameDft2  = FnameDE406
)

// House system constants.
const (
	Equal                = 'E' // also 'A'
	Alcabitus            = 'B'
	Campanus             = 'C'
	Gauquelin            = 'G'
	Azimuthal            = 'H'
	Koch                 = 'K'
	Morinus              = 'M'
	Porphyrius           = 'O'
	Placidus             = 'P'
	Regiomontanus        = 'R'
	PolichPage           = 'T'
	KrusinskiPisaGoelzer = 'U'
	VehlowEqual          = 'V'
	WholeSign            = 'W'
	AxialRotation        = 'X'
	APCHouses            = 'Y'
)
