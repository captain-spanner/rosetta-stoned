package shapefile

import (
	"fmt"
	"os"
)

const (
	idebug = false
)

type intwrap struct {
	n int
}

func (q *Quad) Search(pt *Point) int {
	r := q.searchFirst(pt)
	if r == nil {
		return -1
	}
	return r.Result()
}

func (q *Quad) SearchDebug(pt *Point) int {
	r := q.searchDebug(pt)
	if r == nil {
		return -1
	}
	return r.Result()
}

func (q *Quad) SearchEps(pt *Point) int {
	r := q.searchEps(pt)
	if r == nil {
		return 0
	}
	return r.Result()
}

func (c *intwrap) Result() int {
	return c.n
}

func (q *Quad) search(pt *Point, proc func(q *Quad, pt *Point) Qres) Qres {
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

func finddebug(q *Quad, pt *Point) Qres {
	fmt.Print("find ")
	pt.print()
	qa := q.box.area()
	fmt.Printf("area %f\nbox: ", qa)
	q.box.print(os.Stdout)
	fmt.Printf("in %s\n", q.box.encloseds(pt))
	if q.only != nil {
		a := q.only.box.area()
		fmt.Printf("only %f (%f%%)\nbox: ", a, 100.*a/qa)
		q.only.box.print(os.Stdout)
		fmt.Printf("in %s\n", q.only.box.encloseds(pt))
		x := q.only.region(pt)
		if x < 0 {
			fmt.Println("pirate")
		} else {
			fmt.Printf("Region %d\n", x)
		}
		return wrap(x)
	} else if q.full != nil {
		fmt.Println("full")
		for i, s := range q.full {
			a := s.box.area()
			fmt.Printf("%d: %f (%f%%)\nbox: ", i, a, 100.*a/qa)
			s.box.print(os.Stdout)
			fmt.Printf("in %s\n", s.box.encloseds(pt))
			r := s.region(pt)
			if r < 0 {
				fmt.Println("pirate")
			} else {
				fmt.Printf("Region %d\n", r)
				return wrap(r)
			}
		}
	}
	fmt.Println("down")
	return nil
}

func findeps(q *Quad, pt *Point) Qres {
	n := 0
	if q.only != nil {
		if q.only.box.enclosed(pt) {
			n = 1
		}
	} else if q.full != nil {
		for _, s := range q.full {
			if s.box.enclosed(pt) {
				n++
			}
		}
	}
	if n == 0 {
		return nil
	}
	return wrap(n)
}

func findfirst(q *Quad, pt *Point) Qres {
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

func (q *Quad) searchDebug(pt *Point) Qres {
	return q.search(pt, finddebug)
}

func (q *Quad) searchFirst(pt *Point) Qres {
	return q.search(pt, findfirst)
}

func (q *Quad) searchEps(pt *Point) Qres {
	return q.search(pt, findeps)
}
