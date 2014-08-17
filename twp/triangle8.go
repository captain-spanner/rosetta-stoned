package main

import (
	"fmt"
)

// traingle wave, shape is u8 - phase is u24

func main() {
	// 0 maps to 0
	// 255 maps to 1/2
	s := .5 / 255
	knees := make([]float64, 256, 256)
	scaled := make([]uint32, 256, 256)
	for i := 0; i < 256; i++ {
		knees[i] = float64(i) * s
		scaled[i] = uint32(float64(1 << 24) * knees[i])
		// fmt.Printf("%d %f 0x%X\n", i, knees[i], scaled[i])
	}
	sc := make([]float64, 256, 256)
	for i := 0; i < 256; i++ {
		sc[i] = .5 / knees[i]
		// fmt.Printf("%d %f\n", i, sc[i])
	}
	// output cut
	fmt.Printf("\tconstant tricut: map_8_24 :=\r\n");
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
	// output sc
	fmt.Printf("\r\n");
	fmt.Printf("\tconstant triscale: map_8_24 :=\r\n");
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


// input is 0 to 255.00
func bits24(v float64) uint32 {
	b := uint32(v * float64(1 << 16))
	return b
}
