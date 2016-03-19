package swecgo

import (
	"testing"

	"github.com/dwlnetnl/swego"
)

func inDelta(a, b []float64, delta float64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, lhs := range a {
		d := lhs - b[i]

		if d < -delta || d > delta {
			return false
		}
	}

	return true
}

func TestVersion(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.Version()

		if got != Version {
			t.Errorf("Version() != %s, got: %q", Version, got)
		}
	})
}

func TestClose(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		swe.Close()
		swe.SetPath(DefaultPath)
	})
}

func TestCalc(t *testing.T) {
	t.Parallel()

	type result struct {
		xx  [6]float64
		cfl int
	}

	loc := swego.TopoLoc{Lat: 52.083333, Long: 5.116667, Alt: 0}

	cases := []struct {
		fn   func(float64, int, swego.CalcFlags) ([6]float64, int, error)
		in   swego.CalcFlags
		want result
	}{
		{gWrapper.Calc,
			swego.CalcFlags{Flags: 1},
			result{[6]float64{279.858461, .000229, .983331}, 1}},
		{gWrapper.CalcUT,
			swego.CalcFlags{Flags: 1},
			result{[6]float64{279.859214, .000229, .983331}, 1}},
		{gWrapper.Calc,
			swego.CalcFlags{Flags: 1 | flgTopo, TopoLoc: loc},
			result{[6]float64{279.858426, -.000966, .983369}, 32769}},
		{gWrapper.CalcUT,
			swego.CalcFlags{Flags: 1 | flgTopo, TopoLoc: loc},
			result{[6]float64{279.859186, -.000966, .983369}, 32769}},
		{gWrapper.Calc,
			swego.CalcFlags{Flags: 1 | flgSidereal, SidMode: swego.SidMode{Mode: 0}},
			result{[6]float64{255.121938, .000229, .983331}, 65601}},
		{gWrapper.CalcUT,
			swego.CalcFlags{Flags: 1 | flgSidereal, SidMode: swego.SidMode{Mode: 0}},
			result{[6]float64{255.122691, .000229, .983331}, 65601}},
		{gWrapper.Calc,
			swego.CalcFlags{Flags: 1, FileNameJPL: "de406.eph"},
			result{[6]float64{279.858461, .000230, .983331}, 1}},
		{gWrapper.CalcUT,
			swego.CalcFlags{Flags: 1, FileNameJPL: "de406.eph"},
			result{[6]float64{279.859214, .000230, .983331}, 1}},
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			xx, cfl, err := c.fn(2451544.5, 0, c.in)

			if err != nil {
				t.Errorf("err != nil, got: %q", err)
			}

			if !inDelta(xx[:], c.want.xx[:], 1e-6) {
				t.Errorf("xx != %v ± 1e-6, got: %v", c.want.xx, xx)
			}

			if cfl != c.want.cfl {
				t.Errorf("cfl != %d, got: %d", c.want.cfl, cfl)
			}
		})
	}
}

func TestCalc_error(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fn  func(float64, int, swego.CalcFlags) ([6]float64, int, error)
		err string
	}{
		{gWrapper.Calc, "swecgo: jd 99999999.000000 outside JPL eph. range -3027215.50 .. 7930192.50;"},
		{gWrapper.CalcUT, "swecgo: jd 100002561.779707 outside JPL eph. range -3027215.50 .. 7930192.50;"},
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			xx, cfl, err := c.fn(99999999.0, 0, swego.CalcFlags{Flags: 1})

			if err.Error() != c.err {
				t.Fatalf("err != %q, got: %q", err, c.err)
			}

			if xx != ([6]float64{}) {
				t.Error("xx != [6]float{}, got:", xx)
			}

			if cfl != -1 {
				t.Error("xx != -1, got:", cfl)
			}
		})
	}
}

func TestPlanetName(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		name := swe.PlanetName(0)
		if name != "Sun" {
			t.Error("PlanetName(0) != Sun, got:", name)
		}
	})
}

func TestGetAyanamsa(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fn   func(float64, swego.SidMode) float64
		want float64
	}{
		{gWrapper.GetAyanamsa, 24.740393},
		{gWrapper.GetAyanamsaUT, 24.740393},
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			got := c.fn(2451544.5, swego.SidMode{Mode: 0})

			if !inDelta([]float64{got}, []float64{c.want}, 1e-6) {
				t.Errorf("deltaT != %f, got: %f", c.want, got)
			}
		})
	}
}

