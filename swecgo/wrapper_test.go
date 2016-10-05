package swecgo

import (
	"math"
	"reflect"
	"testing"

	"github.com/dwlnetnl/swego"
)

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

func Test_wrapper_Version(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.Version()

		if got != Version {
			t.Errorf("Version() = %q, want: %q", got, Version)
		}
	})
}

func Test_wrapper_GetLibraryPath(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		libpath := swe.GetLibraryPath()
		if libpath == "" {
			t.Error(`libpath == ""`)
		}
	})
}

func Test_wrapper_Close(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		swe.Close()
		swe.SetPath(DefaultPath)
	})
}

func Test_wrapper_PlanetName(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		name := swe.PlanetName(swego.Sun)
		if name != "Sun" {
			t.Errorf("PlanetName(Sun) = %q, want: \"Sun\"", name)
		}
	})
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
		{gWrapper.Calc,
			&swego.CalcFlags{Flags: swego.FlagEphJPL},
			result{[]float64{279.858461, .000229, .983331, .0, .0, .0}, 1}},
		{gWrapper.CalcUT,
			&swego.CalcFlags{Flags: swego.FlagEphJPL},
			result{[]float64{279.859214, .000229, .983331, .0, .0, .0}, 1}},
		{gWrapper.Calc,
			&swego.CalcFlags{Flags: swego.FlagEphJPL, FileNameJPL: swego.FnameDft2},
			result{[]float64{279.858461, .000230, .983331, .0, .0, .0}, 1}},
		{gWrapper.CalcUT,
			&swego.CalcFlags{Flags: swego.FlagEphJPL, FileNameJPL: swego.FnameDft2},
			result{[]float64{279.859214, .000230, .983331, .0, .0, .0}, 1}},
		{gWrapper.Calc,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagTopo,
				TopoLoc: &swego.GeoLoc{Lat: 52.083333, Long: 5.116667, Alt: 0},
			},
			result{[]float64{279.858426, -.000966, .983369, .0, .0, .0}, 32769}},
		{gWrapper.CalcUT,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagTopo,
				TopoLoc: &swego.GeoLoc{Lat: 52.083333, Long: 5.116667, Alt: 0},
			},
			result{[]float64{279.859186, -.000966, .983369, .0, .0, .0}, 32769}},
		{gWrapper.Calc,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagSidereal,
				SidMode: &swego.SidMode{},
			},
			result{[]float64{255.121938, .000229, .983331, .0, .0, .0}, 65601}},
		{gWrapper.CalcUT,
			&swego.CalcFlags{
				Flags:   swego.FlagEphJPL | swego.FlagSidereal,
				SidMode: &swego.SidMode{},
			},
			result{[]float64{255.122691, .000229, .983331, .0, .0, .0}, 65601}},
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			xx, cfl, err := c.fn(2451544.5, swego.Sun, c.in)

			if err != nil {
				t.Errorf("err = %q, want: nil", err)
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
		err string
	}{
		{gWrapper.Calc, "swecgo: jd 99999999.000000 outside JPL eph. range -3027215.50 .. 7930192.50;"},
		{gWrapper.CalcUT, "swecgo: jd 100002561.779707 outside JPL eph. range -3027215.50 .. 7930192.50;"},
	}

	fl := &swego.CalcFlags{
		Flags: swego.FlagEphJPL,
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			xx, cfl, err := c.fn(99999999., swego.Sun, fl)

			if err == nil {
				t.Fatal("err == nil")
			}

			if err.Error() != c.err {
				t.Errorf("err = %q, want: %q", c.err, err)
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
		{gWrapper.NodAps, swego.NodbitMean, result{
			[]float64{125.067162, .0, .002461, .0, .0, .0},
			[]float64{305.067162, .0, .002671, .0, .0, .0},
			[]float64{83.408587, -3.425232, .002428, .0, .0, .0},
			[]float64{263.408587, 3.425232, .002710, .0, .0, .0},
		}},
		{gWrapper.NodApsUT, swego.NodbitMean, result{
			[]float64{125.067123, .0, .002461, .0, .0, .0},
			[]float64{305.067123, .0, .002671, .0, .0, .0},
			[]float64{83.408669, -3.425224, .002428, .0, .0, .0},
			[]float64{263.408669, 3.425224, .002710, .0, .0, .0},
		}},
		{gWrapper.NodAps, swego.NodbitMean | swego.NodbitFoPoint, result{
			[]float64{125.067162, .0, .002461, .0, .0, .0},
			[]float64{305.067162, .0, .002671, .0, .0, .0},
			[]float64{83.408587, -3.425232, .002428, .0, .0, .0},
			[]float64{263.408587, 3.425232, .000282, .0, .0, .0}, // different
		}},
		{gWrapper.NodApsUT, swego.NodbitMean | swego.NodbitFoPoint, result{
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
		Call(nil, func(_ swego.Interface) {
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
		err string
	}{
		{gWrapper.NodAps, "swecgo: jd 99999999.000000 outside JPL eph. range -3027215.50 .. 7930192.50;"},
		{gWrapper.NodApsUT, "swecgo: jd 100002561.779707 outside JPL eph. range -3027215.50 .. 7930192.50;"},
	}

	fl := &swego.CalcFlags{
		Flags: swego.FlagEphJPL,
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			nasc, ndsc, peri, aphe, err := c.fn(99999999., swego.Moon, fl, swego.NodbitMean)

			if err == nil {
				t.Fatalf("err = nil, want: %v", c.err)
			}

			if err.Error() != c.err {
				t.Errorf("err = %q, want: %q", err, c.err)
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

func Test_wrapper_GetAyanamsa(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fn   func(float64, *swego.SidMode) float64
		want float64
	}{
		{gWrapper.GetAyanamsa, 24.740393},
		{gWrapper.GetAyanamsaUT, 24.740393},
	}

	fl := &swego.SidMode{}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			got := c.fn(2451544.5, fl)

			if math.Abs(got-c.want) >= 1e-6 {
				t.Errorf("deltaT = %f, want: %f", got, c.want)
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
		{gWrapper.GetAyanamsaEx, 24.740393},
		{gWrapper.GetAyanamsaExUT, 24.740393},
	}

	fl := &swego.AyanamsaExFlags{
		Flags:   1,
		SidMode: &swego.SidMode{},
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			got, err := c.fn(2451544.5, fl)

			if err != nil {
				t.Fatalf("err = %q, want: nil", err)
			}

			if !inDelta(got, c.want, 1e-6) {
				t.Errorf("deltaT = %f, want: %f", got, c.want)
			}
		})
	}
}

func Test_wrapper_GetAyanamsaName(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.GetAyanamsaName(0)
		want := "Fagan/Bradley"

		if got != want {
			t.Errorf("GetAyanamsaName(0) = %q, want: %q", got, want)
		}
	})
}

func Test_wrapper_JulDay(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.JulDay(2000, 1, 1, 0, swego.Gregorian)

		if !inDelta(got, 2451544.5, 1e-6) {
			t.Errorf("JD = %f, want: 2451544.5", got)
		}
	})
}

