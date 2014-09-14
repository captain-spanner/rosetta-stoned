package shapefile

import (
	"fmt"
	"io"
)

type polygons struct {
	bounds bbox
	count int
	polys []*polygon
}

type polygon struct {
	count  int
	points []point
}

type point struct {
	x float64
	y float64
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
	}
	s.decpolys(out)
	return nil
}

func (s *Shapefile) decpolys(out io.Writer) {
	n:= s.shx.nrecs
	s.nrecs = n
	v := make([]*polygons, n, n)
	for i:= 0; i < n; i++ {
		b := s.getrec(i)
		rn := int(sb32(b[0:]))
		data := b[8:]
		v[rn-1] = makepolys(data)
	}
	s.polys = v
}

func (s *Shapefile) getrec(n int) []byte {
	if n < 0 || n >= s.nrecs {
		return nil
	}
	d := s.shx.recs[n]
	o := d.off
	x := o + 8 + d.size
	return s.shp.body[o:x]
}

func makepolys(b []byte) *polygons {
	p := new(polygons)
	makebbox(b[0:], &p.bounds)
	n := int(sl32(b[32:]))
	p.count = n
	zo := 40
	po := 40 + 4 * n
	v := make([]*polygon, n, n)
	for i := 0; i < n; i++ {
		c := int(sl32(b[zo:]))
		zo += 4
		g := new(polygon)
		g.count = c
		ps := make([]point, c, c)
		for j := 0; j < c; j++ {
			makepoint(b[po:], &ps[j])
			po += 16
		}
		g.points = ps
		v[i] = g
	}
	p.polys = v
	return p
}

func makepoint(b []byte, p *point) {
	p.x = fl64(b[0:])
	p.y = fl64(b[8:])
}
