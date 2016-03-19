package swego

// Calendar types.
const (
	Julian    CalType = 0
	Gregorian CalType = 1
)

// Planet, fictional body and asteroid constants.
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
	SidmUser               Ayanamsa = 255
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
	Equal                HSys = 'E' // also 'A'
	Alcabitus            HSys = 'B'
	Campanus             HSys = 'C'
	Gauquelin            HSys = 'G'
	Azimuthal            HSys = 'H'
	Koch                 HSys = 'K'
	Morinus              HSys = 'M'
	Porphyrius           HSys = 'O'
	Placidus             HSys = 'P'
	Regiomontanus        HSys = 'R'
	PolichPage           HSys = 'T'
	KrusinskiPisaGoelzer HSys = 'U'
	VehlowEqual          HSys = 'V'
	WholeSign            HSys = 'W'
	AxialRotation        HSys = 'X'
	APCHouses            HSys = 'Y'
)
