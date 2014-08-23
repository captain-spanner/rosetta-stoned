package main

import (
	"fmt"
)

func main() {
	fmt.Printf("%f\n", float64(sinz((1 << 22) - 1)) / float64(1 << 23))
	fmt.Printf("%f\n", float64(sinz(1 << 21)) / float64(1 << 23))
}

// first quadrant
// z * (3 - z*z*) / 2
// 22 bit signed in [0, 1), 23 bit out [0, 1)
func sinz(a uint32) int32 {
	// z:22 [0, 1)
	z := uint64(a)
	// z2:44 [0, 1)
	z2 := z * z
	// z3:22 [0, 1)
	z3 := z2 >> 22
	// z4:24 [0, 3)
	z4 := (3 << 22) - z3
	// z5:46 [0, 3)
	z5 := z * z4
	return int32(z5 >> 22)
}
