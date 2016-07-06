package swecgo

import (
	"fmt"

	"github.com/dwlnetnl/swego"
)

func ExampleCall_calcUT() {
	Call(nil, func(swe swego.Interface) {
		xx, cfl, err := swe.CalcUT(2451544.5, swego.Sun, &swego.CalcFlags{})
		if err != nil {
			fmt.Println("Calculation error: ", err)
			return
		}

		fmt.Println("xx[0]", xx[0])
		fmt.Println("xx[1]", xx[1])
		fmt.Println("xx[2]", xx[2])
		fmt.Println("xx[3]", xx[3])
		fmt.Println("xx[4]", xx[4])
		fmt.Println("xx[5]", xx[5])
		fmt.Println("cfl", cfl)

		// Output:
		// xx[0] 279.8592144230897
		// xx[1] 0.0002296532779708713
		// xx[2] 0.9833318568951199
		// xx[3] 0
		// xx[4] 0
		// xx[5] 0
		// cfl 2
	})
}
