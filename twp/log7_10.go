package main

import (
	"math"
	"fmt"
)

const (
	a = -1489.77
	b = 210.972
	c = 1171.57

	fold = true
)

func main() {
	var r float64
	fmt.Printf("\ttype map_7_10 is array(0 to 127)\r\n")
	fmt.Printf("\t\tof std_logic_vector(9 downto 0);\r\n");
	fmt.Printf("\tconstant log7_10: map_7_10 :=\r\n");
	fmt.Printf("\t(\r\n");
	for i, s := 0, "\t\t"; i < 128; i++ {
		x := i
		if fold {
			x = 127 - x
		}
		if x == 0 {
			r = 0.0
		} else {
			r = a + b * math.Log(c * float64(x))
		}
		v := uint(r + 0.5)
		if fold {
			v = 1023 - v
		}
		fmt.Printf("%s\"%010b\"", s, v)
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n");
}
