package shapefile

import (
	"fmt"
)

type axdesc struct {
	v0 float64
	vz float64
}

func (s *Shapefile) StochEps(n int) {
	fmt.Println("StochEps:")
}

func mkAxdesc(min, max float64) *axdesc {
	return &axdesc{v0: min, vz: max - min}
}
