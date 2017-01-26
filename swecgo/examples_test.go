package swecgo

import (
	"fmt"

	"github.com/astrotools/swego"
)

func ExampleCall_calcUT() {
	Call(func(swe swego.Interface) {
		xx, cfl, err := swe.CalcUT(2451544.5, swego.Sun, &swego.CalcFlags{})
		if err != nil {
			fmt.Println("Calculation error: ", err)
			return
		}

		fmt.Printf("xx[0] %f\n", xx[0])
		fmt.Printf("xx[1] %f\n", xx[1])
		fmt.Printf("xx[2] %f\n", xx[2])
		fmt.Printf("xx[3] %f\n", xx[3])
		fmt.Printf("xx[4] %f\n", xx[4])
		fmt.Printf("xx[5] %f\n", xx[5])
		fmt.Println("cfl", cfl)

		// Output:
		// xx[0] 279.859214
		// xx[1] 0.000230
		// xx[2] 0.983332
		// xx[3] 0.000000
		// xx[4] 0.000000
		// xx[5] 0.000000
		// cfl 2
	})
}
