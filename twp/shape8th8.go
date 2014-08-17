package main

import (
	"fmt"
)

// shaper map, 1/8th lo, 7/8ths hi, shape is u8 - phase is u24

func main() {
	// 0 maps to 0
	// 255 maps to 3/4
	// 3/4 / 255
	s := 3./4 / 255
	knees := make([]float64, 256, 256)
	scaled := make([]uint32, 256, 256)
	for i := 0; i < 256; i++ {
		knees[i] = 1./8 + float64(i) * s
		scaled[i] = uint32(float64(1 << 24) * knees[i])
		// fmt.Printf("%d %f 0x%X\n", i, knees[i], scaled[i])
	}
	// rise rate left is [0, 2^23) -> [0, scaled[i])
	// rise rate right is [2^23, 2^24) -> [scaled[i], 2^24)
	// use knees and scale by 2^24
	sc := make([]float64, 256, 256)
	for i := 0; i < 256; i++ {
		sc[i] = 0.5 / knees[i]
		// fmt.Printf("%d %f\n", i, sc[i])
	}
	// output scaled
	fmt.Printf("\ttype map_8_24 is array(0 to 255)\r\n")
	fmt.Printf("\t\tof std_logic_vector(23 downto 0);\r\n");
	fmt.Printf("\tconstant sincut: map_8_24 :=\r\n");
	fmt.Printf("\t(\r\n");
	for i, s := 0, "\t\t"; i < 256; i++ {
		fmt.Printf("%s\"%024b\"", s, scaled[i])
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n");
	// output scale
	fmt.Printf("\r\n");
	fmt.Printf("\tconstant sinscale: map_8_24 :=\r\n");
	fmt.Printf("\t(\r\n");
	for i, s := 0, "\t\t"; i < 256; i++ {
		fmt.Printf("%s\"%024b\"", s, bits24(sc[i]))
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n");
}

// input is 0 to 4.00, we make 4 a little less to save a bit
func bits24(v float64) uint32 {
	b := uint32(v * float64(1 << 22))
	if b >= 1 << 24 {
		b = (1 << 24) - 1
	}
	return b
}
