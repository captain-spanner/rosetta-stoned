package main

import (
	"fmt"
	"math"
)

const (
	RANGE  = 128
	M440   = 69
	FreqA  = 44100.
	FreqB  = 48000.
	pref   = "midi"
	debugb = false
	debugm = false
)

var (
	midif [RANGE]float64
)

type table struct {
	sr    float64
	delta [RANGE]float64
	scale [RANGE]uint32
	top   [RANGE]uint32
	low   [RANGE]uint32
	low2  [RANGE]uint32
	bits  int
	head  int
	tail  int
}

func main() {
	mkmidi()
	a := mktable(FreqA)
	b := mktable(FreqB)
	a.deltas()
	b.deltas()
	a.scales()
	b.scales()
	a.trail()
	b.trail()
	fmt.Printf("-- a: bits = %d, tail = %d\r\n", a.bits, a.tail)
	fmt.Printf("-- b: bits = %d, tail = %d\r\n", b.bits, b.tail)
	a.split()
	b.split()
	m := a.deftype(0)
	b.deftype(m)
	a.mkvhh()
	b.mkvhh()
}

// 2^((m-69)/12)*(440 Hz)
func mkmidi() {
	for i := 0; i < RANGE; i++ {
		midif[i] = math.Pow(2., float64(i-M440)/12.0) * 440.
	}
	if debugm {
		fmt.Println(midif)
	}
}

func mktable(sr float64) *table {
	return &table{sr: sr}
}

func (t *table) deltas() {
	for i := 0; i < RANGE; i++ {
		t.delta[i] = midif[i] / t.sr
	}
}

func (t *table) scales() {
	for i := 0; i < RANGE; i++ {
		s := t.delta[i] * float64(1<<24)
		v := uint32(s + 0.5)
		t.scale[i] = v
	}
	t.bits = bits(t.scale[RANGE-1])
}

func (t *table) trail() {
	var max uint = 0
	for i := 0; i < RANGE; i++ {
		u := t.scale[i]
		v := u + 1
		d := u ^ v
		var j uint = 24
		for ; j >= 0; j-- {
			if d & (1 << j) != 0 {
				break
			}
		}
		if j > max {
			max = j
		}
	}
	t.tail = int(max) + 1
}

func (t *table) split() {
	t.head = t.bits - t.tail
	s := uint32(t.tail)
	m := uint32(1 << s) - 1
	for i, v := range t.scale {
		top := v >> s
		low := v & m
		low2 := (low + 1) >> 1
		t.top[i] = top
		t.low[i] = low
		t.low2[i] = low2
	}
}

func bits(v uint32) int {
	i := 0
	for l := uint32(1); v > l; i, l = i+1, l<<1 {
	}
	if debugb {
		fmt.Printf("max = %d, bits = %d\n", v, i)
	}
	return i
}

func (t *table) deftype(m int) int {
	m = deftype(t.head, m)
	m = deftype(t.tail, m)
	m = deftype(t.tail-1, m)
	return m
}

func deftype(z, m int) int {
	b := 1 << uint(z)
	if b&m != 0 {
		return m
	}
	fmt.Printf("\ttype %s%d is array(0 to %d)\r\n", pref, z, RANGE-1)
	fmt.Printf("\t\tof std_logic_vector(%d downto 0);\r\n", z-1)
	m |= b
	return m
}

func binary(v uint32, b int) string {
	s := fmt.Sprintf("%%0%db", b)
	return fmt.Sprintf(s, v)
}

func (t *table) mkvhh() {
	mkvhh(t.top[:], t.head, t.sr, "head")
	mkvhh(t.low[:], t.tail, t.sr, "low0")
	mkvhh(t.low2[:], t.tail-1, t.sr, "low1")
}

func mkvhh(tab []uint32, b int, sr float64, tag string) {
	m := 3
	if b < 8 {
		m = 7
	}
	fmt.Printf("\r\n")
	fmt.Printf("\tconstant %s%d_%s: %s%d :=\r\n", pref, uint(sr), tag, pref, b)
	fmt.Printf("\t(\r\n")
	for i, s := 0, "\t\t"; i < RANGE; i++ {
		fmt.Printf("%s\"%s\"", s, binary(tab[i], b))
		if (i & m) == m {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n")
}