func Test_wrapper_RevJul(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		y, m, d, h := swe.RevJul(2451544.5, swego.Gregorian)

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
	})
}

func Test_wrapper_RevJul_bce(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		y, m, d, h := swe.RevJul(990574.5, swego.Gregorian)

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
	})
}

func Test_wrapper_UTCToJD(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		et, ut, err := swe.UTCToJD(2000, 1, 1, 0, 0, 0, swego.Gregorian)

		got := []float64{et, ut}
		want := []float64{2451544.500743, 2451544.500004}

		if err != nil {
			t.Fatalf("err = nil, want: %q", err)
		}

		if !inDeltaSlice(got, want, 1e-6) {
			t.Errorf("[et, ut] = %v, want: %f", got, want)
		}
	})
}

func Test_wrapper_JDETToUTC(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		y, m, d, h, i, s := swe.JdETToUTC(2451544.5, swego.Gregorian)

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
	})
}

func Test_wrapper_JDUT1ToUTC(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		y, m, d, h, i, s := swe.JdUT1ToUTC(2451544.5, swego.Gregorian)

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

		if !inDelta(s, 59.645586, 1e-6) {
			t.Errorf("s = %f ± 1e-6, want: 59.645586", s)
		}
	})
}

func Test_wrapper_Houses(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   swego.HSys
	}

	type result struct {
		cusps []float64
		ascmc []float64
		err   string
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, swego.Placidus}, result{
			[]float64{0,
				190.553653, 215.538288, 246.822987, 283.886819, 319.373115, 348.152982,
				10.553653, 35.538288, 66.822987, 103.886819, 139.373115, 168.152982,
			},
			[]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
				.0, .0,
			},
			"",
		}},
		{input{82.083333, swego.Koch}, result{
			[]float64{0,
				183.972931, 217.277560, 250.582190, 283.886819, 310.582190, 337.277560,
				3.972931, 37.277560, 70.582190, 103.886819, 130.582190, 157.277560,
			},
			[]float64{
				183.972931, 103.886819, 105.080915, 17.393326,
				196.367263, 352.493044, 195.452718, 172.493044,
				.0, .0,
			},
			"swecgo: error calculating houses",
		}},
		{input{52.083333, swego.Gauquelin}, result{
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
			"",
		}},
		// SunshineAlt is the only lower case house system letter.
		// It is introduced in Swiss Ephemeris version 2.05.
		{input{52.083333, swego.SunshineAlt}, result{
			[]float64{0,
				190.553653, 216.974402, 246.972073, 283.886819, 318.522699, 345.395439,
				10.553653, 26.204261, 52.666679, 103.886819, 151.275154, 175.374467,
			},
			[]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
				.0, .0,
			},
			"",
		}},
	}

	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			cusps, ascmc, err := swe.Houses(2451544.5, c.in.geolat, 5.116667, c.in.hsys)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err = %q, want: %q",
					c.in.geolat, c.in.hsys, err, c.want.err)
			}

			if !inDeltaSlice(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps = %v, want: %v",
					c.in.geolat, c.in.hsys, cusps, c.want.cusps)
			}

			if !inDeltaSlice(ascmc, c.want.ascmc, 1e-6) {
				t.Errorf("(%f, %c) ascmc = %v, want: %v",
					c.in.geolat, c.in.hsys, ascmc, c.want.ascmc)
			}
		})
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
		err   string
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, swego.Placidus, &swego.HousesExFlags{}}, result{
			[]float64{0,
				190.553653, 215.538288, 246.822987, 283.886819, 319.373115, 348.152982,
				10.553653, 35.538288, 66.822987, 103.886819, 139.373115, 168.152982,
			},
			[]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
				.0, .0,
			},
			"",
		}},
		{input{82.083333, swego.Koch, &swego.HousesExFlags{}}, result{
			[]float64{0,
				183.972931, 217.277560, 250.582190, 283.886819, 310.582190, 337.277560,
				3.972931, 37.277560, 70.582190, 103.886819, 130.582190, 157.277560,
			},
			[]float64{
				183.972931, 103.886819, 105.080915, 17.393326,
				196.367263, 352.493044, 195.452718, 172.493044,
				.0, .0,
			},
			"swecgo: error calculating houses",
		}},
		{input{52.083333, swego.Gauquelin, &swego.HousesExFlags{}}, result{
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
			"",
		}},
		{input{52.083333, swego.Placidus, &swego.HousesExFlags{
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
				"",
			}},
		// SunshineAlt is the only lower case house system letter.
		// It is introduced in Swiss Ephemeris version 2.05.
		{input{52.083333, swego.SunshineAlt, &swego.HousesExFlags{}}, result{
			[]float64{0,
				190.553653, 216.974402, 246.972073, 283.886819, 318.522699, 345.395439,
				10.553653, 26.204261, 52.666679, 103.886819, 151.275154, 175.374467,
			},
			[]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
				.0, .0,
			},
			"",
		}},
	}

	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			cusps, ascmc, err := swe.HousesEx(2451544.5, c.in.flags, c.in.geolat,
				5.116667, c.in.hsys)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err = %q, want: %q",
					c.in.geolat, c.in.hsys, err, c.want.err)
			}

			if !inDeltaSlice(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps = %v, want: %v",
					c.in.geolat, c.in.hsys, cusps, c.want.cusps)
			}

			if !inDeltaSlice(ascmc, c.want.ascmc, 1e-6) {
				t.Errorf("(%f, %c) ascmc = %v, want: %v",
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
		err   string
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
			"",
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
			"swecgo: error calculating houses",
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
			"",
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
			"",
		}},
	}

	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			cusps, ascmc, err := swe.HousesARMC(105.080916, c.in.geolat, 23.439279, c.in.hsys)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err = %q, want: %q",
					c.in.geolat, c.in.hsys, err, c.want.err)
			}

			if !inDeltaSlice(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps = %v, want: %v",
					c.in.geolat, c.in.hsys, cusps, c.want.cusps)
			}

			if !inDeltaSlice(ascmc, c.want.ascmc, 1e-6) {
				t.Errorf("(%f, %c) ascmc = %v, want: %v",
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
		{input{52.083333, swego.SunshineAlt}, result{6.509577, ""}},
	}
	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			// Houses, HousesEx or HousesARMC should be called before HousePos.
			_, _, _ = swe.HousesARMC(105.080916, c.in.geolat, 23.439279, c.in.hsys)

			pos, err := swe.HousePos(105.080916, c.in.geolat, 23.439279, c.in.hsys,
				279.858461, 0.000229)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err = %q, want: %q",
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

	Call(nil, func(swe swego.Interface) {
		got := swe.HouseName('P')
		want := "Placidus"

		if got != want {
			t.Errorf("HouseName('P') = %q, want: %q", got, want)
		}
	})
}

