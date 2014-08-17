package main

import (
	"math"
	"fmt"
)

const (
	a = -36266.9
	b = 3434.95
	c = 38498.9

	fold = true
)

func fit(x int) float64 {
	return a + b * math.Log(c * float64(x + 1))
}

func test() {
	v0 := fit(0)
	v1 := fit(16383)
	i0 := uint(v0 + 0.5)
	i1 := uint(v1 + 0.5)
	fmt.Printf("-- 0 -> %d, 16383 -> %d\n", i0, i1)
	fmt.Printf("--\t scale is 200 / 3, 2 guard\n")
	fmt.Printf("--\t range is %f -> %f\n", v0*3./200., v1*3./200.)
}

func main() {
	var r float64
	test()
	fmt.Printf("\t-- type freq16 is array(0 to 16383)\r\n")
	fmt.Printf("\t-- \tof std_logic_vector(17 downto 0);\r\n");
	fmt.Printf("\tconstant freq16_l: freq16 :=\r\n");
	fmt.Printf("\t(\r\n");
	for i, s := 0, "\t\t"; i < 16384; i++ {
		x := i
		if fold {
			x = 16383 - x
		}
		r = a + b * math.Log(c * float64(x + 1))
		r = r * 4.
		v := int(r + 0.5)
		if fold {
			v = 133334 - v
		}
		if v <= 0 {
			v = 1
		}
		fmt.Printf("%s\"%018b\"", s, v)
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n");
}