func TestGetAyanamsaEx(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fn   func(float64, swego.AyanamsaExFlags) (float64, error)
		want float64
	}{
		{gWrapper.GetAyanamsaEx, 24.740393},
		{gWrapper.GetAyanamsaExUT, 24.740393},
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			got, err := c.fn(2451544.5, swego.AyanamsaExFlags{
				Flags: 1,
				SidMode: swego.SidMode{
					Mode:   0,
					T0:     0,
					AyanT0: 0,
				},
			})

			if err != nil {
				t.Fatalf("err != nil, got: %q", err)
			}

			if !inDelta([]float64{got}, []float64{c.want}, 1e-6) {
				t.Errorf("deltaT != %f, got: %f", c.want, got)
			}
		})
	}
}

func TestGetAyanamsaName(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		name := swe.GetAyanamsaName(0)
		if name != "Fagan/Bradley" {
			t.Error("GetAyanamsaName(0) != Fagan/Bradley, got:", name)
		}
	})
}

func TestJulDay(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.JulDay(2000, 1, 1, 0, swego.Gregorian)

		if !inDelta([]float64{got}, []float64{2451544.5}, 1e-6) {
			t.Errorf("JD != 2451544.5, got: %f", got)
		}
	})
}

func TestRevJul(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		y, m, d, h := swe.RevJul(2451544.5, swego.Gregorian)

		if y != 2000 {
			t.Errorf("y != 2000, got: %d", y)
		}

		if m != 1 {
			t.Errorf("m != 1, got: %d", m)
		}

		if d != 1 {
			t.Errorf("d != 1, got: %d", d)
		}

		if h != 0 {
			t.Errorf("h != 0, got: %f", h)
		}
	})
}

func TestUTCToJD(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		et, ut, err := swe.UTCToJD(2000, 1, 1, 0, 0, 0, swego.Gregorian)

		got := []float64{et, ut}
		want := []float64{2451544.500743, 2451544.500004}

		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta(got, want, 1e-6) {
			t.Errorf("[et, ut] != %v, got: %f", want, got)
		}
	})
}

func TestJDETToUTC(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		y, m, d, h, i, s := swe.JdETToUTC(2451544.5, swego.Gregorian)

		if y != 1999 {
			t.Errorf("y != 1999, got: %d", y)
		}

		if m != 12 {
			t.Errorf("m != 12, got: %d", m)
		}

		if d != 31 {
			t.Errorf("d != 31, got: %d", d)
		}

		if h != 23 {
			t.Errorf("h != 23, got: %d", h)
		}

		if i != 58 {
			t.Errorf("i != 58, got: %d", i)
		}

		if !inDelta([]float64{s}, []float64{55.815999}, 1e-6) {
			t.Errorf("s != 55.815999 ± 1e-6, got: %f", s)
		}
	})
}

func TestJDUT1ToUTC(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		y, m, d, h, i, s := swe.JdUT1ToUTC(2451544.5, swego.Gregorian)

		if y != 1999 {
			t.Errorf("y != 1999, got: %d", y)
		}

		if m != 12 {
			t.Errorf("m != 12, got: %d", m)
		}

		if d != 31 {
			t.Errorf("d != 31, got: %d", d)
		}

		if h != 23 {
			t.Errorf("h != 23, got: %d", h)
		}

		if i != 59 {
			t.Errorf("i != 59, got: %d", i)
		}

		if !inDelta([]float64{s}, []float64{59.645586}, 1e-6) {
			t.Errorf("s != 59.645586 ± 1e-6, got: %f", s)
		}
	})
}

func TestHouses(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   int
	}

	type result struct {
		cusps []float64
		ascmc [10]float64
		err   string
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, 'P'}, result{
			[]float64{0,
				190.553653, 215.538288, 246.822987, 283.886819, 319.373115, 348.152982,
				10.553653, 35.538288, 66.822987, 103.886819, 139.373115, 168.152982,
			},
			[10]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
			},
			"",
		}},
		{input{82.083333, 'K'}, result{
			[]float64{0,
				183.972931, 217.277560, 250.582190, 283.886819, 310.582190, 337.277560,
				3.972931, 37.277560, 70.582190, 103.886819, 130.582190, 157.277560,
			},
			[10]float64{
				183.972931, 103.886819, 105.080915, 17.393326,
				196.367263, 352.493044, 195.452718, 172.493044,
			},
			"swecgo: error calculating houses",
		}},
		{input{52.083333, 'G'}, result{
			[]float64{0,
				190.553653, 183.704634, 176.258623, 168.152982, 159.330891, 149.746713,
				139.373115, 128.213369, 116.328153, 103.886819, 91.215325, 78.741976,
				66.822987, 55.634583, 45.214174, 35.538288, 26.565986, 18.252908,
				10.553653, 3.704634, 356.258623, 348.152982, 339.330891, 329.746713,
				319.373115, 308.213369, 296.328153, 283.886819, 271.215325, 258.741976,
				246.822987, 235.634583, 225.214174, 215.538288, 206.565986, 198.252908,
			},
			[10]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
			},
			"",
		}},
	}

	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			cusps, ascmc, err := swe.Houses(2451544.5, c.in.geolat, 5.116667, c.in.hsys)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err != %q, got: %q",
					c.in.geolat, c.in.hsys, c.want.err, err)
			}

			if !inDelta(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps != %v, got: %v",
					c.in.geolat, c.in.hsys, c.want.cusps, cusps)
			}

			if !inDelta(ascmc[:], c.want.ascmc[:], 1e-6) {
				t.Errorf("(%f, %c) ascmc != %v, got: %v",
					c.in.geolat, c.in.hsys, c.want.ascmc, ascmc)
			}
		})
	}
}

