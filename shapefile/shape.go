package shapefile

import (
	"errors"
	"fmt"
	"io"
)

const (
	spoly = 5
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
	if out != nil {
		fmt.Fprintln(out, "shapes:")
	}
	s.body = body
	s.size = len(body)
	err = lencheck(Hdrsize, s.size, "header")
	if err != nil {
		return nil, err
	}
	s.hdr = MakeHeader(body, out)
	if s.hdr.shape != spoly {
		m := fmt.Sprintf("shape type %d not supported", s.hdr.shape)
		return nil, errors.New(m)
	}
	return s, nil
}

func (s *Shapes) Error() string {
	return "shp: "+s.err
}
