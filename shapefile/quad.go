package shapefile

type Quad struct {
	box  bbox
	qbox []*bbox
	down []*Quad
	only *subreg
	full []*Region
}

type subreg struct {
	depth int
	box   bbox
	reg   *Region
}

func MakeQuad(b *bbox) *Quad {
	q := new(Quad)
	q.box = *b
	return q
}

func (q *Quad) AddRegion(r *Region) {
	s := new(subreg)
	s.box = r.poly.bounds
	s.reg = r
	q.addsubreg(s)
}

func (q *Quad) addsubreg(s *subreg) {
	if q.down == nil {
		if q.only == nil {
			q.only = s
			return
		}
		q.populate()
		q.addsubreg(s)
		q.only = nil
	}
	if q.box.equal(&s.box) {
		if q.full == nil {
			q.full = make([]*Region, 0)
		}
		q.full = append(q.full, s.reg)
		return
	}
}

func (q *Quad) populate() {
	q.qbox = q.box.divide()
	q.down = make([]*Quad, 4, 4)
	for i := 0; i < 4; i++ {
		q.down[i] = MakeQuad(q.qbox[i])
	}
}
