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
	bits  int
}

func main() {
	mkmidi()
	a := mktable(FreqA)
	b := mktable(FreqB)
	a.deltas()
	b.deltas()
	a.scales()
	b.scales()
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
	b := 1 << uint(t.bits)
	if b&m != 0 {
		return m
	}
	fmt.Printf("\ttype %s%d is array(0 to %d)\r\n", pref, t.bits, RANGE-1)
	fmt.Printf("\t\tof std_logic_vector(%d downto 0);\r\n", t.bits-1)
	m |= b
	return m
}

func (t *table) binary(v uint32) string {
	s := fmt.Sprintf("%%0%db", t.bits)
	return fmt.Sprintf(s, v)
}

func (t *table) mkvhh() {
	fmt.Printf("\r\n")
	fmt.Printf("\tconstant %s%d: %s%d :=\r\n", pref, int(t.sr), pref, t.bits)
	fmt.Printf("\t(\r\n")
	for i, s := 0, "\t\t"; i < RANGE; i++ {
		fmt.Printf("%s\"%s\"", s, t.binary(t.scale[i]))
		if (i & 3) == 3 {
			s = ",\r\n\t\t"
		} else {
			s = ", "
		}
	}
	fmt.Printf("\r\n\t);\r\n")
}
