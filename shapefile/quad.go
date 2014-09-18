package shapefile

type quad struct {
	box  bbox
	qbox []*bbox
	down []*quad
	only *subreg
	full []*Region
}

type subreg struct {
	box bbox
	reg *Region
}

func MakeQuad(b *bbox) *quad {
	q := new(quad)
	q.box = *b
	return q
}

func (q *quad) AddRegion(r *Region) {
	s := new(subreg)
	s.box = r.poly.bounds
	s.reg = r
	q.addsubreg(s)
}

func (q *quad) addsubreg(s *subreg) {
}

func (q *quad) populate() {
	q.qbox = q.box.divide()
	q.down = make([]*quad, 4, 4)
	for i := 0; i < 4; i++ {
		q.down[i] = MakeQuad(q.qbox[i])
	}
}
