package rose

import (
	"container/heap"
	"fmt"
)

var (
	pcacheq	chan *pcreq
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
	plimit	= 1024 * 1024

	pcget	= iota
	pcput
)

type pcreq struct {
	corpus	*corpus
	resp	chan(*pcache)
}

type pkreq struct {
	cmd	int
	key	string
	value	part
	resp	chan(part)
}

type pmaker func(partc, []byte, *Petal) part

type pcachev struct {
	cmap	map[string]*pentry
	queue	*pindexpq
	count	int
	pkq	chan *pkreq
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
	tag	psd
	pos	partd
	index	uint32
	ptr	uint32
}

type dframe struct {
	kind	byte
	xword	int
	nwords	int
}

type pdata struct {
	value	[]byte
	error	string
	lex	int
	pos	partd
	ptroz	byte
	words	[]string
	ptrs	[]*dptr
	ptrp	[]part
	rec	int
	frames	[]*dframe
	extra	[]string
}

func make_dptr(v []string, m map[string]psd, rose *Petal) *dptr {
	d := new(dptr)
	t := m[v[0]]
	if t == dNone && verbose {
		fmt.Fprintf(rose.wr, "unknown psd: %#q\n", v[0]) 
	}
	d.tag = t
	s := v[2][0]
	ps, f := posm[s]
	if !f {
		fmt.Fprintf(rose.wr, "bad pos '%c'", s)
		return nil
	}
	d.pos = ps
	d.index = str_uint(v[1])
	d.ptr = str_uint(v[3])
	return d
}

func (d *dptr) pop_dptr(r int, z int, rose *Petal) (part, string, int) {
	dpart := posdx[d.pos]
	v, m, e := part_get(dpart, uint_strz(d.index, z), rose)
	if v != nil {
		v.Populate(r, rose)
	}
	return v, m, e
}

