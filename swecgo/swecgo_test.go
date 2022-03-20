//go:build (linux && cgo) || (darwin && cgo)
// +build linux,cgo darwin,cgo

package swecgo

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/howesteve/swego"
)

func TestVersionHeaderFile(t *testing.T) {
	const got = Version

	want := fmt.Sprintf("%d.%02d", VersionMajor, VersionMinor)
	if strings.Count(Version, ".") == 2 {
		want += fmt.Sprintf(".%02d", VersionPatch)
	}

	if got != want {
		t.Errorf("sweph.h = %q, sweversion.h = %q", got, want)
	}
}

func TestConstantCheck(t *testing.T) {
	if swego.FlagSidereal != flgSidereal {
		t.Errorf("swego.FlagSidereal = %d, flgSidereal = %d", swego.FlagSidereal, flgSidereal)
	}

	if swego.FlagTopo != flgTopo {
		t.Errorf("swego.FlagTopo = %d, flgTopo = %d", swego.FlagTopo, flgTopo)
	}
}

func inDelta(lhs, rhs, delta float64) bool {
	return math.Abs(lhs-rhs) < delta
}

func inDeltaSlice(lhs, rhs []float64, delta float64) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i, lhs := range lhs {
		if !inDelta(lhs, rhs[i], delta) {
			return false
		}
	}

	return true
}

var swe = Open()

func Test_wrapper_Version(t *testing.T) {
	t.Parallel()

	got, err := swe.Version()
	if err != nil {
		t.Errorf("err = %v, want: nil", err)
	}

	if got != Version {
		t.Errorf("Version() = %q, want: %q", got, Version)
	}
}

func Test_wrapper_Close(t *testing.T) {
	t.Parallel()
	Locked(swe, func(swe Library) {
		swe.SetPath(DefaultPath)
		swe.Close()
		swe.SetPath(DefaultPath)
	})
}

func Test_wrapper_PlanetName(t *testing.T) {
	t.Parallel()

	name, err := swe.PlanetName(swego.Sun)
	if err != nil {
		t.Errorf("err = %v, want: nil", err)
	}

	if name != "Sun" {
		t.Errorf("PlanetName(Sun) = %q, want: \"Sun\"", name)
	}
}

func Test_wrapper_Calc(t *testing.T) {
	t.Parallel()

	type result struct {
		xx  []float64
		cfl int
	}

	cases := []struct {
		fn   func(float64, swego.Planet, *swego.CalcFlags) ([]float64, int, error)
		in   *swego.CalcFlags
		want result
	}{
		{swe.Calc,
			nil,
			result{[]float64{279.858461, .000229, .983331, .0, .0, .0}, 0}},
		{swe.CalcUT,
			nil,
			result{[]float64{279.859216, .000229, .983331, .0, .0, .0}, 2}},
		{swe.Calc,
			&swego.CalcFlags{Flags: swego.FlagEphJPL},
			result{[]float64{279.858461, .000229, .983331, .0, .0, .0}, 1}},
		{swe.CalcUT,
			&swego.CalcFlags{Flags: swego.FlagEphJPL},
			result{[]float64{279.859216, .000229, .983331, .0, .0, .0}, 1}},
		{swe.Calc,
			&swego.CalcFlags{Flags: swego.FlagEphJPL, JPLFile: swego.FnameDft2},
			result{[]float64{279.858461, .000230, .983331, .0, .0, .0}, 1}},
		{swe.CalcUT,
			&swego.CalcFlags{Flags: swego.FlagEphJPL, JPLFile: swego.FnameDft2},
			result{[]float64{279.859216, .000230, .983331, .0, .0, .0}, 1}},
		{swe.Calc,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagTopo,
				TopoLoc: &swego.GeoLoc{Lat: 52.083333, Long: 5.116667, Alt: 0},
			},
			result{[]float64{279.858426, -.000966, .983369, .0, .0, .0}, 32772}},
		{swe.CalcUT,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagTopo,
				TopoLoc: &swego.GeoLoc{Lat: 52.083333, Long: 5.116667, Alt: 0},
			},
			result{[]float64{279.859186, -.000966, .983369, .0, .0, .0}, 32772}},
		{swe.Calc,
			&swego.CalcFlags{Flags: swego.FlagEphJPL | swego.FlagSidereal},
			result{[]float64{255.121938, .000229, .983331, .0, .0, .0}, 65604}},
		{swe.CalcUT,
			&swego.CalcFlags{Flags: swego.FlagEphJPL | swego.FlagSidereal},
			result{[]float64{255.12280449619868, 0.0002346766083462040, 0.9833318780303111, .0, .0, .0}, 65604}},
		{swe.Calc,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagSidereal,
				SidMode: &swego.SidMode{Mode: 1},
			},
			result{[]float64{256.005296, .000229, .983331, .0, .0, .0}, 65604}},
		{swe.CalcUT,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagSidereal,
				SidMode: &swego.SidMode{Mode: 1},
			},
			result{[]float64{256.0060121369246, 0.00023467660834620405, 0.983331878030311, .0, .0, .0}, 65604}},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			xx, cfl, err := c.fn(2451544.5, swego.Sun, c.in)
			if err != nil {
				t.Errorf("err = %v, want: nil", err)
			}

			if !inDeltaSlice(xx, c.want.xx, 1e-6) {
				t.Errorf("xx = %v ± 1e-6, want: %v", xx, c.want.xx)
			}

			if cfl != c.want.cfl {
				t.Errorf("cfl = %d, want: %d", cfl, c.want.cfl)
			}
		})
	}
}

