package shapefile

import (
	"fmt"
	"io"
)

type Index struct {
	path    string
	size    int
	body    []byte
	nrecs	int
	recs	[]desc
	hdr	*Header
	err	string
}

type desc struct {
	off	int
	size	int
}

func MakeIndex(n string, out io.Writer) (*Index, error) {
	s := new(Index)
	s.path = n
	body, err := ReadFile(n)
	if err != nil {
		return nil, err
	}
	s.body = body
	s.size = len(body)
	err = lencheck(Hdrsize, s.size, "shx: header")
	if err != nil {
		return nil, err
	}
	z := (s.size - Hdrsize) / 8
	s.nrecs = z
	if out != nil {
		fmt.Fprintln(out, "index:")
		fmt.Fprintf(out, "nrecs\t%d\n", s.nrecs)
	}
	s.hdr = MakeHeader(body, out)
	v := make([]desc, z, z)
	o := Hdrsize
	for i := 0; i < z; i++ {
		v[i].off = int(sb32(body[o:]))
		o += 4
		v[i].size = int(sb32(body[o:]))
		o += 4
	}
	s.recs = v
	return s, nil
}

func (s *Index) Error() string {
	return "shx: "+s.err
}

func (s *Index) Nrecs() int {
	return s.nrecs
}
