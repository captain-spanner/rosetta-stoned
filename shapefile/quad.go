package shapefile

type quad struct {
	box  bbox
	qbox []*bbox
	down []*quad
	only *subreg
	full []*region
}

type subreg struct {
	box bbox
	reg *region
}

func MakeQuad(b *bbox) *quad {
	q := new(quad)
	q.box = *b
	return q
}

func (q *quad) populate() {
	q.qbox = q.box.divide()
	q.down = make([]*quad, 4, 4)
	for i := 0; i < 4; i++ {
		q.down[i] = MakeQuad(q.qbox[i])
	}
}
