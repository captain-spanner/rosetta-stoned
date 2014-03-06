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

func make_icache() *icache {
	c := new(icache)
	c.cmap = make(map[string]*entry)
	c.queue = &indexpq{}
	heap.Init(c.queue)
	c.count = 0
	return c
}

func (c *icache) Get(s string) []byte {
	r, ok := c.cmap[s]
	if !ok {
		return nil
	}
	c.queue.update(r)
	return r.value
}

func (c *icache) Put(s string, v []byte) {
	e := new(entry)
	e.key = s
	e.value = v
	e.stamp = tstamp
	tstamp++
	if c.count < limit {
		heap.Push(c.queue, e)
		c.cmap[s] = e
		c.count++
	} else {
		// LRU
		o := heap.Pop(c.queue).(*entry)
		delete(c.cmap, o.key)
		heap.Push(c.queue, e)
		c.cmap[s] = e
	}
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

/*
func main() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue and put the items in it.
	pq := &indexpq{}
	heap.Init(pq)
	for value, priority := range items {
		item := &entry{
			value:    value,
			priority: priority,
		}
		heap.Push(pq, item)
	}

	// Insert a new item and then modify its priority.
	item := &entry{
		value:    "orange",
		priority: 1,
	}
	heap.Push(pq, item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(pq).(*entry)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
*/
