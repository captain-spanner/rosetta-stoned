package shapefile

import (
	"fmt"
	"io"
	"os"
)

const (
	hdrsize = 68
)

type Dbase struct {
	path	string
	size	int
	body	[]byte
	tag	byte
	nrecs	int
	hdrsize	int
	recsize	int
	err	string
}

func MakeDbase(n string, out io.Writer) (*Dbase, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	d := new(Dbase)
	d.path = n
	if err != nil {
		return nil, err
	}
	m := fi.Mode()
	if !m.IsRegular() {
		d.err = "not a file"
		return nil, d
	}
	d.size = int(fi.Size())
	d.body = make([]byte, d.size, d.size)
	z, err := io.ReadFull(f, d.body)
	if err != nil {
		return nil, err
	}
	if z != d.size {
		d.err = fmt.Sprintf("read mismatch: size %d, %d read", d.size, z)
		return nil, d
	}
	return d, d.decode(out)
}

func (d *Dbase) Error() string {
	return d.err
}

func (d *Dbase) lencheck(n int, s string) error {
	if n > d.size {
		d.err = fmt.Sprintf("need %d for %s, have %d", n, s, d.size)
		return d
	}
	return nil
}

func (d *Dbase) decode(out io.Writer) error {
	err := d.lencheck(hdrsize, "header")
	if err != nil {
		return err
	}
	body := d.body
	d.tag = body[0]
	d.nrecs = int(sb32(body[4:]))
	d.hdrsize = int(sb16(body[8:]))
	d.recsize = int(sb16(body[10:]))
	return nil
}
