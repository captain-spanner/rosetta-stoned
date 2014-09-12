package shapefile

import (
	"io"
)

type Shapes struct {
	path    string
	size    int
	body    []byte
	hdr	*Header
	err	string
}

func MakeShapes(n string, out io.Writer) (*Shapes, error) {
	s := new(Shapes)
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

func (s *Shapes) Error() string {
	return "shp: "+s.err
}
