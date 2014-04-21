package rose

import (
	"fmt"
	"io/ioutil"
)

var (
	indexm	map[string]*index = make(map[string]*index)
	indexv	[]*index = make([]*index , 0, 0)
	indexr	map[string]string = make(map[string]string)
	corpi	[]*corpus = make([]*corpus , 0, 0)

	wordlist string = "«words»"
)

const (
)

type part interface {
	Content() []byte
	Describe() string
	Error() string
	Populate(int) ([]string, int)
	Print()
}

type index struct {
	name	string
	path	string
	count	int
	date	string
	format	string
	hash	hashc
	arg	int
	argx	int
	imap	[]byte
	imapz	int
	cache	indexer
	fetch	fetchf
	words	map[string]bool
	ok	bool
}

func (ix *index) print_index() string {
	w := "-"
	if ix.words != nil {
		w = fmt.Sprintf("%d", len(ix.words))
	}
	return fmt.Sprintf("%8d%10s%5d%6t%7s  %s",
		ix.count, hashes[ix.hash], ix.arg, ix.ok, w, ix.name)
}

func print_header() {
	fmt.Printf("%8s%10s%5s%6s%7s  %s\n",
		"Count", "Encoding", "Arg", "Ready", "Words", "Name")
}

func print_indexes() []string {
	if len(indexm) == 0 {
		m := "No indexes"
		if message {
			fmt.Println(m)
		}
		return none
	}
	if message {
		print_header()
	}
	v := none
	for _, x := range indexm {
		m := x.print_index()
		if message {
			fmt.Println(m)
		}
		v = append(v, m)
	}
	return v
}

func list_ixword(w string) {
	for _, x := range indexm {
		m := x.words
		if m == nil {
			continue
		}
		if m[w] {
			fmt.Println(x.name)
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

func make_index(s string) (int, string) {
	ix := new(index)
	ix.ok = true
	ix.name = s
	p := root + "/" + s
	m := checkdir(p)
	if m != "" {
		if message {
			fmt.Printf("%s: %s\n", s, m)
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
	if ix.ok {
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
	indexm[s] = ix
	indexv = append(indexv, ix)
	indexr[s] = s
	return 0, ""
}

type corpus struct {
	name	string
	base	bool
	parts	[]*index
	pcaches	*pcache
}

func make_corpus(s string, opt string) (int, string) {
	c := new(corpus)
	c.name = s
	c.base = false
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
		e, m := make_index(p)
		if message && m != "" {
			fmt.Printf("%s: %s\n", p, m)
		}
		if m != "" {
			return e, m
		}
	}
	isbase := false
	if opt == "base" {
		isbase = true
	}
	e := 0
	m := ""
	if isbase {
		if base != nil {
			e = 1
			m = s + ": base is set"
			isbase = false
		} else {
			base = c
			c.base = true
		}
	}
	d := make([]*index, int(pMax), int(pMax))
	for _, f := range l {
		p := s + "/" + f
		indexr[f] = p
		c := partm[f]
		if c != pNone {
			d[c] = indexm[p]
			pf := parts[c]
			indexr[s + "." + pf] = p
			if isbase {
				indexr[pf] = p
			}
		}
	}
	c.parts = d
	corpi = append(corpi, c)
	return e, m
}

func make_collection(s string) (int, string) {
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
		e, m := make_index(p)
		if m != "" {
			return e, m
		}
	}
	return 0, ""
}

func (ix *index) get(s string, h uint32) []byte {
	if verbose {
		fmt.Printf("ix.get(%q)\n", s)
	}
	b := ix.cache.get(s)
	if b != nil {
		return b
	}
	if verbose {
		fmt.Println("not in cache")
	}
	if !checkmap(ix.imap, h) {
		if verbose {
			fmt.Println("not in map")
		}
		return nil
	} else {
		if verbose {
			fmt.Println("in map")
		}
	}
	b = fetch_string(ix, s, h)
	if b == nil {
		return nil
	}
	ix.cache.put(s, b)
	return b
}
