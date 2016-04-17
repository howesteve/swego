package swego

import "testing"

func TestCalcFlags_SetEphemeris(t *testing.T) {
	fl := CalcFlags{}
	fl.SetEphemeris(JPL)

	if fl.Flags != int32(JPL) {
		t.Error("flags value does not contain ephemeris flag")
	}
}
