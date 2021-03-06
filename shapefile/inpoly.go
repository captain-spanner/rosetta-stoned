package shapefile

import (
	"fmt"
	"sort"
)

type deployreq struct {
	poly *polygon
	resp chan bool
}

type indata struct {
	runs []*run
}

func mkseg(p *Point, q *Point) *seg {
	s := new(seg)
	if p.x < q.x {
		s.xmin = p.x
		s.ymin = p.y
		s.xmax = q.x
		s.ymax = q.y
	} else {
		s.xmin = q.x
		s.ymin = q.y
		s.xmax = p.x
		s.ymax = p.y
	}
	if s.xmin != s.xmax {
		s.grad = (s.ymax - s.ymin) / (s.xmax - s.xmin)
	}
	return s
}

func (s *Shapefile) dploysrv() {
	for {
		r := <-s.deployq
		r.resp <- r.poly.deploy()
	}
}

func (s *Shapefile) inside(p *polygon, t *Point) bool {
	if !p.bounds.enclosed(t) {
		if idebug {
			fmt.Println("what, not enclosed?")
		}
		return false
	}
	if p.ind == nil {
		r := &deployreq{poly: p, resp: make(chan bool)}
		s.deployq <- r
		<-r.resp
	}
	return p.inside(t)
}

func (p *polygon) deploy() bool {
	if p.ind == nil {
		p.mkindata()
	}
	return true
}

type seg struct {
	bbox
	grad float64
}

type endpt struct {
	x float64
	s *seg
	l bool
}

type runx struct {
	x float64
	e *endpt
}

type run struct {
	x float64
	e []*endpt
}

type byX []*endpt

func (a byX) Len() int           { return len(a) }
func (a byX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byX) Less(i, j int) bool { return a[i].x < a[j].x }

func (p *polygon) mkindata() {
	segs := p.mksegs()
	endpts := make([]*endpt, 0)
	for _, s := range segs {
		if s.xmin != s.xmax {
			endpts = append(endpts, &endpt{x: s.xmin, s: s, l: true})
			endpts = append(endpts, &endpt{x: s.xmax, s: s, l: false})
		}
	}
	sort.Sort(byX(endpts))
	p.ind = &indata{runs: scan(endpts)}
}

func scan(es []*endpt) []*run {
	// invariants:
	//	r is array of completed runs
	//		finalized at end of loop
	//	x is current x
	//		initialized when i == 0
	//	rx is current list
	x := 0.
	r := make([]*run, 0)
	rx := make([]*runx, 0)
	if rdebug2 {
		rendpts += len(es)
	}
	for i, e := range es {
		if i == 0 {
			x = e.x
		} else if x != e.x {
			r = append(r, mkrun(rx, x))
			rx = cull(rx, e.x)
			x = e.x
		}
		if e.l {
			rx = append(rx, &runx{e: e})
		}
	}
	r = append(r, mkrun(rx, x))
	return r
}

func (s *seg) xintercept(x float64) float64 {
	return s.xmin + (x-s.xmin)*s.grad
}

type byXG []*runx

func (a byXG) Len() int      { return len(a) }
func (a byXG) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a byXG) Less(i, j int) bool {
	if a[i].x == a[j].x {
		return a[i].e.s.grad < a[j].e.s.grad
	} else {
		return a[i].x < a[j].x
	}
}

func mkrun(rx []*runx, x float64) *run {
	for _, t := range rx {
		t.x = t.e.s.xintercept(x)
	}
	sort.Sort(byXG(rx))
	e := make([]*endpt, 0)
	for _, t := range rx {
		e = append(e, t.e)
	}
	if rdebug2 {
		l := len(e)
		if l < rmin {
			rmin = l
		}
		if l > rmax {
			rmax = l
		}
	}
	return &run{x: x, e: e}
}

func cull(rx []*runx, x float64) []*runx {
	r := make([]*runx, 0)
	for _, t := range rx {
		if t.e.s.xmax != x {
			r = append(r, t)
		}
	}
	if rdebug2 {
		rcull += len(rx) - len(r)
	}
	return r
}

func (p *polygon) search(pt *Point) int {
	r := p.ind.runs
	x := pt.x
	return sort.Search(len(r), func(i int) bool { return r[i].x >= x })
}

func (p *polygon) inside(pt *Point) bool {
	i := p.search(pt)
	if i < 0 {
		if idebug {
			fmt.Println("search failed")
		}
		return false
	}
	r := p.ind.runs
	if idebug {
		fmt.Printf("index %d %f ", i, r[i].x)
	}
	if r[i].inside(pt) {
		return true
	}
	if pt.x == r[i].x && i != 0 {
		if idebug {
			fmt.Print("check lower ")
		}
		return r[i-1].inside(pt)
	}
	return false
}

func (r *run) inside(pt *Point) bool {
	if idebug {
		pt.print()
		xs := make([]float64, len(r.e), len(r.e))
		for i, e := range r.e {
			xs[i] = e.s.xintercept(pt.x)
		}
		fmt.Println(xs)
	}
	in := false
	for i, e := range r.e {
		x := e.s.xintercept(pt.x)
		if idebug {
			fmt.Printf("e[%d] = %f\n", i, x)
		}
		if in {
			if idebug {
				fmt.Printf("in %f %f\n", pt.x, x)
			}
			if pt.x <= x {
				if idebug {
					fmt.Println("IN")
				}
				return true
			}
		} else {
			if idebug {
				fmt.Printf("out %f %f\n", pt.x, x)
			}
			if pt.x < x {
				if idebug {
					fmt.Println("OUT")
				}
				return false
			}
		}
		in = !in
	}
	if idebug {
		fmt.Printf("end in %t\n", in)
	}
	return false
}