func TestHousesEx(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   int
		flags  swego.HousesExFlags
	}

	type result struct {
		cusps []float64
		ascmc [10]float64
		err   string
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, 'P', swego.HousesExFlags{}}, result{
			[]float64{0,
				190.553653, 215.538288, 246.822987, 283.886819, 319.373115, 348.152982,
				10.553653, 35.538288, 66.822987, 103.886819, 139.373115, 168.152982,
			},
			[10]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
			},
			"",
		}},
		{input{82.083333, 'K', swego.HousesExFlags{}}, result{
			[]float64{0,
				183.972931, 217.277560, 250.582190, 283.886819, 310.582190, 337.277560,
				3.972931, 37.277560, 70.582190, 103.886819, 130.582190, 157.277560,
			},
			[10]float64{
				183.972931, 103.886819, 105.080915, 17.393326,
				196.367263, 352.493044, 195.452718, 172.493044,
			},
			"swecgo: error calculating houses",
		}},
		{input{52.083333, 'G', swego.HousesExFlags{}}, result{
			[]float64{0,
				190.553653, 183.704634, 176.258623, 168.152982, 159.330891, 149.746713,
				139.373115, 128.213369, 116.328153, 103.886819, 91.215325, 78.741976,
				66.822987, 55.634583, 45.214174, 35.538288, 26.565986, 18.252908,
				10.553653, 3.704634, 356.258623, 348.152982, 339.330891, 329.746713,
				319.373115, 308.213369, 296.328153, 283.886819, 271.215325, 258.741976,
				246.822987, 235.634583, 225.214174, 215.538288, 206.565986, 198.252908,
			},
			[10]float64{
				190.553653, 103.886819, 105.080915, 24.306488,
				196.367263, 214.734661, 192.275917, 34.734661,
			},
			"",
		}},
		{input{52.083333, 'P', swego.HousesExFlags{
			Flags:   flgSidereal,
			SidMode: swego.SidMode{Mode: 0},
		}},
			result{
				[]float64{0,
					165.817130, 190.801765, 222.086464, 259.150296, 294.636593, 323.416459,
					345.817130, 10.801765, 42.086464, 79.150296, 114.636593, 143.416459,
				},
				[10]float64{
					165.817130, 79.150296, 105.080915, 359.569965,
					171.630740, 189.998138, 167.539394, 9.998138,
				},
				"",
			}},
	}

	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			cusps, ascmc, err := swe.HousesEx(2451544.5, c.in.flags, c.in.geolat,
				5.116667, c.in.hsys)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err != %q, got: %q",
					c.in.geolat, c.in.hsys, c.want.err, err)
			}

			if !inDelta(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps != %v, got: %v",
					c.in.geolat, c.in.hsys, c.want.cusps, cusps)
			}

			if !inDelta(ascmc[:], c.want.ascmc[:], 1e-6) {
				t.Errorf("(%f, %c) ascmc != %v, got: %v",
					c.in.geolat, c.in.hsys, c.want.ascmc, ascmc)
			}
		})
	}
}

