package shapefile

type deployreq struct {
	poly *polygon
	resp chan bool
}

type seg bbox

type endpt struct {
	x float64
	s *seg
	l bool
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

func (p *polygon) mkindata() {
	segs := p.mksegs()
	c := 2 * len(segs)
	endpts := make([]*endpt, c, c)
	for i, s := range segs {
		x := 2 * i
		endpts[x] = &endpt{x: s.xmin, s: s, l: true}
		endpts[x+1] = &endpt{x: s.xmax, s: s, l: false}
	}
}

func (p *polygon) inside(pt *point) bool {
	return false
}