func Test_wrapper_Calc_error(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fn  func(float64, swego.Planet, *swego.CalcFlags) ([]float64, int, error)
		err swego.Error
	}{
		{swe.Calc, "jd 99999999.000000 outside JPL eph. range -3027215.50 .. 7930192.50;"},
		{swe.CalcUT, "jd 100002682.057840 outside JPL eph. range -3027215.50 .. 7930192.50;"},
	}

	fl := &swego.CalcFlags{
		Flags: swego.FlagEphJPL,
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			xx, cfl, err := c.fn(99999999., swego.Sun, fl)
			if err != c.err {
				t.Errorf("err = %v, want: %q", err, c.err)
			}

			if !reflect.DeepEqual(xx, make([]float64, 6)) {
				t.Errorf("xx = %v, want: []float{}", xx)
			}

			if cfl != -1 {
				t.Errorf("cfl = %d, want: -1", cfl)
			}
		})
	}
}

func Test_wrapper_NodAps(t *testing.T) {
	t.Parallel()

	type result struct {
		nasc, ndsc, peri, aphe []float64
	}

	cases := []struct {
		fn   func(float64, swego.Planet, *swego.CalcFlags, swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error)
		in   swego.NodApsMethod
		want result
	}{
		{swe.NodAps, swego.NodbitMean, result{
			[]float64{125.067162, .0, .002461, .0, .0, .0},
			[]float64{305.067162, .0, .002671, .0, .0, .0},
			[]float64{83.408587, -3.425232, .002428, .0, .0, .0},
			[]float64{263.408587, 3.425232, .002710, .0, .0, .0},
		}},
		{swe.NodApsUT, swego.NodbitMean, result{
			[]float64{125.067123, .0, .002461, .0, .0, .0},
			[]float64{305.067123, .0, .002671, .0, .0, .0},
			[]float64{83.408669, -3.425224, .002428, .0, .0, .0},
			[]float64{263.408669, 3.425224, .002710, .0, .0, .0},
		}},
		{swe.NodAps, swego.NodbitMean | swego.NodbitFoPoint, result{
			[]float64{125.067162, .0, .002461, .0, .0, .0},
			[]float64{305.067162, .0, .002671, .0, .0, .0},
			[]float64{83.408587, -3.425232, .002428, .0, .0, .0},
			[]float64{263.408587, 3.425232, .000282, .0, .0, .0}, // different
		}},
		{swe.NodApsUT, swego.NodbitMean | swego.NodbitFoPoint, result{
			[]float64{125.067123, .0, .002461, .0, .0, .0},
			[]float64{305.067123, .0, .002671, .0, .0, .0},
			[]float64{83.408669, -3.425224, .002428, .0, .0, .0},
			[]float64{263.408669, 3.425224, .000282, .0, .0, .0}, // different
		}},
	}

	fl := &swego.CalcFlags{
		Flags: swego.FlagEphJPL,
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			nasc, ndsc, peri, aphe, err := c.fn(2451544.5, swego.Moon, fl, c.in)
			if err != nil {
				t.Errorf("err = %v, want: nil", err)
			}

			if !inDeltaSlice(nasc, c.want.nasc, 1e-6) {
				t.Errorf("nasc = %v ± 1e-6, want: %v", nasc, c.want.nasc)
			}

			if !inDeltaSlice(ndsc, c.want.ndsc, 1e-6) {
				t.Errorf("ndsc = %v ± 1e-6, want: %v", ndsc, c.want.ndsc)
			}

			if !inDeltaSlice(peri, c.want.peri, 1e-6) {
				t.Errorf("peri = %v ± 1e-6, want: %v", peri, c.want.peri)
			}

			if !inDeltaSlice(aphe, c.want.aphe, 1e-6) {
				t.Errorf("aphe = %v ± 1e-6, want: %v", aphe, c.want.aphe)
			}
		})
	}
}

