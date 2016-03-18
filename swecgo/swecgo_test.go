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
				t.Errorf("cfl != %d, got: %d", cfl, c.want.cfl)
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

			if xx != ([6]float64{}) {
				t.Error("xx != [6]float{}, got:", xx)
			}

			if cfl != -1 {
				t.Error("xx != -1, got:", cfl)
			}

			if err.Error() != c.err {
				t.Errorf("err != %q, got: %q", err, c.err)
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
