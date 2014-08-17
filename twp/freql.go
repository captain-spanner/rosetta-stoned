package main

import (
	"fmt"
	"strconv"
)

const (
	f0 = "1"
	f1 = "100000100011010100"
//	scale is 200 / 3, 2 guard
//	sr = 44100.
//	sr = 48000.
//	sr = 88200.
	sr = 96000.
	fm = 200. * (1 << 2)
	fd = 3.
)

func main() {
	d0, _ := strconv.ParseUint(f0, 2, 32)
	d1, _ := strconv.ParseUint(f1, 2, 32)
	sc := sr * fm / fd
	fmt.Printf("d1/sr %9f\n", float64(d1) / sr)
	v0 := float64(d0) / sc
	v1 := float64(d1) / sc
	fmt.Printf("lo %9f hi %9f\n", v0, v1)
	fmt.Printf("freq lo %9f hi %9f\n", v0 * sr, v1 * sr)
	l, rm := leading(sc)
	fmt.Printf("1/sc leading %d, rm 0x%X\n", l, rm)
	s1 := d1 * rm
	fmt.Printf("unscaled 0x%X\n", s1)
	fmt.Printf("hi %d rm X\"%X\"\n", s1 >> 19, rm)
}

func leading(sc float64) (int, uint64) {
	v := uint64(float64(1 << 43) * (1.0 / sc))
	c := 0
	for i := 43; i > 0; i-- {
		if v & (1 << uint(i)) == 0 {
			c++
		} else {
			break
		}
	}
	return c, v
}