func Test_wrapper_NodAps_error(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fn  func(float64, swego.Planet, *swego.CalcFlags, swego.NodApsMethod) (nasc, ndsc, peri, aphe []float64, err error)
		err swego.Error
	}{
		{swe.NodAps, "jd 99999999.000000 outside JPL eph. range -3027215.50 .. 7930192.50;"},
		{swe.NodApsUT, "jd 100002682.057840 outside JPL eph. range -3027215.50 .. 7930192.50;"},
	}

	fl := &swego.CalcFlags{
		Flags: swego.FlagEphJPL,
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			nasc, ndsc, peri, aphe, err := c.fn(99999999., swego.Moon, fl, swego.NodbitMean)
			if err != c.err {
				t.Errorf("err = %v, want: %q", err, c.err)
			}

			if !reflect.DeepEqual(nasc, make([]float64, 6)) {
				t.Errorf("nasc = %v, want: []float{}", nasc)
			}

			if !reflect.DeepEqual(ndsc, make([]float64, 6)) {
				t.Errorf("ndsc = %v, want: []float{}", ndsc)
			}

			if !reflect.DeepEqual(peri, make([]float64, 6)) {
				t.Errorf("peri = %v, want: []float{}", peri)
			}

			if !reflect.DeepEqual(aphe, make([]float64, 6)) {
				t.Errorf("aphe = %v, want: []float{}", aphe)
			}
		})
	}
}

func Test_wrapper_GetAyanamsaEx(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fn   func(float64, *swego.AyanamsaExFlags) (float64, error)
		want float64
	}{
		{swe.GetAyanamsaEx, 24.736411},
		{swe.GetAyanamsaExUT, 24.736411},
	}

	fl := &swego.AyanamsaExFlags{
		Flags:   1,
		SidMode: new(swego.SidMode),
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			got, err := c.fn(2451544.5, fl)
			if err != nil {
				t.Fatalf("err = %v, want: nil", err)
			}

			if !inDelta(got, c.want, 1e-6) {
				t.Errorf("deltaT = %f, want: %f", got, c.want)
			}
		})
	}
}

func Test_wrapper_GetAyanamsaName(t *testing.T) {
	t.Parallel()

	got, err := swe.GetAyanamsaName(0)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	const want = "Fagan/Bradley"
	if got != want {
		t.Errorf("GetAyanamsaName(0) = %q, want: %q", got, want)
	}
}

func Test_wrapper_JulDay(t *testing.T) {
	t.Parallel()

	got, err := swe.JulDay(2000, 1, 1, 0, swego.Gregorian)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if !inDelta(got, 2451544.5, 1e-6) {
		t.Errorf("JD = %f, want: 2451544.5", got)
	}
}

func Test_wrapper_RevJul(t *testing.T) {
	t.Parallel()

	y, m, d, h, err := swe.RevJul(2451544.5, swego.Gregorian)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if y != 2000 {
		t.Errorf("y = %d, want: 2000", y)
	}

	if m != 1 {
		t.Errorf("m = %d, want: 1", m)
	}

	if d != 1 {
		t.Errorf("d = %d, want: 1", d)
	}

	if h != 0 {
		t.Errorf("h = %f, want: 0", h)
	}
}

