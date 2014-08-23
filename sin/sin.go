package main

import (
	"fmt"
	"math"
)

// Test out the Cody and Waite algorithm on 24 bits

func main() {
	const test = 1
	switch test {
	case 0:
		var nv, xv uint32
		var min, max float64
		min, max = 0.0, 0.0
		for i := uint32(0); i < 0x400000; i++ {
			f := sin(i)
			if f < min {
				min = f
				nv = i
			}
			if f > max {
				max = f
				xv = i
			}
		}
		fmt.Printf("min %f at %x\n", min, nv)
		fmt.Printf("max %f at %x\n", max, xv)
	case 1:
		fmt.Printf("pi / 2 = %f\n", math.Pi / 2.0)
		sin_of(44)
		sin_of(45)
		sin_of(46)
		sin_of(89)
		sin_of(90)
		sin_of(91)
	}
}

func sin_of(d int) {
	v := sin(deg(d))
	fmt.Printf("%d - %f\n", d, v)
}

// degrees
func deg(d int) uint32 {
	return uint32((uint64(1) << 24) * uint64(d) / 360)
}

// input: 24 bits, 0 to 2 Pi
// output: 24 bits, nominally signed +-5v

func sin(x uint32) float64 {
	const b22 = 0x3FFFFF
	const r1 = -0.6666662674
	const r2 =  0.1333284022
	const r3 = -0.0126767480
	const r4 =  0.0006660872
	const eps = -1.0 / float64(1 << 12)
	var n, s bool
	var v float64
	q := x >> 22
	b := x & b22
	switch q {
	case 0:
		n = false
		s = false
	case 1:
		n = true
		s = true
	case 2:
		n = true
		s = false
	case 3:
		n = false
		s = true
	}
	if s {
		b = b22 - b
	}
	fmt.Printf("n = %t, s = %t\n", n, s)
	f := (math.Pi / 2.0) * (float64(b) / float64(1 << 22))
	if f < eps {
		if s {
			v = -f
		} else {
			v = f
		}
	} else {
		if s {
			f = -f
		}
		fmt.Printf("f = %f\n", f)
		g := f * f
		// unrolled for FPGA
		t0 := r4 * g + r3
		t1 := t0 * g + r2
		t2 := t1 * g + r1
		r := t2 * g
		v = f + f * r
	}
	if n {
		v = -v
	}
	return v
}
