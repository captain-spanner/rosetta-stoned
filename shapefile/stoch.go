package shapefile

import (
	"fmt"
	"math/rand"
)

const (
	epsz = 4
)

type axdesc struct {
	v0 float64
	vz float64
}

func (s *Shapefile) StochEps(n int) {
	fmt.Println("StochEps:")
	xd := mkAxdesc(s.box.xmin, s.box.xmax)
	yd := mkAxdesc(s.box.ymin, s.box.ymax)
	r := rand.New(rand.NewSource(666))
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
}

func mkAxdesc(min, max float64) *axdesc {
	return &axdesc{v0: min, vz: max - min}
}

func (a *axdesc) choose(r *rand.Rand) float64 {
	return a.v0 + r.Float64()*a.vz
}
