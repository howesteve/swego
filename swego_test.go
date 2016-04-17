package swego

import "testing"

func TestNewHSys(t *testing.T) {
	cases := []struct {
		in   byte
		want bool
	}{
		{'A', true}, // Equal
		{'B', true}, // Alcabitius
		{'C', true}, // Campanus
		{'E', true}, // Equal
		{'G', true}, // Gauquelin sectors
		{'H', true}, // Azimuthal
		{'K', true}, // Koch
		{'M', true}, // Morinus
		{'O', true}, // Porphyrius
		{'P', true}, // Placidus
		{'R', true}, // Regiomontanus
		{'T', true}, // Polich-Page
		{'U', true}, // Krusinski-Pisa-Goelzer
		{'V', true}, // Vehlow equal
		{'W', true}, // Whole sign
		{'X', true}, // Axial rotation
		{'Y', true}, // APC houses
		{'_', false},
		{'S', false},
		{'Z', false},
	}

	for _, c := range cases {
		hsys, ok := NewHSys(c.in)

		if ok && byte(hsys) != c.in {
			t.Errorf("input %q is not returned", c.in)
		}

		if !ok && byte(hsys) != '\000' {
			t.Errorf("invalid input does not return '\\000' as hsys")
		}

		if ok != c.want {
			t.Errorf("%q is no valid house system", c.in)
		}
	}
}

func TestCalcFlags_SetEphemeris(t *testing.T) {
	fl := CalcFlags{}
	fl.SetEphemeris(JPL)

	if fl.Flags != int32(JPL) {
		t.Error("flags value does not contain ephemeris flag")
	}
}
