package shapefile

import (
	"fmt"
	"io"
	"os"
)

const (
	hdrsize = 68
	fdescsz	= 48
)

type Dbase struct {
	path	string
	size	int
	body	[]byte
	tag	byte
	nrecs	int
	hdrsize	int
	recsize	int
	nfdescs	int
	fdraw	[][]byte
	err	string
}

func MakeDbase(n string, out io.Writer) (*Dbase, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	d := new(Dbase)
	d.path = n
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
	if out != nil {
		fmt.Fprintf(out, "path\t%q\n", d.path)
		fmt.Fprintf(out, "size\t%d\n", d.size)
		fmt.Fprintf(out, "tag\t0x%02X\n", d.tag)
		fmt.Fprintf(out, "nrecs\t%d\n", d.nrecs)
		fmt.Fprintf(out, "hdrsize\t%d\n", d.hdrsize)
		fmt.Fprintf(out, "recsize\t%d\n", d.recsize)
	}
	err = d.lencheck(hdrsize + d.hdrsize, "field desc")
	if err != nil {
		return err
	}
	n := d.hdrsize
	if n % fdescsz != 1 {
		d.err = fmt.Sprintf("hdrsize %d not %d * n + 1", n, fdescsz)
		return d
	}
	n = (n - 1) / fdescsz
	d.nfdescs = n
	if out != nil {
		fmt.Fprintf(out, "nfdescs\t%d\n", n)
	}
	v := make([][]byte, n, n)
	d.fdraw = v
	o := hdrsize
	for i := 0; i < n; i++ {
		v[i] = body[o : o + fdescsz]
		o += fdescsz
	}
	return nil
}
