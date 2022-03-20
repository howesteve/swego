//go:build (linux && cgo) || (darwin && cgo)
// +build linux,cgo darwin,cgo

package swecgo

import (
	"fmt"

	"github.com/howesteve/swego"
)

func Example() {
	swe := Open()

	xx, cfl, err := swe.CalcUT(2451544.5, swego.Sun, &swego.CalcFlags{})
	if err != nil {
		fmt.Println("calculation error: ", err)
		return
	}

	fmt.Printf("xx[0] %f\n", xx[0])
	fmt.Printf("xx[1] %f\n", xx[1])
	fmt.Printf("xx[2] %f\n", xx[2])
	fmt.Printf("xx[3] %f\n", xx[3])
	fmt.Printf("xx[4] %f\n", xx[4])
	fmt.Printf("xx[5] %f\n", xx[5])
	fmt.Printf("cfl %d\n", cfl)

	// Output:
	// xx[0] 279.859214
	// xx[1] 0.000230
	// xx[2] 0.983332
	// xx[3] 0.000000
	// xx[4] 0.000000
	// xx[5] 0.000000
	// cfl 2
}

func ExampleLocked() {
	swe := Open()

	// Locked guarantees that both calls to swe.CalcUT will be executed in order
	// as a single unit without the interference of other calls to the library.
	Locked(swe, func(swe Library) {
		// Expensive call because the result is computed.
		fl := new(swego.CalcFlags)
		swe.CalcUT(2451544.5, swego.Sun, fl)

		// Cheap call because results are cached.
		fl.Flags |= swego.FlagEquatorial
		swe.CalcUT(2451544.5, swego.Sun, fl)
	})
}