func Test_wrapper_RevJul_bce(t *testing.T) {
	t.Parallel()

	y, m, d, h, err := swe.RevJul(990574.5, swego.Gregorian)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if y != -2000 {
		t.Errorf("y = %d, want: -2000", y)
	}

	if m != 1 {
		t.Errorf("m = %d, want: 1", m)
	}

	if d != 1 {
		t.Errorf("d = %d, want: 1", d)
	}

	if h != 0 {
		t.Errorf("h = %f, want: 0", h)
	}
}

func Test_wrapper_UTCToJD(t *testing.T) {
	t.Parallel()

	fl := &swego.DateConvertFlags{Calendar: swego.Gregorian}
	et, ut, err := swe.UTCToJD(2000, 1, 1, 0, 0, 0, fl)
	got := []float64{et, ut}
	want := []float64{2451544.500743, 2451544.500004}

	if err != nil {
		t.Fatalf("err = nil, want: %q", err)
	}

	if !inDeltaSlice(got, want, 1e-6) {
		t.Errorf("[et, ut] = %v, want: %f", got, want)
	}
}

func Test_wrapper_JDETToUTC(t *testing.T) {
	t.Parallel()

	fl := &swego.DateConvertFlags{Calendar: swego.Gregorian}
	y, m, d, h, i, s, err := swe.JdETToUTC(2451544.5, fl)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if y != 1999 {
		t.Errorf("y = %d, want: 1999", y)
	}

	if m != 12 {
		t.Errorf("m = %d, want: 12", m)
	}

	if d != 31 {
		t.Errorf("d = %d, want: 31", d)
	}

	if h != 23 {
		t.Errorf("h = %d, want: 23", h)
	}

	if i != 58 {
		t.Errorf("i = %d, want: 58", i)
	}

	if !inDelta(s, 55.815999, 1e-6) {
		t.Errorf("s = %f ± 1e-6, want: 55.815999", s)
	}
}

func Test_wrapper_JDUT1ToUTC(t *testing.T) {
	t.Parallel()

	fl := &swego.DateConvertFlags{Calendar: swego.Gregorian}
	y, m, d, h, i, s, err := swe.JdUT1ToUTC(2451544.5, fl)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if y != 1999 {
		t.Errorf("y = %d, want: 1999", y)
	}

	if m != 12 {
		t.Errorf("m = %d, want: 12", m)
	}

	if d != 31 {
		t.Errorf("d = %d, want: 31", d)
	}

	if h != 23 {
		t.Errorf("h = %d, want: 23", h)
	}

	if i != 59 {
		t.Errorf("i = %d, want: 59", i)
	}

	const want = 59.644500
	if !inDelta(s, want, 1e-6) {
		t.Errorf("s = %f ± 1e-6, want: %f", s, want)
	}
}

