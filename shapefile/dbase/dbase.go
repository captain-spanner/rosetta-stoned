package shapefile

import (
	"io"
	"os"
)

type Dbase struct {
	path	string
	size	int
	body	[]byte
}

func MakeDbase(n string) (*Dbase, string) {
	f, err := os.Open(n)
	if err != nil {
		return nil, "open failed"
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, "stat failed"
	}
	m := fi.Mode()
	if !m.IsRegular() {
		return nil, "not a file"
	}
	d := new(Dbase)
	d.path = n
	d.size = int(fi.Size())
	d.body = make([]byte, d.size, d.size)
	z, err := io.ReadFull(f, d.body)
	if err != nil {
		return nil, "read error"
	}
	if z != d.size {
		return nil, "short read"
	}
	return d, d.decode()
}

func (d *Dbase) decode() string {
	return ""
}
