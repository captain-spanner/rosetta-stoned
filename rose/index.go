package rose

import (
	"fmt"
	"io/ioutil"
)

var (
	indexm  map[string]*index  = make(map[string]*index)
	indexv  []*index           = make([]*index, 0, 0)
	indexr  map[string]string  = make(map[string]string)
	corpi   []*corpus          = make([]*corpus, 0, 0)
	corpn   map[string]*corpus = make(map[string]*corpus)
	addixq  chan *index
	addcorq chan *corpus
)

const (
	wordlist string = "«words»"
)

type part interface {
	Content() []byte
	Describe() string
	Error() string
	Populate(int, *Petal) ([]string, int)
	Print(*Petal)
}

type index struct {
	name   string
	file   string
	path   string
	count  int
	date   string
	format string
	hash   hashc
	arg    int
	argx   int
	imap   []byte
	imapz  int
	cache  indexer
	fetch  fetchf
	words  map[string]bool
	ok     bool
}

func addixsrv() {
	for {
		ix := <-addixq
		s := ix.name
		if indexm[s] != nil {
			continue
		}
		indexm[s] = ix
		indexv = append(indexv, ix)
		indexr[s] = s
	}
}

func (ix *index) print_index() string {
	w := "-"
	if ix.words != nil {
		w = fmt.Sprintf("%d", len(ix.words))
	}
	return fmt.Sprintf("%9d%9s%5d%6t%8s  %s",
		ix.count, hashes[ix.hash], ix.arg, ix.ok, w, ix.name)
}

func (rose *Petal) print_header() {
	fmt.Fprintf(rose.wr, "%9s%9s%5s%6s%8s  %s\n",
		"Count", "Encoding", "Arg", "Ready", "Words", "Name")
}

func (rose *Petal) print_indexes() []string {
	if len(indexm) == 0 {
		m := "No indexes"
		if message {
			fmt.Fprintln(rose.wr, m)
		}
		return none
	}
	if message {
		rose.print_header()
	}
	v := none
	for _, x := range indexm {
		m := x.print_index()
		if message {
			fmt.Fprintln(rose.wr, m)
		}
		v = append(v, m)
	}
	return v
}

func (rose *Petal) list_ixword(w string) {
	for _, x := range indexm {
		m := x.words
		if m == nil {
			continue
		}
		if m[w] {
			fmt.Fprintln(rose.wr, x.name)
		}
	}
}

func (ix *index) decode_fmt() bool {
	v := smash_cmd(ix.format)
	ix.hash = hError
	ix.arg = 0
	if len(v) < 2 {
		return false
	}
	switch v[0] {
	case "fsrec":
		ix.hash = hFsRec
	case "hash":
		ix.hash = hHashed
	case "literal":
		ix.hash = hLiteral
	case "seek":
		ix.hash = hIndexed
	default:
		return false
	}
	ix.arg = str_int(v[1])
	if len(v) > 2 {
		ix.argx = str_int(v[2])
	}
	ix.fetch = fetchtv[ix.hash]
	return true
}

func read_words(p string) map[string]bool {
	v := readwordlist(p)
	if v == nil {
		return nil
	}
	m := make(map[string]bool)
	for _, w := range v {
		m[w] = true
	}
	return m
}

func make_index(s string, f string, rose *Petal) (int, string) {
	if indexm[s] != nil {
		return 0, ""
	}
	ix := new(index)
	ix.ok = true
	ix.name = s
	ix.file = f
	p := root + "/" + s
	m := checkdir(p)
	if m != "" {
		if message {
			fmt.Fprintf(rose.wr, "%s: %s\n", s, m)
		}
		return 1, m
	}
	ix.path = p
	m, err := readpstr(p, "«count»")
	if err == "" {
		ix.count = str_int(m)
	}
	m, err = readpstr(p, "«date»")
	if err != "" {
		ix.date = "Unknown"
	} else {
		ix.date = m
	}
	m, err = readpstr(p, "«format»")
	if err != "" {
		return 1, err
	} else {
		ix.format = m
	}
	b := ix.decode_fmt()
	if !b || ix.count == 0 || ix.arg == 0 {
		ix.ok = false
	}
	if ix.ok && ix.hash != hFsRec {
		b, err := readpbytes(p, "«map»")
		if err != "" {
			ix.ok = false
			return 1, err
		} else {
			ix.imap = b
			ix.imapz = len(b)
		}
		ix.cache = make_indexer(ix.count)
	}
	ix.words = read_words(p)
	addixq <- ix
	return 0, ""
}

func (rose *Petal) set_base(s string) {
	c := corpn[s]
	if c == nil {
		fmt.Fprintf(rose.wr, "%s: corpus not loaded\n", s)
	} else {
		rose.base = c
	}
}

func (rose *Petal) make_corpus(s string) (int, string) {
	if corpn[s] != nil {
		return 0, ""
	}
	c := new(corpus)
	c.name = s
	p := root + "/" + s
	v, err := ioutil.ReadDir(p)
	if err != nil {
		return 1, p + ": readdir failed"
	}
	l := make([]string, 0, 0)
	for _, x := range v {
		l = append(l, x.Name())
	}
	for _, f := range l {
		if f == wordlist {
			continue
		}
		p := s + "/" + f
		e, m := make_index(p, f, rose)
		if message && m != "" {
			fmt.Fprintf(rose.wr, "%s: %s\n", p, m)
		}
		if m != "" {
			return e, m
		}
	}
	e := 0
	m := ""
	d := make([]*index, int(pMax), int(pMax))
	for _, f := range l {
		p := s + "/" + f
		indexr[f] = p
		c := partm[f]
		if c != pNone {
			d[c] = indexm[p]
			pf := parts[c]
			indexr[s+"."+pf] = p
			/*
				if isbase {
					indexr[pf] = p
				}
			*/
		}
	}
	c.parts = d
	addcorq <- c
	if rose.base == nil {
		rose.base = c
	}
	return e, m
}

func (rose *Petal) print_corpi() {
	for k, _ := range corpn {
		fmt.Fprintln(rose.wr, k)
	}
}

func make_collection(s string, rose *Petal) (int, string) {
	p := root + "/" + s
	v, err := ioutil.ReadDir(p)
	if err != nil {
		return 1, p + ": readdir failed"
	}
	l := make([]string, 0, 0)
	for _, x := range v {
		l = append(l, x.Name())
	}
	for _, f := range l {
		p := s + "/" + f
		e, m := make_index(p, f, rose)
		if m != "" {
			return e, m
		}
	}
	return 0, ""
}

func (ix *index) get(s string, h uint32) []byte {
	b := ix.cache.get(s)
	if b != nil {
		return b
	}
	if !checkmap(ix.imap, h) {
		return nil
	}
	b = fetch_string(ix, s, h)
	if b == nil {
		return nil
	}
	ix.cache.put(s, b)
	return b
}