func Test_wrapper_HousesEx(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   swego.HSys
		flags  *swego.HousesExFlags
	}

	type result struct {
		cusps []float64
		ascmc []float64
		err   error
	}

	cases := []struct {
		in   input
		want result
	}{
		{
			input{52.083333, swego.Placidus, nil},
			result{
				[]float64{0,
					190.553653, 215.538288, 246.822987, 283.886819, 319.373115, 348.152982,
					10.553653, 35.538288, 66.822987, 103.886819, 139.373115, 168.152982,
				},
				[]float64{
					190.553653, 103.886819, 105.080915, 24.306488,
					196.367263, 214.734661, 192.275917, 34.734661,
					.0, .0,
				},
				nil,
			}},
		{
			input{82.083333, swego.Koch, nil},
			result{
				[]float64{0,
					183.972931, 217.277560, 250.582190, 283.886819, 310.582190, 337.277560,
					3.972931, 37.277560, 70.582190, 103.886819, 130.582190, 157.277560,
				},
				[]float64{
					183.972931, 103.886819, 105.080915, 17.393326,
					196.367263, 352.493044, 195.452718, 172.493044,
					.0, .0,
				},
				swego.Error(""),
			}},
		{
			input{52.083333, swego.Gauquelin, nil},
			result{
				[]float64{0,
					190.553653, 183.704634, 176.258623, 168.152982, 159.330891, 149.746713,
					139.373115, 128.213369, 116.328153, 103.886819, 91.215325, 78.741976,
					66.822987, 55.634583, 45.214174, 35.538288, 26.565986, 18.252908,
					10.553653, 3.704634, 356.258623, 348.152982, 339.330891, 329.746713,
					319.373115, 308.213369, 296.328153, 283.886819, 271.215325, 258.741976,
					246.822987, 235.634583, 225.214174, 215.538288, 206.565986, 198.252908,
				},
				[]float64{
					190.553653, 103.886819, 105.080915, 24.306488,
					196.367263, 214.734661, 192.275917, 34.734661,
					.0, .0,
				},
				nil,
			}},
		{
			input{52.083333, swego.Placidus, &swego.HousesExFlags{
				Flags:   flgSidereal,
				SidMode: &swego.SidMode{},
			}},
			result{
				[]float64{0,
					165.817130, 190.801765, 222.086464, 259.150296, 294.636593, 323.416459,
					345.817130, 10.801765, 42.086464, 79.150296, 114.636593, 143.416459,
				},
				[]float64{
					165.817130, 79.150296, 105.080915, 359.569965,
					171.630740, 189.998138, 167.539394, 9.998138,
					.0, .0,
				},
				nil,
			}},
		// SunshineAlt is the only lower case house system letter.
		// It is introduced in Swiss Ephemeris version 2.05.
		{
			input{52.083333, swego.SunshineAlt, nil},
			result{
				[]float64{0,
					190.553653, 216.974402, 246.972073, 283.886819, 318.522699, 345.395439,
					10.553653, 26.204261, 52.666679, 103.886819, 151.275154, 175.374467,
				},
				[]float64{
					190.553653, 103.886819, 105.080915, 24.306488,
					196.367263, 214.734661, 192.275917, 34.734661,
					.0, -23.071122,
				},
				nil,
			}},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			cusps, ascmc, err := swe.HousesEx(2451544.5, c.in.flags, c.in.geolat, 5.116667, c.in.hsys)
			if err != c.want.err {
				t.Fatalf("(%f, %c) err = %v, want: %q",
					c.in.geolat, c.in.hsys, err, c.want.err)
			}

			if !inDeltaSlice(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps:\n\t%v\nwant:\n\t%v",
					c.in.geolat, c.in.hsys, cusps, c.want.cusps)
			}

			if !inDeltaSlice(ascmc, c.want.ascmc, 1e-6) {
				t.Errorf("(%f, %c) ascmc:\n\t%v\nwant:\n\t%v",
					c.in.geolat, c.in.hsys, ascmc, c.want.ascmc)
			}
		})
	}
}

