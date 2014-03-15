package rose

import (
	"container/heap"
	"fmt"
)

var (
	pstamp	uint64
	pmakers	[pMax]pmaker = [pMax]pmaker {
		pDj:	make_pdata,
		pDn:	make_pdata,
		pDr:	make_pdata,
		pDv:	make_pdata,
		pIj:	make_pindex,
		pIn:	make_pindex,
		pIr:	make_pindex,
		pIv:	make_pindex,
	}
)

const (
	plimit	= 512
)

type pmaker func([]byte) part

type pcachev struct {
	cmap	map[string]*pentry
	queue	*pindexpq
	count	int
}

type pcache struct {
	caches	[]*pcachev
}

type pentry struct {
	key	string
	epart	part
	stamp	uint64
	index	int
}

type pindexpq []*pentry

type dptr struct {
	tag	byte
	pos	byte
	index	uint32
	ptr	uint32
}

type pdata struct {
	value	[]byte
	error	string
	lex	int
	pos	byte
	ptroz	byte
	words	[]string
	ptrs	[]dptr
	extra	[]string
}

func make_dptr(v []string) *dptr {
	d := new(dptr)
	d.tag = v[0][0]
	d.index = str_uint(v[1])
	d.pos = v[2][0]
	d.ptr = str_uint(v[3])
	return d
}

func (p *pdata) Content() []byte {
	return p.value
}

func (*pdata) Describe() string {
	return "Data"
}

func (p *pdata) Error() string {
	return p.error
}

func (p *pdata) Print() {
}

func make_pdata(b []byte) part {
	p := new(pdata)
	p.value = b
	v := smash_cmd(string(b))
	l := len(v)
	if l < 4 {
		p.error = "short index"
		return p
	}
	p.lex = str_int(v[0])
	p.pos = v[1][0]
	wc := str_int(v[2])
	x := 3
	if x + 2 * wc + 1 >= l {
		p.error = "index too short for words"
		return p
	}
	w := make([]string, 0, wc)
	for i := 0; i < wc; i++ {
		w = append(w, v[x])
		x += 2
	}
	p.words = w
	pc := str_int(v[x])
	x++
	if x + 4 * pc >= l {
		p.error = "index too short for pointers"
		return p
	}
	ps := make([]*dptr, 0, pc)
	for i := 0; i < pc; i++ {
		if i == 0 {
			p.ptroz = byte(len(v[x+1]))
		}
		ps = append(ps, make_dptr(v[x:x+4]))
		x += 4
	}
	if x < l {
		if v[x][0] != '|' {
			p.error = "| expected in index"
			return p
		}
		x++
		xz := l - x
		xv := make([]string, xz, xz)
		copy(xv, v[x:])
		p.extra = xv
	}
	return p
}

type pindex struct {
	value	[]byte
	error	string
	pos	byte
	pvect	[]byte
	sensez	int
	senses	[]uint32
}

func (p *pindex) Content() []byte {
	return p.value
}

func (*pindex) Describe() string {
	return "Index"
}

func (p *pindex) Error() string {
	return p.error
}

func (p *pindex) Print() {
	fmt.Printf("pos %c\n", p.pos)
	fmt.Printf("rels:\n\t{%s }\n", chars_str(p.pvect))
	fmt.Printf("offz %d\n", p.sensez)
	fmt.Printf("senses:\n\t{%s }\n", uints_str(p.senses))
}

func make_pindex(b []byte) part {
	p := new(pindex)
	p.value = b
	v := smash_cmd(string(b))
	l := len(v)
	if l < 3 + 2 {
		p.error = "short index"
		return p
	}
	p.pos = byte(v[0][0])
	sz := str_int(v[1])
	pz := str_int(v[2])
	tz := 3 + pz + 2 + sz
	if l != tz {
		p.error = fmt.Sprintf("size mismatch - %d %d", l, tz)
		return p
	}
	pv := make([]byte, pz, pz)
	for i := 0; i < pz; i++ {
		pv[i] = byte(v[3+i][0])
	}
	p.pvect = pv
	o := 3 + pz + 2
	sv := make([]uint32, sz, sz)
	for i := 0; i < sz; i++ {
		if i == 0 {
			p.sensez = len(v[o+i])
		}
		sv[i] = str_uint(v[o+i])
	}
	p.senses = sv
	return p
}

// should allow other corpi than base
func part_get(p string, s string) (part, string, int) {
	q, f := partt[p]
	if !f {
		return nil, p + " is not a part", 1
	}
	c := base
	if c == nil {
		return nil, "base not set", 1
	}
	if c.pcaches == nil {
		c.pcaches = make_pcache(c.parts)
	}
	k := c.pcaches.caches[q]
	if k == nil  {
		return nil, "no part maker for " + p, 1
	}
	r := k.get(s)
	if r != nil {
		return r, "", 0
	}
	b := fetch_getx(c.parts[q], s)
	if b == nil  {
		return nil, s + ": not found", 1
	}
	r = pmakers[q](b)
	k.put(s, r)
	return r, "", 0
}

func make_pcache(parts []*index) *pcache {
	c := new(pcache)
	c.caches = make([]*pcachev, pMax, pMax)
	for p := pNone; p < pMax; p++ {
		if parts[p] == nil || pmakers[p] == nil {
			continue
		}
		v := new(pcachev)
		v.cmap = make(map[string]*pentry)
		v.queue = &pindexpq{}
		heap.Init(v.queue)
		v.count = 0
		c.caches[p] = v
	}
	return c
}

func (c *pcache) get(s string, p partc) part {
	return c.caches[p].get(s)
}

func (c *pcache) put(s string, p partc, v part) {
	c.caches[p].put(s, v)
}

func (c *pcachev) get(s string) part {
	r, ok := c.cmap[s]
	if !ok {
		return nil
	}
	c.queue.update(r)
	return r.epart
}

func (c *pcachev) put(s string, p part) {
	e := new(pentry)
	e.key = s
	e.epart = p
	e.stamp = pstamp
	pstamp++
	if c.count < plimit {
		heap.Push(c.queue, e)
		c.count++
	} else {
		// LRU
		o := heap.Pop(c.queue).(*entry)
		delete(c.cmap, o.key)
		heap.Push(c.queue, e)
	}
	c.cmap[s] = e
}

func (pq pindexpq) Len() int { return len(pq) }

func (pq pindexpq) Less(i, j int) bool {
	return pq[i].stamp < pq[j].stamp
}

func (pq pindexpq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *pindexpq) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pentry)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *pindexpq) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *pindexpq) update(item *pentry) {
	heap.Remove(pq, item.index)
	item.stamp = pstamp
	heap.Push(pq, item)
	pstamp++
}
