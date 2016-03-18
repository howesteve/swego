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
		fn   func(float64) float64
		want float64
	}{
		{gWrapper.GetAyanamsa, 24.740393},
		{gWrapper.GetAyanamsaUT, 24.740393},
	}

	for _, c := range cases {
		Call(nil, func(_ swego.Interface) {
			got := c.fn(2451544.5)

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

			if !inDelta([]float64{c.want}, []float64{got}, 1e-6) {
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

		if !inDelta([]float64{2451544.5}, []float64{got}, 1e-6) {
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

		if !inDelta(want, got, 1e-6) {
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

		if !inDelta([]float64{55.815999}, []float64{s}, 1e-6) {
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

		if !inDelta([]float64{59.645586}, []float64{s}, 1e-6) {
			t.Errorf("s != 59.645586 ± 1e-6, got: %f", s)
		}
	})
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

func TestTimeEqu(t *testing.T) {
	t.Parallel()

	Call(nil, func(swe swego.Interface) {
		got, err := swe.TimeEqu(2451544.5)

		if err != nil {
			t.Fatalf("err != nil, got: %q", err)
		}

		if !inDelta([]float64{-0.002116}, []float64{got}, 1e-6) {
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
