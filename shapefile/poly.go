package shapefile

import (
	"fmt"
	"io"
)

type polygons struct {
	count	int
	polys	[]polygon
}

type polygon struct {
	count	int
	points	[]point
}

type point struct {
	x	float64
	y	float64
}

func (s *Shapefile) decode(out io.Writer) error {
	xmin := 1 << 30
	xmax := -1
	n := s.shx.nrecs
	v := s.shx.recs
	for i := 0; i < n; i++ {
		off := v[i].off
		if off < xmin {
			xmin = off
		}
		ext := off + v[i].size + 8
		if ext > xmax {
			xmax = ext
		}
	}
	if out != nil {
		fmt.Fprintf(out, "min\t%d\n", xmin)
		fmt.Fprintf(out, "max\t%d\n", xmax)
		fmt.Fprintf(out, "size\t%d\n", s.shp.size)
		fmt.Fprintf(out, "off0\t%d\n", v[0].off)
		fmt.Fprintf(out, "size0\t%d\n", v[0].size)
	}
	return nil
}
