package rose

import (
	"fmt"
)

var (
	indexm	map[string]*index = make(map[string]*index)
	indexv	[]*index = make([]*index , 0, 0)
)

const (
)

type index struct {
	name	string
	path	string
	count	int
	date	string
	format	string
	hash	hashc
	arg	int
	imap	[]byte
	ok	bool
}

func (ix *index) print_index() string {
	return fmt.Sprintf("%12s%8d%10s%5d%6t  %s",
			ix.name, ix.count, hashes[ix.hash], ix.arg, ix.ok, ix.path)
}

func print_header() {
	fmt.Printf("%12s%8s%10s%5s%6s  %s\n",
			"Name", "Count", "Encoding", "Arg", "Ok", "Path")
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

func (ix *index) decode_fmt() {
	v := smash_cmd(ix.format)
	ix.hash = hError
	ix.arg = 0
	if len(v) < 2 {
		return;
	}
	switch v[0] {
	case "hash":
		ix.hash = hHashed
	case "literal":
		ix.hash = hLiteral
	case "seek":
		ix.hash = hIndexed
	default:
		return
	}
	ix.arg = str_int(v[1])
}

func make_index(s string) (int, string) {
	ix := new(index)
	ix.ok = true
	ix.name = s
	p := root + "/" + s
	m := checkdir(p)
	if m != "" {
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
	if ix.count != 0 {
		b, err := readpbytes(p, "«map»")
		if err != "" {
			return 1, err
		} else {
			ix.imap = b
		}
	}
	ix.decode_fmt()
	if ix.count == 0 || ix.arg == 0 {
		ix.ok = false
	}
	indexm[s] = ix
	indexv = append(indexv, ix)
	return 0, ""
}
