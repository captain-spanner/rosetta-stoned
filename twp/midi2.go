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
	HamG   = 20
	HamS   = 720720.
	pref   = "midi"
	debugb = false
	debugm = false
)

var (
	midif  [RANGE]float64
	midifh [RANGE]float64
	tc     int  = 'a' - 1
	deft   int
)

type table struct {
	tc    int
	sr    float64
	ham   bool
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
	ts := make([]*table, 0)
	ts = append(ts, mktable(FreqA, false))
	ts = append(ts, mktable(FreqB, false))
	ts = append(ts, mktable(FreqA, true))
	ts = append(ts, mktable(FreqB, true))
	for _, t := range ts {
		t.process()
	}
}

func (t *table) process() {
	t.deltas()
	t.scales()
	t.trail()
	fmt.Printf("-- %c: bits = %d, tail = %d, ham = %t\r\n", t.tc, t.bits, t.tail, t.ham)
	t.split()
	t.deftype()
/*
	t.mkvhh()
*/
}

// 2^((m-69)/12)*(440 Hz)
func mkmidi() {
	for i := 0; i < RANGE; i++ {
		v := math.Pow(2., float64(i-M440)/12.0) * 440.
		midif[i] = v
		midifh[i] = v * float64(1 << HamG) / HamS
	}
	if debugm {
		fmt.Println(midif)
	}
}

func mktable(sr float64, h bool) *table {
	tc++
	return &table{tc: tc, sr: sr, ham: h}
}

func (t *table) deltas() {
	for i := 0; i < RANGE; i++ {
		if t.ham {
			t.delta[i] = midifh[i] / t.sr
		} else {
			t.delta[i] = midif[i] / t.sr
		}
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
	if t.ham {
		return
	}
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

func (t *table) deftype() {
	if t.ham {
		deftype(t.bits)
	} else {
		deftype(t.head)
		deftype(t.tail)
		deftype(t.tail-1)
	}
}

func deftype(z int) {
	b := 1 << uint(z)
	if b & deft != 0 {
		return
	}
	fmt.Printf("\ttype %s%d is array(0 to %d)\r\n", pref, z, RANGE-1)
	fmt.Printf("\t\tof std_logic_vector(%d downto 0);\r\n", z-1)
	deft |= b
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
