package shapefile

import (
	"io"
)

type Shp struct {
	path    string
	size    int
	body    []byte
	hdr	*Header
	err	string
}

func MakeShp(n string, out io.Writer) (*Shp, error) {
	s := new(Shp)
	s.path = n
	body, err := ReadFile(n)
	if err != nil {
		return nil, err
	}
	s.body = body
	s.size = len(body)
	err = lencheck(Hdrsize, s.size, "header")
	if err != nil {
		return nil, err
	}
	s.hdr = MakeHeader(body, out)
	return s, nil
}

func (s *Shp) Error() string {
	return "shp: "+s.err
}
