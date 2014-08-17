package main

import (
	"math"
	"fmt"
)

const (
	a = -6141.51
	b = 738.81
	c = 4080.4

	fold = true
)

func main() {
	var r float64
	fmt.Printf("\ttype map_8_12 is array(0 to 255)\r\n")
	fmt.Printf("\t\tof std_logic_vector(11 downto 0);\r\n");
	fmt.Printf("\tconstant log8_12: map_8_12 :=\r\n");
	fmt.Printf("\t(\r\n");
	for i, s := 0, "\t\t"; i < 256; i++ {
		x := i
		if fold {
			x = 255 - x
		}
		if x == 0 {
			r = 0.0
		} else {
			r = a + b * math.Log(c * float64(x))
		}
		v := uint(r + 0.5)
		if fold {
			v = 4095 - v
		}
		fmt.Printf("%s\"%012b\"", s, v)
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n");
}
