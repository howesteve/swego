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
		{'D', true}, // Equal houses, where cusp 10 = MC
		{'E', true}, // Equal
		{'F', true}, // Carter poli-equatorial
		{'G', true}, // Gauquelin sectors
		{'H', true}, // Azimuthal
		{'I', true}, // Sunshine (Treindl)
		{'i', true}, // Sunshine (Makransky)
		{'K', true}, // Koch
		{'L', true}, // Pullen SD (sinusoidal delta) = ex Neo-Porphyry
		{'M', true}, // Morinus
		{'N', true}, // Equal houses, where cusp 1 = 0 Aries
		{'O', true}, // Porphyrius
		{'P', true}, // Placidus
		{'Q', true}, // Pullen SR (sinusoidal ratio)
		{'R', true}, // Regiomontanus
		{'S', true}, // Sripati
		{'T', true}, // Polich-Page
		{'U', true}, // Krusinski-Pisa-Goelzer
		{'V', true}, // Vehlow equal
		{'W', true}, // Whole sign
		{'X', true}, // Axial rotation
		{'Y', true}, // APC houses
		{'_', false},
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

func TestNewHSys_lowerToUpper(t *testing.T) {
	got, ok := NewHSys('a')
	if !ok {
		t.Error("ok = false, want: true")
	}

	want := HSys('A')
	if got != want {
		t.Errorf("hsys = %c, want: %c", got, want)
	}
}

func TestCalcFlags_Copy(t *testing.T) {
	fl := new(CalcFlags)
	fl.Flags = FlagSpeed
	got := fl.Copy()

	if got == fl {
		t.Errorf("%p = %p, want copy", got, fl)
	}
}

func TestCalcFlags_SetEphemeris(t *testing.T) {
	fl := new(CalcFlags)
	fl.SetEphemeris(JPL)

	if fl.Flags != int32(JPL) {
		t.Error("flags value does not contain ephemeris flag")
	}
}

type testInterface struct{ Interface }
type testExclLocker struct{ Interface }
type testLockedIface struct{ Interface }

func (l *testLockedIface) ExclusiveUnlock() {}
func (el *testExclLocker) ExclusiveLock() LockedInterface {
	return &testLockedIface{el}
}

func TestLocked(t *testing.T) {
	t.Run("Interface", func(t *testing.T) {
		called := make(chan struct{}, 1)
		swe := new(testInterface)
		Locked(swe, func(i Interface) {
			if i != swe {
				t.Error("passed Interface no equal to input Interface")
			}

			called <- struct{}{}
		})

		select {
		case <-called:
		default:
			t.Error("callback not called")
		}
	})

	t.Run("ExclusiveLocker", func(t *testing.T) {
		called := make(chan struct{}, 1)
		swe := (*testExclLocker)(nil)
		Locked(swe, func(i Interface) {
			if i.(*testLockedIface).Interface != swe {
				t.Error("passed Interface no equal to input Interface")
			}

			called <- struct{}{}
		})

		select {
		case <-called:
		default:
			t.Error("callback not called")
		}
	})
}
