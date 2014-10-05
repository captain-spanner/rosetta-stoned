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

func (s *Shapefile) StochEps(n int, seed int64) {
	fmt.Printf("StochEps: count %d, seed %d\n", n, seed)
	xd := mkAxdesc(s.box.xmin, s.box.xmax)
	yd := mkAxdesc(s.box.ymin, s.box.ymax)
	r := rand.New(rand.NewSource(seed))
	e := make([]int, epsz+1, epsz+1)
	for i := 0; i < n; i++ {
		x := xd.choose(r)
		y := yd.choose(r)
		c := s.quad.SearchEps(&point{x: x, y: y})
		if c < epsz {
			e[c]++
		} else {
			e[epsz]++
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
		case epsz:
			fmt.Printf("Nutty: %d\n", v)
		default:
			fmt.Printf("Multi %d: %d\n", i, v)
		}
	}
}

func mkAxdesc(min, max float64) *axdesc {
	return &axdesc{v0: min, vz: max - min}
}

func (a *axdesc) choose(r *rand.Rand) float64 {
	return a.v0 + r.Float64()*a.vz
}
