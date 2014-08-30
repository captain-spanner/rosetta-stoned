package rose

import (
	"fmt"
)

var (
	fetchtv	[]fetchf = []fetchf {
		hHashed:	fetch_hashed,
		hIndexed:	fetch_indexed,
		hLiteral:	fetch_literal,
	}
)

type fetchf func(*index, string, uint32) []byte

func fetch_file(p string, s string) []byte {
	b, _ := readpbytes(p, s)
	return b
}

func fetch_info(x string, s string, rose *Petal) (*index, uint32, string) {
	var ix *index
	r, f := indexr[x]
	if f {
		ix, f = indexm[r]
	}
	if !f {
		m := x + ": not an index"
		if verbose {
			fmt.Fprintln(rose.wr, m)
		}
		return nil, 0, m
	}
	h := hashs(s)
	return ix, h, ""
}

func fetch_get(x string, s string, rose *Petal) ([]byte, []string, int) {
	ix, h, m := fetch_info(x, s, rose)
	if m != "" {
		return nil, strv(m), 1
	}
	b := ix.get(s, h)
	return b, none, 0
}

func fetch_getx(ix *index, s string) []byte {
	h := hashs(s)
	return ix.get(s, h)
}

func fetch_raw(x string, s string, rose *Petal) ([]byte, []string, int) {
	ix, h, m := fetch_info(x, s, rose)
	if m != "" {
		return nil, strv(m), 1
	}
	b := fetch_string(ix, s, h)
	return b, none, 0
}

func fetch_string(ix *index, s string, h uint32) []byte {
	return fetchtv[ix.hash](ix, s, h)
}

func fetch_literal(ix *index, s string, h uint32) []byte {
	return fetch_file(ix.path, s)
}

func fetch_indexed(ix *index, s string, h uint32) []byte {
	// todo: check that s is an index
	if len(s) < 8 {
		return nil
	}
	p := ""
	if ix.argx == 8 {
		p = ix.path + "/" + s[:ix.arg]
	} else if ix.argx == 9 {
		p = ix.path + "/" + s[3:ix.arg+3]
	} else {
		return nil
	}
	return fetch_file(p, s)
}

func fetch_hashed(ix *index, s string, h uint32) []byte {
	x := fmt.Sprintf("%02X", h % uint32(ix.arg))
	p := ix.path + "/" + x
	return fetch_file(p, s)
}
