package shapefile

import (
	"fmt"
	"io"
)

const (
	rdebug  = false
	rdebug2 = false
	rpop	= true
)

var (
	rendpts int
	rcull   int
	rmin    int = 0
	rmax    int = 16
)

type polygons struct {
	bounds bbox
	holes  int
	count  int
	polys  []*polygon
}

type polygon struct {
	bounds bbox
	cw     bool
	count  int
	points []point
	ind    *indata
}

type Region struct {
	sf   *Shapefile
	poly *polygon
	i    int
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
	n := s.shx.nrecs
	s.nrecs = n
	v := make([]*polygons, n, n)
	for i := 0; i < n; i++ {
		b := s.getrec(i)
		rn := int(sb32(b[0:]))
		data := b[8:]
		p := makepolys(data)
		v[rn-1] = p
		s.holes += p.holes
		if mkplot {
			s.mkplotfile(i, p)
		}
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
	makebbox(b[4:], &p.bounds)
	n := int(sl32(b[36:]))
	np := int(sl32(b[40:]))
	p.count = n
	zo := 44
	po := 44 + 4*n
	v := make([]*polygon, n, n)
	o := make([]int, n+1, n+1)
	for i := 0; i < n; i++ {
		o[i] = int(sl32(b[zo:]))
		zo += 4
	}
	o[n] = np
	for i := 0; i < n; i++ {
		c := o[i+1] - o[i]
		g := new(polygon)
		g.count = c
		ps := make([]point, c, c)
		for j := 0; j < c; j++ {
			makepoint(b[po:], &ps[j])
			po += 16
		}
		g.points = ps
		g.calc()
		v[i] = g
		if !g.cw {
			p.holes++
		}
	}
	p.polys = v
	return p
}

func makepoint(b []byte, p *point) {
	p.x = fl64(b[0:])
	p.y = fl64(b[8:])
}

func (p *polygon) calc() {
	c := p.count
	ps := p.points
	xmin := 360.
	ymin := 360.
	xmax := -360.
	ymax := -360.
	area := 0.
	for i := 0; i < c; i++ {
		x := ps[i].x
		y := ps[i].y
		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
		if y < ymin {
			ymin = y
		}
		if y > ymax {
			ymax = y
		}
		var xn, yn float64
		if i == c-1 {
			xn = ps[0].x
			yn = ps[0].y
		} else {
			xn = ps[i+1].x
			yn = ps[i+1].y
		}
		area += x*yn - xn*y
	}
	p.bounds.xmin = xmin
	p.bounds.ymin = ymin
	p.bounds.xmax = xmax
	p.bounds.ymax = ymax
	p.cw = area < 0.
}

func (p *polygon) mksegs() []*seg {
	c := p.count
	ps := p.points
	s := make([]*seg, c, c)
	for i := 0; i < c; i++ {
		j := i + 1
		if j == c {
			j = 0
		}
		s[i] = mkseg(&ps[i], &ps[j])
	}
	return s
}

func (s *Shapefile) makeregions(p *polygons, i int) {
	for _, q := range p.polys {
		r := new(Region)
		r.sf = s
		r.poly = q
		r.i = i
		s.regs = append(s.regs, r)
	}
}

func (s *Shapefile) analyze() error {
	for i, p := range s.polys {
		if !p.bounds.normal() {
			s.err = fmt.Sprintf("ps %d not normalized", i)
			return s
		}
		if !p.bounds.inside(&s.box) {
			s.err = fmt.Sprintf("ps %d not contained", i)
			return s
		}
		for j, q := range p.polys {
			if !q.bounds.normal() {
				s.err = fmt.Sprintf("ps (%d %d) not normalized", i, j)
				return s
			}
			if !q.bounds.inside(&p.bounds) {
				s.err = fmt.Sprintf("ps (%d %d) not contained", i, j)
				return s
			}
		}
		s.makeregions(p, i)
		if s.err != "" {
			return s
		}
	}
	q := MakeQuad(&s.box)
	s.quad = q
	for i, r := range s.regs {
		if Qdebug {
			fmt.Printf("%d:\n", i)
		}
		q.AddRegion(r)
	}
	s.deployq = make(chan *deployreq)
	go s.dploysrv()
	if rpop {
		s.populate()
	}
	if rdebug {
		if rdebug2 {
			rstats()
		}
	}
	return nil
}

func (s *Shapefile) populate() {
	for i, p := range s.polys {
		for j, q := range p.polys {
			if rdebug {
				fmt.Printf("pop: (%d, %d)\n", i, j)
			}
			s.inside(q, &point{x: q.bounds.xmin, y: q.bounds.ymin})
		}
	}
}

func (pt *point) print() {
	fmt.Printf("pt(%f %f)\n", pt.x, pt.y)
}

func rstats() {
	fmt.Printf("rendpts: %d\n", rendpts)
	fmt.Printf("rcull: %d\n", rcull)
	fmt.Printf("rmin: %d\n", rmin)
	fmt.Printf("rmax: %d\n", rmax)
}