func TestHousesArmc(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   int
	}

	type result struct {
		cusps []float64
		ascmc [10]float64
		err   string
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, 'P'}, result{
			[]float64{0,
				190.553489, 215.537888, 246.822499, 283.886657, 319.373244, 348.153088,
				10.553489, 35.537888, 66.822499, 103.886657, 139.373244, 168.153088,
			},
			[10]float64{
				190.553489, 103.886657, 105.080916, 24.307632,
				196.367450, 214.737779, 192.275825, 34.737779,
			},
			"",
		}},
		{input{82.083333, 'K'}, result{
			[]float64{0,
				183.972748, 217.277384, 250.582021, 283.886657, 310.582021, 337.277384,
				3.972748, 37.277384, 70.582021, 103.886657, 130.582021, 157.277384,
			},
			[10]float64{
				183.972748, 103.886657, 105.080916, 17.393607,
				196.367450, 352.493777, 195.452830, 172.493777,
			},
			"swecgo: error calculating houses",
		}},
		{input{52.083333, 'G'}, result{
			[]float64{0,
				190.553489, 183.704585, 176.258665, 168.153088, 159.331033, 149.746863,
				139.373244, 128.213442, 116.328129, 103.886657, 91.215011, 78.741543,
				66.822499, 55.634096, 45.213720, 35.537888, 26.565652, 18.252653,
				10.553489, 3.704585, 356.258665, 348.153088, 339.331033, 329.746863,
				319.373244, 308.213442, 296.328129, 283.886657, 271.215011, 258.741543,
				246.822499, 235.634096, 225.213720, 215.537888, 206.565652, 198.252653,
			},
			[10]float64{
				190.553489, 103.886657, 105.080916, 24.307632,
				196.367450, 214.737779, 192.275825, 34.737779,
			},
			"",
		}},
	}

	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			cusps, ascmc, err := swe.HousesArmc(105.080916, c.in.geolat, 23.439279, c.in.hsys)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err != %q, got: %q",
					c.in.geolat, c.in.hsys, c.want.err, err)
			}

			if !inDelta(cusps, c.want.cusps, 1e-6) {
				t.Errorf("(%f, %c) cusps != %v, got: %v",
					c.in.geolat, c.in.hsys, c.want.cusps, cusps)
			}

			if !inDelta(ascmc[:], c.want.ascmc[:], 1e-6) {
				t.Errorf("(%f, %c) ascmc != %v, got: %v",
					c.in.geolat, c.in.hsys, c.want.ascmc, ascmc)
			}
		})
	}
}

func TestHousePos(t *testing.T) {
	t.Parallel()

	type input struct {
		geolat float64
		hsys   int
	}

	type result struct {
		pos float64
		err string
	}

	cases := []struct {
		in   input
		want result
	}{
		{input{52.083333, 'P'}, result{6.355326, ""}},
		{input{82.083333, 'K'}, result{6.355326, ""}},
		{input{52.083333, 'G'}, result{20.934023, ""}},
	}

	for _, c := range cases {
		Call(nil, func(swe swego.Interface) {
			pos, err := swe.HousePos(105.080916, c.in.geolat, 23.439279, c.in.hsys,
				279.858461, 0.000229)

			if c.want.err != "" && c.want.err != err.Error() {
				t.Fatalf("(%f, %c) err != %q, got: %q",
					c.in.geolat, c.in.hsys, c.want.err, err)
			}

			if !inDelta([]float64{pos}, []float64{c.want.pos}, 1e-6) {
				t.Errorf("(%f, %c) pos != %f, got: %f",
					c.in.geolat, c.in.hsys, c.want.pos, pos)
			}
		})
	}
}

func TestHouseName(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		name := swe.HouseName('P')

		if name != "Placidus" {
			t.Error("HouseName('P') != Placidus, got:", name)
		}
	})
}

func TestDeltaT(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.DeltaT(2451544.5)

		if !inDelta([]float64{got}, []float64{0.000739}, 1e-6) {
			t.Errorf("DeltaT(2451544.5) != 0.000739, got: %f", got)
		}
	})
}

func TestDeltaTEx(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.DeltaTEx(2451544.5, 2)

		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta([]float64{got}, []float64{0.000739}, 1e-6) {
			t.Errorf("DeltaTEx(2451544.5, 2) != 0.000739, got: %f", got)
		}
	})
}

func TestTimeEqu(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.TimeEqu(2451544.5)

		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta([]float64{got}, []float64{-0.002116}, 1e-6) {
			t.Errorf("TimeEqu(2451544.5) != -0.002116, got: %f", got)
		}
	})
}

func TestLMTToLAT(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.LMTToLAT(2451544.5, 5.116667)

		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta([]float64{got}, []float64{2451544.497889}, 1e-6) {
			t.Errorf("LMTToLAT(2451544.5, 5.116667) != 2451544.497889, got: %f", got)
		}
	})
}

func TestLATToLMT(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.LATToLMT(2451544.5, 5.116667)

		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta([]float64{got}, []float64{2451544.502112}, 1e-6) {
			t.Errorf("LATToLMT(2451544.5, 5.116667) != 2451544.502112, got: %f", got)
		}
	})
}

func TestSidTime0(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.SidTime0(2451544.5, 23.439279, -0.003869)

		if !inDelta([]float64{got}, []float64{6.664283}, 1e-6) {
			t.Errorf("SidTime0(2451544.5, 23.439279, -0.003869) != 6.664283, got: %f", got)
		}
	})
}

func TestSidTime(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got := swe.SidTime(2451544.5)

		if !inDelta([]float64{got}, []float64{6.664283}, 1e-6) {
			t.Errorf("SidTime(2451544.5) != 6.664283, got: %f", got)
		}
	})
}
