package rose

import (
	"container/heap"
)

var (
	tstamp	uint64
)

const (
	limit	= 256
	worthy	= 376
)

type indexer interface {
	get(s string) []byte
	put(s string, v []byte)
	white() bool
}

type imap struct {
	cmap	map[string][]byte
}

type icache struct {
	cmap	map[string]*entry
	queue	*indexpq
	count	int
}

type entry struct {
	key	string
	value	[]byte
	stamp	uint64
	index	int
}

type indexpq []*entry

func make_indexer(n int) indexer {
	if n < worthy {
		return make_imap()
	} else {
		return make_icache()
	}
}

func make_imap() *imap {
	c := new(imap)
	c.cmap = make(map[string][]byte)
	return c
}

func (c *imap) get(s string) []byte {
	r, ok := c.cmap[s]
	if !ok {
		return nil
	}
	return r
}

func (c *imap) put(s string, v []byte) {
	c.cmap[s] = v
}

func (c *imap) white() bool {
	return false
}

func make_icache() *icache {
	c := new(icache)
	c.cmap = make(map[string]*entry)
	c.queue = &indexpq{}
	heap.Init(c.queue)
	c.count = 0
	return c
}

func (c *icache) get(s string) []byte {
	r, ok := c.cmap[s]
	if !ok {
		return nil
	}
	c.queue.update(r)
	return r.value
}

func (c *icache) put(s string, v []byte) {
	e := new(entry)
	e.key = s
	e.value = v
	e.stamp = tstamp
	tstamp++
	if c.count < limit {
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

func (c *icache) white() bool {
	return true
}

func (pq indexpq) Len() int { return len(pq) }

func (pq indexpq) Less(i, j int) bool {
	return pq[i].stamp < pq[j].stamp
}

func (pq indexpq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *indexpq) Push(x interface{}) {
	n := len(*pq)
	item := x.(*entry)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *indexpq) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *indexpq) update(item *entry) {
	heap.Remove(pq, item.index)
	item.stamp = tstamp
	heap.Push(pq, item)
	tstamp++
}
