package rose

import (
	"container/heap"
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

type pdata struct {
}

func (*pdata) Describe() string {
	return "Data"
}

func make_pdata(b []byte) part {
	return new(pdata)
}

type pindex struct {
}

func (*pindex) Describe() string {
	return "Index"
}

func make_pindex(b []byte) part {
	return new(pindex)
}

// should allow other corpi than base
func part_get(p string, s string) (part, string, int) {
	q, f := partt[p]
	if !f {
		return nil, p + " is not a part", 1
	}
	c := base
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
