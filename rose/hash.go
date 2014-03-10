package rose

import (
)

const (
	hashp	= 16777619
	nprimes	= 8
	hrot	= 5
)

var (
	primes	[nprimes]uint32 = [nprimes]uint32 {
		73,
		167,
		379,
		661,
		1129,
		2039,
		3797,
		6563,
	}
)

// from disk/hash.c
func hashs(s string) uint32 {
	var h uint32 = 0
	v := []byte(s)
	for _, b := range v {
		h *= hashp
		h ^= uint32(b)
	}
	return h
}

// from disk/map.c
func hashrot(h uint32) uint32 {
	return (h << hrot) | (h >> (32 - hrot))
}

func issetmap(m []byte, x uint32, i int) bool {
	b := m[x]
	t := b & (1 << uint(i))
	set := t != 0
	return set
}

func checkmap(m []byte, h uint32) bool {
	z := uint32(len(m))
	for i, p := range primes {
		h *= p
		x := h % z
		if !issetmap(m, x, i) {
			return false
		}
		h = hashrot(h)
	}
	return true
}
