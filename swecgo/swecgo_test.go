// +build linux,cgo darwin,cgo

package swecgo

import (
	"testing"

	"github.com/astrotools/swego"
)

func TestConstantCheck(t *testing.T) {
	if swego.FlagSidereal != flgSidereal {
		t.Errorf("Sidereal = %d, want: %d", swego.FlagSidereal, flgSidereal)
	}

	if swego.FlagTopo != flgTopo {
		t.Errorf("Topo = %d, want: %d", swego.FlagTopo, flgTopo)
	}
}

func TestInvoker(t *testing.T) {
	inv := NewInvoker(nil)

	invoked := false
	inv.Invoke(func(swe swego.Interface) {
		invoked = true
	})

	if !invoked {
		t.Error("closure is not invoked")
	}
}
