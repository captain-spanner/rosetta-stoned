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
	imap	[]byte
}

func print_index() {
	if message {
		for k, _ := range indexm {
			fmt.Println(k)
		}
	}
}

func make_index(s string) (int, string) {
	ix := new(index)
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
		b, err := readpbytes(p, "map", ix.count)
		if err != "" {
			return 1, err
		} else {
			ix.imap = b
		}
	}
	indexm[s] = ix
	indexv = append(indexv, ix)
	return 0, ""
}
