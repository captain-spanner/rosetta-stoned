package main

import (
	"fmt"
)

// Test the order 3 find

func main() {
	// 2^qN is pi/2
	for i := 0; i < 8; i++ {
		r := isin_S3(int32(i << (qN -1)))
		fmt.Printf("%d %x %f\n", 45 * i, r, float32(r) / float32(1 << qA))
	}
}

// S(x) = x * ( (3<<p) - (x*x>>r) ) >> s
// n : Q-pos for quarter circle             12
// A : Q-pos for output                     13
// p : Q-pos for parentheses intermediate   15
// r = 2n-p                                 11
// s = A-1-p-n                              17

const (
	qN = 12
	qA = 13
	qP = 15
//	qN = 22
//	qA = 24
//	qP = 20
	qR = 2*qN - qP
	qS = qN + qP + 1 - qA
)

func isin_S3(x int32) int32 {
	x = x << (30 - qN) // shift to full s32 range (Q13->Q30)

	if (x ^ (x << 1)) < 0 { // test for quadrant 1 or 2
		x = int32((1 << 31) - uint32(x))
	}

	x = x >> (30 - qN)

	return x * ((3 << qP) - (x * x >> qR)) >> qS
}

const (
	Z = 32
	Z1 = Z - 1
	Z2 = Z - 2
)

func isin_S3x(a int32) int32 {
	x := int64(a)
	x = x << (Z2 - qN) // shift to full Z range

	if (x ^ (x << 1)) < 0 { // test for quadrant 1 or 2
		x = (1 << Z1) - x
	}

	x = x >> (Z2 - qN)
	t0 := x * x
	t1 := (3 << qP) - (t0 >> qR)
	r := x * t1 >> qS
	return int32(r)
}
