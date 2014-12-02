package rose

import (
	"fmt"
	"fsrec"
	//	"io/ioutil"
)

type corpus struct {
	name    string
	parts   []*index
	pcaches *pcache
	coordm  *fsrec.Fsrec
}

func addcorsrv() {
	for {
		c := <-addcorq
		if corpn[c.name] != nil {
			continue
		}
		corpi = append(corpi, c)
		corpn[c.name] = c
	}
}

func (c *corpus) coord(n string, rose *Petal) string {
	var ix *index
	r, f := indexr[n]
	if f {
		ix, f = indexm[r]
	}
	if !f {
		m := n + ": not an index"
		if verbose {
			fmt.Fprintln(rose.wr, m)
		}
		return m
	}
	if ix.hash != hFsRec {
		return "not a coord file"
	}
	p := root + "/" + ix.name + "/" + ix.file
	fs, e := fsrec.MakeFsrec(p, ix.arg, ix.argx)
	if e != nil {
		return e.Error()
	}
	c.coordm = fs
	return ""
}

func (c *corpus) parsecoord(s string) (x float64, y float64, p bool) {
	z := c.coordm.KeySize()
	x, y, p = 0., 0., true
	if s[z+1] == 'U' {
		return
	}
	r, err := fmt.Sscanf(s[z:], " %f %f", &y, &x)
	if r != 2 || err != nil {
		return
	}
	p = false
	return
}