func Test_wrapper_HousesArmc(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   swego.HSys
	}

	type result struct {
		cusps []float64
		ascmc []float64
		err   error
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, swego.Placidus}, result{
			[]float64{0,
				190.553489, 215.537888, 246.822499, 283.886657, 319.373244, 348.153088,
				10.553489, 35.537888, 66.822499, 103.886657, 139.373244, 168.153088,
			},
			[]float64{
				190.553489, 103.886657, 105.080916, 24.307632,
				196.367450, 214.737779, 192.275825, 34.737779,
				.0, .0,
			},
			nil,
		}},
		{input{82.083333, swego.Koch}, result{
			[]float64{0,
				183.972748, 217.277384, 250.582021, 283.886657, 310.582021, 337.277384,
				3.972748, 37.277384, 70.582021, 103.886657, 130.582021, 157.277384,
			},
			[]float64{
				183.972748, 103.886657, 105.080916, 17.393607,
				196.367450, 352.493777, 195.452830, 172.493777,
				.0, .0,
			},
			swego.Error(""),
		}},
		{input{52.083333, swego.Gauquelin}, result{
			[]float64{0,
				190.553489, 183.704585, 176.258665, 168.153088, 159.331033, 149.746863,
				139.373244, 128.213442, 116.328129, 103.886657, 91.215011, 78.741543,
				66.822499, 55.634096, 45.213720, 35.537888, 26.565652, 18.252653,
				10.553489, 3.704585, 356.258665, 348.153088, 339.331033, 329.746863,
				319.373244, 308.213442, 296.328129, 283.886657, 271.215011, 258.741543,
				246.822499, 235.634096, 225.213720, 215.537888, 206.565652, 198.252653,
			},
			[]float64{
				190.553489, 103.886657, 105.080916, 24.307632,
				196.367450, 214.737779, 192.275825, 34.737779,
				.0, .0,
			},
			nil,
		}},
		// SunshineAlt is the only lower case house system letter.
		// It is introduced in Swiss Ephemeris version 2.05.
		{input{52.083333, swego.SunshineAlt}, result{
			[]float64{0,
				190.553489, 213.007193, 243.038107, 283.886657, 322.037715, 349.034810,
				10.553489, 33.007193, 63.038107, 103.886657, 142.037715, 169.034810,
			},
			[]float64{
				190.553489, 103.886657, 105.080916, 24.307632,
				196.367450, 214.737779, 192.275825, 34.737779,
				.0, .0,
			},
			nil,
		}},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			cusps, ascmc, err := swe.HousesARMC(105.080916, c.in.geolat, 23.439279, c.in.hsys)
			if err != c.want.err {
				t.Fatalf("(%f, %c) err = %v, want: %q",
					c.in.geolat, c.in.hsys, err, c.want.err)
			}

			if !inDeltaSlice(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps:\n\t%v\nwant:\n\t%v",
					c.in.geolat, c.in.hsys, cusps, c.want.cusps)
			}

			if !inDeltaSlice(ascmc, c.want.ascmc, 1e-6) {
				t.Errorf("(%f, %c) ascmc:\n\t%v\nwant:\n\t%v",
					c.in.geolat, c.in.hsys, ascmc, c.want.ascmc)
			}
		})
	}
}

func Test_wrapper_HousePos(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   swego.HSys
	}

	type result struct {
		pos float64
		err string
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, swego.Placidus}, result{6.355326, ""}},
		{input{82.083333, swego.Koch}, result{6.355326, ""}},
		{input{52.083333, swego.Gauquelin}, result{20.934023, ""}},
		// SunshineAlt is the only lower case house system letter.
		// It is introduced in Swiss Ephemeris version 2.05.
		{input{52.083333, swego.SunshineAlt}, result{4.597296, ""}},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			// Houses, HousesEx or HousesARMC should be called before HousePos.
			_, _, _ = swe.HousesARMC(105.080916, c.in.geolat, 23.439279, c.in.hsys)
			pos, err := swe.HousePos(105.080916, c.in.geolat, 23.439279, c.in.hsys, 279.858461, 0.000229)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err = %v, want: %q",
					c.in.geolat, c.in.hsys, err, c.want.err)
			}

			if !inDelta(pos, c.want.pos, 1e-6) {
				t.Errorf("(%f, %c) pos = %f, want: %f",
					c.in.geolat, c.in.hsys, pos, c.want.pos)
			}
		})
	}
}

func Test_wrapper_HouseName(t *testing.T) {
	t.Parallel()

	got, err := swe.HouseName('P')
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	want := "Placidus"
	if got != want {
		t.Errorf("HouseName('P') = %q, want: %q", got, want)
	}
}

func Test_wrapper_DeltaTEx(t *testing.T) {
	t.Parallel()

	got, err := swe.DeltaTEx(2451544.5, swego.Swiss)
	if err != nil {
		t.Fatalf("err = %v, want: nil", err)
	}

	if !inDelta(got, 0.000739, 1e-6) {
		t.Errorf("DeltaTEx(2451544.5, Swiss) = %f, want: 0.000739", got)
	}
}

