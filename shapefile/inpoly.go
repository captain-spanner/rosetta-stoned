package shapefile

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
		r := new(deployreq)
		r.poly = p
		r.resp = make(chan bool)
		s.deployq <- r
		<-r.resp
	}
	r := new(inreq)
	r.pt = t
	r.resp = make(chan bool)
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
}

func (p *polygon) inside(pt *point) bool {
	return false
}
