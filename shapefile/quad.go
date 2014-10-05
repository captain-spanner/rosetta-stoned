package shapefile

import (
	"fmt"
	"os"
)

const (
	Qdebug  = false
	Qdebug2 = false
)

type Quad struct {
	box  bbox
	qbox []*bbox
	down []*Quad
	only *subreg
	full []*subreg
}

type Qres interface {
	Result() int
}

type subreg struct {
	depth int
	box   bbox
	reg   *Region
}

const (
	edepth = 13
	eps    = 360.0 / float64(1<<edepth)
	eps2   = eps * eps
)

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

func (s *subreg) mksubreg(b *bbox) *subreg {
	r := new(subreg)
	r.depth = s.depth + 1
	r.box = *b
	r.reg = s.reg
	return r
}

func (q *Quad) addsubreg(s *subreg) {
	if Qdebug {
		fmt.Print("quad: ")
		q.box.print(os.Stdout)
		fmt.Print("sub: ")
		s.box.print(os.Stdout)
	}
	if q.down == nil {
		if q.only == nil {
			if Qdebug {
				fmt.Println("set only")
			}
			q.only = s
			return
		}
		q.populate()
		if Qdebug {
			fmt.Println("proc only")
		}
		q.addsubreg(s)
		q.only = nil
	}
	if q.box.full(&s.box, eps2) {
		if q.full == nil {
			q.full = make([]*subreg, 0)
		}
		q.full = append(q.full, s)
		return
	}
	for i := 0; i < 4; i++ {
		if s.box.inside(q.qbox[i]) {
			s.depth++
			q.down[i].addsubreg(s)
			return
		}
	}
	for i := 0; i < 4; i++ {
		b := q.qbox[i].intersection(&s.box)
		if b != nil {
			q.down[i].addsubreg(s.mksubreg(b))
		}
	}
}

func (s *subreg) region(pt *point) int {
	if s.reg.sf.inside(s.reg.poly, pt) {
		return s.reg.i
	} else {
		return -1
	}
}

func (q *Quad) populate() {
	q.qbox = q.box.divide()
	q.down = make([]*Quad, 4, 4)
	for i := 0; i < 4; i++ {
		q.down[i] = MakeQuad(q.qbox[i])
	}
}