func Test_wrapper_SetDeltaTUserDef(t *testing.T) {
	t.Parallel()

	Locked(swe, func(swe Library) {
		const want = 63.8496
		fl := new(swego.DateConvertFlags)
		fl.SetDeltaT(want)
		swe.JdETToUTC(0, fl) // call swe_set_delta_t_userdef

		got, err := swe.DeltaTEx(2451544.5, swego.Swiss)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		if got != want {
			t.Errorf("user defined ΔT not set correctly; ΔT = %f, want: %f", got, want)
		}

		fl.DeltaT = nil      // reset ΔT library value
		swe.JdETToUTC(0, fl) // call swe_set_delta_t_userdef

		got, err = swe.DeltaTEx(2451544.5, swego.Swiss)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		const libValue = 0.000739
		if !inDelta(got, libValue, 1e-6) {
			t.Errorf("user defined ΔT not reset correctly; ΔT = %f, want: %f", got, libValue)
		}
	})
}

func Test_wrapper_TimeEqu(t *testing.T) {
	t.Parallel()

	Locked(swe, func(swe Library) {
		fl := new(swego.TimeEquFlags)
		fl.SetDeltaT(0)

		got, err := swe.TimeEqu(2451544.5, fl)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		if !inDelta(got, -0.002114, 1e-6) {
			t.Errorf("TimeEqu(2451544.5) = %f, want: -0.002114", got)
		}

		got, err = swe.TimeEqu(2451544.5, nil)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		if !inDelta(got, -0.002116, 1e-6) {
			t.Errorf("TimeEqu(2451544.5) = %f, want: -0.002116", got)
		}
	})
}

func Test_wrapper_LMTToLAT(t *testing.T) {
	t.Parallel()

	Locked(swe, func(swe Library) {
		fl := new(swego.TimeEquFlags)
		fl.SetDeltaT(0)

		got, err := swe.LMTToLAT(2451544.5, 5.116667, fl)
		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta(got, 2451544.497891, 1e-6) {
			t.Errorf("LMTToLAT(2451544.5, 5.116667) = %f, want: 2451544.497891", got)
		}

		got, err = swe.LMTToLAT(2451544.5, 5.116667, nil)
		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta(got, 2451544.497889, 1e-6) {
			t.Errorf("LMTToLAT(2451544.5, 5.116667) = %f, want: 2451544.497889", got)
		}
	})
}

func Test_wrapper_LATToLMT(t *testing.T) {
	t.Parallel()

	Locked(swe, func(swe Library) {
		fl := new(swego.TimeEquFlags)
		fl.SetDeltaT(0)

		got, err := swe.LATToLMT(2451544.5, 5.116667, fl)
		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta(got, 2451544.502110, 1e-6) {
			t.Errorf("LATToLMT(2451544.5, 5.116667) = %f, want: 2451544.502110", got)
		}

		got, err = swe.LATToLMT(2451544.5, 5.116667, nil)
		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta(got, 2451544.502112, 1e-6) {
			t.Errorf("LATToLMT(2451544.5, 5.116667) = %f, want: 2451544.502112", got)
		}
	})
}

func Test_wrapper_SidTime0(t *testing.T) {
	t.Parallel()

	Locked(swe, func(swe Library) {
		fl := new(swego.SidTimeFlags)
		fl.SetDeltaT(0)

		got, err := swe.SidTime0(2451544.5, 23.439279, -0.003869, fl)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		if !inDelta(got, 6.664283, 1e-6) {
			t.Errorf("SidTime0(2451544.5, 23.439279, -0.003869) = %f, want: 6.664283", got)
		}

		got, err = swe.SidTime0(2451544.5, 23.439279, -0.003869, nil)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		if !inDelta(got, 6.664283, 1e-6) {
			t.Errorf("SidTime0(2451544.5, 23.439279, -0.003869) = %f, want: 6.664283", got)
		}
	})
}

func Test_wrapper_SidTime(t *testing.T) {
	t.Parallel()

	Locked(swe, func(swe Library) {
		fl := new(swego.SidTimeFlags)
		fl.SetDeltaT(0)

		got, err := swe.SidTime(2451544.5, fl)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		if !inDelta(got, 6.664283, 1e-6) {
			t.Errorf("SidTime(2451544.5) = %f, want: 6.664283", got)
		}

		got, err = swe.SidTime(2451544.5, nil)
		if err != nil {
			t.Fatalf("err = %v, want: nil", err)
		}

		if !inDelta(got, 6.664283, 1e-6) {
			t.Errorf("SidTime(2451544.5) = %f, want: 6.664283", got)
		}
	})
}
