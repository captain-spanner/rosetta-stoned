package shapefile

import (
	"sort"
)

type deployreq struct {
	poly *polygon
	resp chan bool
}

type indata struct {
}

type inreq struct {
	pt   *point
	resp chan bool
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
	if p.inq == nil {
		r := &deployreq{poly: p, resp: make(chan bool)}
		s.deployq <- r
		<-r.resp
	}
	r := &inreq{pt: t, resp: make(chan bool)}
	p.inq <- r
	return <-r.resp
}

func (p *polygon) deploy() bool {
	if p.inq == nil {
		p.inq = make(chan *inreq)
		go p.insrv()
	}
	return true
}

func (p *polygon) insrv() {
	p.mkindata()
	for {
		r := <-p.inq
		r.resp <- p.inside(r.pt)
	}
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
	y []float64
	e []*endpt
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
}

func (p *polygon) inside(pt *point) bool {
	return false
}
