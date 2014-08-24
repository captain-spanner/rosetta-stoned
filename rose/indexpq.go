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

	ixget	= iota
	ixput
)

type indexer interface {
	get(s string) []byte
	put(s string, v []byte)
}

type imap struct {
	cmap	map[string][]byte
}

type icache struct {
	cmap	map[string]*entry
	queue	*indexpq
	count	int
	ixq	chan *ixreq
}

type ixreq struct {
	cmd	int
	key	string
	value	[]byte
	res	chan []byte
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

func make_icache() *icache {
	c := new(icache)
	c.cmap = make(map[string]*entry)
	c.queue = &indexpq{}
	heap.Init(c.queue)
	c.count = 0
	c.ixq = make(chan *ixreq)
	go c.ixsrv()
	return c
}

func (c *icache) get(s string) []byte {
	req := new(ixreq)
	req.cmd = ixget
	req.key = s
	req.res = make(chan []byte)
	c.ixq <- req
	return <- req.res
}

func (c *icache) put(s string, v []byte) {
	req := new(ixreq)
	req.cmd = ixput
	req.key = s
	req.value = v
	c.ixq <- req
}

func (c *icache) ixsrv() {
	for {
		req := <- c.ixq
		if req.cmd == ixget {
			req.res <- c.getx(req.key)
		} else {
			c.putx(req.key, req.value)
		}
	}
}

func (c *icache) getx(s string) []byte {
	r, ok := c.cmap[s]
	if !ok {
		return nil
	}
	c.queue.update(r)
	return r.value
}

func (c *icache) putx(s string, v []byte) {
	_, ok := c.cmap[s]
	if ok {
		return
	}
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
