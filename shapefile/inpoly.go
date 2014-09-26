package shapefile

import (
	"sort"
)

type deployreq struct {
	poly *polygon
	resp chan bool
}

type indata struct {
	runs	[]*run
}

func mkseg(p *point, q *point) *seg {
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

func (s *Shapefile) inside(p *polygon, t *point) bool {
	if !p.bounds.enclosed(t) {
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
	y float64
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
	for i, e := range es {
		if i == 0 {
			x = e.x
		} else if x != e.x {
			r = append(r, mkrun(rx, x))
			rx = cull(rx, e.x)
			x = e.x
		}
		rx = append(rx, &runx{e: e})
	}
	r = append(r, mkrun(rx, x))
	return r
}

func (s *seg) intercept(x float64) float64 {
	return s.xmin + (x-s.xmin)*s.grad
}

type byYG []*runx

func (a byYG) Len() int      { return len(a) }
func (a byYG) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a byYG) Less(i, j int) bool {
	if a[i].y == a[j].y {
		return a[i].e.s.grad < a[j].e.s.grad
	} else {
		return a[i].y < a[j].y
	}
}

func mkrun(rx []*runx, x float64) *run {
	for _, t := range rx {
		t.y = t.e.s.intercept(x)
	}
	sort.Sort(byYG(rx))
	e := make([]*endpt, 0)
	for _, t := range rx {
		e = append(e, t.e)
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
	return r
}

func (p *polygon) inside(pt *point) bool {
	return false
}
