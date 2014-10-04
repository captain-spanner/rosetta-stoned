package shapefile

type count struct {
	n int
}

func (q *Quad) Search(pt *point, proc func(q *Quad, pt *point) Qres) Qres {
	r := proc(q, pt)
	if r != nil {
		return r
	}
	if q.down == nil {
		return nil
	}
	for i, b := range q.qbox {
		if b.enclosed(pt) {
			return q.down[i].Search(pt, proc)
		}
	}
	return nil
}

func findeps(q *Quad, pt *point) Qres {
	n := 0
	if q.only != nil {
		n = 1
	} else if q.full != nil {
		n = len(q.full)
	}
	if n == 0 {
		return nil
	}
	return Qres(&count{n: n})
}

func (q *Quad) SearchEps(pt *point) Qres {
	return q.Search(pt, findeps)
}
