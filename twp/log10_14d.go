package main

import (
	"math"
	"fmt"
)

const (
	a = -23358.8
	b = 2363.73
	c = 19583.8

	fold = true
)

func main() {
	var r float64
	fmt.Printf("\ttype map_10_14 is array(0 to 1023)\r\n")
	fmt.Printf("\t\tof std_logic_vector(13 downto 0);\r\n")
	fmt.Printf("\tconstant log10_14: map_10_14 :=\r\n")
	fmt.Printf("\t(\r\n")
	z := make([]uint, 1025, 1025)
	for i, s := 0, "\t\t"; i < 1024; i++ {
		x := i
		if fold {
			x = (1 << 10) - 1 - x
		}
		if x == 0 {
			r = 0.0
		} else {
			r = a + b * math.Log(c * float64(x))
		}
		v := uint(r + 0.5)
		if fold {
			v = (1 << 14) - 1 - v
		}
		z[i] = v
		fmt.Printf("%s\"%014b\"", s, v)
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	z[1024] = 1 << 14
	fmt.Printf("\r\n\t);\r\n")
	d := make([]uint, 1024, 1024)
	max := uint(0)
	for i := 0; i < 1024; i++ {
		v := z[i+1] -z[i]
		d[i] = v
		if v > max {
			max = v
		}
	}
	// fmt.Printf("max %d\n", max) // -> 1638 -> 11 bits
	fmt.Printf("\r\n")
	fmt.Printf("\ttype map_10_11 is array(0 to 1023)\r\n")
	fmt.Printf("\t\tof std_logic_vector(10 downto 0);\r\n")
	fmt.Printf("\tconstant diff10_11: map_10_11 :=\r\n")
	fmt.Printf("\t(\r\n")
	for i, s := 0, "\t\t"; i < 1024; i++ {
		fmt.Printf("%s\"%011b\"", s, d[i])
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n")
}