func make_dframe(v []string, m map[string]psd) *dframe {
	d := new(dframe)
	d.kind = v[0][0]
	d.xword = str_intx(v[1])
	d.nwords = str_intx(v[2])
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

func (p *pdata) Populate(r int, rose *Petal) ([]string, int) {
	if r <= p.rec {
		return none, 0
	}
	p.rec = r
	r--
	errs := 0
	mesgs := none
	if p.ptrp == nil {
		z := len(p.ptrs)
		v := make([]part, z, z)
		p.ptrp = v
		oz := int(p.ptroz)
		for i, d := range p.ptrs {
			t, m, e := d.pop_dptr(r, oz, rose)
			if e != 0 {
				errs++
				mesgs = append(mesgs, m)
			} else {
				v[i] = t
			}
		}
	} else {
		for _, t := range p.ptrp {
			if t != nil {
				t.Populate(r, rose)
			}
		}
	}
	return mesgs, errs
}

func (p *pdata) Print(rose *Petal) {
	fmt.Fprintf(rose.wr, "lex %d\n", p.lex)
	fmt.Fprintf(rose.wr, "pos %s\n", poss[p.pos])
	fmt.Fprintf(rose.wr, "words %d\n", len(p.words))
	for _, w := range p.words {
		fmt.Fprintf(rose.wr, "\t%q\n", w)
	}
	for _, d := range p.ptrs {
		fmt.Fprintf(rose.wr, "\t{ %s }\n", dptr_str(d))
	}
	if p.ptrp != nil {
		fmt.Fprintln(rose.wr, "Populated")
	}
}

func make_pdata(c partc, b []byte, rose *Petal) part {
	p := new(pdata)
	p.value = b
	v := smash_cmd(string(b))
	l := len(v)
	if l < 4 {
		p.error = "short index"
		return p
	}
	p.lex = str_int(v[0])
	s := v[1][0]
	ps, f := posm[s]
	if !f {
		p.error = fmt.Sprintf("bad pos '%c'", s)
		return p
	}
	p.pos = ps
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
	ptrs := make([]*dptr, pc, pc)
	psm := posmv[p.pos]
	for i := 0; i < pc; i++ {
		if i == 0 {
			p.ptroz = byte(len(v[x+1]))
		}
		ptrs[i] = make_dptr(v[x:x+4], psm, rose)
		x += 4
	}
	p.ptrs = ptrs
	if x == l {
		return p
	}
	if v[x][0] != '|' {
		m := psdmv[c]
		fc := str_int(v[x])
		x++
		if x + 3*fc >= l {
			p.error = "index too short for frames"
		}
		fs := make([]*dframe, fc, fc)
		for i := 0; i < fc; i++ {
			fs[i] = make_dframe(v[x:x+3], m)
			x += 3
		}
		p.frames = fs
	}
	if x < l {
		if v[x][0] != '|' {
			p.error = "| expected"
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
	pos	partd
	pvect	[]psd
	sensez	int
	senses	[]uint32
	sensep	[]part
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

func (p *pindex) Populate(r int, rose *Petal) ([]string, int) {
	if r == 0 {
		return none, 0
	}
	errs := 0
	mesgs := none
	if p.sensep == nil {
		dpart := posdx[p.pos]
		z := len(p.senses)
		v := make([]part, z, z)
		p.sensep = v
		for i, s := range p.senses {
			t, m, e := part_get(dpart, uint_strz(s, p.sensez), rose)
			if e != 0 {
				errs++
				mesgs = append(mesgs, m)
			} else {
				v[i] = t
				t.Populate(r, rose)
			}
		}
		if errs != 0 && message {
			fmt.Fprintf(rose.wr, "errors: %d\n", errs)
		}
	} else {
		for _, t := range p.sensep {
			if t != nil {
				t.Populate(r, rose)
			}
		}
			
	}
	return mesgs, errs
}

func (p *pindex) Print(rose *Petal) {
	fmt.Fprintf(rose.wr, "pos %s\n", poss[p.pos])
	fmt.Fprintf(rose.wr, "rels:\n\t{%s }\n", psds_str(p.pvect))
	fmt.Fprintf(rose.wr, "senses:\n\t{%s }\n", uints_strz(p.senses, p.sensez))
	if p.sensep != nil {
		fmt.Fprintln(rose.wr, "Populated")
	}
}

func make_pindex(c partc, b []byte, rose *Petal) part {
	p := new(pindex)
	p.value = b
	v := smash_cmd(string(b))
	l := len(v)
	if l < 3 + 2 {
		p.error = "short index"
		return p
	}
	s := byte(v[0][0])
	ps, f := posm[s]
	if !f {
		p.error = fmt.Sprintf("bad pos '%c'", s)
		return p
	}
	p.pos = ps
	sz := str_int(v[1])
	pz := str_int(v[2])
	tz := 3 + pz + 2 + sz
	if l != tz {
		p.error = fmt.Sprintf("size mismatch - %d %d", l, tz)
		return p
	}
	m := psdmv[c]
	pv := make([]psd, pz, pz)
	for i := 0; i < pz; i++ {
		t := m[v[3+i]]
		if t == dNone && verbose {
			fmt.Fprintf(rose.wr, "unknown psd: %#q\n", v[3+i]) 
		}
		pv[i] = t
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

func part_get(p string, s string, rose *Petal) (part, string, int) {
	q, f := partt[p]
	if !f {
		return nil, p + " is not a part", 1
	}
	c := rose.base
	if c == nil {
		return nil, "base not set", 1
	}
	make_pcache(c)
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
	r = pmakers[q](q, b, rose)
	k.put(s, r)
	return r, "", 0
}

func pcachesrv() {
	for {
		req := <- pcacheq
		if req.corpus.pcaches != nil {
			req.resp <- req.corpus.pcaches
		} else {
			req.resp <- make_pcachex(req.corpus.parts)
		}
	}
}

func make_pcache(c *corpus) {
	if c.pcaches == nil {
		req := new(pcreq)
		req.corpus = c
		req.resp = make(chan *pcache)
		pcacheq <- req
		p := <- req.resp
		c.pcaches = p
	}
}

func make_pcachex(parts []*index) *pcache {
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
		v.pkq = make(chan *pkreq)
		go v.pcsrv()
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

func (c *pcachev) pcsrv() {
	for {
		req := <- c.pkq
		if req.cmd == pcget {
			req.resp <- c.getx(req.key)
		} else {
			c.putx(req.key, req.value)
		}
	}
}

func (c *pcachev) get(s string) part {
	req := new(pkreq)
	req.cmd = pcget
	req.key = s
	req.resp = make(chan part)
	c.pkq <- req
	return <- req.resp
}

func (c *pcachev) put(s string, p part) {
	req := new(pkreq)
	req.cmd = pcget
	req.key = s
	req.value = p
	c.pkq <- req
}

func (c *pcachev) getx(s string) part {
	r, ok := c.cmap[s]
	if !ok {
		return nil
	}
	c.queue.update(r)
	return r.epart
}

func (c *pcachev) putx(s string, p part) {
	if c.cmap[s] != nil {
		return
	}
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
		o := heap.Pop(c.queue).(*pentry)
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
