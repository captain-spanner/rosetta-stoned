package shapefile

type intwrap struct {
	n int
}

func (q *Quad) Search(pt *point) int {
	r := q.searchFirst(pt)
	if r == nil {
		return -1
	}
	return r.Result()
}

func (q *Quad) SearchEps(pt *point) int {
	r := q.searchEps(pt)
	if r == nil {
		return 0
	}
	return r.Result()
}

func (c *intwrap) Result() int {
	return c.n
}

func (q *Quad) search(pt *point, proc func(q *Quad, pt *point) Qres) Qres {
	r := proc(q, pt)
	if r != nil {
		return r
	}
	if q.down == nil {
		return nil
	}
	for i, b := range q.qbox {
		if b.enclosed(pt) {
			return q.down[i].search(pt, proc)
		}
	}
	return nil
}

func wrap(n int) Qres {
	return Qres(&intwrap{n: n})
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
	return wrap(n)
}

func findfirst(q *Quad, pt *point) Qres {
	if q.only != nil {
		return wrap(q.only.region(pt))
	} else if q.full != nil {
		for _, s := range q.full {
			r := s.region(pt)
			if r >= 0 {
				return wrap(r)
			}
		}
	}
	return nil
}

func (q *Quad) searchFirst(pt *point) Qres {
	return q.search(pt, findfirst)
}

func (q *Quad) searchEps(pt *point) Qres {
	return q.search(pt, findeps)
}
