package shapefile

import (
	"fmt"
	"math/rand"
)

const (
	epsz = 8
)

type axdesc struct {
	v0 float64
	vz float64
}

func (s *Shapefile) Stoch(n int, seed int64) {
	sz := len(s.regs)
	xd, yd, r, e := s.mkdata("", sz, n, seed)
	for i := 0; i < n; i++ {
		x := xd.choose(r)
		y := yd.choose(r)
		c := s.Where(&Point{x: x, y: y})
		if c < 0 {
			e[sz]++
		} else {
			e[c]++
		}
	}
	for i, v := range e {
		if v == 0 {
			continue
		}
		if i == sz {
			fmt.Printf("Pirate land: %d (%d%%)\n", v, v*100/n)
		} else {
			fmt.Printf("Reg %d: %d\n", i, v)
		}
	}
}

func (s *Shapefile) StochDebug(n int, seed int64) {
	sz := len(s.regs)
	xd, yd, r, e := s.mkdata("", sz, n, seed)
	for i := 0; i < n; i++ {
		x := xd.choose(r)
		y := yd.choose(r)
		c := s.Where(&Point{x: x, y: y})
		if c < 0 {
			e[sz]++
		} else {
			e[c]++
		}
	}
	for i, v := range e {
		if v == 0 {
			continue
		}
		if i == sz {
			fmt.Printf("Pirate land: %d (%d%%)\n", v, v*100/n)
		} else {
			fmt.Printf("Reg %d: %d\n", i, v)
		}
	}
	fmt.Println("=========")
}

func (s *Shapefile) StochEps(n int, seed int64) {
	sz := epsz
	xd, yd, r, e := s.mkdata("Eps", sz, n, seed)
	for i := 0; i < n; i++ {
		x := xd.choose(r)
		y := yd.choose(r)
		c := s.Where(&Point{x: x, y: y})
		if c < sz {
			e[c]++
		} else {
			e[sz]++
		}
	}
	for i, v := range e {
		if v == 0 {
			continue
		}
		switch i {
		case 0:
			fmt.Printf("Pirate land: %d (%d%%)\n", v, v*100/n)
		case 1:
			fmt.Printf("Claimed: %d\n", v)
		case sz:
			fmt.Printf("Nutty: %d\n", v)
		default:
			fmt.Printf("Multi %d: %d\n", i, v)
		}
	}
}

func (s *Shapefile) mkdata(t string, z int, n int, seed int64) (*axdesc, *axdesc, *rand.Rand, []int) {
	fmt.Printf("Stoch%s: count %d, seed %d\n", t, n, seed)
	xd := mkAxdesc(s.box.xmin, s.box.xmax)
	yd := mkAxdesc(s.box.ymin, s.box.ymax)
	r := rand.New(rand.NewSource(seed))
	e := make([]int, z+1, z+1)
	return xd, yd, r, e
}

func mkAxdesc(min, max float64) *axdesc {
	return &axdesc{v0: min, vz: max - min}
}

func (a *axdesc) choose(r *rand.Rand) float64 {
	return a.v0 + r.Float64()*a.vz
}