func Test_wrapper_DeltaT(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.DeltaT(2451544.5)

		if !inDelta(got, 0.000739, 1e-6) {
			t.Errorf("DeltaT(2451544.5) = %f, want: 0.000739", got)
		}
	})
}

func Test_wrapper_DeltaTEx(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.DeltaTEx(2451544.5, swego.Swiss)

		if err != nil {
			t.Fatalf("err = %q, want: nil", err)
		}

		if !inDelta(got, 0.000739, 1e-6) {
			t.Errorf("DeltaTEx(2451544.5, Swiss) = %f, want: 0.000739", got)
		}
	})
}

func Test_wrapper_SetDeltaTUserDef(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		want := 63.8496
		swe.SetDeltaTUserDef(want)

		// Normally you should not use DeltaT but in this test case the call to
		// swe_set_delta_t_userdef() is tested, not the ΔT value.
		got := swe.DeltaT(0.)

		if got != want {
			t.Errorf("user defined ΔT not set correctly; ΔT = %f, want: %f", got, want)
		}

		swe.SetDeltaTUserDef(swego.ResetDeltaT) // restore ΔT calculation
	})
}

func Test_wrapper_TimeEqu(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.TimeEqu(2451544.5)

		if err != nil {
			t.Fatalf("err = %q, want: nil", err)
		}

		if !inDelta(got, -0.002116, 1e-6) {
			t.Errorf("TimeEqu(2451544.5) = %f, want: -0.002116", got)
		}
	})
}

func Test_wrapper_LMTToLAT(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.LMTToLAT(2451544.5, 5.116667)

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

	Call(nil, func(swe swego.Interface) {
		got, err := swe.LATToLMT(2451544.5, 5.116667)

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

	Call(nil, func(swe swego.Interface) {
		got := swe.SidTime0(2451544.5, 23.439279, -0.003869)

		if !inDelta(got, 6.664283, 1e-6) {
			t.Errorf("SidTime0(2451544.5, 23.439279, -0.003869) = %f, want: 6.664283", got)
		}
	})
}

func Test_wrapper_SidTime(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.SidTime(2451544.5)

		if !inDelta(got, 6.664283, 1e-6) {
			t.Errorf("SidTime(2451544.5) = %f, want: 6.664283", got)
		}
	})
}
